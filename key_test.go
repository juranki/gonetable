package gonetable_test

import (
	"testing"

	"github.com/juranki/gonetable"
)

func TestJoinKey(t *testing.T) {
	type args struct {
		segments []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				segments: []string{"a", "b"},
			},
			want:    "a#b",
			wantErr: false,
		},
		{
			name: "t2",
			args: args{
				segments: []string{"#", "b"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "t3",
			args: args{
				segments: []string{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "t4",
			args: args{
				segments: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gonetable.JoinKeySegments(tt.args.segments)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
