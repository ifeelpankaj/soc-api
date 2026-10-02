package models

import "strings"

type LoginRequest struct {
	Email       string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,phone_intl,max=20"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

func (r *LoginRequest) Sanitize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
}

type RegisterRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=100,alphanumeric_space"`
	LastName    string `json:"last_name" validate:"omitempty,max=100,alphanumeric_space"`
	Email       string `json:"email" validate:"required,email,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,phone_intl,max=20"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

func (r *RegisterRequest) Sanitize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
}

type ResidentRegisterRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=100,alphanumeric_space"`
	LastName    string `json:"last_name" validate:"omitempty,max=100,alphanumeric_space"`
	Email       string `json:"email" validate:"required,email,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,phone_intl,max=20"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

func (r *ResidentRegisterRequest) Sanitize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
	OTP   string `json:"otp" validate:"required,len=6,numeric"`
}

func (r *VerifyOTPRequest) Sanitize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.OTP = strings.TrimSpace(r.OTP)
}

type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (r *ResendOTPRequest) Sanitize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (r *ForgotPasswordRequest) Sanitize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email,max=255"`
	OTP             string `json:"otp" validate:"required,len=6,numeric"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=72"`
}

func (r *ResetPasswordRequest) Sanitize() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.OTP = strings.TrimSpace(r.OTP)
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=72"`
}

func (r *ChangePasswordRequest) Sanitize() {}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
