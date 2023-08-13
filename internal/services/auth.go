package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	initializers "github.com/souvik150/golang-fiber/config"
	database "github.com/souvik150/golang-fiber/internal/database"
	"github.com/souvik150/golang-fiber/internal/models"
	token "github.com/souvik150/golang-fiber/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"time"
)

func SignupUser(payload *models.RegisterUserSchema) (models.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.AuthResponse{}, err
	}

	newUser := models.User{
		Username:     payload.Username,
		Email:        payload.Email,
		Password:     string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ProfileImage: payload.ProfileImage,
	}

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return models.AuthResponse{}, result.Error
	}

	authResponse, err := generateAuthResponse(&newUser)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return authResponse, nil
}

func UploadProfilePic(fileReader io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return "", err
	}

	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(config.AWSRegion),
			Credentials: credentials.NewStaticCredentials(config.AWSAccessKey, config.AWSSecretKey, ""),
		},
	})

	if err != nil {
		panic(err)
	}

	uploader := s3manager.NewUploader(awsSession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   fileReader,
	})
	if err != nil {
		return "", err
	}

	// Get the URL of the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.AWSBucketName, fileHeader.Filename)

	return url, nil
}

func LoginUser(payload *models.LoginUserSchema) (models.AuthResponse, error) {
	var user models.User
	result := database.DB.Where("username = ?", payload.Username).First(&user)
	if result.Error != nil {
		return models.AuthResponse{}, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return models.AuthResponse{}, err
	}

	authResponse, err := generateAuthResponse(&user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return authResponse, nil
}

func generateAuthResponse(user *models.User) (models.AuthResponse, error) {
	accessToken, err := token.GenerateAccessToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	refreshToken, err := token.GenerateRefreshToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Store the refresh token in the database
	refreshTokenEntry := models.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
	}
	refreshTokenResult := database.DB.Create(&refreshTokenEntry)
	if refreshTokenResult.Error != nil {
		return models.AuthResponse{}, refreshTokenResult.Error
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

func RefreshAccessToken(refreshToken string) (models.AuthResponse, error) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Check if the provided refresh token exists in the database
	var refreshTokenEntry models.RefreshToken
	result := database.DB.Where("token = ?", refreshToken).First(&refreshTokenEntry)
	if result.Error != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Parse and validate the access token
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.RefreshTokenSecret), nil
	})
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Generate a new access token
	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	var user models.User
	user.ID = userID
	accessToken, err := token.GenerateAccessToken(&user)
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}
