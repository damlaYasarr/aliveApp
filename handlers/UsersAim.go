package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// I dont use it
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
		Name              string `json:"name"`
		COMPLETE_DAYS     string `json:"complete_days"`
		Startday          string `json:"startday"`
		Endday            string `json:"endday"`
		NotificationHour  string `json:"notificationhour"`
		CompleteDaysCount int    `json:"complete_days_count"`
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

		// Check if completeDaysStr is null or empty
		if completeDaysStr == "" || completeDaysStr == "null" {
			aims[i].CompleteDaysCount = 0
		} else {
			// Remove curly braces and split the string by commas
			dates := strings.Split(strings.Trim(completeDaysStr, "{}"), ",")
			// Count the number of dates
			aims[i].CompleteDaysCount = len(dates)
		}

	}
	// Filter active aims
	var activeAims []Aim
	for _, aim := range aims {
		if IsAimActive(aim.Startday, aim.Endday) {
			activeAims = append(activeAims, aim)
		}
	}
	return c.Status(fiber.StatusOK).JSON(activeAims)
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

	_, err := GetAIMIDByNAME(req.AimName, req.UserID)
	if err == nil {

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "aim name already exists for this user"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "aim name is available"})
}

// AddNewAim handles adding a new aim and scheduling a notification
func AddNewAim(c *fiber.Ctx) error {
	// Parse request body
	var requestBody struct {
		Email        string `json:"email"`
		Aim          string `json:"name"`
		AimDate      string `json:"startday"`
		Endday       string `json:"endday"`
		Notification string `json:"notification_hour"`
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

	newAim := models.Aim{
		USERID:           int64(userID),
		Name:             requestBody.Aim,
		Startday:         requestBody.AimDate,
		Endday:           requestBody.Endday,
		NotificationHour: requestBody.Notification,
	}

	// Insert the new aim into the database
	if err := database.DB.Db.Create(&newAim).Error; err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to add new aim")
	}

	// Create a new aim ID in the time table
	aimid, err := GetAIMIDByNAME(requestBody.Aim, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get aim ID")
	}
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
		ID            int64  `json:"id"`
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

// ListUsersActiveAim, belirtilen e-posta adresine sahip kullanıcının aktif hedeflerini listeler
func ListUsersActiveAim(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	type Aim struct {
		ID               int64  `json:"id"`
		Name             string `json:"name"`
		Startday         string `json:"startday"`
		Endday           string `json:"endday"`
		NotificationHour string `json:"notification_hour"`
	}

	var aims []Aim
	if err := database.DB.Db.Raw(`
        SELECT a.id, a.name, a.startday, a.endday, a.notification_hour
        FROM aims a
        WHERE a.user_id = ?`, userID).Scan(&aims).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch user aims"})
	}

	// Filter active aims
	var activeAims []Aim
	for _, aim := range aims {
		if IsAimActive(aim.Startday, aim.Endday) {
			activeAims = append(activeAims, aim)
		}
	}

	return c.Status(fiber.StatusOK).JSON(activeAims)
}

// IsAimActiveAtCurrentTime checks if the current time is within the specified aim's active period and matches the notification hour.
func IsAimActiveAtCurrentTime(startday, endday, notificationHour string) bool {
	// Set the local time zone to Europe/Istanbul (UTC+3)
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		fmt.Println("Error loading location:", err)
		// If location loading fails, default to UTC
		loc = time.UTC
	}

	now := time.Now().In(loc) // Get the current time in Istanbul timezone

	// Parse the start day
	startTime, err := time.ParseInLocation("02.01.2006", startday, loc)
	if err != nil {
		fmt.Println("Error parsing start day:", err)
		return false
	}

	// Parse the end day
	endTime, err := time.ParseInLocation("02.01.2006", endday, loc)
	if err != nil {
		fmt.Println("Error parsing end day:", err)
		return false
	}

	// Check if the current date is within the specified range
	if now.Before(startTime) || now.After(endTime) {
		return false
	}

	// Parse the notification hour
	notificationTime, err := time.ParseInLocation("3:04 PM", notificationHour, loc)
	if err != nil {
		fmt.Println("Error parsing notification hour:", err)
		return false
	}

	// Debugging: Print parsed and current times
	fmt.Printf("Current time: %02d:%02d %s\n", now.Hour(), now.Minute(), now.Format("PM"))
	fmt.Printf("Notification time: %02d:%02d %s\n", notificationTime.Hour(), notificationTime.Minute(), notificationTime.Format("PM"))

	// Check if the current time matches the notification hour and minute, considering AM/PM
	return now.Hour() == notificationTime.Hour() && now.Minute() == notificationTime.Minute() && now.Format("PM") == notificationTime.Format("PM")
}

