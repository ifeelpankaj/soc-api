package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-server/internal/models"
	authsvc "go-server/internal/services/authSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SessionUserReader interface {
	GetByID(context.Context, int64) (*models.User, error)
}

func AccessAuthMiddleware(jwtSecret, jwtIssuer string, reader SessionUserReader) gin.HandlerFunc {
	return authMiddleware(jwtSecret, jwtIssuer, "access_token", authsvc.TokenTypeAccess, "Authentication required. Please login", nil, reader)
}

func RefreshAuthMiddleware(jwtSecret, jwtIssuer string, reader SessionUserReader) gin.HandlerFunc {
	return authMiddleware(jwtSecret, jwtIssuer, "refresh_token", authsvc.TokenTypeRefresh, "Refresh token missing. Please login again", resolveRefreshToken, reader)
}

type tokenResolver func(*gin.Context) string

func authMiddleware(jwtSecret, jwtIssuer, cookieName, tokenType, missingMessage string, resolver tokenResolver, reader SessionUserReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := resolveAuthToken(c, cookieName, resolver)
		if token == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, models.ErrCodeUnauthorized, missingMessage, nil)
			c.Abort()
			return
		}

		claims, err := authsvc.ValidateToken(token, jwtSecret, jwtIssuer, tokenType)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, models.ErrCodeInvalidToken, "Invalid authentication token", err)
			c.Abort()
			return
		}

		if reader == nil {
			utils.InternalServerErrorResponse(c, errors.New("session verifier unavailable"))
			c.Abort()
			return
		}
		user, lookupErr := reader.GetByID(c.Request.Context(), claims.UserID)
		if lookupErr != nil {
			utils.InternalServerErrorResponse(c, lookupErr)
			c.Abort()
			return
		}
		if user == nil || user.SessionVersion != claims.SessionVersion {
			utils.AppErrorResponse(c, authsvc.ErrSessionRevoked)
			c.Abort()
			return
		}
		setAuthContext(c, claims)
		c.Next()
	}
}

func resolveAuthToken(c *gin.Context, cookieName string, resolver tokenResolver) string {
	if token, err := c.Cookie(cookieName); err == nil && strings.TrimSpace(token) != "" {
		return strings.TrimSpace(token)
	}

	if resolver != nil {
		if token := resolver(c); token != "" {
			return token
		}
	}

	if token := bearerToken(c.GetHeader("Authorization")); token != "" {
		return token
	}

	return ""
}

func resolveRefreshToken(c *gin.Context) string {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ""
	}

	return strings.TrimSpace(req.RefreshToken)
}

func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func setAuthContext(c *gin.Context, claims *authsvc.Claims) {
	c.Set("user_id", claims.UserID)
	c.Set("session_version", claims.SessionVersion)
	c.Set("email", claims.Email)
	c.Set("role", claims.Role)
	c.Set("token_type", claims.TokenType)
	c.Set("phone_number", claims.PhoneNumber)
	c.Set("email_verified", claims.EmailVerified)
	c.Set("phone_verified", claims.PhoneVerified)
}

func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	value, ok := email.(string)
	return value, ok
}

func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	value, ok := role.(string)
	return value, ok
}

func GetUserPhoneNumberFromContext(c *gin.Context) (string, bool) {
	phoneNumber, exists := c.Get("phone_number")
	if !exists {
		return "", false
	}
	value, ok := phoneNumber.(string)
	return value, ok
}

func GetUserEmailVerifiedFromContext(c *gin.Context) (bool, bool) {
	emailVerified, exists := c.Get("email_verified")
	if !exists {
		return false, false
	}
	value, ok := emailVerified.(bool)
	return value, ok
}

func GetUserPhoneVerifiedFromContext(c *gin.Context) (bool, bool) {
	phoneVerified, exists := c.Get("phone_verified")
	if !exists {
		return false, false
	}
	value, ok := phoneVerified.(bool)
	return value, ok
}
