package crawler

import (
	"errors"
	"fmt"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	hashTable map[string]string //database hash -> hash:path
)

// Walk through filesystem and enumerate found files in db, hashing them
func EnumerateDocuments() error {
	var (
		err     error
		docRoot string //Documents disk path
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

	hashTable, err = database.Document{}.GetHashTable()
	if err != nil {
		return errors.New(fmt.Sprintf("crawler/EnumerateDocuments returned error while retrieving hash table from db: %v\n", err))
	}

	//TODO: Implement filewalker to enum existent files
	if err = filepath.Walk(docRoot, addDir); err != nil {
		return err
	}

	enumerateDuration := time.Since(enumerateStartTime) //calculate total startup time
	log.Printf("Document enumerated in %s", enumerateDuration)

	return nil
}

//TODO: Implement function to add file to db after hashing it
func addDir(path string, fi os.FileInfo, err error) error {
	var (
		category     string //Calculated category based on path
		displayName  string //Calculated display name based on filename
		sha1Checksum string //SHA-1 filename hash
	)
	//TODO: Remove printf
	fmt.Printf("Path: %v\n", path)

	// Check if file is not dir and process it
	if !fi.IsDir() {

		//Call helper func to extract category from path as relative path from document root
		category, err = getCategory(path)
		if err != nil {
			log.Printf("error retrieving category: %v\n", err)
		}
		//TODO: Remove printf
		fmt.Printf("category: %v\n", category)

		//Call helper func to extract display name from path
		displayName = getDisplayNameFromPath(path)
		if err != nil {
			log.Printf("error retrieving display name: %v\n", err)
		}
		//TODO: Remove printf
		fmt.Printf("displayName: %v\n", displayName)

		//Call helper func to calculate SHA-1 hash from file
		sha1Checksum, err = getSha1(path)
		if err != nil {
			log.Printf("error calculating SHA-1 from file")
		}
		//TODO: Remove printf
		fmt.Printf("SHA-1: %v\n", sha1Checksum)
	}

	return nil
}
