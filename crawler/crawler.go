package crawler

import (
	"errors"
	"fmt"
	"internal-backend/database"
	"internal-backend/utils"
	"log"
	"os"
	"path/filepath"
)

var (
	hashTable map[string]string //database hash -> hash:path
)

// Walk through filesystem and enumerate found files in db, hashing them
func EnumerateDocuments() error {
	var (
		err      error
		docsRoot string //Documents disk path
	)

	// -------------------------
	// Read .env
	// -------------------------
	docsRoot, err = utils.ReadEnv("DOC_ROOT")
	if err != nil {
		log.Fatalf("Error retrieving documents root from env")
	}

	hashTable, err = database.Document{}.GetHashTable()
	if err != nil {
		return errors.New(fmt.Sprintf("crawler/EnumerateDocuments returned error while retrieving hash table from db: %v\n", err))
	}

	//TODO: Implement filewalker to enum existent files
	if err = filepath.Walk(docsRoot, addDir); err != nil {
		return err
	}

	return nil
}

//TODO: Implement function to add file to db after hashing it
func addDir(path string, fi os.FileInfo, err error) error {
	fmt.Printf("path: %v\n", path)
	fmt.Printf("fi: %v\n", fi)
	return nil
}
