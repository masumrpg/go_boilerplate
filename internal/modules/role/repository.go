package role

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	Create(role *Role) error
	FindByID(id uuid.UUID) (*Role, error)
	FindBySlug(slug string) (*Role, error)
	FindAll(offset, limit int) ([]Role, int64, error)
	Update(role *Role) error
	Delete(id uuid.UUID) error
	ExistsBySlug(slug string) (bool, error)
	ExistsByName(name string) (bool, error)
}

// roleRepository implements RoleRepository interface
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create creates a new role
func (r *roleRepository) Create(role *Role) error {
	return r.db.Create(role).Error
}

// FindByID finds a role by ID
func (r *roleRepository) FindByID(id uuid.UUID) (*Role, error) {
	var role Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindBySlug finds a role by slug
func (r *roleRepository) FindBySlug(slug string) (*Role, error) {
	var role Role
	err := r.db.Where("slug = ?", slug).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &role, nil
}

// FindAll finds all roles with pagination
func (r *roleRepository) FindAll(offset, limit int) ([]Role, int64, error) {
	var roles []Role
	var total int64

	// Count total
	if err := r.db.Model(&Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Find roles with pagination
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Update updates a role
func (r *roleRepository) Update(role *Role) error {
	return r.db.Save(role).Error
}

// Delete deletes a role by ID
func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Role{}, "id = ?", id).Error
}

// ExistsBySlug checks if a role with the given slug exists
func (r *roleRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&Role{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

// ExistsByName checks if a role with the given name exists
func (r *roleRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
