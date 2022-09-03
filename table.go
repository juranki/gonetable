package gonetable

import (
	"errors"
	"regexp"
)

var (
	ShortNameError   = errors.New("table name too short (min: 3)")
	LongNameError    = errors.New("table name too long (max: 255)")
	InvalidCharError = errors.New("invalid character in tablename (required pattern: [a-zA-Z0-9_.-]+)")

	KeyDelimiter = "#"

	re = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

type RecordType interface {
	KeyPrefix() string
}
type Table struct {
	name  string
	types []RecordType
}

// New table
func New(name string, types []RecordType) (*Table, error) {
	if nameError := checkTablename(name); nameError != nil {
		return nil, nameError
	}
	return &Table{
		name:  name,
		types: types,
	}, nil
}

// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html#API_CreateTable_RequestSyntax
func checkTablename(name string) error {
	if len(name) < 3 {
		return ShortNameError
	}
	if len(name) > 255 {
		return LongNameError
	}
	if !re.MatchString(name) {
		return InvalidCharError
	}
	return nil
}
