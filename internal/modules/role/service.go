package role

import (
	"errors"
	"math"

	"go_boilerplate/internal/modules/role/dto"
	"go_boilerplate/internal/shared/utils"

	"github.com/google/uuid"
)

// RoleService defines the interface for role business logic
type RoleService interface {
	GetRole(roleID uuid.UUID) (*dto.RoleResponse, error)
	GetRoleBySlug(slug string) (*dto.RoleResponse, error)
	GetAllRoles(page, limit int) (*dto.RolesResponse, error)
	CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	UpdateRole(roleID uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	DeleteRole(roleID uuid.UUID) error
	SeedInitialRoles() error
}

// roleService implements RoleService interface
type roleService struct {
	repo RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(repo RoleRepository) RoleService {
	return &roleService{repo: repo}
}

// GetRole gets a role by ID
func (s *roleService) GetRole(roleID uuid.UUID) (*dto.RoleResponse, error) {
	roleModel, err := s.repo.FindByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	response := s.modelToResponse(roleModel)
	return &response, nil
}

// GetRoleBySlug gets a role by slug
func (s *roleService) GetRoleBySlug(slug string) (*dto.RoleResponse, error) {
	roleModel, err := s.repo.FindBySlug(slug)
	if err != nil || roleModel == nil {
		return nil, errors.New("role not found")
	}

	response := s.modelToResponse(roleModel)
	return &response, nil
}

// GetAllRoles gets all roles with pagination
func (s *roleService) GetAllRoles(page, limit int) (*dto.RolesResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Find roles
	roles, total, err := s.repo.FindAll(offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, roleModel := range roles {
		roleResponses[i] = s.modelToResponse(&roleModel)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.RolesResponse{
		Roles: roleResponses,
		Meta: utils.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

// CreateRole creates a new role
func (s *roleService) CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// Check if slug already exists
	exists, err := s.repo.ExistsBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("role with this slug already exists")
	}

	// Check if name already exists
	exists, err = s.repo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("role with this name already exists")
	}

	// Create role model
	roleModel := &Role{
		Name:        req.Name,
		Slug:        req.Slug,
		Permissions: StringSlice(req.Permissions),
		Description: req.Description,
	}

	// Save role
	if err := s.repo.Create(roleModel); err != nil {
		return nil, err
	}

	response := s.modelToResponse(roleModel)
	return &response, nil
}

// UpdateRole updates a role
func (s *roleService) UpdateRole(roleID uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	// Find role
	roleModel, err := s.repo.FindByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if new name already exists (excluding current role)
		existingRole, _ := s.repo.FindBySlug(roleModel.Slug)
		if existingRole != nil && existingRole.ID != roleID {
			return nil, errors.New("role with this name already exists")
		}
		roleModel.Name = req.Name
	}

	if len(req.Permissions) > 0 {
		roleModel.Permissions = StringSlice(req.Permissions)
	}

	if req.Description != "" {
		roleModel.Description = req.Description
	}

	// Save changes
	if err := s.repo.Update(roleModel); err != nil {
		return nil, err
	}

	response := s.modelToResponse(roleModel)
	return &response, nil
}

// DeleteRole deletes a role
func (s *roleService) DeleteRole(roleID uuid.UUID) error {
	// Check if role exists
	_, err := s.repo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Delete role
	if err := s.repo.Delete(roleID); err != nil {
		return err
	}

	return nil
}

// SeedInitialRoles seeds the database with initial roles
func (s *roleService) SeedInitialRoles() error {
	initialRoles := []*dto.CreateRoleRequest{
		{
			Name:        "SuperAdmin",
			Slug:        "super_admin",
			Permissions: []string{"*"},
			Description: "Full system access with all permissions",
		},
		{
			Name:        "Admin",
			Slug:        "admin",
			Permissions: []string{
				"users.create",
				"users.read",
				"users.update",
				"users.delete",
				"roles.read",
				"roles.assign",
			},
			Description: "Administrative access for user and role management",
		},
		{
			Name:        "User",
			Slug:        "user",
			Permissions: []string{
				"users.read",
				"users.update",
			},
			Description: "Standard user access with self-profile management",
		},
	}

	for _, roleReq := range initialRoles {
		existing, _ := s.repo.FindBySlug(roleReq.Slug)
		if existing == nil {
			roleModel := &Role{
				Name:        roleReq.Name,
				Slug:        roleReq.Slug,
				Permissions: StringSlice(roleReq.Permissions),
				Description: roleReq.Description,
			}
			if err := s.repo.Create(roleModel); err != nil {
				return err
			}
		}
	}

	return nil
}

// modelToResponse converts Role model to RoleResponse
func (s *roleService) modelToResponse(role *Role) dto.RoleResponse {
	return dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Slug:        role.Slug,
		Permissions: []string(role.Permissions),
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
