package handler

import "fmt"

const (
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

func ErrString(m string) string {
	return fmt.Sprintf("%s[ERR] %s%s", colorRed, m, colorReset)
}

func ErrStringMsg(m string, e error) string {
	return fmt.Sprintf("%s[ERR] %s:%s %s", colorRed, m, colorReset, e)
}
