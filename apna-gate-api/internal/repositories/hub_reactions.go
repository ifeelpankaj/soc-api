package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubAddReaction(ctx context.Context, arg contracts.HubAddReactionInput) (int64, error) {
	value, err := GetQueries(ctx, r.database).HubAddReaction(ctx, db.HubAddReactionParams{UserID: arg.UserID, PostID: arg.PostID, CommentID: arg.CommentID, ReactionType: arg.ReactionType})
	return value, persistenceError(err)
}

func (r *HubRepository) HubRemoveReaction(ctx context.Context, arg contracts.HubRemoveReactionInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubRemoveReaction(ctx, db.HubRemoveReactionParams{UserID: arg.UserID, PostID: arg.PostID, CommentID: arg.CommentID, ReactionType: arg.ReactionType}))
}
