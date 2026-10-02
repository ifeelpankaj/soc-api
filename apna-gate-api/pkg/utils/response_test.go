package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-server/internal/models"

	"github.com/gin-gonic/gin"
)

func performResponse(handler gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	SetResponseConfig(ResponseConfig{EnableDetailedErrors: true, IncludeRequestID: true})
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.Set("request_id", "req-1")
		c.Set("user_id", "7")
		handler(c)
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))
	return w
}

func TestSuccessResponseHelpers(t *testing.T) {
	w := performResponse(func(c *gin.Context) {
		CreatedResponse(c, "created", gin.H{"id": 1}, "/items/1")
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
	if got := w.Header().Get("Location"); got != "/items/1" {
		t.Fatalf("Location = %q", got)
	}
	if got := w.Header().Get("X-Request-ID"); got != "req-1" {
		t.Fatalf("X-Request-ID = %q", got)
	}

	var body models.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if !body.Success || body.Message != "created" {
		t.Fatalf("unexpected response body: %#v", body)
	}

	w = performResponse(func(c *gin.Context) { AcceptedResponse(c, "accepted", nil) })
	if w.Code != http.StatusAccepted {
		t.Fatalf("accepted status = %d", w.Code)
	}

	w = performResponse(NoContentResponse)
	if w.Code != http.StatusNoContent {
		t.Fatalf("no content status = %d", w.Code)
	}
}

func TestErrorResponseHelpers(t *testing.T) {
	tests := []struct {
		name       string
		handler    gin.HandlerFunc
		wantStatus int
		wantCode   string
		wantMsg    string
	}{
		{"unauthorized default", func(c *gin.Context) { UnauthorizedResponse(c, "") }, http.StatusUnauthorized, models.ErrCodeUnauthorized, "Authentication required"},
		{"forbidden default", func(c *gin.Context) { ForbiddenResponse(c, "") }, http.StatusForbidden, models.ErrCodeForbidden, "Access forbidden"},
		{"not found resource", func(c *gin.Context) { NotFoundResponse(c, "Flat") }, http.StatusNotFound, models.ErrCodeNotFound, "Flat not found"},
		{"conflict custom", func(c *gin.Context) { ConflictResponse(c, "duplicate visitor") }, http.StatusConflict, models.ErrCodeConflict, "duplicate visitor"},
		{"bad request default", func(c *gin.Context) { BadRequestResponse(c, "") }, http.StatusBadRequest, models.ErrCodeBadRequest, "Invalid request"},
		{"service unavailable default", func(c *gin.Context) { ServiceUnavailableResponse(c, "") }, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service temporarily unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performResponse(tt.handler)
			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			var body models.APIResponse
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("invalid JSON response: %v", err)
			}
			if body.Success || body.Error == nil || body.Error.Code != tt.wantCode || body.Error.Message != tt.wantMsg {
				t.Fatalf("unexpected body: %#v", body)
			}
		})
	}
}

func TestDetailedErrorAndValidationResponses(t *testing.T) {
	w := performResponse(func(c *gin.Context) {
		InternalServerErrorResponse(c, errors.New("database down"))
	})
	var body models.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if w.Code != http.StatusInternalServerError || body.Error == nil || body.Error.Details["error"] != "database down" {
		t.Fatalf("unexpected internal error body: status=%d body=%#v", w.Code, body)
	}

	w = performResponse(func(c *gin.Context) {
		ValidationErrorResponse(c, map[string]interface{}{"flat_id": "must be positive"})
	})
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body.Error == nil || body.Error.Details["flat_id"] != "must be positive" {
		t.Fatalf("validation details missing: %#v", body)
	}

	w = performResponse(func(c *gin.Context) {
		AppErrorResponse(c, models.NewAppError(models.ErrCodeConflict, "duplicate", http.StatusConflict, nil))
	})
	if w.Code != http.StatusConflict {
		t.Fatalf("AppErrorResponse status = %d", w.Code)
	}
}

func TestPaginatedAndRateLimitResponses(t *testing.T) {
	w := performResponse(func(c *gin.Context) {
		PaginatedResponse(c, http.StatusOK, "ok", &models.PaginationResponse{
			Data: []string{"a"}, Page: 1, PageSize: 10, TotalPages: 1, TotalItems: 1,
		})
	})
	if w.Code != http.StatusOK {
		t.Fatalf("paginated status = %d", w.Code)
	}

	w = performResponse(func(c *gin.Context) { PaginatedResponse(c, http.StatusOK, "ok", nil) })
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("nil pagination status = %d", w.Code)
	}

	w = performResponse(func(c *gin.Context) { TooManyRequestsResponse(c, 30) })
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("rate limit status = %d", w.Code)
	}
	if got := w.Header().Get("Retry-After"); got == "" {
		t.Fatal("expected Retry-After header")
	}
}
