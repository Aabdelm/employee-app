package regex

import (
	"log"
	"regexp"
)

func IsDigit(digitString string, l *log.Logger) bool {
	l.Printf("[INFO] Function isDigit started")
	regex, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		l.Printf("[INFO] Failed to compile regex. Error %s", err)
		return false
	}
	return regex.Match([]byte(digitString))
}
