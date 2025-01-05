package auth

import (
	"errors"
	"fmt"

	"github.com/bigxxby/dream-test-task/internal/api/repo/auth"
	"github.com/bigxxby/dream-test-task/internal/api/repo/user"
	"github.com/bigxxby/dream-test-task/internal/models"
	"github.com/bigxxby/dream-test-task/internal/utils"
	"github.com/google/uuid"
)

type IAuthService interface {
	Login(username, password string) (string, int, error)
	Register(username, password string) (*models.User, int, error)
	WHOAMI(userId *uuid.UUID) (*models.User, int, error)
}

func (as AuthService) WHOAMI(userId *uuid.UUID) (*models.User, int, error) {
	user, err := as.UserRepo.GetUserById(userId)
	if err != nil {
		return nil, 404, err
	}
	return user, 200, nil
}
func (as AuthService) Login(username, password string) (string, int, error) {
	// Проверяем наличие пользователя в базе данных
	user, _ := as.UserRepo.GetUserByName(username)
	if user == nil {
		return "", 404, errors.New("user not found")
	}

	// Проверяем пароль
	if !user.ComparePassword(password) {
		return "", 401, errors.New("invalid password")
	}
	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return "", 500, err
	}

	return token, 200, nil
}

type AuthService struct {
	//repo
	AuthRepo auth.IAuthRepo
	UserRepo user.IUserRepo
}

func NewAuthService(authRepo auth.IAuthRepo, userRepo user.IUserRepo) IAuthService {
	return AuthService{
		AuthRepo: authRepo,
		UserRepo: userRepo,
	}
}

func (as AuthService) Register(username, password string) (*models.User, int, error) {
	newUser := models.User{
		Username: username,
		Password: password,
	}
	err := newUser.ValidatePassword()
	if err != nil {
		return nil, 400, err
	}

	user, _ := as.UserRepo.GetUserByName(username)
	if user != nil {
		return nil, 409, errors.New("user already exists")
	}
	fmt.Println("user")

	newUser.HashPassword()

	// Сохраняем пользователя в базе данных
	createdUser, err := as.UserRepo.CreateUser(newUser)
	if err != nil {
		return nil, 500, err
	}

	return createdUser, 200, nil
}
