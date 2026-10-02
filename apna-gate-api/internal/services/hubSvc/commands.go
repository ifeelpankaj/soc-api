package hubsvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

//	func (s *Service) SavePost(ctx context.Context, society, user, channel, id int64, req models.HubPostRequest) (models.HubContent, error) {
//		var out models.HubContent
//		if !validAttachments(req.AttachmentIDs) {
//			return out, ErrInvalid
//		}
//		err := s.mutation(ctx, society, user, func(tx context.Context, role string) error {
//			p := contracts.HubLockPostRecord{CommentsEnabled: true}
//			creating := id == 0
//			if creating {
//				// Serialize publication order so a late commit cannot fall behind an observed read watermark.
//				c, e := s.repo.HubLockChannel(tx, contracts.HubLockChannelInput{SocietyID: society, ID: channel})
//				if e != nil {
//					return missing(e)
//				}
//				p.ChannelType = c.Type
//				if e = s.rate(tx, society, user, false); e != nil {
//					return e
//				}
//			} else {
//				var e error
//				p, e = s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: id})
//				if e != nil {
//					return missing(e)
//				}
//				if p.Status != "active" {
//					return ErrNotFound
//				}
//				if p.AuthorID != user {
//					return ErrForbidden
//				}
//			}
//			if p.ChannelType == "announcement" && !admin(role) {
//				return ErrForbidden
//			}
//			if req.Title != nil {
//				v := strings.TrimSpace(*req.Title)
//				if utf8.RuneCountInString(v) > 200 {
//					return ErrInvalid
//				}
//				p.Title = nil
//				if v != "" {
//					p.Title = &v
//				}
//			}
//			if req.Body != nil {
//				p.Body = strings.TrimSpace(*req.Body)
//			}
//			if !validBody(p.Body, 20000) {
//				return ErrInvalid
//			}
//			if p.ChannelType == "community" {
//				if req.IsPinned != nil || req.IsImportant != nil || req.CommentsEnabled != nil {
//					return ErrForbidden
//				}
//				if creating || req.CategoryID != nil {
//					category, e := s.repo.HubCategory(tx, req.CategoryID)
//					if e != nil {
//						return ErrInvalid
//					}
//					p.CategoryID = &category
//				}
//			} else {
//				if req.CategoryID != nil {
//					return ErrInvalid
//				}
//				if req.CommentsEnabled != nil {
//					p.CommentsEnabled = *req.CommentsEnabled
//				}
//				if req.IsPinned != nil {
//					p.IsPinned = *req.IsPinned
//				}
//				if req.IsImportant != nil {
//					p.IsImportant = *req.IsImportant
//				}
//			}
//			if creating {
//				var e error
//				id, e = s.repo.HubCreatePost(tx, contracts.HubCreatePostInput{SocietyID: society, ChannelID: channel, AuthorID: user, Title: p.Title, Body: p.Body, CategoryID: p.CategoryID, CommentsEnabled: p.CommentsEnabled, IsPinned: p.IsPinned, IsImportant: p.IsImportant})
//				if e != nil {
//					return e
//				}
//			} else {
//				if e := s.repo.HubUpdatePost(tx, contracts.HubUpdatePostInput{SocietyID: society, ID: id, Title: p.Title, Body: p.Body, CategoryID: p.CategoryID, CommentsEnabled: p.CommentsEnabled, IsPinned: p.IsPinned, IsImportant: p.IsImportant}); e != nil {
//					return e
//				}
//			}
//			if e := s.attach(tx, society, user, &id, nil, req.AttachmentIDs); e != nil {
//				return e
//			}
//			if creating && p.ChannelType == "announcement" {
//				users, e := s.repo.HubResidents(tx, society)
//				if e != nil {
//					return e
//				}
//				if e = s.notify(tx, society, user, id, users, "announcement", fmt.Sprint(id), true, p.IsImportant); e != nil {
//					return e
//				}
//			}
//			var e error
//			out, e = s.content(tx, society, user, id, false)
//			return e
//		})
//		return out, err
//	}
func (s *Service) SavePost(
	ctx context.Context,
	society, user, channel, id int64,
	req models.HubPostRequest,
) (models.HubContent, error) {
	var out models.HubContent

	if !validAttachments(req.AttachmentIDs) {
		return out, ErrInvalid
	}

	err := s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		p := contracts.HubLockPostRecord{
			CommentsEnabled: true,
		}

		creating := id == 0

		if creating {
			// Serialize publication order so a late commit cannot fall behind
			// an observed read watermark.
			c, err := s.repo.HubLockChannel(tx, contracts.HubLockChannelInput{
				SocietyID: society,
				ID:        channel,
			})
			if err != nil {
				return missing(err)
			}

			p.ChannelType = c.Type

			if err := s.rate(tx, society, user, false); err != nil {
				return err
			}
		} else {
			existing, err := s.repo.HubLockPost(tx, contracts.HubLockPostInput{
				SocietyID: society,
				ID:        id,
			})
			if err != nil {
				return missing(err)
			}

			p = existing

			if p.Status != "active" {
				return ErrNotFound
			}

			if p.AuthorID != user {
				return ErrForbidden
			}
		}

		if p.ChannelType == "announcement" && !admin(role) {
			return ErrForbidden
		}

		if req.Title != nil {
			title := strings.TrimSpace(*req.Title)

			if utf8.RuneCountInString(title) > 200 {
				return ErrInvalid
			}

			p.Title = nil
			if title != "" {
				p.Title = &title
			}
		}

		if req.Body != nil {
			p.Body = strings.TrimSpace(*req.Body)
		}

		if !validBody(p.Body, 20000) {
			return ErrInvalid
		}

		switch p.ChannelType {
		case "community":
			if req.IsPinned != nil ||
				req.IsImportant != nil ||
				req.CommentsEnabled != nil {
				return ErrForbidden
			}

			if creating || req.CategoryID != nil {
				categoryID, err := s.repo.HubCategory(tx, req.CategoryID)
				if err != nil {
					return ErrInvalid
				}

				p.CategoryID = &categoryID
			}

		case "announcement":
			if req.CategoryID != nil {
				return ErrInvalid
			}

			if req.CommentsEnabled != nil {
				p.CommentsEnabled = *req.CommentsEnabled
			}

			if req.IsPinned != nil {
				p.IsPinned = *req.IsPinned
			}

			if req.IsImportant != nil {
				p.IsImportant = *req.IsImportant
			}

		default:
			return ErrInvalid
		}

		if creating {
			postID, err := s.repo.HubCreatePost(tx, contracts.HubCreatePostInput{
				SocietyID:       society,
				ChannelID:       channel,
				AuthorID:        user,
				Title:           p.Title,
				Body:            p.Body,
				CategoryID:      p.CategoryID,
				CommentsEnabled: p.CommentsEnabled,
				IsPinned:        p.IsPinned,
				IsImportant:     p.IsImportant,
			})
			if err != nil {
				return err
			}

			id = postID
		} else {
			if err := s.repo.HubUpdatePost(tx, contracts.HubUpdatePostInput{
				SocietyID:       society,
				ID:              id,
				Title:           p.Title,
				Body:            p.Body,
				CategoryID:      p.CategoryID,
				CommentsEnabled: p.CommentsEnabled,
				IsPinned:        p.IsPinned,
				IsImportant:     p.IsImportant,
			}); err != nil {
				return err
			}
		}

		if err := s.attach(
			tx,
			society,
			user,
			&id,
			nil,
			req.AttachmentIDs,
		); err != nil {
			return err
		}

		if creating && p.ChannelType == "announcement" {
			if err := s.repo.HubQueueAnnouncementFanout(tx, society, id, user, p.IsImportant); err != nil {
				return err
			}
		}

		content, err := s.content(tx, society, user, id, false)
		if err != nil {
			return err
		}

		out = content
		return nil
	})

	return out, err
}
func (s *Service) SaveComment(ctx context.Context, society, user, post, id int64, req models.HubCommentRequest) (models.HubContent, error) {
	var out models.HubContent
	if !validAttachments(req.AttachmentIDs) {
		return out, ErrInvalid
	}
	err := s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		creating := id == 0
		var old contracts.ChannelComment
		if !creating {
			var e error
			old, e = s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: id})
			if e != nil {
				return missing(e)
			}
			post = old.PostID
		}
		p, e := s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: post})
		if e != nil {
			return missing(e)
		}
		if p.Status != "active" {
			return ErrNotFound
		}
		if !p.CommentsEnabled {
			return ErrConflict
		}
		if !creating {
			old, e = s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: id})
			if e != nil {
				return missing(e)
			}
			if old.Status != "active" {
				return ErrNotFound
			}
			if old.AuthorID != user {
				return ErrForbidden
			}
			if req.ParentID != nil {
				return ErrInvalid
			}
		}
		body := old.Body
		if req.Body != nil {
			body = strings.TrimSpace(*req.Body)
		}
		if !validBody(body, 5000) {
			return ErrInvalid
		}
		recipients := []int64{p.AuthorID}
		if creating {
			if e = s.rate(tx, society, user, true); e != nil {
				return e
			}
			if req.ParentID != nil {
				parent, e := s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: *req.ParentID})
				if e != nil {
					return missing(e)
				}
				if parent.PostID != post || parent.ParentID != nil || parent.Status != "active" {
					return ErrInvalid
				}
				recipients = append(recipients, parent.AuthorID)
			}
			id, e = s.repo.HubCreateComment(tx, contracts.HubCreateCommentInput{PostID: post, AuthorID: user, ParentID: req.ParentID, Body: body})
			if e != nil {
				return e
			}
		} else {
			if e = s.repo.HubUpdateComment(tx, contracts.HubUpdateCommentInput{ID: id, PostID: post, Body: body}); e != nil {
				return e
			}
		}
		if e = s.attach(tx, society, user, nil, &id, req.AttachmentIDs); e != nil {
			return e
		}
		if creating {
			if e = s.notify(tx, society, user, post, recipients, "reply", fmt.Sprint(id), true, false); e != nil {
				return e
			}
		}
		out, e = s.content(tx, society, user, id, true)
		return e
	})
	return out, err
}

