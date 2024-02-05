package crawler

import (
	"errors"
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/xuri/excelize/v2"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"time"
)

func newsDbUpdate(file string) error {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}

	// Read sheet list name
	sheetList := f.GetSheetList()

	// Read all rows of the 1st sheet
	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	var res []database.NewsFeed
	for i, row := range rows {
		// Skip first row
		if i == 0 {
			continue
		}

		temp := database.NewsFeed{
			Timestamp: row[0],
			Title:     row[1],
			News:      row[2],
		}

		//temp := map[string]interface{}{
		//	"timestamp": row[0],
		//	"title":     row[1],
		//	"news":      row[2],
		//}
		res = append(res, temp)
	}

	// Truncate news table for fresh start
	err = database.NewsFeed{}.TruncateTable()
	if err != nil {
		return errors.New(fmt.Sprintf("newsfeedwatcher/newsfeedDbUpdate returned error while truncating table: %v\n", err))
	}

	// Add content to db
	err = database.NewsFeed{}.BulkCreate(res)
	if err != nil {
		return errors.New(fmt.Sprintf("newsfeedwatcher/newsfeedDbUpdate returned error while bulk creating table: %v\n", err))
	}

	return nil
}

func WatchNewsTableFromEnv() error {

	// -------------------------
	// Read .env
	// -------------------------
	newsTablePath, err := utils.ReadEnv("NEWS_ROOT")
	if err != nil {
		log.Fatalf("Error retrieving data table from env")
	}

	// -------------------------
	// Create starting DB
	// -------------------------
	err = newsDbUpdate(fmt.Sprintf("%s/news.xlsx", newsTablePath))
	if err != nil {
		return err
	}

	w := watcher.New()

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
				log.Printf(" - | NewsFeed watcher event |\t%v", event)
				newsDbUpdate(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Add watch to folder
	err = w.Add(fmt.Sprintf("%s/news.xlsx", newsTablePath))
	if err != nil {
		return err
	}

	// Start the watcher
	go w.Start(time.Second * 10)

	return nil
}
