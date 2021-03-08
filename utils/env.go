package utils

import (
	"errors"
	"os"
)

//ReadEnv read s env variable and return it or error if failed
func ReadEnv(s string) (string, error) {
	r := os.Getenv(s)
	if r != "" {
		return r, nil
	}
	return "", errors.New("passed env variable not set")
}
