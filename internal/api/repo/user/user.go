package user

import (
	"github.com/bigxxby/dream-test-task/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepo interface {
	CreateUser(user models.User) (*models.User, error)
	GetUserByName(username string) (*models.User, error)
	GetUserById(userId *uuid.UUID) (*models.User, error)
	DeleteUser(username string) error
	UpdateUser(user models.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (ur UserRepo) GetUserById(userId *uuid.UUID) (*models.User, error) {
	var user models.User
	result := ur.db.Where("id = ?", userId.String()).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (ur UserRepo) CreateUser(user models.User) (*models.User, error) {
	result := ur.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur UserRepo) GetUserByName(username string) (*models.User, error) {
	var user models.User
	result := ur.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur UserRepo) DeleteUser(username string) error {

	result := ur.db.Where("username = ?", username).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur UserRepo) UpdateUser(user models.User) error {
	result := ur.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
