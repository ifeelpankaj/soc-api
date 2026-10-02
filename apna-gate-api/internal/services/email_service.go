package service

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"go-server/internal/config"

	"github.com/resend/resend-go/v3"
)

//go:embed templates/*.html
var templateFS embed.FS

// EmailService handles all email operations
type EmailService interface {
	SendOTP(ctx context.Context, to, otp, name string) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendPasswordResetEmail(ctx context.Context, to, resetLink, name string) error
	SendForgetPasswordEmail(ctx context.Context, to, otp, name string) error
	SendMonthlyVisitorReport(ctx context.Context, recipients []string, societyName string, reportMonth time.Time, filename string, content []byte, idempotencyKey string) (string, error)
	Close() error
}

type emailService struct {
	config    *emailConfig
	templates *emailTemplates
	client    *resend.Client
}

type emailConfig struct {
	apiKey    string
	fromEmail string
	fromName  string
	appURL    string
}

type emailTemplates struct {
	otp            *template.Template
	welcome        *template.Template
	passwordReset  *template.Template
	forgetPassword *template.Template
	monthlyReport  *template.Template
}

// NewEmailService creates a new email service instance
func NewEmailService(cfg *config.Config) (EmailService, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid email configuration: %w", err)
	}

	emailCfg := &emailConfig{
		apiKey:    cfg.ResendAPIKey,
		fromEmail: cfg.EmailFrom,
		fromName:  cfg.EmailFromName,
		appURL:    cfg.PublicAppURL,
	}

	templates, err := loadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	client := resend.NewClient(emailCfg.apiKey)

	return &emailService{
		config:    emailCfg,
		templates: templates,
		client:    client,
	}, nil
}

// validateConfig ensures all required email configuration is present
func validateConfig(cfg *config.Config) error {
	if strings.TrimSpace(cfg.ResendAPIKey) == "" {
		return fmt.Errorf("RESEND_API_KEY is required")
	}
	if strings.TrimSpace(cfg.EmailFrom) == "" {
		return fmt.Errorf("EMAIL_FROM is required")
	}
	return nil
}

// SendOTP sends an OTP verification email
func (s *emailService) SendOTP(ctx context.Context, to, otp, name string) error {
	data := struct {
		Name   string
		OTP    string
		AppURL string
		Year   int
	}{
		Name:   name,
		OTP:    otp,
		AppURL: s.config.appURL,
		Year:   time.Now().Year(),
	}

	body, err := s.renderTemplate(s.templates.otp, data)
	if err != nil {
		return fmt.Errorf("failed to render OTP template: %w", err)
	}

	return s.sendEmail(ctx, to, "Your Apna Gate verification code", body)
}

// SendForgetPasswordEmail sends a forget-password OTP email
func (s *emailService) SendForgetPasswordEmail(ctx context.Context, to, otp, name string) error {
	data := struct {
		Name   string
		OTP    string
		AppURL string
		Year   int
	}{
		Name:   name,
		OTP:    otp,
		AppURL: s.config.appURL,
		Year:   time.Now().Year(),
	}

	body, err := s.renderTemplate(s.templates.forgetPassword, data)
	if err != nil {
		return fmt.Errorf("failed to render forget password template: %w", err)
	}

	return s.sendEmail(ctx, to, "Reset your Apna Gate password", body)
}

// SendWelcomeEmail sends a welcome email after successful verification
func (s *emailService) SendWelcomeEmail(ctx context.Context, to, name string) error {
	data := struct {
		Name   string
		AppURL string
		Year   int
	}{
		Name:   name,
		AppURL: s.config.appURL,
		Year:   time.Now().Year(),
	}

	body, err := s.renderTemplate(s.templates.welcome, data)
	if err != nil {
		return fmt.Errorf("failed to render welcome template: %w", err)
	}

	return s.sendEmail(ctx, to, "Welcome to Apna Gate", body)
}

