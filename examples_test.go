package gonetable_test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/juranki/gonetable"
)

const TABLENAME = "SchemaExample"

type ExampleDocument struct {
	ID   string
	Name string
}

func (ed *ExampleDocument) Gonetable_TypeID() string {
	return "ed"
}

func (ed *ExampleDocument) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"ed", ed.ID},
		RangeSegments: []string{"ed"},
	}
}

func ExampleSchema() {
	cfg := MustLoadLocalDDBConfig()
	client := dynamodb.NewFromConfig(cfg)
	DeleteTableIfExists(context.Background(), client, TABLENAME)

	schema, err := gonetable.NewSchema([]gonetable.Document{
		&ExampleDocument{},
	})
	if err != nil {
		panic(err)
	}

	_, err = client.CreateTable(
		context.Background(),
		&dynamodb.CreateTableInput{
			TableName:              aws.String(TABLENAME),
			BillingMode:            types.BillingModePayPerRequest,
			AttributeDefinitions:   schema.AttributeDefinitions(),
			KeySchema:              schema.KeySchema(),
			GlobalSecondaryIndexes: schema.GlobalSecondaryIndexes(),
		},
	)
	if err != nil {
		panic(err)
	}

	ed := &ExampleDocument{
		ID:   "123456",
		Name: "Example",
	}

	marshaled, err := schema.Marshal(ed)
	if err != nil {
		panic(err)
	}
	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(TABLENAME),
		Item:      marshaled,
	})
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(marshaled)
	// Output:
	// {"ID":{"Value":"123456"},"Name":{"Value":"Example"},"PK":{"Value":"ed#123456"},"SK":{"Value":"ed"},"_Type":{"Value":"ed"}}
}
