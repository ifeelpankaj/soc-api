package contracts

import (
	"context"
	"go-server/internal/models"
)

type VerificationRepository interface {
	CreateVerification(ctx context.Context, verification *models.UserVerification) error
	GetActiveVerification(ctx context.Context, userID int64, purpose models.VerificationPurpose, target string) (*models.UserVerification, error)
	MarkAsUsed(ctx context.Context, verificationID int64) error
	IncrementAttempts(ctx context.Context, verificationID int64) error
	DeleteActiveByPurpose(ctx context.Context, userID int64, purpose models.VerificationPurpose, target string) error
	DeleteUsedOrExpired(ctx context.Context) error
}
