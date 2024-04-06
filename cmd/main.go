package main

import (
   
    "github.com/damlaYasarr/aliveApp/database"
    "github.com/gofiber/fiber/v2"
)

func main() {
    database.ConnectDb()

	app := fiber.New()

	xRoutes(app)
 // Register routes
    app.Get("/", initialize)
	app.Listen(":3000")



}

func initialize(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

