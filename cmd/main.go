package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/handlers"
	"github.com/damlaYasarr/aliveApp/utils"
	"github.com/gofiber/fiber/v2"
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

	xRoutes(app)
	// Start listening on port 3000
	app.Get("/", initialize)
	app.Listen(":3000")
}

func sendNotificationHandler(c *fiber.Ctx) error {

	email := c.Query("email")
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token is required"})
	}

	activeAims, err := handlers.ListActiveHabits(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to retrieve active aims: %v", err)})
	}

	for _, aim := range activeAims {

		if handlers.IsAimActiveAtCurrentTime(aim.Startday, aim.Endday, aim.NotificationHour) {

			err := sendNotification(ctx, client, token, "ALive!!", fmt.Sprintf(aim.Name))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to send notification: %v", err)})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notification sent successfully"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "No active aim at this time"})
}

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
