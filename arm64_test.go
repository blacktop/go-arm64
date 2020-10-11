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
			name: "addg x20, x3, #0x330, #0x5",
			args: args{
				instructionValue: 0x91B31474,
				address:          0,
			},
			want: "addg	x20, x3, #0x330, #0x5",
			wantErr: false,
		},
		{
			name: "irg x20, x21, x29",
			args: args{
				instructionValue: 0x9ADD12B4,
				address:          0,
			},
			want: "irg	x20, x21, x29",
			wantErr: false,
		},
		{
			name: "st2g x16, [x10, #0x280]",
			args: args{
				instructionValue: 0xD9A28950,
				address:          0,
			},
			want: "st2g	x16, [x10, #0x280]",
			wantErr: false,
		},
		{
			name: "add w2, w3, #4095",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x3f, 0x11}),
				address:          0,
			},
			want: "add	w2, w3, #0xfff",
			wantErr: false,
		},
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
