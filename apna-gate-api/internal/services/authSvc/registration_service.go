package authsvc

import (
	"context"
	"strings"
	"time"

	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	service "go-server/internal/services"
	"go-server/pkg/utils"
)

type RegistrationSvc interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*RegistrationResult, error)
	RegisterResident(ctx context.Context, req *models.ResidentRegisterRequest) (*LoginResult, error)
	CreateResidentUser(ctx context.Context, req *models.ResidentRegisterRequest) (*models.UserResponse, error)
}
type RegistrationResult struct {
	User   *models.UserResponse `json:"user"`
	DevOTP string               `json:"dev_otp,omitempty"`
}

type registrationSvc struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	txManager        repository.TransactionManager
	emailSvc         EmailService
	authCfg          *config.AuthConfig
}

func NewRegistrationService(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	txManager repository.TransactionManager,
	emailSvc EmailService,
	authCfg *config.AuthConfig,
) RegistrationSvc {
	return &registrationSvc{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		txManager:        txManager,
		emailSvc:         emailSvc,
		authCfg:          authCfg,
	}
}

func (s *registrationSvc) Register(ctx context.Context, req *models.RegisterRequest) (*RegistrationResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	fullName := buildFullName(req.FirstName, req.LastName)

	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, ErrCheckEmail.WithCause(err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	exists, err = s.userRepo.PhoneExists(ctx, req.PhoneNumber)
	if err != nil {
		return nil, ErrCheckPhone.WithCause(err)
	}
	if exists {
		return nil, ErrPhoneExists
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, ErrHashPassword.WithCause(err)
	}

	otpSecret := s.authCfg.OTPSecret

	otp, err := GenerateOTP()
	if err != nil {
		return nil, ErrGenerateOTP.WithCause(err)
	}
	otpHash := HashOTP(otp, otpSecret)

	firstName := req.FirstName
	lastName := nullableTrimmed(req.LastName)
	userEmail := req.Email
	userPhone := req.PhoneNumber
	password := passwordHash

	user := &models.User{
		FirstName:     &firstName,
		LastName:      lastName,
		FullName:      fullName,
		Email:         &userEmail,
		PhoneNumber:   &userPhone,
		PasswordHash:  &password,
		AuthProvider:  models.AuthProviderEmail,
		GlobalRole:    models.GlobalRoleUser,
		EmailVerified: false,
		PhoneVerified: false,
		IsActive:      true,
		IsBlocked:     false,
		Timezone:      "Asia/Kolkata",
		Language:      "en",
		Metadata:      map[string]any{},
	}

	verification := &models.UserVerification{
		Purpose:     models.VerificationPurposeEmailVerification,
		Target:      req.Email,
		OTPHash:     otpHash,
		Attempts:    0,
		MaxAttempts: 3,
		IsUsed:      false,
		ExpiresAt:   time.Now().UTC().Add(s.authCfg.OTPExpiry),
	}

	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.Create(txCtx, user); err != nil {
			return ErrCreateUser.WithCause(err)
		}

		verification.UserID = user.ID
		if err := s.verificationRepo.DeleteActiveByPurpose(txCtx, user.ID, verification.Purpose, verification.Target); err != nil {
			return ErrCreateVerification.WithCause(err)
		}
		if err := s.verificationRepo.CreateVerification(txCtx, verification); err != nil {
			return ErrCreateVerification.WithCause(err)
		}

		return nil
	}); err != nil {
		return nil, ErrRegistrationTx.WithCause(err)
	}

	emailCtx, emailCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer emailCancel()

	if err := s.emailSvc.SendOTP(emailCtx, req.Email, otp, fullName); err != nil {
		return nil, ErrSendEmail.WithCause(err)
	}

	result := &RegistrationResult{
		User: user.ToResponse(),
	}

	if !s.authCfg.IsProduction {
		result.DevOTP = otp
	}

	return result, nil
}

func (s *registrationSvc) RegisterResident(ctx context.Context, req *models.ResidentRegisterRequest) (*LoginResult, error) {
	user, err := s.createResidentUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, ErrLastLoginUpdate.WithCause(err)
	}

	accessToken, err := GenerateToken(
		user.ID,
		user.EmailValue(),
		user.PhoneNumberValue(),
		string(user.GlobalRole),
		user.EmailVerified,
		user.PhoneVerified,
		TokenTypeAccess,
		s.authCfg.JWTSecret,
		s.authCfg.JWTIssuer,
		s.authCfg.AccessExpiry,
		user.SessionVersion,
	)
	if err != nil {
		return nil, ErrGenerateToken.WithCause(err)
	}

	refreshToken, err := GenerateToken(
		user.ID,
		user.EmailValue(),
		user.PhoneNumberValue(),
		string(user.GlobalRole),
		user.EmailVerified,
		user.PhoneVerified,
		TokenTypeRefresh,
		s.authCfg.JWTSecret,
		s.authCfg.JWTIssuer,
		s.authCfg.RefreshExpiry,
		user.SessionVersion,
	)
	if err != nil {
		return nil, ErrGenerateToken.WithCause(err)
	}

	return &LoginResult{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *registrationSvc) CreateResidentUser(ctx context.Context, req *models.ResidentRegisterRequest) (*models.UserResponse, error) {
	user, err := s.createResidentUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *registrationSvc) createResidentUser(ctx context.Context, req *models.ResidentRegisterRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	fullName := buildFullName(req.FirstName, req.LastName)

	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, ErrCheckEmail.WithCause(err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	exists, err = s.userRepo.PhoneExists(ctx, req.PhoneNumber)
	if err != nil {
		return nil, ErrCheckPhone.WithCause(err)
	}
	if exists {
		return nil, ErrPhoneExists
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, ErrHashPassword.WithCause(err)
	}

	firstName := req.FirstName
	lastName := nullableTrimmed(req.LastName)
	userEmail := req.Email
	userPhone := req.PhoneNumber
	password := passwordHash

	user := &models.User{
		FirstName:     &firstName,
		LastName:      lastName,
		FullName:      fullName,
		Email:         &userEmail,
		PhoneNumber:   &userPhone,
		PasswordHash:  &password,
		AuthProvider:  models.AuthProviderEmail,
		GlobalRole:    models.GlobalRoleUser,
		EmailVerified: true,
		PhoneVerified: false,
		IsActive:      true,
		IsBlocked:     false,
		Timezone:      "Asia/Kolkata",
		Language:      "en",
		Metadata: map[string]any{
			"registered_from": "member_invite_flow",
		},
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, ErrCreateUser.WithCause(err)
	}

	return user, nil
}

func buildFullName(firstName, lastName string) string {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if lastName == "" {
		return firstName
	}
	return firstName + " " + lastName
}

func nullableTrimmed(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
