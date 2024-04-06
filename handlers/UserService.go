package handlers

import (
	"net/http"
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
)

//gooogle auth will added in the end of the flutter design

func RegisterAppWithEmail(c *fiber.Ctx) error {
    // Parse request body to get user's email
    var requestBody struct {
        Email string `json:"email"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid request body")
    }

    // Check if email is empty
    if requestBody.Email == "" {
        return c.Status(http.StatusBadRequest).SendString("Email is required")
    }


    newUser := models.User{
        Email: requestBody.Email,
    }

    if err := database.DB.Db.Create(&newUser).Error; err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to register user")
    }

    return c.Status(http.StatusOK).JSON(newUser)
}
