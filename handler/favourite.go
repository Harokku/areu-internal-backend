package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"log"
)

type Favourite struct {
}

// GetAll get all favourites info from db
func (f Favourite) GetAll(ctx *fiber.Ctx) error {
	var (
		err        error
		favourites []database.Favourite
	)

	// Retrieve all favourites
	err = database.Favourite{}.GetAll(&favourites)
	if err != nil {
		log.Printf(ErrStringMsg("favourites/GetAll while retrieving all favourites", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all favourites",
		"retrieved": len(favourites),
		"data":      favourites,
	})
}

// GetAggregatedByIp get aggregated favourites info from db by ip
// pass "own" as parameter to get only favourites from your own ip
func (f Favourite) GetAggregatedByIp(ctx *fiber.Ctx) error {
	var (
		err        error
		ip         string // ip to retrieve favourites from
		favourites []database.Favourite
	)

	// Try to parse requested ip, if it's not "own"
	ip = ctx.Params("ip")
	if ip == "own" {
		ip = ctx.IP()
	} else {
		// Validate ip
		if !validateIp(ip) {
			// invalid ip - log it and return 400
			log.Printf(ErrStringMsg("favourites/GetAggregatedByIp while validating ip", err))
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}

	// retrieve all favourites aggregated by ip
	err = database.Favourite{}.GetAggregatedByConsoleIp(ip, &favourites)
	if err != nil {
		log.Printf(ErrStringMsg("favourites/GetAggregatedByIp while retrieving all favourites aggregated by ip", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all favourites aggregated by ip",
		"retrieved": len(favourites),
		"data":      favourites,
	})
}

// GetAggregatedByFunctionFromIp get aggregated favourites info from db by function using passed ip
// pass "own" as parameter to get only favourites from your own ip
func (f Favourite) GetAggregatedByFunctionFromIp(ctx *fiber.Ctx) error {
	var (
		err        error
		ip         string // ip to retrieve favourites from
		function   string // function to retrieve favourites from
		favourites []database.Favourite
	)

	// Try to parse requested ip, if it's not "own"
	ip = ctx.Params("ip")
	if ip == "own" {
		ip = ctx.IP()
	} else {
		// Validate ip
		if !validateIp(ip) {
			// invalid ip - log it and return 400
			log.Printf(ErrStringMsg("favourites/GetAggregatedByFunctionFromIp while validating ip", err))
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
	}

	// retrieve function from ip
	function, err = database.Favourite{}.GetFunctionByIp(ip)
	if err != nil {
		log.Printf(ErrStringMsg("favourites/GetAggregatedByFunctionFromIp while retrieving function from ip", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// retrieve all favourites aggregated by function
	err = database.Favourite{}.GetAggregatedByFunction(function, &favourites)
	if err != nil {
		log.Printf(ErrStringMsg("favourites/GetAggregatedByFunctionFromIp while retrieving all favourites aggregated by function", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all favourites aggregated by function",
		"retrieved": len(favourites),
		"data":      favourites,
	})
}
