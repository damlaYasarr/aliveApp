package handlers

import (
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
    "gorm.io/gorm"
    "errors"
	
)


func ListUsersAim(c *fiber.Ctx) error {
	facts := []models.User{}
	database.DB.Db.Find(&facts)
 //ana ekranda aynı userin bilgileri time bilgisi ile yayınlanacak
	return c.Status(200).JSON(facts)
}

func AddNewAim(c *fiber.Ctx) error {
    // Parse request body
    var requestBody struct {
        Email          string `json:"email"`
        Aim            string `json:"aim"`
        AimDate        string `json:"aim_date"`
        Endday         string `json:"endday"`
        Notification   string `json:"notif"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid request body")
    }

    // Get the user ID from the email
    userID, err := GetUserIDByEmail(requestBody.Email)
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to get user ID")
    }

    // Check if the user exists
    var user models.User
    if err := database.DB.Db.First(&user, userID).Error; err != nil {
        return c.Status(http.StatusNotFound).SendString("User not found")
    }

    // Create a new aim object
    newAim := models.Aim{
        USERID:           int64(userID),
        Name:              requestBody.Aim,
        Startday:          requestBody.AimDate,
        Endday:            requestBody.Endday,
        NotificationHour:  requestBody.Notification,
    }

    // Insert the new aim into the database
    if err := database.DB.Db.Create(&newAim).Error; err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to add new aim")
    }

    return c.Status(http.StatusOK).JSON(newAim)
}
// GetUserIDByEmail, e-posta adresine göre kullanıcı kimliğini alır
func GetUserIDByEmails(email string) (uint, error) {
    var user models.User
    if err := database.DB.Db.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return 0, errors.New("User not found")
        }
        return 0, err
    }
    return uint(user.ID), nil
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

// get habit by name
func GetHabitByName(db *gorm.DB, name string) (*models.Aim, error) {
    var habit models.Aim
    if err := db.Where("name = ?", name).First(&habit).Error; err != nil {
        return nil, err
    }

    return &habit, nil
}
//delete habit its name
func DeleteHabitByName(c *fiber.Ctx) error {
    // Parse the request body
    var requestBody struct {
        Name string `json:"name"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
    }

    // Retrieve the habit name from the request body
    name := requestBody.Name

    // Retrieve the database instance from Fiber context
    db := c.Locals("db").(*gorm.DB)

    // Find the habit by its name
    var habit models.Aim
    if err := db.Where("name = ?", name).First(&habit).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Habit not found")
    }

    // Delete the habit
    if err := db.Delete(&habit).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete habit")
    }

    return c.SendString("Habit deleted successfully")
}
// edit habit by name_ use gethabit name func here 
func EditHabitName(db *gorm.DB, currentName string, newName string) error {
    // Find the habit by its current name
    var habit models.Aim
    if err := db.Where("name = ?", currentName).First(&habit).Error; err != nil {
        return err // Return error if habit with given current name is not found
    }

    // Update the habit's name
    habit.Name = newName

    // Save the changes to the database
    if err := db.Save(&habit).Error; err != nil {
        return err // Return error if saving fails
    }

    return nil // Return nil if editing is successful
}



//arrange calendar
//arrange user as a premium
//get the all resuult



//conversation with premium part
func PremiumUser(c *fiber.Ctx) error {
	facts := []models.User{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}