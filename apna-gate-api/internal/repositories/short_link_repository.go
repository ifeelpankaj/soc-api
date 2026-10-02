package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"time"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
)

type ShortLinkRepository = contracts.ShortLinkRepository

type shortLinkRepository struct {
	db *database.Database
}

func NewShortLinkRepository(db *database.Database) ShortLinkRepository {
	return &shortLinkRepository{db: db}
}

func (r *shortLinkRepository) Create(ctx context.Context, code string, resourceType models.ShortLinkResourceType, resourceID int64, expiresAt *time.Time, createdBy *int64, metadata map[string]any) (*models.ShortLink, error) {
	rawMetadata, err := jsonMap(metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).CreateShortLink(ctx, db.CreateShortLinkParams{
		ShortCode:    code,
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
		ExpiresAt:    timePtrToPgTimestamptz(expiresAt),
		CreatedBy:    createdBy,
		Metadata:     rawMetadata,
	})
	return shortLinkFromDBNoRows(row, persistenceError(err))
}

func (r *shortLinkRepository) GetByCode(ctx context.Context, code string) (*models.ShortLink, error) {
	row, err := GetQueries(ctx, r.db).GetShortLinkByCode(ctx, code)
	return shortLinkFromDBNoRows(row, err)
}

func (r *shortLinkRepository) GetByResource(ctx context.Context, resourceType models.ShortLinkResourceType, resourceID int64) (*models.ShortLink, error) {
	row, err := GetQueries(ctx, r.db).GetShortLinkByResource(ctx, db.GetShortLinkByResourceParams{
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
	})
	return shortLinkFromDBNoRows(row, err)
}

func (r *shortLinkRepository) Revoke(ctx context.Context, code string) (*models.ShortLink, error) {
	row, err := GetQueries(ctx, r.db).RevokeShortLink(ctx, code)
	return shortLinkFromDBNoRows(row, err)
}

func shortLinkFromDBNoRows(row db.ShortLink, err error) (*models.ShortLink, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return shortLinkFromDB(row), nil
}

func shortLinkFromDB(row db.ShortLink) *models.ShortLink {
	return &models.ShortLink{
		ID:           row.ID,
		ShortCode:    row.ShortCode,
		ResourceType: models.ShortLinkResourceType(row.ResourceType),
		ResourceID:   row.ResourceID,
		ExpiresAt:    pgTimestamptzToTimePtr(row.ExpiresAt),
		RevokedAt:    pgTimestamptzToTimePtr(row.RevokedAt),
		CreatedBy:    row.CreatedBy,
		Metadata:     metadataFromJSON(row.Metadata),
		CreatedAt:    pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:    pgTimestamptzToTime(row.UpdatedAt),
	}
}
