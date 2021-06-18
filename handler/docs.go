package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"internal-backend/websocket"
	"log"
)

type Docs struct {
}

//GetAll get all documents info from db
func (d Docs) GetAll(ctx *fiber.Ctx) error {
	var (
		err       error
		documents []database.Document
	)

	// Retrieve all documents
	err = database.Document{}.GetAll(&documents)
	if err != nil {
		log.Printf(ErrStringMsg("docs/GetAll while retrieving all documents", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	//TODO: Remove in production DEBUG only
	websocket.Broadcast <- fiber.Map{
		"data": documents,
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all docs",
		"retrieved": len(documents),
		"data":      documents,
	})
}

//GetById get single document info by id (id from param)
func (d Docs) GetById(ctx *fiber.Ctx) error {
	var (
		err      error
		document database.Document
	)

	// Retrieve document by id
	err = document.GetById(ctx.Params("id"))
	if err != nil {
		log.Printf(ErrStringMsg("docs/GetById while retrieving document", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Retrieved document",
		"data":    document,
	})
}

//ServeById actually retrieve file from server by DB id (id from param)
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

	// Send file to client
	return ctx.Download(dInfo.FileName, dInfo.DisplayName)
}