func (s *Service) Delete(ctx context.Context, society, user, id int64, comment bool, reason string) error {
	reason = strings.TrimSpace(reason)
	if utf8.RuneCountInString(reason) > 2000 {
		return ErrInvalid
	}
	return s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		post := id
		var c contracts.ChannelComment
		var e error
		if comment {
			c, e = s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: id})
			if e != nil {
				return missing(e)
			}
			post = c.PostID
		}
		p, e := s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: post})
		if e != nil {
			return missing(e)
		}
		author, status := p.AuthorID, p.Status
		if comment {
			c, e = s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: id})
			if e != nil {
				return missing(e)
			}
			author, status = c.AuthorID, c.Status
		}
		if author != user && !admin(role) {
			return ErrForbidden
		}
		if status != "active" {
			return nil
		}
		state := "deleted"
		if author != user {
			state = "removed"
			if reason == "" {
				return ErrInvalid
			}
		}
		if comment {
			e = s.repo.HubDeleteComment(tx, contracts.HubDeleteCommentInput{ID: id, PostID: post, Status: state, DeletedBy: &user, RemovalReason: &reason})
		} else {
			e = s.repo.HubDeletePost(tx, contracts.HubDeletePostInput{SocietyID: society, ID: id, Status: state, DeletedBy: &user, RemovalReason: &reason})
		}
		if e != nil {
			return e
		}
		if admin(role) && !comment {
			if e = s.repo.HubResolvePostReports(tx, contracts.HubResolvePostReportsInput{PostID: post, ReviewerID: &user, ResolutionNote: &reason}); e != nil {
				return e
			}
		}
		if state == "removed" {
			return s.notify(tx, society, user, post, []int64{author}, "removed", fmt.Sprintf("%t/%d", comment, id), false, false)
		}
		return nil
	})
}

