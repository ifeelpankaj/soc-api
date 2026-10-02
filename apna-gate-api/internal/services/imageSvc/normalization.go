package imagesvc

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"time"

	"github.com/disintegration/imaging"
)

const visitorImageMaxEdge = 1600

// Keep recovery local to untrusted image processing. No storage mutation has
// happened here, and a malformed decoder panic must not terminate the server.
func imageProcessingBoundary(process func() ([]byte, error)) (data []byte, err error) {
	defer func() {
		if recover() != nil {
			data, err = nil, ErrImageInvalid
		}
	}()
	return process()
}

func normalizeVisitorImage(data []byte) ([]byte, error) {
	return imageProcessingBoundary(func() ([]byte, error) {
		cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil || (format != "jpeg" && format != "png") {
			return nil, ErrImageInvalid
		}
		if !dimensionsAllowed(cfg.Width, cfg.Height) {
			return nil, ErrImageTooLarge
		}
		src, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
		if err != nil {
			return nil, ErrImageInvalid
		}
		w, h := src.Bounds().Dx(), src.Bounds().Dy()
		if !dimensionsAllowed(w, h) {
			return nil, ErrImageTooLarge
		}
		if (w != cfg.Width || h != cfg.Height) && (w != cfg.Height || h != cfg.Width) {
			return nil, ErrImageInvalid
		}
		if w > visitorImageMaxEdge || h > visitorImageMaxEdge {
			src = imaging.Fit(src, visitorImageMaxEdge, visitorImageMaxEdge, imaging.Lanczos)
		}
		flat := image.NewRGBA(src.Bounds())
		draw.Draw(flat, flat.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
		draw.Draw(flat, flat.Bounds(), src, src.Bounds().Min, draw.Over)
		var out bytes.Buffer
		if err = jpeg.Encode(&out, flat, &jpeg.Options{Quality: 80}); err != nil {
			return nil, ErrImageInvalid
		}
		return out.Bytes(), nil
	})
}

// prepareVisitorImage keeps the common app-generated JPEG on a cheap,
// byte-preserving path. Anything that cannot be proven safe by the conservative
// marker scan falls back to the existing full decode and normalization path.
func prepareVisitorImage(data []byte) ([]byte, error) {
	return imageProcessingBoundary(func() ([]byte, error) {
		validationStart := time.Now()
		cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil || (format != "jpeg" && format != "png") {
			imageStageDuration.WithLabelValues("validation").Observe(time.Since(validationStart).Seconds())
			return nil, ErrImageInvalid
		}
		if !dimensionsAllowed(cfg.Width, cfg.Height) {
			imageStageDuration.WithLabelValues("validation").Observe(time.Since(validationStart).Seconds())
			return nil, ErrImageTooLarge
		}
		passthrough := format == "jpeg" && cfg.Width <= visitorImageMaxEdge && cfg.Height <= visitorImageMaxEdge && jpegFastPathAllowed(data)
		imageStageDuration.WithLabelValues("validation").Observe(time.Since(validationStart).Seconds())
		if passthrough {
			imageNormalization.WithLabelValues("passthrough").Inc()
			return data, nil
		}
		normalizationStart := time.Now()
		normalized, normalizeErr := normalizeVisitorImage(data)
		imageStageDuration.WithLabelValues("normalization").Observe(time.Since(normalizationStart).Seconds())
		if normalizeErr == nil {
			imageNormalization.WithLabelValues("reencoded").Inc()
		}
		return normalized, normalizeErr
	})
}

// jpegFastPathAllowed validates the JPEG container without expanding its pixel
// buffer. APP metadata and comments deliberately force full normalization.
func jpegFastPathAllowed(data []byte) bool {
	if len(data) < 4 || data[0] != 0xff || data[1] != 0xd8 {
		return false
	}
	pos := 2
	seenFrame, seenScan := false, false
	inScan := false
	for pos < len(data) {
		if inScan {
			if data[pos] != 0xff {
				pos++
				continue
			}
			markerStart := pos
			for pos < len(data) && data[pos] == 0xff {
				pos++
			}
			if pos >= len(data) {
				return false
			}
			marker := data[pos]
			if marker == 0x00 || marker >= 0xd0 && marker <= 0xd7 {
				pos++
				continue
			}
			pos = markerStart
			inScan = false
		}
		if data[pos] != 0xff {
			return false
		}
		for pos < len(data) && data[pos] == 0xff {
			pos++
		}
		if pos >= len(data) {
			return false
		}
		marker := data[pos]
		pos++
		switch marker {
		case 0xd9:
			return seenFrame && seenScan && pos == len(data)
		case 0xd8, 0x00, 0x01, 0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7:
			return false
		}
		if pos+2 > len(data) {
			return false
		}
		length := int(data[pos])<<8 | int(data[pos+1])
		if length < 2 || pos+length > len(data) {
			return false
		}
		payload := data[pos+2 : pos+length]
		switch {
		case marker == 0xe0:
			// Permit only an ordinary JFIF header without an embedded thumbnail.
			if len(payload) != 14 || !bytes.Equal(payload[:5], []byte("JFIF\x00")) || payload[12] != 0 || payload[13] != 0 {
				return false
			}
		case marker >= 0xe1 && marker <= 0xef || marker == 0xfe:
			return false
		case marker == 0xc0 || marker == 0xc1 || marker == 0xc2:
			if seenFrame {
				return false
			}
			seenFrame = true
		case marker == 0xda:
			if !seenFrame || len(payload) == 0 {
				return false
			}
			seenScan = true
			inScan = true
		case marker == 0xc4 || marker == 0xdb || marker == 0xdd:
			// Huffman, quantization, and restart-interval tables are structural.
		default:
			return false
		}
		pos += length
	}
	return false
}
