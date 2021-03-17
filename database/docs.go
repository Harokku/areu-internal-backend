package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Document struct {
	Id          string `json:"id"`           //Document UUID
	Hash        string `json:"hash"`         //Document hash
	FileName    string `json:"file_name"`    //Document filename (full path)
	DisplayName string `json:"display_name"` //Document displayed name
	Category    string `json:"category"`     //Document category (based on path)
}

// Get all documents
func (d Document) GetAll(dest *[]Document) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id,hash,filename,displayname,category from docs`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving documents: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var d Document
		err = rows.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, d)
	}

	return nil
}

// Get document by hash
func (d *Document) GetByHash(hash string) error {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
	)

	sqlStatement = `select id,hash,filename,displayname,category from docs where hash=$1`

	row = DbConnection.QueryRow(sqlStatement, hash)
	switch err = row.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category); err {
	case sql.ErrNoRows:
		return errors.New("no row where retrieved")
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("error retrieving doc from db: %v\n", err))
	}
}

// Get document by id
func (d *Document) GetById(id string) error {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
	)

	sqlStatement = `select id,hash,filename,displayname,category from docs where id=$1`

	row = DbConnection.QueryRow(sqlStatement, id)
	switch err = row.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category); err {
	case sql.ErrNoRows:
		return errors.New("no row where retrieved")
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("error retrieving doc from db: %v\n", err))
	}
}

// Build hash table in form of hash:path
func (d Document) GetHashTable() (map[string]string, error) {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
		hashTable    map[string]string
	)

	sqlStatement = `select hash,filename from docs`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving hash/filename: %v\n", err))
	}

	defer rows.Close()

	hashTable = make(map[string]string)

	for rows.Next() {
		var d Document
		err = rows.Scan(&d.Hash, &d.FileName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning row %v\n", err))
		}
		hashTable[d.Hash] = d.FileName
	}

	return hashTable, nil
}
