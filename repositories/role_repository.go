package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: config.GetDB(),
	}
}

func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("is_active = ?", true).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}
