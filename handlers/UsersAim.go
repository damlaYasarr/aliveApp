package handlers

import (
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
)


func ListUsersAim(c *fiber.Ctx) error {
	facts := []models.Fact{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}
//add new aim
//delete aim
//arrange calendar
//arrange user as a premium
//get the all resuult
func PremiumUser(c *fiber.Ctx) error {
	facts := []models.Fact{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}