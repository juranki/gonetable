package gonetable_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/juranki/gonetable"
)

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

	schema, err := gonetable.NewSchema([]gonetable.Document{
		&ExampleDocument{},
	})
	if err != nil {
		panic(err)
	}

	table, err := client.CreateTable(
		context.Background(),
		&dynamodb.CreateTableInput{
			TableName:              aws.String("ExampleTable"),
			BillingMode:            types.BillingModePayPerRequest,
			AttributeDefinitions:   schema.AttributeDefinitions(),
			KeySchema:              schema.KeySchema(),
			GlobalSecondaryIndexes: schema.GlobalSecondaryIndexes(),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(*table.TableDescription.TableName)
	// Output:
	// ExampleTable
}
