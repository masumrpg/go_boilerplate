package role

import (
	"go_boilerplate/internal/modules/role/dto"
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoleHandler defines the interface for role HTTP handlers
type RoleHandler interface {
	GetRoles(c *fiber.Ctx) error
	GetRole(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
}

// roleHandler implements RoleHandler interface
type roleHandler struct {
	service RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(service RoleService) RoleHandler {
	return &roleHandler{service: service}
}

// GetRoles gets all roles with pagination
func (h *roleHandler) GetRoles(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get roles
	response, err := h.service.GetAllRoles(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get roles", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "Roles retrieved successfully")
}

// GetRole gets a role by ID
func (h *roleHandler) GetRole(c *fiber.Ctx) error {
	// Parse role ID
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", err)
	}

	// Get role
	role, err := h.service.GetRole(roleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Role not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, role, "Role retrieved successfully")
}

// CreateRole creates a new role
func (h *roleHandler) CreateRole(c *fiber.Ctx) error {
	// Get validated body
	validatedBody := c.Locals("validatedBody").(*dto.CreateRoleRequest)

	// Create role
	role, err := h.service.CreateRole(validatedBody)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, role, "Role created successfully")
}

// UpdateRole updates a role
func (h *roleHandler) UpdateRole(c *fiber.Ctx) error {
	// Parse role ID
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", err)
	}

	// Get validated body
	validatedBody := c.Locals("validatedBody").(*dto.UpdateRoleRequest)

	// Update role
	role, err := h.service.UpdateRole(roleID, validatedBody)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, role, "Role updated successfully")
}

// DeleteRole deletes a role
func (h *roleHandler) DeleteRole(c *fiber.Ctx) error {
	// Parse role ID
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID", err)
	}

	// Delete role
	if err := h.service.DeleteRole(roleID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete role", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Role deleted successfully")
}
