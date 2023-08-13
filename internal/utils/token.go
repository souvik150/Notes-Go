package usercontrollers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/souvik150/golang-fiber/config"
	"github.com/souvik150/golang-fiber/internal/models"
	"time"

	"github.com/google/uuid"
)

func generateToken(userID uuid.UUID, secretKey string, expiration time.Duration, user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID.String(),
		"exp":       time.Now().Add(expiration).Unix(),
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt.Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func GenerateAccessToken(user *models.User) (string, error) {
	config, _ := initializers.LoadConfig(".")
	return generateToken(user.ID, config.AccessTokenSecret, config.AccessTokenExpiry, user)
}

func GenerateRefreshToken(user *models.User) (string, error) {
	config, _ := initializers.LoadConfig(".")
	return generateToken(user.ID, config.RefreshTokenSecret, config.RefreshTokenExpiry, user)
}
