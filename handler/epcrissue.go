package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"internal-backend/database"
	"log"
	"net"
	"strconv"
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

func (i EpcrIssue) GenerateAndDownloadReport(ctx *fiber.Ctx) error {
	issues, err := database.EpcrIssue{}.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	f := excelize.NewFile()
	index, err := f.NewSheet("ePCR-Issues")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	f.SetActiveSheet(index)

	// Write headers to xlsx
	f.SetCellValue("ePCR-Issues", "A1", "Data e ora")
	f.SetCellValue("ePCR-Issues", "B1", "MSB")
	f.SetCellValue("ePCR-Issues", "C1", "Problema riscontrato")
	f.SetCellValue("ePCR-Issues", "D1", "Ip Operatore che rileva")

	for i2, issue := range issues {
		f.SetCellValue("ePCR-Issues", "A"+strconv.Itoa(i2+2), issue.Timestamp)
		f.SetCellValue("ePCR-Issues", "B"+strconv.Itoa(i2+2), issue.VehicleId)
		f.SetCellValue("ePCR-Issues", "C"+strconv.Itoa(i2+2), issue.Text)
		f.SetCellValue("ePCR-Issues", "D"+strconv.Itoa(i2+2), issue.IpAddress.To4().String())
	}

	// Delete default Sheet1
	err = f.DeleteSheet("Sheet1")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Set headers for file download
	ctx.Set(fiber.HeaderContentDisposition, `attachment; filename="EpcrIssues.xlsx"`)
	ctx.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write excel file content to fiber response
	if err := f.Write(ctx.Response().BodyWriter()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}
