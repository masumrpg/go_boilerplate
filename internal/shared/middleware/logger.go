package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// HTTPLogger is a middleware that logs HTTP requests
func HTTPLogger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		method := c.Method()
		path := c.Path()
		status := c.Response().StatusCode()
		ip := c.IP()
		userAgent := c.Get("User-Agent")

		// Create log entry
		entry := logger.WithFields(logrus.Fields{
			"method":     method,
			"path":       path,
			"status":     status,
			"latency":    latency.String(),
			"ip":         ip,
			"user_agent": userAgent,
		})

		// Log based on status code
		if err != nil {
			entry.Error("Request completed with error")
		} else if status >= 500 {
			entry.Error("Server error")
		} else if status >= 400 {
			entry.Warn("Client error")
		} else if status >= 300 {
			entry.Info("Redirect")
		} else {
			entry.Info("Request completed")
		}

		return err
	}
}
