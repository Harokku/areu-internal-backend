package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Issue struct {
	Id        string        `json:"id,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Operator  string        `json:"operator,omitempty"`
	Priority  int           `json:"priority,omitempty"`
	Note      string        `json:"note,omitempty"`
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

// GetAll get all issue
// {mode}: Optional, if set to full return also all issue details
func (i Issue) GetAll(mode string, dest *[]Issue) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	// Query to retrieve all issues
	sqlStatement = `select id,timestamp,operator,priority,note
					from issue
					order by priority desc, timestamp desc;`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError retrieving Issue:\t%v", err))
	}

	defer rows.Close()

	for rows.Next() {
		var i Issue
		err = rows.Scan(&i.Id, &i.Timestamp, &i.Operator, &i.Priority, &i.Note)
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
		INSERT INTO issue (operator, priority, note)
		VALUES ($1,$2,$3)
		RETURNING id, timestamp
`
	err = DbConnection.QueryRow(sqlStatement, i.Operator, i.Priority, i.Note).Scan(&i.Id, &i.Timestamp)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERR]\tError inserting issue in db:\t%v", err))
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
					order by timestamp desc;
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

	return nil
}
