package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/damlaYasarr/aliveApp/database"

	"github.com/gofiber/fiber/v2"
)

// listUserAllHabit retrieves the user's habits by their email, using database.DB.Db as the connection.
func listUserAllHabit(email string) (string, error) {
	// Fetch the user ID by email
	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return "", err // Return empty string and the error
	}

	// Aims represents the structure for user habits.
	type Aimss struct {
		Name              string `json:"name"`
		CompleteDays      string `json:"complete_days"`
		StartDay          string `json:"startday"`
		EndDay            string `json:"endday"`
		NotificationHour  string `json:"notification_hour"`
		CompleteDaysCount int    `json:"complete_days_count"`
	}

	// Fetch aims and time information associated with a user ID
	var aims []Aimss
	rows, err := database.DB.Db.Raw(`
		SELECT a.name, t.complete_days, a.startday, a.endday, a.notification_hour
		FROM aims a
		JOIN times t ON a.id = t.aim_id
		WHERE a.user_id = ?`, userID).Rows()
	if err != nil {
		return "null", err
	}
	defer rows.Close()

	for rows.Next() {
		var aim Aimss
		if err := rows.Scan(&aim.Name, &aim.CompleteDays, &aim.StartDay, &aim.EndDay, &aim.NotificationHour); err != nil {
			return "", err
		}

		// Check if CompleteDays is null or empty
		if aim.CompleteDays == "" {
			return "null", nil // Return null for all values if CompleteDays is null
		}

		// Calculate the number of complete days
		aim.CompleteDaysCount = len(strings.Split(strings.Trim(aim.CompleteDays, "{}"), ","))
		aims = append(aims, aim)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return "", err
	}

	// Convert the habit list into a string format to send to the AI
	var feedback string
	for _, aim := range aims {
		feedback += fmt.Sprintf("Habit: %s, CompleteDays: %s, Start: %s, End: %s\n", aim.Name, aim.CompleteDays, aim.StartDay, aim.EndDay)
	}
	return feedback, nil
}

// FeedbackResponse represents the structure returned by the Python script
type FeedbackResponse struct {
	Response string `json:"response"`
}

// GetFeedBackByUsingAI calls the Python script and returns the AI feedback
func GetFeedBackByUsingAI(feedback string) (string, error) {
	// Prepare the Python command with the feedback as an argument
	cmd := exec.Command("/usr/src/app/venv/bin/python", "/usr/src/app/middleware/ai.py", feedback)

	// Create buffers to capture stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		// If there's an error, log stderr for debugging
		return "", fmt.Errorf("failed to execute AI script: %w: %s", err, stderr.String())
	}

	// Capture stdout (the actual output)
	output := out.Bytes()

	// Parse the JSON response from the Python script
	var aiResponse FeedbackResponse
	if err := json.Unmarshal(output, &aiResponse); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	return aiResponse.Response, nil
}

// ReturnFeedBack fetches user habits, sends them to AI, and returns feedback
func ReturnFeedBack(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	// Get user habit data
	feedback, err := listUserAllHabit(email)
	fmt.Print(feedback)
	if err != nil {
		log.Println("Error fetching user habits:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user habits"})
	}

	// Get AI feedback
	aiFeedback, err := GetFeedBackByUsingAI(feedback + "observe these scredule, duration and completed days. Can you help me how can I improve myself more?")
	fmt.Print(aiFeedback)
	if err != nil {
		log.Println("Error getting AI feedback:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get AI feedback"})
	}

	// Return AI feedback as JSON
	return c.JSON(fiber.Map{
		"ai_feedback": aiFeedback,
	})
}
