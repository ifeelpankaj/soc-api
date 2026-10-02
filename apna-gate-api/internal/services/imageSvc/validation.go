package imagesvc

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go-server/internal/models"
)

type validatedImage struct {
	data     []byte
	filename string
}

func dimensionsAllowed(width, height int) bool {
	return width > 0 && height > 0 && int64(width) <= models.ImageMaxPixels/int64(height)
}
func validateImage(data []byte) (extension string, err error) {
	defer func() {
		if recover() != nil {
			extension, err = "", ErrImageInvalid
		}
	}()
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "jpeg" && format != "png") {
		return "", ErrImageInvalid
	}
	if !dimensionsAllowed(cfg.Width, cfg.Height) {
		return "", ErrImageTooLarge
	}
	decoded, decodedFormat, err := image.Decode(bytes.NewReader(data))
	if err != nil || decodedFormat != format || decoded.Bounds().Dx() != cfg.Width || decoded.Bounds().Dy() != cfg.Height {
		return "", ErrImageInvalid
	}
	if format == "jpeg" {
		return ".jpg", nil
	}
	return ".png", nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
func readMultipart(ctx context.Context, writer http.ResponseWriter, request *http.Request, normalize ...bool) (validatedImage, error) {
	request.Body = http.MaxBytesReader(writer, request.Body, models.ImageMaxBodyBytes)
	request.Body = struct {
		io.Reader
		io.Closer
	}{contextReader{ctx, request.Body}, request.Body}
	reader, err := request.MultipartReader()
	if err != nil {
		return validatedImage{}, ErrImageInvalid
	}
	part, err := reader.NextPart()
	if err != nil {
		return validatedImage{}, multipartError(err)
	}
	if part.FormName() != "file" || part.FileName() == "" {
		return validatedImage{}, ErrImageInvalid
	}
	data, err := io.ReadAll(io.LimitReader(part, models.ImageMaxUploadBytes+1))
	if err != nil {
		return validatedImage{}, multipartError(err)
	}
	if int64(len(data)) > models.ImageMaxUploadBytes {
		return validatedImage{}, ErrImageTooLarge
	}
	if len(data) == 0 {
		return validatedImage{}, ErrImageInvalid
	}
	if _, err = reader.NextPart(); !errors.Is(err, io.EOF) {
		return validatedImage{}, multipartError(err)
	}
	// Consume any epilogue too, so the body cap also covers bytes after the boundary.
	if _, err = io.Copy(io.Discard, request.Body); err != nil {
		return validatedImage{}, multipartError(err)
	}
	var extension string
	if len(normalize) > 0 && normalize[0] {
		data, err = prepareVisitorImage(data)
		extension = ".jpg"
	} else {
		validationStart := time.Now()
		extension, err = validateImage(data)
		imageStageDuration.WithLabelValues("validation").Observe(time.Since(validationStart).Seconds())
	}
	if err != nil {
		return validatedImage{}, err
	}
	return validatedImage{data: data, filename: uuid.NewString() + extension}, nil
}
func multipartError(err error) error {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return ErrImageTooLarge
	}
	return ErrImageInvalid
}
