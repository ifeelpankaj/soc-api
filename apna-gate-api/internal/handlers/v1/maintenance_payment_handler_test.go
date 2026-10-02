package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-server/internal/models"

	"github.com/gin-gonic/gin"
)

type paymentHandlerFake struct {
	maintenancePayments
	society, user, bill int64
	key                 string
	filter              models.UPIListFilter
	calls               int
	err                 error
}

func (s *paymentHandlerFake) SubmitClaim(_ context.Context, society, user, bill int64, key string, _ models.UPISubmitClaim) (models.UPIClaim, error) {
	s.society, s.user, s.bill, s.key = society, user, bill, key
	s.calls++
	return models.UPIClaim{}, s.err
}
func (s *paymentHandlerFake) Payments(_ context.Context, f models.UPIListFilter) (models.UPIPaymentsPage, error) {
	s.filter = f
	s.calls++
	return models.UPIPaymentsPage{Items: []models.UPIPayment{}}, s.err
}
func (s *paymentHandlerFake) Claims(_ context.Context, f models.UPIListFilter) (models.UPIClaimsPage, error) {
	s.filter = f
	s.calls++
	return models.UPIClaimsPage{Items: []models.UPIClaim{}}, s.err
}
func (s *paymentHandlerFake) QR(_ context.Context, society, user int64, id string) ([]byte, error) {
	s.society, s.user, s.key = society, user, id
	s.calls++
	return []byte("PNG"), s.err
}
func paymentTestRouter(s *paymentHandlerFake, authenticated bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if authenticated {
		r.Use(func(c *gin.Context) { c.Set("user_id", int64(7)) })
	}
	h := NewMaintenancePaymentHandler(s)
	r.GET("/societies/:societyId/maintenance/payments", h.Payments)
	r.POST("/societies/:societyId/maintenance/my/bills/:id/payment-claims", h.SubmitPaymentClaim)
	r.GET("/societies/:societyId/maintenance/my/payments", h.MyPayments)
	r.GET("/societies/:societyId/maintenance/my/payment-claims", h.MyPaymentClaims)
	r.GET("/societies/:societyId/maintenance/my/payment-requests/:requestId/qr", h.PaymentQR)
	return r
}

func TestAdminPaymentHandlerFilters(t *testing.T) {
	s := &paymentHandlerFake{}
	r := paymentTestRouter(s, true)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/payments?bill_id=12&flat_id=4&billing_month=2026-10&block=left-wing&flat_number=G-02&cursor=99&limit=10", nil))
	if w.Code != 200 || s.filter.Resident || s.filter.BillID != 12 || s.filter.FlatID != 4 || s.filter.BillingMonth != "2026-10" || s.filter.Block != "left-wing" || s.filter.FlatNumber != "G-02" || s.filter.SocietyID != 42 || s.filter.UserID != 7 {
		t.Fatalf("filters: %d %+v", w.Code, s.filter)
	}
	for _, suffix := range []string{"payments?flat_id=0", "payments?bill_id=-1", "payments?bill_id=abc", "my/payments?flat_id=0"} {
		before := s.calls
		w = httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/"+suffix, nil))
		if w.Code != 400 || s.calls != before {
			t.Fatalf("invalid filter reached service: %s %d", suffix, w.Code)
		}
	}
}
func TestResidentPaymentFiltersRemainAuthenticated(t *testing.T) {
	s := &paymentHandlerFake{}
	r := paymentTestRouter(s, true)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/payments?bill_id=12&flat_id=4&billing_month=2026-10&user_id=999", nil))
	if w.Code != 200 || !s.filter.Resident || s.filter.UserID != 7 || s.filter.BillID != 12 || s.filter.FlatID != 4 || s.filter.BillingMonth != "2026-10" {
		t.Fatalf("resident ownership scope changed: %d %+v", w.Code, s.filter)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/payment-claims?bill_id=12&flat_id=4&billing_month=2026-10&user_id=999", nil))
	if w.Code != 200 || !s.filter.Resident || s.filter.UserID != 7 || s.filter.BillID != 12 || s.filter.FlatID != 4 || s.filter.BillingMonth != "2026-10" {
		t.Fatalf("resident claim ownership scope changed: %d %+v", w.Code, s.filter)
	}
}
func TestPaymentHandlerScopeIdempotencyAndConflicts(t *testing.T) {
	s := &paymentHandlerFake{err: models.NewAppError("BILL_ALREADY_PAID", "Report additional transfer", 409, nil)}
	r := paymentTestRouter(s, true)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/societies/42/maintenance/my/bills/123/payment-claims", strings.NewReader(`{"payment_request_id":"saved-request","reference":"0012345678","payment_date":"2026-09-10","amount_paise":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "claim-retry-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 409 || !strings.Contains(w.Body.String(), "BILL_ALREADY_PAID") || s.society != 42 || s.user != 7 || s.bill != 123 || s.key != "claim-retry-key" {
		t.Fatalf("%d %s %+v", w.Code, w.Body, s)
	}
	s.err = nil
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/payments?cursor=99&limit=10", nil))
	if w.Code != 200 || !s.filter.Resident || s.filter.SocietyID != 42 || s.filter.UserID != 7 || s.filter.BeforeID != 99 || s.filter.Limit != 10 {
		t.Fatalf("scope: %+v", s.filter)
	}
}
func TestPaymentQRAuthenticationAndNoCache(t *testing.T) {
	for _, authenticated := range []bool{false, true} {
		s := &paymentHandlerFake{}
		r := paymentTestRouter(s, authenticated)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/payment-requests/request-id/qr", nil))
		if !authenticated {
			if w.Code != 401 || s.calls != 0 {
				t.Fatalf("unauthenticated: %d %d", w.Code, s.calls)
			}
			continue
		}
		if w.Code != 200 || w.Header().Get("Content-Type") != "image/png" || w.Header().Get("Cache-Control") != "private, no-store" || s.user != 7 || s.society != 42 || s.key != "request-id" {
			t.Fatalf("QR: %d %+v %+v", w.Code, w.Header(), s)
		}
		s.err = models.NewAppError("PAYMENT_REQUEST_OBSOLETE", "Closed request", 409, nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/payment-requests/request-id/qr", nil))
		if w.Code != 409 || strings.Contains(w.Body.String(), "PNG") {
			t.Fatal("obsolete QR returned")
		}
	}
}
