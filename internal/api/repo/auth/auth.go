package auth

import (
	"gorm.io/gorm"
)

type IAuthRepo interface {
}

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}
