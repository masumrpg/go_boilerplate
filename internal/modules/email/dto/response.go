package dto

import "time"

// EmailResponse represents an email response
type EmailResponse struct {
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
}
