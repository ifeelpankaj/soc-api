package repository

import (
	"context"
	"errors"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerificationRepository = contracts.VerificationRepository

type verificationRepository struct {
	db *database.Database
}

func NewVerificationRepository(db *database.Database) VerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) CreateVerification(ctx context.Context, verification *models.UserVerification) error {
	row, err := GetQueries(ctx, r.db).CreateUserVerification(ctx, db.CreateUserVerificationParams{
		UserID:      verification.UserID,
		Purpose:     db.VerificationPurpose(verification.Purpose),
		Target:      verification.Target,
		OtpHash:     verification.OTPHash,
		Attempts:    int32(verification.Attempts),    //nolint:gosec // Attempts are bounded by validation before persistence.
		MaxAttempts: int32(verification.MaxAttempts), //nolint:gosec // Max attempts are bounded by validation before persistence.
		IsUsed:      verification.IsUsed,
		ExpiresAt:   timeToPgTimestamptz(verification.ExpiresAt),
	})
	if err != nil {
		return err
	}

	*verification = *verificationFromDB(row)
	return nil
}

func (r *verificationRepository) GetActiveVerification(ctx context.Context, userID int64, purpose models.VerificationPurpose, target string) (*models.UserVerification, error) {
	row, err := GetQueries(ctx, r.db).GetActiveUserVerification(ctx, db.GetActiveUserVerificationParams{
		UserID:  userID,
		Purpose: db.VerificationPurpose(purpose),
		Target:  target,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return verificationFromDB(row), nil
}

func (r *verificationRepository) MarkAsUsed(ctx context.Context, verificationID int64) error {
	return GetQueries(ctx, r.db).MarkUserVerificationUsed(ctx, verificationID)
}

func (r *verificationRepository) IncrementAttempts(ctx context.Context, verificationID int64) error {
	return GetQueries(ctx, r.db).IncrementUserVerificationAttempts(ctx, verificationID)
}

func (r *verificationRepository) DeleteActiveByPurpose(ctx context.Context, userID int64, purpose models.VerificationPurpose, target string) error {
	return GetQueries(ctx, r.db).DeleteActiveUserVerificationByPurpose(ctx, db.DeleteActiveUserVerificationByPurposeParams{
		UserID:  userID,
		Purpose: db.VerificationPurpose(purpose),
		Target:  target,
	})
}

func (r *verificationRepository) DeleteUsedOrExpired(ctx context.Context) error {
	return GetQueries(ctx, r.db).DeleteUsedOrExpiredUserVerifications(ctx)
}

func verificationFromDB(row db.UserVerification) *models.UserVerification {
	return &models.UserVerification{
		ID:          row.ID,
		UserID:      row.UserID,
		Purpose:     models.VerificationPurpose(row.Purpose),
		Target:      row.Target,
		OTPHash:     row.OtpHash,
		Attempts:    int(row.Attempts),
		MaxAttempts: int(row.MaxAttempts),
		IsUsed:      row.IsUsed,
		ExpiresAt:   pgTimestamptzToTime(row.ExpiresAt),
		UsedAt:      pgTimestamptzToTimePtr(row.UsedAt),
		CreatedAt:   pgTimestamptzToTime(row.CreatedAt),
	}
}

func timeToPgTimestamptz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: !value.IsZero()}
}
