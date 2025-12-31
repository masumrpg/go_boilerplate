package user

import (
	"errors"
	"math"

	"go_boilerplate/internal/shared/utils"
	"go_boilerplate/internal/modules/role"
	userdto "go_boilerplate/internal/modules/user/dto"

	"github.com/google/uuid"
)

// UserService defines the interface for user business logic
type UserService interface {
	GetProfile(userID uuid.UUID) (*userdto.UserResponse, error)
	GetProfileWithRole(userID uuid.UUID) (*userdto.UserRoleResponse, error)
	GetAll(page, limit int) (*userdto.UsersResponse, error)
	CreateUser(req *userdto.CreateUserRequest) (*userdto.UserResponse, error)
	UpdateUser(userID uuid.UUID, req *userdto.UpdateUserRequest) (*userdto.UserRoleResponse, error)
	DeleteUser(userID uuid.UUID) error
	ValidatePassword(email, password string) (*User, error)
	AssignRole(userID uuid.UUID, roleID uuid.UUID) (*userdto.UserRoleResponse, error)
	HasPermission(userID uuid.UUID, permission string) (bool, error)
	HasRole(userID uuid.UUID, roleSlug string) (bool, error)
}

// userService implements UserService interface
type userService struct {
	repo      UserRepository
	roleRepo  role.RoleRepository
}

// NewUserService creates a new user service
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// NewUserServiceWithRole creates a new user service with role repository
func NewUserServiceWithRole(repo UserRepository, roleRepo role.RoleRepository) UserService {
	return &userService{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

// GetProfile gets a user profile by ID
func (s *userService) GetProfile(userID uuid.UUID) (*userdto.UserResponse, error) {
	userModel, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := userModel.ToResponse()
	return &response, nil
}

// GetProfileWithRole gets a user profile with role information
func (s *userService) GetProfileWithRole(userID uuid.UUID) (*userdto.UserRoleResponse, error) {
	userModel, err := s.repo.FindByIDWithRole(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := userModel.ToResponseWithRole()
	return &response, nil
}

// GetAll gets all users with pagination
func (s *userService) GetAll(page, limit int) (*userdto.UsersResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Find users
	users, total, err := s.repo.FindAll(offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	userResponses := make([]userdto.UserResponse, len(users))
	for i, userModel := range users {
		userResponses[i] = userModel.ToResponse()
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &userdto.UsersResponse{
		Users: userResponses,
		Meta: userdto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

// CreateUser creates a new user with specified role (defaults to "user" role if not provided)
// Only allows creating "user" or "admin" roles, not "super_admin"
func (s *userService) CreateUser(req *userdto.CreateUserRequest) (*userdto.UserResponse, error) {
	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Determine role ID to assign
	var roleID uuid.UUID
	if req.RoleID != nil {
		// RoleID provided in request - validate it's user or admin role only
		role, err := s.roleRepo.FindByID(*req.RoleID)
		if err != nil {
			return nil, errors.New("role not found")
		}

		// Only allow "user" or "admin" roles to be assigned during creation
		if role.Slug != "user" && role.Slug != "admin" {
			return nil, errors.New("can only assign 'user' or 'admin' role during user creation")
		}

		roleID = role.ID
	} else {
		// No role specified - default to "user" role
		userRole, err := s.roleRepo.FindBySlug("user")
		if err != nil || userRole == nil {
			return nil, errors.New("default user role not found")
		}
		roleID = userRole.ID
	}

	// Create user model
	userModel := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed in BeforeCreate hook
		RoleID:   roleID, // Assign specified or default role
	}

	// Save user
	if err := s.repo.Create(userModel); err != nil {
		return nil, err
	}

	response := userModel.ToResponse()
	return &response, nil
}

// UpdateUser updates a user
// Only allows updating role to "user" or "admin", not "super_admin"
func (s *userService) UpdateUser(userID uuid.UUID, req *userdto.UpdateUserRequest) (*userdto.UserRoleResponse, error) {
	// Find user
	userModel, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being changed and if it already exists
	if req.Email != "" && req.Email != userModel.Email {
		exists, err := s.repo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
		userModel.Email = req.Email
	}

	// Update name if provided
	if req.Name != "" {
		userModel.Name = req.Name
	}

	// Update role if provided
	if req.RoleID != nil {
		// Verify role exists
		role, err := s.roleRepo.FindByID(*req.RoleID)
		if err != nil {
			return nil, errors.New("role not found")
		}

		// Only allow "user" or "admin" roles to be assigned during update
		// "super_admin" role can only be assigned via AssignRole endpoint (SuperAdmin only)
		if role.Slug != "user" && role.Slug != "admin" {
			return nil, errors.New("can only assign 'user' or 'admin' role during user update")
		}

		userModel.RoleID = role.ID
	}

	// Save changes
	if err := s.repo.Update(userModel); err != nil {
		return nil, err
	}

	// Load user with role to return complete response
	userWithRole, err := s.repo.FindByIDWithRole(userID)
	if err != nil {
		return nil, err
	}

	response := userWithRole.ToResponseWithRole()
	return &response, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(userID uuid.UUID) error {
	// Check if user exists
	_, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Delete user
	if err := s.repo.Delete(userID); err != nil {
		return err
	}

	return nil
}

// ValidatePassword validates user credentials
func (s *userService) ValidatePassword(email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password
	if !utils.ComparePassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// AssignRole assigns a role to a user
func (s *userService) AssignRole(userID uuid.UUID, roleID uuid.UUID) (*userdto.UserRoleResponse, error) {
	// Find user
	userModel, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify role exists
	if s.roleRepo != nil {
		_, err := s.roleRepo.FindByID(roleID)
		if err != nil {
			return nil, errors.New("role not found")
		}
	}

	// Assign role
	userModel.RoleID = roleID

	// Save changes
	if err := s.repo.Update(userModel); err != nil {
		return nil, err
	}

	// Load user with role
	userWithRole, err := s.repo.FindByIDWithRole(userID)
	if err != nil {
		return nil, err
	}

	response := userWithRole.ToResponseWithRole()
	return &response, nil
}

// HasPermission checks if a user has a specific permission
func (s *userService) HasPermission(userID uuid.UUID, permission string) (bool, error) {
	user, err := s.repo.FindByIDWithRole(userID)
	if err != nil {
		return false, err
	}

	if user.Role == nil {
		return false, nil
	}

	// Check for wildcard permission
	for _, p := range user.Role.Permissions {
		if p == "*" {
			return true, nil
		}
	}

	// Check specific permission
	for _, p := range user.Role.Permissions {
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}

// HasRole checks if a user has a specific role (by slug)
func (s *userService) HasRole(userID uuid.UUID, roleSlug string) (bool, error) {
	user, err := s.repo.FindByIDWithRole(userID)
	if err != nil {
		return false, err
	}

	if user.Role == nil {
		return false, nil
	}

	return user.Role.Slug == roleSlug, nil
}
