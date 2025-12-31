package middleware

import (
	"go_boilerplate/internal/shared/config"

	"github.com/gofiber/fiber/v2"
)

// CORS returns a CORS middleware
func CORS(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Allow all origins in development
		if cfg.Server.IsDevelopment() {
			c.Set("Access-Control-Allow-Origin", "*")
		} else {
			// In production, you should specify allowed origins
			origin := c.Get("Origin")
			// You can add your own logic here to validate origin
			c.Set("Access-Control-Allow-Origin", origin)
		}

		// Allow methods
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")

		// Allow headers
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")

		// Allow credentials
		c.Set("Access-Control-Allow-Credentials", "true")

		// Max age
		c.Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
