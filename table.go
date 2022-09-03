package gonetable

import (
	"errors"
	"regexp"
)

var (
	NoRecTypesError  = errors.New("at least one record type required")
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
	if len(types) == 0 {
		return nil, NoRecTypesError
	}
	return &Table{
		name:  name,
		types: types,
	}, nil
}

func checkTablename(name string) error {
	// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html#API_CreateTable_RequestSyntax
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
