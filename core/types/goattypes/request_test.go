package goattypes

import (
	"reflect"
	"testing"
)

func TestUint64Codec(t *testing.T) {
	type args struct {
		n []uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1",
			args: args{[]uint64{100, 1e4}},
			want: []byte{100, 0, 0, 0, 0, 0, 0, 0, 16, 39, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "2",
			args: args{[]uint64{4294967297}},
			want: []byte{1, 0, 0, 0, 1, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeUint64(tt.args.n...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUint64() = %v, want %v", got, tt.want)
				return
			}

			decoded, err := DecodeUint64(got, len(tt.args.n))
			if err != nil {
				t.Errorf("DecodeUint64(): error = %v", err)
				return
			}

			if !reflect.DeepEqual(decoded, tt.args.n) {
				t.Errorf("not deepEqual = %v, want %v", decoded, tt.args.n)
				return
			}
		})
	}
}