func (s *Service) Control(ctx context.Context, society, user, id int64, pin, locked *bool) error {
	if (pin == nil) == (locked == nil) {
		return ErrInvalid
	}
	return s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		if !admin(role) {
			return ErrForbidden
		}
		p, e := s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: id})
		if e != nil {
			return missing(e)
		}
		if p.Status != "active" {
			return ErrNotFound
		}
		if pin != nil {
			if p.ChannelType != "announcement" {
				return ErrInvalid
			}
			p.IsPinned = *pin
		}
		if locked != nil {
			p.CommentsEnabled = !*locked
		}
		return s.repo.HubControlPost(tx, contracts.HubControlPostInput{SocietyID: society, ID: id, IsPinned: p.IsPinned, CommentsEnabled: p.CommentsEnabled})
	})
}

func (s *Service) Reaction(ctx context.Context, society, user, id int64, comment, remove bool, kind string) error {
	if kind != "like" && kind != "love" && kind != "helpful" {
		return ErrInvalid
	}
	return s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		post := id
		var postID, commentID *int64
		if comment {
			c, e := s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: id})
			if e != nil {
				return missing(e)
			}
			post = c.PostID
			commentID = &id
		} else {
			postID = &id
		}
		p, e := s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: post})
		if e != nil {
			return missing(e)
		}
		if p.Status != "active" {
			return ErrNotFound
		}
		if comment {
			c, e := s.repo.HubComment(tx, contracts.HubCommentInput{SocietyID: society, ID: id})
			if e != nil {
				return missing(e)
			}
			if c.Status != "active" {
				return ErrNotFound
			}
		}
		if remove {
			return s.repo.HubRemoveReaction(tx, contracts.HubRemoveReactionInput{UserID: user, PostID: postID, CommentID: commentID, ReactionType: kind})
		}
		rid, e := s.repo.HubAddReaction(tx, contracts.HubAddReactionInput{UserID: user, PostID: postID, CommentID: commentID, ReactionType: kind})
		if errors.Is(e, contracts.ErrNotFound) {
			return nil
		}
		if e != nil {
			return e
		}
		if !comment {
			return s.notify(tx, society, user, post, []int64{p.AuthorID}, "reaction", fmt.Sprint(rid), false, false)
		}
		return nil
	})
}

