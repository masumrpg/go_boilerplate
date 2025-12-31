package email

import (
	"time"

	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/modules/email/dto"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendWelcomeEmail(to, name string) error
	SendPasswordResetEmail(to, resetLink string) error
}

// emailService implements EmailService interface
type emailService struct {
	cfg    *config.Config
	dialer *gomail.Dialer
	logger *logrus.Logger
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config, logger *logrus.Logger) EmailService {
	dialer := gomail.NewDialer(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPassword,
	)

	return &emailService{
		cfg:    cfg,
		dialer: dialer,
		logger: logger,
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
	// Get template
	emailTemplate := WelcomeEmailTemplate(name)

	// Send email
	return s.SendEmail(to, emailTemplate.Subject, emailTemplate.Body)
}

// SendPasswordResetEmail sends a password reset email
func (s *emailService) SendPasswordResetEmail(to, resetLink string) error {
	// Get template
	emailTemplate := PasswordResetEmailTemplate(resetLink)

	// Send email
	return s.SendEmail(to, emailTemplate.Subject, emailTemplate.Body)
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
