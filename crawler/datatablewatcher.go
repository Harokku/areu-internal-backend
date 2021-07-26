package crawler

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/radovskyb/watcher"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"regexp"
	"time"
)

// sanitizeLink take a string and strip all whitespaces and special characters to be used as url part
func sanitizeLink(i string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		return ""
	}
	return re.ReplaceAllString(i, "")
}

// dataTableDbUpdate take an XLSX and extract content table reading all sheets
func dataTableDbUpdate(file string) error {
	var (
		f             *excelize.File     // excel file to read
		sheetList     []string           // List of sheets in excel file
		contentObject []database.Content //DB array to be bulk added
		err           error
	)
	log.Printf("\033[32mStarting\033[0m enumerating content...")
	enumerateStartTime := time.Now()

	// Open excel file
	f, err = excelize.OpenFile(file)
	if err != nil {
		return err
	}

	// Read all sheets name
	sheetList = f.GetSheetList()

	// Cycle all sheets and create content array
	for i, s := range sheetList {
		row := database.Content{}
		row.DisplayName = s
		row.Link = sanitizeLink(s)
		row.SheetNumber = i
		contentObject = append(contentObject, row)
	}

	// Truncate content table for fresh start
	err = database.Content{}.TruncateTable()
	if err != nil {
		return errors.New(fmt.Sprintf("datatablewatcher/dataTableDbUpdate returned error while truncating table: %v\n", err))
	}

	// Add contents do db
	err = database.Content{}.BulkCreate(contentObject)
	if err != nil {
		return errors.New(fmt.Sprintf("datatablewatcher/dataTableDbUpdate returned error while bulk creating table: %v\n", err))
	}

	enumerateDuration := time.Since(enumerateStartTime) //calculate total startup time
	log.Printf("Content \033[31menumerated\033[0m in %s", enumerateDuration)

	return nil
}

func WatchDataTableFromEnv() error {
	var (
		w             *watcher.Watcher
		err           error
		dataTablePath string
	)

	// -------------------------
	// Read .env
	// -------------------------
	dataTablePath, err = utils.ReadEnv("DATA_TABLE")
	if err != nil {
		log.Fatalf("Error retrieving data table from env")
	}

	// -------------------------
	// Create starting DB
	// -------------------------
	err = dataTableDbUpdate(fmt.Sprintf("%s/data-table.xlsx", dataTablePath))
	if err != nil {
		return err
	}

	w = watcher.New()

	// -------------------------
	// Watcher config
	// -------------------------
	w.SetMaxEvents(1)
	w.IgnoreHiddenFiles(true)
	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Printf(" - | DataTable watcher event |\t%v", event)
				dataTableDbUpdate(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Add watch to folder
	err = w.Add(fmt.Sprintf("%s/data-table.xlsx", dataTablePath))
	if err != nil {
		return err
	}

	// Start the watcher
	go w.Start(time.Second * 10)

	return nil
}
