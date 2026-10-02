package handlers

import (
	"net/http"

	"go-server/internal/config"
	middleware "go-server/internal/middlewares"
	"go-server/internal/models"
	authsvc "go-server/internal/services/authSvc"
	"go-server/pkg/utils"
	"go-server/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	registrationSvc authsvc.RegistrationSvc
	verificationSvc authsvc.VerificationSvc
	sessionSvc      authsvc.SessionSvc
	passwordSvc     authsvc.PasswordSvc
	authCfg         *config.AuthConfig
}

func NewAuthHandler(
	registrationSvc authsvc.RegistrationSvc,
	verificationSvc authsvc.VerificationSvc,
	sessionSvc authsvc.SessionSvc,
	passwordSvc authsvc.PasswordSvc,
	authCfg *config.AuthConfig,
) *AuthHandler {
	return &AuthHandler{
		registrationSvc: registrationSvc,
		verificationSvc: verificationSvc,
		sessionSvc:      sessionSvc,
		passwordSvc:     passwordSvc,
		authCfg:         authCfg,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a user account and sends an email verification OTP.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration payload"
// @Success 201 {object} models.RegisterAPIResponse "Account created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 409 {object} models.ErrorResponseDoc "Email or phone already exists"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	result, err := h.registrationSvc.Register(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}

	data := gin.H{
		"user":    result.User,
		"message": "Please verify your email using the OTP sent to your email address",
	}

	if result.DevOTP != "" {
		data["dev_otp"] = result.DevOTP
	}

	utils.SuccessResponse(c, http.StatusCreated, "Account created successfully", data)
}

// RegisterResident godoc
// @Summary Register resident without OTP
// @Description Creates a resident user for the public flat claim flow, marks the email verified, and sets auth cookies immediately. This does not create a society membership or flat resident.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ResidentRegisterRequest true "Resident registration payload"
// @Success 201 {object} models.LoginAPIResponse "Resident account created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 409 {object} models.ErrorResponseDoc "Email or phone already exists"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/resident/register [post]
func (h *AuthHandler) RegisterResident(c *gin.Context) {
	var req models.ResidentRegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	result, err := h.registrationSvc.RegisterResident(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}

	setAccessTokenCookie(c, h.authCfg, result.AccessToken)
	setRefreshTokenCookie(c, h.authCfg, result.RefreshToken)

	utils.SuccessResponse(c, http.StatusCreated, "Resident account created successfully",
		buildAuthSessionData(h.authCfg, result.User, result.AccessToken, result.RefreshToken, true),
	)
}

// VerifyOTP godoc
// @Summary Verify email OTP
// @Description Verifies the registration email OTP and marks the user's email as verified.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Email verification payload"
// @Success 200 {object} models.VerifyOTPAPIResponse "Email verified successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid OTP or validation error"
// @Failure 404 {object} models.ErrorResponseDoc "User or OTP not found"
// @Failure 429 {object} models.ErrorResponseDoc "Too many OTP attempts"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	result, err := h.verificationSvc.VerifyEmail(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Email verified successfully", gin.H{
		"user": result.User,
	})
}

// ResendOTP godoc
// @Summary Resend email verification OTP
// @Description Sends a fresh email verification OTP for an unverified user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ResendOTPRequest true "Resend OTP payload"
// @Success 200 {object} models.ResendOTPAPIResponse "OTP sent successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 404 {object} models.ErrorResponseDoc "User not found"
// @Failure 409 {object} models.ErrorResponseDoc "Email already verified"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req models.ResendOTPRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	result, err := h.verificationSvc.ResendEmailOTP(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}

	data := gin.H{
		"message": result.Message,
	}

	if result.DevOTP != "" {
		data["dev_otp"] = result.DevOTP
	}

	utils.SuccessResponse(c, http.StatusOK, result.Message, data)
}

