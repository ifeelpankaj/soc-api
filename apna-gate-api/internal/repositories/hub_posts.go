package repository

import (
	"context"
	"encoding/json"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubCategories(ctx context.Context) ([]contracts.CommunityCategory, error) {
	value, err := GetQueries(ctx, r.database).HubCategories(ctx)
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.CommunityCategory, len(value))
	for i := range value {
		out[i] = mapCommunityCategory(value[i])
	}
	return out, nil
}

func (r *HubRepository) HubCategory(ctx context.Context, id *int16) (int16, error) {
	value, err := GetQueries(ctx, r.database).HubCategory(ctx, id)
	return value, persistenceError(err)
}

func (r *HubRepository) HubChannel(ctx context.Context, arg contracts.HubChannelInput) (contracts.SocietyChannel, error) {
	value, err := GetQueries(ctx, r.database).HubChannel(ctx, db.HubChannelParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapSocietyChannel(value), persistenceError(err)
}

func (r *HubRepository) HubChannels(ctx context.Context, arg contracts.HubChannelsInput) ([]contracts.HubChannelsRecord, error) {
	value, err := GetQueries(ctx, r.database).HubChannels(ctx, db.HubChannelsParams{UserID: arg.UserID, SocietyID: arg.SocietyID})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.HubChannelsRecord, len(value))
	for i := range value {
		out[i] = mapHubChannelsRow(value[i])
	}
	return out, nil
}

func (r *HubRepository) HubControlPost(ctx context.Context, arg contracts.HubControlPostInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubControlPost(ctx, db.HubControlPostParams{SocietyID: arg.SocietyID, ID: arg.ID, IsPinned: arg.IsPinned, CommentsEnabled: arg.CommentsEnabled}))
}

func (r *HubRepository) HubCreatePost(ctx context.Context, arg contracts.HubCreatePostInput) (int64, error) {
	value, err := GetQueries(ctx, r.database).HubCreatePost(ctx, db.HubCreatePostParams{SocietyID: arg.SocietyID, ChannelID: arg.ChannelID, AuthorID: arg.AuthorID, Title: arg.Title, Body: arg.Body, CategoryID: arg.CategoryID, CommentsEnabled: arg.CommentsEnabled, IsPinned: arg.IsPinned, IsImportant: arg.IsImportant})
	return value, persistenceError(err)
}

func (r *HubRepository) HubDeletePost(ctx context.Context, arg contracts.HubDeletePostInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubDeletePost(ctx, db.HubDeletePostParams{SocietyID: arg.SocietyID, ID: arg.ID, Status: arg.Status, DeletedBy: arg.DeletedBy, RemovalReason: arg.RemovalReason}))
}

func (r *HubRepository) HubLockChannel(ctx context.Context, arg contracts.HubLockChannelInput) (contracts.SocietyChannel, error) {
	value, err := GetQueries(ctx, r.database).HubLockChannel(ctx, db.HubLockChannelParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapSocietyChannel(value), persistenceError(err)
}

func (r *HubRepository) HubLockPost(ctx context.Context, arg contracts.HubLockPostInput) (contracts.HubLockPostRecord, error) {
	value, err := GetQueries(ctx, r.database).HubLockPost(ctx, db.HubLockPostParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapHubLockPostRow(value), persistenceError(err)
}

func (r *HubRepository) HubPinnedIDs(ctx context.Context, arg contracts.HubPinnedIDsInput) ([]int64, error) {
	value, err := GetQueries(ctx, r.database).HubPinnedIDs(ctx, db.HubPinnedIDsParams{SocietyID: arg.SocietyID, ChannelID: arg.ChannelID})
	return value, persistenceError(err)
}

func (r *HubRepository) HubPost(ctx context.Context, arg contracts.HubPostInput) (contracts.HubPostRecord, error) {
	value, err := GetQueries(ctx, r.database).HubPost(ctx, db.HubPostParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapHubPostRow(value), persistenceError(err)
}

func (r *HubRepository) HubPostIDs(ctx context.Context, arg contracts.HubPostIDsInput) ([]int64, error) {
	value, err := GetQueries(ctx, r.database).HubPostIDs(ctx, db.HubPostIDsParams{SocietyID: arg.SocietyID, ChannelID: arg.ChannelID, CategoryID: arg.CategoryID, BeforeAt: timePtrToPgTimestamptz(arg.BeforeAt), BeforeID: arg.BeforeID, PageLimit: arg.PageLimit})
	return value, persistenceError(err)
}

func (r *HubRepository) HubPostRate(ctx context.Context, arg contracts.HubPostRateInput) (contracts.HubPostRateRecord, error) {
	value, err := GetQueries(ctx, r.database).HubPostRate(ctx, db.HubPostRateParams{SocietyID: arg.SocietyID, AuthorID: arg.AuthorID, WindowSeconds: arg.WindowSeconds})
	return mapHubPostRateRow(value), persistenceError(err)
}

func (r *HubRepository) HubPostView(ctx context.Context, arg contracts.HubPostViewInput) (models.HubContent, error) {
	value, err := GetQueries(ctx, r.database).HubPostView(ctx, db.HubPostViewParams{UserID: arg.UserID, SocietyID: arg.SocietyID, ID: arg.ID})
	var result models.HubContent
	if err != nil {
		return result, persistenceError(err)
	}
	err = json.Unmarshal(value, &result)
	return result, err
}

func (r *HubRepository) HubRead(ctx context.Context, arg contracts.HubReadInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubRead(ctx, db.HubReadParams{UserID: arg.UserID, SocietyID: arg.SocietyID, ChannelID: arg.ChannelID, PostID: arg.PostID}))
}

func (r *HubRepository) HubUpdatePost(ctx context.Context, arg contracts.HubUpdatePostInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubUpdatePost(ctx, db.HubUpdatePostParams{SocietyID: arg.SocietyID, ID: arg.ID, Title: arg.Title, Body: arg.Body, CategoryID: arg.CategoryID, CommentsEnabled: arg.CommentsEnabled, IsPinned: arg.IsPinned, IsImportant: arg.IsImportant}))
}
