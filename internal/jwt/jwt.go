package jwt

import (
	"fmt"

	"github.com/charlieroth/reminders/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(conf *config.Config, userID uuid.UUID, expUtc int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expUtc,
	})
	tokenString, err := token.SignedString([]byte(conf.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(conf *config.Config, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func RefreshToken(conf *config.Config, tokenString string) (string, error) {
	return "", nil
}

func UserID(conf *config.Config, tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JwtSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.MustParse(token.Claims.(jwt.MapClaims)["sub"].(string)), nil
}

func ExpiresAt(conf *config.Config, tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	return int64(token.Claims.(jwt.MapClaims)["exp"].(int64)), nil
}
