package models

import (
	"context"
	"io"
	"time"
)

const (
	ImageMaxUploadBytes    int64 = 5 * 1024 * 1024
	ImageMaxBodyBytes      int64 = ImageMaxUploadBytes + 64*1024
	ImageMaxPixels         int64 = 25_000_000
	ImageSignedURLTTL            = 15 * time.Minute
	ImageProcessingTimeout       = 45 * time.Second
	ImageUploadTimeout           = 20 * time.Second
	ImageDatabaseTimeout         = 5 * time.Second
	ImageCleanupTimeout          = 5 * time.Second
	ImageConcurrency             = 2
)

type StoredImage struct{ FileID, Path, URL string }

func (s StoredImage) Managed() bool { return s.FileID != "" && s.Path != "" }

type ImageUploadInput struct {
	Reader           io.Reader
	FileName, Folder string
}
type ImageStorage interface {
	Upload(context.Context, ImageUploadInput) (*StoredImage, error)
	Delete(context.Context, string) error
	SignedURL(string, time.Time, ...ImageVariant) (string, error)
}

// ImageTarget is built from authenticated identity and validated route parameters,
// never decoded from client JSON. VisitorID is resolved by authorization.
type ImageTarget struct{ ActorID, SocietyID, EntryID, FlatID, VisitorID int64 }

func (t ImageTarget) Avatar() bool { return t.EntryID == 0 }

type ImageView struct {
	URL       string       `json:"url"`
	ExpiresAt time.Time    `json:"expires_at"`
	Variant   ImageVariant `json:"variant"`
}

type ImageVariant string

const (
	ImageOriginal ImageVariant = "original"
	ImageList     ImageVariant = "list"
	ImageAvatar   ImageVariant = "avatar"
	ImageDetail   ImageVariant = "detail"
)

type ImageViewAPIResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    ImageView `json:"data"`
}
