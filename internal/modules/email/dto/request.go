package dto

// SendEmailRequest represents an email sending request
type SendEmailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

// SendWelcomeEmailRequest represents a welcome email request
type SendWelcomeEmailRequest struct {
	To   string `json:"to" validate:"required,email"`
	Name string `json:"name" validate:"required"`
}

// SendPasswordResetRequest represents a password reset email request
type SendPasswordResetRequest struct {
	To           string `json:"to" validate:"required,email"`
	ResetToken   string `json:"reset_token" validate:"required"`
	ResetLink    string `json:"reset_link" validate:"required"`
}
