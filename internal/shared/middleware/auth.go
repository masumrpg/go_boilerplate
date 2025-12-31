package middleware

import (
	"strings"

	"go_boilerplate/internal/shared/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth returns a JWT authentication middleware
func JWTAuth(cfg *config.Config) fiber.Handler {
	// Using Fiber's contrib JWT middleware
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(cfg.JWT.Secret)},
		ErrorHandler: jwtError,
	})
}

// jwtError handles JWT errors
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing or malformed JWT",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"error":   "Invalid or expired JWT",
	})
}

// OptionalAuth is a middleware that checks for JWT but doesn't require it
// If JWT is present and valid, it sets the user context
// If JWT is missing, it continues without setting user context
func OptionalAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		// No authorization header, continue without auth
		if authHeader == "" {
			return c.Next()
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Create JWT middleware instance to validate
		jwtMiddleware := jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWT.Secret)},
			ContextKey: "jwt",
		})

		// Create a fake context to test the token
		app := fiber.New()
		app.Use(jwtMiddleware)

		// Try to parse and validate the token
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())

		token, err := parser.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or expired JWT",
			})
		}

		// Token is valid, store it in context
		c.Locals("jwt", token)
		c.Locals("user", token.Claims.(jwt.MapClaims))

		return c.Next()
	}
}

// GetUserIDFromContext extracts user ID from JWT context
func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	user := c.Locals("user")
	if user == nil {
		return "", false
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", false
	}

	return userID, true
}

// GetEmailFromContext extracts email from JWT context
func GetEmailFromContext(c *fiber.Ctx) (string, bool) {
	user := c.Locals("user")
	if user == nil {
		return "", false
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", false
	}

	return email, true
}
