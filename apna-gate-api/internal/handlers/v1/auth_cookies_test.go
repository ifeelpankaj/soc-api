package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-server/internal/config"

	"github.com/gin-gonic/gin"
)

func TestSetAccessTokenCookieUsesConfiguredDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	cfg := &config.AuthConfig{
		AccessExpiry: time.Minute,
		IsProduction: true,
		CookieDomain: "apnagate.org",
	}

	setAccessTokenCookie(ctx, cfg, "access-token")

	header := rec.Header().Values("Set-Cookie")[0]
	if !strings.Contains(header, "Domain=apnagate.org") {
		t.Fatalf("Set-Cookie = %q, want Domain=apnagate.org", header)
	}
	if !strings.Contains(header, "Secure") {
		t.Fatalf("Set-Cookie = %q, want Secure", header)
	}
	if !strings.Contains(header, "HttpOnly") {
		t.Fatalf("Set-Cookie = %q, want HttpOnly", header)
	}
}

func TestSetAccessTokenCookieOmitsEmptyDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	cfg := &config.AuthConfig{
		AccessExpiry: time.Minute,
	}

	setAccessTokenCookie(ctx, cfg, "access-token")

	header := rec.Header().Values("Set-Cookie")[0]
	if strings.Contains(header, "Domain=") {
		t.Fatalf("Set-Cookie = %q, want no Domain attribute", header)
	}
}

func TestClearAuthCookiesUsesConfiguredDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	cfg := &config.AuthConfig{
		IsProduction: true,
		CookieDomain: "apnagate.org",
	}

	clearAuthCookies(ctx, cfg)

	headers := rec.Header().Values("Set-Cookie")
	if len(headers) != 4 {
		t.Fatalf("Set-Cookie header count = %d, want 4", len(headers))
	}

	for _, header := range headers[:2] {
		if !strings.Contains(header, "Domain=apnagate.org") {
			t.Fatalf("Set-Cookie = %q, want Domain=apnagate.org", header)
		}
	}

	for _, header := range headers[2:] {
		if strings.Contains(header, "Domain=") {
			t.Fatalf("Set-Cookie = %q, want no Domain attribute", header)
		}
	}
}

func TestClearAuthCookiesOmitsExtraHostOnlyCookiesWithoutConfiguredDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	cfg := &config.AuthConfig{}

	clearAuthCookies(ctx, cfg)

	headers := rec.Header().Values("Set-Cookie")
	if len(headers) != 2 {
		t.Fatalf("Set-Cookie header count = %d, want 2", len(headers))
	}

	for _, header := range headers {
		if strings.Contains(header, "Domain=") {
			t.Fatalf("Set-Cookie = %q, want no Domain attribute", header)
		}
	}
}

func TestSetRefreshTokenCookieUsesRefreshPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	cfg := &config.AuthConfig{
		RefreshExpiry: time.Hour,
		CookieDomain:  "apnagate.org",
	}

	setRefreshTokenCookie(ctx, cfg, "refresh-token")

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d, want 1", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "refresh_token" {
		t.Fatalf("cookie name = %q, want refresh_token", cookie.Name)
	}
	if cookie.Path != "/api/v1/auth/refresh" {
		t.Fatalf("cookie path = %q, want /api/v1/auth/refresh", cookie.Path)
	}
	if cookie.Domain != "apnagate.org" {
		t.Fatalf("cookie domain = %q, want apnagate.org", cookie.Domain)
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie SameSite = %v, want Lax", cookie.SameSite)
	}
}
