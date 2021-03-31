package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"log"
)

type Docs struct {
}

func (d Docs) ServeById(ctx *fiber.Ctx) error {
	var (
		err   error
		id    string            //Document id to retrieve
		dInfo database.Document //Document info retrieved from bd
	)

	id = ctx.Params("id")
	// Try to parse id url, return bad request otherwise
	if id == "" {
		log.Printf(ErrString("docs/ServeById while parsing input from body"))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Retrieve document meta from db
	err = dInfo.GetById(id)
	if err != nil {
		log.Printf(ErrStringMsg("docs/ServeById while retrieving document", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return nil
}