// Login godoc
// @Summary Login
// @Description Authenticates a verified active user and sets access_token and refresh_token HTTP-only cookies.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login payload"
// @Success 200 {object} models.LoginAPIResponse "Login successful"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Invalid credentials"
// @Failure 403 {object} models.ErrorResponseDoc "Email not verified or account disabled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	if req.Email == "" && req.PhoneNumber == "" {
		utils.ValidationErrorResponse(c, map[string]interface{}{
			"identifier": "email or phone_number is required",
		})
		return
	}
	if req.Email != "" && req.PhoneNumber != "" {
		utils.ValidationErrorResponse(c, map[string]interface{}{
			"identifier": "use either email or phone_number, not both",
		})
		return
	}

	result, err := h.sessionSvc.Login(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}

	setAccessTokenCookie(c, h.authCfg, result.AccessToken)
	setRefreshTokenCookie(c, h.authCfg, result.RefreshToken)

	utils.SuccessResponse(c, http.StatusOK, "Login successful",
		buildAuthSessionData(h.authCfg, result.User, result.AccessToken, result.RefreshToken, true),
	)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Validates the refresh token from either the refresh_token cookie or the request body.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} models.RefreshTokenAPIResponse
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 403 {object} models.ErrorResponseDoc
// @Failure 500 {object} models.ErrorResponseDoc
// @Security RefreshToken
// @Router /v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	result, err := h.sessionSvc.Refresh(c.Request.Context(), userID, c.GetInt64("session_version"))
	if handleServiceError(c, err) {
		return
	}

	setAccessTokenCookie(c, h.authCfg, result.AccessToken)

	utils.SuccessResponse(c, http.StatusOK, "Access token refreshed successfully",
		buildRefreshSessionData(h.authCfg, result.AccessToken),
	)
}

// Logout godoc
// @Summary Logout
// @Description Clears access_token and refresh_token cookies. This route is public and succeeds even if cookies are absent.
// @Tags Auth
// @Produce json
// @Success 200 {object} models.LogoutAPIResponse "Logout successful"
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	clearAuthCookies(c, h.authCfg)

	utils.SuccessResponse(c, http.StatusOK, "Logout successful", gin.H{
		"message": "Logged out successfully",
	})
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Returns the authenticated user's profile. Does not refresh or rotate tokens.
// @Tags Auth
// @Produce json
// @Success 200 {object} models.GetProfileAPIResponse "Profile fetched successfully"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Account disabled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	user, err := h.sessionSvc.GetProfile(c.Request.Context(), userID)
	if handleServiceError(c, err) {
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile fetched successfully", gin.H{
		"user": user,
	})
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Updates authenticated user profile fields such as name, phone, and date of birth.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.UpdateUserRequest true "Profile update payload"
// @Success 200 {object} models.GetProfileAPIResponse "Profile updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Account disabled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/auth/profile [patch]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	var req models.UpdateUserRequest
	if !bindJSON(c, &req) {
		return
	}
	req.ID = userID

	user, err := h.sessionSvc.UpdateProfile(c.Request.Context(), userID, &req)
	if handleServiceError(c, err) {
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", gin.H{
		"user": user,
	})
}

// ForgotPassword godoc
// @Summary Request password reset OTP
// @Description Always returns a generic success message. If the email exists, sends a password reset OTP.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Forgot password payload"
// @Success 200 {object} models.ForgotPasswordAPIResponse "Password reset instructions response"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	result, err := h.passwordSvc.ForgotPassword(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}

	data := gin.H{
		"message": result.Message,
	}
	if result.DevOTP != "" {
		data["dev_otp"] = result.DevOTP
	}

	utils.SuccessResponse(c, http.StatusOK, result.Message, data)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Resets a user's password using a valid password reset OTP and clears auth cookies on success.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} models.ResetPasswordAPIResponse "Password reset successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid OTP, password reuse, or validation error"
// @Failure 404 {object} models.ErrorResponseDoc "User or active OTP not found"
// @Failure 429 {object} models.ErrorResponseDoc "Too many OTP attempts"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	if err := h.passwordSvc.ResetPassword(c.Request.Context(), &req); err != nil {
		if handleServiceError(c, err) {
			return
		}
	}

	clearAuthCookies(c, h.authCfg)
	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", gin.H{
		"message": "Password reset successfully. Please login again.",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Changes the authenticated user's password after verifying the current password, invalidates all access and refresh sessions, then clears auth cookies.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Change password payload"
// @Success 200 {object} models.ChangePasswordAPIResponse "Password changed successfully"
// @Failure 400 {object} models.ErrorResponseDoc "CURRENT_PASSWORD_INCORRECT, PASSWORD_NOT_SET, PASSWORD_TOO_LONG, mismatch, reuse, or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	if err := h.passwordSvc.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		if handleServiceError(c, err) {
			return
		}
	}

	clearAuthCookies(c, h.authCfg)
	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", gin.H{
		"message": "Password changed successfully. Please login again.",
	})
}
