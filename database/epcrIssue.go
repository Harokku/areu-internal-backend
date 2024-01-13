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
