package handler

import (
	"fmt"
	"strings"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	//colorGreen  = "\033[32m"
	//colorYellow = "\033[33m"
	//colorBlue   = "\033[34m"
	//colorPurple = "\033[35m"
	//colorCyan   = "\033[36m"
	//colorWhite  = "\033[37m"
)

func ErrString(m string) string {
	return fmt.Sprintf("%s[ERR] %s%s", colorRed, m, colorReset)
}

func ErrStringMsg(m string, e error) string {
	return fmt.Sprintf("%s[ERR] %s:%s %s", colorRed, m, colorReset, e)
}

// ExtractLotto Take a fleet entry and return vehicle callsign and lotto
func ExtractLotto(item string) (string, string) {
	var (
		callsign string
		lotto    string
		found    bool
	)
	callsign, lotto, found = strings.Cut(strings.TrimSpace(item), ".")
	if found {
		return callsign, lotto
	}
	return "", ""
}

// JoinCallsignLotto return joined vehicle id
func JoinCallsignLotto(callsign, lotto string) string {
	return strings.Join([]string{strings.TrimSpace(callsign), strings.TrimSpace(lotto)}, ".")
}