func (s *Service) Read(ctx context.Context, society, user, channel, post int64) error {
	return s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		p, e := s.repo.HubPost(tx, contracts.HubPostInput{SocietyID: society, ID: post})
		if e != nil {
			return missing(e)
		}
		if p.ChannelID != channel {
			return ErrNotFound
		}
		return s.repo.HubRead(tx, contracts.HubReadInput{SocietyID: society, UserID: user, ChannelID: channel, PostID: post})
	})
}

func (s *Service) notify(ctx context.Context, society, actor, post int64, users []int64, event, suffix string, push, important bool) error {
	seen := map[int64]bool{}
	key := "hub." + event + ":" + suffix
	for _, user := range users {
		if user == actor || seen[user] {
			continue
		}
		seen[user] = true
		if _, e := s.repo.HubRole(ctx, contracts.HubRoleInput{SocietyID: society, UserID: user}); errors.Is(e, contracts.ErrNotFound) {
			continue
		} else if e != nil {
			return e
		}
		title, body := "Society Hub update", "Open Society Hub to view the update."
		switch event {
		case "announcement":
			title = "New society announcement"
		case "reply":
			title = "New reply to your discussion"
		case "reaction":
			title = "Someone reacted to your post"
		case "removed":
			title = "Your content was removed by society management"
		}
		priority := "normal"
		if important {
			priority = "high"
			title = "Important society announcement"
		}
		n := models.NotificationCreate{UserID: user, SocietyID: &society, Type: "hub." + event, Title: title, Body: body, EventKey: &key, Data: map[string]any{"type": "hub." + event, "society_id": fmt.Sprint(society), "post_id": fmt.Sprint(post), "priority": priority, "event_key": key}}
		payload, e := json.Marshal(n)
		if e != nil {
			return e
		}
		if e = s.repo.HubEnqueue(ctx, contracts.HubEnqueueInput{UserID: user, SocietyID: society, EventKey: key, Payload: payload, PushEnabled: push}); e != nil {
			return e
		}
	}
	return nil
}

