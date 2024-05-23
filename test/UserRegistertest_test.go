package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterUser(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Register the handler function for the /postemail route
	app.Post("/postemail", func(c *fiber.Ctx) error {
		// Define the structure of the request body
		var requestBody struct {
			Email string `json:"email"`
		}

		// Parse the request body into the defined structure
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid request body")
		}

		// Return the parsed email as a JSON response
		return c.JSON(fiber.Map{"email": requestBody.Email})
	})

	// Example test user
	testUser := `{"email": "damla@gmail.com"}`

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/postemail", bytes.NewBufferString(testUser))
	req.Header.Set("Content-Type", "application/json")

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
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check the expected response
	expectedEmail := "damla@gmail.com"
	actualEmail, ok := responseBody["email"].(string)
	if !ok {
		t.Fatalf("Expected email field is not a string")
	}
	if actualEmail != expectedEmail {
		t.Errorf("Expected email %q; got %q", expectedEmail, actualEmail)
	}
}
