package gonetable_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/juranki/gonetable"
)

type SimpleRecord struct {
	Name string
}

func (sr *SimpleRecord) GoneTable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"sr", "a"},
		RangeSegments: []string{"sr"},
	}
}
func (sr *SimpleRecord) GoneTable_Prefix() string { return "SR" }

func TestTable_GetCreateTableInputHasTablename(t *testing.T) {
	table, err := gonetable.New(&gonetable.Schema{
		Tablename: "tablename",
		RecordTypes: map[string]gonetable.RecordType{
			"SR": &SimpleRecord{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	createTableInput := table.GetCreateTableInput()
	if *createTableInput.TableName != "tablename" {
		t.Fatalf("expected tablename, got %s", *createTableInput.TableName)
	}
	json.NewEncoder(os.Stdout).Encode(createTableInput)
}

func TestTable_GetCreateTableInputHasAttributes(t *testing.T) {
	table, err := gonetable.New(&gonetable.Schema{
		Tablename: "tablename",
		RecordTypes: map[string]gonetable.RecordType{
			"SR": &SimpleRecord{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	createTableInput := table.GetCreateTableInput()
	requiredAttributes := []string{"PK", "SK"}
	for _, attr := range requiredAttributes {
		found := false
		for _, k := range createTableInput.AttributeDefinitions {
			if *k.AttributeName == attr {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s not found in attributes", attr)
		}
	}
}
