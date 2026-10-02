package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-server/internal/models"
	flatsvc "go-server/internal/services/flatSvc"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"

	"github.com/gin-gonic/gin"
)

type generationHandlerFlatService struct {
	flatsvc.FlatService
	request *models.GenerateFlatsRequest
	err     error
}

func (s *generationHandlerFlatService) GenerateFlats(_ context.Context, _ int64, _ int64, req *models.GenerateFlatsRequest) (*models.BulkCreateFlatsResponse, error) {
	s.request = req
	if s.err != nil {
		return nil, s.err
	}
	return &models.BulkCreateFlatsResponse{Items: []*models.FlatResponse{}, Total: 3}, nil
}

func generationHandlerRouter(service *generationHandlerFlatService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/societies/:societyId/flats/generate", func(c *gin.Context) {
		c.Set("user_id", int64(7))
		NewFlatHandler(service).GenerateFlats(c)
	})
	return router
}

func TestGenerateFlatsHandlerSuccess(t *testing.T) {
	service := &generationHandlerFlatService{}
	router := generationHandlerRouter(service)
	body := `{"block":" A ","numbering_mode":"continuous","flats_per_floor":3,"sequence_start":"001","total_flats":3}`
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/v1/societies/42/flats/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"total":3`) {
		t.Fatalf("status/body = %d/%s", rec.Code, rec.Body.String())
	}
	if service.request == nil || service.request.Block != "A" || service.request.StartFloorValue() != 1 {
		t.Fatalf("sanitized/default request = %#v", service.request)
	}
}

func TestGenerateFlatsHandlerRejectsInvalidRequest(t *testing.T) {
	service := &generationHandlerFlatService{}
	router := generationHandlerRouter(service)
	body := `{"block":"A","numbering_mode":"continuous","flats_per_floor":3,"sequence_start":"001","total_flats":0}`
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/v1/societies/42/flats/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest || service.request != nil {
		t.Fatalf("status/service request = %d/%#v; body=%s", rec.Code, service.request, rec.Body.String())
	}
}

func TestGenerateFlatsHandlerMapsServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "quota", err: subscriptionsvc.ErrQuotaExceeded, want: http.StatusPaymentRequired},
		{name: "conflict", err: flatsvc.ErrFlatConflict, want: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &generationHandlerFlatService{err: tt.err}
			router := generationHandlerRouter(service)
			body := `{"block":"A","numbering_mode":"continuous","flats_per_floor":3,"sequence_start":"001","total_flats":3}`
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/v1/societies/42/flats/generate", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tt.want, rec.Body.String())
			}
		})
	}
}