// SendPasswordResetEmail sends a password reset link email
func (s *emailService) SendPasswordResetEmail(ctx context.Context, to, resetLink, name string) error {
	data := struct {
		Name      string
		ResetLink string
		AppURL    string
		Year      int
	}{
		Name:      name,
		ResetLink: resetLink,
		AppURL:    s.config.appURL,
		Year:      time.Now().Year(),
	}

	body, err := s.renderTemplate(s.templates.passwordReset, data)
	if err != nil {
		return fmt.Errorf("failed to render password reset template: %w", err)
	}

	return s.sendEmail(ctx, to, "Reset your Apna Gate password", body)
}

// SendMonthlyVisitorReport sends a society-scoped CSV using the existing Resend client.
func (s *emailService) SendMonthlyVisitorReport(
	ctx context.Context,
	recipients []string,
	societyName string,
	reportMonth time.Time,
	filename string,
	content []byte,
	idempotencyKey string,
) (string, error) {
	if len(recipients) == 0 {
		return "", fmt.Errorf("at least one visitor report recipient is required")
	}
	period := reportMonth.Format("January 2006")
	data := struct {
		SocietyName string
		Period      string
		Filename    string
		AppURL      string
		Year        int
	}{
		SocietyName: societyName,
		Period:      period,
		Filename:    filename,
		AppURL:      s.config.appURL,
		Year:        time.Now().Year(),
	}
	body, err := s.renderTemplate(s.templates.monthlyReport, data)
	if err != nil {
		return "", fmt.Errorf("failed to render monthly visitor report template: %w", err)
	}
	params := &resend.SendEmailRequest{
		From:    formatFromAddress(s.config.fromName, s.config.fromEmail),
		To:      recipients,
		Subject: fmt.Sprintf("%s visitor report | %s", societyName, period),
		Html:    body,
		Attachments: []*resend.Attachment{{
			Content: content, Filename: filename, ContentType: "text/csv; charset=utf-8",
		}},
	}

	const maxRetries = 3
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", fmt.Errorf("visitor report email cancelled: %w", ctx.Err())
			case <-time.After(time.Second * time.Duration(attempt)):
			}
		}
		response, err := s.client.Emails.SendWithOptions(ctx, params, &resend.SendEmailOptions{IdempotencyKey: idempotencyKey})
		if err == nil {
			return response.Id, nil
		}
		lastErr = err
	}
	return "", fmt.Errorf("failed to send visitor report after %d attempts: %w", maxRetries, lastErr)
}

// sendEmail sends an email with retry logic and context support
func (s *emailService) sendEmail(ctx context.Context, to, subject, body string) error {
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("email sending cancelled: %w", ctx.Err())
			case <-time.After(time.Second * time.Duration(attempt)):
			}
		}

		if err := s.sendEmailAttempt(ctx, to, subject, body); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

// sendEmailAttempt performs a single email send attempt
func (s *emailService) sendEmailAttempt(ctx context.Context, to, subject, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	params := &resend.SendEmailRequest{
		From:    formatFromAddress(s.config.fromName, s.config.fromEmail),
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	if _, err := s.client.Emails.SendWithContext(ctx, params); err != nil {
		return fmt.Errorf("resend send failed: %w", err)
	}

	return nil
}

func formatFromAddress(name, email string) string {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

// renderTemplate renders an email template with data
func (s *emailService) renderTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}
	return buf.String(), nil
}

// loadTemplates loads all email templates from embedded filesystem
func loadTemplates() (*emailTemplates, error) {
	otp, err := template.ParseFS(templateFS, "templates/otp_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse OTP template: %w", err)
	}

	welcome, err := template.ParseFS(templateFS, "templates/welcome_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse welcome template: %w", err)
	}

	reset, err := template.ParseFS(templateFS, "templates/reset_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse password reset template: %w", err)
	}

	forgetPassword, err := template.ParseFS(templateFS, "templates/forget_password_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse forget password template: %w", err)
	}

	monthlyReport, err := template.ParseFS(templateFS, "templates/monthly_report_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse monthly visitor report template: %w", err)
	}

	return &emailTemplates{
		otp:            otp,
		welcome:        welcome,
		passwordReset:  reset,
		forgetPassword: forgetPassword,
		monthlyReport:  monthlyReport,
	}, nil
}

// Close gracefully shuts down the email service
func (s *emailService) Close() error {
	return nil
}
