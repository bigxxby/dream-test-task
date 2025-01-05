package utils

import (
	"time"

	"github.com/bigxxby/dream-test-task/internal/config"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа из конфигурации
	return token.SignedString((config.JwtSecret))
}
