package gonetable

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrNoDocSamples    = errors.New("at least one document sample required")
	ErrDuplicateTypeID = errors.New("multiple document samples with same type id")
	ErrIndexName       = errors.New("invalid index name, must match ^[a-zA-Z0-9_.-]{3,255}$")

	keyMethodRE = regexp.MustCompile(`^Gonetable_([a-zA-Z0-9]+)Key$`)
	indexRE     = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,255}$`)
)

type Schema struct {
	docTypes map[string]reflect.Type
	indeces  []string
}

func NewSchema(docSamples []Document) (*Schema, error) {
	if len(docSamples) == 0 {
		return nil, ErrNoDocSamples
	}
	s := Schema{
		docTypes: map[string]reflect.Type{},
		indeces:  []string{},
	}
	for _, d := range docSamples {
		docType := reflect.TypeOf(d)
		docTypeID := d.Gonetable_TypeID()
		if _, exists := s.docTypes[docTypeID]; exists {
			return nil, ErrDuplicateTypeID
		}
		s.docTypes[docTypeID] = docType
		s.indeces = append(s.indeces, getIndexNames(docType)...)
	}
	uniqueIndeces := map[string]bool{}
	for _, idx := range s.indeces {
		if !indexRE.MatchString(idx) {
			return nil, ErrIndexName
		}
		if _, exists := uniqueIndeces[idx]; !exists {
			uniqueIndeces[idx] = true
		}
	}
	s.indeces = make([]string, len(uniqueIndeces))
	i := 0
	for idx := range uniqueIndeces {
		s.indeces[i] = idx
		i++
	}
	return &s, nil
}

// Returns attribute definitions for all partition and sort keys fields
// of the table and GSIs
func (s *Schema) AttributeDefinitions() []types.AttributeDefinition {
	rv := makeIndexAttributes("PK", "SK")
	for _, idx := range s.indeces {
		rv = append(rv, makeIndexAttributes(
			fmt.Sprintf("%sPK", idx),
			fmt.Sprintf("%sSK", idx),
		)...)
	}
	return rv
}

// Returns definitions for GSIs
func (s *Schema) GlobalSecondaryIndexes() []types.GlobalSecondaryIndex {
	rv := []types.GlobalSecondaryIndex{}
	for _, idx := range s.indeces {
		rv = append(rv, types.GlobalSecondaryIndex{
			IndexName: aws.String(idx),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String(fmt.Sprintf("%sPK", idx)),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String(fmt.Sprintf("%sSK", idx)),
					KeyType:       types.KeyTypeRange,
				},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
		})
	}
	return rv
}

// Returns key schema that is always the same.
// Hash and range keys named PK and SK.
func (s *Schema) KeySchema() []types.KeySchemaElement {
	return []types.KeySchemaElement{
		{
			AttributeName: aws.String("PK"),
			KeyType:       types.KeyTypeHash,
		},
		{
			AttributeName: aws.String("SK"),
			KeyType:       types.KeyTypeRange,
		},
	}
}

func getIndexNames(documentType reflect.Type) []string {
	indeces := []string{}
	for i := 0; i < documentType.NumMethod(); i++ {
		name := documentType.Method(i).Name
		if matches := keyMethodRE.FindStringSubmatch(name); matches != nil {
			indeces = append(indeces, matches[1])
		}
	}
	return indeces
}

func makeIndexAttributes(pk, sk string) []types.AttributeDefinition {
	return []types.AttributeDefinition{
		{
			AttributeName: aws.String(pk),
			AttributeType: types.ScalarAttributeTypeS,
		},
		{
			AttributeName: aws.String(sk),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}
}
