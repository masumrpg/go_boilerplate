package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"go_boilerplate/internal/modules/email/dto"
	"go_boilerplate/internal/shared/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

//go:embed templates/*.html
var templatesFS embed.FS

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendWelcomeEmail(to, name string) error
	SendPasswordResetEmail(to, resetLink string) error
	SendVerificationEmail(to, code string) error
	SendTwoFactorEmail(to, code string) error
}

// emailService implements EmailService interface
type emailService struct {
	cfg       *config.Config
	dialer    *gomail.Dialer
	logger    *logrus.Logger
	templates *template.Template
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config, logger *logrus.Logger) EmailService {
	dialer := gomail.NewDialer(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPassword,
	)

	// Parse templates from embedded FS
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		logger.Errorf("Failed to parse email templates: %v", err)
	}

	return &emailService{
		cfg:       cfg,
		dialer:    dialer,
		logger:    logger,
		templates: tmpl,
	}
}

// SendEmail sends an email
func (s *emailService) SendEmail(to, subject, body string) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Email.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send email
	if err := s.dialer.DialAndSend(m); err != nil {
		s.logger.Errorf("Failed to send email to %s: %v", to, err)
		return err
	}

	s.logger.Infof("Email sent successfully to %s", to)
	return nil
}

// SendWelcomeEmail sends a welcome email
func (s *emailService) SendWelcomeEmail(to, name string) error {
	body, err := s.renderTemplate("welcome.html", map[string]interface{}{
		"Name": name,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(to, "Welcome to Our Platform!", body)
}

// SendPasswordResetEmail sends a password reset email
func (s *emailService) SendPasswordResetEmail(to, resetLink string) error {
	body, err := s.renderTemplate("password_reset.html", map[string]interface{}{
		"ResetLink": resetLink,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(to, "Password Reset Request", body)
}

// SendVerificationEmail sends an account verification email
func (s *emailService) SendVerificationEmail(to, code string) error {
	body, err := s.renderTemplate("verification_code.html", map[string]interface{}{
		"Code": code,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(to, "Verify Your Account", body)
}

// SendTwoFactorEmail sends a 2FA verification email
func (s *emailService) SendTwoFactorEmail(to, code string) error {
	body, err := s.renderTemplate("2fa_code.html", map[string]interface{}{
		"Code": code,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(to, "Your Login Verification Code", body)
}

// renderTemplate renders an HTML template with data
func (s *emailService) renderTemplate(name string, data interface{}) (string, error) {
	if s.templates == nil {
		return "", fmt.Errorf("templates not initialized")
	}

	var buf bytes.Buffer
	if err := s.templates.ExecuteTemplate(&buf, name, data); err != nil {
		s.logger.Errorf("Failed to render template %s: %v", name, err)
		return "", err
	}

	return buf.String(), nil
}

// BuildEmailResponse creates an email response
func BuildEmailResponse(to, subject string) *dto.EmailResponse {
	return &dto.EmailResponse{
		Message: "Email sent successfully",
		SentAt:  time.Now(),
		To:      to,
		Subject: subject,
	}
}
