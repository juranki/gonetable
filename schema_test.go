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
			name: "no doc samples",
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
			name:       "minimal",
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
