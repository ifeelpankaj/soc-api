package authsvc

import (
	"context"
	"errors"
	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/pkg/utils"
	"strings"
	"sync"
	"testing"
	"time"
)

type passwordUsers struct {
	repository.UserRepository
	mu       sync.Mutex
	user     models.User
	writeErr error
}

func (r *passwordUsers) GetByID(context.Context, int64) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u := r.user
	return &u, nil
}
func (r *passwordUsers) GetByEmail(ctx context.Context, _ string) (*models.User, error) {
	return r.GetByID(ctx, 1)
}
func (r *passwordUsers) UpdateLastLogin(context.Context, int64) error { return nil }
func (r *passwordUsers) UpdatePasswordHash(_ context.Context, _ int64, hash string, version int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.writeErr != nil {
		return r.writeErr
	}
	if r.user.SessionVersion != version {
		return models.NewAppError("PASSWORD_CHANGE_CONFLICT", "Changed concurrently", 409, nil)
	}
	r.user.PasswordHash = &hash
	r.user.SessionVersion++
	return nil
}
func userFixture(t *testing.T) *passwordUsers {
	t.Helper()
	hash, err := utils.HashPassword("old-password")
	if err != nil {
		t.Fatal(err)
	}
	email := "user@example.com"
	return &passwordUsers{user: models.User{ID: 1, Email: &email, PasswordHash: &hash, IsActive: true, EmailVerified: true}}
}
func codeOf(err error) string {
	var app *models.AppError
	if errors.As(err, &app) {
		return app.Code
	}
	return ""
}

func TestChangePasswordFailuresDoNotModifySession(t *testing.T) {
	for _, tc := range []struct{ name, current, next, confirm, hash, code string }{
		{name: "incorrect", current: "wrong-password", next: "new-password", confirm: "new-password", code: "CURRENT_PASSWORD_INCORRECT"},
		{name: "missing", current: "old-password", next: "new-password", confirm: "new-password", hash: "nil", code: "PASSWORD_NOT_SET"},
		{name: "malformed", current: "old-password", next: "new-password", confirm: "new-password", hash: "broken", code: "PASSWORD_VERIFICATION_FAILED"},
		{name: "mismatch", current: "old-password", next: "new-password", confirm: "different-password", code: "PASSWORD_MISMATCH"},
		{name: "reuse", current: "old-password", next: "old-password", confirm: "old-password", code: "PASSWORD_REUSE"},
		{name: "unicode", current: "old-password", next: strings.Repeat("😀", 19), confirm: strings.Repeat("😀", 19), code: "PASSWORD_TOO_LONG"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			users := userFixture(t)
			if tc.hash == "nil" {
				users.user.PasswordHash = nil
			} else if tc.hash != "" {
				users.user.PasswordHash = &tc.hash
			}
			svc := NewPasswordService(users, nil, nil, nil, &config.AuthConfig{})
			err := svc.ChangePassword(context.Background(), 1, &models.ChangePasswordRequest{CurrentPassword: tc.current, NewPassword: tc.next, ConfirmPassword: tc.confirm})
			if codeOf(err) != tc.code {
				t.Fatalf("error = %v, want %s", err, tc.code)
			}
			if users.user.SessionVersion != 0 {
				t.Fatal("failed request revoked session")
			}
		})
	}
}

