package handlers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	hubsvc "go-server/internal/services/hubSvc"
)

func TestHubBindingRejectsIdentityAndTrailingJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, body := range []string{`{"body":"test","author_id":123}`, `{"body":"test","society_id":12}`, `{"body":"test"} {}`, `{`} {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequestWithContext(context.Background(), "POST", "/", strings.NewReader(body))
		var req models.HubPostRequest
		if hubBind(c, &req) || w.Code != 400 {
			t.Fatalf("accepted %s", body)
		}
	}
}
func TestHubRateResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequestWithContext(context.Background(), "POST", "/", nil)
	hubResult(c, nil, &hubsvc.RateError{RetryAfter: 7})
	if w.Code != 429 || w.Header().Get("Retry-After") != "7" {
		t.Fatal(w.Code, w.Header())
	}
}
