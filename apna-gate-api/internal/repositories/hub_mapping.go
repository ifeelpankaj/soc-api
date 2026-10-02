package repository

import (
	"github.com/google/uuid"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func mapChannelComment(v db.ChannelComment) contracts.ChannelComment {
	return contracts.ChannelComment{
		ID:            v.ID,
		PostID:        v.PostID,
		AuthorID:      v.AuthorID,
		ParentID:      v.ParentID,
		Body:          v.Body,
		Status:        v.Status,
		EditedAt:      pgTimestamptzToTimePtr(v.EditedAt),
		RemovalReason: v.RemovalReason,
		DeletedAt:     pgTimestamptzToTimePtr(v.DeletedAt),
		DeletedBy:     v.DeletedBy,
		CreatedAt:     pgTimestamptzToTime(v.CreatedAt),
		UpdatedAt:     pgTimestamptzToTime(v.UpdatedAt),
	}
}

func mapChannelReport(v db.ChannelReport) contracts.ChannelReport {
	return contracts.ChannelReport{
		ID:             v.ID,
		PostID:         v.PostID,
		ReporterID:     v.ReporterID,
		Reason:         v.Reason,
		Status:         v.Status,
		ReviewerID:     v.ReviewerID,
		ReviewedAt:     pgTimestamptzToTimePtr(v.ReviewedAt),
		ResolutionNote: v.ResolutionNote,
		CreatedAt:      pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapChannelUpload(v db.ChannelUpload) contracts.ChannelUpload {
	return contracts.ChannelUpload{
		ID:               v.ID,
		SocietyID:        v.SocietyID,
		UploaderID:       v.UploaderID,
		FileID:           v.FileID,
		FilePath:         v.FilePath,
		OriginalFilename: v.OriginalFilename,
		MimeType:         v.MimeType,
		FileSize:         v.FileSize,
		Width:            v.Width,
		Height:           v.Height,
		Status:           v.Status,
		ExpiresAt:        pgTimestamptzToTime(v.ExpiresAt),
		LeaseUntil:       pgTimestamptzToTimePtr(v.LeaseUntil),
		LeaseToken:       uuid.UUID(v.LeaseToken.Bytes),
		CreatedAt:        pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapCommunityCategory(v db.CommunityCategory) contracts.CommunityCategory {
	return contracts.CommunityCategory{
		ID:           v.ID,
		Code:         v.Code,
		Name:         v.Name,
		DisplayOrder: v.DisplayOrder,
		IsActive:     v.IsActive,
		CreatedAt:    pgTimestamptzToTime(v.CreatedAt),
		UpdatedAt:    pgTimestamptzToTime(v.UpdatedAt),
	}
}

func mapHubAttachmentTargetRow(v db.HubAttachmentTargetRow) contracts.HubAttachmentTargetRecord {
	return contracts.HubAttachmentTargetRecord{
		PostID:        v.PostID,
		CommentID:     v.CommentID,
		SocietyID:     v.SocietyID,
		PostStatus:    v.PostStatus,
		CommentStatus: v.CommentStatus,
	}
}

func mapHubChannelsRow(v db.HubChannelsRow) contracts.HubChannelsRecord {
	return contracts.HubChannelsRecord{
		ID:              v.ID,
		SocietyID:       v.SocietyID,
		Type:            v.Type,
		Name:            v.Name,
		Description:     v.Description,
		CreatedAt:       pgTimestamptzToTime(v.CreatedAt),
		UpdatedAt:       pgTimestamptzToTime(v.UpdatedAt),
		UnreadCount:     v.UnreadCount,
		ImportantPostID: v.ImportantPostID,
	}
}

func mapHubCommentRateRow(v db.HubCommentRateRow) contracts.HubCommentRateRecord {
	return contracts.HubCommentRateRecord{
		Count:  v.Count,
		Oldest: pgTimestamptzToTimePtr(v.Oldest),
	}
}

func mapHubLockPostRow(v db.HubLockPostRow) contracts.HubLockPostRecord {
	return contracts.HubLockPostRecord{
		ID:              v.ID,
		SocietyID:       v.SocietyID,
		ChannelID:       v.ChannelID,
		AuthorID:        v.AuthorID,
		Title:           v.Title,
		Body:            v.Body,
		CategoryID:      v.CategoryID,
		CommentsEnabled: v.CommentsEnabled,
		IsPinned:        v.IsPinned,
		IsImportant:     v.IsImportant,
		Status:          v.Status,
		EditedAt:        pgTimestamptzToTimePtr(v.EditedAt),
		RemovalReason:   v.RemovalReason,
		DeletedAt:       pgTimestamptzToTimePtr(v.DeletedAt),
		DeletedBy:       v.DeletedBy,
		CreatedAt:       pgTimestamptzToTime(v.CreatedAt),
		UpdatedAt:       pgTimestamptzToTime(v.UpdatedAt),
		ChannelType:     v.ChannelType,
	}
}

func mapHubPostRateRow(v db.HubPostRateRow) contracts.HubPostRateRecord {
	return contracts.HubPostRateRecord{
		Count:  v.Count,
		Oldest: pgTimestamptzToTimePtr(v.Oldest),
	}
}

func mapHubPostRow(v db.HubPostRow) contracts.HubPostRecord {
	return contracts.HubPostRecord{
		ID:              v.ID,
		SocietyID:       v.SocietyID,
		ChannelID:       v.ChannelID,
		AuthorID:        v.AuthorID,
		Title:           v.Title,
		Body:            v.Body,
		CategoryID:      v.CategoryID,
		CommentsEnabled: v.CommentsEnabled,
		IsPinned:        v.IsPinned,
		IsImportant:     v.IsImportant,
		Status:          v.Status,
		EditedAt:        pgTimestamptzToTimePtr(v.EditedAt),
		RemovalReason:   v.RemovalReason,
		DeletedAt:       pgTimestamptzToTimePtr(v.DeletedAt),
		DeletedBy:       v.DeletedBy,
		CreatedAt:       pgTimestamptzToTime(v.CreatedAt),
		UpdatedAt:       pgTimestamptzToTime(v.UpdatedAt),
		ChannelType:     v.ChannelType,
	}
}

func mapSocietyChannel(v db.SocietyChannel) contracts.SocietyChannel {
	return contracts.SocietyChannel{
		ID:          v.ID,
		SocietyID:   v.SocietyID,
		Type:        v.Type,
		Name:        v.Name,
		Description: v.Description,
		CreatedAt:   pgTimestamptzToTime(v.CreatedAt),
		UpdatedAt:   pgTimestamptzToTime(v.UpdatedAt),
	}
}
