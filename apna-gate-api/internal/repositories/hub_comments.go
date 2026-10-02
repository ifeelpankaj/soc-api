package repository

import (
	"context"
	"encoding/json"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubComment(ctx context.Context, arg contracts.HubCommentInput) (contracts.ChannelComment, error) {
	value, err := GetQueries(ctx, r.database).HubComment(ctx, db.HubCommentParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapChannelComment(value), persistenceError(err)
}

func (r *HubRepository) HubCommentIDs(ctx context.Context, arg contracts.HubCommentIDsInput) ([]int64, error) {
	value, err := GetQueries(ctx, r.database).HubCommentIDs(ctx, db.HubCommentIDsParams{SocietyID: arg.SocietyID, PostID: arg.PostID, AfterAt: timePtrToPgTimestamptz(arg.AfterAt), AfterID: arg.AfterID, PageLimit: arg.PageLimit})
	return value, persistenceError(err)
}

func (r *HubRepository) HubCommentRate(ctx context.Context, arg contracts.HubCommentRateInput) (contracts.HubCommentRateRecord, error) {
	value, err := GetQueries(ctx, r.database).HubCommentRate(ctx, db.HubCommentRateParams{SocietyID: arg.SocietyID, AuthorID: arg.AuthorID, WindowSeconds: arg.WindowSeconds})
	return mapHubCommentRateRow(value), persistenceError(err)
}

func (r *HubRepository) HubCommentView(ctx context.Context, arg contracts.HubCommentViewInput) (models.HubContent, error) {
	value, err := GetQueries(ctx, r.database).HubCommentView(ctx, db.HubCommentViewParams{UserID: arg.UserID, SocietyID: arg.SocietyID, ID: arg.ID})
	var result models.HubContent
	if err != nil {
		return result, persistenceError(err)
	}
	err = json.Unmarshal(value, &result)
	return result, err
}

func (r *HubRepository) HubCreateComment(ctx context.Context, arg contracts.HubCreateCommentInput) (int64, error) {
	value, err := GetQueries(ctx, r.database).HubCreateComment(ctx, db.HubCreateCommentParams{PostID: arg.PostID, AuthorID: arg.AuthorID, ParentID: arg.ParentID, Body: arg.Body})
	return value, persistenceError(err)
}

func (r *HubRepository) HubDeleteComment(ctx context.Context, arg contracts.HubDeleteCommentInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubDeleteComment(ctx, db.HubDeleteCommentParams{ID: arg.ID, PostID: arg.PostID, Status: arg.Status, DeletedBy: arg.DeletedBy, RemovalReason: arg.RemovalReason}))
}

func (r *HubRepository) HubUpdateComment(ctx context.Context, arg contracts.HubUpdateCommentInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubUpdateComment(ctx, db.HubUpdateCommentParams{ID: arg.ID, PostID: arg.PostID, Body: arg.Body}))
}
