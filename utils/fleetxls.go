package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"time"
)

type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operator  string    `json:"operator"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
}

func buildFilePath() string {
	var (
		path  string // Filepath calculated from env
		err   error
		year  int        // Actual year
		month time.Month // Actual month
	)
	// Initialize actual year and month
	year, month, _ = time.Now().Date()

	path, err = ReadEnv("FLEET_TABLE")
	if err != nil {
		log.Fatalf("Error retrieving fleet table from env")
	}
	// Convert path from env to forward slash (os agnostic) and add destination folder and filename from actual date
	path = fmt.Sprintf("%scompilati/%d-%s.xlsx", filepath.ToSlash(path), year, month)

	// Return os specific path format
	return filepath.FromSlash(path)
}

func (e Entry) WriteEntry() error {
	var (
		f         *excelize.File
		sheetList []string
		err       error
	)

	// Open actual file, if not exist, create a new one
	f, err = excelize.OpenFile(buildFilePath())
	if err != nil {
		f = excelize.NewFile()
		sheetList = f.GetSheetList()
		_ = f.SetSheetRow(sheetList[0], "A1", &[]interface{}{"Timestamp", "Operatore", "Stato controllo", "Note"})
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
		})
		_ = f.SetCellStyle(sheetList[0], "A1", "D1", style)
	}

	// Get all sheet list
	sheetList = f.GetSheetList()

	// Insert a new row on top of 1st sheet (after header)
	err = f.InsertRows(sheetList[0], 2, 1)
	if err != nil {
		return err
	}

	// Write entry to file
	err = f.SetSheetRow(sheetList[0], "A2", &[]interface{}{e.Timestamp.Unix(), e.Operator, e.Status, e.Note})
	if err != nil {
		return err
	}

	// Write file to disk
	err = f.SaveAs(buildFilePath())
	if err != nil {
		return err
	}

	return nil
}
