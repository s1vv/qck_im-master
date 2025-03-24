package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"qckim-backend/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var cfg = config.GetConfig()

// Генерация JWT-токена
func GenerateJWT(userID int64) (string, error) {
	claims := Claims{
		UserID: strconv.FormatInt(userID, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// Валидация JWT-токена
func ValidateJWT(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("пустой токен func ValidateJWT")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("неверный токен")
	}

	// Проверяем срок действия
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("срок действия токена истёк")
	}

	return claims, nil
}

// GenerateToken создает случайную строку длиной 32 символа
func GenerateToken32() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
