package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
)

type News struct {
}

// GetAll retrieve all news
func (n News) GetAll(ctx *fiber.Ctx) error {
	var nf []database.NewsFeed

	// Retrieve all contents
	err := database.NewsFeed{}.GetAll(&nf)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving news",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all news",
		"retrieved": len(nf),
		"data":      nf,
	})
}
