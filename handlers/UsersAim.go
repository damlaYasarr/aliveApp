package handlers

import (
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
    "gorm.io/gorm"
    "errors"
    "strings"

	
)


func ListUsersAim(c *fiber.Ctx) error {
    email := c.Query("email")
    if email == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
    }

    userID, err := GetUserIDByEmail(email)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
    }

    type Aim struct {
        Name             string   `json:"name"`
        COMPLETE_DAYS    string   `json:"complete_days"`
        Startday         string   `json:"startday"`
        Endday           string   `json:"endday"`
        NotificationHour string   `json:"notificationhour"`
        CompleteDaysCount int      `json:"complete_days_count"`
    }
    
    // Fetch aims and time information associated with a user ID
    var aims []Aim
    if err := database.DB.Db.Raw(`
        SELECT a.name, t.complete_days, a.startday, a.endday, a.notification_hour
        FROM aims a
        JOIN times t ON a.id = t.aim_id
        WHERE a.user_id = ? `, userID).Scan(&aims).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
    }
    
    // Calculate the count of complete days for each aim
    for i := range aims {
        // Extract complete days from the string
        completeDaysStr := aims[i].COMPLETE_DAYS
        // Remove curly braces and split the string by commas
        dates := strings.Split(strings.Trim(completeDaysStr, "{}"), ",")
        // Count the number of dates
        aims[i].CompleteDaysCount = len(dates)
    }
    
    return c.Status(fiber.StatusOK).JSON(aims)
}
// donotduplicate aim name --USE IT 
func Donotduplicatename(c *fiber.Ctx) error {
    type Request struct {
        AimName string `json:"aimname"`
        UserID  uint   `json:"userid"`
    }
    
    var req Request
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
    }

    // Aim adını ve kullanıcı ID'sini kontrol et
    _, err := GetAIMIDByNAME(req.AimName, req.UserID)
    if err == nil {
        // Aim mevcutsa, hata döndür
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "aim name already exists for this user"})
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        // Başka bir hata oluşmuşsa, hatayı döndür
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
    }

    // Aim mevcut değilse, işlemi devam ettir (örneğin, yeni bir aim oluşturabilirsiniz)
    // Yeni aim oluşturma işlemi burada yapılabilir

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "aim name is available"})
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

    //create new aim id in the time table
    aimid, err :=GetAIMIDByNAME(requestBody.Aim,userID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to get aim ID")
    }
// Create a new Aim record and assign the user ID
    newaimid := models.Time{
        AIM_ID: int64(aimid),
    }
    if err := database.DB.Db.Create(&newaimid).Error; err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to create aim_id in time table record")
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

// GetUserIDByEmail, e-posta adresine göre kullanıcı kimliğini alır
func GetAIMIDByNAME(aimname string, userid uint) (uint, error) {
    var aim models.Aim
    if err := database.DB.Db.Where("name = ? AND user_id = ?", aimname, userid).First(&aim).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return 0, errors.New("aim not found")
        }
        return 0, err
    }
    return uint(aim.ID), nil
}

func GETActivedays(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Kullanıcı ID'si ile ilişkili tüm aim ve zaman bilgilerini al
	var aims []struct {
		ID            int64    `json:"id"`
		COMPLETE_DAYS string `json:"complete_days"`
	}
	if err := database.DB.Db.Raw(`
        SELECT a.id, t.complete_days 
        FROM aims a
        JOIN times t ON a.id = t.aim_id
        WHERE a.user_id = ? `, userID).Scan(&aims).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

  

	return c.Status(fiber.StatusOK).JSON(aims)
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


//get the all resuult
//getuser as a premium  payment microvervice


//conversation with premium part
func PremiumUser(c *fiber.Ctx) error {
	facts := []models.User{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}