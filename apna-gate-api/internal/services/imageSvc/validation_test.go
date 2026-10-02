package imagesvc

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-server/internal/models"
)

func testImage(t *testing.T, format string) []byte {
	t.Helper()
	var b bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 4, 3))
	var err error
	if format == "jpeg" {
		err = jpeg.Encode(&b, img, nil)
	} else {
		err = png.Encode(&b, img)
	}
	if err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
func multipartRequest(t *testing.T, data []byte, extra bool) *http.Request {
	t.Helper()
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", "untrusted.exe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatal(err)
	}
	if extra {
		if err = writer.WriteField("folder", "untrusted"); err != nil {
			t.Fatal(err)
		}
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestMultipartValidation(t *testing.T) {
	for _, format := range []string{"png", "jpeg"} {
		req := multipartRequest(t, testImage(t, format), false)
		result, err := readMultipart(context.Background(), httptest.NewRecorder(), req)
		if err != nil {
			t.Fatal(err)
		}
		extension := ".png"
		if format == "jpeg" {
			extension = ".jpg"
		}
		if !strings.HasSuffix(result.filename, extension) || strings.Contains(result.filename, "untrusted") {
			t.Fatal(result.filename)
		}
	}
	for _, tc := range []struct {
		name  string
		data  []byte
		extra bool
		want  error
	}{
		{"empty", nil, false, ErrImageInvalid},
		{"spoof", []byte("not an image"), false, ErrImageInvalid},
		{"truncated", testImage(t, "png")[:40], false, ErrImageInvalid},
		{"extra part", testImage(t, "png"), true, ErrImageInvalid},
		{"oversized chunked", bytes.Repeat([]byte{0}, int(models.ImageMaxUploadBytes+1)), false, ErrImageTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := multipartRequest(t, tc.data, tc.extra)
			req.ContentLength = -1
			if _, err := readMultipart(context.Background(), httptest.NewRecorder(), req); !errors.Is(err, tc.want) {
				t.Fatalf("got %v want %v", err, tc.want)
			}
		})
	}
}

func TestImageDimensionsBeforeDecode(t *testing.T) {
	for _, wh := range [][2]int{{0, 1}, {1, 0}, {-1, 1}, {5001, 5000}, {int(^uint(0) >> 1), 2}} {
		if dimensionsAllowed(wh[0], wh[1]) {
			t.Fatalf("accepted %v", wh)
		}
	}
	if !dimensionsAllowed(5000, 5000) {
		t.Fatal("rejected boundary")
	}
	data := testImage(t, "png")
	binary.BigEndian.PutUint32(data[16:20], 100_000)
	binary.BigEndian.PutUint32(data[20:24], 100_000)
	binary.BigEndian.PutUint32(data[29:33], crc32.ChecksumIEEE(data[12:29]))
	if _, err := validateImage(data); !errors.Is(err, ErrImageTooLarge) {
		t.Fatalf("expected pixel limit before decoding: %v", err)
	}
}

func TestMultipartBodyLimitIncludesEpilogue(t *testing.T) {
	req := multipartRequest(t, testImage(t, "png"), false)
	req.Body = io.NopCloser(io.MultiReader(req.Body, strings.NewReader(strings.Repeat("x", int(models.ImageMaxBodyBytes)))))
	req.ContentLength = -1
	if _, err := readMultipart(context.Background(), httptest.NewRecorder(), req); !errors.Is(err, ErrImageTooLarge) {
		t.Fatalf("got %v", err)
	}
	for _, body := range []string{"", "--broken\r\n"} {
		req = httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=broken")
		if _, err := readMultipart(context.Background(), httptest.NewRecorder(), req); !errors.Is(err, ErrImageInvalid) {
			t.Fatalf("got %v", err)
		}
	}
}
