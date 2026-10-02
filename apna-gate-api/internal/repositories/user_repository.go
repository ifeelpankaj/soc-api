package repository

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository = contracts.UserRepository

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	metadata, err := json.Marshal(user.Metadata)
	if err != nil {
		return err
	}
	if string(metadata) == "null" {
		metadata = []byte("{}")
	}

	row, err := GetQueries(ctx, r.db).CreateUser(ctx, db.CreateUserParams{
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.FullName,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		PasswordHash:  user.PasswordHash,
		AuthProvider:  db.AuthProvider(user.AuthProvider),
		GlobalRole:    db.GlobalRole(user.GlobalRole),
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		IsActive:      user.IsActive,
		IsBlocked:     user.IsBlocked,
		Timezone:      user.Timezone,
		Language:      user.Language,
		Metadata:      metadata,
	})
	if err != nil {
		return err
	}

	*user = *userFromDB(row)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row, err := GetQueries(ctx, r.db).GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return userFromDB(row), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := GetQueries(ctx, r.db).GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return userFromDB(row), nil
}

func (r *userRepository) GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error) {
	row, err := GetQueries(ctx, r.db).GetUserByPhoneNumber(ctx, &phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return userFromDB(row), nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return GetQueries(ctx, r.db).EmailExists(ctx, email)
}

func (r *userRepository) PhoneExists(ctx context.Context, phone string) (bool, error) {
	return GetQueries(ctx, r.db).PhoneExists(ctx, &phone)
}

func userFromDB(row db.User) *models.User {
	metadata := map[string]any{}
	if len(row.Metadata) > 0 {
		_ = json.Unmarshal(row.Metadata, &metadata)
	}

	return &models.User{
		ID:                row.ID,
		FirstName:         row.FirstName,
		LastName:          row.LastName,
		FullName:          row.FullName,
		Email:             row.Email,
		PhoneNumber:       row.PhoneNumber,
		PasswordHash:      row.PasswordHash,
		SessionVersion:    row.SessionVersion,
		AuthProvider:      models.AuthProvider(row.AuthProvider),
		ProviderID:        row.ProviderID,
		GlobalRole:        models.GlobalRole(row.GlobalRole),
		EmailVerified:     row.EmailVerified,
		PhoneVerified:     row.PhoneVerified,
		IsActive:          row.IsActive,
		IsBlocked:         row.IsBlocked,
		BlockedReason:     row.BlockedReason,
		AvatarURL:         row.AvatarUrl,
		DateOfBirth:       pgDateToTimePtr(row.DateOfBirth),
		Gender:            row.Gender,
		Timezone:          row.Timezone,
		Language:          row.Language,
		LastLoginAt:       pgTimestamptzToTimePtr(row.LastLoginAt),
		PasswordChangedAt: pgTimestamptzToTimePtr(row.PasswordChangedAt),
		DeletedAt:         pgTimestamptzToTimePtr(row.DeletedAt),
		Metadata:          metadata,
		CreatedAt:         pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:         pgTimestamptzToTime(row.UpdatedAt),
	}
}

func pgDateToTimePtr(date pgtype.Date) *time.Time {
	if !date.Valid {
		return nil
	}
	value := date.Time
	return &value
}

func pgTimestamptzToTimePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}

func pgTimestamptzToTime(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}
func (r *userRepository) MarkEmailVerified(ctx context.Context, userID int64) error {
	return GetQueries(ctx, r.db).MarkUserEmailVerified(ctx, userID)
}

func (r *userRepository) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string, version int64) error {
	rows, err := GetQueries(ctx, r.db).UpdateUserPasswordHash(ctx, db.UpdateUserPasswordHashParams{
		ID:             userID,
		PasswordHash:   &passwordHash,
		SessionVersion: version,
	})
	if err != nil {
		return err
	}
	if rows != 1 {
		return models.NewAppError("PASSWORD_CHANGE_CONFLICT", "Your password changed during this request. Please sign in again.", 409, nil)
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	return GetQueries(ctx, r.db).UpdateUserLastLogin(ctx, userID)
}

func (r *userRepository) UpdateProfile(ctx context.Context, userID int64, req *contracts.UpdateUserInput) (*models.User, error) {
	row, err := GetQueries(ctx, r.db).UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:          userID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: timePtrToPgDate(req.DateOfBirth),
		Gender:      req.Gender,
		Timezone:    req.Timezone,
		Language:    req.Language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return userFromUpdateProfileRow(row), nil
}

func userFromUpdateProfileRow(row db.UpdateUserProfileRow) *models.User {
	metadata := map[string]any{}
	if len(row.Metadata) > 0 {
		_ = json.Unmarshal(row.Metadata, &metadata)
	}
	return &models.User{
		ID: row.ID, FirstName: row.FirstName, LastName: row.LastName, FullName: row.FullName,
		Email: row.Email, PhoneNumber: row.PhoneNumber, PasswordHash: row.PasswordHash, SessionVersion: row.SessionVersion,
		AuthProvider: models.AuthProvider(row.AuthProvider), ProviderID: row.ProviderID,
		GlobalRole: models.GlobalRole(row.GlobalRole), EmailVerified: row.EmailVerified,
		PhoneVerified: row.PhoneVerified, IsActive: row.IsActive, IsBlocked: row.IsBlocked,
		AvatarURL: row.AvatarUrl, DateOfBirth: pgDateToTimePtr(row.DateOfBirth), Gender: row.Gender,
		Timezone: row.Timezone, Language: row.Language, LastLoginAt: pgTimestamptzToTimePtr(row.LastLoginAt),
		PasswordChangedAt: pgTimestamptzToTimePtr(row.PasswordChangedAt), DeletedAt: pgTimestamptzToTimePtr(row.DeletedAt),
		Metadata: metadata, CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func timePtrToPgDate(value *time.Time) pgtype.Date {
	if value == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *value, Valid: true}
}
