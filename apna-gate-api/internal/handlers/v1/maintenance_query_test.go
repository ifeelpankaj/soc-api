package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMaintenanceFlatQueryRejectsLongValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		if _, ok := maintenanceFlatQuery(c); !ok {
			return
		}
		c.Status(http.StatusOK)
	})
	long := strings.Repeat("a", 51)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test?block="+long, nil))
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
