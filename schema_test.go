package gonetable_test

import (
	"testing"

	"github.com/juranki/gonetable"
)

type InvalidIndex struct {
	Name string
}

func (sd1 *InvalidIndex) Gonetable_TypeID() string { return "sd1" }
func (sd1 *InvalidIndex) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}
func (sd1 *InvalidIndex) Gonetable_AKey() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}

type MinimalDoc struct {
	Name string
}

func (sd1 *MinimalDoc) Gonetable_TypeID() string { return "sd1" }
func (sd1 *MinimalDoc) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}

type WithIndex struct {
	Name string
}

func (sd1 *WithIndex) Gonetable_TypeID() string { return "sd1" }
func (sd1 *WithIndex) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}
func (sd1 *WithIndex) Gonetable_GSI1Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}

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
