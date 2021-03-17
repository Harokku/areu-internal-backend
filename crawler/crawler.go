package crawler

import (
	"internal-backend/database"
	"log"
	"os"
)

var (
// dbMap map[string]string //database hash -> hash:path
)

// Walk through filesystem and enumerate found files in db, hashing them
func EnumerateDocuments() error {
	var (
		err error
		d   []database.Document //Documents retrieved from db
	)

	// -------------------------
	// Dummy db interaction test
	// -------------------------
	dbDocument := database.Document{} //Db interaction object
	err = dbDocument.GetAll(&d)
	if err != nil {
		log.Printf("Error in crawler/EnumerateDocuments: %v\n", err)
		return err
	}
	log.Printf("Retrieved documents: %v\n", d)
	return nil

	//populate dbMap with actual db state

}

func addDir(path string, fi os.FileInfo, err error) error {
	return nil
}
