package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-server/internal/app"
	handlers "go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"
	"go-server/internal/models"
	authsvc "go-server/internal/services/authSvc"
	societysvc "go-server/internal/services/societySvc"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"

	"github.com/gin-gonic/gin"
)

const (
	routeTestSigningMaterial = "route-test-secret"
	routeTestJWTIssuer       = "route-test-issuer"
)

type routeSessionUsers struct{}

func (routeSessionUsers) GetByID(_ context.Context, id int64) (*models.User, error) {
	return &models.User{ID: id}, nil
}

type adminRouteSocietyService struct {
	societysvc.SocietyService
	ensureRoleErr   error
	onboardingCalls int
}

func (s *adminRouteSocietyService) EnsureMembershipRole(context.Context, int64, int64, ...string) error {
	return s.ensureRoleErr
}

func (s *adminRouteSocietyService) EnsureRole(context.Context, int64, int64, ...string) error {
	return s.ensureRoleErr
}

func (s *adminRouteSocietyService) GetOnboardingBootstrap(context.Context, int64) (*models.SocietyOnboardingBootstrapResponse, error) {
	s.onboardingCalls++
	return &models.SocietyOnboardingBootstrapResponse{
		Society:  &models.SocietyResponse{ID: 42, Status: models.SocietyStatusPending},
		NextPath: "/onboarding/soc16",
	}, nil
}

type adminRouteSubscriptionGuard struct {
	subscriptionsvc.SubscriptionGuardService
	operationalCalls int
	err              error
}

func (s *adminRouteSubscriptionGuard) EnsureSocietyOperational(context.Context, int64) error {
	s.operationalCalls++
	return s.err
}

func setupAdminRouteTest(t *testing.T, societySvc *adminRouteSocietyService, subscriptionGuard *adminRouteSubscriptionGuard) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handlerSet := &app.V1Handlers{Society: handlers.NewSocietyHandler(societySvc)}
	guardSet := guards.New(
		routeTestSigningMaterial,
		routeTestJWTIssuer,
		societySvc,
		guards.WithSubscriptionService(subscriptionGuard),
		guards.WithSessionUsers(routeSessionUsers{}),
	)
	SetupAdminRoutesV1(router.Group("/v1"), handlerSet, guardSet)
	return router
}

func adminRouteToken(t *testing.T) string {
	t.Helper()
	token, err := authsvc.GenerateToken(
		100,
		"owner@example.com",
		"+919999999999",
		string(models.GlobalRoleUser),
		true,
		false,
		authsvc.TokenTypeAccess,
		routeTestSigningMaterial,
		routeTestJWTIssuer,
		time.Hour,
	)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	return token
}

func performAdminRouteRequest(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+adminRouteToken(t))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestOnboardingBootstrapRequiresMembershipButNotOperationalSociety(t *testing.T) {
	societySvc := &adminRouteSocietyService{}
	subscriptionGuard := &adminRouteSubscriptionGuard{err: subscriptionsvc.ErrSubscriptionRequired}
	router := setupAdminRouteTest(t, societySvc, subscriptionGuard)

	rec := performAdminRouteRequest(t, router, http.MethodGet, "/v1/societies/42/onboarding/bootstrap", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if societySvc.onboardingCalls != 1 {
		t.Fatalf("onboarding calls = %d, want 1", societySvc.onboardingCalls)
	}
	if subscriptionGuard.operationalCalls != 0 {
		t.Fatalf("operational guard calls = %d, want 0", subscriptionGuard.operationalCalls)
	}
}

func TestOnboardingBootstrapRejectsNonAdminMember(t *testing.T) {
	societySvc := &adminRouteSocietyService{ensureRoleErr: societysvc.ErrForbiddenSociety}
	subscriptionGuard := &adminRouteSubscriptionGuard{}
	router := setupAdminRouteTest(t, societySvc, subscriptionGuard)

	rec := performAdminRouteRequest(t, router, http.MethodGet, "/v1/societies/42/onboarding/bootstrap", "")

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusForbidden, rec.Body.String())
	}
	if societySvc.onboardingCalls != 0 {
		t.Fatalf("onboarding calls = %d, want 0", societySvc.onboardingCalls)
	}
}

func TestDashboardAndOnboardingMutationsRemainOperationallyProtected(t *testing.T) {
	tests := []struct {
		method string
		path   string
		body   string
	}{
		{method: http.MethodGet, path: "/v1/societies/42/dashboard/bootstrap"},
		{method: http.MethodPost, path: "/v1/societies/42/flats", body: `{}`},
		{method: http.MethodPost, path: "/v1/societies/42/flats/generate", body: `{}`},
		{method: http.MethodPost, path: "/v1/societies/42/guards", body: `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			societySvc := &adminRouteSocietyService{}
			subscriptionGuard := &adminRouteSubscriptionGuard{err: subscriptionsvc.ErrSubscriptionRequired}
			router := setupAdminRouteTest(t, societySvc, subscriptionGuard)

			rec := performAdminRouteRequest(t, router, tt.method, tt.path, tt.body)

			if rec.Code != http.StatusPaymentRequired {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusPaymentRequired, rec.Body.String())
			}
			if subscriptionGuard.operationalCalls != 1 {
				t.Fatalf("operational guard calls = %d, want 1", subscriptionGuard.operationalCalls)
			}
		})
	}
}
