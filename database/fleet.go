package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"internal-backend/utils"
	"time"
)

type Fleet struct {
	Id            string    `json:"id"`            // Vehicle UUID
	Ente          string    `json:"ente"`          // Vehicle callsign
	Stazionamento string    `json:"stazionamento"` // Vechicle position
	Convenzione   string    `json:"convenzione"`   // Vehicle convention type
	Minimum       string    `json:"minimum"`       // Minimum number of personnel on board
	ActiveFrom    time.Time `json:"active_from"`   // Time interval to check for availability
}

type BacoSnapshoot struct {
	Id            string `json:"id"`            // Vehicle UUID
	Ente          string `json:"ente"`          // Vehicle callsign
	Mezzo         string `json:"mezzo"`         // Vehicle name
	Stazionamento string `json:"stazionamento"` // Vechicle position
	Convenzione   string `json:"convenzione"`   // Vehicle convention type
	Radio         string `json:"radio"`         // Vechicle radio id

}

// -------------------------
// Main db function
// -------------------------

// GetAll Get all fleet data
func (c Fleet) GetAll(dest *[]Fleet) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id,convenzione,ente,minimum,active_from from check_convenzioni order by convenzione desc, ente asc`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving fleet info: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var c Fleet
		err = rows.Scan(&c.Id, &c.Convenzione, &c.Ente, &c.Minimum, &c.ActiveFrom)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, c)
	}

	return nil
}

func (c Fleet) GetActiveNow(dest *[]Fleet) error {
	var (
		err                           error
		rows                          *sql.Rows
		actualRange                   time.Time //Actual time range to retrieve from db built based on now
		sqlStatement                  string
		sqlActualRangeStatement       string // Query to get actual time range based on now
		sqlActualRangeStatementIfNull string // Query to get actual range if precedent is null
	)

	sqlActualRangeStatement = `	select active_from
								from check_convenzioni
								where active_from < $1
								order by active_from desc
								limit 1`

	sqlActualRangeStatementIfNull = `	select active_from
										from check_convenzioni
										where active_from > $1
										order by active_from desc
										limit 1`

	sqlStatement = `select id,convenzione,ente,stazionamento,minimum,active_from from check_convenzioni where active_from=$1 order by convenzione desc, ente asc`

	// Look for actual time range
	nowTime, err := utils.ConvertTimestampToTime(time.Now())
	if err != nil {
		return err
	}
	row := DbConnection.QueryRow(sqlActualRangeStatement, nowTime)
	switch err = row.Scan(&actualRange); err {
	case sql.ErrNoRows:
		row = DbConnection.QueryRow(sqlActualRangeStatementIfNull, nowTime)
		switch err = row.Scan(&actualRange); err {
		case sql.ErrNoRows:
			return errors.New("no row where retrieved")
		case nil:
		default:
			return errors.New(fmt.Sprintf("error retrieving actual time range from db: %v\n", err))
		}
	case nil:
	default:
		return errors.New(fmt.Sprintf("error retrieving actual time range from db: %v\n", err))
	}

	// Retrieve actual active vehicles
	rows, err = DbConnection.Query(sqlStatement, actualRange)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving fleet info: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var c Fleet
		err = rows.Scan(&c.Id, &c.Convenzione, &c.Ente, &c.Stazionamento, &c.Minimum, &c.ActiveFrom)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, c)
	}

	return nil
}

// TruncateTable Truncate (clean) actual table
func (c Fleet) TruncateTable() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `TRUNCATE TABLE check_convenzioni`

	_, err = DbConnection.Exec(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error truncating Fleet table"))
	}

	return nil
}

// BulkCreate Bulk create passed in content array
func (c Fleet) BulkCreate(contentToAdd []Fleet) error {
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
	sqlStatement, err = txn.Prepare(pq.CopyIn("check_convenzioni", "convenzione", "ente", "active_from", "stazionamento", "minimum"))

	//Exec insert for every passed content
	for _, content := range contentToAdd {
		_, err = sqlStatement.Exec(content.Convenzione, content.Ente, content.ActiveFrom, content.Stazionamento, content.Minimum)
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

// -------------------------
// Baco db function
// -------------------------

// GetSnapshoot return actual fleet state reading last system snapshoot
func (b BacoSnapshoot) GetSnapshoot(dest *[]BacoSnapshoot) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id, ente,mezzo,stazionamento,convenzione,radio from "db_Baco" order by convenzione desc , ente asc `

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving BaCo Db snapshot: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var c BacoSnapshoot
		err = rows.Scan(&c.Id, &c.Ente, &c.Mezzo, &c.Stazionamento, &c.Convenzione, &c.Radio)
		if err != nil {
			return errors.New(fmt.Sprintf("Error retrieving row: %v\n", err))
		}
		*dest = append(*dest, c)
	}

	return nil
}

// TruncateBaCoTable Truncate (clean) actual table
func (b BacoSnapshoot) TruncateBaCoTable() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `TRUNCATE TABLE "db_Baco"`

	_, err = DbConnection.Exec(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error truncating BaCo Snapshoot table"))
	}

	return nil
}

// BulkCreateBaCo Bulk create passed in content array
func (b BacoSnapshoot) BulkCreateBaCo(contentToAdd []BacoSnapshoot) error {
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
	sqlStatement, err = txn.Prepare(pq.CopyIn("db_Baco", "ente", "mezzo", "stazionamento", "convenzione", "radio"))

	//Exec insert for every passed content
	for _, content := range contentToAdd {
		_, err = sqlStatement.Exec(content.Ente, content.Mezzo, content.Stazionamento, content.Convenzione, content.Radio)
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
