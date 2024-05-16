package main


import (
	"github.com/damlaYasarr/aliveApp/handlers"
	"github.com/gofiber/fiber/v2"
)
func xRoutes(app *fiber.App) {
	app.Get("/listallAims", handlers.ListUsersAim)

	app.Post("/postemail", handlers.RegisterUser)

	app.Post("/addnewaim", handlers.AddNewAim)

	app.Get("/getall", handlers.ListUsersAllAim)
	
	//app.Put("/edithabit",  handlers.EditHabitName)
	app.Delete("/deletehabit/:name", handlers.DeleteHabitByName)

}