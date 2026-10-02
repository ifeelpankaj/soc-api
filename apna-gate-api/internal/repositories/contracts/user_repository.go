package contracts

import (
	"context"
	"go-server/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	PhoneExists(ctx context.Context, phone string) (bool, error)
	MarkEmailVerified(ctx context.Context, userID int64) error
	UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string, version int64) error
	UpdateLastLogin(ctx context.Context, userID int64) error
	UpdateProfile(ctx context.Context, userID int64, req *UpdateUserInput) (*models.User, error)
}
