package handlers

import (
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
	
)


func ListUsersAim(c *fiber.Ctx) error {
	facts := []models.User{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}

//Add new user
func AddnewAim(c *fiber.Ctx) error {

    var requestBody struct {
        UserID    int64      `json:"user_id"`
        Aim       string    `json:"aim"`
        AimDate   string `json:"aim_date"`
		Endday    string `json:"endday"`
		Notification string`json:"notif"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid request body")
    }

    // Check if the user exists
    var user models.User
    if err := database.DB.Db.First(&user, requestBody.UserID).Error; err != nil {
        return c.Status(http.StatusNotFound).SendString("User not found")
    }

    // Create a new aim object
    newAim := models.Aim{
        USERID:   requestBody.UserID,
        Name:      requestBody.Aim,
        Startday:  requestBody.AimDate,
        Endday: requestBody.Endday,
        NotificationHour: requestBody.Notification,
    }

    // Insert the new aim into the database
    if err := database.DB.Db.Create(&newAim).Error; err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to add new aim")
    }
    return c.Status(http.StatusOK).JSON(newAim)
}
// ListUsersAllAim lists all aims for a specific user
func ListUsersAllAim(c *fiber.Ctx) error {
    // Parse request body to get user ID
    var requestBody struct {
        UserID int64 `json:"user_id"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid request body")
    }

    // Check if the user exists
    var user models.User
    if err := database.DB.Db.First(&user, requestBody.UserID).Error; err != nil {
        return c.Status(http.StatusNotFound).SendString("User not found")
    }

    // Retrieve all aims for the user
    var aims []models.Aim
    if err := database.DB.Db.Where("user_id = ?", requestBody.UserID).Find(&aims).Error; err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to retrieve aims")
    }

    // Return the list of aims
    return c.Status(http.StatusOK).JSON(aims)
}

//delete aim
//arrange calendar
//arrange user as a premium
//get the all resuult

//conversation with premium part
func PremiumUser(c *fiber.Ctx) error {
	facts := []models.User{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}