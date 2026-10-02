package imagesvc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	"go-server/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	storage models.ImageStorage
	repo    repository.ImageRepository
	auth    Authorizer
	prefix  string
	now     func() time.Time
	slots   chan struct{}
}

func New(storage models.ImageStorage, repo repository.ImageRepository, auth Authorizer, prefix string, now func() time.Time) *Service {
	if repo == nil || auth == nil {
		panic("imageSvc: required dependency is nil")
	}
	if now == nil {
		now = time.Now
	}
	return &Service{storage: storage, repo: repo, auth: auth, prefix: prefix, now: now, slots: make(chan struct{}, models.ImageConcurrency)}
}

// Execute keeps authentication/authorization ahead of multipart reads. HTTP
// parsing is delegated to validation.go; persistence never receives an HTTP body.
func (s *Service) Execute(ctx context.Context, method string, target models.ImageTarget, writer http.ResponseWriter, request *http.Request) (view *models.ImageView, err error) {
	start := time.Now()
	defer func() {
		imageDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
		if err != nil {
			imageFailures.WithLabelValues(method).Inc()
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, models.ImageProcessingTimeout)
	defer cancel()
	if s.storage == nil {
		return nil, ErrImageUnavailable
	}
	authCtx, authCancel := context.WithTimeout(ctx, models.ImageDatabaseTimeout)
	target, err = s.auth.Authorize(authCtx, target, method == http.MethodPut)
	authCancel()
	if err != nil {
		return nil, safeError(err)
	}
	switch method {
	case http.MethodGet:
		variant, variantErr := ParseImageVariant(request.URL.Query().Get("variant"))
		if variantErr != nil {
			return nil, variantErr
		}
		dbCtx, done := context.WithTimeout(ctx, models.ImageDatabaseTimeout)
		stored, readErr := s.repo.Read(dbCtx, target)
		done()
		if readErr != nil {
			return nil, safeError(readErr)
		}
		return s.view(stored, variant)
	case http.MethodDelete:
		if target.FlatID != 0 {
			return nil, ErrImageForbidden
		}
		return nil, s.replace(ctx, target, models.StoredImage{}, request.Header.Get("X-Request-ID"))
	case http.MethodPut:
		body := request.Body
		stopClosing := context.AfterFunc(ctx, func() { _ = body.Close() })
		defer stopClosing()
		select {
		case s.slots <- struct{}{}:
			defer func() { <-s.slots }()
		default:
			return nil, ErrImageBusy
		}
		file, readErr := readMultipart(ctx, writer, request, !target.Avatar())
		if readErr != nil {
			return nil, readErr
		}
		if ctx.Err() != nil {
			return nil, ErrImageUnavailable
		}
		folder := fmt.Sprintf("%s/users/%d/avatars", s.prefix, target.ActorID)
		if !target.Avatar() {
			folder = fmt.Sprintf("%s/societies/%d/visitors/%d", s.prefix, target.SocietyID, target.VisitorID)
		}
		uploadCtx, done := context.WithTimeout(ctx, models.ImageUploadTimeout)
		uploadStart := time.Now()
		stored, uploadErr := s.storage.Upload(uploadCtx, models.ImageUploadInput{Reader: bytes.NewReader(file.data), FileName: file.filename, Folder: folder})
		imageStageDuration.WithLabelValues("provider_upload").Observe(time.Since(uploadStart).Seconds())
		done()
		if uploadErr != nil || stored == nil || !stored.Managed() || stored.URL == "" {
			s.logCleanup(target, request.Header.Get("X-Request-ID"), "upload_outcome_unconfirmed", models.StoredImage{Path: folder + "/" + file.filename})
			if stored != nil {
				s.cleanup(target, *stored, request.Header.Get("X-Request-ID"))
			}
			return nil, ErrImageProvider
		}
		replaceStart := time.Now()
		err = s.replace(ctx, target, *stored, request.Header.Get("X-Request-ID"))
		imageStageDuration.WithLabelValues("database_replace").Observe(time.Since(replaceStart).Seconds())
		if err != nil {
			return nil, err
		}
		return s.view(*stored)
	default:
		return nil, ErrImageInvalid
	}
}

func (s *Service) view(stored models.StoredImage, variants ...models.ImageVariant) (*models.ImageView, error) {
	if !stored.Managed() {
		return nil, ErrImageNotFound
	}
	expires := s.now().UTC().Add(models.ImageSignedURLTTL).Truncate(time.Second)
	variant := models.ImageOriginal
	if len(variants) > 0 {
		variant = variants[0]
	}
	signStart := time.Now()
	signed, err := s.storage.SignedURL(stored.Path, expires, variant)
	imageStageDuration.WithLabelValues("url_sign").Observe(time.Since(signStart).Seconds())
	if err != nil {
		return nil, ErrImageProvider
	}
	return &models.ImageView{URL: signed, ExpiresAt: expires, Variant: variant}, nil
}

func (s *Service) replace(ctx context.Context, target models.ImageTarget, next models.StoredImage, requestID string) error {
	dbCtx, done := context.WithTimeout(ctx, models.ImageDatabaseTimeout)
	old, err := s.repo.Replace(dbCtx, target, next, func(txCtx context.Context) error {
		checked, authErr := s.auth.Authorize(txCtx, target, next.Managed())
		if authErr == nil && checked.VisitorID != target.VisitorID {
			return ErrImageTargetNotFound
		}
		return authErr
	})
	done()
	if err != nil {
		var uncertain *repository.ImageCommitUncertain
		if errors.As(err, &uncertain) {
			resolveCtx, cancel := context.WithTimeout(context.Background(), models.ImageDatabaseTimeout)
			current, resolveErr := s.repo.Resolve(resolveCtx, target)
			cancel()
			if resolveErr != nil {
				s.logCleanup(target, requestID, "commit_reconciliation_required", next)
				return ErrImagePersistence
			}
			if current == next {
				s.cleanup(target, old, requestID)
				return nil
			}
			// The new file is not referenced and cannot be reattached by this API.
			// Do not delete old: an uncertain commit may have rolled back.
			s.cleanup(target, next, requestID)
			s.logCleanup(target, requestID, "commit_superseded_or_rolled_back", old)
			return ErrImagePersistence
		}
		s.cleanup(target, next, requestID)
		return safeError(err)
	}
	s.cleanup(target, old, requestID)
	return nil
}

func (s *Service) cleanup(target models.ImageTarget, stored models.StoredImage, requestID string) {
	if stored.FileID == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), models.ImageCleanupTimeout)
	defer cancel()
	if err := s.storage.Delete(ctx, stored.FileID); err != nil {
		s.logCleanup(target, requestID, "delete_failed", stored)
	}
}
func (s *Service) logCleanup(target models.ImageTarget, requestID, category string, stored models.StoredImage) {
	cleanupFailures.Inc()
	logger.Warn("Image reconciliation required", zap.String("request_id", requestID), zap.String("category", category),
		zap.Int64("actor_id", target.ActorID), zap.Int64("society_id", target.SocietyID), zap.Int64("visitor_id", target.VisitorID),
		zap.String("file_id", stored.FileID), zap.String("file_path", stored.Path))
}
func safeError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return ErrImageTargetNotFound
	}
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		if appErr.Internal == nil {
			return appErr
		}
		return models.NewAppError(appErr.Code, appErr.Message, appErr.StatusCode, nil)
	}
	return ErrImagePersistence
}
