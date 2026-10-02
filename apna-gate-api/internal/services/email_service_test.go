package service

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"go-server/internal/config"
)

func TestFormatFromAddress(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{
			name:  "Apna Gate",
			email: "no-reply@apnagate.org",
			want:  "Apna Gate <no-reply@apnagate.org>",
		},
		{
			name:  "",
			email: "no-reply@apnagate.org",
			want:  "no-reply@apnagate.org",
		},
		{
			name:  "  Apna Gate  ",
			email: "  no-reply@apnagate.org  ",
			want:  "Apna Gate <no-reply@apnagate.org>",
		},
	}

	for _, tt := range tests {
		got := formatFromAddress(tt.name, tt.email)
		if got != tt.want {
			t.Fatalf("formatFromAddress(%q, %q) = %q, want %q", tt.name, tt.email, got, tt.want)
		}
	}
}

func TestValidateConfigRequiresResendAPIKey(t *testing.T) {
	cfg := &config.Config{
		EmailFrom: "no-reply@apnagate.org",
	}

	err := validateConfig(cfg)
	if err == nil || !strings.Contains(err.Error(), "RESEND_API_KEY is required") {
		t.Fatalf("validateConfig() error = %v, want missing RESEND_API_KEY", err)
	}
}

func TestValidateConfigRequiresEmailFrom(t *testing.T) {
	cfg := &config.Config{
		ResendAPIKey: "re_test",
	}

	err := validateConfig(cfg)
	if err == nil || !strings.Contains(err.Error(), "EMAIL_FROM is required") {
		t.Fatalf("validateConfig() error = %v, want missing EMAIL_FROM", err)
	}
}

func TestValidateConfigAcceptsResendConfig(t *testing.T) {
	cfg := &config.Config{
		ResendAPIKey: "re_test",
		EmailFrom:    "no-reply@apnagate.org",
	}

	if err := validateConfig(cfg); err != nil {
		t.Fatalf("validateConfig() error = %v, want nil", err)
	}
}

func TestEmailTemplatesRenderWithApnaGateBranding(t *testing.T) {
	templates, err := loadTemplates()
	if err != nil {
		t.Fatalf("loadTemplates() error = %v", err)
	}

	tests := []struct {
		name string
		tmpl *template.Template
		data any
		want []string
	}{
		{
			name: "verification",
			tmpl: templates.otp,
			data: struct {
				Name, OTP, AppURL string
				Year              int
			}{"Asha", "123456", "https://apnagate.app", 2026},
			want: []string{"Apna Gate", "Asha", "123456", "Expires in 5 minutes"},
		},
		{
			name: "forgot password",
			tmpl: templates.forgetPassword,
			data: struct {
				Name, OTP, AppURL string
				Year              int
			}{"Asha", "654321", "https://apnagate.app", 2026},
			want: []string{"Apna Gate", "Asha", "654321", "Expires in 5 minutes"},
		},
		{
			name: "password reset link",
			tmpl: templates.passwordReset,
			data: struct {
				Name, ResetLink, AppURL string
				Year                    int
			}{"Asha", "https://apnagate.app/reset?token=safe", "https://apnagate.app", 2026},
			want: []string{"Apna Gate", "Asha", "Reset my password", "token=safe"},
		},
		{
			name: "welcome",
			tmpl: templates.welcome,
			data: struct {
				Name, AppURL string
				Year         int
			}{"Asha", "https://apnagate.app", 2026},
			want: []string{"Welcome to Apna Gate", "Asha", "Get started"},
		},
		{
			name: "monthly report",
			tmpl: templates.monthlyReport,
			data: struct {
				SocietyName, Period, Filename, AppURL string
				Year                                  int
			}{"Sunrise <Residency>", "August 2026", "visitor-report.csv", "https://apnagate.app", 2026},
			want: []string{"Sunrise &lt;Residency&gt;", "August 2026", "visitor-report.csv", "Handle resident data carefully"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := tt.tmpl.Execute(&output, tt.data); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			body := output.String()
			for _, want := range tt.want {
				if !strings.Contains(body, want) {
					t.Errorf("rendered template does not contain %q", want)
				}
			}
			for _, unwanted := range []string{"Your App", "Your Platform", "ï¿½", "â"} {
				if strings.Contains(body, unwanted) {
					t.Errorf("rendered template contains stale or malformed text %q", unwanted)
				}
			}
		})
	}
}
