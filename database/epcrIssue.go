package database

import (
	"fmt"
	"net"
	"time"
)

type EpcrIssue struct {
	Id        string    `json:"id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	IpAddress net.IP    `json:"ipAddress,omitempty"`
	VehicleId string    `json:"vehicleId,omitempty"`
	Text      string    `json:"text,omitempty"`
}

// TruncateTable truncates the epcr_issues table in the database.
// It accepts a *sql.DB parameter and performs the TRUNCATE TABLE operation.
// Returns an error if the operation fails.
func (e EpcrIssue) TruncateTable() error {
	_, err := DbConnection.Exec("TRUNCATE TABLE  epcr_issues")
	return err
}

// GetAll retrieves all ePCR issues from the database.
// It executes the SQL query "SELECT ip_address, vehicle_id, issue, timestamp FROM epcr_issues"
// and returns a slice of EpcrIssue objects and an error if the query fails.
// Each row returned from the query is scanned into an EpcrIssue object, and the objects are appended to the issues slice.
// The rows are closed before returning the result.
func (e EpcrIssue) GetAll() ([]EpcrIssue, error) {
	sqlStatement := `
		SELECT ip_address,vehicle_id,issue,timestamp from epcr_issues
`
	rows, err := DbConnection.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("[ERR]\tError querying ePCR issues from db:\t%v", err)
	}
	defer rows.Close()

	var issues []EpcrIssue
	for rows.Next() {
		var issue EpcrIssue
		ip := ""
		err := rows.Scan(&ip, &issue.VehicleId, &issue.Text, &issue.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("[ERR]\tError scanning ePCR issue from db:\t%v", err)
		}
		issue.IpAddress = net.ParseIP(ip)

		issues = append(issues, issue)
	}

	return issues, nil
}

// PostIssue posts an ePCR issue to the database.
func (e EpcrIssue) PostIssue() error {
	var (
		err          error
		sqlStatement string
	)
	sqlStatement = `
		INSERT INTO epcr_issues (ip_address, vehicle_id, issue)
		VALUES ($1,$2,$3)
`

	_, err = DbConnection.Exec(sqlStatement, e.IpAddress, e.VehicleId, e.Text)
	if err != nil {
		return fmt.Errorf("[ERR]\tError inserting ePCR issue in db:\t%v", err)
	}

	return nil
}
