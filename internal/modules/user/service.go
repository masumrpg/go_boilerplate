package user

import (
	"errors"
	"math"

	"go_boilerplate/internal/shared/utils"
	userdto "go_boilerplate/internal/modules/user/dto"

	"github.com/google/uuid"
)

// UserService defines the interface for user business logic
type UserService interface {
	GetProfile(userID uuid.UUID) (*userdto.UserResponse, error)
	GetAll(page, limit int) (*userdto.UsersResponse, error)
	CreateUser(req *userdto.CreateUserRequest) (*userdto.UserResponse, error)
	UpdateUser(userID uuid.UUID, req *userdto.UpdateUserRequest) (*userdto.UserResponse, error)
	DeleteUser(userID uuid.UUID) error
	ValidatePassword(email, password string) (*User, error)
}

// userService implements UserService interface
type userService struct {
	repo UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
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

// CreateUser creates a new user
func (s *userService) CreateUser(req *userdto.CreateUserRequest) (*userdto.UserResponse, error) {
	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create user model
	userModel := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed in BeforeCreate hook
	}

	// Save user
	if err := s.repo.Create(userModel); err != nil {
		return nil, err
	}

	response := userModel.ToResponse()
	return &response, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(userID uuid.UUID, req *userdto.UpdateUserRequest) (*userdto.UserResponse, error) {
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

	// Save changes
	if err := s.repo.Update(userModel); err != nil {
		return nil, err
	}

	response := userModel.ToResponse()
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
