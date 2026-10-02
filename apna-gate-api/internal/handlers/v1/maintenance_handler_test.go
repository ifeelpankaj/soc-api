package handlers

import (
	"context"
	"go-server/internal/models"
	maintenancesvc "go-server/internal/services/maintenanceSvc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type maintenanceHandlerFake struct {
	maintenanceService
	filter      models.MaintenanceBillFilter
	calls       int
	generateErr error
}

func (s *maintenanceHandlerFake) List(_ context.Context, f models.MaintenanceBillFilter) (models.MaintenanceBillList, error) {
	s.filter = f
	s.calls++
	return models.MaintenanceBillList{Items: []models.MaintenanceBill{}}, nil
}
func (s *maintenanceHandlerFake) Get(_ context.Context, f models.MaintenanceBillFilter) (models.MaintenanceBill, error) {
	s.filter = f
	s.calls++
	return models.MaintenanceBill{}, models.NewAppError("MAINTENANCE_BILL_NOT_FOUND", "Bill not found", 404, nil)
}
func (s *maintenanceHandlerFake) GenerateCommand(context.Context, int64, int64, string, models.MaintenanceMonthRequest) (models.MaintenanceRunResult, error) {
	s.calls++
	return models.MaintenanceRunResult{}, s.generateErr
}
func maintenanceTestRouter(s *maintenanceHandlerFake) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_id", int64(7)) })
	h := NewMaintenanceHandler(s)
	r.GET("/societies/:societyId/maintenance/bills", h.ListBills)
	r.GET("/societies/:societyId/maintenance/my/bills", h.MyBills)
	r.GET("/societies/:societyId/maintenance/my/bills/:id", h.MyBill)
	r.POST("/societies/:societyId/maintenance/bills/generate", h.GenerateBills)
	return r
}
func TestMaintenanceResidentHandlerScope(t *testing.T) {
	s := &maintenanceHandlerFake{}
	r := maintenanceTestRouter(s)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/bills?flat_id=9&cursor=100&limit=10", nil))
	if w.Code != 200 || !s.filter.Resident || s.filter.UserID != 7 || s.filter.SocietyID != 42 || s.filter.FlatID != 9 || s.filter.BeforeID != 100 || s.filter.Limit != 10 {
		t.Fatalf("%d %+v", w.Code, s.filter)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/bills/999", nil))
	if w.Code != 404 || s.filter.ID != 999 || !s.filter.Resident {
		t.Fatalf("%d %+v", w.Code, s.filter)
	}
}
func TestMaintenanceHandlerParsesFlatSearchFilters(t *testing.T) {
	s := &maintenanceHandlerFake{}
	r := maintenanceTestRouter(s)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/bills?block=left-wing&flat_number=G-02&search=g", nil))
	if w.Code != 200 || s.filter.Block != "left-wing" || s.filter.FlatNumber != "G-02" || s.filter.Search != "g" || s.filter.Resident {
		t.Fatalf("%d %+v", w.Code, s.filter)
	}
}

func TestMaintenanceHandlerRejectsMalformedFilters(t *testing.T) {
	for _, query := range []string{"limit=0", "limit=101", "limit=bad", "cursor=-1", "flat_id=bad"} {
		s := &maintenanceHandlerFake{}
		r := maintenanceTestRouter(s)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/societies/42/maintenance/my/bills?"+query, nil))
		if w.Code != 400 || s.calls != 0 {
			t.Fatalf("%s: %d", query, w.Code)
		}
	}
}
func TestMaintenanceHandlerReturnsFlatValidationDetails(t *testing.T) {
	s := &maintenanceHandlerFake{generateErr: &maintenancesvc.ValidationIssues{Issues: []models.MaintenanceFlatIssue{{FlatID: 10, Reason: "area_sqft is required"}}}}
	r := maintenanceTestRouter(s)
	w := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/societies/42/maintenance/bills/generate", strings.NewReader(`{"billing_month":"2026-10"}`))
	request.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, request)
	if w.Code != 422 || !strings.Contains(w.Body.String(), `"flat_id":10`) || !strings.Contains(w.Body.String(), `"success":false`) {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
}
