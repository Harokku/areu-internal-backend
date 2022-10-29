package crawler

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"internal-backend/utils"
	"internal-backend/websocket"
	"log"
	"strings"
	"time"
)

func WatchRootFromEnv() error {
	var (
		w            *watcher.Watcher
		err          error
		docRoot      string   //Documents disk path from env
		docRootArray []string //Documents array og path to check
	)

	// -------------------------
	// Read .env
	// -------------------------
	docRoot, err = utils.ReadEnv("DOC_ROOT")
	if err != nil {
		log.Fatalf("Error retrieving documents root from env")
	}
	docRootArray = strings.Split(docRoot, "|")

	w = watcher.New()

	// -------------------------
	// watcher config
	// -------------------------

	w.SetMaxEvents(2)
	w.IgnoreHiddenFiles(true)
	w.FilterOps(watcher.Create, watcher.Remove, watcher.Rename, watcher.Move)

	// start the event loop
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Printf(" - | FileWatcher Event |\t%v", event)
				// Re-enumerate watched documents
				EnumerateDocuments()
				// Calculate hash based on filepath to enable file retrieval
				hash, err := getSha1(event.Path)
				if err != nil {
					log.Printf("[ERR]\tError calculating hash after filwatche event")
				}
				websocket.Broadcast <- map[string]interface{}{
					"id":        "Filewatcher event",
					"operation": fmt.Sprint(event.Op),
					"filename":  fmt.Sprint(event.FileInfo.Name()),
					"hash":      hash,
				}
			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// TODO: Implement not hardcoded version
	// Set folder to avoid
	err = w.Ignore(`y:\Docs\Docs_SRLombardia\Doc per RAR SOREU`)
	if err != nil {
		log.Printf("[ERR] - Error adding path to file watcher ignored list:\t%v", err)
	}

	// Cycle through every folder and add it and subfolder to watcher
	for _, rootFolder := range docRootArray {
		if err = w.AddRecursive(rootFolder); err != nil {
			log.Printf("[WARN] - Error adding folder to watch list: %v", err)
		}
	}

	// Start the watching process - it'll check for changes every 1000ms.
	go w.Start(time.Millisecond * 1000)

	return nil
}
