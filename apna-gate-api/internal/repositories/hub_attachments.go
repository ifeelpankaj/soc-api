package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubAttach(ctx context.Context, arg contracts.HubAttachInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubAttach(ctx, db.HubAttachParams{UploadID: arg.UploadID, PostID: arg.PostID, CommentID: arg.CommentID}))
}

func (r *HubRepository) HubAttachmentTarget(ctx context.Context, uploadID int64) (contracts.HubAttachmentTargetRecord, error) {
	value, err := GetQueries(ctx, r.database).HubAttachmentTarget(ctx, uploadID)
	return mapHubAttachmentTargetRow(value), persistenceError(err)
}

func (r *HubRepository) HubAttachments(ctx context.Context, arg contracts.HubAttachmentsInput) ([]int64, error) {
	value, err := GetQueries(ctx, r.database).HubAttachments(ctx, db.HubAttachmentsParams{PostID: arg.PostID, CommentID: arg.CommentID})
	return value, persistenceError(err)
}

func (r *HubRepository) HubClaimCleanup(ctx context.Context, leaseToken uuid.UUID) (contracts.ChannelUpload, error) {
	value, err := GetQueries(ctx, r.database).HubClaimCleanup(ctx, pgtype.UUID{Bytes: leaseToken, Valid: leaseToken != uuid.Nil})
	return mapChannelUpload(value), persistenceError(err)
}

func (r *HubRepository) HubClaimUpload(ctx context.Context, arg contracts.HubClaimUploadInput) (int64, error) {
	value, err := GetQueries(ctx, r.database).HubClaimUpload(ctx, db.HubClaimUploadParams{ID: arg.ID, SocietyID: arg.SocietyID, UploaderID: arg.UploaderID})
	return value, persistenceError(err)
}

func (r *HubRepository) HubCreateUpload(ctx context.Context, arg contracts.HubCreateUploadInput) (contracts.ChannelUpload, error) {
	value, err := GetQueries(ctx, r.database).HubCreateUpload(ctx, db.HubCreateUploadParams{SocietyID: arg.SocietyID, UploaderID: arg.UploaderID, FileID: arg.FileID, FilePath: arg.FilePath, OriginalFilename: arg.OriginalFilename, MimeType: arg.MimeType, FileSize: arg.FileSize, Width: arg.Width, Height: arg.Height})
	return mapChannelUpload(value), persistenceError(err)
}

func (r *HubRepository) HubDetach(ctx context.Context, uploadID int64) error {
	return persistenceError(GetQueries(ctx, r.database).HubDetach(ctx, uploadID))
}

func (r *HubRepository) HubFinishCleanup(ctx context.Context, arg contracts.HubFinishCleanupInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubFinishCleanup(ctx, db.HubFinishCleanupParams{ID: arg.ID, LeaseToken: pgtype.UUID{Bytes: arg.LeaseToken, Valid: arg.LeaseToken != uuid.Nil}}))
}

func (r *HubRepository) HubLockUpload(ctx context.Context, arg contracts.HubLockUploadInput) (contracts.ChannelUpload, error) {
	value, err := GetQueries(ctx, r.database).HubLockUpload(ctx, db.HubLockUploadParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapChannelUpload(value), persistenceError(err)
}

func (r *HubRepository) HubQueueUploadCleanup(ctx context.Context, id int64) error {
	return persistenceError(GetQueries(ctx, r.database).HubQueueUploadCleanup(ctx, id))
}

func (r *HubRepository) HubUpload(ctx context.Context, arg contracts.HubUploadInput) (contracts.ChannelUpload, error) {
	value, err := GetQueries(ctx, r.database).HubUpload(ctx, db.HubUploadParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapChannelUpload(value), persistenceError(err)
}

func (r *HubRepository) HubUploadByFile(ctx context.Context, fileID string) (contracts.ChannelUpload, error) {
	value, err := GetQueries(ctx, r.database).HubUploadByFile(ctx, fileID)
	return mapChannelUpload(value), persistenceError(err)
}
