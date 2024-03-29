package crawler

import (
	"errors"
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/xuri/excelize/v2"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"strings"
	"time"
)

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

	// Read 2nd row (headers)
	rows.Next()
	rows.Next()
	headers, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Populate header map with relative column index (build using odd column text)
	// substitute even text with previous column text
	for i, j := range headers {
		if i%2 == 0 {
			columnKeyMap[i] = getFirstNChar(strings.TrimSpace(j), 3) // If populated add to array
		} else {
			columnKeyMap[i] = getFirstNChar(strings.TrimSpace(headers[i-1]), 3) // Else set equal to previous cell
		}
	}

	// For each sequent row aggregate by header and populate map
	// For each value EVEN index are vehicle callsign and ODD index are daily availability
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
			for i, convItem := range valueConv {
				if i%2 == 0 {
					row := database.Fleet{}
					row.Ente, row.Stazionamento, _, row.Minimum = parseVehicle(convItem)
					row.Convenzione = keyConv
					row.ActiveFrom, err = utils.ConvertLabelToTime(s)
					row.ActiveDays = valueConv[i+1]
					if err != nil {
						return err
					}
					contentObject = append(contentObject, row)
				}
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
	err = fleetDbUpdate(fmt.Sprintf("%sfleet.xlsx", fleetPath))
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
				var err error
				log.Printf(" - | Fleet watcher event |\t%v", event)
				if event.FileInfo.Name() == "fleet.xlsx" {
					err = fleetDbUpdate(event.Path)
					if err != nil {
						log.Printf("Error updating fleet db: %s", err)
					}
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Add watch to folder
	err = w.Add(fleetPath)
	if err != nil {
		return err
	}

	// Start the watcher
	go w.Start(time.Second * 10)

	return nil
}