func (s *Service) Report(ctx context.Context, society, user, post int64, reason string) (models.HubReport, error) {
	var out contracts.ChannelReport
	reason = strings.TrimSpace(reason)
	if !validBody(reason, 2000) {
		return models.HubReport{}, ErrInvalid
	}
	err := s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		p, e := s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: post})
		if e != nil {
			return missing(e)
		}
		if p.Status != "active" {
			return ErrNotFound
		}
		out, e = s.repo.HubReport(tx, contracts.HubReportInput{PostID: post, ReporterID: user, Reason: reason})
		return e
	})
	return reportView(out), err
}

func (s *Service) Reports(ctx context.Context, society, user int64, status, raw string, limit int) (models.HubReportPage, error) {
	out := models.HubReportPage{Items: []models.HubReport{}}
	role, e := s.access(ctx, society, user)
	if e != nil {
		return out, e
	}
	if !admin(role) {
		return out, ErrForbidden
	}
	if status != "" && status != "pending" && status != "dismissed" && status != "actioned" {
		return out, ErrInvalid
	}
	scope := fmt.Sprintf("reports/%d/%s", society, status)
	c, e := decodeCursor(raw, scope)
	if e != nil {
		return out, e
	}
	size, e := pageSize(limit)
	if e != nil {
		return out, e
	}
	rows, e := s.repo.HubReports(ctx, contracts.HubReportsInput{SocietyID: society, Status: status, AfterID: c.ID, PageLimit: size + 1})
	if e != nil {
		return out, e
	}
	if len(rows) > int(size) {
		rows = rows[:size]
		last := rows[len(rows)-1]
		out.NextCursor = encodeCursor(scope, last.CreatedAt, last.ID)
	}
	for _, r := range rows {
		out.Items = append(out.Items, reportView(r))
	}
	return out, nil
}
func (s *Service) Resolve(ctx context.Context, society, user, id int64, req models.HubResolveRequest) error {
	if (req.Status != "dismissed" && req.Status != "actioned") || !validBody(req.Note, 2000) {
		return ErrInvalid
	}
	return s.mutation(ctx, society, user, func(tx context.Context, role string) error {
		if !admin(role) {
			return ErrForbidden
		}
		r, e := s.repo.HubReportByID(tx, contracts.HubReportByIDInput{SocietyID: society, ID: id})
		if e != nil {
			return missing(e)
		}
		p, e := s.repo.HubLockPost(tx, contracts.HubLockPostInput{SocietyID: society, ID: r.PostID})
		if e != nil {
			return missing(e)
		}
		r, e = s.repo.HubReportByID(tx, contracts.HubReportByIDInput{SocietyID: society, ID: id})
		if e != nil {
			return missing(e)
		}
		if r.Status != "pending" {
			if r.Status == req.Status {
				return nil
			}
			return ErrConflict
		}
		if req.Status == "actioned" {
			if p.Status == "active" {
				if e = s.repo.HubDeletePost(tx, contracts.HubDeletePostInput{SocietyID: society, ID: p.ID, Status: "removed", DeletedBy: &user, RemovalReason: &req.Note}); e != nil {
					return e
				}
				if e = s.notify(tx, society, user, p.ID, []int64{p.AuthorID}, "removed", fmt.Sprintf("false/%d", p.ID), false, false); e != nil {
					return e
				}
			}
			return s.repo.HubResolvePostReports(tx, contracts.HubResolvePostReportsInput{PostID: p.ID, ReviewerID: &user, ResolutionNote: &req.Note})
		}
		return s.repo.HubResolveReport(tx, contracts.HubResolveReportInput{ID: id, Status: req.Status, ReviewerID: &user, ResolutionNote: &req.Note})
	})
}

func reportView(r contracts.ChannelReport) models.HubReport {
	v := models.HubReport{ID: r.ID, PostID: r.PostID, ReporterID: r.ReporterID, Reason: r.Reason, Status: r.Status, ReviewerID: r.ReviewerID, ResolutionNote: r.ResolutionNote, CreatedAt: r.CreatedAt}
	if r.ReviewedAt != nil {
		v.ReviewedAt = r.ReviewedAt
	}
	return v
}
