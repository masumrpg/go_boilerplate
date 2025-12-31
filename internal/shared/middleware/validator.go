package middleware

import (
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
)

// BodyValidator validates request body against a struct
func BodyValidator(v interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse body
		if err := c.BodyParser(v); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to parse request body",
			})
		}

		// Validate struct
		validator := utils.NewValidator()
		if err := validator.ValidateStruct(v); err != nil {
			errors := utils.GetValidationErrors(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Validation failed",
				"details": errors,
			})
		}

		// Store validated body in context for later use
		c.Locals("validatedBody", v)

		return c.Next()
	}
}
