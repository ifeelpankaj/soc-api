package imagesvc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"testing"
)

func TestVisitorNormalization(t *testing.T) {
	for _, size := range []image.Point{{3200, 1600}, {400, 800}, {3, 4}} {
		var input bytes.Buffer
		if err := png.Encode(&input, image.NewNRGBA(image.Rect(0, 0, size.X, size.Y))); err != nil {
			t.Fatal(err)
		}
		out, err := normalizeVisitorImage(input.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		img, err := jpeg.Decode(bytes.NewReader(out))
		if err != nil {
			t.Fatal(err)
		}
		want := size
		if want.X > 1600 {
			want = image.Pt(1600, 800)
		}
		if img.Bounds().Size() != want {
			t.Fatalf("got %v want %v", img.Bounds(), want)
		}
		r, g, b, _ := img.At(0, 0).RGBA()
		if r < 65000 || g < 65000 || b < 65000 {
			t.Fatal("alpha was not flattened onto white")
		}
	}
}

func TestPreparedVisitorJPEGUsesBytePreservingFastPath(t *testing.T) {
	input := testImage(t, "jpeg")
	out, err := prepareVisitorImage(input)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, input) {
		t.Fatal("compliant JPEG was re-encoded")
	}
}

func TestPreparedVisitorImageFallsBackWhenNormalizationIsRequired(t *testing.T) {
	pngInput := testImage(t, "png")
	pngOutput, err := prepareVisitorImage(pngInput)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(pngOutput, pngInput) {
		t.Fatal("PNG bypassed normalization")
	}
	if _, format, decodeErr := image.DecodeConfig(bytes.NewReader(pngOutput)); decodeErr != nil || format != "jpeg" {
		t.Fatalf("fallback output format=%q err=%v", format, decodeErr)
	}

	jpegInput := testImage(t, "jpeg")
	metadata := []byte{0xff, 0xe1, 0, 8, 'E', 'x', 'i', 'f', 0, 0}
	withMetadata := append(append(append([]byte{}, jpegInput[:2]...), metadata...), jpegInput[2:]...)
	metadataOutput, err := prepareVisitorImage(withMetadata)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(metadataOutput, withMetadata) || bytes.Contains(metadataOutput, []byte("Exif")) {
		t.Fatal("metadata-bearing JPEG bypassed normalization")
	}

	var oversized bytes.Buffer
	if err = jpeg.Encode(&oversized, image.NewRGBA(image.Rect(0, 0, visitorImageMaxEdge+1, 20)), nil); err != nil {
		t.Fatal(err)
	}
	oversizedOutput, err := prepareVisitorImage(oversized.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	oversizedConfig, err := jpeg.DecodeConfig(bytes.NewReader(oversizedOutput))
	if err != nil || oversizedConfig.Width != visitorImageMaxEdge {
		t.Fatalf("oversized fallback width=%d err=%v", oversizedConfig.Width, err)
	}
}

func TestJPEGFastPathRejectsUncertainStructure(t *testing.T) {
	input := testImage(t, "jpeg")
	for _, candidate := range [][]byte{
		input[:len(input)-2],
		append(append([]byte{}, input...), 0),
		{0xff, 0xd8, 0xff, 0xe1, 0, 1, 0xff, 0xd9},
	} {
		if jpegFastPathAllowed(candidate) {
			t.Fatal("uncertain JPEG accepted by fast path")
		}
	}
}

func TestVisitorOrientationAndMetadata(t *testing.T) {
	for orientation := uint16(1); orientation <= 8; orientation++ {
		var input bytes.Buffer
		src := image.NewRGBA(image.Rect(0, 0, 40, 20))
		for y := 0; y < 20; y++ {
			for x := 0; x < 40; x++ {
				src.Set(x, y, color.RGBA{255, 0, 0, 255})
			}
		}
		if err := jpeg.Encode(&input, src, nil); err != nil {
			t.Fatal(err)
		}
		// Minimal EXIF/TIFF IFD containing the orientation tag.
		exif := []byte{'E', 'x', 'i', 'f', 0, 0, 'I', 'I', 42, 0, 8, 0, 0, 0, 1, 0, 0x12, 1, 3, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint16(exif[24:26], orientation)
		segment := []byte{0xff, 0xe1, 0, 0}
		segmentLength := len(exif) + 2
		if segmentLength > math.MaxUint16 {
			t.Fatal("test EXIF segment exceeds JPEG length limit")
		}
		binary.BigEndian.PutUint16(segment[2:], uint16(segmentLength)) //nolint:gosec // Bounds checked above.
		withEXIF := append(append(append([]byte{}, input.Bytes()[:2]...), segment...), exif...)
		withEXIF = append(withEXIF, input.Bytes()[2:]...)
		out, err := normalizeVisitorImage(withEXIF)
		if err != nil {
			t.Fatal(err)
		}
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(out))
		if err != nil {
			t.Fatal(err)
		}
		w, h := 40, 20
		if orientation >= 5 {
			w, h = h, w
		}
		if cfg.Width != w || cfg.Height != h {
			t.Fatalf("orientation %d: %dx%d", orientation, cfg.Width, cfg.Height)
		}
		if bytes.Contains(out, []byte("Exif")) {
			t.Fatal("EXIF preserved")
		}
	}
}

func TestImageProcessingFailureBoundary(t *testing.T) {
	oversized := testImage(t, "png")
	binary.BigEndian.PutUint32(oversized[16:20], 100_000)
	binary.BigEndian.PutUint32(oversized[20:24], 100_000)
	binary.BigEndian.PutUint32(oversized[29:33], crc32.ChecksumIEEE(oversized[12:29]))
	if _, err := normalizeVisitorImage(oversized); !errors.Is(err, ErrImageTooLarge) {
		t.Fatalf("pixel cap must precede decode: %v", err)
	}
	for _, input := range [][]byte{nil, []byte("not an image"), testImage(t, "png")[:40]} {
		if _, err := normalizeVisitorImage(input); !errors.Is(err, ErrImageInvalid) {
			t.Fatalf("got %v", err)
		}
	}
	if data, err := imageProcessingBoundary(func() ([]byte, error) { panic("malformed image") }); data != nil || !errors.Is(err, ErrImageInvalid) {
		t.Fatalf("%v %v", data, err)
	}
}
