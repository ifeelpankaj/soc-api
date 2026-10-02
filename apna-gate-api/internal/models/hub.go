package models

import "time"

type HubCategory struct {
	ID           int16  `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	DisplayOrder int16  `json:"display_order"`
}
type HubReport struct {
	ID             int64      `json:"id"`
	PostID         int64      `json:"post_id"`
	ReporterID     int64      `json:"reporter_id"`
	Reason         string     `json:"reason"`
	Status         string     `json:"status"`
	ReviewerID     *int64     `json:"reviewer_id"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
	ResolutionNote *string    `json:"resolution_note"`
	CreatedAt      time.Time  `json:"created_at"`
}
type HubReportPage struct {
	Items      []HubReport `json:"items"`
	NextCursor string      `json:"next_cursor,omitempty"`
}

type HubPostRequest struct {
	Title           *string  `json:"title"`
	Body            *string  `json:"body"`
	CategoryID      *int16   `json:"category_id"`
	CommentsEnabled *bool    `json:"comments_enabled"`
	IsPinned        *bool    `json:"is_pinned"`
	IsImportant     *bool    `json:"is_important"`
	AttachmentIDs   *[]int64 `json:"attachment_ids"`
}
type HubCommentRequest struct {
	Body          *string  `json:"body"`
	ParentID      *int64   `json:"parent_id"`
	AttachmentIDs *[]int64 `json:"attachment_ids"`
}
type HubAuthor struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type HubContent struct {
	ID              int64            `json:"id"`
	SocietyID       int64            `json:"society_id,omitempty"`
	ChannelID       int64            `json:"channel_id,omitempty"`
	PostID          int64            `json:"post_id,omitempty"`
	ParentID        *int64           `json:"parent_id,omitempty"`
	Author          *HubAuthor       `json:"author"`
	Title           *string          `json:"title,omitempty"`
	Body            string           `json:"body"`
	CategoryID      *int16           `json:"category_id,omitempty"`
	CommentsEnabled bool             `json:"comments_enabled"`
	IsPinned        bool             `json:"is_pinned"`
	IsImportant     bool             `json:"is_important"`
	Status          string           `json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	EditedAt        *time.Time       `json:"edited_at"`
	ReplyCount      int64            `json:"reply_count"`
	Reactions       map[string]int64 `json:"reactions"`
	MyReactions     []string         `json:"my_reactions"`
	AttachmentIDs   []int64          `json:"attachment_ids"`
}
type HubPage struct {
	Items      []HubContent `json:"items"`
	Pinned     []HubContent `json:"pinned,omitempty"`
	NextCursor string       `json:"next_cursor,omitempty"`
}
type HubChannelSummary struct {
	ID                    int64       `json:"id"`
	Type                  string      `json:"type"`
	Name                  string      `json:"name"`
	Description           *string     `json:"description"`
	UnreadCount           int64       `json:"unread_count"`
	ImportantAnnouncement *HubContent `json:"important_announcement,omitempty"`
}
type HubUploadView struct {
	ID        int64     `json:"id"`
	Filename  string    `json:"filename"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	ExpiresAt time.Time `json:"expires_at"`
}
type HubAttachmentView struct {
	HubUploadView
	URL string `json:"url"`
}
type HubControlRequest struct {
	IsPinned *bool `json:"is_pinned"`
	Locked   *bool `json:"locked"`
}
type HubReactionRequest struct {
	ReactionType string `json:"reaction_type" enums:"like,love,helpful" binding:"required"`
}
type HubReadRequest struct {
	PostID int64 `json:"post_id" minimum:"1"`
}
type HubReasonRequest struct {
	Reason string `json:"reason"`
}
type HubResolveRequest struct {
	Status string `json:"status" enums:"dismissed,actioned" binding:"required"`
	Note   string `json:"note"`
}

// HubResponse is the standard API envelope. Each handler refines the data schema.
type HubResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
