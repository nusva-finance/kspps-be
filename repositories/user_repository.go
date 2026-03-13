package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: config.GetDB(),
	}
}

// GetDB returns the database instance for transactions
func GetDB() *gorm.DB {
	return config.GetDB()
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(offset, limit int, search string) ([]models.User, int64, error) {
    var users []models.User
    var total int64

    // 1. Tambahkan .Debug() di sini
    countQuery := r.db.Debug().Model(&models.User{})
    if search != "" {
        countQuery = countQuery.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
            "%"+search+"%", "%"+search+"%", "%"+search+"%")
    }

    err := countQuery.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    // 2. Tambahkan .Debug() di sini juga
    dataQuery := r.db.Debug().Model(&models.User{})
    if search != "" {
        dataQuery = dataQuery.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
            "%"+search+"%", "%"+search+"%", "%"+search+"%")
    }

    var userIDs []uint
    err = dataQuery.Offset(offset).Limit(limit).Pluck("id", &userIDs).Error
    if err != nil {
        return nil, 0, err
    }

    if len(userIDs) > 0 {
        err = r.db.Preload("Roles").
            Where("id IN ?", userIDs).
            Order("id").
            Find(&users).Error
    }

    return users, total, nil
}

func (r *UserRepository) AssignRole(userID, roleID uint) error {
	return r.db.Create(&models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *UserRepository) RemoveRole(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
}
