package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"time"
)

type Fleet struct {
}

func (f Fleet) GetAll(ctx *fiber.Ctx) error {
	var (
		d   []database.Fleet
		err error
	)

	// Retrieve all content
	err = database.Fleet{}.GetAll(&d)
	if err != nil {
		log.Printf(ErrStringMsg("fleet/GetAll while retrieving all content", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all fleet data",
		"retrieved": len(d),
		"data":      d,
	})
}

func (f Fleet) GetActualTimeRange(ctx *fiber.Ctx) error {
	var (
		d   []database.Fleet
		err error
	)
	err = database.Fleet{}.GetActiveNow(&d)
	if err != nil {
		log.Printf(ErrStringMsg("fleet/GetActualTimeRAnge while retrieving actual time range", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved actual time range fleet data",
		"retrieved": len(d),
		"data":      d,
	})
}

func (f Fleet) LogExecutedCheck(ctx *fiber.Ctx) error {
	var (
		e   utils.Entry
		err error
	)

	err = ctx.BodyParser(&e)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Add actual timestamp
	e.Timestamp = time.Now()

	// Log entry
	err = e.WriteEntry()
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
