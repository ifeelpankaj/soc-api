package contracts

import (
	"github.com/google/uuid"
	"time"
)

type ChannelComment struct {
	ID            int64
	PostID        int64
	AuthorID      int64
	ParentID      *int64
	Body          string
	Status        string
	EditedAt      *time.Time
	RemovalReason *string
	DeletedAt     *time.Time
	DeletedBy     *int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ChannelReport struct {
	ID             int64
	PostID         int64
	ReporterID     int64
	Reason         string
	Status         string
	ReviewerID     *int64
	ReviewedAt     *time.Time
	ResolutionNote *string
	CreatedAt      time.Time
}

type ChannelUpload struct {
	ID               int64
	SocietyID        int64
	UploaderID       int64
	FileID           string
	FilePath         string
	OriginalFilename string
	MimeType         string
	FileSize         int64
	Width            *int32
	Height           *int32
	Status           string
	ExpiresAt        time.Time
	LeaseUntil       *time.Time
	LeaseToken       uuid.UUID
	CreatedAt        time.Time
}

type CommunityCategory struct {
	ID           int16
	Code         string
	Name         string
	DisplayOrder int16
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type HubActorLockInput struct {
	SocietyID int64
	UserID    int64
}

type HubAddReactionInput struct {
	UserID       int64
	PostID       *int64
	CommentID    *int64
	ReactionType string
}

type HubAttachInput struct {
	UploadID  int64
	PostID    *int64
	CommentID *int64
}

type HubAttachmentTargetRecord struct {
	PostID        *int64
	CommentID     *int64
	SocietyID     int64
	PostStatus    string
	CommentStatus *string
}

type HubAttachmentsInput struct {
	PostID    *int64
	CommentID *int64
}

type HubChannelInput struct {
	SocietyID int64
	ID        int64
}

type HubChannelsInput struct {
	UserID    int64
	SocietyID int64
}

type HubChannelsRecord struct {
	ID              int64
	SocietyID       int64
	Type            string
	Name            string
	Description     *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UnreadCount     int64
	ImportantPostID int64
}

type HubClaimUploadInput struct {
	ID         int64
	SocietyID  int64
	UploaderID int64
}

type HubCommentIDsInput struct {
	SocietyID int64
	PostID    int64
	AfterAt   *time.Time
	AfterID   int64
	PageLimit int32
}

type HubCommentInput struct {
	SocietyID int64
	ID        int64
}

type HubCommentRateInput struct {
	SocietyID     int64
	AuthorID      int64
	WindowSeconds int32
}

type HubCommentRateRecord struct {
	Count  int64
	Oldest *time.Time
}

type HubCommentViewInput struct {
	UserID    int64
	SocietyID int64
	ID        int64
}

type HubControlPostInput struct {
	SocietyID       int64
	ID              int64
	IsPinned        bool
	CommentsEnabled bool
}

type HubCreateCommentInput struct {
	PostID   int64
	AuthorID int64
	ParentID *int64
	Body     string
}

type HubCreatePostInput struct {
	SocietyID       int64
	ChannelID       int64
	AuthorID        int64
	Title           *string
	Body            string
	CategoryID      *int16
	CommentsEnabled bool
	IsPinned        bool
	IsImportant     bool
}

type HubCreateUploadInput struct {
	SocietyID        int64
	UploaderID       int64
	FileID           string
	FilePath         string
	OriginalFilename string
	MimeType         string
	FileSize         int64
	Width            *int32
	Height           *int32
}

type HubDeleteCommentInput struct {
	ID            int64
	PostID        int64
	Status        string
	DeletedBy     *int64
	RemovalReason *string
}

type HubDeletePostInput struct {
	SocietyID     int64
	ID            int64
	Status        string
	DeletedBy     *int64
	RemovalReason *string
}

type HubEnqueueInput struct {
	UserID      int64
	SocietyID   int64
	EventKey    string
	Payload     []byte
	PushEnabled bool
}

type HubFinishCleanupInput struct {
	ID         int64
	LeaseToken uuid.UUID
}

type HubLockChannelInput struct {
	SocietyID int64
	ID        int64
}

type HubLockPostInput struct {
	SocietyID int64
	ID        int64
}

type HubLockPostRecord struct {
	ID              int64
	SocietyID       int64
	ChannelID       int64
	AuthorID        int64
	Title           *string
	Body            string
	CategoryID      *int16
	CommentsEnabled bool
	IsPinned        bool
	IsImportant     bool
	Status          string
	EditedAt        *time.Time
	RemovalReason   *string
	DeletedAt       *time.Time
	DeletedBy       *int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ChannelType     string
}

type HubLockUploadInput struct {
	SocietyID int64
	ID        int64
}

type HubPinnedIDsInput struct {
	SocietyID int64
	ChannelID int64
}

type HubPostIDsInput struct {
	SocietyID  int64
	ChannelID  int64
	CategoryID *int16
	BeforeAt   *time.Time
	BeforeID   int64
	PageLimit  int32
}

type HubPostInput struct {
	SocietyID int64
	ID        int64
}

type HubPostRateInput struct {
	SocietyID     int64
	AuthorID      int64
	WindowSeconds int32
}

type HubPostRateRecord struct {
	Count  int64
	Oldest *time.Time
}

type HubPostRecord struct {
	ID              int64
	SocietyID       int64
	ChannelID       int64
	AuthorID        int64
	Title           *string
	Body            string
	CategoryID      *int16
	CommentsEnabled bool
	IsPinned        bool
	IsImportant     bool
	Status          string
	EditedAt        *time.Time
	RemovalReason   *string
	DeletedAt       *time.Time
	DeletedBy       *int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ChannelType     string
}

type HubPostViewInput struct {
	UserID    int64
	SocietyID int64
	ID        int64
}

type HubReadInput struct {
	UserID    int64
	SocietyID int64
	ChannelID int64
	PostID    int64
}

type HubRemoveReactionInput struct {
	UserID       int64
	PostID       *int64
	CommentID    *int64
	ReactionType string
}

type HubReportByIDInput struct {
	SocietyID int64
	ID        int64
}

type HubReportInput struct {
	PostID     int64
	ReporterID int64
	Reason     string
}

type HubReportsInput struct {
	SocietyID int64
	Status    string
	AfterID   int64
	PageLimit int32
}

type HubResolvePostReportsInput struct {
	PostID         int64
	ReviewerID     *int64
	ResolutionNote *string
}

type HubResolveReportInput struct {
	ID             int64
	Status         string
	ReviewerID     *int64
	ResolutionNote *string
}

type HubRoleInput struct {
	SocietyID int64
	UserID    int64
}

type HubUpdateCommentInput struct {
	ID     int64
	PostID int64
	Body   string
}

type HubUpdatePostInput struct {
	SocietyID       int64
	ID              int64
	Title           *string
	Body            string
	CategoryID      *int16
	CommentsEnabled bool
	IsPinned        bool
	IsImportant     bool
}

type HubUploadInput struct {
	SocietyID int64
	ID        int64
}

type SocietyChannel struct {
	ID          int64
	SocietyID   int64
	Type        string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
