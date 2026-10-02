package authsvc

import (
	"context"
	"net/http"
	"strings"

	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	service "go-server/internal/services"
	"go-server/pkg/utils"
	"go-server/pkg/validator"
)

type SessionSvc interface {
	Login(ctx context.Context, req *models.LoginRequest) (*LoginResult, error)
	AuthenticateCredentials(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, error)
	Refresh(ctx context.Context, userID int64, version int64) (*RefreshResult, error)
	GetProfile(ctx context.Context, userID int64) (*models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID int64, req *models.UpdateUserRequest) (*models.UserResponse, error)
}

type LoginResult struct {
	User         *models.UserResponse `json:"user"`
	AccessToken  string               `json:"-"`
	RefreshToken string               `json:"-"`
}

type RefreshResult struct {
	AccessToken string `json:"-"`
}

type sessionSvc struct {
	userRepo repository.UserRepository
	authCfg  *config.AuthConfig
}

func NewSessionService(
	userRepo repository.UserRepository,
	authCfg *config.AuthConfig,
) SessionSvc {
	return &sessionSvc{
		userRepo: userRepo,
		authCfg:  authCfg,
	}
}

func (s *sessionSvc) Login(ctx context.Context, req *models.LoginRequest) (*LoginResult, error) {
	user, err := s.authenticateUser(ctx, req)
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

func (s *sessionSvc) AuthenticateCredentials(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, error) {
	user, err := s.authenticateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *sessionSvc) authenticateUser(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	email := strings.TrimSpace(strings.ToLower(req.Email))
	phoneNumber := strings.TrimSpace(req.PhoneNumber)

	var (
		user *models.User
		err  error
	)
	if email != "" {
		user, err = s.userRepo.GetByEmail(ctx, email)
	} else {
		user, err = s.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	}
	if err != nil {
		return nil, ErrInvalidCredentials.WithCause(err)
	}
	if user == nil || user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.CheckPassword(req.Password, *user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if user.IsBlocked {
		return nil, ErrUserBlocked
	}

	return user, nil
}

func (s *sessionSvc) Refresh(ctx context.Context, userID int64, version int64) (*RefreshResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.SessionVersion != version {
		return nil, ErrSessionRevoked
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if user.IsBlocked {
		return nil, ErrUserBlocked
	}

	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
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

	return &RefreshResult{
		AccessToken: accessToken,
	}, nil
}

func (s *sessionSvc) GetProfile(ctx context.Context, userID int64) (*models.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	if user.IsBlocked {
		return nil, ErrUserBlocked
	}

	return user.ToResponse(), nil
}

func (s *sessionSvc) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	if req.AvatarURL != nil && strings.TrimSpace(*req.AvatarURL) != "" {
		return nil, models.ErrAvatarUploadRequired
	}
	req.AvatarURL = nil
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, models.NewAppError(
			models.ErrCodeValidation,
			validationErrors.Error(),
			http.StatusBadRequest,
			validationErrors,
		)
	}
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidName.WithCause(err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound.WithCause(err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	if user.IsBlocked {
		return nil, ErrUserBlocked
	}

	updated, err := s.userRepo.UpdateProfile(ctx, userID, (*repository.UpdateUserInput)(req))
	if err != nil {
		return nil, ErrUpdateUser.WithCause(err)
	}
	if updated == nil {
		return nil, ErrUserNotFound
	}
	return updated.ToResponse(), nil
}
