package main

import (
	"context"
	"fmt"
	"log"

	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/utils"
    "github.com/damlaYasarr/aliveApp/handlers"
	"github.com/gofiber/fiber/v2"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

var app *firebase.App
var ctx context.Context
var client *messaging.Client

func main() {
	// Connect to the database
	database.ConnectDb()

	// Initialize Firebase
	app, ctx, client = utils.SetupFirebase()

	// Create a new Fiber app
	app := fiber.New()

	// Define routes
	defineRoutes(app)
    xRoutes(app)
	// Start listening on port 3000
	app.Get("/", initialize)
	app.Listen(":3000")
}

// Define HTTP routes
func defineRoutes(app *fiber.App) {
	app.Get("/send-notification", sendNotificationHandler)
}

// Send a notification to a specific token
func sendNotificationHandler(c *fiber.Ctx) error {
	// Retrieve token from query parameter
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token is required"})
	}

	// Retrieve active aims
	activeAims, err := handlers.ListActiveHabits("damlaprotel17@gmail.com")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to retrieve active aims: %v", err)})
	}

	// Iterate through active aims
	for _, aim := range activeAims {
		// Check if aim is active at current time
		if handlers.IsAimActiveAtCurrentTime(aim.Startday, aim.Endday, aim.NotificationHour) {
			// Send notification immediately
			err := sendNotification(ctx, client, token, "ALive!!", fmt.Sprintf(aim.Name))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to send notification: %v", err)})
			}
			// Assuming you want to send notifications for each aim separately, you can return here if desired
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notification sent successfully"})
		}
	}

	// If no active aim found at current time
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "No active aim at this time"})
}

// Send a notification to a specific device token
func sendNotification(ctx context.Context, client *messaging.Client, token, title, body string) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
		return err
	}

	fmt.Printf("Message successfully sent: %v\n", response)
	return nil
}

// Initialize route handler
func initialize(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}