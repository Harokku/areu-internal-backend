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

// CheckActualAvailability Check actual fleet state vs theoretical state and return if all ok
func (f Fleet) CheckActualAvailability(ctx *fiber.Ctx) error {
	var (
		err                       error
		actualTheoreticFleetState []database.Fleet
		actualFleetState          []database.BacoSnapshoot
		missingFleet              []database.BacoSnapshoot
	)

	// Retrieve theoretical actual fleet state based on time frame
	err = database.Fleet{}.GetActiveNow(&actualTheoreticFleetState)
	if err != nil {
		log.Printf(ErrStringMsg("fleet/CheckActualAvailability while retrieving actual time range", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Retrive actual fleet state from Baco server
	err = database.BacoSnapshoot{}.GetSnapshoot(&actualFleetState)
	if err != nil {
		log.Printf(ErrStringMsg("fleet/CheckActualAvailability while retrieving actual fleet state", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Cycle actualTheoreticFleetState anf check if all entry are present in actualFleetState
	// If an entry is present in both, delete from actualTheoreticFleetState
	// If an entry is missing from actualFleetState add it to missingFleet and delete from actualTheoreticFleetState

	for _, fleetItem := range actualTheoreticFleetState {
		var (
			found          bool   // if actual item had been found in actualFleetState
			foundedIndex   int    // actualFleetState founded item index
			theoreticLotto string // actual theoretic item's lotto to be searched
		)
		_, theoreticLotto = ExtractLotto(fleetItem.Ente)

		for iActual, snapshotItem := range actualFleetState {
			var actualLotto string
			_, actualLotto = ExtractLotto(snapshotItem.Ente)
			// If found set flag to true and index to actual item to be removed
			if actualLotto == theoreticLotto {
				found = true
				foundedIndex = iActual
				// Item found stop searching the rest of slice
				break
			}
		}

		// If found remove the item from snapshot, or add to missing list
		if found {
			actualFleetState = append(actualFleetState[:foundedIndex], actualFleetState[foundedIndex+1:]...)
		} else {
			missingFleet = append(missingFleet, database.BacoSnapshoot{
				Convenzione: fleetItem.Convenzione,
				Ente:        fleetItem.Ente,
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"status":     "success",
		"message":    "Fleet check completed",
		"additional": actualFleetState,
		"missing":    missingFleet,
	})
}
