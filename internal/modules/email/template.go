package email

import (
	"fmt"
)

// EmailTemplate represents an email template
type EmailTemplate struct {
	Subject string
	Body    string
}

// WelcomeEmailTemplate returns the welcome email template
func WelcomeEmailTemplate(name string) EmailTemplate {
	subject := "Welcome to Our Platform!"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		.container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; }
		.footer { background-color: #f1f1f1; padding: 10px; text-align: center; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Welcome!</h1>
		</div>
		<div class="content">
			<p>Hello %s,</p>
			<p>Welcome to our platform! We're excited to have you on board.</p>
			<p>If you have any questions, feel free to reach out to us.</p>
			<p>Best regards,<br>The Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2024 Our Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`, name)

	return EmailTemplate{
		Subject: subject,
		Body:     body,
	}
}

// PasswordResetEmailTemplate returns the password reset email template
func PasswordResetEmailTemplate(resetLink string) EmailTemplate {
	subject := "Password Reset Request"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		.container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; }
		.button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; display: inline-block; margin: 20px 0; }
		.footer { background-color: #f1f1f1; padding: 10px; text-align: center; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Password Reset</h1>
		</div>
		<div class="content">
			<p>You requested a password reset.</p>
			<p>Click the button below to reset your password:</p>
			<center><a href="%s" class="button">Reset Password</a></center>
			<p>Or copy and paste this link into your browser:</p>
			<p>%s</p>
			<p>This link will expire in 1 hour.</p>
			<p>If you didn't request this, please ignore this email.</p>
		</div>
		<div class="footer">
			<p>&copy; 2024 Our Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`, resetLink, resetLink)

	return EmailTemplate{
		Subject: subject,
		Body:     body,
	}
}
