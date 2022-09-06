package gonetable_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/juranki/gonetable"
)

func TestCompositeKey_Marshal(t *testing.T) {
	type fields struct {
		HashSegments  []string
		RangeSegments []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]types.AttributeValue
		wantErr bool
	}{
		{
			name: "simple",
			fields: fields{
				HashSegments:  []string{"a", "b"},
				RangeSegments: []string{"a", "b"},
			},
			want: map[string]types.AttributeValue{
				"PK": MustMarshal("a#b"),
				"SK": MustMarshal("a#b"),
			},
			wantErr: false,
		},
		{
			name: "error hash",
			fields: fields{
				HashSegments:  []string{"a#", "b"},
				RangeSegments: []string{"a", "b"},
			},
			wantErr: true,
		},
		{
			name: "error range",
			fields: fields{
				HashSegments:  []string{"a", "b"},
				RangeSegments: []string{"a#", "b"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := gonetable.CompositeKey{
				HashSegments:  tt.fields.HashSegments,
				RangeSegments: tt.fields.RangeSegments,
			}
			got, err := k.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompositeKey.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompositeKey.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
