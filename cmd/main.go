package main

import ("github.com/gofiber/fiber/v2"

"github.com/damlaYasarr/aliveApp/database")

func main() {
	database.ConnectDb()
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("hüloğğ annem")
    })

    app.Listen(":3000")
}