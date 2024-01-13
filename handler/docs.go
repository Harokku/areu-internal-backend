package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
)

type Docs struct {
}

// GetAll get all documents info from db
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

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all docs",
		"retrieved": len(documents),
		"data":      documents,
	})
}

// GetById get single document info by id (id from param)
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

// ServeById actually retrieve file from server by DB id (id from param)
func (d Docs) ServeById(ctx *fiber.Ctx) error {
	var (
		err      error
		id       string            //Document id to retrieve
		dInfo    database.Document //Document info retrieved from bd
		blobName string            //Blob filename extracted from path for mime association
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

	// Build blob name
	blobName = fmt.Sprintf("%s%s", dInfo.DisplayName, filepath.Ext(dInfo.FileName))
	fmt.Println(blobName)

	// Send file to client
	return ctx.Download(filepath.FromSlash(dInfo.FileName), blobName)
}

// ServeByHash actually retrieve file from server by DB hash (hash from param)
func (d Docs) ServeByHash(ctx *fiber.Ctx) error {
	var (
		err        error
		hash       string             //Document id to retrieve
		dInfo      database.Document  //Document info retrieved from bd
		dFavourite database.Favourite //Favourite info retrieved from bd
		blobName   string             //Blob filename extracted from path for mime association
	)

	hash = ctx.Params("id")
	// Try to parse id url, return bad request otherwise
	if hash == "" {
		log.Printf(ErrString("docs/ServeByHash while parsing input from body"))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Retrieve document meta from db
	err = dInfo.GetByHash(hash)
	if err != nil {
		log.Printf(ErrStringMsg("docs/ServeByHash while retrieving document", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Build dFavourite item to post to db
	dFavourite.ConsoleIp = ctx.IP()
	dFavourite.Filename = dInfo.DisplayName

	// Add favourite to db
	err = dFavourite.PostFavourite()
	if err != nil {
		log.Printf(ErrStringMsg("docs/ServeByHash while adding favourite", err))
		// Don't return error, just log it and serve file to client
	}

	// Build blob name
	blobName = fmt.Sprintf("%s%s", dInfo.DisplayName, filepath.Ext(dInfo.FileName))

	// Send file to client
	return ctx.Download(filepath.FromSlash(dInfo.FileName), blobName)
}

// ServeByFilename actually retrieve file from server by filename (filename url encoded from param)
func (d Docs) ServeByFilename(ctx *fiber.Ctx) error {
	var (
		err        error
		filename   string             //Document filename to retrieve
		dInfo      database.Document  //Document info retrieved from bd
		dFavourite database.Favourite //Favourite info retrieved from bd
	)

	filename, err = url.QueryUnescape(ctx.Params("filename"))
	if err != nil {
		log.Printf(ErrString("docs/ServeByFilename while parsing input from body"))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Retrieve document meta from db
	err = dInfo.GetByFilename(filename)
	if err != nil {
		log.Printf(ErrStringMsg("docs/ServeByFilename while retrieving document", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Build dFavourite item to post to db
	dFavourite.ConsoleIp = ctx.IP()
	dFavourite.Filename = dInfo.DisplayName

	// Add favourite to db
	err = dFavourite.PostFavourite()
	if err != nil {
		log.Printf(ErrStringMsg("docs/ServeByHash while adding favourite", err))
		// Don't return error, just log it and serve file to client
	}

	// Send file to client
	return ctx.Download(filepath.FromSlash(dInfo.FileName), dInfo.DisplayName)
}

// GetRecent get most {num} recent documents
func (d Docs) GetRecent(ctx *fiber.Ctx) error {
	var (
		err       error
		num       int                 //How many document to retrieve
		mode      string              //Define which aggregation method to use (all, by category...), default all
		documents []database.Document //Document info retrieved from db
	)

	// Try to parse num url, return bad request otherwise
	num, err = strconv.Atoi(ctx.Params("num"))
	if err != nil {
		log.Printf(ErrString("docs/GetRecent while parsing input from body"))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Try to parse aggregation mode, default all if null or malformed
	mode = ctx.Query("mode", "all")
	if mode != "all" && mode != "split" {
		mode = "all"
	}

	// Retrieve last {num} documents
	err = database.Document{}.GetRecent(num, mode, &documents)
	if err != nil {
		log.Printf(ErrStringMsg("docs/GetRecent while retrieving recent documents", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved most recent documents",
		"retrieved": len(documents),
		"data":      documents,
	})
}
