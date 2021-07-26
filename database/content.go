package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type Content struct {
	Id          string `json:"id"`           //Content UUID
	DisplayName string `json:"display_name"` //Content display name as should appear in frontend
	Link        string `json:"link"`         //Content link to construct resource URI
	SheetNumber int    `json:"sheet_number"` //Content XLSX sheet number to read from
}

// GetAll Get all content links
func (c Content) GetAll(dest *[]Content) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id, display_name, link, sheet_number from content_links`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrievinf content links: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var c Content
		err = rows.Scan(&c.Id, &c.DisplayName, &c.Link, &c.SheetNumber)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, c)
	}

	return nil
}

// GetSheetNumber return sheet number by link
func (c *Content) GetSheetNumber(link string) error {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
	)

	sqlStatement = `select sheet_number from content_links where link=$1 limit 1`

	row = DbConnection.QueryRow(sqlStatement, link)
	c.Link = link
	switch err = row.Scan(&c.SheetNumber); err {
	case sql.ErrNoRows:
		return errors.New("no row where retrieved")
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("error retrieving link from db: %v\n", err))
	}
}

// TruncateTable Truncate (clean) actual table
func (c Content) TruncateTable() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `TRUNCATE TABLE content_links`

	_, err = DbConnection.Exec(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error truncating Content table"))
	}

	return nil
}

// BulkCreate Bulk create passed in content array
func (c Content) BulkCreate(contentToAdd []Content) error {
	var (
		err          error
		sqlStatement *sql.Stmt //Prepared sql statement
		txn          *sql.Tx   //DB transaction
	)

	//Begin new transaction
	txn, err = DbConnection.Begin()
	if err != nil {
		return err
	}

	//Prepare insert statement
	sqlStatement, err = txn.Prepare(pq.CopyIn("content_links", "display_name", "link", "sheet_number"))

	//Exec insert for every passed content
	for _, content := range contentToAdd {
		_, err = sqlStatement.Exec(content.DisplayName, content.Link, content.SheetNumber)
		if err != nil {
			return err
		}
	}

	//Flush actual data
	_, err = sqlStatement.Exec()
	if err != nil {
		return err
	}

	err = sqlStatement.Close()
	if err != nil {
		return err
	}

	//Execute transaction
	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
