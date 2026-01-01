package user

import (
	"time"

	roleModule "go_boilerplate/internal/modules/role"
	"go_boilerplate/internal/modules/user/dto"
	"go_boilerplate/internal/shared/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string                 `json:"name" gorm:"type:varchar(100);not null"`
	Email     string                 `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string                 `json:"-" gorm:"type:varchar(255);not null"` // Never expose password in JSON
	RoleID    uuid.UUID              `json:"role_id" gorm:"type:uuid;not null"`   // Foreign key to m_roles
	Role      *roleModule.Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"` // Role relationship (eager load)
	IsVerified bool                  `json:"is_verified" gorm:"default:false"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	DeletedAt gorm.DeletedAt         `json:"-" gorm:"index"` // Soft delete support
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "m_users"
}

// BeforeCreate hook runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return u.hashPassword()
}

// BeforeUpdate hook runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.hashPassword()
}

// hashPassword handles password hashing if needed
func (u *User) hashPassword() error {
	if u.Password != "" && !utils.IsHashed(u.Password) {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

// ToResponse converts User to UserResponse (without password)
func (u *User) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// ToResponseWithRole converts User to UserResponse with role information
func (u *User) ToResponseWithRole() dto.UserRoleResponse {
	response := dto.UserRoleResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}

	if u.Role != nil {
		response.Role = &dto.RoleInfo{
			ID:          u.Role.ID,
			Name:        u.Role.Name,
			Slug:        u.Role.Slug,
			Permissions: []string(u.Role.Permissions),
		}
	}

	return response
}
