package main

import (
    "os"
    "github.com/damlaYasarr/aliveApp/database"
    "github.com/gofiber/fiber/v2"
    "github.com/damlaYasarr/aliveApp/middleware"
)

func main() {
    database.ConnectDb()

	app := fiber.New()

	xRoutes(app)
 // Register routes
    app.Get("/", initialize)
	app.Listen(":3000")
    key := os.Getenv("NOTIF_TOKEN")
    deviceTokens := []string{
		key,
		// Add more tokens as needed
	}

	err := SendPushNotification(deviceTokens)
	if err != nil {
		log.Fatalf("Error sending push notification: %v", err)
	}

}

func initialize(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