func TestPasswordChangeRequiresNewLoginAndRejectsOldRefresh(t *testing.T) {
	users := userFixture(t)
	cfg := &config.AuthConfig{JWTSecret: "test-password-flow", JWTIssuer: "test", AccessExpiry: time.Hour, RefreshExpiry: time.Hour}
	sessions := NewSessionService(users, cfg)
	oldLogin, err := sessions.Login(context.Background(), &models.LoginRequest{Email: "user@example.com", Password: "old-password"})
	if err != nil {
		t.Fatal(err)
	}
	oldClaims, err := ValidateToken(oldLogin.RefreshToken, cfg.JWTSecret, cfg.JWTIssuer, TokenTypeRefresh)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewPasswordService(users, nil, nil, nil, cfg)
	err = svc.ChangePassword(context.Background(), 1, &models.ChangePasswordRequest{CurrentPassword: "old-password", NewPassword: "new-password", ConfirmPassword: "new-password"})
	if err != nil {
		t.Fatal(err)
	}
	if users.user.SessionVersion != 1 {
		t.Fatal("session not invalidated")
	}
	if _, err = sessions.Refresh(context.Background(), 1, oldClaims.SessionVersion); codeOf(err) != "SESSION_REVOKED" {
		t.Fatalf("old refresh accepted: %v", err)
	}
	if _, err = sessions.Login(context.Background(), &models.LoginRequest{Email: "user@example.com", Password: "old-password"}); err == nil {
		t.Fatal("old password accepted")
	}
	login, err := sessions.Login(context.Background(), &models.LoginRequest{Email: "user@example.com", Password: "new-password"})
	if err != nil {
		t.Fatal(err)
	}
	for kind, token := range map[string]string{TokenTypeAccess: login.AccessToken, TokenTypeRefresh: login.RefreshToken} {
		claims, err := ValidateToken(token, cfg.JWTSecret, cfg.JWTIssuer, kind)
		if err != nil || claims.SessionVersion != 1 {
			t.Fatalf("new token: %v", err)
		}
	}
}

func TestNewPasswordByteBoundaries(t *testing.T) {
	for _, value := range []string{strings.Repeat("a", 72), strings.Repeat("😀", 18), "12345678"} {
		if err := validateNewPassword(value); err != nil {
			t.Fatal(err)
		}
	}
	for _, value := range []string{strings.Repeat("a", 73), strings.Repeat("😀", 19), "1234567"} {
		if err := validateNewPassword(value); err == nil {
			t.Fatal("invalid password accepted")
		}
	}
}

func TestConcurrentPasswordUpdatesOnlyOneWins(t *testing.T) {
	users := userFixture(t)
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); results <- users.UpdatePasswordHash(context.Background(), 1, "new hash", 0) }()
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		} else if codeOf(err) != "PASSWORD_CHANGE_CONFLICT" {
			t.Fatal(err)
		}
	}
	if success != 1 || users.user.SessionVersion != 1 {
		t.Fatal("concurrent overwrite")
	}
}

type resetVerifications struct {
	repository.VerificationRepository
	verification models.UserVerification
	used         bool
}

func (r *resetVerifications) GetActiveVerification(context.Context, int64, models.VerificationPurpose, string) (*models.UserVerification, error) {
	return &r.verification, nil
}
func (r *resetVerifications) MarkAsUsed(context.Context, int64) error { r.used = true; return nil }

type resetTransaction struct {
	users        *passwordUsers
	verification *resetVerifications
}

func (tx resetTransaction) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	before := tx.users.user
	used := tx.verification.used
	if err := fn(ctx); err != nil {
		tx.users.user = before
		tx.verification.used = used
		return err
	}
	return nil
}
func TestResetPasswordRevokesSessionsAndRollsBackOTPOnFailure(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "rollback"}[fail], func(t *testing.T) {
			users := userFixture(t)
			verification := &resetVerifications{verification: models.UserVerification{ID: 1, OTPHash: HashOTP("123456", "test-otp"), MaxAttempts: 3, ExpiresAt: time.Now().Add(time.Hour)}}
			if fail {
				users.writeErr = errors.New("database write failed")
			}
			svc := NewPasswordService(users, verification, resetTransaction{users, verification}, nil, &config.AuthConfig{OTPSecret: "test-otp"})
			err := svc.ResetPassword(context.Background(), &models.ResetPasswordRequest{Email: "user@example.com", OTP: "123456", NewPassword: "reset-password", ConfirmPassword: "reset-password"})
			if fail {
				if err == nil || users.user.SessionVersion != 0 || verification.used {
					t.Fatal("reset failure was not rolled back")
				}
			} else {
				if err != nil || users.user.SessionVersion != 1 || !verification.used {
					t.Fatalf("reset failed: %v", err)
				}
				if err := utils.CheckPassword("reset-password", *users.user.PasswordHash); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
