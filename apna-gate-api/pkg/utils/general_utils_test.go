package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		value   string
		want    int64
		wantErr bool
	}{
		{name: "valid", value: "42", want: 42},
		{name: "zero", value: "0", wantErr: true},
		{name: "negative", value: "-1", wantErr: true},
		{name: "not number", value: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/items/:id", func(c *gin.Context) {
				got, err := GetIDParam(c, "id")
				if tt.wantErr {
					if err == nil && got != 0 {
						t.Fatalf("expected invalid id to return zero/error, got id=%d err=%v", got, err)
					}
					c.Status(http.StatusBadRequest)
					return
				}
				if err != nil || got != tt.want {
					t.Fatalf("GetIDParam = %d, %v; want %d, nil", got, err, tt.want)
				}
				c.Status(http.StatusNoContent)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/items/"+tt.value, nil)
			router.ServeHTTP(w, req)
		})
	}
}
