package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "go-server/internal/handlers/v1"
	"go-server/internal/jobs"

	"github.com/gin-gonic/gin"
)

const testJobWebhookSecret = "0123456789abcdef0123456789abcdef"

type jobTriggererFake struct {
	names []string
}

func (f *jobTriggererFake) Trigger(name string) (jobs.Trigger, error) {
	f.names = append(f.names, name)
	return jobs.Trigger{Job: name, TriggerID: "trigger-123", Status: "accepted"}, nil
}

func TestInternalJobRoutesRequireBearerSecret(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
	}{
		{name: "missing"},
		{name: "malformed", authorization: "Token " + testJobWebhookSecret},
		{name: "wrong", authorization: "Bearer 0123456789abcdef0123456789abcdeg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triggerer := &jobTriggererFake{}
			router := internalJobTestRouter(triggerer)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/internal/jobs/maintenance-reminders", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", response.Code)
			}
			if len(triggerer.names) != 0 {
				t.Fatalf("triggered jobs = %#v, want none", triggerer.names)
			}
		})
	}
}

func TestInternalJobRoutesTriggerOnlySelectedJob(t *testing.T) {
	tests := []struct {
		path string
		job  string
	}{
		{path: "/api/internal/jobs/cleanup", job: jobs.JobCleanup},
		{path: "/api/internal/jobs/monthly-visitor-report", job: jobs.JobMonthlyVisitorReport},
		{path: "/api/internal/jobs/maintenance-billing", job: jobs.JobMaintenanceBilling},
		{path: "/api/internal/jobs/maintenance-reminders", job: jobs.JobMaintenanceReminders},
	}
	for _, tt := range tests {
		t.Run(tt.job, func(t *testing.T) {
			triggerer := &jobTriggererFake{}
			router := internalJobTestRouter(triggerer)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+testJobWebhookSecret)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != http.StatusAccepted {
				t.Fatalf("status = %d, want 202; body=%s", response.Code, response.Body.String())
			}
			if len(triggerer.names) != 1 || triggerer.names[0] != tt.job {
				t.Fatalf("triggered jobs = %#v, want %q", triggerer.names, tt.job)
			}
			var body jobs.Trigger
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Job != tt.job || body.TriggerID != "trigger-123" || body.Status != "accepted" {
				t.Fatalf("response = %+v", body)
			}
		})
	}
}

func internalJobTestRouter(triggerer *jobTriggererFake) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupInternalJobRoutes(router.Group("/api/internal"), handlers.NewJobWebhookHandler(triggerer), testJobWebhookSecret)
	return router
}
