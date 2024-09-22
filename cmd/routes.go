package main

import (
	"github.com/damlaYasarr/aliveApp/handlers"
	"github.com/gofiber/fiber/v2"
)

func xRoutes(app *fiber.App) {
	// notification
	app.Get("/send-notification", sendNotificationHandler)
	// approve the completed days
	app.Put("/approvaltime", handlers.ApprovalHabitDate)

	// list all aims
	app.Get("/listallAims", handlers.ListUsersAim)

	// register user with email or sign in
	app.Post("/postemail", handlers.RegisterUser)

	// do not duplicate aim name -- it is inactive
	app.Get("/donotduplicatename", handlers.Donotduplicatename)

	// get active days of users
	app.Get("/activedays", handlers.GETActivedays)

	// delete the habits
	app.Delete("/deletehabit", handlers.DeleteUserAim)

	// add new aim
	app.Post("/addnewaim", handlers.AddNewAim)

	// list all ACTIVE habits
	app.Get("/listactivedays", handlers.ListUsersActiveAim)

	// send the users info to the ai and get the interpretation
	app.Get("/getfeedback", handlers.ReturnFeedBack)

}
