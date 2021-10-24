package utils

import (
	"fmt"
	"time"
)

// ConvertLabelToTime take a string from Excel sheet name and convert to a valid time object
func ConvertLabelToTime(label string) (time.Time, error) {
	var (
		parsedTime time.Time
		err        error
	)
	const timeOnly = "15.4"

	parsedTime, err = time.Parse(timeOnly, label)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func ConvertTimestampToTime(t time.Time) (time.Time, error) {
	var (
		h int //Hours
		m int //Minutes
	)

	h, m, _ = t.Clock()
	return ConvertLabelToTime(fmt.Sprintf("%v.%v", h, m))
}
