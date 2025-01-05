package models

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       *uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Username string     `json:"username" gorm:"unique;not null"`
	Password string     `json:"-" gorm:"not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	new := uuid.New()
	u.ID = &new
	return
}

// hash password with bcrypt
func (u *User) HashPassword() error {
	oldPassword := u.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) ValidatePassword() error {
	if u.Password == "" {
		return errors.New("password is required")
	}

	// Validate password length
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Validate password complexity
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	for _, char := range u.Password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '@' || char == '$' || char == '!' || char == '%' || char == '*' || char == '?' || char == '&':
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain one uppercase letter, one lowercase letter, one digit, and one special character")
	}

	return nil // Пароль соответствует требованиям
}
