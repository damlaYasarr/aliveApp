package handlers

import (
	"errors"
	"net/http"

	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUser(c *fiber.Ctx) error {
	var requestBody struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}

	// Check if the user with the provided email already exists
	var existingUser models.User
	if err := database.DB.Db.Where("email = ?", requestBody.Email).First(&existingUser).Error; err == nil {
		// User with the provided email already exists, return the existing user as a response
		return c.Status(http.StatusConflict).JSON(existingUser)
	}

	// Create a new user
	newUser := models.User{
		Email: requestBody.Email,
	}
	if err := database.DB.Db.Create(&newUser).Error; err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to register user")
	}

	// Get the ID of the newly created user
	userID, err := GetUserIDByEmail(requestBody.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get user ID")
	}

	// Create a new Aim record and assign the user ID
	newAim := models.Aim{
		USERID: int64(userID),
	}
	if err := database.DB.Db.Create(&newAim).Error; err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to create Aim record")
	}

	// Return the newly registered user as a JSON response
	return c.Status(http.StatusOK).JSON(newUser)
}

// GetUserIDByEmail, e-posta adresine göre kullanıcı kimliğini alır
func GetUserIDByEmail(email string) (uint, error) {
	var user models.User
	if err := database.DB.Db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("User not found")
		}
		return 0, err
	}
	return uint(user.ID), nil
}
