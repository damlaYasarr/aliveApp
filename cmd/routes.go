package main


import (
	"github.com/damlaYasarr/aliveApp/handlers"
	"github.com/gofiber/fiber/v2"
)
func xRoutes(app *fiber.App) {
	app.Get("/listallAims", handlers.ListUsersAim)

	app.Post("/postemail", handlers.RegisterAppWithEmail)

	app.Post("/addnewaim", handlers.AddnewAim)
}