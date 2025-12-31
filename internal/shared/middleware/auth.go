package middleware

import (
	"strings"

	"go_boilerplate/internal/shared/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
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

		token, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
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

// getClaims extracts JWT claims from context handling both *jwt.Token and jwt.MapClaims
func getClaims(c *fiber.Ctx) (jwt.MapClaims, bool) {
	user := c.Locals("user")
	if user == nil {
		return nil, false
	}

	// Case 1: stored as *jwt.Token (standard behavior of gofiber/contrib/jwt)
	if token, ok := user.(*jwt.Token); ok {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims, true
		}
	}

	// Case 2: stored directly as jwt.MapClaims (e.g. from OptionalAuth or custom override)
	if claims, ok := user.(jwt.MapClaims); ok {
		return claims, true
	}

	return nil, false
}

// GetUserIDFromContext extracts user ID from JWT context
func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	claims, ok := getClaims(c)
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
	claims, ok := getClaims(c)
	if !ok {
		return "", false
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", false
	}

	return email, true
}

// RequireRole checks if the authenticated user has one of the required roles
func RequireRole(cfg *config.Config, roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := getClaims(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Unauthorized",
			})
		}

		// Extract role_slug from claims
		userRoleSlug, ok := claims["role_slug"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Role information not found in token",
			})
		}

		// Check if user has any of the required roles
		for _, role := range roles {
			if userRoleSlug == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success":       false,
			"error":         "Insufficient permissions",
			"required_roles": roles,
			"user_role":      userRoleSlug,
		})
	}
}

// RequirePermission checks if the authenticated user has a specific permission
func RequirePermission(cfg *config.Config, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := getClaims(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Unauthorized",
			})
		}

		// Get permissions from claims
		permissionsInterface, ok := claims["permissions"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Permissions not found in token",
			})
		}

		// Convert to string slice
		permissions := make([]string, len(permissionsInterface))
		for i, p := range permissionsInterface {
			if str, ok := p.(string); ok {
				permissions[i] = str
			}
		}

		// Check for wildcard permission
		for _, p := range permissions {
			if p == "*" {
				return c.Next()
			}
		}

		// Check specific permission
		for _, p := range permissions {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success":   false,
			"error":     "Insufficient permissions",
			"required":  permission,
		})
	}
}

// GetRoleSlugFromContext extracts role slug from JWT context
func GetRoleSlugFromContext(c *fiber.Ctx) (string, bool) {
	claims, ok := getClaims(c)
	if !ok {
		return "", false
	}

	roleSlug, ok := claims["role_slug"].(string)
	return roleSlug, ok
}

// GetPermissionsFromContext extracts permissions from JWT context
func GetPermissionsFromContext(c *fiber.Ctx) ([]string, bool) {
	claims, ok := getClaims(c)
	if !ok {
		return nil, false
	}

	permissionsInterface, ok := claims["permissions"].([]interface{})
	if !ok {
		return nil, false
	}

	permissions := make([]string, len(permissionsInterface))
	for i, p := range permissionsInterface {
		if str, ok := p.(string); ok {
			permissions[i] = str
		}
	}

	return permissions, true
}
