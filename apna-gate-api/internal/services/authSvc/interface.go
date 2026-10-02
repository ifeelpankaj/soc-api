package authsvc

import "context"

type EmailService interface {
	SendOTP(ctx context.Context, to, otp, name string) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendForgetPasswordEmail(ctx context.Context, to, otp, name string) error
	Close() error
}
