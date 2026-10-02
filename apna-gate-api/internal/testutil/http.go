package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func NewRecorderRequest(t *testing.T, method string, target string, body any) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}

	req := httptest.NewRequestWithContext(context.Background(), method, target, reader)
	req.Header.Set("Content-Type", "application/json")
	return httptest.NewRecorder(), req
}

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}
