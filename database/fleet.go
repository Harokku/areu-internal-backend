package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Fleet struct {
	Id         string    `json:"id"`          // Vehicle UUID
	ConvType   string    `json:"conv_type"`   // Vehicle convention type
	Name       string    `json:"name"`        // Vehicle callsign
	ActiveFrom time.Time `json:"active_from"` // Time interval to check for availability
}

// GetAll Get all fleet data
func (c Fleet) GetAll(dest *[]Fleet) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id,conv_type,name,active_from from check_convenzioni`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrievinf content links: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var c Fleet
		err = rows.Scan(&c.Id, &c.ConvType, &c.Name, &c.ActiveFrom)
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
		actualRange                   string //Actual time range to retrieve from db built based on now
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

	// Look for actual time range

	sqlStatement = `select id,conv_type,name,active_from from check_convenzioni`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrievinf content links: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var c Fleet
		err = rows.Scan(&c.Id, &c.ConvType, &c.Name, &c.ActiveFrom)
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
	sqlStatement, err = txn.Prepare(pq.CopyIn("check_convenzioni", "conv_type", "name", "active_from"))

	//Exec insert for every passed content
	for _, content := range contentToAdd {
		_, err = sqlStatement.Exec(content.ConvType, content.Name, content.ActiveFrom)
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
