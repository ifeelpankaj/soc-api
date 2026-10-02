package shortlinksvc

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"strings"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	service "go-server/internal/services"
)

const (
	codeLength       = 14
	createMaxRetries = 5
	codeAlphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type ShortLinkService interface {
	Create(ctx context.Context, resourceType models.ShortLinkResourceType, resourceID int64, expiresAt *time.Time, createdBy *int64) (*models.ShortLinkResponse, error)
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	GetActiveByCode(ctx context.Context, code string) (*models.ShortLink, error)
	Response(link *models.ShortLink) *models.ShortLinkResponse
}

type Service struct {
	repo      repository.ShortLinkRepository
	publicURL string
}

func NewShortLinkService(repo repository.ShortLinkRepository, publicURL string) ShortLinkService {
	return &Service{repo: repo, publicURL: strings.TrimRight(publicURL, "/")}
}

func (s *Service) Create(ctx context.Context, resourceType models.ShortLinkResourceType, resourceID int64, expiresAt *time.Time, createdBy *int64) (*models.ShortLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if !resourceType.IsValid() || resourceID <= 0 {
		return nil, ErrInvalidShortLink
	}

	for i := 0; i < createMaxRetries; i++ {
		code, err := randomCode(codeLength)
		if err != nil {
			return nil, err
		}
		link, err := s.repo.Create(ctx, code, resourceType, resourceID, expiresAt, createdBy, nil)
		if err == nil {
			return s.Response(link), nil
		}
		if !isUniqueViolation(err) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("short link code collision after %d attempts", createMaxRetries)
}

func (s *Service) GetActiveByCode(ctx context.Context, code string) (*models.ShortLink, error) {
	link, err := s.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if link.RevokedAt != nil {
		return nil, ErrShortLinkUnavailable
	}
	if link.ExpiresAt != nil && !link.ExpiresAt.After(time.Now()) {
		return nil, ErrShortLinkUnavailable
	}
	return link, nil
}

func (s *Service) GetByCode(ctx context.Context, code string) (*models.ShortLink, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	code = models.NormalizeShortCode(code)
	if code == "" {
		return nil, ErrInvalidShortLink
	}
	link, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, ErrShortLinkNotFound
	}
	return link, nil
}

func (s *Service) Response(link *models.ShortLink) *models.ShortLinkResponse {
	if link == nil {
		return nil
	}
	return &models.ShortLinkResponse{
		ShortCode: link.ShortCode,
		URL:       s.publicShortLinkURL(link.ShortCode),
		ExpiresAt: link.ExpiresAt,
	}
}

func (s *Service) publicShortLinkURL(code string) string {
	return s.publicURL + "/link/" + code
}

func randomCode(length int) (string, error) {
	var b strings.Builder
	b.Grow(length)
	max := big.NewInt(int64(len(codeAlphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b.WriteByte(codeAlphabet[n.Int64()])
	}
	return b.String(), nil
}

func isUniqueViolation(err error) bool {
	return errors.Is(err, repository.ErrAlreadyExists)
}
