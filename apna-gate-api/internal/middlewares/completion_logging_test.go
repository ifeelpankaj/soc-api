package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	"go-server/pkg/logger"
	"go-server/pkg/utils"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCompletionContract(t *testing.T) {
	dir := os.Getenv("COMPLETION_TEST_DIR")
	if dir == "" {
		dir = t.TempDir()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestCompletionContract$") //nolint:gosec // Relaunch this test binary with a fixed test selector; no application input.
		cmd.Env = append(os.Environ(), "COMPLETION_TEST_DIR="+dir)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
		return
	}
	cfg := logger.DefaultConfig()
	cfg.LogDir = dir
	if err := logger.InitLogger(cfg); err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	r.Use(RequestID(), Metrics(), Logger(), Recovery())
	r.GET("/success/:id", func(c *gin.Context) { utils.SuccessResponse(c, 201, "Created", nil) })
	r.GET("/page", func(c *gin.Context) { utils.PaginatedResponse(c, 200, "Listed", &models.PaginationResponse{}) })
	r.GET("/client", func(c *gin.Context) { utils.ErrorResponse(c, 400, "INVALID", "Bad input", nil) })
	r.GET("/server", func(c *gin.Context) { utils.ErrorResponse(c, 500, "INTERNAL", "Failed", nil) })
	r.GET("/panic", func(c *gin.Context) { panic("password=secret") })
	r.GET("/committed", func(c *gin.Context) { c.String(202, "accepted"); panic("late failure") })
	r.GET("/fallback", func(c *gin.Context) { c.Status(204) })
	r.GET("/diagnostic", func(c *gin.Context) {
		_ = c.Error(errors.New("diagnostic"))
		utils.ErrorResponse(c, 500, "INTERNAL", "Failed", nil)
	})
	paths := []string{"/success/42", "/page", "/client", "/server", "/panic", "/committed", "/fallback", "/missing", "/diagnostic"}
	for _, path := range paths {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), "GET", path+"?token=secret-query", nil))
	}
	_ = logger.Sync()
	data, err := os.ReadFile(filepath.Join(dir, "info.production.log")) //nolint:gosec // Parent test supplies its private temporary directory to the child.
	if err != nil {
		t.Fatal(err)
	}
	completions := map[string]map[string]interface{}{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatal(err)
		}
		if entry["event"] != "request_completed" {
			continue
		}
		path := entry["path"].(string)
		if completions[path] != nil {
			t.Fatalf("duplicate completion %s", path)
		}
		completions[path] = entry
		for _, field := range []string{"@timestamp", "status", "method", "route", "response_message", "error_code", "internal_error", "duration_ms", "request_id"} {
			if _, ok := entry[field]; !ok {
				t.Errorf("missing %s", field)
			}
		}
	}
	if len(completions) != len(paths) {
		t.Fatalf("got %d completions", len(completions))
	}
	for path, expected := range map[string]string{"/success/42": "Created", "/page": "Listed", "/client": "Bad input", "/server": "Failed", "/fallback": "No Content"} {
		if completions[path]["response_message"] != expected {
			t.Errorf("wrong message for %s", path)
		}
	}
	if completions["/success/42"]["route"] != "/success/:id" || completions["/missing"]["route"] != "unmatched" {
		t.Fatal("route templates missing")
	}
	if completions["/panic"]["status"] != float64(500) || completions["/committed"]["status"] != float64(202) {
		t.Fatal("panic status incorrect")
	}
	if strings.Contains(string(data), "secret-query") || strings.Contains(string(data), "password=secret") {
		t.Fatal("secret leaked")
	}
}
