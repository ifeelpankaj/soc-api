package models

import "time"

type UserResponse struct {
	ID int64 `json:"id"`

	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	FullName  string  `json:"full_name"`

	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`

	AuthProvider AuthProvider `json:"auth_provider"`
	GlobalRole   GlobalRole   `json:"global_role"`

	EmailVerified bool `json:"email_verified"`
	PhoneVerified bool `json:"phone_verified"`

	IsActive      bool    `json:"is_active"`
	IsBlocked     bool    `json:"is_blocked"`
	BlockedReason *string `json:"blocked_reason,omitempty"`

	AvatarURL   *string   `json:"avatar_url,omitempty"`
	DateOfBirth *DateOnly `json:"date_of_birth,omitempty"`
	Gender      *string   `json:"gender,omitempty"`

	Timezone string `json:"timezone"`
	Language string `json:"language"`

	LastLoginAt *time.Time `json:"last_login_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthUserResponse struct {
	User *UserResponse `json:"user"`
}
