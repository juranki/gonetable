package gonetable_test

import (
	"testing"

	"github.com/juranki/gonetable"
)

func TestNewTableName(t *testing.T) {
	_, err := gonetable.New("", []gonetable.RecordType{})
	if err != gonetable.ShortNameError {
		t.Fatal("expected too short")
	}
	_, err = gonetable.New("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789", []gonetable.RecordType{})
	if err != gonetable.LongNameError {
		t.Fatal("expected too long")
	}
	_, err = gonetable.New("asdf(", []gonetable.RecordType{})
	if err != gonetable.InvalidCharError {
		t.Fatal("expected invalid character")
	}

}
