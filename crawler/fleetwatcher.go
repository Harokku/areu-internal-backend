package crawler

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/radovskyb/watcher"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"strings"
	"time"
)

// convertLabelToTime take a string from Excel sheet name and convert to a valid time object
func convertLabelToTime(label string) (time.Time, error) {
	var (
		parsedTime time.Time
		err        error
	)
	const timeOnly = "15.04"

	parsedTime, err = time.Parse(timeOnly, label)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func convertTimestampToTime(t time.Time) (time.Time, error) {
	var (
		h int //Hours
		m int //Minutes
	)

	h, m, _ = t.Clock()
	return convertLabelToTime(fmt.Sprintf("%v.%v", h, m))
}

// parseSheet parse Excel sheet and aggregate data by column header (1st row)
func parseSheet(f *excelize.File, sheetName string) (map[string][]string, error) {
	columnKeyMap := make(map[int]string)        // Hold column number mapped to aggregation header
	colAggregation := make(map[string][]string) // Actual aggregated data by column first 3 char header
	var err error

	// Read all sheet rows
	rows, err := f.Rows(sheetName)
	if err != nil {
		fmt.Println(err)
	}

	// Read 1st row (headers)
	rows.Next()
	headers, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Populate header map with relative column index (build using first 3 char of each column)
	for i, j := range headers {
		columnKeyMap[i] = getFirstNChar(strings.TrimSpace(j), 3)
	}

	// For each sequent row aggregate by header and populate map
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		for i, s := range row {
			if strings.TrimSpace(s) != "" {
				colAggregation[columnKeyMap[i]] = append(colAggregation[columnKeyMap[i]], s)
			}
		}
	}

	return colAggregation, nil
}

// fleetDbUpdate take an XLSX and extract content table reading all sheets
func fleetDbUpdate(file string) error {
	var (
		f             *excelize.File   // Excel file to read
		sheetList     []string         // List of sheets in Excel file
		contentObject []database.Fleet // DB array to be bulk added
		err           error
	)

	log.Printf("\033[32mStarting\033[0m enumerating fleet content...")
	enumerateStartTime := time.Now()

	// Open Excel file
	f, err = excelize.OpenFile(file)
	if err != nil {
		return err
	}

	// Read all sheets name
	sheetList = f.GetSheetList()

	// Cycle all sheets and create content array
	for _, s := range sheetList {
		// Parse give sheet and extract aggregation by column header
		sheetValueMap, err := parseSheet(f, s)
		if err != nil {
			return err
		}

		for keyConv, valueConv := range sheetValueMap {
			for _, convItem := range valueConv {
				row := database.Fleet{}
				row.Name = convItem
				row.ConvType = keyConv
				row.ActiveFrom, err = convertLabelToTime(s)
				if err != nil {
					return err
				}
				contentObject = append(contentObject, row)
			}
		}
	}

	// Truncate content table for fresh start
	err = database.Fleet{}.TruncateTable()
	if err != nil {
		return errors.New(fmt.Sprintf("fleetwatcher/fleetDbUpdate returned error while truncating table: %v\n", err))
	}

	// Add contents do db
	err = database.Fleet{}.BulkCreate(contentObject)
	if err != nil {
		return errors.New(fmt.Sprintf("fleetwatcher/fleetDbUpdate returned error while bulk creating table: %v\n", err))
	}

	enumerateDuration := time.Since(enumerateStartTime) //calculate total startup time
	log.Printf("Fleet \033[31menumerated\033[0m in %s", enumerateDuration)

	return nil
}

func WatchFleetFromEnv() error {
	var (
		w         *watcher.Watcher
		err       error
		fleetPath string
	)

	// -------------------------
	// Read .env
	// -------------------------
	fleetPath, err = utils.ReadEnv("FLEET_TABLE")
	if err != nil {
		log.Fatalf("Error retrieving fleet table from env")
	}

	// -------------------------
	// Create starting DB
	// -------------------------
	err = fleetDbUpdate(fmt.Sprintf("%s/fleet.xlsx", fleetPath))
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
				log.Printf(" - | Fleet watcher event |\t%v", event)
				fleetDbUpdate(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Add watch to folder
	err = w.Add(fmt.Sprintf("%s/fleet.xlsx", fleetPath))
	if err != nil {
		return err
	}

	// Start the watcher
	go w.Start(time.Second * 10)

	return nil
}
