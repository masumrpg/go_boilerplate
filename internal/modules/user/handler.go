package user

import (
	"strconv"

	sharedmiddleware "go_boilerplate/internal/shared/middleware"
	"go_boilerplate/internal/shared/utils"
	userdto "go_boilerplate/internal/modules/user/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler defines the interface for user HTTP handlers
type UserHandler interface {
	GetUser(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetCurrentUser(c *fiber.Ctx) error
}

// userHandler implements UserHandler interface
type userHandler struct {
	service UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service UserService) UserHandler {
	return &userHandler{service: service}
}

// GetUser gets a user by ID
func (h *userHandler) GetUser(c *fiber.Ctx) error {
	// Get user ID from params
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Get user
	user, err := h.service.GetProfile(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user, "User retrieved successfully")
}

// GetUsers gets all users with pagination
func (h *userHandler) GetUsers(c *fiber.Ctx) error {
	// Get pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Default values
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get users
	users, err := h.service.GetAll(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, users, "Users retrieved successfully")
}

// CreateUser creates a new user
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	// Get validated body from context
	validatedBody := c.Locals("validatedBody").(*userdto.CreateUserRequest)

	// Create user
	user, err := h.service.CreateUser(validatedBody)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create user", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, user, "User created successfully")
}

// UpdateUser updates a user
func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	// Get user ID from params
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Get authenticated user ID from context
	authUserIDStr, ok := sharedmiddleware.GetUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	authUserID, err := uuid.Parse(authUserIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid auth user ID", err)
	}

	// Check if user is updating their own profile
	if authUserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only update your own profile", nil)
	}

	// Get validated body from context
	validatedBody := c.Locals("validatedBody").(*userdto.UpdateUserRequest)

	// Update user
	user, err := h.service.UpdateUser(userID, validatedBody)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update user", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user, "User updated successfully")
}

// DeleteUser deletes a user
func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	// Get user ID from params
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Delete user
	if err := h.service.DeleteUser(userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete user", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "User deleted successfully")
}

// GetCurrentUser gets the authenticated user's profile
func (h *userHandler) GetCurrentUser(c *fiber.Ctx) error {
	// Get user ID from context
	userIDStr, ok := sharedmiddleware.GetUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Get user
	user, err := h.service.GetProfile(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user, "User profile retrieved successfully")
}
