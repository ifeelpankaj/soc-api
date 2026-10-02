package models

import "time"

type VerificationPurpose string

const (
	VerificationPurposeEmailVerification VerificationPurpose = "email_verification"
	VerificationPurposePhoneVerification VerificationPurpose = "phone_verification"
	VerificationPurposePasswordReset     VerificationPurpose = "password_reset"
	VerificationPurposeLoginOTP          VerificationPurpose = "login_otp"
)

func (p VerificationPurpose) IsValid() bool {
	switch p {
	case VerificationPurposeEmailVerification,
		VerificationPurposePhoneVerification,
		VerificationPurposePasswordReset,
		VerificationPurposeLoginOTP:
		return true
	default:
		return false
	}
}

type UserVerification struct {
	ID int64 `json:"id" db:"id"`

	UserID int64 `json:"user_id" db:"user_id"`

	Purpose VerificationPurpose `json:"purpose" db:"purpose"`
	Target  string              `json:"target" db:"target"`

	OTPHash string `json:"-" db:"otp_hash"`

	Attempts    int `json:"attempts" db:"attempts"`
	MaxAttempts int `json:"max_attempts" db:"max_attempts"`

	IsUsed bool `json:"is_used" db:"is_used"`

	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
