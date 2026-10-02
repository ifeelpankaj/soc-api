package hubsvc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"strings"
	"time"
	"unicode/utf8"

	"go-server/internal/repositories/contracts"

	"go-server/internal/requestctx"
)

func (s *Service) access(ctx context.Context, society, user int64) (string, error) {
	if society <= 0 || user <= 0 {
		return "", ErrForbidden
	}
	role, err := s.repo.HubRole(ctx, contracts.HubRoleInput{SocietyID: society, UserID: user})
	if errors.Is(err, contracts.ErrNotFound) {
		return "", ErrForbidden
	}
	if err != nil {
		return "", err
	}
	if s.operational == nil {
		return "", errors.New("hub operational guard unavailable")
	}
	if err = s.operational.EnsureSocietyOperational(requestctx.WithoutDeveloperGuardBypass(ctx), society); err != nil {
		return "", err
	}
	return role, nil
}
func admin(role string) bool { return role == "owner" || role == "admin" }
func missing(err error) error {
	if errors.Is(err, contracts.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
func ptr[T any](v T) *T { return &v }
func validBody(body string, max int) bool {
	return strings.TrimSpace(body) != "" && utf8.RuneCountInString(body) <= max
}
func validAttachments(ids *[]int64) bool {
	if ids == nil {
		return true
	}
	if len(*ids) > 5 {
		return false
	}
	seen := map[int64]bool{}
	for _, id := range *ids {
		if id <= 0 || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}
func (s *Service) mutation(ctx context.Context, society, user int64, fn func(context.Context, string) error) error {
	return s.repo.WithTransaction(ctx, func(tx context.Context) error {
		if err := s.repo.HubActorLock(tx, contracts.HubActorLockInput{SocietyID: society, UserID: user}); err != nil {
			return err
		}
		role, err := s.access(tx, society, user)
		if err != nil {
			return err
		}
		return fn(tx, role)
	})
}
func (s *Service) rate(ctx context.Context, society, user int64, comment bool) error {
	var count int64
	var oldest *time.Time
	limit, window := s.limits.PostsPerWindow, s.limits.PostWindowSeconds
	if comment {
		limit, window = s.limits.CommentsPerWindow, s.limits.CommentWindowSeconds
		r, e := s.repo.HubCommentRate(ctx, contracts.HubCommentRateInput{SocietyID: society, AuthorID: user, WindowSeconds: window})
		if e != nil {
			return e
		}
		count, oldest = r.Count, r.Oldest
	} else {
		r, e := s.repo.HubPostRate(ctx, contracts.HubPostRateInput{SocietyID: society, AuthorID: user, WindowSeconds: window})
		if e != nil {
			return e
		}
		count, oldest = r.Count, r.Oldest
	}
	if count >= int64(limit) {
		retry := int(time.Until(oldest.Add(time.Duration(window)*time.Second)).Seconds()) + 1
		return &RateError{RetryAfter: max(1, retry)}
	}
	return nil
}

type cursor struct {
	Scope string    `json:"s"`
	Time  time.Time `json:"t"`
	ID    int64     `json:"i"`
}

func decodeCursor(raw, scope string) (cursor, error) {
	c := cursor{Scope: scope}
	if raw == "" {
		return c, nil
	}
	if len(raw) > 1024 {
		return c, ErrInvalid
	}
	b, e := base64.RawURLEncoding.DecodeString(raw)
	if e != nil {
		return c, ErrInvalid
	}
	if json.Unmarshal(b, &c) != nil || c.Scope != scope || c.ID <= 0 || c.Time.IsZero() {
		return c, ErrInvalid
	}
	return c, nil
}
func encodeCursor(scope string, t time.Time, id int64) string {
	b, _ := json.Marshal(cursor{scope, t, id})
	return base64.RawURLEncoding.EncodeToString(b)
}
func pageSize(n int) (int32, error) {
	if n == 0 {
		return 20, nil
	}
	if n < 1 || n > 100 {
		return 0, ErrInvalid
	}
	return int32(n), nil
}
