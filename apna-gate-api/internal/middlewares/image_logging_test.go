package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go-server/pkg/logger"
)

func TestImageLogsExcludeQueriesBodiesAndRequestDumps(t *testing.T) {
	dir := os.Getenv("APNA_GATE_IMAGE_LOG_TEST_DIR")
	if dir == "" {
		// The production rotating logger owns open handles for the process lifetime.
		// Isolate it so Windows can remove the temporary logs after the child exits.
		dir = t.TempDir()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestImageLogsExcludeQueriesBodiesAndRequestDumps$") //nolint:gosec // Relaunch this test binary with a fixed test selector; no application input.
		command.Env = append(os.Environ(), "APNA_GATE_IMAGE_LOG_TEST_DIR="+dir)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("log verification: %v\n%s", err, output)
		}
		return
	}
	cfg := logger.DefaultConfig()
	cfg.LogDir = dir
	if err := logger.InitLogger(cfg); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(LoggerWithConfig(LoggerConfig{EnableBody: true}))
	router.Use(RecoveryWithConfig(RecoveryConfig{EnableRequestDump: true}))
	router.PUT("/api/v1/auth/profile/avatar", func(c *gin.Context) { c.JSON(200, gin.H{"url": "https://signed.example/image?ik-s=response-secret"}) })
	router.GET("/api/v1/societies/:societyId/visitor-entries/:entryId/photo", func(c *gin.Context) { panic("test image panic") })
	for _, tc := range []struct{ method, path string }{
		{http.MethodPut, "/api/v1/auth/profile/avatar"},
		{http.MethodGet, "/api/v1/societies/2/visitor-entries/3/photo"},
	} {
		req := httptest.NewRequestWithContext(context.Background(), tc.method, tc.path+"?ik-s=query-secret", strings.NewReader("body-secret"))
		req.Header.Set("Authorization", "Bearer header-secret")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}
	_ = logger.Sync()
	data, err := os.ReadFile(filepath.Join(dir, "info.production.log")) //nolint:gosec // Parent test supplies its private temporary directory to the child.
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "Request completed") || !strings.Contains(text, "Panic recovered") {
		t.Fatal("expected actual request/recovery logs")
	}
	for _, secret := range []string{"query-secret", "body-secret", "header-secret", "response-secret", "request_dump"} {
		if strings.Contains(text, secret) {
			t.Fatalf("logs exposed %s", secret)
		}
	}
}
