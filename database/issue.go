package database

import (
	"database/sql"
	"errors"
	"fmt"
	"internal-backend/websocket"
	"time"
)

type Issue struct {
	Id        string        `json:"id,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Operator  string        `json:"operator,omitempty"`
	Priority  int           `json:"priority,omitempty"`
	Title     string        `json:"title,omitempty"`
	Note      string        `json:"note,omitempty"`
	Open      bool          `json:"open,omitempty"`
	Detail    []IssueDetail `json:"detail,omitempty"`
}

type IssueDetail struct {
	Id        string    `json:"id,omitempty"`
	IssueID   string    `json:"issue_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Operator  string    `json:"operator,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// -------------------------
// Issue methods
// -------------------------

// GetAll get all open issue
// {mode}: Optional, if set to full return also all issue details
func (i Issue) GetAll(mode string, dest *[]Issue) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	// Query to retrieve all issues
	sqlStatement = `select id,timestamp,operator,priority,title,note
					from issue
					where open = true
					order by priority desc, timestamp desc;`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError retrieving Issue:\t%v", err))
	}

	defer rows.Close()

	for rows.Next() {
		var i Issue
		err = rows.Scan(&i.Id, &i.Timestamp, &i.Operator, &i.Priority, &i.Title, &i.Note)
		if err != nil {
			return errors.New(fmt.Sprintf("[ERR]\tError scanning row:\t%v", err))
		}

		// Check mode flag, if 'full' also retrieve details
		if mode == "full" {
			err := IssueDetail{}.GetDetailsByIssue(i.Id, &i.Detail)
			if err != nil {
				return errors.New(fmt.Sprintf("[ERR]\tError retrieving issue detail for issue {%s}:\t%v", i.Id, err))
			}
		}

		*dest = append(*dest, i)
	}

	return nil
}

// PostIssue insert new issue in db
func (i *Issue) PostIssue() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `
		INSERT INTO issue (operator, priority,title, note)
		VALUES ($1,$2,$3,$4)
		RETURNING id, timestamp
`
	err = DbConnection.QueryRow(sqlStatement, i.Operator, i.Priority, i.Title, i.Note).Scan(&i.Id, &i.Timestamp)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError inserting issue in db:\t%v", err))
	}

	websocket.Broadcast <- map[string]interface{}{
		"id":        websocket.Issue,
		"operation": "Issue created",
		"data":      i,
	}

	return nil
}

// CloseIssue close an issue setting open state to false
func (i *Issue) CloseIssue() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `
		UPDATE issue
		SET open = false
		WHERE id = $1
`
	_, err = DbConnection.Exec(sqlStatement, i.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError closing issue in db:\t%v", err))
	}

	websocket.Broadcast <- map[string]interface{}{
		"id":        websocket.Issue,
		"operation": "Issue closed",
		"data":      i,
	}

	return nil
}

// -------------------------
// IssueDetail methods
// -------------------------

// GetDetailsByIssue retrieve detail for passed issue id
func (i IssueDetail) GetDetailsByIssue(issueId string, dest *[]IssueDetail) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `
					select id,timestamp,operator,note
					from issue_detail
					where issue_id = $1
					order by timestamp asc;
`
	rows, err = DbConnection.Query(sqlStatement, issueId)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError retrieving issue detail list:\t%v", err))
	}

	defer rows.Close()

	for rows.Next() {
		var i IssueDetail
		err = rows.Scan(&i.Id, &i.Timestamp, &i.Operator, &i.Note)
		if err != nil {
			errors.New(fmt.Sprintf("[ERR]\tError scanning issue detail row:\t%v", err))
		}
		i.IssueID = issueId
		*dest = append(*dest, i)
	}

	return nil
}

// PostIssueDetail insert new issue detail to db
func (i *IssueDetail) PostIssueDetail(issueId string) error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `
		INSERT INTO issue_detail (issue_id, operator, note) 
		VALUES ($1,$2,$3)
		RETURNING id, timestamp
`

	err = DbConnection.QueryRow(sqlStatement, issueId, i.Operator, i.Note).Scan(&i.Id, &i.Timestamp)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError inserting issue detail in db:\t%v", err))
	}

	websocket.Broadcast <- map[string]interface{}{
		"id":        websocket.Issue,
		"operation": "Detail created",
		"data":      i,
	}

	return nil
}
