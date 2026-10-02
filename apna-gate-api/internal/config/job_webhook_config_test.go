package config

import (
	"strings"
	"testing"
	"time"
)

func TestJobWebhookSecretValidation(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{name: "missing", wantErr: true},
		{name: "short after trimming", secret: "   too-short   ", wantErr: true},
		{name: "valid", secret: strings.Repeat("a", 32)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validJobWebhookConfig()
			cfg.JobWebhookSecret = tt.secret
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validJobWebhookConfig() *Config {
	return &Config{
		Environment: "development", Port: "8080", DBName: "test", OTPExpiry: time.Minute,
		Auth: AuthConfig{OnboardingExpiry: time.Minute}, JobTimezone: "Asia/Kolkata",
		VisitorEntryRetentionDays: 90, VisitorInviteRetentionDays: 90, FlatInviteRetentionDays: 90,
		NotificationRetentionDays: 30, CleanupBatchSize: 500, VisitorReportBatchSize: 50,
		VisitorReportProcessingLease: 15 * time.Minute,
	}
}

func TestVisitorReportResendIsDevelopmentOnly(t *testing.T) {
	cfg := validJobWebhookConfig()
	cfg.JobWebhookSecret = strings.Repeat("a", 32)
	cfg.VisitorReportAllowResend = true
	if err := cfg.Validate(); err != nil {
		t.Fatalf("development rejected resend: %v", err)
	}
	for _, environment := range []string{"production", "prod", "test"} {
		cfg.Environment = environment
		if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "VISITOR_REPORT_ALLOW_RESEND") {
			t.Fatalf("environment %q accepted resend: %v", environment, err)
		}
	}
}
