package authsvc

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode/utf8"

	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	service "go-server/internal/services"
	"go-server/pkg/logger"
	"go-server/pkg/utils"

	"go.uber.org/zap"
)

const ForgotPasswordSuccessMessage = "if this email exists, password reset instructions have been sent" //nolint:gosec // Generic public response text, not a credential.

type PasswordSvc interface {
	ForgotPassword(ctx context.Context, req *models.ForgotPasswordRequest) (*ForgotPasswordResult, error)
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID int64, req *models.ChangePasswordRequest) error
}

type ForgotPasswordResult struct {
	Message string `json:"message"`
	DevOTP  string `json:"dev_otp,omitempty"`
}

type passwordSvc struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	txManager        repository.TransactionManager
	emailSvc         EmailService
	authCfg          *config.AuthConfig
}

func NewPasswordService(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	txManager repository.TransactionManager,
	emailSvc EmailService,
	authCfg *config.AuthConfig,
) PasswordSvc {
	return &passwordSvc{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		txManager:        txManager,
		emailSvc:         emailSvc,
		authCfg:          authCfg,
	}
}

func (s *passwordSvc) ForgotPassword(ctx context.Context, req *models.ForgotPasswordRequest) (*ForgotPasswordResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	email := strings.TrimSpace(strings.ToLower(req.Email))
	result := &ForgotPasswordResult{Message: ForgotPasswordSuccessMessage}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Warn("failed to lookup forgot-password user", zap.Error(err))
		return result, nil
	}
	if user == nil {
		return result, nil
	}

	otpSecret, err := s.effectiveOTPSecret()
	if err != nil {
		return nil, ErrGenerateOTP.WithCause(err)
	}

	otp, err := GenerateOTP()
	if err != nil {
		return nil, ErrGenerateOTP.WithCause(err)
	}

	verification := &models.UserVerification{
		UserID:      user.ID,
		Purpose:     models.VerificationPurposePasswordReset,
		Target:      email,
		OTPHash:     HashOTP(otp, otpSecret),
		Attempts:    0,
		MaxAttempts: 3,
		IsUsed:      false,
		ExpiresAt:   time.Now().UTC().Add(s.authCfg.OTPExpiry),
	}

	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.verificationRepo.DeleteActiveByPurpose(txCtx, user.ID, verification.Purpose, email); err != nil {
			return ErrUpdateVerification.WithCause(err)
		}
		if err := s.verificationRepo.CreateVerification(txCtx, verification); err != nil {
			return ErrCreateVerification.WithCause(err)
		}
		return nil
	}); err != nil {
		var appErr *models.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, ErrVerificationTx.WithCause(err)
	}

	if s.emailSvc != nil {
		emailCtx, emailCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer emailCancel()

		if err := s.emailSvc.SendForgetPasswordEmail(emailCtx, email, otp, user.FullName); err != nil {
			logger.Warn("failed to send forgot-password email", zap.Error(err))
		}
	}

	if !s.authCfg.IsProduction {
		result.DevOTP = otp
	}

	return result, nil
}

func (s *passwordSvc) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := validateNewPassword(req.NewPassword); err != nil {
		return err
	}
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := s.rejectPasswordReuse(user, req.NewPassword); err != nil {
		return err
	}

	verification, err := s.verificationRepo.GetActiveVerification(ctx, user.ID, models.VerificationPurposePasswordReset, email)
	if err != nil {
		return ErrGetVerification.WithCause(err)
	}
	if verification == nil {
		return ErrNoActiveOTP
	}
	if verification.IsUsed {
		return ErrInvalidOTP
	}
	if time.Now().UTC().After(verification.ExpiresAt) {
		return ErrOTPExpired
	}
	if verification.Attempts >= verification.MaxAttempts {
		return ErrToManyAttempts
	}

	otpSecret, err := s.effectiveOTPSecret()
	if err != nil {
		return ErrGenerateOTP.WithCause(err)
	}

	if !CompareOTP(req.OTP, verification.OTPHash, otpSecret) {
		if err := s.verificationRepo.IncrementAttempts(ctx, verification.ID); err != nil {
			return ErrUpdateVerification.WithCause(err)
		}
		return ErrInvalidOTP
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return ErrHashPassword.WithCause(err)
	}

	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.verificationRepo.MarkAsUsed(txCtx, verification.ID); err != nil {
			return ErrMarkOTPUsed.WithCause(err)
		}
		if err := s.userRepo.UpdatePasswordHash(txCtx, user.ID, passwordHash, user.SessionVersion); err != nil {
			var appErr *models.AppError
			if errors.As(err, &appErr) {
				return appErr
			}
			return ErrUpdatePassword.WithCause(err)
		}
		return nil
	}); err != nil {
		var appErr *models.AppError
		if errors.As(err, &appErr) {
			return appErr
		}
		return ErrVerificationTx.WithCause(err)
	}

	return nil
}

func (s *passwordSvc) ChangePassword(ctx context.Context, userID int64, req *models.ChangePasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := validateNewPassword(req.NewPassword); err != nil {
		return err
	}
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.PasswordHash == nil {
		return ErrPasswordNotSet
	}

	if err := utils.CheckPassword(req.CurrentPassword, *user.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrCurrentPasswordIncorrect
		}
		return ErrStoredPassword.WithCause(errors.New("stored password hash is invalid"))
	}
	if err := s.rejectPasswordReuse(user, req.NewPassword); err != nil {
		return err
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return ErrHashPassword.WithCause(err)
	}

	if err := s.userRepo.UpdatePasswordHash(ctx, user.ID, passwordHash, user.SessionVersion); err != nil {
		var appErr *models.AppError
		if errors.As(err, &appErr) {
			return appErr
		}
		return ErrUpdatePassword.WithCause(err)
	}

	return nil
}

func (s *passwordSvc) rejectPasswordReuse(user *models.User, newPassword string) error {
	if user == nil {
		return ErrUserNotFound
	}
	if user.PasswordHash == nil {
		return nil
	}
	if err := utils.CheckPassword(newPassword, *user.PasswordHash); err == nil {
		return ErrPasswordReuse
	}
	return nil
}

func (s *passwordSvc) effectiveOTPSecret() (string, error) {
	secret := strings.TrimSpace(s.authCfg.OTPSecret)
	if secret != "" {
		return secret, nil
	}
	if s.authCfg.IsProduction {
		return "", errors.New("OTP_SECRET is required in production")
	}
	return "development-otp-secret-change-me", nil
}

func validateNewPassword(password string) error {
	if len(password) > 72 {
		return ErrPasswordTooLong
	}
	if utf8.RuneCountInString(password) < 8 {
		return models.NewAppError("VALIDATION_ERROR", "New password must contain at least 8 characters.", 400, nil)
	}
	return nil
}
