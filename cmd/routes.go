package main


import (
	"github.com/damlaYasarr/aliveApp/handlers"
	"github.com/gofiber/fiber/v2"
)
func xRoutes(app *fiber.App) {
	app.Get("/listallAims", handlers.ListUsersAim)

	app.Post("/postemail", handlers.RegisterUser)

	app.Post("/addnewaim", handlers.AddNewAim)

	//here is control function
    app.Get("/donotduplicatename", handlers.Donotduplicatename)
	
   app.Get("/activedays", handlers.GETActivedays)


	//app.Put("/edithabit",  handlers.EditHabitName)
	app.Delete("/deletehabit/:name", handlers.DeleteHabitByName)

}