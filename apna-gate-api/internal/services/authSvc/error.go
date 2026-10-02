package authsvc

import (
	"net/http"

	"go-server/internal/constants"
	"go-server/internal/models"
)

var (
	ErrCurrentPasswordIncorrect = models.NewAppError("CURRENT_PASSWORD_INCORRECT", "Current password is incorrect", http.StatusBadRequest, nil)
	ErrPasswordNotSet           = models.NewAppError("PASSWORD_NOT_SET", "No password is set. Please use password reset to create one.", http.StatusBadRequest, nil)
	ErrPasswordTooLong          = models.NewAppError("PASSWORD_TOO_LONG", "New password must use at most 72 bytes. Some characters use more than one byte.", http.StatusBadRequest, nil)
	ErrSessionRevoked           = models.NewAppError("SESSION_REVOKED", "Your session has ended. Please sign in again.", http.StatusUnauthorized, nil)
	ErrStoredPassword           = models.NewAppError("PASSWORD_VERIFICATION_FAILED", "Could not verify your password. Please try password reset or contact support.", http.StatusInternalServerError, nil)
	ErrInvalidCredentials       = models.NewAppError(
		"INVALID_CREDENTIALS",
		"invalid email or password",
		http.StatusUnauthorized,
		nil,
	)

	ErrEmailNotVerified = models.NewAppError(
		"EMAIL_NOT_VERIFIED",
		"please verify your email before login",
		http.StatusForbidden,
		nil,
	)

	ErrUserInactive = models.NewAppError(
		"USER_INACTIVE",
		"user account is inactive",
		http.StatusForbidden,
		nil,
	)

	ErrUserBlocked = models.NewAppError(
		"USER_BLOCKED",
		"user account is blocked",
		http.StatusForbidden,
		nil,
	)

	ErrGenerateToken = models.NewAppError(
		"GENERATE_TOKEN_FAILED",
		"failed to generate authentication token",
		http.StatusInternalServerError,
		nil,
	)

	ErrInvalidToken = models.NewAppError(
		"INVALID_TOKEN",
		"invalid or expired token",
		http.StatusUnauthorized,
		nil,
	)

	ErrWrongTokenType = models.NewAppError(
		"WRONG_TOKEN_TYPE",
		"invalid token type",
		http.StatusUnauthorized,
		nil,
	)
	ErrGetVerification = models.NewAppError(
		models.ErrCodeDatabaseError,
		"failed to get verification record",
		500,
		nil,
	)
	ErrToManyAttempts = models.NewAppError(
		models.ErrCodeTooManyRequests,
		"too many incorrect OTP attempts. Please request a new OTP.",
		http.StatusTooManyRequests,
		nil,
	)

	ErrInvalidName = models.NewAppError(
		models.ErrCodeValidation,
		"First name is required.",
		http.StatusBadRequest,
		nil,
	)
	ErrCheckEmail = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrCheckEmail,
		http.StatusInternalServerError,
		nil,
	)
	ErrEmailExists = models.NewAppError(
		models.ErrCodeEmailExists,
		constants.ErrEmailExists,
		http.StatusConflict,
		nil,
	)
	ErrCheckPhone = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrCheckPhone,
		http.StatusInternalServerError,
		nil,
	)
	ErrPhoneExists = models.NewAppError(
		models.ErrCodePhoneExists,
		constants.ErrPhoneExists,
		http.StatusConflict,
		nil,
	)
	ErrHashPassword = models.NewAppError(
		models.ErrCodeInternalServer,
		constants.ErrHashPassword,
		http.StatusInternalServerError,
		nil,
	)
	ErrGenerateOTP = models.NewAppError(
		models.ErrCodeInternalServer,
		constants.ErrGenerateOTP,
		http.StatusInternalServerError,
		nil,
	)
	ErrCreateUser = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrCreateUser,
		http.StatusInternalServerError,
		nil,
	)
	ErrCreateVerification = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrCreateVerification,
		http.StatusInternalServerError,
		nil,
	)
	ErrRegistrationTx = models.NewAppError(
		models.ErrCodeTransactionFailed,
		constants.ErrRegistrationTx,
		http.StatusInternalServerError,
		nil,
	)
	ErrSendEmail = models.NewAppError(
		models.ErrCodeEmailSendFailed,
		constants.ErrSendEmail,
		http.StatusServiceUnavailable,
		nil,
	)
	ErrUnknownID = models.NewAppError(
		models.ErrCodeBadRequest,
		constants.ErrUnknownID,
		http.StatusBadRequest,
		nil,
	)

	ErrAuthNotFound = models.NewAppError(
		models.ErrCodeNotFound,
		"Record not found. Please check and try again.",
		http.StatusNotFound,
		nil,
	)

	// =========================
	// User
	// =========================

	ErrUserNotFound = models.NewAppError(
		models.ErrCodeNotFound,
		constants.ErrUserNotFound,
		http.StatusNotFound,
		nil,
	)

	ErrFetchUser = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrFetchUser,
		http.StatusInternalServerError,
		nil,
	)

	ErrUpdateUser = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrUpdateUser,
		http.StatusInternalServerError,
		nil,
	)

	ErrUserUnverified = models.NewAppError(
		models.ErrCodeUserNotVerified,
		constants.ErrUserUnverified,
		http.StatusUnauthorized,
		nil,
	)

	ErrAlreadyVerified = models.NewAppError(
		models.ErrCodeEmailNotVerified,
		constants.ErrAlreadyVerified,
		http.StatusConflict,
		nil,
	)

	// =========================
	// Login / Password
	// =========================

	ErrInvalidLogin = models.NewAppError(
		models.ErrCodeInvalidCredentials,
		constants.ErrInvalidLogin,
		http.StatusUnauthorized,
		nil,
	)

	ErrPasswordMismatch = models.NewAppError(
		"PASSWORD_MISMATCH",
		constants.ErrPasswordMismatch,
		http.StatusBadRequest,
		nil,
	)

	ErrPasswordReuse = models.NewAppError(
		"PASSWORD_REUSE",
		"new password must be different from the current password",
		http.StatusBadRequest,
		nil,
	)

	ErrUpdatePassword = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrUpdateUser,
		http.StatusInternalServerError,
		nil,
	)

	ErrLastLoginUpdate = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.MsgLastLoginSkipped,
		http.StatusInternalServerError,
		nil,
	)

	// =========================
	// OTP / Verification
	// =========================

	ErrInvalidOTP = models.NewAppError(
		models.ErrCodeInvalidOTP,
		constants.ErrInvalidOTP,
		http.StatusBadRequest,
		nil,
	)

	ErrNoActiveOTP = models.NewAppError(
		models.ErrCodeNotFound,
		constants.ErrNoActiveOTP,
		http.StatusNotFound,
		nil,
	)

	ErrOTPExpired = models.NewAppError(
		models.ErrCodeOTPExpired,
		"The OTP has expired. Please request a new one.",
		http.StatusBadRequest,
		nil,
	)
	ErrFindVerificationUser = models.NewAppError(
		models.ErrCodeNotFound,
		constants.ErrFindVerificationUser,
		http.StatusNotFound,
		nil,
	)

	ErrUpdateVerification = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrUpdateVerification,
		http.StatusInternalServerError,
		nil,
	)

	ErrMarkOTPUsed = models.NewAppError(
		models.ErrCodeDatabaseError,
		constants.ErrMarkOTPUsed,
		http.StatusInternalServerError,
		nil,
	)

	ErrVerificationTx = models.NewAppError(
		models.ErrCodeTransactionFailed,
		constants.ErrVerificationTx,
		http.StatusInternalServerError,
		nil,
	)

	// =========================
	// Email
	// =========================

	ErrSendResetEmail = models.NewAppError(
		models.ErrCodeEmailSendFailed,
		constants.ErrSendEmail,
		http.StatusInternalServerError,
		nil,
	)
)
