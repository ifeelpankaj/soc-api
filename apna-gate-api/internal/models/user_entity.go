package models

import "time"

type GlobalRole string

const (
	GlobalRoleUser       GlobalRole = "user"
	GlobalRoleDeveloper  GlobalRole = "developer"
	GlobalRoleSuperAdmin GlobalRole = "super_admin"
)

func (r GlobalRole) IsValid() bool {
	switch r {
	case GlobalRoleUser, GlobalRoleDeveloper, GlobalRoleSuperAdmin:
		return true
	default:
		return false
	}
}

type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderApple  AuthProvider = "apple"
	AuthProviderPhone  AuthProvider = "phone"
)

func (p AuthProvider) IsValid() bool {
	switch p {
	case AuthProviderEmail, AuthProviderGoogle, AuthProviderApple, AuthProviderPhone:
		return true
	default:
		return false
	}
}

type User struct {
	SessionVersion int64 `json:"-" db:"session_version"`
	ID             int64 `json:"id" db:"id"`

	FirstName *string `json:"first_name,omitempty" db:"first_name"`
	LastName  *string `json:"last_name,omitempty" db:"last_name"`
	FullName  string  `json:"full_name" db:"full_name"`

	Email        *string `json:"email,omitempty" db:"email"`
	PhoneNumber  *string `json:"phone_number,omitempty" db:"phone_number"`
	PasswordHash *string `json:"-" db:"password_hash"`

	AuthProvider AuthProvider `json:"auth_provider" db:"auth_provider"`
	ProviderID   *string      `json:"provider_id,omitempty" db:"provider_id"`

	GlobalRole GlobalRole `json:"global_role" db:"global_role"`

	EmailVerified bool `json:"email_verified" db:"email_verified"`
	PhoneVerified bool `json:"phone_verified" db:"phone_verified"`

	IsActive      bool    `json:"is_active" db:"is_active"`
	IsBlocked     bool    `json:"is_blocked" db:"is_blocked"`
	BlockedReason *string `json:"blocked_reason,omitempty" db:"blocked_reason"`

	AvatarURL   *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender      *string    `json:"gender,omitempty" db:"gender"`

	Timezone string `json:"timezone" db:"timezone"`
	Language string `json:"language" db:"language"`

	LastLoginAt       *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" db:"password_changed_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	Metadata map[string]any `json:"metadata" db:"metadata"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) EmailValue() string {
	if u == nil || u.Email == nil {
		return ""
	}
	return *u.Email
}

func (u *User) PhoneNumberValue() string {
	if u == nil || u.PhoneNumber == nil {
		return ""
	}
	return *u.PhoneNumber
}

func (u *User) ToResponse() *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		FullName:      u.FullName,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		AuthProvider:  u.AuthProvider,
		GlobalRole:    u.GlobalRole,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		IsActive:      u.IsActive,
		IsBlocked:     u.IsBlocked,
		BlockedReason: u.BlockedReason,
		AvatarURL:     u.AvatarURL,
		DateOfBirth:   NewDateOnly(u.DateOfBirth),
		Gender:        u.Gender,
		Timezone:      u.Timezone,
		Language:      u.Language,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
