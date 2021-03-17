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
