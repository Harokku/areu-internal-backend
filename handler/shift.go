package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"internal-backend/utils"
	"log"
)

type Shift struct {
}

//checkParam check if passed params are present in the list of accepted
func checkParam(pName, pType string) bool {
	var (
		nameValidated bool
		typeValidated bool
	)
	nameValidated = false
	typeValidated = false

	// Check if name is valid
	switch pName {
	case "tecnici",
		"infermieri",
		"medici":
		nameValidated = true
	}

	// Check if type is valid
	switch pType {
	case "turni",
		"postazioni":
		typeValidated = true
	}

	// If both are validated return true
	if nameValidated && typeValidated {
		return true
	}

	// If something are not valid return false
	return false
}

func (s Shift) ServeByPath(ctx *fiber.Ctx) error {
	var (
		err          error
		shiftName    string //Shift name to retrieve
		shiftType    string //Shift type to retrieve
		shiftRoot    string //Shift root from env
		downloadPath string //Path to download file from
		fileName     string //Name of downloaded file
	)

	shiftRoot, err = utils.ReadEnv("SHIFT_ROOT")
	if err != nil {
		log.Printf(ErrStringMsg("shift/ServeByPath while reading env", err))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// Retrieve name and type from url
	shiftName = ctx.Params("name")
	shiftType = ctx.Params("type")

	// Check if params are valid
	if !checkParam(shiftName, shiftType) {
		log.Printf(ErrString("shift/ServeByPath while checking for param validity"))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Build download path and filename
	downloadPath = fmt.Sprintf("%s/%s/%s.pdf", shiftRoot, shiftName, shiftType)
	fileName = fmt.Sprintf("%s_%s.pdf", shiftType, shiftName)

	// Send file to the client
	return ctx.Download(downloadPath, fileName)
}
