package gonetable_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
				"PK":     MustMarshal("a#b"),
				"SK":     MustMarshal("a#b"),
				"GSI1PK": MustMarshal("a#b"),
				"GSI1SK": MustMarshal("a#b"),
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
