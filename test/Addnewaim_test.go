package test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
)

// Aim struct should match the one in your main code
type Aim struct {
    Name              string `json:"name"`
    CompleteDays      string `json:"complete_days"`
    Startday          string `json:"startday"`
    Endday            string `json:"endday"`
    NotificationHour  string `json:"notificationhour"`
    CompleteDaysCount int    `json:"complete_days_count"`
}

// Mock function to simulate GetUserIDByEmail
func GetUserIDByEmail(email string) (int, error) {
    return 1, nil
}

// Mock function to simulate database interaction
func AddNewAim(c *fiber.Ctx) error {
    var requestBody struct {
        Email         string `json:"email"`
        Aim           string `json:"aim"`
        AimDate       string `json:"aim_date"`
        Endday        string `json:"endday"`
        Notification  string `json:"notif"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid request body")
    }

    _, err := GetUserIDByEmail(requestBody.Email)
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to get user ID")
    }

    // Mock user check and aim creation
    newAim := Aim{
        Name:             requestBody.Aim,
        Startday:         requestBody.AimDate,
        Endday:           requestBody.Endday,
        NotificationHour: requestBody.Notification,
        CompleteDaysCount: 0,
    }

    return c.Status(http.StatusOK).JSON(newAim)
}

func TestAddNewAim(t *testing.T) {
    app := fiber.New()

    // Register the handler function for the /addaim route
    app.Post("/addnewaim", AddNewAim)

    // Example test aim data
    testAim := `{
        "email": "damlaprotel17@gmail.com",
        "aim": "Read a book",
        "aim_date": "2024-01-01",
        "endday": "2024-12-31",
        "notif": "18:00"
    }`

    req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/addnewaim?email=damlaprotel17@gmail.com", bytes.NewBufferString(testAim))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("Failed to perform request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK; got %d", resp.StatusCode)
    }

    var responseBody Aim
    if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }

    expectedAim := "Read a book"
    if responseBody.Name != expectedAim {
        t.Errorf("Expected aim %q; got %q", expectedAim, responseBody.Name)
    }

    expectedAimDate := "2024-01-01"
    if responseBody.Startday != expectedAimDate {
        t.Errorf("Expected start day %q; got %q", expectedAimDate, responseBody.Startday)
    }

    expectedEndday := "2024-12-31"
    if responseBody.Endday != expectedEndday {
        t.Errorf("Expected end day %q; got %q", expectedEndday, responseBody.Endday)
    }

    expectedNotification := "18:00"
    if responseBody.NotificationHour != expectedNotification {
        t.Errorf("Expected notification hour %q; got %q", expectedNotification, responseBody.NotificationHour)
    }
}
