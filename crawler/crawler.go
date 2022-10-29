package crawler

import (
	"errors"
	"fmt"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	hashTable      map[string]string   //database hash -> hash:path
	documentObject []database.Document //Db array to be bulk added
)

// EnumerateDocuments Walk through filesystem and enumerate found files in doc root, hashing them and building dictionary
func EnumerateDocuments() error {
	var (
		err          error
		docRoot      string   //Documents disk path from env
		docRootArray []string //Documents array of path to check
	)
	log.Printf("Starting enumerating documents...")
	enumerateStartTime := time.Now()

	// -------------------------
	// Read .env
	// -------------------------
	docRoot, err = utils.ReadEnv("DOC_ROOT")
	if err != nil {
		log.Fatalf("Error retrieving documents root from env")
	}
	docRootArray = strings.Split(docRoot, "|")

	// -------------------------
	// Walk documents hierarchy and process files
	// -------------------------

	// Clean existent pending array
	documentObject = []database.Document{}

	// Enumerate files and create document object to be added to db
	//
	// Expand input array and check all path
	for _, p := range docRootArray {
		if err = filepath.Walk(p, addFile(p)); err != nil {
			log.Printf("Error walking documents: %s", err)
			return filepath.SkipDir
		}
	}

	// Truncate document table for fresh start
	err = database.Document{}.TruncateTable()
	if err != nil {
		return errors.New(fmt.Sprintf("crawler/EnumerateDocuments returned error while truncating table: %v\n", err))
	}

	// Add hierarchy to db
	err = database.Document{}.BulkCreate(documentObject)
	if err != nil {
		return errors.New(fmt.Sprintf("crawler/EnumerateDocuments returned error while bulk creating table: %v\n", err))
	}

	enumerateDuration := time.Since(enumerateStartTime) //calculate total startup time
	log.Printf("Document enumerated in %s", enumerateDuration)

	return nil
}

// TODO: implement not hard coded version
// Check if filename is on the excluded list
func checkExcludedFile(filename string) bool {
	excludedFiles := map[string]bool{
		"Thumbs.db":         true,
		".DS_Store":         true,
		"sync_ffs.lock":     true,
		"Doc per RAR SOREU": true,
	}

	if excludedFiles[filename] {
		return true
	}
	return false
}

// If file isn't a dir process it, extracting: category, display name, path and SHA-1
func addFile(r string) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		var (
			category     string            //Calculated category based on path
			displayName  string            //Calculated display name based on filename
			sha1Checksum string            //SHA-1 filename hash
			newDoc       database.Document //New document entry to append
		)

		// nil point deference guard on filepath not accessible
		if fi == nil || path == "" {
			return nil
		}

		// Acces error handling
		if fi.IsDir() && fi.Name() == "Doc per RAR SOREU" {
			fmt.Printf("skipping a dir without errors: %+v \n", fi.Name())
			return nil
		}
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return nil
		}

		// Check if file is in avoidance list and return
		if checkExcludedFile(filepath.Base(path)) {
			return nil
		}

		// Check if file is not dir and process it
		if fi.IsDir() {
			//Call helper func to extract category from path as relative path from document root
			category, err = getCategory(r, path)
			if err != nil {
				log.Printf("error retrieving category: %v\n", err)
			}

			//Call helper func to extract display name from path
			displayName = getDisplayNameFromPath(path)
			if err != nil {
				log.Printf("error retrieving display name: %v\n", err)
			}

			//Add info to documentObject array
			newDoc.FileName = filepath.ToSlash(path)
			newDoc.DisplayName = displayName
			newDoc.Category = strings.ReplaceAll(category, "\\", "/")
			newDoc.IsDir = true
			newDoc.CreationTime = fi.ModTime()
			documentObject = append(documentObject, newDoc)
		} else {
			//Call helper func to extract category from path as relative path from document root
			category, err = getCategory(r, path)
			if err != nil {
				log.Printf("error retrieving category: %v\n", err)
			}

			//Call helper func to extract display name from path
			displayName = getDisplayNameFromPath(path)
			if err != nil {
				log.Printf("error retrieving display name: %v\n", err)
			}

			//Call helper func to calculate SHA-1 hash from file
			sha1Checksum, err = getSha1(path)
			if err != nil {
				log.Printf("error calculating SHA-1 from file")
			}

			//Add info to documentObject array
			newDoc.Hash = sha1Checksum
			newDoc.FileName = filepath.ToSlash(path)
			newDoc.DisplayName = displayName
			newDoc.Category = strings.ReplaceAll(category, "\\", "/")
			newDoc.CreationTime = fi.ModTime()
			documentObject = append(documentObject, newDoc)
		}

		return nil
	}
}
