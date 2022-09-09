package gonetable_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/juranki/gonetable"
)

func TestNewSchema(t *testing.T) {
	type args struct {
		docSamples []gonetable.Document
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "invalid index",
			args: args{
				docSamples: []gonetable.Document{&InvalidIndex{}},
			},
			wantErr: gonetable.ErrIndexName,
		},
		{
			name: "no doc samples",
			args: args{
				docSamples: []gonetable.Document{},
			},
			wantErr: gonetable.ErrNoDocSamples,
		},
		{
			name: "duplicate type id",
			args: args{
				docSamples: []gonetable.Document{&InvalidIndex{}, &InvalidIndex{}},
			},
			wantErr: gonetable.ErrDuplicateTypeID,
		},
		{
			name: "simple doc",
			args: args{
				docSamples: []gonetable.Document{&MinimalDoc{}},
			},
			wantErr: nil,
		},
		{
			name: "one index",
			args: args{
				docSamples: []gonetable.Document{&WithIndex{}},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gonetable.NewSchema(tt.args.docSamples)
			if err == nil {
				if tt.wantErr != nil {
					t.Error("Expected error")
				}
				return
			}
			if tt.wantErr.Error() == err.Error() {
				return
			}
			t.Errorf("error = '%v', want '%v'", err.Error(), tt.wantErr.Error())
		})
	}
}

func TestSchema_AttributeDefinitions(t *testing.T) {
	tests := []struct {
		name       string
		docSamples []gonetable.Document
		want       []types.AttributeDefinition
	}{
		{
			name:       "minimal",
			docSamples: []gonetable.Document{&MinimalDoc{}},
			want: []types.AttributeDefinition{
				{
					AttributeName: aws.String("PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
		},
		{
			name:       "withIndex",
			docSamples: []gonetable.Document{&WithIndex{}},
			want: []types.AttributeDefinition{
				{
					AttributeName: aws.String("PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI1PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI1SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := gonetable.NewSchema(tt.docSamples)
			if err != nil {
				t.Fatal(err)
			}
			if got := s.AttributeDefinitions(); !reflect.DeepEqual(got, tt.want) {
				json.NewEncoder(os.Stdout).Encode(got)
				json.NewEncoder(os.Stdout).Encode(tt.want)
				t.Errorf("Schema.AttributeDefinitions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchema_KeySchema(t *testing.T) {
	tests := []struct {
		name       string
		docSamples []gonetable.Document
		want       []types.KeySchemaElement
	}{
		{
			name:       "minimal",
			docSamples: []gonetable.Document{&MinimalDoc{}},
			want: []types.KeySchemaElement{
				{
					AttributeName: aws.String("PK"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("SK"),
					KeyType:       types.KeyTypeRange,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := gonetable.NewSchema(tt.docSamples)
			if err != nil {
				t.Fatal(err)
			}
			if got := s.KeySchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Schema.KeySchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchema_GlobalSecondaryIndexes(t *testing.T) {
	tests := []struct {
		name       string
		docSamples []gonetable.Document
		want       []types.GlobalSecondaryIndex
	}{
		{
			name:       "minimal",
			docSamples: []gonetable.Document{&MinimalDoc{}},
			want:       nil,
		},
		{
			name:       "with index",
			docSamples: []gonetable.Document{&WithIndex{}},
			want: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("GSI1"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI1PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI1SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := gonetable.NewSchema(tt.docSamples)
			if err != nil {
				t.Fatal(err)
			}
			if got := s.GlobalSecondaryIndexes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Schema.GlobalSecondaryIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchema_Marshal(t *testing.T) {
	type args struct {
		docSamples []gonetable.Document
		doc        gonetable.Document
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]types.AttributeValue
		wantErr bool
	}{
		{
			name: "minimal",
			args: args{
				docSamples: []gonetable.Document{&MinimalDoc{}},
				doc: &MinimalDoc{
					Name: "hiihaa",
				},
			},
			want: map[string]types.AttributeValue{
				"Name":  MustMarshal("hiihaa"),
				"PK":    MustMarshal("a#b"),
				"SK":    MustMarshal("a#b"),
				"_Type": MustMarshal("sd1"),
			},
			wantErr: false,
		},
		{
			name: "with index",
			args: args{
				docSamples: []gonetable.Document{&WithIndex{}},
				doc: &WithIndex{
					Name: "hiihaa",
				},
			},
			want: map[string]types.AttributeValue{
				"Name":   MustMarshal("hiihaa"),
				"PK":     MustMarshal("wi#hiihaa"),
				"SK":     MustMarshal("wi"),
				"GSI1PK": MustMarshal("wi#hiihaa"),
				"GSI1SK": MustMarshal("wi"),
				"_Type":  MustMarshal("wi1"),
			},
			wantErr: false,
		},
		{
			name: "wrong document type",
			args: args{
				docSamples: []gonetable.Document{&WithIndex{}},
				doc: &MinimalDoc{
					Name: "hiihaa",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := gonetable.NewSchema(tt.args.docSamples)
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.Marshal(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schema.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				json.NewEncoder(os.Stdout).Encode(got)
				t.Errorf("Schema.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSchema_Marshal(b *testing.B) {
	s, err := gonetable.NewSchema([]gonetable.Document{&WithIndex{}})
	if err != nil {
		b.Fatal(err)
	}
	d := &WithIndex{
		Name: "withindex",
	}
	for i := 0; i < 1000; i++ {
		_, err = s.Marshal(d)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSchema_MarshalAlternative(b *testing.B) {
	s := &AlternativeSchema{}
	d := &WithIndex{
		Name: "withindex",
	}
	for i := 0; i < 1000; i++ {
		_, err := s.alternativeMarshal(d)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type AlternativeSchema struct {
}

// this would be generated for each document type
func (s *AlternativeSchema) alternativeMarshal(doc *WithIndex) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(doc)
	if err != nil {
		return nil, err
	}
	docType := doc.Gonetable_TypeID()
	av["_Type"], err = attributevalue.Marshal(docType)
	if err != nil {
		return nil, err
	}
	cKey := doc.Gonetable_Key()
	cKeyAV, err := cKey.Marshal()
	if err != nil {
		return nil, err
	}
	for k, v := range cKeyAV {
		av[k] = v
	}
	cGSI1Key := doc.Gonetable_GSI1Key()
	cGSI1KeyAV, err := cGSI1Key.Marshal()
	if err != nil {
		return nil, err
	}
	for k, v := range cGSI1KeyAV {
		av[fmt.Sprintf("GSI1%s", k)] = v
	}

	return av, nil
}
