package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type NewsFeed struct {
	Id        string `json:"id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Title     string `json:"title,omitempty"`
	News      string `json:"news,omitempty"`
}

func (nf NewsFeed) GetAll(dest *[]NewsFeed) error {

	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id, timestamp, title,news from news_feed`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving news: %v\n", err))
	}

	for rows.Next() {
		var nf NewsFeed
		err = rows.Scan(&nf.Id, &nf.Timestamp, &nf.Title, &nf.News)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, nf)
	}

	return nil
}

// TruncateTable Truncate (clean) actual table
func (nf NewsFeed) TruncateTable() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `TRUNCATE TABLE news_feed`

	_, err = DbConnection.Exec(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error truncating news feed table"))
	}

	return nil
}

func (nf NewsFeed) BulkCreate(contentToAdd []NewsFeed) error {
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
	sqlStatement, err = txn.Prepare(pq.CopyIn("news_feed", "timestamp", "title", "news"))

	//Exec insert for every passed content
	for _, content := range contentToAdd {
		_, err = sqlStatement.Exec(content.Timestamp, content.Title, content.News)
		if err != nil {
			txn.Rollback()
			return err
		}
	}

	// Flush actual data
	_, err = sqlStatement.Exec()
	if err != nil {
		txn.Rollback()
		return err
	}

	// Close statement before commit
	err = sqlStatement.Close()
	if err != nil {
		return err
	}

	// Commit the transaction
	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
