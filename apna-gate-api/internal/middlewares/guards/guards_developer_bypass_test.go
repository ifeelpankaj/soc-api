package guards_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-server/internal/middlewares/guards"
	"go-server/internal/models"
	"go-server/internal/requestctx"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type trackingSubscriptionGuard struct {
	called bool
	err    error
}

func (s *trackingSubscriptionGuard) EnsureSocietyOperational(_ context.Context, _ int64) error {
	s.called = true
	return s.err
}

func (s *trackingSubscriptionGuard) EnsureActiveSubscription(_ context.Context, _ int64) error {
	return s.err
}

func (s *trackingSubscriptionGuard) EnsureFeatureEnabled(_ context.Context, _ int64, _ string) error {
	return s.err
}

func (s *trackingSubscriptionGuard) CanAddFlat(_ context.Context, _ int64, _ int64) error {
	return s.err
}

func (s *trackingSubscriptionGuard) CanAddAdmin(_ context.Context, _ int64, _ int64) error {
	return s.err
}

func (s *trackingSubscriptionGuard) CanAddStaff(_ context.Context, _ int64, _ int64) error {
	return s.err
}

func (s *trackingSubscriptionGuard) CanAddResident(_ context.Context, _ int64, _ int64) error {
	return s.err
}

func (s *trackingSubscriptionGuard) CanAddResidentWithLock(_ context.Context, _ int64, _ int64) error {
	return s.err
}

func runOperationalGuard(t *testing.T, role string) (*httptest.ResponseRecorder, *gin.Context, *trackingSubscriptionGuard) {
	t.Helper()

	subGuard := &trackingSubscriptionGuard{err: subscriptionsvc.ErrSubscriptionExpired}
	g := guards.New("secret", "issuer", guards.WithSubscriptionService(subGuard))

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/societies/42/dashboard/bootstrap", nil)
	c.Params = gin.Params{{Key: "societyId", Value: "42"}}
	if role != "" {
		c.Set("role", role)
	}

	handler := g.OperationalSocietyForParam("societyId")[0]
	handler(c)
	if !c.IsAborted() {
		c.Next()
	}

	return rec, c, subGuard
}

func TestRequireOperationalSocietyDeveloperBypassesSubscription(t *testing.T) {
	rec, c, subGuard := runOperationalGuard(t, string(models.GlobalRoleDeveloper))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if subGuard.called {
		t.Fatal("expected subscription guard to be skipped for developer")
	}
	if !requestctx.HasDeveloperGuardBypass(c.Request.Context()) {
		t.Fatal("expected developer bypass flag on request context")
	}
}

func TestRequireOperationalSocietySuperAdminBypassesSubscription(t *testing.T) {
	rec, _, subGuard := runOperationalGuard(t, string(models.GlobalRoleSuperAdmin))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if subGuard.called {
		t.Fatal("expected subscription guard to be skipped for super_admin")
	}
}

func TestRequireOperationalSocietyRegularUserChecksSubscription(t *testing.T) {
	rec, _, subGuard := runOperationalGuard(t, string(models.GlobalRoleUser))

	if rec.Code != http.StatusPaymentRequired {
		t.Fatalf("expected status 402, got %d", rec.Code)
	}
	if !subGuard.called {
		t.Fatal("expected subscription guard to run for regular users")
	}
}

func TestRequireOperationalSocietyMissingRoleChecksSubscription(t *testing.T) {
	rec, _, subGuard := runOperationalGuard(t, "")

	if rec.Code != http.StatusPaymentRequired {
		t.Fatalf("expected status 402, got %d", rec.Code)
	}
	if !subGuard.called {
		t.Fatal("expected subscription guard to run when role is missing")
	}
}
