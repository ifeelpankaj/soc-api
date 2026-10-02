package routes

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-server/internal/backup"
	handlers "go-server/internal/handlers/v1"
)

type backupFake struct {
	calls int
	err   error
	id    string
}

func (f *backupFake) Trigger(_ context.Context, kind string) (*backup.Run, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return &backup.Run{ID: f.id, Type: kind, Status: "queued"}, nil
}
func (f *backupFake) Get(context.Context, string) (*backup.Run, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &backup.Run{ID: f.id, Status: "completed"}, nil
}
func backupRouter(f *backupFake) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	SetupInternalJobRoutes(r.Group("/api/internal"), handlers.NewJobWebhookHandler(&jobTriggererFake{}, f), testJobWebhookSecret)
	return r
}
func TestBackupStrictBodyAndAuth(t *testing.T) {
	for _, body := range []string{`{}`, `null`, `[]`, `{"type":null}`, `{"type":123}`, `{"type":"monthly"}`, `{"Type":"daily"}`, `{"type":"daily","type":"weekly"}`, `{"type":"daily","database":"other"}`, `{"type":"daily"} {}`, `{"type":"daily"`, strings.Repeat(" ", 1100) + `{"type":"daily"}`} {
		t.Run(body[:min(len(body), 60)], func(t *testing.T) {
			f := &backupFake{}
			r := backupRouter(f)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/internal/jobs/db-backup", strings.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+testJobWebhookSecret)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != 400 || f.calls != 0 {
				t.Fatalf("status %d calls %d", w.Code, f.calls)
			}
		})
	}
	for _, method := range []string{http.MethodPost, http.MethodGet} {
		f := &backupFake{}
		r := backupRouter(f)
		path := "/api/internal/jobs/db-backup"
		if method == http.MethodGet {
			path += "/" + uuid.NewString()
		}
		req := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(`{"type":"daily"}`))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 401 || f.calls != 0 {
			t.Fatalf("unauthenticated %s: %d", method, w.Code)
		}
	}
}
func TestBackupAdmissionAndStatus(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want int
	}{{nil, 202}, {backup.ErrDisabled, 503}, {errors.New("password=secret"), 500}} {
		f := &backupFake{id: uuid.NewString(), err: tc.err}
		r := backupRouter(f)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/internal/jobs/db-backup", strings.NewReader(`{"type":"weekly"}`))
		req.Header.Set("Authorization", "Bearer "+testJobWebhookSecret)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want || strings.Contains(w.Body.String(), "secret") {
			t.Fatalf("%d: %s", w.Code, w.Body.String())
		}
		if tc.err == nil && !strings.Contains(w.Body.String(), f.id) {
			t.Fatal("missing persisted run ID")
		}
	}
	for _, tc := range []struct {
		id   string
		err  error
		want int
	}{{"invalid", nil, 400}, {uuid.NewString(), backup.ErrNotFound, 404}, {uuid.NewString(), nil, 200}} {
		f := &backupFake{id: tc.id, err: tc.err}
		r := backupRouter(f)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/internal/jobs/db-backup/"+tc.id, nil)
		req.Header.Set("Authorization", "Bearer "+testJobWebhookSecret)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Fatalf("status %d want %d", w.Code, tc.want)
		}
	}
}
