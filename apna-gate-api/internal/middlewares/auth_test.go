package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authsvc "go-server/internal/services/authSvc"

	"github.com/gin-gonic/gin"
)

const (
	testJWTSecret = "test-secret"
	testJWTIssuer = "test-issuer"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func issueToken(t *testing.T, tokenType string) string {
	t.Helper()

	token, err := authsvc.GenerateToken(
		42,
		"user@example.com",
		"+911234567890",
		"user",
		true,
		false,
		tokenType,
		testJWTSecret,
		testJWTIssuer,
		time.Hour,
	)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	return token
}

func TestAccessAuthMiddlewareCookie(t *testing.T) {
	token := issueToken(t, authsvc.TokenTypeAccess)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/auth/profile", nil)
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: token, Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})

	mw := AccessAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected success, got status %d", rec.Code)
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok || userID != 42 {
		t.Fatalf("expected user_id 42, got %d ok=%v", userID, ok)
	}
}

func TestAccessAuthMiddlewareBearer(t *testing.T) {
	token := issueToken(t, authsvc.TokenTypeAccess)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/auth/profile", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	mw := AccessAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected success, got status %d", rec.Code)
	}
}

func TestAccessAuthMiddlewarePrefersCookieOverBearer(t *testing.T) {
	cookieToken := issueToken(t, authsvc.TokenTypeAccess)
	bearerTokenValue := issueToken(t, authsvc.TokenTypeAccess)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/auth/profile", nil)
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: cookieToken, Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	c.Request.Header.Set("Authorization", "Bearer "+bearerTokenValue)

	mw := AccessAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected success, got status %d", rec.Code)
	}
}

func TestAccessAuthMiddlewareMissingToken(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/auth/profile", nil)

	mw := AccessAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got status %d", rec.Code)
	}
}

func TestRefreshAuthMiddlewareCookie(t *testing.T) {
	token := issueToken(t, authsvc.TokenTypeRefresh)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/auth/refresh", nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: token, Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})

	mw := RefreshAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected success, got status %d", rec.Code)
	}
}

func TestRefreshAuthMiddlewareBody(t *testing.T) {
	token := issueToken(t, authsvc.TokenTypeRefresh)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/v1/auth/refresh",
		strings.NewReader(`{"refresh_token":"`+token+`"}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")

	mw := RefreshAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected success, got status %d", rec.Code)
	}
}

func TestRefreshAuthMiddlewareBearer(t *testing.T) {
	token := issueToken(t, authsvc.TokenTypeRefresh)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/auth/refresh", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	mw := RefreshAuthMiddleware(testJWTSecret, testJWTIssuer, &versionUsers{})
	mw(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected success, got status %d", rec.Code)
	}
}
