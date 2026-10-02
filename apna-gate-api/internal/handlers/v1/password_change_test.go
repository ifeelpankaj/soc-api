package handlers

import (
	"context"
	"go-server/internal/config"
	"go-server/internal/models"
	authsvc "go-server/internal/services/authSvc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type changePasswordService struct {
	authsvc.PasswordSvc
	err error
}

func (s changePasswordService) ChangePassword(context.Context, int64, *models.ChangePasswordRequest) error {
	return s.err
}

func TestChangePasswordHTTPContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"success", nil, 200, ""},
		{"incorrect", authsvc.ErrCurrentPasswordIncorrect, 400, "CURRENT_PASSWORD_INCORRECT"},
		{"missing", authsvc.ErrPasswordNotSet, 400, "PASSWORD_NOT_SET"},
		{"stored-hash", authsvc.ErrStoredPassword, 500, "PASSWORD_VERIFICATION_FAILED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewAuthHandler(nil, nil, nil, changePasswordService{err: tc.err}, &config.AuthConfig{})
			router := gin.New()
			router.POST("/change-password", func(c *gin.Context) { c.Set("user_id", int64(1)); handler.ChangePassword(c) })
			req := httptest.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				"/change-password",
				strings.NewReader(`{"current_password":"old-password","new_password":"new-password","confirm_password":"new-password"}`),
			)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.status || (tc.code != "" && !strings.Contains(w.Body.String(), tc.code)) {
				t.Fatalf("response: %d %s", w.Code, w.Body.String())
			}
			cookies := w.Result().Cookies()
			if tc.err != nil && len(cookies) != 0 {
				t.Fatal("failed change cleared authentication")
			}
			if tc.err == nil {
				if len(cookies) < 2 {
					t.Fatal("success must clear access and refresh cookies")
				}
				for _, cookie := range cookies {
					if cookie.Value != "" || cookie.MaxAge >= 0 {
						t.Fatal("cookie was not cleared")
					}
				}
			}
		})
	}
}
