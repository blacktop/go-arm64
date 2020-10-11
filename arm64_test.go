package arm64

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func Test_decompose(t *testing.T) {
	type args struct {
		instructionValue uint32
		address          uint64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ldr w0, #1048572",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xff, 0x7f, 0x18}),
				address:          0,
			},
			want: "ldr	w0, #0xffffc",
			wantErr: false,
		},
		{
			name: "ldr x10, #-1048576",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0a, 0x00, 0x80, 0x58}),
				address:          0,
			},
			want: "ldr	x10, #0xfffffffffff00000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			out, _ := got.disassemble()
			fmt.Println(out)
			if (err != nil) != tt.wantErr {
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(out, tt.want) {
				t.Errorf("disassemble() = %v, want %v", out, tt.want)
			}
		})
	}
}
