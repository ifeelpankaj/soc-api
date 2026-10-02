package routes

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	handlers "go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"
)

func TestHubRoutesRequireAuthenticationAndHaveSwagger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupHubRoutesV1(router.Group("/api/v1"), handlers.NewHubHandler(nil), guards.New(routeTestSigningMaterial, routeTestJWTIssuer, guards.WithSessionUsers(routeSessionUsers{})))
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	raw, e := os.ReadFile("../../docs/swagger.json")
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw, &spec); e != nil {
		t.Fatal(e)
	}
	params := regexp.MustCompile(`:([A-Za-z]+)`)
	routes := router.Routes()
	if len(routes) != 24 {
		t.Fatalf("routes %d", len(routes))
	}
	for _, route := range routes {
		path := params.ReplaceAllString(route.Path, "1")
		req := httptest.NewRequestWithContext(context.Background(), route.Method, path, nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != 401 {
			t.Fatalf("unauthenticated %s %s: %d", route.Method, path, res.Code)
		}
		swaggerPath := params.ReplaceAllString(strings.TrimPrefix(route.Path, "/api"), "{$1}")
		if len(spec.Paths[swaggerPath][strings.ToLower(route.Method)]) == 0 {
			t.Errorf("missing Swagger: %s %s", route.Method, swaggerPath)
		}
	}
}
