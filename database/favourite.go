package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Favourite struct {
	Id        string    `json:"id,omitempty"`         //Favourite UUID
	Timestamp time.Time `json:"timestamp,omitempty"`  //Favourite timestamp
	ConsoleIp string    `json:"console_ip,omitempty"` //User console IP
	Filename  string    `json:"filename,omitempty"`   //Favourite filename
	Count     int       `json:"count,omitempty"`      //Favourite count
}

// GetAll Get all favourites
func (f Favourite) GetAll(dest *[]Favourite) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id,timestamp,console_ip,filename from favourite`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving favourites: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var f Favourite
		err = rows.Scan(&f.Id, &f.Timestamp, &f.ConsoleIp, &f.Filename)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, f)
	}

	return nil
}

// GetAggregatedByConsoleIp Get favourite by console IP, aggregated by filename
func (f Favourite) GetAggregatedByConsoleIp(ip string, dest *[]Favourite) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select filename, count(filename) as count
					from favourite
					where console_ip = $1
					group by filename
					order by count desc
					limit 10
					`

	rows, err = DbConnection.Query(sqlStatement, ip)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving aggregated favourites by ip: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var f Favourite
		err = rows.Scan(&f.Filename, &f.Count)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, f)
	}

	return nil
}

// PostFavourite append new favourite entry to db
func (f Favourite) PostFavourite() error {
	var (
		err              error  // Error variable
		sqlStatement     string // SQL statement
		strippedFilename string // Filename without "rev" suffix
	)

	sqlStatement = `INSERT INTO favourite(console_ip, filename) 
					VALUES($1, $2)`

	// Cut filename at "rev" and take the first part, then trim spaces
	strippedFilename, _, _ = iCut(f.Filename, "rev")
	strippedFilename = strings.TrimSpace(strippedFilename)

	_, err = DbConnection.Exec(sqlStatement, f.ConsoleIp, strippedFilename)
	if err != nil {
		return errors.New(fmt.Sprintf("Error inserting favourite: %v\n", err))
	}

	// TODO: Evaluate if ws broadcast should be done here or in handler

	return nil
}

// TruncateTable Truncate favourite table
func (f Favourite) TruncateTable() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `truncate table favourite`

	_, err = DbConnection.Exec(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error truncating favourite table: %v\n", err))
	}

	return nil
}
