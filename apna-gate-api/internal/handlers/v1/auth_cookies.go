package handlers

import (
	"net/http"
	"time"

	"go-server/internal/config"
	"go-server/internal/models"

	"github.com/gin-gonic/gin"
)

func setAccessTokenCookie(c *gin.Context, cfg *config.AuthConfig, token string) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"access_token",
		token,
		int(cfg.AccessExpiry/time.Second),
		"/",
		cfg.CookieDomain,
		cfg.IsProduction,
		true,
	)
}

func setRefreshTokenCookie(c *gin.Context, cfg *config.AuthConfig, token string) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"refresh_token",
		token,
		int(cfg.RefreshExpiry/time.Second),
		"/api/v1/auth/refresh",
		cfg.CookieDomain,
		cfg.IsProduction,
		true,
	)
}

func clearAuthCookies(c *gin.Context, cfg *config.AuthConfig) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("access_token", "", -1, "/", cfg.CookieDomain, cfg.IsProduction, true)
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth/refresh", cfg.CookieDomain, cfg.IsProduction, true)
	if cfg.CookieDomain != "" {
		c.SetCookie("access_token", "", -1, "/", "", cfg.IsProduction, true)
		c.SetCookie("refresh_token", "", -1, "/api/v1/auth/refresh", "", cfg.IsProduction, true)
	}
}

func buildAuthSessionData(cfg *config.AuthConfig, user *models.UserResponse, accessToken, refreshToken string, includeRefresh bool) gin.H {
	now := time.Now().UTC()
	data := gin.H{
		"user":                    user,
		"access_token":            accessToken,
		"access_token_expires_at": now.Add(cfg.AccessExpiry).Format(time.RFC3339),
	}

	if includeRefresh {
		data["refresh_token"] = refreshToken
		data["refresh_token_expires_at"] = now.Add(cfg.RefreshExpiry).Format(time.RFC3339)
	}

	return data
}

func buildRefreshSessionData(cfg *config.AuthConfig, accessToken string) gin.H {
	now := time.Now().UTC()
	return gin.H{
		"message":                 "Access token refreshed successfully",
		"access_token":            accessToken,
		"access_token_expires_at": now.Add(cfg.AccessExpiry).Format(time.RFC3339),
	}
}
