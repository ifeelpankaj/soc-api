package middleware

import (
	"context"
	"errors"
	"go-server/internal/models"
	authsvc "go-server/internal/services/authSvc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type versionUsers struct {
	version int64
	err     error
}

func (r *versionUsers) GetByID(context.Context, int64) (*models.User, error) {
	return &models.User{ID: 1, SessionVersion: r.version}, r.err
}

func TestRevocationRejectsEveryDeviceAndTokenType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := &versionUsers{}
	for _, kind := range []string{authsvc.TokenTypeAccess, authsvc.TokenTypeRefresh} {
		for _, device := range []string{"phone", "tablet"} {
			t.Run(kind+device, func(t *testing.T) {
				token, err := authsvc.GenerateToken(1, "user@example.com", "", "user", true, false, kind, "session-test", "session-test", time.Hour)
				if err != nil {
					t.Fatal(err)
				}
				router := gin.New()
				mw := AccessAuthMiddleware("session-test", "session-test", users)
				if kind == authsvc.TokenTypeRefresh {
					mw = RefreshAuthMiddleware("session-test", "session-test", users)
				}
				router.GET("/protected", mw, func(c *gin.Context) { c.Status(204) })
				call := func() *httptest.ResponseRecorder {
					w := httptest.NewRecorder()
					req := httptest.NewRequestWithContext(
						context.Background(),
						http.MethodGet,
						"/protected",
						nil,
					)
					req.Header.Set("Authorization", "Bearer "+token)
					router.ServeHTTP(w, req)
					return w
				}
				users.version = 0
				if w := call(); w.Code != 204 {
					t.Fatalf("legacy zero-version token rejected: %s", w.Body.String())
				}
				users.version = 1
				if w := call(); w.Code != 401 || !strings.Contains(w.Body.String(), "SESSION_REVOKED") {
					t.Fatalf("revoked token accepted: %d %s", w.Code, w.Body.String())
				}
				users.err = errors.New("database unavailable")
				if w := call(); w.Code != 500 {
					t.Fatal("lookup failure must fail closed")
				}
				users.err = nil
			})
		}
	}
}

func TestRefreshBodyTakesPrecedenceOverAccessBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	access, _ := authsvc.GenerateToken(1, "", "", "user", true, false, authsvc.TokenTypeAccess, "test", "test", time.Hour)
	refresh, _ := authsvc.GenerateToken(1, "", "", "user", true, false, authsvc.TokenTypeRefresh, "test", "test", time.Hour)
	router := gin.New()
	router.POST("/refresh", RefreshAuthMiddleware("test", "test", &versionUsers{}), func(c *gin.Context) { c.Status(204) })
	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/refresh",
		strings.NewReader(`{"refresh_token":"`+refresh+`"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("refresh failed: %d %s", w.Code, w.Body.String())
	}
}
