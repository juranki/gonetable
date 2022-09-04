package gonetable

import (
	"errors"
	"regexp"
)

var (
	ErrNoRecTypes  = errors.New("at least one record type required")
	ErrShortName   = errors.New("table name too short (min: 3)")
	ErrLongName    = errors.New("table name too long (max: 255)")
	ErrInvalidChar = errors.New("invalid character in tablename (required pattern: [a-zA-Z0-9_.-]+)")

	re = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

type Schema struct {
	Tablename   string
	RecordTypes map[string]Document
}
type Table struct {
	schema *Schema
}

// New table
func New(schema *Schema) (*Table, error) {
	if nameError := checkTablename(schema.Tablename); nameError != nil {
		return nil, nameError
	}
	if len(schema.RecordTypes) == 0 {
		return nil, ErrNoRecTypes
	}
	return &Table{
		schema: schema,
	}, nil
}

func checkTablename(name string) error {
	// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html#API_CreateTable_RequestSyntax
	if len(name) < 3 {
		return ErrShortName
	}
	if len(name) > 255 {
		return ErrLongName
	}
	if !re.MatchString(name) {
		return ErrInvalidChar
	}
	return nil
}
