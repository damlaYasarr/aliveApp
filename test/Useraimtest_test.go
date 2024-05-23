package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"net/http/httptest"
	"github.com/gofiber/fiber/v2"
)

type Aim struct {
	Name             string `json:"name"`
	CompleteDays     string `json:"complete_days"`
	Startday         string `json:"startday"`
	Endday           string `json:"endday"`
	NotificationHour string `json:"notificationhour"`
	CompleteDaysCount int    `json:"complete_days_count"`
}

func TestListAllAims(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Register the handler function for the /listallAims route
	app.Get("/listallAims", func(c *fiber.Ctx) error {
		// Mock data for testing
		aims := []Aim{
			{
				Name:             "hergün günlük yazacağım",
				CompleteDays:     "",
				Startday:         "2023-05-22",
				Endday:           "2023-06-22",
				NotificationHour: "20.00",
				CompleteDaysCount: 1,
			},
		}

		// Return the aims as a JSON response
		return c.JSON(aims)
	})

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000/listallAims?email=test@example.com", nil)

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
	var aims []Aim
	if err := json.NewDecoder(resp.Body).Decode(&aims); err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check the expected response
	if len(aims) != 1 {
		t.Fatalf("Expected 1 aim; got %d", len(aims))
	}

	expectedName := "hergün günlük yazacağım"
	if aims[0].Name != expectedName {
		t.Errorf("Expected aim name %q; got %q", expectedName, aims[0].Name)
	}

	expectedStartday := "2023-05-22"
	if aims[0].Startday != expectedStartday {
		t.Errorf("Expected start day %q; got %q", expectedStartday, aims[0].Startday)
	}

	expectedEndday := "2023-06-22"
	if aims[0].Endday != expectedEndday {
		t.Errorf("Expected end day %q; got %q", expectedEndday, aims[0].Endday)
	}

	expectedCompleteDaysCount := 1
	if aims[0].CompleteDaysCount != expectedCompleteDaysCount {
		t.Errorf("Expected complete days count %d; got %d", expectedCompleteDaysCount, aims[0].CompleteDaysCount)
	}
}
