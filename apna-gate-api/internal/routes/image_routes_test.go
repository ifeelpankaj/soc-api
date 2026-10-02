package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	handlers "go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"
	authsvc "go-server/internal/services/authSvc"
	imagesvc "go-server/internal/services/imageSvc"
)

func TestImageRoutesAuthenticateAndDisable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupImageRoutesV1(router.Group("/api/v1"), handlers.NewImageHandler(imagesvc.New(nil, nil, nil, "dev/apna-gate", nil)), guards.New(routeTestSigningMaterial, routeTestJWTIssuer, guards.WithSessionUsers(routeSessionUsers{})))
	token, err := authsvc.GenerateToken(42, "test@example.com", "+911234567890", "user", true, false, authsvc.TokenTypeAccess, routeTestSigningMaterial, routeTestJWTIssuer, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/v1/auth/profile/avatar", "/api/v1/societies/2/visitor-entries/3/photo", "/api/v1/societies/2/flats/4/visitor-entries/3/photo"} {
		for _, authenticated := range []bool{false, true} {
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
			if authenticated {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			want := 401
			if authenticated {
				want = 503
			}
			if res.Code != want || res.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("%s: %d %s", path, res.Code, res.Body.String())
			}
		}
	}
}
