package hubsvc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/prometheus/client_golang/prometheus"
	"go-server/internal/config"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	imagesvc "go-server/internal/services/imageSvc"
	"go-server/pkg/logger"
	"go.uber.org/zap"
)

var cleanupFailures = prometheus.NewCounter(prometheus.CounterOpts{Name: "apna_gate_hub_upload_cleanup_failures_total", Help: "Society Hub upload deletions requiring retry."})

func init() { prometheus.MustRegister(cleanupFailures) }

func validateFile(data []byte, limits config.HubConfig) (mime, extension string, width, height *int32, err error) {
	if len(data) == 0 || int64(len(data)) > limits.UploadMaxBytes {
		return "", "", nil, nil, imagesvc.ErrImageTooLarge
	}
	mime = http.DetectContentType(data)
	if mime == "application/pdf" {
		if !bytes.HasPrefix(data, []byte("%PDF-")) || !bytes.Contains(data[max(0, len(data)-1024):], []byte("%%EOF")) {
			return "", "", nil, nil, ErrInvalid
		}
		return mime, ".pdf", nil, nil, nil
	}
	if len(data) > int(models.ImageMaxUploadBytes) {
		return "", "", nil, nil, imagesvc.ErrImageTooLarge
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", "", nil, nil, ErrInvalid
	}
	if cfg.Width < 1 || cfg.Width > 32767 || cfg.Height < 1 || cfg.Height > 32767 || cfg.Width > int(limits.ImageMaxWidth) || cfg.Height > int(limits.ImageMaxHeight) {
		return "", "", nil, nil, imagesvc.ErrImageTooLarge
	}
	extension, err = imagesvc.ValidateImage(data)
	if err != nil {
		return "", "", nil, nil, err
	}
	if extension == ".jpg" {
		mime = "image/jpeg"
	} else {
		mime = "image/png"
	}
	return mime, extension, ptr(int32(cfg.Width)), ptr(int32(cfg.Height)), nil
}

func (s *Service) Upload(ctx context.Context, society, user int64, w http.ResponseWriter, r *http.Request) (models.HubUploadView, error) {
	var out models.HubUploadView
	if _, e := s.access(ctx, society, user); e != nil {
		return out, e
	}
	if s.storage == nil {
		return out, imagesvc.ErrImageUnavailable
	}
	select {
	case s.slots <- struct{}{}:
		defer func() { <-s.slots }()
	default:
		return out, imagesvc.ErrImageBusy
	}
	ctx, cancel := context.WithTimeout(ctx, models.ImageProcessingTimeout)
	defer cancel()
	stop := context.AfterFunc(ctx, func() { _ = r.Body.Close() })
	defer stop()
	r.Body = http.MaxBytesReader(w, r.Body, s.limits.UploadMaxBytes+64*1024)
	reader, e := r.MultipartReader()
	if e != nil {
		return out, ErrInvalid
	}
	part, e := reader.NextPart()
	if e != nil || part.FormName() != "file" || part.FileName() == "" {
		return out, ErrInvalid
	}
	filename := path.Base(strings.ReplaceAll(part.FileName(), "\\", "/"))
	if len(filename) > 255 || strings.ContainsAny(filename, "\r\n\x00") {
		return out, ErrInvalid
	}
	data, e := io.ReadAll(io.LimitReader(part, s.limits.UploadMaxBytes+1))
	if e != nil {
		return out, ErrInvalid
	}
	mime, ext, width, height, e := validateFile(data, s.limits)
	if e != nil {
		return out, e
	}
	if _, e = reader.NextPart(); !errors.Is(e, io.EOF) {
		return out, ErrInvalid
	}
	if _, e = io.Copy(io.Discard, r.Body); e != nil {
		return out, ErrInvalid
	}
	folder := fmt.Sprintf("%s/societies/%d/hub/%d", s.prefix, society, user)
	stored, e := s.storage.Upload(ctx, models.ImageUploadInput{Reader: bytes.NewReader(data), FileName: uuid.NewString() + ext, Folder: folder})
	if e != nil || stored == nil || !stored.Managed() {
		if stored != nil && stored.FileID != "" {
			s.cleanupUnclaimed(*stored)
		}
		return out, imagesvc.ErrImageProvider
	}
	var upload contracts.ChannelUpload
	e = s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		var err error
		upload, err = s.repo.HubCreateUpload(tx, contracts.HubCreateUploadInput{SocietyID: society, UploaderID: user, FileID: stored.FileID, FilePath: stored.Path, OriginalFilename: filename, MimeType: mime, FileSize: int64(len(data)), Width: width, Height: height})
		return err
	})
	if e != nil {
		// Reconcile a possibly successful COMMIT before deleting a provider file.
		check, done := context.WithTimeout(context.WithoutCancel(ctx), models.ImageDatabaseTimeout)
		defer done()
		found, checkErr := s.repo.HubUploadByFile(check, stored.FileID)
		if checkErr == nil {
			upload = found
		} else {
			if errors.Is(checkErr, contracts.ErrNotFound) {
				s.cleanupUnclaimed(*stored)
			} else {
				logger.Error("hub upload needs reconciliation", zap.String("file_id", stored.FileID))
			}
			return out, e
		}
	}
	return uploadView(upload), nil
}
func uploadView(u contracts.ChannelUpload) models.HubUploadView {
	return models.HubUploadView{ID: u.ID, Filename: u.OriginalFilename, MimeType: u.MimeType, Size: u.FileSize, ExpiresAt: u.ExpiresAt}
}
func (s *Service) cleanupUnclaimed(stored models.StoredImage) {
	ctx, cancel := context.WithTimeout(context.Background(), models.ImageCleanupTimeout)
	defer cancel()
	if e := s.storage.Delete(ctx, stored.FileID); e != nil {
		cleanupFailures.Inc()
		logger.Error("hub upload cleanup requires reconciliation", zap.String("file_id", stored.FileID))
	}
}

