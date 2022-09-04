package gonetable_test

import (
	"testing"

	"github.com/juranki/gonetable"
)

func TestNewInvalidInput(t *testing.T) {
	_, err := gonetable.New(&gonetable.Schema{
		Tablename: "a",
	})
	if err != gonetable.ErrShortName {
		t.Fatal("expected too short")
	}
	_, err = gonetable.New(&gonetable.Schema{
		Tablename: "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
	})
	if err != gonetable.ErrLongName {
		t.Fatal("expected too long")
	}
	_, err = gonetable.New(&gonetable.Schema{
		Tablename: "asdf(",
	})
	if err != gonetable.ErrInvalidChar {
		t.Fatal("expected invalid character")
	}
	_, err = gonetable.New(&gonetable.Schema{
		Tablename: "asdf",
	})
	if err != gonetable.ErrNoRecTypes {
		t.Fatal("expected norecordtypeerror")
	}
	_, err = gonetable.New(&gonetable.Schema{
		Tablename:   "asdf",
		RecordTypes: map[string]gonetable.Document{},
	})
	if err != gonetable.ErrNoRecTypes {
		t.Fatal("expected norecordtypeerror")
	}
}
