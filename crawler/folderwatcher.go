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
				EnumerateDocuments()
				websocket.Broadcast <- map[string]interface{}{
					"id":        "Filewatcher event",
					"operation": fmt.Sprint(event.Op),
					"filename":  fmt.Sprint(event.FileInfo.Name()),
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Cycle trough every folder and add it and subfolder to watcher
	for _, rootFolder := range docRootArray {
		if err = w.AddRecursive(rootFolder); err != nil {
			return err
		}
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	// Start the watching process - it'll check for changes every 1000ms.
	go w.Start(time.Millisecond * 1000)

	return nil
}
