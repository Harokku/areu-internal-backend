package handler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"strings"
)

type Content struct {
}

//GetAll retrieve content index (all links and display name)
func (c Content) GetAll(ctx *fiber.Ctx) error {
	var (
		d   []database.Content
		err error
	)

	//Retrieve all contents
	err = database.Content{}.GetAll(&d)
	if err != nil {
		log.Printf(ErrStringMsg("content/GetAll while retrieving all content", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all content indexes",
		"retrieved": len(d),
		"data":      d,
	})
}

//GetContent retrieve content from XLSX sheet from link
func (c Content) GetContent(ctx *fiber.Ctx) error {
	var (
		f       *excelize.File
		path    string                 //Data table path on disk (from env)
		link    string                 //XLSX sheet number to read from
		d       database.Content       //Database content calls
		keys    []string               //JSON keys from query result (column names)
		dataRow map[string]interface{} //A single row from result, to be added to final JSON
		data    []interface{}          //Final JSON marshalable to be sent to client
		err     error
	)

	// -------------------------
	// Read data table path from env
	// -------------------------

	path, err = utils.ReadEnv("DATA_TABLE")
	if err != nil {
		log.Fatalf("Error retrieving data table from env")
	}
	path = fmt.Sprintf("%s/data-table.xlsx", path)

	// -------------------------
	// Open XLSX file and extract information from url param
	// -------------------------

	//Try to parse link from url, return bad param otherwise
	link = ctx.Params("link")
	if link == "" {
		log.Printf(ErrString("content/GetContent while parsing input from body"))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	//Retrieve sheet name from DB
	err = d.GetDisplayName(link)
	if err != nil {
		log.Printf(ErrString("content/GetContent while retrieving content"))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	//Read sheet and return data
	f, err = excelize.OpenFile(path)
	if err != nil {
		log.Printf(ErrString("content/GetContent while opening XLSX"))
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	rows, err := f.GetRows(d.DisplayName)
	if err != nil {
		log.Printf(ErrString("content/GetContent while reading data from XLSX"))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	//Extract keys (columns names) from rows
	//According to format definition 2nd row of XLSX file contain column names
	for _, s := range rows[1] {
		keys = append(keys, strings.ToLower(s))
	}

	//Extract JSON item from each row and build map
	for i, row := range rows {
		//According to format definition skip first 2 row (heading and column definition)
		//Data start on 3rd row
		if i > 1 {
			dataRow = make(map[string]interface{})
			// For each cell create key/value with Column name and cell content
			for i2, cell := range row {
				dataRow[keys[i2]] = cell
			}
			data = append(data, dataRow)
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Retrieved content",
		"keys":    keys,
		"data":    data,
	})
}