func (s *Service) attach(ctx context.Context, society, user int64, post, comment *int64, ids *[]int64) error {
	if ids == nil {
		return nil
	}
	if !validAttachments(ids) {
		return ErrInvalid
	}
	old, e := s.repo.HubAttachments(ctx, contracts.HubAttachmentsInput{PostID: post, CommentID: comment})
	if e != nil {
		return e
	}
	previous := map[int64]bool{}
	for _, id := range old {
		previous[id] = true
	}
	wanted := map[int64]bool{}
	for _, id := range *ids {
		wanted[id] = true
	}
	all := append(append([]int64{}, old...), (*ids)...)
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	seen := map[int64]bool{}
	for _, id := range all {
		if seen[id] {
			continue
		}
		seen[id] = true
		u, e := s.repo.HubLockUpload(ctx, contracts.HubLockUploadInput{SocietyID: society, ID: id})
		if e != nil {
			return missing(e)
		}
		if wanted[id] && previous[id] {
			continue
		}
		if wanted[id] {
			if u.UploaderID != user {
				return ErrForbidden
			}
			n, e := s.repo.HubClaimUpload(ctx, contracts.HubClaimUploadInput{ID: id, SocietyID: society, UploaderID: user})
			if e != nil {
				return e
			}
			if n != 1 {
				return ErrConflict
			}
			if e = s.repo.HubAttach(ctx, contracts.HubAttachInput{UploadID: id, PostID: post, CommentID: comment}); e != nil {
				return e
			}
		} else {
			if e = s.repo.HubDetach(ctx, id); e != nil {
				return e
			}
			if e = s.repo.HubQueueUploadCleanup(ctx, id); e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *Service) DeleteUpload(ctx context.Context, society, user, id int64) error {
	return s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		u, e := s.repo.HubLockUpload(tx, contracts.HubLockUploadInput{SocietyID: society, ID: id})
		if e != nil {
			return missing(e)
		}
		if u.UploaderID != user {
			return ErrForbidden
		}
		if u.Status == "claimed" {
			return ErrConflict
		}
		return s.repo.HubQueueUploadCleanup(tx, id)
	})
}
func (s *Service) Attachment(ctx context.Context, society, user, id int64) (models.HubAttachmentView, error) {
	var out models.HubAttachmentView
	if _, e := s.access(ctx, society, user); e != nil {
		return out, e
	}
	if s.storage == nil {
		return out, imagesvc.ErrImageUnavailable
	}
	u, e := s.repo.HubUpload(ctx, contracts.HubUploadInput{SocietyID: society, ID: id})
	if e != nil {
		return out, missing(e)
	}
	switch u.Status {
	case "pending":
		if u.UploaderID != user || !u.ExpiresAt.After(time.Now()) {
			return out, ErrNotFound
		}
	case "claimed":
		target, e := s.repo.HubAttachmentTarget(ctx, id)
		if e != nil {
			return out, missing(e)
		}
		if target.SocietyID != society || target.PostStatus != "active" || (target.CommentStatus != nil && *target.CommentStatus != "active") {
			return out, ErrNotFound
		}
	default:
		return out, ErrNotFound
	}
	out.HubUploadView = uploadView(u)
	out.ExpiresAt = time.Now().UTC().Add(models.ImageSignedURLTTL).Truncate(time.Second)
	out.URL, e = s.storage.SignedURL(u.FilePath, out.ExpiresAt, models.ImageOriginal)
	if e != nil {
		return models.HubAttachmentView{}, imagesvc.ErrImageProvider
	}
	return out, nil
}

// CleanupUploads is called by the existing cleanup job. Each claim is committed
// before provider I/O; its lease permits recovery after process failure.
func (s *Service) CleanupUploads(ctx context.Context) error {
	if s.storage == nil {
		return nil
	}
	var failures []error
	for i := 0; i < 500; i++ {
		if e := ctx.Err(); e != nil {
			return e
		}
		token := uuid.New()
		u, e := s.repo.HubClaimCleanup(ctx, token)
		if errors.Is(e, contracts.ErrNotFound) {
			break
		}
		if e != nil {
			return e
		}
		work, done := context.WithTimeout(ctx, models.ImageCleanupTimeout)
		e = s.storage.Delete(work, u.FileID)
		done()
		if e == nil {
			e = s.repo.HubFinishCleanup(ctx, contracts.HubFinishCleanupInput{ID: u.ID, LeaseToken: token})
		}
		if e != nil {
			cleanupFailures.Inc()
			failures = append(failures, errors.New("hub upload cleanup failed"))
		}
	}
	return errors.Join(failures...)
}
