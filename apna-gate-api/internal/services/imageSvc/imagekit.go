package imagesvc

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	imagekit "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"github.com/imagekit-developer/imagekit-go/v2/packages/param"
	"github.com/imagekit-developer/imagekit-go/v2/shared"
	"go-server/internal/config"
	"go-server/internal/models"
)

type ImageKitStorage struct {
	client           imagekit.Client
	endpoint, prefix string
}

func NewImageKitStorage(cfg config.ImageKitConfig) *ImageKitStorage {
	return &ImageKitStorage{
		client:   imagekit.NewClient(option.WithPrivateKey(cfg.PrivateKey), option.WithMaxRetries(0)),
		endpoint: strings.TrimRight(cfg.URLEndpoint, "/"), prefix: "/" + cfg.FolderPrefix + "/",
	}
}
func (s *ImageKitStorage) Upload(ctx context.Context, in models.ImageUploadInput) (*models.StoredImage, error) {
	result, err := s.client.Files.Upload(ctx, imagekit.FileUploadParams{
		File: in.Reader, FileName: in.FileName, Folder: param.NewOpt(in.Folder),
		IsPrivateFile: param.NewOpt(true), UseUniqueFileName: param.NewOpt(false),
		OverwriteFile: param.NewOpt(false), ResponseFields: []string{"isPrivateFile"},
	})
	// Never propagate SDK errors: they can contain request headers and bodies.
	if err != nil {
		return nil, ErrImageProvider
	}
	if result == nil {
		return nil, ErrImageProvider
	}
	stored := &models.StoredImage{FileID: result.FileID, Path: result.FilePath, URL: result.URL}
	if !stored.Managed() || !result.IsPrivateFile || !s.validPath(stored.Path) || stored.Path != "/"+strings.Trim(in.Folder, "/")+"/"+in.FileName {
		return stored, ErrImageProvider
	}
	u, parseErr := url.Parse(stored.URL)
	if parseErr != nil || u.Scheme != "https" || u.Host == "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return stored, ErrImageProvider
	}
	return stored, nil
}
func (s *ImageKitStorage) Delete(ctx context.Context, fileID string) error {
	err := s.client.Files.Delete(ctx, fileID)
	var providerErr *imagekit.Error
	if errors.As(err, &providerErr) && providerErr.StatusCode == 404 {
		return nil
	}
	if err != nil {
		return ErrImageProvider
	}
	return nil
}
func (s *ImageKitStorage) validPath(p string) bool {
	return strings.HasPrefix(p, s.prefix) && path.Clean(p) == p && !strings.ContainsAny(p, "?%#\\\r\n")
}
func (s *ImageKitStorage) SignedURL(p string, expiresAt time.Time, variants ...models.ImageVariant) (string, error) {
	if !s.validPath(p) {
		return "", ErrImageProvider
	}
	variant := models.ImageOriginal
	if len(variants) > 0 {
		var err error
		variant, err = ParseImageVariant(string(variants[0]))
		if err != nil {
			return "", err
		}
	}
	// Build the allowlisted transformation into the path BEFORE signing it.
	switch variant {
	case models.ImageList:
		p = "/tr:w-80,h-80,fo-face" + p
	case models.ImageAvatar:
		p = "/tr:w-300,h-300,fo-face" + p
	case models.ImageDetail:
		p = "/tr:w-600" + p
	}
	// SDK accepts a relative TTL. Verify the emitted absolute expiry to avoid a
	// one-second boundary crossing between our clock read and its internal clock.
	for range 3 {
		seconds := expiresAt.Unix() - time.Now().Unix()
		if seconds <= 0 {
			return "", ErrImageProvider
		}
		signed := s.client.Helper.BuildURL(shared.SrcOptionsParam{
			URLEndpoint: s.endpoint, Src: p, Signed: param.NewOpt(true), ExpiresIn: param.NewOpt(float64(seconds)),
		})
		parsed, err := url.Parse(signed)
		if err == nil && parsed.Query().Get("ik-t") == strconv.FormatInt(expiresAt.Unix(), 10) && parsed.Query().Get("ik-s") != "" {
			return signed, nil
		}
	}
	return "", ErrImageProvider
}
