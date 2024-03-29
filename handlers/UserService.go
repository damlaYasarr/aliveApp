package handlers

import (
	"github.com/damlaYasarr/aliveApp/database"
	"github.com/damlaYasarr/aliveApp/models"
	"github.com/gofiber/fiber/v2"
)


func RegisterApp(c *fiber.Ctx) error {
	facts := []models.Fact{}
	database.DB.Db.Find(&facts)

	return c.Status(200).JSON(facts)
}


