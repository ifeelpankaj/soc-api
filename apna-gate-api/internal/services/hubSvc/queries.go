package hubsvc

import (
	"context"

	"errors"
	"fmt"

	"time"

	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

func (s *Service) content(ctx context.Context, society, user, id int64, comment bool) (models.HubContent, error) {
	var result models.HubContent
	var err error
	if comment {
		result, err = s.repo.HubCommentView(ctx, contracts.HubCommentViewInput{SocietyID: society, UserID: user, ID: id})
	} else {
		result, err = s.repo.HubPostView(ctx, contracts.HubPostViewInput{SocietyID: society, UserID: user, ID: id})
	}
	return result, missing(err)
}
func (s *Service) Channels(ctx context.Context, society, user int64) ([]models.HubChannelSummary, error) {
	if _, err := s.access(ctx, society, user); err != nil {
		return nil, err
	}
	rows, err := s.repo.HubChannels(ctx, contracts.HubChannelsInput{SocietyID: society, UserID: user})
	if err != nil {
		return nil, err
	}
	out := make([]models.HubChannelSummary, 0, len(rows))
	for _, r := range rows {
		v := models.HubChannelSummary{ID: r.ID, Type: r.Type, Name: r.Name, Description: r.Description, UnreadCount: r.UnreadCount}
		if r.ImportantPostID > 0 {
			p, e := s.content(ctx, society, user, r.ImportantPostID, false)
			if e != nil && !errors.Is(e, ErrNotFound) {
				return nil, e
			}
			if e == nil {
				v.ImportantAnnouncement = &p
			}
		}
		out = append(out, v)
	}
	return out, nil
}
func (s *Service) Categories(ctx context.Context, society, user int64) ([]models.HubCategory, error) {
	if _, err := s.access(ctx, society, user); err != nil {
		return nil, err
	}
	rows, err := s.repo.HubCategories(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]models.HubCategory, 0, len(rows))
	for _, r := range rows {
		out = append(out, models.HubCategory{ID: r.ID, Code: r.Code, Name: r.Name, DisplayOrder: r.DisplayOrder})
	}
	return out, nil
}
func (s *Service) Post(ctx context.Context, society, user, id int64) (models.HubContent, error) {
	if _, err := s.access(ctx, society, user); err != nil {
		return models.HubContent{}, err
	}
	return s.content(ctx, society, user, id, false)
}
func (s *Service) Feed(ctx context.Context, society, user, target int64, raw string, limit int, category *int16, comments bool) (models.HubPage, error) {
	out := models.HubPage{Items: []models.HubContent{}}
	if _, err := s.access(ctx, society, user); err != nil {
		return out, err
	}
	scope := fmt.Sprintf("%d/%d/%t/%v", society, target, comments, categoryValue(category))
	c, err := decodeCursor(raw, scope)
	if err != nil {
		return out, err
	}
	size, err := pageSize(limit)
	if err != nil {
		return out, err
	}
	var ids []int64
	var at *time.Time
	if !c.Time.IsZero() {
		at = &c.Time
	}
	if comments {
		p, e := s.repo.HubPost(ctx, contracts.HubPostInput{SocietyID: society, ID: target})
		if e != nil {
			return out, missing(e)
		}
		if p.Status != "active" {
			return out, ErrNotFound
		}
		ids, err = s.repo.HubCommentIDs(ctx, contracts.HubCommentIDsInput{SocietyID: society, PostID: target, AfterAt: at, AfterID: c.ID, PageLimit: size + 1})
	} else {
		if _, e := s.repo.HubChannel(ctx, contracts.HubChannelInput{SocietyID: society, ID: target}); e != nil {
			return out, missing(e)
		}
		ids, err = s.repo.HubPostIDs(ctx, contracts.HubPostIDsInput{SocietyID: society, ChannelID: target, CategoryID: category, BeforeAt: at, BeforeID: c.ID, PageLimit: size + 1})
		if err == nil {
			pinned, e := s.repo.HubPinnedIDs(ctx, contracts.HubPinnedIDsInput{SocietyID: society, ChannelID: target})
			if e != nil {
				return out, e
			}
			for _, id := range pinned {
				p, e := s.content(ctx, society, user, id, false)
				if errors.Is(e, ErrNotFound) {
					continue
				}
				if e != nil {
					return out, e
				}
				out.Pinned = append(out.Pinned, p)
			}
		}
	}
	if err != nil {
		return out, err
	}
	more := len(ids) > int(size)
	if more {
		ids = ids[:size]
	}
	for _, id := range ids {
		p, e := s.content(ctx, society, user, id, comments)
		if errors.Is(e, ErrNotFound) {
			continue
		}
		if e != nil {
			return out, e
		}
		out.Items = append(out.Items, p)
	}
	if more && len(ids) > 0 {
		last := ids[len(ids)-1]
		var t time.Time
		if comments {
			r, e := s.repo.HubComment(ctx, contracts.HubCommentInput{SocietyID: society, ID: last})
			if e != nil {
				return out, e
			}
			t = r.CreatedAt
		} else {
			r, e := s.repo.HubPost(ctx, contracts.HubPostInput{SocietyID: society, ID: last})
			if e != nil {
				return out, e
			}
			t = r.CreatedAt
		}
		out.NextCursor = encodeCursor(scope, t, last)
	}
	return out, nil
}
func categoryValue(c *int16) int16 {
	if c == nil {
		return 0
	}
	return *c
}
