package regex_test

import (
	"log"
	"testing"

	"github.com/Aabdelm/employee-app/regex"
)

func TestRegexDigits(t *testing.T) {
	query := "19487"
	compiles := regex.IsDigit(query, log.Default())
	if !compiles {
		t.Fatal("Expect to get true")
	}
}

func TestRegexEmpty(t *testing.T) {
	query := ""
	compiles := regex.IsDigit(query, log.Default())
	if compiles {
		t.Fatalf("Expected to get false")
	}
}

func TestRegexDigitsAndText(t *testing.T) {
	query := "He110"
	compiles := regex.IsDigit(query, log.Default())

	if compiles {
		t.Fatalf("Expected to get false")
	}
}

func TestRegexText(t *testing.T) {
	query := "word"
	compiles := regex.IsDigit(query, log.Default())
	if compiles {
		t.Fatalf("Expected to get false")
	}
}
