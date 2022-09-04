package gonetable

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	baseAttributes = []types.AttributeDefinition{
		{
			AttributeName: aws.String("PK"),
			AttributeType: types.ScalarAttributeTypeS,
		},
		{
			AttributeName: aws.String("SK"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}
)

func (table *Table) GetCreateTableInput() *dynamodb.CreateTableInput {
	input := &dynamodb.CreateTableInput{
		TableName:            aws.String(table.schema.Tablename),
		BillingMode:          types.BillingModePayPerRequest,
		TableClass:           types.TableClassStandard,
		AttributeDefinitions: []types.AttributeDefinition{},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{},
		SSESpecification: &types.SSESpecification{
			Enabled: aws.Bool(true),
		},
		StreamSpecification: &types.StreamSpecification{},
		Tags:                []types.Tag{},
	}
	input.AttributeDefinitions = append(input.AttributeDefinitions, baseAttributes...)
	table.reflectAttributes()
	return input
}

func (table *Table) reflectAttributes() {
	for _, recordType := range table.schema.RecordTypes {
		// s := reflect.ValueOf(recordType).Elem()
		typeOfRecord := reflect.TypeOf(recordType)
		// for i := 0; i < s.NumField(); i++ {
		// 	f := s.Field(i)
		// 	fmt.Printf("%d: %s %s\n", i, typeOfRecord.Field(i).Name, f.Type())
		// }
		for i := 0; i < typeOfRecord.NumMethod(); i++ {
			m := typeOfRecord.Method(i)
			fmt.Printf("Method %d: %s\n", i, m.Name)
		}
	}
}
