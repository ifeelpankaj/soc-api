package contracts

import (
	"context"
	"github.com/google/uuid"
	"go-server/internal/models"
)

type HubAccess interface {
	HubActorLock(ctx context.Context, arg HubActorLockInput) error
	HubResidents(ctx context.Context, societyID int64) ([]int64, error)
	HubRole(ctx context.Context, arg HubRoleInput) (string, error)
}

type HubPosts interface {
	HubCategories(ctx context.Context) ([]CommunityCategory, error)
	HubCategory(ctx context.Context, id *int16) (int16, error)
	HubChannel(ctx context.Context, arg HubChannelInput) (SocietyChannel, error)
	HubChannels(ctx context.Context, arg HubChannelsInput) ([]HubChannelsRecord, error)
	HubControlPost(ctx context.Context, arg HubControlPostInput) error
	HubCreatePost(ctx context.Context, arg HubCreatePostInput) (int64, error)
	HubDeletePost(ctx context.Context, arg HubDeletePostInput) error
	HubLockChannel(ctx context.Context, arg HubLockChannelInput) (SocietyChannel, error)
	HubLockPost(ctx context.Context, arg HubLockPostInput) (HubLockPostRecord, error)
	HubPinnedIDs(ctx context.Context, arg HubPinnedIDsInput) ([]int64, error)
	HubPost(ctx context.Context, arg HubPostInput) (HubPostRecord, error)
	HubPostIDs(ctx context.Context, arg HubPostIDsInput) ([]int64, error)
	HubPostRate(ctx context.Context, arg HubPostRateInput) (HubPostRateRecord, error)
	HubPostView(ctx context.Context, arg HubPostViewInput) (models.HubContent, error)
	HubRead(ctx context.Context, arg HubReadInput) error
	HubUpdatePost(ctx context.Context, arg HubUpdatePostInput) error
}

type HubComments interface {
	HubComment(ctx context.Context, arg HubCommentInput) (ChannelComment, error)
	HubCommentIDs(ctx context.Context, arg HubCommentIDsInput) ([]int64, error)
	HubCommentRate(ctx context.Context, arg HubCommentRateInput) (HubCommentRateRecord, error)
	HubCommentView(ctx context.Context, arg HubCommentViewInput) (models.HubContent, error)
	HubCreateComment(ctx context.Context, arg HubCreateCommentInput) (int64, error)
	HubDeleteComment(ctx context.Context, arg HubDeleteCommentInput) error
	HubUpdateComment(ctx context.Context, arg HubUpdateCommentInput) error
}

type HubAttachments interface {
	HubAttach(ctx context.Context, arg HubAttachInput) error
	HubAttachmentTarget(ctx context.Context, uploadID int64) (HubAttachmentTargetRecord, error)
	HubAttachments(ctx context.Context, arg HubAttachmentsInput) ([]int64, error)
	HubClaimCleanup(ctx context.Context, leaseToken uuid.UUID) (ChannelUpload, error)
	HubClaimUpload(ctx context.Context, arg HubClaimUploadInput) (int64, error)
	HubCreateUpload(ctx context.Context, arg HubCreateUploadInput) (ChannelUpload, error)
	HubDetach(ctx context.Context, uploadID int64) error
	HubFinishCleanup(ctx context.Context, arg HubFinishCleanupInput) error
	HubLockUpload(ctx context.Context, arg HubLockUploadInput) (ChannelUpload, error)
	HubQueueUploadCleanup(ctx context.Context, id int64) error
	HubUpload(ctx context.Context, arg HubUploadInput) (ChannelUpload, error)
	HubUploadByFile(ctx context.Context, fileID string) (ChannelUpload, error)
}

type HubModeration interface {
	HubReport(ctx context.Context, arg HubReportInput) (ChannelReport, error)
	HubReportByID(ctx context.Context, arg HubReportByIDInput) (ChannelReport, error)
	HubReports(ctx context.Context, arg HubReportsInput) ([]ChannelReport, error)
	HubResolvePostReports(ctx context.Context, arg HubResolvePostReportsInput) error
	HubResolveReport(ctx context.Context, arg HubResolveReportInput) error
}

type HubReactions interface {
	HubAddReaction(ctx context.Context, arg HubAddReactionInput) (int64, error)
	HubRemoveReaction(ctx context.Context, arg HubRemoveReactionInput) error
}

type HubNotifications interface {
	HubEnqueue(ctx context.Context, arg HubEnqueueInput) error
	HubQueueAnnouncementFanout(ctx context.Context, societyID, postID, actorID int64, important bool) error
}

// HubStore composes independently mockable persistence capabilities.
type HubStore interface {
	HubAccess
	HubPosts
	HubComments
	HubAttachments
	HubModeration
	HubReactions
	HubNotifications
	WithTransaction(context.Context, func(context.Context) error) error
}
