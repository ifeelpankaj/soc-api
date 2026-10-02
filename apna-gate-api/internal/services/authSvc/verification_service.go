package authsvc

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	service "go-server/internal/services"
)

type VerificationSvc interface {
	VerifyEmail(ctx context.Context, req *models.VerifyOTPRequest) (*VerificationResult, error)
	ResendEmailOTP(ctx context.Context, req *models.ResendOTPRequest) (*ResendOTPResult, error)
}

type VerificationResult struct {
	User *models.UserResponse `json:"user"`
}

type ResendOTPResult struct {
	Message string `json:"message"`
	DevOTP  string `json:"dev_otp,omitempty"`
}

type verificationSvc struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	txManager        repository.TransactionManager
	emailSvc         EmailService
	authCfg          *config.AuthConfig
}

func NewVerificationService(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	txManager repository.TransactionManager,
	emailSvc EmailService,
	authCfg *config.AuthConfig,
) VerificationSvc {
	return &verificationSvc{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		txManager:        txManager,
		emailSvc:         emailSvc,
		authCfg:          authCfg,
	}
}

func (s *verificationSvc) VerifyEmail(
	ctx context.Context,
	req *models.VerifyOTPRequest,
) (*VerificationResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.EmailVerified {
		return &VerificationResult{
			User: user.ToResponse(),
		}, nil
	}

	verification, err := s.verificationRepo.GetActiveVerification(
		ctx,
		user.ID,
		models.VerificationPurposeEmailVerification,
		req.Email,
	)
	if err != nil {
		return nil, ErrGetVerification.WithCause(err)
	}
	if verification == nil {
		return nil, ErrInvalidOTP
	}

	if verification.IsUsed {
		return nil, ErrInvalidOTP
	}

	if time.Now().UTC().After(verification.ExpiresAt) {
		return nil, ErrOTPExpired
	}

	if verification.Attempts >= verification.MaxAttempts {
		return nil, ErrToManyAttempts
	}

	otpSecret, err := s.effectiveOTPSecret()
	if err != nil {
		return nil, ErrGenerateOTP.WithCause(err)
	}

	otpHash := HashOTP(req.OTP, otpSecret)
	if otpHash != verification.OTPHash {
		if err := s.verificationRepo.IncrementAttempts(ctx, verification.ID); err != nil {
			return nil, ErrUpdateVerification.WithCause(err)
		}
		return nil, ErrInvalidOTP
	}

	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.verificationRepo.MarkAsUsed(txCtx, verification.ID); err != nil {
			return ErrUpdateVerification.WithCause(err)
		}

		if err := s.userRepo.MarkEmailVerified(txCtx, user.ID); err != nil {
			return ErrUpdateUser.WithCause(err)
		}

		return nil
	}); err != nil {
		var appErr *models.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, ErrVerificationTx.WithCause(err)
	}

	user.EmailVerified = true

	return &VerificationResult{
		User: user.ToResponse(),
	}, nil
}

func (s *verificationSvc) ResendEmailOTP(
	ctx context.Context,
	req *models.ResendOTPRequest,
) (*ResendOTPResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.EmailVerified {
		return &ResendOTPResult{
			Message: "email already verified",
		}, nil
	}

	otpSecret, err := s.effectiveOTPSecret()
	if err != nil {
		return nil, ErrGenerateOTP.WithCause(err)
	}

	otp, err := GenerateOTP()
	if err != nil {
		return nil, ErrGenerateOTP.WithCause(err)
	}

	otpHash := HashOTP(otp, otpSecret)

	verification := &models.UserVerification{
		UserID:      user.ID,
		Purpose:     models.VerificationPurposeEmailVerification,
		Target:      email,
		OTPHash:     otpHash,
		Attempts:    0,
		MaxAttempts: 3,
		IsUsed:      false,
		ExpiresAt:   time.Now().UTC().Add(s.authCfg.OTPExpiry),
	}

	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.verificationRepo.DeleteActiveByPurpose(
			txCtx,
			user.ID,
			verification.Purpose,
			email,
		); err != nil {
			return ErrCreateVerification.WithCause(err)
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

		if err := s.emailSvc.SendOTP(emailCtx, email, otp, user.FullName); err != nil {
			return nil, ErrSendEmail.WithCause(err)
		}
	}

	result := &ResendOTPResult{
		Message: "verification OTP sent successfully",
	}

	if !s.authCfg.IsProduction {
		result.DevOTP = otp
	}

	return result, nil
}

func (s *verificationSvc) effectiveOTPSecret() (string, error) {
	secret := strings.TrimSpace(s.authCfg.OTPSecret)
	if secret != "" {
		return secret, nil
	}

	if s.authCfg.IsProduction {
		return "", errors.New("OTP_SECRET is required in production")
	}

	return "development-otp-secret-change-me", nil
}
