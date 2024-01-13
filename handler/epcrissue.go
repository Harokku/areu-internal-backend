package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"log"
	"net"
)

type EpcrIssue struct {
}

// PostEpctIssue receives a POST request and adds an ePCR issue to the database.
// It takes a fiber.Ctx argument.
// The request body is parsed into an EpcrIssue struct.
// If parsing fails, a fiber.StatusBadRequest response is sent.
// The IP address from the request context is assigned to the issue.
// The issue is then added to the database using the PostIssue method.
// If adding the record to the database fails, an error is logged.
// Finally, a JSON response is returned indicating the success and message.
func (i EpcrIssue) PostEpctIssue(ctx *fiber.Ctx) error {
	var (
		err   error
		issue database.EpcrIssue
	)

	err = ctx.BodyParser(&issue)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	issue.IpAddress = net.IP(ctx.IP())

	err = issue.PostIssue()
	if err != nil {
		log.Printf(ErrStringMsg("epcrissue/PostEpctIssue while adding record to db", err))
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "epcr issue added to db",
	})
}
