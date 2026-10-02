package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
)

// ImageCommitUncertain deliberately carries no driver text: callers must resolve
// the reference before compensation, and must not expose connection details.
type ImageCommitUncertain = contracts.ImageCommitUncertain

type ImageRepository = contracts.ImageRepository

type imageRepository struct{ database *database.Database }

func NewImageRepository(database *database.Database) ImageRepository {
	return &imageRepository{database: database}
}

func imageValue(url, id, path *string) models.StoredImage {
	var result models.StoredImage
	if url != nil {
		result.URL = *url
	}
	if id != nil {
		result.FileID = *id
	}
	if path != nil {
		result.Path = *path
	}
	return result
}
func imageReadError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return contracts.ErrNotFound
	}
	return err
}

func (r *imageRepository) Read(ctx context.Context, target models.ImageTarget) (models.StoredImage, error) {
	q := GetQueries(ctx, r.database)
	if target.Avatar() {
		row, err := q.GetUserImage(ctx, target.ActorID)
		return imageValue(row.AvatarUrl, row.AvatarImagekitFileID, row.AvatarImagekitFilePath), imageReadError(err)
	}
	row, err := q.GetVisitorImage(ctx, target.VisitorID)
	return imageValue(row.PhotoUrl, row.PhotoImagekitFileID, row.PhotoImagekitFilePath), imageReadError(err)
}

func lockImage(ctx context.Context, q *db.Queries, target models.ImageTarget) (models.StoredImage, error) {
	if target.Avatar() {
		row, err := q.LockUserImage(ctx, target.ActorID)
		return imageValue(row.AvatarUrl, row.AvatarImagekitFileID, row.AvatarImagekitFilePath), imageReadError(err)
	}
	row, err := q.LockVisitorImage(ctx, target.VisitorID)
	return imageValue(row.PhotoUrl, row.PhotoImagekitFileID, row.PhotoImagekitFilePath), imageReadError(err)
}

func rollbackImage(tx pgx.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), models.ImageDatabaseTimeout)
	defer cancel()
	_ = tx.Rollback(ctx)
}

func (r *imageRepository) Replace(ctx context.Context, target models.ImageTarget, next models.StoredImage, authorize func(context.Context) error) (models.StoredImage, error) {
	tx, err := r.database.Pool.Begin(ctx)
	if err != nil {
		return models.StoredImage{}, err
	}
	defer rollbackImage(tx)
	txCtx := context.WithValue(ctx, txKey{}, tx)
	q := db.New(tx)
	// Always lock the entry before its visitor, including photo deletion.
	if !target.Avatar() {
		entry, lockErr := q.GetVisitorEntryForUpdate(txCtx, db.GetVisitorEntryForUpdateParams{SocietyID: target.SocietyID, ID: target.EntryID})
		if lockErr != nil {
			return models.StoredImage{}, imageReadError(lockErr)
		}
		if entry.VisitorID != target.VisitorID {
			return models.StoredImage{}, contracts.ErrNotFound
		}
	}
	old, err := lockImage(txCtx, q, target)
	if err != nil {
		return old, err
	}
	if err = authorize(txCtx); err != nil {
		return old, err
	}
	var url, id, path *string
	if next.Managed() {
		url, id, path = &next.URL, &next.FileID, &next.Path
	}
	var count int64
	if target.Avatar() {
		count, err = q.SetUserImage(txCtx, db.SetUserImageParams{ID: target.ActorID, Url: url, FileID: id, FilePath: path})
	} else {
		count, err = q.SetVisitorImage(txCtx, db.SetVisitorImageParams{ID: target.VisitorID, Url: url, FileID: id, FilePath: path})
	}
	if err != nil {
		return old, err
	}
	if count != 1 {
		return old, contracts.ErrNotFound
	}
	if err = tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return old, contracts.ErrConflict
		}
		return old, &ImageCommitUncertain{}
	}
	return old, nil
}

// Resolve uses an owner row lock to wait for any in-flight commit to settle.
// It never authorizes signing and never runs network storage operations.
func (r *imageRepository) Resolve(ctx context.Context, target models.ImageTarget) (models.StoredImage, error) {
	tx, err := r.database.Pool.Begin(ctx)
	if err != nil {
		return models.StoredImage{}, err
	}
	defer rollbackImage(tx)
	current, err := lockImage(ctx, db.New(tx), target)
	if errors.Is(err, contracts.ErrNotFound) {
		return models.StoredImage{}, nil
	}
	return current, err
}
