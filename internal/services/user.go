package services

import (
	"github.com/google/uuid"
	database "github.com/souvik150/golang-fiber/internal/database"
	"github.com/souvik150/golang-fiber/internal/models"
)

func GetUserByID(userID uuid.UUID) (models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}