// Check if an aim is active
func IsAimActive(startday string, endday string) bool {
	now := time.Now()

	startTime, err := time.Parse("02.01.2006", startday)
	if err != nil {
		return false
	}
	endTime, err := time.Parse("02.01.2006", endday)
	if err != nil {
		return false
	}

	return now.After(startTime) && now.Before(endTime)
}

type Aim struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Startday         string `json:"startday"`
	Endday           string `json:"endday"`
	NotificationHour string `json:"notification_hour"`
}

func ListActiveHabits(email string) ([]Aim, error) {
	// Get user ID by email
	userID, err := GetUserIDByEmail(email)
	if err != nil {
		// Handle user not found error
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	// Query active aims from database
	var aims []Aim
	if err := database.DB.Db.Raw(`
        SELECT id, name, startday, endday, notification_hour
        FROM aims
        WHERE user_id = ?`, userID).Scan(&aims).Error; err != nil {
		// Handle database query error
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch user aims")
	}

	// Filter active aims
	var activeAims []Aim
	for _, aim := range aims {
		if IsAimActive(aim.Startday, aim.Endday) {
			activeAims = append(activeAims, aim)
		}
	}

	// Return active aims
	return activeAims, nil
}

// it is not used
func ListActiveHabitsTrial(users []models.User) (map[string][]Aim, error) {
	var wg sync.WaitGroup
	activeAimsMap := make(map[string][]Aim)

	for _, user := range users {
		wg.Add(1) // Her kullanıcı için bir iş parçacığı ekle

		go func(user models.User) {
			defer wg.Done() // İş parçacığı tamamlandığında sayacı azalt

			// Kullanıcının ID'sini al
			userID, err := GetUserIDByEmail(user.Email)
			if err != nil {
				log.Printf("User not found: %s, error: %v", user.Email, err)
				return
			}

			// Veritabanından aktif hedefleri sorgula
			var aims []Aim
			if err := database.DB.Db.Raw(`
                SELECT id, name, startday, endday, notification_hour
                FROM aims
                WHERE user_id = ?`, userID).Scan(&aims).Error; err != nil {
				log.Printf("Failed to fetch user aims for %s: %v", user.Email, err)
				return
			}

			// Aktif hedefleri filtrele
			var activeAims []Aim
			for _, aim := range aims {
				if IsAimActive(aim.Startday, aim.Endday) {
					activeAims = append(activeAims, aim)
				}
			}

			// Aktif hedefleri haritaya ekle
			activeAimsMap[user.Email] = activeAims
		}(user)
	}

	wg.Wait()                 // Tüm iş parçacıklarının tamamlanmasını bekle
	return activeAimsMap, nil // Kullanıcı e-postasına göre aktif hedefleri döndür
}

// ApprovalHabitDate handles the request to update the habit date
func ApprovalHabitDate(c *fiber.Ctx) error {

	var requestBody struct {
		Email   string `json:"email"`
		Aim     string `json:"name"`
		AimDate string `json:"complete_days"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	aimID, err := GetHabitIdByName(requestBody.Aim, requestBody.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve aim_id: " + err.Error())
	}

	var existingTime models.Time
	if err := database.DB.Db.Where("aim_id = ?", aimID).First(&existingTime).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			newTime := models.Time{
				AIM_ID:        aimID,
				COMPLETE_DAYS: []string{requestBody.AimDate},
			}
			if err := database.DB.Db.Create(&newTime).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to create new time table record")
			}
			return c.Status(fiber.StatusOK).SendString("Time table record created successfully")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve time record")
	}

	if existingTime.COMPLETE_DAYS == nil {
		existingTime.COMPLETE_DAYS = make([]string, 0)
	}

	// Check if AimDate already exists in COMPLETE_DAYS array
	var dateExists bool
	for _, day := range existingTime.COMPLETE_DAYS {
		if day == requestBody.AimDate {
			dateExists = true
			break
		}
	}

	if !dateExists {
		existingTime.COMPLETE_DAYS = append(existingTime.COMPLETE_DAYS, requestBody.AimDate)

		arrayData := pq.Array(existingTime.COMPLETE_DAYS)

		query := `
            UPDATE times 
            SET complete_days = $1 
            WHERE aim_id = $2
        `
		// Execute the SQL query
		if err := database.DB.Db.Exec(query, arrayData, aimID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to update time table record")
		}
	}

	// Return success status
	return c.Status(fiber.StatusOK).SendString("Time table record updated successfully")
}

// get habit's id with email
func GetHabitIdByName(name string, email string) (int64, error) {
	// Retrieve user ID by email
	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return 0, err
	}

	// Find the habit (aim) by name and user ID
	var aim models.Aim
	if err := database.DB.Db.Where("name = ? AND user_id = ?", name, userID).First(&aim).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("Habit not found")
		}
		return 0, err
	}

	// Return the habit's ID
	fmt.Printf("Habit ID: %d\n", aim.ID)
	return aim.ID, nil
}

// deleete user aim from aim table --> It works
func DeleteUserAim(c *fiber.Ctx) error {
	// Extract the email and aim name from query parameters
	email := c.Query("email")
	aimname := c.Query("name")

	// Check if the email and aim name are provided
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}
	if aimname == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "aim name is required"})
	}

	// Retrieve aim_id using the provided Aim name and Email
	aim_id, err := GetHabitIdByName(aimname, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve aim_id: " + err.Error()})
	}

	// Get the user ID by email
	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Check if the aim exists and belongs to the user
	var aimExists bool
	if err := database.DB.Db.Raw(`
        SELECT EXISTS(
            SELECT 1
            FROM aims
            WHERE id = ? AND user_id = ?
        )`, aim_id, userID).Scan(&aimExists).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if !aimExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "aim not found or does not belong to user"})
	}

	// Delete the aim from the database
	if err := database.DB.Db.Exec(`DELETE FROM aims WHERE id = ?`, aim_id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	// Delete the times associated with the aim from the database
	if err := database.DB.Db.Exec(`DELETE FROM times WHERE aim_id = ?`, aim_id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	// Return a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "aim deleted successfully"})
}

// list expired habits
func ListExpiredHabits(c *fiber.Ctx) error {

	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	type Aim struct {
		Name              string `json:"name"`
		COMPLETE_DAYS     string `json:"complete_days"`
		Startday          string `json:"startday"`
		Endday            string `json:"endday"`
		NotificationHour  string `json:"notificationhour"`
		CompleteDaysCount int    `json:"complete_days_count"`
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

		// Check if completeDaysStr is null or empty
		if completeDaysStr == "" || completeDaysStr == "null" {
			aims[i].CompleteDaysCount = 0
		} else {
			// Remove curly braces and split the string by commas
			dates := strings.Split(strings.Trim(completeDaysStr, "{}"), ",")
			// Count the number of dates
			aims[i].CompleteDaysCount = len(dates)
		}

	}
	// Filter active aims
	var activeAims []Aim
	for _, aim := range aims {
		if !IsAimActive(aim.Startday, aim.Endday) {
			activeAims = append(activeAims, aim)
		}
	}
	return c.Status(fiber.StatusOK).JSON(activeAims)

}
