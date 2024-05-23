package test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
)

// Mock database and GetUserIDByEmail function
var database struct {
    DB struct {
        Db MockDb
    }
}

type MockDb struct{}

func (db MockDb) Raw(query string, values ...interface{}) *MockDbResult {
    return &MockDbResult{}
}

type MockDbResult struct {
    Error error
}

func (r *MockDbResult) Scan(dest interface{}) *MockDbResult {
    aims := []struct {
        ID            int64  `json:"id"`
        COMPLETE_DAYS string `json:"complete_days"`
    }{
        {ID: 1, COMPLETE_DAYS: "{2023-05-22,2023-05-23}"},
        {ID: 2, COMPLETE_DAYS: "{2023-05-24,2023-05-25}"},
    }
    destPtr := dest.(*[]struct {
        ID            int64  `json:"id"`
        COMPLETE_DAYS string `json:"complete_days"`
    })
    *destPtr = aims
    return r
}

func GetUserIDByEmail(email string) (int, error) {
    return 7, nil // Mocked user ID
}

// Mock function to simulate database interaction
func GETActivedays(c *fiber.Ctx) error {
    email := c.Query("email")
    if email == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
    }

    userID, err := GetUserIDByEmail(email)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
    }

    // Retrieve all aim and time data related to the user ID
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

func TestGETActivedays(t *testing.T) {
    app := fiber.New()

    // Register the handler function for the /activedays route
    app.Get("/activedays", GETActivedays)

    // Create a mock HTTP request
    req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/activedays?email=damlaprotel17@gmail.com", nil)

    // Process the HTTP request
    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("Failed to perform request: %v", err)
    }
    defer resp.Body.Close()

    // Check the expected status code
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK; got %d", resp.StatusCode)
    }

    // Read the response body
    var aims []struct {
        ID            int64  `json:"id"`
        COMPLETE_DAYS string `json:"complete_days"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&aims); err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }

    // Check the expected response
    if len(aims) != 2 {
        t.Fatalf("Expected 2 aims; got %d", len(aims))
    }

    expectedCompleteDays := "{2023-05-22,2023-05-23}"
    if aims[0].COMPLETE_DAYS != expectedCompleteDays {
        t.Errorf("Expected complete days %q; got %q", expectedCompleteDays, aims[0].COMPLETE_DAYS)
    }
}
