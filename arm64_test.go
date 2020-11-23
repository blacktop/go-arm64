package arm64

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_decompose_single_instr(t *testing.T) {
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
			name: "bfdot	v2.2s, v3.4h, v4.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x44, 0x2e}),
				address:          0,
			},
			want: "bfdot	v2.2s, v3.4h, v4.4h",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("want: %s\n", tt.want)
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_1a(t *testing.T) {
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
		//
		// llvm/test/MC/AArch64/armv8.1a-atomic.s
		//
		{
			name: "ldadda	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldadda	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldclrl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeoral	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldeoral	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldset	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldset	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxa	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldsmaxa	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x60, 0x38}),
				address:          0,
			},
			want: "ldsminlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldumaxalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumin	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldumin	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminb	w2, w3, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x50, 0x22, 0x38}),
				address:          0,
			},
			want: "ldsminb	w2, w3, [x5]",
			wantErr: false,
		},
		{
			name: "staddlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x60, 0x38}),
				address:          0,
			},
			want: "staddlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x60, 0x78}),
				address:          0,
			},
			want: "stclrlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steorl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x60, 0xb8}),
				address:          0,
			},
			want: "steorl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsetl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x60, 0xf8}),
				address:          0,
			},
			want: "stsetl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x20, 0x38}),
				address:          0,
			},
			want: "stsmaxb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x20, 0x78}),
				address:          0,
			},
			want: "stsminh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumax	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x20, 0xb8}),
				address:          0,
			},
			want: "stumax	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumin	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x70, 0x20, 0xf8}),
				address:          0,
			},
			want: "stumin	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminl	x29, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x7d, 0xf8}),
				address:          0,
			},
			want: "stsminl	x29, [sp]",
			wantErr: false,
		},
		{
			name: "swp	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x20, 0xf8}),
				address:          0,
			},
			want: "swp	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "swpb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x20, 0x38}),
				address:          0,
			},
			want: "swpb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swplh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x60, 0x78}),
				address:          0,
			},
			want: "swplh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpal	x0, x1, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x83, 0xe0, 0xf8}),
				address:          0,
			},
			want: "swpal	x0, x1, [sp]",
			wantErr: false,
		},
		{
			name: "casp	x0, x1, x2, x3, [x4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x7c, 0x20, 0x48}),
				address:          0,
			},
			want: "casp	x0, x1, x2, x3, [x4]",
			wantErr: false,
		},
		{
			name: "casp	w0, w1, w2, w3, [x4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x7c, 0x20, 0x08}),
				address:          0,
			},
			want: "casp	w0, w1, w2, w3, [x4]",
			wantErr: false,
		},
		//
		// llvm/test/MC/AArch64/armv8.1a-lor.s
		//
		{
			name: "ldlarb	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0xdf, 0x08}),
				address:          0,
			},
			want: "ldlarb	w0, [x1]",
			wantErr: false,
		},
		{
			name: "ldlarh	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0xdf, 0x48}),
				address:          0,
			},
			want: "ldlarh	w0, [x1]",
			wantErr: false,
		},
		{
			name: "ldlar	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0xdf, 0x88}),
				address:          0,
			},
			want: "ldlar	w0, [x1]",
			wantErr: false,
		},
		{
			name: "ldlar	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0xdf, 0xc8}),
				address:          0,
			},
			want: "ldlar	x0, [x1]",
			wantErr: false,
		},
		{
			name: "stllrb	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0x9f, 0x08}),
				address:          0,
			},
			want: "stllrb	w0, [x1]",
			wantErr: false,
		},
		{
			name: "stllrh	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0x9f, 0x48}),
				address:          0,
			},
			want: "stllrh	w0, [x1]",
			wantErr: false,
		},
		{
			name: "stllr	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0x9f, 0x88}),
				address:          0,
			},
			want: "stllr	w0, [x1]",
			wantErr: false,
		},
		{
			name: "stllr	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7c, 0x9f, 0xc8}),
				address:          0,
			},
			want: "stllr	x0, [x1]",
			wantErr: false,
		},
		{
			name: "msr	lorsa_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xa4, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	lorsa_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	lorea_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xa4, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	lorea_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	lorn_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xa4, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	lorn_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	lorc_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xa4, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	lorc_el1, x0",
			wantErr: false,
		},
		{
			name: "mrs	x0, lorid_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xa4, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, lorid_el1",
			wantErr: false,
		},
		//
		// llvm/test/MC/AArch64/armv8.1a-pan.s
		//
		{
			name: "msr	pan, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x40, 0x00, 0xd5}),
				address:          0,
			},
			want: "msr	pan, #0",
			wantErr: false,
		},
		{
			name: "msr	pan, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x41, 0x00, 0xd5}),
				address:          0,
			},
			want: "msr	pan, #1",
			wantErr: false,
		},
		{
			name: "msr	pan, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x65, 0x42, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pan, x5",
			wantErr: false,
		},
		{
			name: "mrs	x13, pan",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6d, 0x42, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x13, pan",
			wantErr: false,
		},
		//
		// llvm/test/MC/AArch64/armv8.1a-rdma.s
		//
		{
			name: "sqrdmlah	v0.4h, v1.4h, v2.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x42, 0x2e}),
				address:          0,
			},
			want: "sqrdmlah	v0.4h, v1.4h, v2.4h",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.4h, v1.4h, v2.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x8c, 0x42, 0x2e}),
				address:          0,
			},
			want: "sqrdmlsh	v0.4h, v1.4h, v2.4h",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.2s, v1.2s, v2.2s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x82, 0x2e}),
				address:          0,
			},
			want: "sqrdmlah	v0.2s, v1.2s, v2.2s",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.2s, v1.2s, v2.2s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x8c, 0x82, 0x2e}),
				address:          0,
			},
			want: "sqrdmlsh	v0.2s, v1.2s, v2.2s",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.4s, v1.4s, v2.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x82, 0x6e}),
				address:          0,
			},
			want: "sqrdmlah	v0.4s, v1.4s, v2.4s",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.4s, v1.4s, v2.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x8c, 0x82, 0x6e}),
				address:          0,
			},
			want: "sqrdmlsh	v0.4s, v1.4s, v2.4s",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.8h, v1.8h, v2.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x42, 0x6e}),
				address:          0,
			},
			want: "sqrdmlah	v0.8h, v1.8h, v2.8h",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.8h, v1.8h, v2.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x8c, 0x42, 0x6e}),
				address:          0,
			},
			want: "sqrdmlsh	v0.8h, v1.8h, v2.8h",
			wantErr: false,
		},
		{
			name: "sqrdmlah	h0, h1, h2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x42, 0x7e}),
				address:          0,
			},
			want: "sqrdmlah	h0, h1, h2",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	h0, h1, h2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x8c, 0x42, 0x7e}),
				address:          0,
			},
			want: "sqrdmlsh	h0, h1, h2",
			wantErr: false,
		},
		{
			name: "sqrdmlah	s0, s1, s2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x82, 0x7e}),
				address:          0,
			},
			want: "sqrdmlah	s0, s1, s2",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	s0, s1, s2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x8c, 0x82, 0x7e}),
				address:          0,
			},
			want: "sqrdmlsh	s0, s1, s2",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.4h, v1.4h, v2.h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd0, 0x72, 0x2f}),
				address:          0,
			},
			want: "sqrdmlah	v0.4h, v1.4h, v2.h[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.4h, v1.4h, v2.h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf0, 0x72, 0x2f}),
				address:          0,
			},
			want: "sqrdmlsh	v0.4h, v1.4h, v2.h[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.2s, v1.2s, v2.s[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd0, 0xa2, 0x2f}),
				address:          0,
			},
			want: "sqrdmlah	v0.2s, v1.2s, v2.s[1]",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.2s, v1.2s, v2.s[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf0, 0xa2, 0x2f}),
				address:          0,
			},
			want: "sqrdmlsh	v0.2s, v1.2s, v2.s[1]",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.8h, v1.8h, v2.h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd0, 0x72, 0x6f}),
				address:          0,
			},
			want: "sqrdmlah	v0.8h, v1.8h, v2.h[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.8h, v1.8h, v2.h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf0, 0x72, 0x6f}),
				address:          0,
			},
			want: "sqrdmlsh	v0.8h, v1.8h, v2.h[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlah	v0.4s, v1.4s, v2.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd8, 0xa2, 0x6f}),
				address:          0,
			},
			want: "sqrdmlah	v0.4s, v1.4s, v2.s[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	v0.4s, v1.4s, v2.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf8, 0xa2, 0x6f}),
				address:          0,
			},
			want: "sqrdmlsh	v0.4s, v1.4s, v2.s[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlah	h0, h1, v2.h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd0, 0x72, 0x7f}),
				address:          0,
			},
			want: "sqrdmlah	h0, h1, v2.h[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	h0, h1, v2.h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf0, 0x72, 0x7f}),
				address:          0,
			},
			want: "sqrdmlsh	h0, h1, v2.h[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlah	s0, s1, v2.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd8, 0xa2, 0x7f}),
				address:          0,
			},
			want: "sqrdmlah	s0, s1, v2.s[3]",
			wantErr: false,
		},
		{
			name: "sqrdmlsh	s0, s1, v2.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf8, 0xa2, 0x7f}),
				address:          0,
			},
			want: "sqrdmlsh	s0, s1, v2.s[3]",
			wantErr: false,
		},
		//
		// llvm/test/MC/AArch64/armv8.1a-vhe.s
		//
		{
			name: "msr	ttbr1_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x20, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr1_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	contextidr_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	contextidr_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	cnthv_tval_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xe3, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthv_tval_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	cnthv_cval_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xe3, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthv_cval_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	cnthv_ctl_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe3, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthv_ctl_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	sctlr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x10, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	sctlr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	cpacr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x10, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cpacr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	ttbr0_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x20, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr0_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	ttbr1_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x20, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr1_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	tcr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x20, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	tcr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	afsr0_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x51, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	afsr0_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	afsr1_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x51, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	afsr1_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	esr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x52, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	esr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	far_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x60, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	far_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	mair_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xa2, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	mair_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	amair_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xa3, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	amair_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	vbar_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xc0, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	vbar_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	contextidr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd0, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	contextidr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	cntkctl_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xe1, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntkctl_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	cntp_tval_el02, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xe2, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntp_tval_el02, x0",
			wantErr: false,
		},
		{
			name: "msr	cntp_ctl_el02, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe2, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntp_ctl_el02, x0",
			wantErr: false,
		},
		{
			name: "msr	cntp_cval_el02, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xe2, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntp_cval_el02, x0",
			wantErr: false,
		},
		{
			name: "msr	cntv_tval_el02, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xe3, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntv_tval_el02, x0",
			wantErr: false,
		},
		{
			name: "msr	cntv_ctl_el02, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe3, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntv_ctl_el02, x0",
			wantErr: false,
		},
		{
			name: "msr	cntv_cval_el02, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xe3, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	cntv_cval_el02, x0",
			wantErr: false,
		},
		{
			name: "msr	spsr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x40, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	elr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x40, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	elr_el12, x0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_1a_LSE(t *testing.T) {
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
		//
		// llvm/test/MC/AArch64/armv8.1a-lse.s
		//
		{
			name: "cas	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xa0, 0x88}),
				address:          0,
			},
			want: "cas	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "cas	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xa2, 0x88}),
				address:          0,
			},
			want: "cas	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casa	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xe0, 0x88}),
				address:          0,
			},
			want: "casa	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casa	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xe2, 0x88}),
				address:          0,
			},
			want: "casa	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xa0, 0x88}),
				address:          0,
			},
			want: "casl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xa2, 0x88}),
				address:          0,
			},
			want: "casl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xe0, 0x88}),
				address:          0,
			},
			want: "casal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xe2, 0x88}),
				address:          0,
			},
			want: "casal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xa0, 0x08}),
				address:          0,
			},
			want: "casb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xa2, 0x08}),
				address:          0,
			},
			want: "casb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "cash	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xa0, 0x48}),
				address:          0,
			},
			want: "cash	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "cash	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xa2, 0x48}),
				address:          0,
			},
			want: "cash	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xe0, 0x08}),
				address:          0,
			},
			want: "casab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xe2, 0x08}),
				address:          0,
			},
			want: "casab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "caslb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xa0, 0x08}),
				address:          0,
			},
			want: "caslb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "caslb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xa2, 0x08}),
				address:          0,
			},
			want: "caslb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xe0, 0x08}),
				address:          0,
			},
			want: "casalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xe2, 0x08}),
				address:          0,
			},
			want: "casalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xe0, 0x48}),
				address:          0,
			},
			want: "casah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xe2, 0x48}),
				address:          0,
			},
			want: "casah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "caslh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xa0, 0x48}),
				address:          0,
			},
			want: "caslh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "caslh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xa2, 0x48}),
				address:          0,
			},
			want: "caslh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "casalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xe0, 0x48}),
				address:          0,
			},
			want: "casalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "casalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xe2, 0x48}),
				address:          0,
			},
			want: "casalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "cas	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xa0, 0xc8}),
				address:          0,
			},
			want: "cas	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "cas	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xa2, 0xc8}),
				address:          0,
			},
			want: "cas	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "casa	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x7c, 0xe0, 0xc8}),
				address:          0,
			},
			want: "casa	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "casa	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0xe2, 0xc8}),
				address:          0,
			},
			want: "casa	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "casl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xa0, 0xc8}),
				address:          0,
			},
			want: "casl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "casl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xa2, 0xc8}),
				address:          0,
			},
			want: "casl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "casal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xe0, 0xc8}),
				address:          0,
			},
			want: "casal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "casal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xe2, 0xc8}),
				address:          0,
			},
			want: "casal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "swp	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x20, 0xb8}),
				address:          0,
			},
			want: "swp	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swp	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x22, 0xb8}),
				address:          0,
			},
			want: "swp	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpa	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xa0, 0xb8}),
				address:          0,
			},
			want: "swpa	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpa	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xa2, 0xb8}),
				address:          0,
			},
			want: "swpa	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x60, 0xb8}),
				address:          0,
			},
			want: "swpl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x62, 0xb8}),
				address:          0,
			},
			want: "swpl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xe0, 0xb8}),
				address:          0,
			},
			want: "swpal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xe2, 0xb8}),
				address:          0,
			},
			want: "swpal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x20, 0x38}),
				address:          0,
			},
			want: "swpb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x22, 0x38}),
				address:          0,
			},
			want: "swpb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swph	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x20, 0x78}),
				address:          0,
			},
			want: "swph	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swph	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x22, 0x78}),
				address:          0,
			},
			want: "swph	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xa0, 0x38}),
				address:          0,
			},
			want: "swpab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xa2, 0x38}),
				address:          0,
			},
			want: "swpab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swplb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x60, 0x38}),
				address:          0,
			},
			want: "swplb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swplb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x62, 0x38}),
				address:          0,
			},
			want: "swplb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xe0, 0x38}),
				address:          0,
			},
			want: "swpalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xe2, 0x38}),
				address:          0,
			},
			want: "swpalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xa0, 0x78}),
				address:          0,
			},
			want: "swpah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xa2, 0x78}),
				address:          0,
			},
			want: "swpah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swplh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x60, 0x78}),
				address:          0,
			},
			want: "swplh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swplh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x62, 0x78}),
				address:          0,
			},
			want: "swplh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swpalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xe0, 0x78}),
				address:          0,
			},
			want: "swpalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "swpalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xe2, 0x78}),
				address:          0,
			},
			want: "swpalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "swp	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x20, 0xf8}),
				address:          0,
			},
			want: "swp	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "swp	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x22, 0xf8}),
				address:          0,
			},
			want: "swp	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "swpa	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xa0, 0xf8}),
				address:          0,
			},
			want: "swpa	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "swpa	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xa2, 0xf8}),
				address:          0,
			},
			want: "swpa	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "swpl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0x60, 0xf8}),
				address:          0,
			},
			want: "swpl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "swpl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0x62, 0xf8}),
				address:          0,
			},
			want: "swpl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "swpal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x80, 0xe0, 0xf8}),
				address:          0,
			},
			want: "swpal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "swpal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x83, 0xe2, 0xf8}),
				address:          0,
			},
			want: "swpal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "casp	w0, w1, w2, w3, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x7c, 0x20, 0x08}),
				address:          0,
			},
			want: "casp	w0, w1, w2, w3, [x5]",
			wantErr: false,
		},
		{
			name: "casp	w4, w5, w6, w7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x7f, 0x24, 0x08}),
				address:          0,
			},
			want: "casp	w4, w5, w6, w7, [sp]",
			wantErr: false,
		},
		{
			name: "casp	x0, x1, x2, x3, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0x7c, 0x20, 0x48}),
				address:          0,
			},
			want: "casp	x0, x1, x2, x3, [x2]",
			wantErr: false,
		},
		{
			name: "casp	x4, x5, x6, x7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x7f, 0x24, 0x48}),
				address:          0,
			},
			want: "casp	x4, x5, x6, x7, [sp]",
			wantErr: false,
		},
		{
			name: "caspa	w0, w1, w2, w3, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x7c, 0x60, 0x08}),
				address:          0,
			},
			want: "caspa	w0, w1, w2, w3, [x5]",
			wantErr: false,
		},
		{
			name: "caspa	w4, w5, w6, w7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x7f, 0x64, 0x08}),
				address:          0,
			},
			want: "caspa	w4, w5, w6, w7, [sp]",
			wantErr: false,
		},
		{
			name: "caspa	x0, x1, x2, x3, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0x7c, 0x60, 0x48}),
				address:          0,
			},
			want: "caspa	x0, x1, x2, x3, [x2]",
			wantErr: false,
		},
		{
			name: "caspa	x4, x5, x6, x7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x7f, 0x64, 0x48}),
				address:          0,
			},
			want: "caspa	x4, x5, x6, x7, [sp]",
			wantErr: false,
		},
		{
			name: "caspl	w0, w1, w2, w3, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0xfc, 0x20, 0x08}),
				address:          0,
			},
			want: "caspl	w0, w1, w2, w3, [x5]",
			wantErr: false,
		},
		{
			name: "caspl	w4, w5, w6, w7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xff, 0x24, 0x08}),
				address:          0,
			},
			want: "caspl	w4, w5, w6, w7, [sp]",
			wantErr: false,
		},
		{
			name: "caspl	x0, x1, x2, x3, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0xfc, 0x20, 0x48}),
				address:          0,
			},
			want: "caspl	x0, x1, x2, x3, [x2]",
			wantErr: false,
		},
		{
			name: "caspl	x4, x5, x6, x7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xff, 0x24, 0x48}),
				address:          0,
			},
			want: "caspl	x4, x5, x6, x7, [sp]",
			wantErr: false,
		},
		{
			name: "caspal	w0, w1, w2, w3, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0xfc, 0x60, 0x08}),
				address:          0,
			},
			want: "caspal	w0, w1, w2, w3, [x5]",
			wantErr: false,
		},
		{
			name: "caspal	w4, w5, w6, w7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xff, 0x64, 0x08}),
				address:          0,
			},
			want: "caspal	w4, w5, w6, w7, [sp]",
			wantErr: false,
		},
		{
			name: "caspal	x0, x1, x2, x3, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0xfc, 0x60, 0x48}),
				address:          0,
			},
			want: "caspal	x0, x1, x2, x3, [x2]",
			wantErr: false,
		},
		{
			name: "caspal	x4, x5, x6, x7, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xff, 0x64, 0x48}),
				address:          0,
			},
			want: "caspal	x4, x5, x6, x7, [sp]",
			wantErr: false,
		},
		{
			name: "ldadd	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldadd	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldadd	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldadd	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldadda	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldadda	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldadda	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldadda	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldaddl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldaddl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldaddal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldaddal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x20, 0x38}),
				address:          0,
			},
			want: "ldaddb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x22, 0x38}),
				address:          0,
			},
			want: "ldaddb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x20, 0x78}),
				address:          0,
			},
			want: "ldaddh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x22, 0x78}),
				address:          0,
			},
			want: "ldaddh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldaddab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldaddab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x60, 0x38}),
				address:          0,
			},
			want: "ldaddlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x62, 0x38}),
				address:          0,
			},
			want: "ldaddlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldaddalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldaddalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldaddah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldaddah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x60, 0x78}),
				address:          0,
			},
			want: "ldaddlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x62, 0x78}),
				address:          0,
			},
			want: "ldaddlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldaddalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldaddalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldadd	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldadd	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldadd	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldadd	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldadda	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldadda	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldadda	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldadda	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldaddl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldaddl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldaddal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x00, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldaddal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldaddal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldaddal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclr	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldclr	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclr	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldclr	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclra	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldclra	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclra	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldclra	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldclrl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldclrl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclral	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldclral	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclral	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldclral	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x20, 0x38}),
				address:          0,
			},
			want: "ldclrb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x22, 0x38}),
				address:          0,
			},
			want: "ldclrb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x20, 0x78}),
				address:          0,
			},
			want: "ldclrh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x22, 0x78}),
				address:          0,
			},
			want: "ldclrh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldclrab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldclrab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x60, 0x38}),
				address:          0,
			},
			want: "ldclrlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x62, 0x38}),
				address:          0,
			},
			want: "ldclrlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclralb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldclralb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclralb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldclralb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldclrah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldclrah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x60, 0x78}),
				address:          0,
			},
			want: "ldclrlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x62, 0x78}),
				address:          0,
			},
			want: "ldclrlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclralh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldclralh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclralh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldclralh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclr	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldclr	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclr	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldclr	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclra	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldclra	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclra	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldclra	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclrl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldclrl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclrl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldclrl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldclral	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldclral	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldclral	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x13, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldclral	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeor	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldeor	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeor	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldeor	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeora	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldeora	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeora	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldeora	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldeorl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldeorl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeoral	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldeoral	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeoral	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldeoral	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x20, 0x38}),
				address:          0,
			},
			want: "ldeorb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x22, 0x38}),
				address:          0,
			},
			want: "ldeorb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x20, 0x78}),
				address:          0,
			},
			want: "ldeorh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x22, 0x78}),
				address:          0,
			},
			want: "ldeorh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldeorab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldeorab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x60, 0x38}),
				address:          0,
			},
			want: "ldeorlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x62, 0x38}),
				address:          0,
			},
			want: "ldeorlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeoralb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldeoralb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeoralb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldeoralb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldeorah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldeorah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x60, 0x78}),
				address:          0,
			},
			want: "ldeorlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x62, 0x78}),
				address:          0,
			},
			want: "ldeorlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeoralh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldeoralh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeoralh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldeoralh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeor	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldeor	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeor	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldeor	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeora	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldeora	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeora	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldeora	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeorl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldeorl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeorl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldeorl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldeoral	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x20, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldeoral	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldeoral	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x23, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldeoral	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldset	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldset	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldset	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldset	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldseta	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldseta	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldseta	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldseta	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldsetl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldsetl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldsetal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldsetal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x20, 0x38}),
				address:          0,
			},
			want: "ldsetb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x22, 0x38}),
				address:          0,
			},
			want: "ldsetb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldseth	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x20, 0x78}),
				address:          0,
			},
			want: "ldseth	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldseth	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x22, 0x78}),
				address:          0,
			},
			want: "ldseth	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldsetab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldsetab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x60, 0x38}),
				address:          0,
			},
			want: "ldsetlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x62, 0x38}),
				address:          0,
			},
			want: "ldsetlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldsetalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldsetalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldsetah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldsetah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x60, 0x78}),
				address:          0,
			},
			want: "ldsetlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x62, 0x78}),
				address:          0,
			},
			want: "ldsetlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldsetalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldsetalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldset	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldset	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldset	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldset	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldseta	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldseta	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldseta	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldseta	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldsetl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldsetl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsetal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x30, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldsetal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsetal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x33, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldsetal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmax	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldsmax	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmax	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldsmax	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxa	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldsmaxa	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxa	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldsmaxa	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldsmaxl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldsmaxl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldsmaxal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldsmaxal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x20, 0x38}),
				address:          0,
			},
			want: "ldsmaxb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x22, 0x38}),
				address:          0,
			},
			want: "ldsmaxb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x20, 0x78}),
				address:          0,
			},
			want: "ldsmaxh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x22, 0x78}),
				address:          0,
			},
			want: "ldsmaxh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldsmaxab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldsmaxab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x60, 0x38}),
				address:          0,
			},
			want: "ldsmaxlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x62, 0x38}),
				address:          0,
			},
			want: "ldsmaxlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldsmaxalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldsmaxalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldsmaxah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldsmaxah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x60, 0x78}),
				address:          0,
			},
			want: "ldsmaxlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x62, 0x78}),
				address:          0,
			},
			want: "ldsmaxlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldsmaxalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldsmaxalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmax	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldsmax	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmax	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldsmax	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxa	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldsmaxa	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxa	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldsmaxa	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldsmaxl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldsmaxl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmaxal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x40, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldsmaxal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmaxal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x43, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldsmaxal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmin	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldsmin	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmin	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldsmin	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmina	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldsmina	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmina	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldsmina	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldsminl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldsminl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldsminal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldsminal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x20, 0x38}),
				address:          0,
			},
			want: "ldsminb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x22, 0x38}),
				address:          0,
			},
			want: "ldsminb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x20, 0x78}),
				address:          0,
			},
			want: "ldsminh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x22, 0x78}),
				address:          0,
			},
			want: "ldsminh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldsminab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldsminab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x60, 0x38}),
				address:          0,
			},
			want: "ldsminlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x62, 0x38}),
				address:          0,
			},
			want: "ldsminlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldsminalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldsminalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldsminah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldsminah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x60, 0x78}),
				address:          0,
			},
			want: "ldsminlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x62, 0x78}),
				address:          0,
			},
			want: "ldsminlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldsminalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldsminalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmin	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldsmin	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmin	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldsmin	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsmina	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldsmina	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsmina	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldsmina	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldsminl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldsminl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldsminal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x50, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldsminal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldsminal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x53, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldsminal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumax	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldumax	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumax	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldumax	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxa	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldumaxa	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxa	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldumaxa	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x60, 0xb8}),
				address:          0,
			},
			want: "ldumaxl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x62, 0xb8}),
				address:          0,
			},
			want: "ldumaxl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xe0, 0xb8}),
				address:          0,
			},
			want: "ldumaxal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xe2, 0xb8}),
				address:          0,
			},
			want: "ldumaxal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x20, 0x38}),
				address:          0,
			},
			want: "ldumaxb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x22, 0x38}),
				address:          0,
			},
			want: "ldumaxb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x20, 0x78}),
				address:          0,
			},
			want: "ldumaxh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x22, 0x78}),
				address:          0,
			},
			want: "ldumaxh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xa0, 0x38}),
				address:          0,
			},
			want: "ldumaxab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xa2, 0x38}),
				address:          0,
			},
			want: "ldumaxab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x60, 0x38}),
				address:          0,
			},
			want: "ldumaxlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x62, 0x38}),
				address:          0,
			},
			want: "ldumaxlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xe0, 0x38}),
				address:          0,
			},
			want: "ldumaxalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xe2, 0x38}),
				address:          0,
			},
			want: "ldumaxalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xa0, 0x78}),
				address:          0,
			},
			want: "ldumaxah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xa2, 0x78}),
				address:          0,
			},
			want: "ldumaxah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x60, 0x78}),
				address:          0,
			},
			want: "ldumaxlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x62, 0x78}),
				address:          0,
			},
			want: "ldumaxlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xe0, 0x78}),
				address:          0,
			},
			want: "ldumaxalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xe2, 0x78}),
				address:          0,
			},
			want: "ldumaxalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumax	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldumax	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumax	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldumax	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxa	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldumaxa	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxa	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldumaxa	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldumaxl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0x62, 0xf8}),
				address:          0,
			},
			want: "ldumaxl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumaxal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x60, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldumaxal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumaxal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x63, 0xe2, 0xf8}),
				address:          0,
			},
			want: "ldumaxal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumin	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x20, 0xb8}),
				address:          0,
			},
			want: "ldumin	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumin	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x22, 0xb8}),
				address:          0,
			},
			want: "ldumin	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumina	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xa0, 0xb8}),
				address:          0,
			},
			want: "ldumina	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumina	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xa2, 0xb8}),
				address:          0,
			},
			want: "ldumina	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminl	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x60, 0xb8}),
				address:          0,
			},
			want: "lduminl	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminl	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x62, 0xb8}),
				address:          0,
			},
			want: "lduminl	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminal	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xe0, 0xb8}),
				address:          0,
			},
			want: "lduminal	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminal	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xe2, 0xb8}),
				address:          0,
			},
			want: "lduminal	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x20, 0x38}),
				address:          0,
			},
			want: "lduminb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x22, 0x38}),
				address:          0,
			},
			want: "lduminb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x20, 0x78}),
				address:          0,
			},
			want: "lduminh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x22, 0x78}),
				address:          0,
			},
			want: "lduminh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminab	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xa0, 0x38}),
				address:          0,
			},
			want: "lduminab	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminab	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xa2, 0x38}),
				address:          0,
			},
			want: "lduminab	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminlb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x60, 0x38}),
				address:          0,
			},
			want: "lduminlb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminlb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x62, 0x38}),
				address:          0,
			},
			want: "lduminlb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminalb	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xe0, 0x38}),
				address:          0,
			},
			want: "lduminalb	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminalb	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xe2, 0x38}),
				address:          0,
			},
			want: "lduminalb	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminah	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xa0, 0x78}),
				address:          0,
			},
			want: "lduminah	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminah	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xa2, 0x78}),
				address:          0,
			},
			want: "lduminah	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminlh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x60, 0x78}),
				address:          0,
			},
			want: "lduminlh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminlh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x62, 0x78}),
				address:          0,
			},
			want: "lduminlh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminalh	w0, w1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xe0, 0x78}),
				address:          0,
			},
			want: "lduminalh	w0, w1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminalh	w2, w3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xe2, 0x78}),
				address:          0,
			},
			want: "lduminalh	w2, w3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumin	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldumin	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumin	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x22, 0xf8}),
				address:          0,
			},
			want: "ldumin	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "ldumina	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldumina	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "ldumina	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xa2, 0xf8}),
				address:          0,
			},
			want: "ldumina	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminl	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0x60, 0xf8}),
				address:          0,
			},
			want: "lduminl	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminl	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0x62, 0xf8}),
				address:          0,
			},
			want: "lduminl	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "lduminal	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x70, 0xe0, 0xf8}),
				address:          0,
			},
			want: "lduminal	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "lduminal	x2, x3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x73, 0xe2, 0xf8}),
				address:          0,
			},
			want: "lduminal	x2, x3, [sp]",
			wantErr: false,
		},
		{
			name: "stadd	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x20, 0xb8}),
				address:          0,
			},
			want: "stadd	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stadd	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x22, 0xb8}),
				address:          0,
			},
			want: "stadd	w2, [sp]",
			wantErr: false,
		},
		{
			name: "staddl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x60, 0xb8}),
				address:          0,
			},
			want: "staddl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "staddl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x62, 0xb8}),
				address:          0,
			},
			want: "staddl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "staddb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x20, 0x38}),
				address:          0,
			},
			want: "staddb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "staddb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x22, 0x38}),
				address:          0,
			},
			want: "staddb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "staddh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x20, 0x78}),
				address:          0,
			},
			want: "staddh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "staddh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x22, 0x78}),
				address:          0,
			},
			want: "staddh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "staddlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x60, 0x38}),
				address:          0,
			},
			want: "staddlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "staddlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x62, 0x38}),
				address:          0,
			},
			want: "staddlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "staddlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x60, 0x78}),
				address:          0,
			},
			want: "staddlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "staddlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x62, 0x78}),
				address:          0,
			},
			want: "staddlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stadd	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x20, 0xf8}),
				address:          0,
			},
			want: "stadd	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stadd	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x22, 0xf8}),
				address:          0,
			},
			want: "stadd	x2, [sp]",
			wantErr: false,
		},
		{
			name: "staddl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x60, 0xf8}),
				address:          0,
			},
			want: "staddl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "staddl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x62, 0xf8}),
				address:          0,
			},
			want: "staddl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stclr	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x20, 0xb8}),
				address:          0,
			},
			want: "stclr	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclr	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x22, 0xb8}),
				address:          0,
			},
			want: "stclr	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stclrl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x60, 0xb8}),
				address:          0,
			},
			want: "stclrl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x62, 0xb8}),
				address:          0,
			},
			want: "stclrl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stclrb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x20, 0x38}),
				address:          0,
			},
			want: "stclrb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x22, 0x38}),
				address:          0,
			},
			want: "stclrb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stclrh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x20, 0x78}),
				address:          0,
			},
			want: "stclrh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x22, 0x78}),
				address:          0,
			},
			want: "stclrh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stclrlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x60, 0x38}),
				address:          0,
			},
			want: "stclrlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x62, 0x38}),
				address:          0,
			},
			want: "stclrlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stclrlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x60, 0x78}),
				address:          0,
			},
			want: "stclrlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x62, 0x78}),
				address:          0,
			},
			want: "stclrlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stclr	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x20, 0xf8}),
				address:          0,
			},
			want: "stclr	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stclr	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x22, 0xf8}),
				address:          0,
			},
			want: "stclr	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stclrl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x10, 0x60, 0xf8}),
				address:          0,
			},
			want: "stclrl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stclrl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x62, 0xf8}),
				address:          0,
			},
			want: "stclrl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "steor	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x20, 0xb8}),
				address:          0,
			},
			want: "steor	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steor	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x22, 0xb8}),
				address:          0,
			},
			want: "steor	w2, [sp]",
			wantErr: false,
		},
		{
			name: "steorl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x60, 0xb8}),
				address:          0,
			},
			want: "steorl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steorl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x62, 0xb8}),
				address:          0,
			},
			want: "steorl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "steorb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x20, 0x38}),
				address:          0,
			},
			want: "steorb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steorb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x22, 0x38}),
				address:          0,
			},
			want: "steorb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "steorh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x20, 0x78}),
				address:          0,
			},
			want: "steorh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steorh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x22, 0x78}),
				address:          0,
			},
			want: "steorh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "steorlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x60, 0x38}),
				address:          0,
			},
			want: "steorlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steorlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x62, 0x38}),
				address:          0,
			},
			want: "steorlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "steorlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x60, 0x78}),
				address:          0,
			},
			want: "steorlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "steorlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x62, 0x78}),
				address:          0,
			},
			want: "steorlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "steor	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x20, 0xf8}),
				address:          0,
			},
			want: "steor	x0, [x2]",
			wantErr: false,
		},
		{
			name: "steor	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x22, 0xf8}),
				address:          0,
			},
			want: "steor	x2, [sp]",
			wantErr: false,
		},
		{
			name: "steorl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x60, 0xf8}),
				address:          0,
			},
			want: "steorl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "steorl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x62, 0xf8}),
				address:          0,
			},
			want: "steorl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stset	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x20, 0xb8}),
				address:          0,
			},
			want: "stset	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stset	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x22, 0xb8}),
				address:          0,
			},
			want: "stset	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsetl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x60, 0xb8}),
				address:          0,
			},
			want: "stsetl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsetl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x62, 0xb8}),
				address:          0,
			},
			want: "stsetl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsetb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x20, 0x38}),
				address:          0,
			},
			want: "stsetb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsetb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x22, 0x38}),
				address:          0,
			},
			want: "stsetb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stseth	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x20, 0x78}),
				address:          0,
			},
			want: "stseth	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stseth	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x22, 0x78}),
				address:          0,
			},
			want: "stseth	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsetlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x60, 0x38}),
				address:          0,
			},
			want: "stsetlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsetlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x62, 0x38}),
				address:          0,
			},
			want: "stsetlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsetlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x60, 0x78}),
				address:          0,
			},
			want: "stsetlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsetlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x62, 0x78}),
				address:          0,
			},
			want: "stsetlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stset	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x20, 0xf8}),
				address:          0,
			},
			want: "stset	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stset	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x22, 0xf8}),
				address:          0,
			},
			want: "stset	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stsetl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x60, 0xf8}),
				address:          0,
			},
			want: "stsetl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsetl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x62, 0xf8}),
				address:          0,
			},
			want: "stsetl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmax	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x20, 0xb8}),
				address:          0,
			},
			want: "stsmax	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmax	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x22, 0xb8}),
				address:          0,
			},
			want: "stsmax	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmaxl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x60, 0xb8}),
				address:          0,
			},
			want: "stsmaxl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x62, 0xb8}),
				address:          0,
			},
			want: "stsmaxl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmaxb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x20, 0x38}),
				address:          0,
			},
			want: "stsmaxb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x22, 0x38}),
				address:          0,
			},
			want: "stsmaxb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmaxh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x20, 0x78}),
				address:          0,
			},
			want: "stsmaxh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x22, 0x78}),
				address:          0,
			},
			want: "stsmaxh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmaxlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x60, 0x38}),
				address:          0,
			},
			want: "stsmaxlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x62, 0x38}),
				address:          0,
			},
			want: "stsmaxlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmaxlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x60, 0x78}),
				address:          0,
			},
			want: "stsmaxlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x62, 0x78}),
				address:          0,
			},
			want: "stsmaxlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmax	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x20, 0xf8}),
				address:          0,
			},
			want: "stsmax	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmax	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x22, 0xf8}),
				address:          0,
			},
			want: "stsmax	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmaxl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x60, 0xf8}),
				address:          0,
			},
			want: "stsmaxl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmaxl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x62, 0xf8}),
				address:          0,
			},
			want: "stsmaxl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmin	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x20, 0xb8}),
				address:          0,
			},
			want: "stsmin	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmin	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x22, 0xb8}),
				address:          0,
			},
			want: "stsmin	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsminl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x60, 0xb8}),
				address:          0,
			},
			want: "stsminl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x62, 0xb8}),
				address:          0,
			},
			want: "stsminl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsminb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x20, 0x38}),
				address:          0,
			},
			want: "stsminb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x22, 0x38}),
				address:          0,
			},
			want: "stsminb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsminh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x20, 0x78}),
				address:          0,
			},
			want: "stsminh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x22, 0x78}),
				address:          0,
			},
			want: "stsminh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsminlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x60, 0x38}),
				address:          0,
			},
			want: "stsminlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x62, 0x38}),
				address:          0,
			},
			want: "stsminlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsminlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x60, 0x78}),
				address:          0,
			},
			want: "stsminlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x62, 0x78}),
				address:          0,
			},
			want: "stsminlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stsmin	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x20, 0xf8}),
				address:          0,
			},
			want: "stsmin	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsmin	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x22, 0xf8}),
				address:          0,
			},
			want: "stsmin	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stsminl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x50, 0x60, 0xf8}),
				address:          0,
			},
			want: "stsminl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stsminl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x62, 0xf8}),
				address:          0,
			},
			want: "stsminl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stumax	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x20, 0xb8}),
				address:          0,
			},
			want: "stumax	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumax	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x22, 0xb8}),
				address:          0,
			},
			want: "stumax	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stumaxl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x60, 0xb8}),
				address:          0,
			},
			want: "stumaxl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumaxl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x62, 0xb8}),
				address:          0,
			},
			want: "stumaxl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stumaxb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x20, 0x38}),
				address:          0,
			},
			want: "stumaxb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumaxb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x22, 0x38}),
				address:          0,
			},
			want: "stumaxb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stumaxh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x20, 0x78}),
				address:          0,
			},
			want: "stumaxh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumaxh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x22, 0x78}),
				address:          0,
			},
			want: "stumaxh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stumaxlb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x60, 0x38}),
				address:          0,
			},
			want: "stumaxlb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumaxlb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x62, 0x38}),
				address:          0,
			},
			want: "stumaxlb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stumaxlh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x60, 0x78}),
				address:          0,
			},
			want: "stumaxlh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumaxlh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x62, 0x78}),
				address:          0,
			},
			want: "stumaxlh	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stumax	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x20, 0xf8}),
				address:          0,
			},
			want: "stumax	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stumax	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x22, 0xf8}),
				address:          0,
			},
			want: "stumax	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stumaxl	x0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x60, 0x60, 0xf8}),
				address:          0,
			},
			want: "stumaxl	x0, [x2]",
			wantErr: false,
		},
		{
			name: "stumaxl	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x62, 0xf8}),
				address:          0,
			},
			want: "stumaxl	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stumin	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x70, 0x20, 0xb8}),
				address:          0,
			},
			want: "stumin	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stumin	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x73, 0x22, 0xb8}),
				address:          0,
			},
			want: "stumin	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stuminl	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x70, 0x60, 0xb8}),
				address:          0,
			},
			want: "stuminl	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stuminl	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x73, 0x62, 0xb8}),
				address:          0,
			},
			want: "stuminl	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stuminb	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x70, 0x20, 0x38}),
				address:          0,
			},
			want: "stuminb	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stuminb	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x73, 0x22, 0x38}),
				address:          0,
			},
			want: "stuminb	w2, [sp]",
			wantErr: false,
		},
		{
			name: "stuminh	w0, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x70, 0x20, 0x78}),
				address:          0,
			},
			want: "stuminh	w0, [x2]",
			wantErr: false,
		},
		{
			name: "stuminh	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x73, 0x22, 0x78}),
				address:          0,
			},
			want: "stuminh	w2, [sp]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_2a(t *testing.T) {
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
		// llvm/test/MC/AArch64/armv8.2a-at.s
		{
			name: "at	s1e1rp, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x79, 0x08, 0xd5}),
				address:          0,
			},
			want: "at	s1e1rp, x1",
			wantErr: false,
		},
		{
			name: "at	s1e1wp, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x22, 0x79, 0x08, 0xd5}),
				address:          0,
			},
			want: "at	s1e1wp, x2",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.2a-crypto-apple.s
		{
			name: "sha512h.2d	q0, q1, v2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x80, 0x62, 0xce}),
				address:          0,
			},
			want: "sha512h.2d	q0, q1, v2",
			wantErr: false,
		},
		{
			name: "sha512h2.2d	q0, q1, v2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x62, 0xce}),
				address:          0,
			},
			want: "sha512h2.2d	q0, q1, v2",
			wantErr: false,
		},
		{
			name: "sha512su0.2d	v11, v12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x81, 0xc0, 0xce}),
				address:          0,
			},
			want: "sha512su0.2d	v11, v12",
			wantErr: false,
		},
		{
			name: "sha512su1.2d	v11, v13, v14",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x89, 0x6e, 0xce}),
				address:          0,
			},
			want: "sha512su1.2d	v11, v13, v14",
			wantErr: false,
		},
		{
			name: "eor3.16b	v25, v12, v7, v2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x99, 0x09, 0x07, 0xce}),
				address:          0,
			},
			want: "eor3.16b	v25, v12, v7, v2",
			wantErr: false,
		},
		{
			name: "rax1.2d	v30, v29, v26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x8f, 0x7a, 0xce}),
				address:          0,
			},
			want: "rax1.2d	v30, v29, v26",
			wantErr: false,
		},
		{
			name: "xar.2d	v26, v21, v27, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xba, 0xfe, 0x9b, 0xce}),
				address:          0,
			},
			want: "xar.2d	v26, v21, v27, #63",
			wantErr: false,
		},
		{
			name: "bcax.16b	v31, v26, v2, v1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x07, 0x22, 0xce}),
				address:          0,
			},
			want: "bcax.16b	v31, v26, v2, v1",
			wantErr: false,
		},
		{
			name: "sm3ss1.4s	v20, v23, v21, v22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x5a, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3ss1.4s	v20, v23, v21, v22",
			wantErr: false,
		},
		{
			name: "sm3tt1a.4s	v20, v23, v21[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xb2, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt1a.4s	v20, v23, v21[3]",
			wantErr: false,
		},
		{
			name: "sm3tt1b.4s	v20, v23, v21[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xb6, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt1b.4s	v20, v23, v21[3]",
			wantErr: false,
		},
		{
			name: "sm3tt2a.4s	v20, v23, v21[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xba, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt2a.4s	v20, v23, v21[3]",
			wantErr: false,
		},
		{
			name: "sm3tt2b.4s	v20, v23, v21[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xbe, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt2b.4s	v20, v23, v21[3]",
			wantErr: false,
		},
		{
			name: "sm3partw1.4s	v30, v29, v26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0xc3, 0x7a, 0xce}),
				address:          0,
			},
			want: "sm3partw1.4s	v30, v29, v26",
			wantErr: false,
		},
		{
			name: "sm3partw2.4s	v30, v29, v26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0xc7, 0x7a, 0xce}),
				address:          0,
			},
			want: "sm3partw2.4s	v30, v29, v26",
			wantErr: false,
		},
		{
			name: "sm4ekey.4s	v11, v11, v19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6b, 0xc9, 0x73, 0xce}),
				address:          0,
			},
			want: "sm4ekey.4s	v11, v11, v19",
			wantErr: false,
		},
		{
			name: "sm4e.4s	v2, v15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x85, 0xc0, 0xce}),
				address:          0,
			},
			want: "sm4e.4s	v2, v15",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.2a-crypto.s
		{
			name: "sha512h	q0, q1, v2.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x80, 0x62, 0xce}),
				address:          0,
			},
			want: "sha512h	q0, q1, v2.2d",
			wantErr: false,
		},
		{
			name: "sha512h2	q0, q1, v2.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x84, 0x62, 0xce}),
				address:          0,
			},
			want: "sha512h2	q0, q1, v2.2d",
			wantErr: false,
		},
		{
			name: "sha512su0	v11.2d, v12.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x81, 0xc0, 0xce}),
				address:          0,
			},
			want: "sha512su0	v11.2d, v12.2d",
			wantErr: false,
		},
		{
			name: "sha512su1	v11.2d, v13.2d, v14.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x89, 0x6e, 0xce}),
				address:          0,
			},
			want: "sha512su1	v11.2d, v13.2d, v14.2d",
			wantErr: false,
		},
		{
			name: "eor3	v25.16b, v12.16b, v7.16b, v2.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x99, 0x09, 0x07, 0xce}),
				address:          0,
			},
			want: "eor3	v25.16b, v12.16b, v7.16b, v2.16b",
			wantErr: false,
		},
		{
			name: "rax1	v30.2d, v29.2d, v26.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x8f, 0x7a, 0xce}),
				address:          0,
			},
			want: "rax1	v30.2d, v29.2d, v26.2d",
			wantErr: false,
		},
		{
			name: "xar	v26.2d, v21.2d, v27.2d, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xba, 0xfe, 0x9b, 0xce}),
				address:          0,
			},
			want: "xar	v26.2d, v21.2d, v27.2d, #63",
			wantErr: false,
		},
		{
			name: "bcax	v31.16b, v26.16b, v2.16b, v1.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x07, 0x22, 0xce}),
				address:          0,
			},
			want: "bcax	v31.16b, v26.16b, v2.16b, v1.16b",
			wantErr: false,
		},
		{
			name: "sm3ss1	v20.4s, v23.4s, v21.4s, v22.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x5a, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3ss1	v20.4s, v23.4s, v21.4s, v22.4s",
			wantErr: false,
		},
		{
			name: "sm3tt1a	v20.4s, v23.4s, v21.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xb2, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt1a	v20.4s, v23.4s, v21.s[3]",
			wantErr: false,
		},
		{
			name: "sm3tt1b	v20.4s, v23.4s, v21.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xb6, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt1b	v20.4s, v23.4s, v21.s[3]",
			wantErr: false,
		},
		{
			name: "sm3tt2a	v20.4s, v23.4s, v21.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xba, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt2a	v20.4s, v23.4s, v21.s[3]",
			wantErr: false,
		},
		{
			name: "sm3tt2b	v20.4s, v23.4s, v21.s[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xbe, 0x55, 0xce}),
				address:          0,
			},
			want: "sm3tt2b	v20.4s, v23.4s, v21.s[3]",
			wantErr: false,
		},
		{
			name: "sm3partw1	v30.4s, v29.4s, v26.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0xc3, 0x7a, 0xce}),
				address:          0,
			},
			want: "sm3partw1	v30.4s, v29.4s, v26.4s",
			wantErr: false,
		},
		{
			name: "sm3partw2	v30.4s, v29.4s, v26.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0xc7, 0x7a, 0xce}),
				address:          0,
			},
			want: "sm3partw2	v30.4s, v29.4s, v26.4s",
			wantErr: false,
		},
		{
			name: "sm4ekey	v11.4s, v11.4s, v19.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6b, 0xc9, 0x73, 0xce}),
				address:          0,
			},
			want: "sm4ekey	v11.4s, v11.4s, v19.4s",
			wantErr: false,
		},
		{
			name: "sm4e	v2.4s, v15.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x85, 0xc0, 0xce}),
				address:          0,
			},
			want: "sm4e	v2.4s, v15.4s",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.2a-dotprod.s
		{
			name: "udot	v0.2s, v1.8b, v2.8b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x94, 0x82, 0x2e}),
				address:          0,
			},
			want: "udot	v0.2s, v1.8b, v2.8b",
			wantErr: false,
		},
		{
			name: "sdot	v0.2s, v1.8b, v2.8b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x94, 0x82, 0x0e}),
				address:          0,
			},
			want: "sdot	v0.2s, v1.8b, v2.8b",
			wantErr: false,
		},
		{
			name: "udot	v0.4s, v1.16b, v2.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x94, 0x82, 0x6e}),
				address:          0,
			},
			want: "udot	v0.4s, v1.16b, v2.16b",
			wantErr: false,
		},
		{
			name: "sdot	v0.4s, v1.16b, v2.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x94, 0x82, 0x4e}),
				address:          0,
			},
			want: "sdot	v0.4s, v1.16b, v2.16b",
			wantErr: false,
		},
		{
			name: "udot	v0.2s, v1.8b, v2.4b[0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe0, 0x82, 0x2f}),
				address:          0,
			},
			want: "udot	v0.2s, v1.8b, v2.4b[0]",
			wantErr: false,
		},
		{
			name: "sdot	v0.2s, v1.8b, v2.4b[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe0, 0xa2, 0x0f}),
				address:          0,
			},
			want: "sdot	v0.2s, v1.8b, v2.4b[1]",
			wantErr: false,
		},
		{
			name: "udot	v0.4s, v1.16b, v2.4b[2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe8, 0x82, 0x6f}),
				address:          0,
			},
			want: "udot	v0.4s, v1.16b, v2.4b[2]",
			wantErr: false,
		},
		{
			name: "sdot	v0.4s, v1.16b, v2.4b[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe8, 0xa2, 0x4f}),
				address:          0,
			},
			want: "sdot	v0.4s, v1.16b, v2.4b[3]",
			wantErr: false,
		},
		{
			name: "udot	v0.2s, v1.8b, v2.4b[0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe0, 0x82, 0x2f}),
				address:          0,
			},
			want: "udot	v0.2s, v1.8b, v2.4b[0]",
			wantErr: false,
		},
		{
			name: "udot	v0.4s, v1.16b, v2.4b[2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe8, 0x82, 0x6f}),
				address:          0,
			},
			want: "udot	v0.4s, v1.16b, v2.4b[2]",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.2a-persistent-memory.s
		{
			name: "dc	cvap, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x27, 0x7c, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cvap, x7",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.2a-statistical-profiling.s
		{
			name: "psb	csync",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x22, 0x03, 0xd5}),
				address:          0,
			},
			want: "psb	csync",
			wantErr: false,
		},
		{
			name: "msr	pmblimitr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x9a, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmblimitr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmbptr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x9a, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmbptr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmbsr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x9a, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmbsr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmscr_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x99, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	pmscr_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	pmscr_el12, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x99, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	pmscr_el12, x0",
			wantErr: false,
		},
		{
			name: "msr	pmscr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x99, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmscr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmsicr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x99, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmsicr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmsirr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x99, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmsirr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmsfcr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x99, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmsfcr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmsevfr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0x99, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmsevfr_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	pmslatfr_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0x99, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmslatfr_el1, x0",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmblimitr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x9a, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmblimitr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmbptr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x9a, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmbptr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmbsr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x9a, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmbsr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmbidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x9a, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmbidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmscr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x99, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmscr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmscr_el12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x99, 0x3d, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmscr_el12",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmscr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmscr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmsicr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmsicr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmsirr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmsirr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmsfcr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmsfcr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmsevfr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmsevfr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmslatfr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmslatfr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, pmsidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x99, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, pmsidr_el1",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.2a-uao.s
		{
			name: "msr	uao, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x40, 0x00, 0xd5}),
				address:          0,
			},
			want: "msr	uao, #0",
			wantErr: false,
		},
		{
			name: "msr	uao, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x41, 0x00, 0xd5}),
				address:          0,
			},
			want: "msr	uao, #1",
			wantErr: false,
		},
		{
			name: "msr	uao, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0x42, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	uao, x1",
			wantErr: false,
		},
		{
			name: "mrs	x2, uao",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x42, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x2, uao",
			wantErr: false,
		},
		//
		// llvm/test/MC/AArch64/armv8a-fpmul.s
		//
		{
			name: "fmlal	v0.2s, v1.2h, v2.2h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xec, 0x22, 0x0e}),
				address:          0,
			},
			want: "fmlal	v0.2s, v1.2h, v2.2h",
			wantErr: false,
		},
		{
			name: "fmlsl	v0.2s, v1.2h, v2.2h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xec, 0xa2, 0x0e}),
				address:          0,
			},
			want: "fmlsl	v0.2s, v1.2h, v2.2h",
			wantErr: false,
		},
		{
			name: "fmlal	v0.4s, v1.4h, v2.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xec, 0x22, 0x4e}),
				address:          0,
			},
			want: "fmlal	v0.4s, v1.4h, v2.4h",
			wantErr: false,
		},
		{
			name: "fmlsl	v0.4s, v1.4h, v2.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xec, 0xa2, 0x4e}),
				address:          0,
			},
			want: "fmlsl	v0.4s, v1.4h, v2.4h",
			wantErr: false,
		},
		{
			name: "fmlal2	v0.2s, v1.2h, v2.2h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xcc, 0x22, 0x2e}),
				address:          0,
			},
			want: "fmlal2	v0.2s, v1.2h, v2.2h",
			wantErr: false,
		},
		{
			name: "fmlsl2	v0.2s, v1.2h, v2.2h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xcc, 0xa2, 0x2e}),
				address:          0,
			},
			want: "fmlsl2	v0.2s, v1.2h, v2.2h",
			wantErr: false,
		},
		{
			name: "fmlal2	v0.4s, v1.4h, v2.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xcc, 0x22, 0x6e}),
				address:          0,
			},
			want: "fmlal2	v0.4s, v1.4h, v2.4h",
			wantErr: false,
		},
		{
			name: "fmlsl2	v0.4s, v1.4h, v2.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xcc, 0xa2, 0x6e}),
				address:          0,
			},
			want: "fmlsl2	v0.4s, v1.4h, v2.4h",
			wantErr: false,
		},
		{
			name: "fmlal	v0.2s, v1.2h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0xb2, 0x0f}),
				address:          0,
			},
			want: "fmlal	v0.2s, v1.2h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlsl	v0.2s, v1.2h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x48, 0xb2, 0x0f}),
				address:          0,
			},
			want: "fmlsl	v0.2s, v1.2h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlal	v0.4s, v1.4h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0xb2, 0x4f}),
				address:          0,
			},
			want: "fmlal	v0.4s, v1.4h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlsl	v0.4s, v1.4h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x48, 0xb2, 0x4f}),
				address:          0,
			},
			want: "fmlsl	v0.4s, v1.4h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlal2	v0.2s, v1.2h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x88, 0xb2, 0x2f}),
				address:          0,
			},
			want: "fmlal2	v0.2s, v1.2h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlsl2	v0.2s, v1.2h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc8, 0xb2, 0x2f}),
				address:          0,
			},
			want: "fmlsl2	v0.2s, v1.2h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlal2	v0.4s, v1.4h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x88, 0xb2, 0x6f}),
				address:          0,
			},
			want: "fmlal2	v0.4s, v1.4h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlsl2	v0.4s, v1.4h, v2.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc8, 0xb2, 0x6f}),
				address:          0,
			},
			want: "fmlsl2	v0.4s, v1.4h, v2.h[7]",
			wantErr: false,
		},
		{
			name: "fmlal	v0.2s, v1.2h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0x92, 0x0f}),
				address:          0,
			},
			want: "fmlal	v0.2s, v1.2h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlsl	v0.2s, v1.2h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x48, 0x92, 0x0f}),
				address:          0,
			},
			want: "fmlsl	v0.2s, v1.2h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlal	v0.4s, v1.4h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0x92, 0x4f}),
				address:          0,
			},
			want: "fmlal	v0.4s, v1.4h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlsl	v0.4s, v1.4h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x48, 0x92, 0x4f}),
				address:          0,
			},
			want: "fmlsl	v0.4s, v1.4h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlal2	v0.2s, v1.2h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x88, 0x92, 0x2f}),
				address:          0,
			},
			want: "fmlal2	v0.2s, v1.2h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlsl2	v0.2s, v1.2h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc8, 0x92, 0x2f}),
				address:          0,
			},
			want: "fmlsl2	v0.2s, v1.2h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlal2	v0.4s, v1.4h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x88, 0x92, 0x6f}),
				address:          0,
			},
			want: "fmlal2	v0.4s, v1.4h, v2.h[5]",
			wantErr: false,
		},
		{
			name: "fmlsl2	v0.4s, v1.4h, v2.h[5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc8, 0x92, 0x6f}),
				address:          0,
			},
			want: "fmlsl2	v0.4s, v1.4h, v2.h[5]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_3a(t *testing.T) {
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
		//
		// llvm/test/MC/AArch64/armv8.3a-complex_nofp16.s
		//
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.4s, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x82, 0x6e}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.4s, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2d, v1.2d, v2.2d, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0xc2, 0x6e}),
				address:          0,
			},
			want: "fcmla	v0.2d, v1.2d, v2.2d, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xcc, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #90",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #180",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #180",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #270",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xdc, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #270",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.4s, v1.4s, v2.4s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x82, 0x6e}),
				address:          0,
			},
			want: "fcadd	v0.4s, v1.4s, v2.4s, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2d, v1.2d, v2.2d, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0xc2, 0x6e}),
				address:          0,
			},
			want: "fcadd	v0.2d, v1.2d, v2.2d, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2s, v1.2s, v2.2s, #270",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.2s, v1.2s, v2.2s, #270",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x30, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #90",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #180",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x50, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #180",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #270",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x70, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #270",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[1], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x18, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[1], #0",
			wantErr: false,
		},

		//
		// llvm/test/MC/AArch64/armv8.3a-complex.s
		//
		{
			name: "fcmla	v0.4h, v1.4h, v2.4h, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x42, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.4h, v1.4h, v2.4h, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.8h, v1.8h, v2.8h, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x42, 0x6e}),
				address:          0,
			},
			want: "fcmla	v0.8h, v1.8h, v2.8h, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.4s, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x82, 0x6e}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.4s, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2d, v1.2d, v2.2d, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0xc2, 0x6e}),
				address:          0,
			},
			want: "fcmla	v0.2d, v1.2d, v2.2d, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xcc, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #90",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #180",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #180",
			wantErr: false,
		},
		{
			name: "fcmla	v0.2s, v1.2s, v2.2s, #270",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xdc, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcmla	v0.2s, v1.2s, v2.2s, #270",
			wantErr: false,
		},
		{
			name: "fcadd	v0.4h, v1.4h, v2.4h, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x42, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.4h, v1.4h, v2.4h, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.8h, v1.8h, v2.8h, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x42, 0x6e}),
				address:          0,
			},
			want: "fcadd	v0.8h, v1.8h, v2.8h, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.4s, v1.4s, v2.4s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x82, 0x6e}),
				address:          0,
			},
			want: "fcadd	v0.4s, v1.4s, v2.4s, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2d, v1.2d, v2.2d, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0xc2, 0x6e}),
				address:          0,
			},
			want: "fcadd	v0.2d, v1.2d, v2.2d, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.2s, v1.2s, v2.2s, #90",
			wantErr: false,
		},
		{
			name: "fcadd	v0.2s, v1.2s, v2.2s, #270",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf4, 0x82, 0x2e}),
				address:          0,
			},
			want: "fcadd	v0.2s, v1.2s, v2.2s, #270",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4h, v1.4h, v2.h[0], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0x42, 0x2f}),
				address:          0,
			},
			want: "fcmla	v0.4h, v1.4h, v2.h[0], #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.8h, v1.8h, v2.h[0], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0x42, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.8h, v1.8h, v2.h[0], #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #90",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x30, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #90",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #180",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x50, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #180",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[0], #270",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x70, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[0], #270",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4h, v1.4h, v2.h[1], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0x62, 0x2f}),
				address:          0,
			},
			want: "fcmla	v0.4h, v1.4h, v2.h[1], #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.8h, v1.8h, v2.h[3], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x18, 0x62, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.8h, v1.8h, v2.h[3], #0",
			wantErr: false,
		},
		{
			name: "fcmla	v0.4s, v1.4s, v2.s[1], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x18, 0x82, 0x6f}),
				address:          0,
			},
			want: "fcmla	v0.4s, v1.4s, v2.s[1], #0",
			wantErr: false,
		},

		//
		// llvm/test/MC/AArch64/armv8.3a-ID_ISAR6_EL1.s
		//
		{
			name: "mrs	x0, id_isar6_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, id_isar6_el1",
			wantErr: false,
		},

		//
		// llvm/test/MC/AArch64/armv8.3a-js.s
		//
		{
			name: "fjcvtzs	w0, d0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x7e, 0x1e}),
				address:          0,
			},
			want: "fjcvtzs	w0, d0",
			wantErr: false,
		},

		// llvm/test/MC/AArch64/armv8.3a-rcpc.s
		{
			name: "ldaprb	w0, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xc0, 0xbf, 0x38}),
				address:          0,
			},
			want: "ldaprb	w0, [x0]",
			wantErr: false,
		},
		{
			name: "ldaprh	w0, [x17]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc2, 0xbf, 0x78}),
				address:          0,
			},
			want: "ldaprh	w0, [x17]",
			wantErr: false,
		},
		{
			name: "ldapr	w0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xc0, 0xbf, 0xb8}),
				address:          0,
			},
			want: "ldapr	w0, [x1]",
			wantErr: false,
		},
		{
			name: "ldapr	x0, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xc0, 0xbf, 0xf8}),
				address:          0,
			},
			want: "ldapr	x0, [x0]",
			wantErr: false,
		},
		{
			name: "ldapr	w18, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x12, 0xc0, 0xbf, 0xb8}),
				address:          0,
			},
			want: "ldapr	w18, [x0]",
			wantErr: false,
		},
		{
			name: "ldapr	x15, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0f, 0xc0, 0xbf, 0xf8}),
				address:          0,
			},
			want: "ldapr	x15, [x0]",
			wantErr: false,
		},

		// llvm/test/MC/AArch64/armv8.3a-signed-pointer.s
		{
			name: "mrs	x0, apiakeylo_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x21, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apiakeylo_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apiakeyhi_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x21, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apiakeyhi_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apibkeylo_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x21, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apibkeylo_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apibkeyhi_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x21, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apibkeyhi_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apdakeylo_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x22, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apdakeylo_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apdakeyhi_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x22, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apdakeyhi_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apdbkeylo_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x22, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apdbkeylo_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apdbkeyhi_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x22, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apdbkeyhi_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apgakeylo_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x23, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apgakeylo_el1",
			wantErr: false,
		},
		{
			name: "mrs	x0, apgakeyhi_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x23, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, apgakeyhi_el1",
			wantErr: false,
		},
		{
			name: "msr	apiakeylo_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x21, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apiakeylo_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apiakeyhi_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x21, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apiakeyhi_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apibkeylo_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x21, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apibkeylo_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apibkeyhi_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x21, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apibkeyhi_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apdakeylo_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x22, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apdakeylo_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apdakeyhi_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x22, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apdakeyhi_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apdbkeylo_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x22, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apdbkeylo_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apdbkeyhi_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x22, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apdbkeyhi_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apgakeylo_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x23, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apgakeylo_el1, x0",
			wantErr: false,
		},
		{
			name: "msr	apgakeyhi_el1, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x23, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	apgakeyhi_el1, x0",
			wantErr: false,
		},
		{
			name: "paciasp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "paciasp",
			wantErr: false,
		},
		{
			name: "autiasp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "autiasp",
			wantErr: false,
		},
		{
			name: "paciaz",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "paciaz",
			wantErr: false,
		},
		{
			name: "autiaz",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "autiaz",
			wantErr: false,
		},
		{
			name: "pacia1716",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x21, 0x03, 0xd5}),
				address:          0,
			},
			want:    "pacia1716",
			wantErr: false,
		},
		{
			name: "autia1716",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x21, 0x03, 0xd5}),
				address:          0,
			},
			want:    "autia1716",
			wantErr: false,
		},
		{
			name: "pacibsp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "pacibsp",
			wantErr: false,
		},
		{
			name: "autibsp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "autibsp",
			wantErr: false,
		},
		{
			name: "pacibz",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "pacibz",
			wantErr: false,
		},
		{
			name: "autibz",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x23, 0x03, 0xd5}),
				address:          0,
			},
			want:    "autibz",
			wantErr: false,
		},
		{
			name: "pacib1716",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x21, 0x03, 0xd5}),
				address:          0,
			},
			want:    "pacib1716",
			wantErr: false,
		},
		{
			name: "autib1716",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x21, 0x03, 0xd5}),
				address:          0,
			},
			want:    "autib1716",
			wantErr: false,
		},
		{
			name: "xpaclri",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "xpaclri",
			wantErr: false,
		},
		{
			name: "pacia	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacia	x0, x1",
			wantErr: false,
		},
		{
			name: "autia	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0xc1, 0xda}),
				address:          0,
			},
			want: "autia	x0, x1",
			wantErr: false,
		},
		{
			name: "pacda	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacda	x0, x1",
			wantErr: false,
		},
		{
			name: "autda	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x18, 0xc1, 0xda}),
				address:          0,
			},
			want: "autda	x0, x1",
			wantErr: false,
		},
		{
			name: "pacib	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacib	x0, x1",
			wantErr: false,
		},
		{
			name: "autib	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x14, 0xc1, 0xda}),
				address:          0,
			},
			want: "autib	x0, x1",
			wantErr: false,
		},
		{
			name: "pacdb	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacdb	x0, x1",
			wantErr: false,
		},
		{
			name: "autdb	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x1c, 0xc1, 0xda}),
				address:          0,
			},
			want: "autdb	x0, x1",
			wantErr: false,
		},
		{
			name: "pacga	x0, x1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x30, 0xc2, 0x9a}),
				address:          0,
			},
			want: "pacga	x0, x1, x2",
			wantErr: false,
		},
		{
			name: "paciza	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x23, 0xc1, 0xda}),
				address:          0,
			},
			want: "paciza	x0",
			wantErr: false,
		},
		{
			name: "autiza	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x33, 0xc1, 0xda}),
				address:          0,
			},
			want: "autiza	x0",
			wantErr: false,
		},
		{
			name: "pacdza	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x2b, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacdza	x0",
			wantErr: false,
		},
		{
			name: "autdza	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x3b, 0xc1, 0xda}),
				address:          0,
			},
			want: "autdza	x0",
			wantErr: false,
		},
		{
			name: "pacizb	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x27, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacizb	x0",
			wantErr: false,
		},
		{
			name: "autizb	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x37, 0xc1, 0xda}),
				address:          0,
			},
			want: "autizb	x0",
			wantErr: false,
		},
		{
			name: "pacdzb	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x2f, 0xc1, 0xda}),
				address:          0,
			},
			want: "pacdzb	x0",
			wantErr: false,
		},
		{
			name: "autdzb	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x3f, 0xc1, 0xda}),
				address:          0,
			},
			want: "autdzb	x0",
			wantErr: false,
		},
		{
			name: "xpaci	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x43, 0xc1, 0xda}),
				address:          0,
			},
			want: "xpaci	x0",
			wantErr: false,
		},
		{
			name: "xpacd	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x47, 0xc1, 0xda}),
				address:          0,
			},
			want: "xpacd	x0",
			wantErr: false,
		},
		{
			name: "braa	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x08, 0x1f, 0xd7}),
				address:          0,
			},
			want: "braa	x0, x1",
			wantErr: false,
		},
		{
			name: "brab	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x0c, 0x1f, 0xd7}),
				address:          0,
			},
			want: "brab	x0, x1",
			wantErr: false,
		},
		{
			name: "blraa	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x08, 0x3f, 0xd7}),
				address:          0,
			},
			want: "blraa	x0, x1",
			wantErr: false,
		},
		{
			name: "blrab	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x0c, 0x3f, 0xd7}),
				address:          0,
			},
			want: "blrab	x0, x1",
			wantErr: false,
		},
		{
			name: "braaz	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x08, 0x1f, 0xd6}),
				address:          0,
			},
			want: "braaz	x0",
			wantErr: false,
		},
		{
			name: "brabz	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x0c, 0x1f, 0xd6}),
				address:          0,
			},
			want: "brabz	x0",
			wantErr: false,
		},
		{
			name: "blraaz	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x08, 0x3f, 0xd6}),
				address:          0,
			},
			want: "blraaz	x0",
			wantErr: false,
		},
		{
			name: "blrabz	x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x0c, 0x3f, 0xd6}),
				address:          0,
			},
			want: "blrabz	x0",
			wantErr: false,
		},
		{
			name: "retaa",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0b, 0x5f, 0xd6}),
				address:          0,
			},
			want:    "retaa",
			wantErr: false,
		},
		{
			name: "retab",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0f, 0x5f, 0xd6}),
				address:          0,
			},
			want:    "retab",
			wantErr: false,
		},
		{
			name: "eretaa",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0b, 0x9f, 0xd6}),
				address:          0,
			},
			want:    "eretaa",
			wantErr: false,
		},
		{
			name: "eretab",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0f, 0x9f, 0xd6}),
				address:          0,
			},
			want:    "eretab",
			wantErr: false,
		},
		{
			name: "ldraa	x0, [x1, #4088]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf4, 0x3f, 0xf8}),
				address:          0,
			},
			want: "ldraa	x0, [x1, #4088]",
			wantErr: false,
		},
		{
			name: "ldraa	x0, [x1, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldraa	x0, [x1, #-4096]",
			wantErr: false,
		},
		{
			name: "ldrab	x0, [x1, #4088]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xf4, 0xbf, 0xf8}),
				address:          0,
			},
			want: "ldrab	x0, [x1, #4088]",
			wantErr: false,
		},
		{
			name: "ldrab	x0, [x1, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldrab	x0, [x1, #-4096]",
			wantErr: false,
		},
		{
			name: "ldraa	x0, [x1, #4088]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xfc, 0x3f, 0xf8}),
				address:          0,
			},
			want: "ldraa	x0, [x1, #4088]!",
			wantErr: false,
		},
		{
			name: "ldraa	x0, [x1, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldraa	x0, [x1, #-4096]!",
			wantErr: false,
		},
		{
			name: "ldrab	x0, [x1, #4088]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xfc, 0xbf, 0xf8}),
				address:          0,
			},
			want: "ldrab	x0, [x1, #4088]!",
			wantErr: false,
		},
		{
			name: "ldrab	x0, [x1, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldrab	x0, [x1, #-4096]!",
			wantErr: false,
		},
		{
			name: "ldraa	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldraa	x0, [x1]",
			wantErr: false,
		},
		{
			name: "ldrab	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldrab	x0, [x1]",
			wantErr: false,
		},
		{
			name: "ldraa	x0, [x1, #0]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0x20, 0xf8}),
				address:          0,
			},
			want: "ldraa	x0, [x1, #0]!",
			wantErr: false,
		},
		{
			name: "ldrab	x0, [x1, #0]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0xa0, 0xf8}),
				address:          0,
			},
			want: "ldrab	x0, [x1, #0]!",
			wantErr: false,
		},
		{
			name: "ldraa	xzr, [sp, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0f, 0x60, 0xf8}),
				address:          0,
			},
			want: "ldraa	xzr, [sp, #-4096]!",
			wantErr: false,
		},
		{
			name: "ldrab	xzr, [sp, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0f, 0xe0, 0xf8}),
				address:          0,
			},
			want: "ldrab	xzr, [sp, #-4096]!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_4a(t *testing.T) {
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
		// llvm/test/MC/AArch64/armv8.4a-ldst.s
		{
			name: "stlurb	w1, [x10]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x01, 0x00, 0x19}),
				address:          0,
			},
			want: "stlurb	w1, [x10]",
			wantErr: false,
		},
		{
			name: "stlurb	w1, [x10, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x01, 0x10, 0x19}),
				address:          0,
			},
			want: "stlurb	w1, [x10, #-256]",
			wantErr: false,
		},
		{
			name: "stlurb	w2, [x11, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf1, 0x0f, 0x19}),
				address:          0,
			},
			want: "stlurb	w2, [x11, #255]",
			wantErr: false,
		},
		{
			name: "stlurb	w3, [sp, #-3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xd3, 0x1f, 0x19}),
				address:          0,
			},
			want: "stlurb	w3, [sp, #-3]",
			wantErr: false,
		},
		{
			name: "ldapurb	wzr, [x12]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x40, 0x19}),
				address:          0,
			},
			want: "ldapurb	wzr, [x12]",
			wantErr: false,
		},
		{
			name: "ldapurb	w4, [x12]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x84, 0x01, 0x40, 0x19}),
				address:          0,
			},
			want: "ldapurb	w4, [x12]",
			wantErr: false,
		},
		{
			name: "ldapurb	w4, [x12, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x84, 0x01, 0x50, 0x19}),
				address:          0,
			},
			want: "ldapurb	w4, [x12, #-256]",
			wantErr: false,
		},
		{
			name: "ldapurb	w5, [x13, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa5, 0xf1, 0x4f, 0x19}),
				address:          0,
			},
			want: "ldapurb	w5, [x13, #255]",
			wantErr: false,
		},
		{
			name: "ldapurb	w6, [sp, #-2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xe3, 0x5f, 0x19}),
				address:          0,
			},
			want: "ldapurb	w6, [sp, #-2]",
			wantErr: false,
		},
		{
			name: "ldapursb	w7, [x14]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc7, 0x01, 0xc0, 0x19}),
				address:          0,
			},
			want: "ldapursb	w7, [x14]",
			wantErr: false,
		},
		{
			name: "ldapursb	w7, [x14, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc7, 0x01, 0xd0, 0x19}),
				address:          0,
			},
			want: "ldapursb	w7, [x14, #-256]",
			wantErr: false,
		},
		{
			name: "ldapursb	w8, [x15, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0xf1, 0xcf, 0x19}),
				address:          0,
			},
			want: "ldapursb	w8, [x15, #255]",
			wantErr: false,
		},
		{
			name: "ldapursb	w9, [sp, #-1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xf3, 0xdf, 0x19}),
				address:          0,
			},
			want: "ldapursb	w9, [sp, #-1]",
			wantErr: false,
		},
		{
			name: "ldapursb	x0, [x16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x02, 0x80, 0x19}),
				address:          0,
			},
			want: "ldapursb	x0, [x16]",
			wantErr: false,
		},
		{
			name: "ldapursb	x0, [x16, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x02, 0x90, 0x19}),
				address:          0,
			},
			want: "ldapursb	x0, [x16, #-256]",
			wantErr: false,
		},
		{
			name: "ldapursb	x1, [x17, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0xf2, 0x8f, 0x19}),
				address:          0,
			},
			want: "ldapursb	x1, [x17, #255]",
			wantErr: false,
		},
		{
			name: "ldapursb	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x03, 0x80, 0x19}),
				address:          0,
			},
			want: "ldapursb	x2, [sp]",
			wantErr: false,
		},
		{
			name: "ldapursb	x2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x03, 0x80, 0x19}),
				address:          0,
			},
			want: "ldapursb	x2, [sp]",
			wantErr: false,
		},
		{
			name: "stlurh	w10, [x18]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0x02, 0x00, 0x59}),
				address:          0,
			},
			want: "stlurh	w10, [x18]",
			wantErr: false,
		},
		{
			name: "stlurh	w10, [x18, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0x02, 0x10, 0x59}),
				address:          0,
			},
			want: "stlurh	w10, [x18, #-256]",
			wantErr: false,
		},
		{
			name: "stlurh	w11, [x19, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6b, 0xf2, 0x0f, 0x59}),
				address:          0,
			},
			want: "stlurh	w11, [x19, #255]",
			wantErr: false,
		},
		{
			name: "stlurh	w12, [sp, #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x13, 0x00, 0x59}),
				address:          0,
			},
			want: "stlurh	w12, [sp, #1]",
			wantErr: false,
		},
		{
			name: "ldapurh	w13, [x20]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8d, 0x02, 0x40, 0x59}),
				address:          0,
			},
			want: "ldapurh	w13, [x20]",
			wantErr: false,
		},
		{
			name: "ldapurh	w13, [x20, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8d, 0x02, 0x50, 0x59}),
				address:          0,
			},
			want: "ldapurh	w13, [x20, #-256]",
			wantErr: false,
		},
		{
			name: "ldapurh	w14, [x21, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xae, 0xf2, 0x4f, 0x59}),
				address:          0,
			},
			want: "ldapurh	w14, [x21, #255]",
			wantErr: false,
		},
		{
			name: "ldapurh	w15, [sp, #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0x23, 0x40, 0x59}),
				address:          0,
			},
			want: "ldapurh	w15, [sp, #2]",
			wantErr: false,
		},
		{
			name: "ldapursh	w16, [x22]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd0, 0x02, 0xc0, 0x59}),
				address:          0,
			},
			want: "ldapursh	w16, [x22]",
			wantErr: false,
		},
		{
			name: "ldapursh	w16, [x22, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd0, 0x02, 0xd0, 0x59}),
				address:          0,
			},
			want: "ldapursh	w16, [x22, #-256]",
			wantErr: false,
		},
		{
			name: "ldapursh	w17, [x23, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xf2, 0xcf, 0x59}),
				address:          0,
			},
			want: "ldapursh	w17, [x23, #255]",
			wantErr: false,
		},
		{
			name: "ldapursh	w18, [sp, #3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0x33, 0xc0, 0x59}),
				address:          0,
			},
			want: "ldapursh	w18, [sp, #3]",
			wantErr: false,
		},
		{
			name: "ldapursh	x3, [x24]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x03, 0x80, 0x59}),
				address:          0,
			},
			want: "ldapursh	x3, [x24]",
			wantErr: false,
		},
		{
			name: "ldapursh	x3, [x24, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x03, 0x90, 0x59}),
				address:          0,
			},
			want: "ldapursh	x3, [x24, #-256]",
			wantErr: false,
		},
		{
			name: "ldapursh	x4, [x25, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x24, 0xf3, 0x8f, 0x59}),
				address:          0,
			},
			want: "ldapursh	x4, [x25, #255]",
			wantErr: false,
		},
		{
			name: "ldapursh	x5, [sp, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x43, 0x80, 0x59}),
				address:          0,
			},
			want: "ldapursh	x5, [sp, #4]",
			wantErr: false,
		},
		{
			name: "stlur	w19, [x26]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x53, 0x03, 0x00, 0x99}),
				address:          0,
			},
			want: "stlur	w19, [x26]",
			wantErr: false,
		},
		{
			name: "stlur	w19, [x26, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x53, 0x03, 0x10, 0x99}),
				address:          0,
			},
			want: "stlur	w19, [x26, #-256]",
			wantErr: false,
		},
		{
			name: "stlur	w20, [x27, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0xf3, 0x0f, 0x99}),
				address:          0,
			},
			want: "stlur	w20, [x27, #255]",
			wantErr: false,
		},
		{
			name: "stlur	w21, [sp, #5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf5, 0x53, 0x00, 0x99}),
				address:          0,
			},
			want: "stlur	w21, [sp, #5]",
			wantErr: false,
		},
		{
			name: "ldapur	w22, [x28]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x96, 0x03, 0x40, 0x99}),
				address:          0,
			},
			want: "ldapur	w22, [x28]",
			wantErr: false,
		},
		{
			name: "ldapur	w22, [x28, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x96, 0x03, 0x50, 0x99}),
				address:          0,
			},
			want: "ldapur	w22, [x28, #-256]",
			wantErr: false,
		},
		{
			name: "ldapur	w23, [x29, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb7, 0xf3, 0x4f, 0x99}),
				address:          0,
			},
			want: "ldapur	w23, [x29, #255]",
			wantErr: false,
		},
		{
			name: "ldapur	w24, [sp, #6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x63, 0x40, 0x99}),
				address:          0,
			},
			want: "ldapur	w24, [sp, #6]",
			wantErr: false,
		},
		{
			name: "ldapursw	x6, [x30]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc6, 0x03, 0x80, 0x99}),
				address:          0,
			},
			want: "ldapursw	x6, [x30]",
			wantErr: false,
		},
		{
			name: "ldapursw	x6, [x30, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc6, 0x03, 0x90, 0x99}),
				address:          0,
			},
			want: "ldapursw	x6, [x30, #-256]",
			wantErr: false,
		},
		{
			name: "ldapursw	x7, [x0, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0xf0, 0x8f, 0x99}),
				address:          0,
			},
			want: "ldapursw	x7, [x0, #255]",
			wantErr: false,
		},
		{
			name: "ldapursw	x8, [sp, #7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0x73, 0x80, 0x99}),
				address:          0,
			},
			want: "ldapursw	x8, [sp, #7]",
			wantErr: false,
		},
		{
			name: "stlur	x9, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x00, 0x00, 0xd9}),
				address:          0,
			},
			want: "stlur	x9, [x1]",
			wantErr: false,
		},
		{
			name: "stlur	x9, [x1, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x00, 0x10, 0xd9}),
				address:          0,
			},
			want: "stlur	x9, [x1, #-256]",
			wantErr: false,
		},
		{
			name: "stlur	x10, [x2, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0xf0, 0x0f, 0xd9}),
				address:          0,
			},
			want: "stlur	x10, [x2, #255]",
			wantErr: false,
		},
		{
			name: "stlur	x11, [sp, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xeb, 0x83, 0x00, 0xd9}),
				address:          0,
			},
			want: "stlur	x11, [sp, #8]",
			wantErr: false,
		},
		{
			name: "ldapur	x12, [x3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0x00, 0x40, 0xd9}),
				address:          0,
			},
			want: "ldapur	x12, [x3]",
			wantErr: false,
		},
		{
			name: "ldapur	x12, [x3, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0x00, 0x50, 0xd9}),
				address:          0,
			},
			want: "ldapur	x12, [x3, #-256]",
			wantErr: false,
		},
		{
			name: "ldapur	x13, [x4, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8d, 0xf0, 0x4f, 0xd9}),
				address:          0,
			},
			want: "ldapur	x13, [x4, #255]",
			wantErr: false,
		},
		{
			name: "ldapur	x14, [sp, #9]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x93, 0x40, 0xd9}),
				address:          0,
			},
			want: "ldapur	x14, [sp, #9]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_5a(t *testing.T) {
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
		// // llvm/test/MC/AArch64/armv8.5a-altnzcv.s
		// {
		// 	name: "xaflag",
		// 	args: args{
		// 		instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x40, 0x00, 0xd5}),
		// 		address:          0,
		// 	},
		// 	want:    "xaflag",
		// 	wantErr: false,
		// },
		// {
		// 	name: "axflag",
		// 	args: args{
		// 		instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x40, 0x00, 0xd5}),
		// 		address:          0,
		// 	},
		// 	want:    "axflag",
		// 	wantErr: false,
		// },
		// llvm/test/MC/AArch64/armv8.5a-bti.s
		{
			name: "bti",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x24, 0x03, 0xd5}),
				address:          0,
			},
			want:    "bti",
			wantErr: false,
		},
		{
			name: "bti	c",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x24, 0x03, 0xd5}),
				address:          0,
			},
			want: "bti	c",
			wantErr: false,
		},
		{
			name: "bti	j",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x24, 0x03, 0xd5}),
				address:          0,
			},
			want: "bti	j",
			wantErr: false,
		},
		{
			name: "bti	jc",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x24, 0x03, 0xd5}),
				address:          0,
			},
			want: "bti	jc",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-frint.s
		{
			name: "frint32z	s0, s1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x40, 0x28, 0x1e}),
				address:          0,
			},
			want: "frint32z	s0, s1",
			wantErr: false,
		},
		{
			name: "frint32z	d0, d1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x40, 0x68, 0x1e}),
				address:          0,
			},
			want: "frint32z	d0, d1",
			wantErr: false,
		},
		{
			name: "frint64z	s2, s3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x40, 0x29, 0x1e}),
				address:          0,
			},
			want: "frint64z	s2, s3",
			wantErr: false,
		},
		{
			name: "frint64z	d2, d3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x40, 0x69, 0x1e}),
				address:          0,
			},
			want: "frint64z	d2, d3",
			wantErr: false,
		},
		{
			name: "frint32x	s4, s5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xc0, 0x28, 0x1e}),
				address:          0,
			},
			want: "frint32x	s4, s5",
			wantErr: false,
		},
		{
			name: "frint32x	d4, d5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xc0, 0x68, 0x1e}),
				address:          0,
			},
			want: "frint32x	d4, d5",
			wantErr: false,
		},
		{
			name: "frint64x	s6, s7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xc0, 0x29, 0x1e}),
				address:          0,
			},
			want: "frint64x	s6, s7",
			wantErr: false,
		},
		{
			name: "frint64x	d6, d7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xc0, 0x69, 0x1e}),
				address:          0,
			},
			want: "frint64x	d6, d7",
			wantErr: false,
		},
		{
			name: "frint32z	v0.2s, v1.2s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe8, 0x21, 0x0e}),
				address:          0,
			},
			want: "frint32z	v0.2s, v1.2s",
			wantErr: false,
		},
		{
			name: "frint32z	v0.2d, v1.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe8, 0x61, 0x4e}),
				address:          0,
			},
			want: "frint32z	v0.2d, v1.2d",
			wantErr: false,
		},
		{
			name: "frint32z	v0.4s, v1.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xe8, 0x21, 0x4e}),
				address:          0,
			},
			want: "frint32z	v0.4s, v1.4s",
			wantErr: false,
		},
		{
			name: "frint64z	v2.2s, v3.2s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x21, 0x0e}),
				address:          0,
			},
			want: "frint64z	v2.2s, v3.2s",
			wantErr: false,
		},
		{
			name: "frint64z	v2.2d, v3.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x61, 0x4e}),
				address:          0,
			},
			want: "frint64z	v2.2d, v3.2d",
			wantErr: false,
		},
		{
			name: "frint64z	v2.4s, v3.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x21, 0x4e}),
				address:          0,
			},
			want: "frint64z	v2.4s, v3.4s",
			wantErr: false,
		},
		{
			name: "frint32x	v4.2s, v5.2s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xe8, 0x21, 0x2e}),
				address:          0,
			},
			want: "frint32x	v4.2s, v5.2s",
			wantErr: false,
		},
		{
			name: "frint32x	v4.2d, v5.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xe8, 0x61, 0x6e}),
				address:          0,
			},
			want: "frint32x	v4.2d, v5.2d",
			wantErr: false,
		},
		{
			name: "frint32x	v4.4s, v5.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xe8, 0x21, 0x6e}),
				address:          0,
			},
			want: "frint32x	v4.4s, v5.4s",
			wantErr: false,
		},
		{
			name: "frint64x	v6.2s, v7.2s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xf8, 0x21, 0x2e}),
				address:          0,
			},
			want: "frint64x	v6.2s, v7.2s",
			wantErr: false,
		},
		{
			name: "frint64x	v6.2d, v7.2d",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xf8, 0x61, 0x6e}),
				address:          0,
			},
			want: "frint64x	v6.2d, v7.2d",
			wantErr: false,
		},
		{
			name: "frint64x	v6.4s, v7.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xf8, 0x21, 0x6e}),
				address:          0,
			},
			want: "frint64x	v6.4s, v7.4s",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-persistent-memory.s
		{
			name: "dc	cvadp, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x27, 0x7d, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cvadp, x7",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-predres.s
		{
			name: "cfp	rctx, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x73, 0x0b, 0xd5}),
				address:          0,
			},
			want: "cfp	rctx, x0",
			wantErr: false,
		},
		{
			name: "dvp	rctx, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa1, 0x73, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dvp	rctx, x1",
			wantErr: false,
		},
		{
			name: "cpp	rctx, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x73, 0x0b, 0xd5}),
				address:          0,
			},
			want: "cpp	rctx, x2",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-rand.s
		{
			name: "mrs	x0, RNDR",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x24, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, RNDR",
			wantErr: false,
		},
		{
			name: "mrs	x1, RNDRRS",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0x24, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x1, RNDRRS",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-sb.s
		{
			name: "sb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x30, 0x03, 0xd5}),
				address:          0,
			},
			want:    "sb",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-specrestrict.s
		{
			name: "mrs	x9, id_pfr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x03, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_pfr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x8, scxtnum_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0xd0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x8, scxtnum_el0",
			wantErr: false,
		},
		{
			name: "mrs	x7, scxtnum_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0xd0, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x7, scxtnum_el1",
			wantErr: false,
		},
		{
			name: "mrs	x6, scxtnum_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xd0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x6, scxtnum_el2",
			wantErr: false,
		},
		{
			name: "mrs	x5, scxtnum_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0xd0, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x5, scxtnum_el3",
			wantErr: false,
		},
		{
			name: "mrs	x4, scxtnum_el12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe4, 0xd0, 0x3d, 0xd5}),
				address:          0,
			},
			want: "mrs	x4, scxtnum_el12",
			wantErr: false,
		},
		{
			name: "msr	scxtnum_el0, x8",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0xd0, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	scxtnum_el0, x8",
			wantErr: false,
		},
		{
			name: "msr	scxtnum_el1, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0xd0, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	scxtnum_el1, x7",
			wantErr: false,
		},
		{
			name: "msr	scxtnum_el2, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xd0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	scxtnum_el2, x6",
			wantErr: false,
		},
		{
			name: "msr	scxtnum_el3, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0xd0, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	scxtnum_el3, x5",
			wantErr: false,
		},
		{
			name: "msr	scxtnum_el12, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe4, 0xd0, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	scxtnum_el12, x4",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.5a-ssbs.s
		{
			name: "mrs	x2, ssbs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc2, 0x42, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x2, ssbs",
			wantErr: false,
		},
		{
			name: "msr	ssbs, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc3, 0x42, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	ssbs, x3",
			wantErr: false,
		},
		{
			name: "msr	ssbs, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x41, 0x03, 0xd5}),
				address:          0,
			},
			want: "msr	ssbs, #1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fmt.Printf("want: %s\n", tt.want)
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_MTE(t *testing.T) {
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
			name: "irg	x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0xdf, 0x9a}),
				address:          0,
			},
			want: "irg	x0, x1",
			wantErr: false,
		},
		{
			name: "irg	sp, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x10, 0xdf, 0x9a}),
				address:          0,
			},
			want: "irg	sp, x1",
			wantErr: false,
		},
		{
			name: "irg	x0, sp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x13, 0xdf, 0x9a}),
				address:          0,
			},
			want: "irg	x0, sp",
			wantErr: false,
		},
		{
			name: "irg	x0, x1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x10, 0xc2, 0x9a}),
				address:          0,
			},
			want: "irg	x0, x1, x2",
			wantErr: false,
		},
		{
			name: "irg	sp, x1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x10, 0xc2, 0x9a}),
				address:          0,
			},
			want: "irg	sp, x1, x2",
			wantErr: false,
		},
		{
			name: "addg	x0, x1, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x80, 0x91}),
				address:          0,
			},
			want: "addg	x0, x1, #0, #1",
			wantErr: false,
		},
		{
			name: "addg	sp, x2, #32, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x0c, 0x82, 0x91}),
				address:          0,
			},
			want: "addg	sp, x2, #32, #3",
			wantErr: false,
		},
		{
			name: "addg	x0, sp, #64, #5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x17, 0x84, 0x91}),
				address:          0,
			},
			want: "addg	x0, sp, #64, #5",
			wantErr: false,
		},
		{
			name: "addg	x3, x4, #1008, #6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x18, 0xbf, 0x91}),
				address:          0,
			},
			want: "addg	x3, x4, #1008, #6",
			wantErr: false,
		},
		{
			name: "addg	x5, x6, #112, #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x3c, 0x87, 0x91}),
				address:          0,
			},
			want: "addg	x5, x6, #112, #15",
			wantErr: false,
		},
		{
			name: "subg	x0, x1, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x80, 0xd1}),
				address:          0,
			},
			want: "subg	x0, x1, #0, #1",
			wantErr: false,
		},
		{
			name: "subg	sp, x2, #32, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x0c, 0x82, 0xd1}),
				address:          0,
			},
			want: "subg	sp, x2, #32, #3",
			wantErr: false,
		},
		{
			name: "subg	x0, sp, #64, #5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x17, 0x84, 0xd1}),
				address:          0,
			},
			want: "subg	x0, sp, #64, #5",
			wantErr: false,
		},
		{
			name: "subg	x3, x4, #1008, #6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x18, 0xbf, 0xd1}),
				address:          0,
			},
			want: "subg	x3, x4, #1008, #6",
			wantErr: false,
		},
		{
			name: "subg	x5, x6, #112, #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x3c, 0x87, 0xd1}),
				address:          0,
			},
			want: "subg	x5, x6, #112, #15",
			wantErr: false,
		},
		{
			name: "gmi	x0, x1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x14, 0xc2, 0x9a}),
				address:          0,
			},
			want: "gmi	x0, x1, x2",
			wantErr: false,
		},
		{
			name: "gmi	x3, sp, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0xc4, 0x9a}),
				address:          0,
			},
			want: "gmi	x3, sp, x4",
			wantErr: false,
		},
		{
			name: "gmi	xzr, x0, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x14, 0xde, 0x9a}),
				address:          0,
			},
			want: "gmi	xzr, x0, x30",
			wantErr: false,
		},
		{
			name: "gmi	x30, x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1e, 0x14, 0xdf, 0x9a}),
				address:          0,
			},
			want: "gmi	x30, x0, xzr",
			wantErr: false,
		},
		{
			name: "stg	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	x0, [x1]",
			wantErr: false,
		},
		{
			name: "stg	x1, [x1, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0x08, 0x30, 0xd9}),
				address:          0,
			},
			want: "stg	x1, [x1, #-4096]",
			wantErr: false,
		},
		{
			name: "stg	x2, [x2, #4080]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0xf8, 0x2f, 0xd9}),
				address:          0,
			},
			want: "stg	x2, [x2, #4080]",
			wantErr: false,
		},
		{
			name: "stg	x3, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x1b, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	x3, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stg	sp, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1b, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	sp, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stzg	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	x0, [x1]",
			wantErr: false,
		},
		{
			name: "stzg	x1, [x1, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0x08, 0x70, 0xd9}),
				address:          0,
			},
			want: "stzg	x1, [x1, #-4096]",
			wantErr: false,
		},
		{
			name: "stzg	x2, [x2, #4080]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0xf8, 0x6f, 0xd9}),
				address:          0,
			},
			want: "stzg	x2, [x2, #4080]",
			wantErr: false,
		},
		{
			name: "stzg	x3, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x1b, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	x3, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stzg	sp, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1b, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	sp, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stg	x0, [x1, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0x30, 0xd9}),
				address:          0,
			},
			want: "stg	x0, [x1, #-4096]!",
			wantErr: false,
		},
		{
			name: "stg	x1, [x2, #4080]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0x2f, 0xd9}),
				address:          0,
			},
			want: "stg	x1, [x2, #4080]!",
			wantErr: false,
		},
		{
			name: "stg	x2, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x1f, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	x2, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stg	sp, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1f, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	sp, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stzg	x0, [x1, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0x70, 0xd9}),
				address:          0,
			},
			want: "stzg	x0, [x1, #-4096]!",
			wantErr: false,
		},
		{
			name: "stzg	x1, [x2, #4080]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0x6f, 0xd9}),
				address:          0,
			},
			want: "stzg	x1, [x2, #4080]!",
			wantErr: false,
		},
		{
			name: "stzg	x2, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x1f, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	x2, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stzg	sp, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1f, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	sp, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stg	x0, [x1], #-4096",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x30, 0xd9}),
				address:          0,
			},
			want: "stg	x0, [x1], #-4096",
			wantErr: false,
		},
		{
			name: "stg	x1, [x2], #4080",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xf4, 0x2f, 0xd9}),
				address:          0,
			},
			want: "stg	x1, [x2], #4080",
			wantErr: false,
		},
		{
			name: "stg	x2, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x17, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	x2, [sp], #16",
			wantErr: false,
		},
		{
			name: "stg	sp, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x17, 0x20, 0xd9}),
				address:          0,
			},
			want: "stg	sp, [sp], #16",
			wantErr: false,
		},
		{
			name: "stzg	x0, [x1], #-4096",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x70, 0xd9}),
				address:          0,
			},
			want: "stzg	x0, [x1], #-4096",
			wantErr: false,
		},
		{
			name: "stzg	x1, [x2], #4080",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xf4, 0x6f, 0xd9}),
				address:          0,
			},
			want: "stzg	x1, [x2], #4080",
			wantErr: false,
		},
		{
			name: "stzg	x2, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x17, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	x2, [sp], #16",
			wantErr: false,
		},
		{
			name: "stzg	sp, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x17, 0x60, 0xd9}),
				address:          0,
			},
			want: "stzg	sp, [sp], #16",
			wantErr: false,
		},
		{
			name: "st2g	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	x0, [x1]",
			wantErr: false,
		},
		{
			name: "st2g	x1, [x1, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0x08, 0xb0, 0xd9}),
				address:          0,
			},
			want: "st2g	x1, [x1, #-4096]",
			wantErr: false,
		},
		{
			name: "st2g	x2, [x2, #4080]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0xf8, 0xaf, 0xd9}),
				address:          0,
			},
			want: "st2g	x2, [x2, #4080]",
			wantErr: false,
		},
		{
			name: "st2g	x3, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x1b, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	x3, [sp, #16]",
			wantErr: false,
		},
		{
			name: "st2g	sp, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1b, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	sp, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stz2g	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x0, [x1]",
			wantErr: false,
		},
		{
			name: "stz2g	x1, [x1, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0x08, 0xf0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x1, [x1, #-4096]",
			wantErr: false,
		},
		{
			name: "stz2g	x2, [x2, #4080]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0xf8, 0xef, 0xd9}),
				address:          0,
			},
			want: "stz2g	x2, [x2, #4080]",
			wantErr: false,
		},
		{
			name: "stz2g	x3, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x1b, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x3, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stz2g	sp, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1b, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	sp, [sp, #16]",
			wantErr: false,
		},
		{
			name: "st2g	x0, [x1, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0xb0, 0xd9}),
				address:          0,
			},
			want: "st2g	x0, [x1, #-4096]!",
			wantErr: false,
		},
		{
			name: "st2g	x1, [x2, #4080]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xaf, 0xd9}),
				address:          0,
			},
			want: "st2g	x1, [x2, #4080]!",
			wantErr: false,
		},
		{
			name: "st2g	x2, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x1f, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	x2, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "st2g	sp, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1f, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	sp, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stz2g	x0, [x1, #-4096]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x0c, 0xf0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x0, [x1, #-4096]!",
			wantErr: false,
		},
		{
			name: "stz2g	x1, [x2, #4080]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xef, 0xd9}),
				address:          0,
			},
			want: "stz2g	x1, [x2, #4080]!",
			wantErr: false,
		},
		{
			name: "stz2g	x2, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x1f, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x2, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stz2g	sp, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x1f, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	sp, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "st2g	x0, [x1], #-4096",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0xb0, 0xd9}),
				address:          0,
			},
			want: "st2g	x0, [x1], #-4096",
			wantErr: false,
		},
		{
			name: "st2g	x1, [x2], #4080",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xf4, 0xaf, 0xd9}),
				address:          0,
			},
			want: "st2g	x1, [x2], #4080",
			wantErr: false,
		},
		{
			name: "st2g	x2, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x17, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	x2, [sp], #16",
			wantErr: false,
		},
		{
			name: "st2g	sp, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x17, 0xa0, 0xd9}),
				address:          0,
			},
			want: "st2g	sp, [sp], #16",
			wantErr: false,
		},
		{
			name: "stz2g	x0, [x1], #-4096",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0xf0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x0, [x1], #-4096",
			wantErr: false,
		},
		{
			name: "stz2g	x1, [x2], #4080",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xf4, 0xef, 0xd9}),
				address:          0,
			},
			want: "stz2g	x1, [x2], #4080",
			wantErr: false,
		},
		{
			name: "stz2g	x2, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x17, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	x2, [sp], #16",
			wantErr: false,
		},
		{
			name: "stz2g	sp, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x17, 0xe0, 0xd9}),
				address:          0,
			},
			want: "stz2g	sp, [sp], #16",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x04, 0x00, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2]",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2, #-1024]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x04, 0x20, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2, #-1024]",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2, #1008]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x84, 0x1f, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2, #1008]",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [sp, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x87, 0x00, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [sp, #16]",
			wantErr: false,
		},
		{
			name: "stgp	xzr, x1, [x2, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x84, 0x00, 0x69}),
				address:          0,
			},
			want: "stgp	xzr, x1, [x2, #16]",
			wantErr: false,
		},
		{
			name: "stgp	x0, xzr, [x2, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xfc, 0x00, 0x69}),
				address:          0,
			},
			want: "stgp	x0, xzr, [x2, #16]",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2, #-1024]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x04, 0xa0, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2, #-1024]!",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2, #1008]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x84, 0x9f, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2, #1008]!",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [sp, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x87, 0x80, 0x69}),
				address:          0,
			},
			want: "stgp	x0, x1, [sp, #16]!",
			wantErr: false,
		},
		{
			name: "stgp	xzr, x1, [x2, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x84, 0x80, 0x69}),
				address:          0,
			},
			want: "stgp	xzr, x1, [x2, #16]!",
			wantErr: false,
		},
		{
			name: "stgp	x0, xzr, [x2, #16]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xfc, 0x80, 0x69}),
				address:          0,
			},
			want: "stgp	x0, xzr, [x2, #16]!",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2], #-1024",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x04, 0xa0, 0x68}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2], #-1024",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [x2], #1008",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x84, 0x9f, 0x68}),
				address:          0,
			},
			want: "stgp	x0, x1, [x2], #1008",
			wantErr: false,
		},
		{
			name: "stgp	x0, x1, [sp], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x87, 0x80, 0x68}),
				address:          0,
			},
			want: "stgp	x0, x1, [sp], #16",
			wantErr: false,
		},
		{
			name: "stgp	xzr, x1, [x2], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x84, 0x80, 0x68}),
				address:          0,
			},
			want: "stgp	xzr, x1, [x2], #16",
			wantErr: false,
		},
		{
			name: "stgp	x0, xzr, [x2], #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xfc, 0x80, 0x68}),
				address:          0,
			},
			want: "stgp	x0, xzr, [x2], #16",
			wantErr: false,
		},
		{
			name: "dc	igvac, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x76, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	igvac, x0",
			wantErr: false,
		},
		{
			name: "dc	igsw, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0x76, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	igsw, x1",
			wantErr: false,
		},
		{
			name: "dc	cgsw, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x7a, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	cgsw, x2",
			wantErr: false,
		},
		{
			name: "dc	cigsw, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x7e, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	cigsw, x3",
			wantErr: false,
		},
		{
			name: "dc	cgvac, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x64, 0x7a, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cgvac, x4",
			wantErr: false,
		},
		{
			name: "dc	cgvap, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x65, 0x7c, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cgvap, x5",
			wantErr: false,
		},
		{
			name: "dc	cgvadp, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x66, 0x7d, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cgvadp, x6",
			wantErr: false,
		},
		{
			name: "dc	cigvac, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x67, 0x7e, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cigvac, x7",
			wantErr: false,
		},
		{
			name: "dc	gva, x8",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x68, 0x74, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	gva, x8",
			wantErr: false,
		},
		{
			name: "dc	igdvac, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x76, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	igdvac, x9",
			wantErr: false,
		},
		{
			name: "dc	igdsw, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x76, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	igdsw, x10",
			wantErr: false,
		},
		{
			name: "dc	cgdsw, x11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcb, 0x7a, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	cgdsw, x11",
			wantErr: false,
		},
		{
			name: "dc	cigdsw, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x7e, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	cigdsw, x12",
			wantErr: false,
		},
		{
			name: "dc	cgdvac, x13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xad, 0x7a, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cgdvac, x13",
			wantErr: false,
		},
		{
			name: "dc	cgdvap, x14",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xae, 0x7c, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cgdvap, x14",
			wantErr: false,
		},
		{
			name: "dc	cgdvadp, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaf, 0x7d, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cgdvadp, x15",
			wantErr: false,
		},
		{
			name: "dc	cigdvac, x16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb0, 0x7e, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cigdvac, x16",
			wantErr: false,
		},
		{
			name: "dc	gzva, x17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x91, 0x74, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	gzva, x17",
			wantErr: false,
		},
		{
			name: "mrs	x0, TCO",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x42, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, TCO",
			wantErr: false,
		},
		{
			name: "mrs	x1, GCR_EL1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc1, 0x10, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x1, GCR_EL1",
			wantErr: false,
		},
		{
			name: "mrs	x2, RGSR_EL1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x10, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x2, RGSR_EL1",
			wantErr: false,
		},
		{
			name: "mrs	x3, TFSR_EL1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x56, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x3, TFSR_EL1",
			wantErr: false,
		},
		{
			name: "mrs	x4, TFSR_EL2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x04, 0x56, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x4, TFSR_EL2",
			wantErr: false,
		},
		{
			name: "mrs	x5, TFSR_EL3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x05, 0x56, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x5, TFSR_EL3",
			wantErr: false,
		},
		{
			name: "mrs	x6, TFSR_EL12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x06, 0x56, 0x3d, 0xd5}),
				address:          0,
			},
			want: "mrs	x6, TFSR_EL12",
			wantErr: false,
		},
		{
			name: "mrs	x7, TFSRE0_EL1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x27, 0x56, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x7, TFSRE0_EL1",
			wantErr: false,
		},
		{
			name: "mrs	x7, GMID_EL1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x87, 0x00, 0x39, 0xd5}),
				address:          0,
			},
			want: "mrs	x7, GMID_EL1",
			wantErr: false,
		},
		{
			name: "msr	TCO, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x40, 0x03, 0xd5}),
				address:          0,
			},
			want: "msr	TCO, #0",
			wantErr: false,
		},
		{
			name: "msr	TCO, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x42, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	TCO, x0",
			wantErr: false,
		},
		{
			name: "msr	GCR_EL1, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc1, 0x10, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	GCR_EL1, x1",
			wantErr: false,
		},
		{
			name: "msr	RGSR_EL1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x10, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	RGSR_EL1, x2",
			wantErr: false,
		},
		{
			name: "msr	TFSR_EL1, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x56, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	TFSR_EL1, x3",
			wantErr: false,
		},
		{
			name: "msr	TFSR_EL2, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x04, 0x56, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	TFSR_EL2, x4",
			wantErr: false,
		},
		{
			name: "msr	TFSR_EL3, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x05, 0x56, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	TFSR_EL3, x5",
			wantErr: false,
		},
		{
			name: "msr	TFSR_EL12, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x06, 0x56, 0x1d, 0xd5}),
				address:          0,
			},
			want: "msr	TFSR_EL12, x6",
			wantErr: false,
		},
		{
			name: "msr	TFSRE0_EL1, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x27, 0x56, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	TFSRE0_EL1, x7",
			wantErr: false,
		},
		{
			name: "subp	x0, x1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0xc2, 0x9a}),
				address:          0,
			},
			want: "subp	x0, x1, x2",
			wantErr: false,
		},
		{
			name: "subp	x0, sp, sp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0xdf, 0x9a}),
				address:          0,
			},
			want: "subp	x0, sp, sp",
			wantErr: false,
		},
		{
			name: "subps	x0, x1, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0xc2, 0xba}),
				address:          0,
			},
			want: "subps	x0, x1, x2",
			wantErr: false,
		},
		{
			name: "subps	x0, sp, sp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0xdf, 0xba}),
				address:          0,
			},
			want: "subps	x0, sp, sp",
			wantErr: false,
		},
		{
			name: "subps	xzr, x0, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0xc1, 0xba}),
				address:          0,
			},
			want: "cmpp	x0, x1",
			wantErr: false,
		},
		{
			name: "subps	xzr, sp, sp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0xdf, 0xba}),
				address:          0,
			},
			want: "cmpp	sp, sp",
			wantErr: false,
		},
		{
			name: "ldg	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0x60, 0xd9}),
				address:          0,
			},
			want: "ldg	x0, [x1]",
			wantErr: false,
		},
		{
			name: "ldg	x2, [sp, #-4096]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x03, 0x70, 0xd9}),
				address:          0,
			},
			want: "ldg	x2, [sp, #-4096]",
			wantErr: false,
		},
		{
			name: "ldg	x3, [x4, #4080]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0xf0, 0x6f, 0xd9}),
				address:          0,
			},
			want: "ldg	x3, [x4, #4080]",
			wantErr: false,
		},
		{
			name: "ldgm	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0xe0, 0xd9}),
				address:          0,
			},
			want: "ldgm	x0, [x1]",
			wantErr: false,
		},
		{
			name: "ldgm	x1, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x03, 0xe0, 0xd9}),
				address:          0,
			},
			want: "ldgm	x1, [sp]",
			wantErr: false,
		},
		{
			name: "ldgm	xzr, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0xe0, 0xd9}),
				address:          0,
			},
			want: "ldgm	xzr, [x2]",
			wantErr: false,
		},
		{
			name: "stgm	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0xa0, 0xd9}),
				address:          0,
			},
			want: "stgm	x0, [x1]",
			wantErr: false,
		},
		{
			name: "stgm	x1, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x03, 0xa0, 0xd9}),
				address:          0,
			},
			want: "stgm	x1, [sp]",
			wantErr: false,
		},
		{
			name: "stgm	xzr, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0xa0, 0xd9}),
				address:          0,
			},
			want: "stgm	xzr, [x2]",
			wantErr: false,
		},
		{
			name: "stzgm	x0, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0x20, 0xd9}),
				address:          0,
			},
			want: "stzgm	x0, [x1]",
			wantErr: false,
		},
		{
			name: "stzgm	x1, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x03, 0x20, 0xd9}),
				address:          0,
			},
			want: "stzgm	x1, [sp]",
			wantErr: false,
		},
		{
			name: "stzgm	xzr, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x20, 0xd9}),
				address:          0,
			},
			want: "stzgm	xzr, [x2]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) {
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_v8_6a(t *testing.T) {
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
		// llvm/test/MC/AArch64/armv8.6a-amvs.s
		{
			name: "msr	amevcntvoff00_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff00_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff01_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff01_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff02_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff02_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff03_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff03_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff04_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff04_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff05_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff05_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff06_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff06_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff07_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xd8, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff07_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff08_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff08_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff09_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff09_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff010_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff010_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff011_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff011_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff012_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff012_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff013_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff013_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff014_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff014_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff015_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xd9, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff015_el2, x0",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff00_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff00_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff01_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff01_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff02_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff02_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff03_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff03_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff04_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff04_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff05_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff05_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff06_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff06_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff07_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xd8, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff07_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff08_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff08_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff09_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff09_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff010_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff010_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff011_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff011_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff012_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff012_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff013_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff013_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff014_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff014_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff015_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xd9, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff015_el2",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff10_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff10_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff11_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff11_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff12_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff12_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff13_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff13_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff14_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff14_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff15_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff15_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff16_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff16_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff17_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xda, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff17_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff18_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff18_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff19_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff19_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff110_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff110_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff111_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff111_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff112_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff112_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff113_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff113_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff114_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff114_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	amevcntvoff115_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xdb, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amevcntvoff115_el2, x0",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff10_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff10_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff11_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff11_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff12_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff12_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff13_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff13_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff14_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff14_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff15_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff15_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff16_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff16_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff17_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xda, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff17_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff18_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff18_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff19_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff19_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff110_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff110_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff111_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff111_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff112_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff112_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff113_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa0, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff113_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff114_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff114_el2",
			wantErr: false,
		},
		{
			name: "mrs	x0, amevcntvoff115_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xdb, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, amevcntvoff115_el2",
			wantErr: false,
		},

		// llvm/test/MC/AArch64/armv8.6a-bf16.s
		{
			name: "bfdot	v2.2s, v3.4h, v4.4h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x44, 0x2e}),
				address:          0,
			},
			want: "bfdot	v2.2s, v3.4h, v4.4h",
			wantErr: false,
		},
		{
			name: "bfdot	v2.4s, v3.8h, v4.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x44, 0x6e}),
				address:          0,
			},
			want: "bfdot	v2.4s, v3.8h, v4.8h",
			wantErr: false,
		},
		{
			name: "bfdot	v2.2s, v3.4h, v4.2h[0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf0, 0x44, 0x0f}),
				address:          0,
			},
			want: "bfdot	v2.2s, v3.4h, v4.2h[0]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.2s, v3.4h, v4.2h[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf0, 0x64, 0x0f}),
				address:          0,
			},
			want: "bfdot	v2.2s, v3.4h, v4.2h[1]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.2s, v3.4h, v4.2h[2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x44, 0x0f}),
				address:          0,
			},
			want: "bfdot	v2.2s, v3.4h, v4.2h[2]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.2s, v3.4h, v4.2h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x64, 0x0f}),
				address:          0,
			},
			want: "bfdot	v2.2s, v3.4h, v4.2h[3]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.4s, v3.8h, v4.2h[0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf0, 0x44, 0x4f}),
				address:          0,
			},
			want: "bfdot	v2.4s, v3.8h, v4.2h[0]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.4s, v3.8h, v4.2h[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf0, 0x64, 0x4f}),
				address:          0,
			},
			want: "bfdot	v2.4s, v3.8h, v4.2h[1]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.4s, v3.8h, v4.2h[2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x44, 0x4f}),
				address:          0,
			},
			want: "bfdot	v2.4s, v3.8h, v4.2h[2]",
			wantErr: false,
		},
		{
			name: "bfdot	v2.4s, v3.8h, v4.2h[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xf8, 0x64, 0x4f}),
				address:          0,
			},
			want: "bfdot	v2.4s, v3.8h, v4.2h[3]",
			wantErr: false,
		},
		{
			name: "bfmmla	v2.4s, v3.8h, v4.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xec, 0x44, 0x6e}),
				address:          0,
			},
			want: "bfmmla	v2.4s, v3.8h, v4.8h",
			wantErr: false,
		},
		{
			name: "bfmmla	v3.4s, v4.8h, v5.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0xec, 0x45, 0x6e}),
				address:          0,
			},
			want: "bfmmla	v3.4s, v4.8h, v5.8h",
			wantErr: false,
		},
		{
			name: "bfcvtn	v5.4h, v5.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa5, 0x68, 0xa1, 0x0e}),
				address:          0,
			},
			want: "bfcvtn	v5.4h, v5.4s",
			wantErr: false,
		},
		{
			name: "bfcvtn2	v5.8h, v5.4s",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa5, 0x68, 0xa1, 0x4e}),
				address:          0,
			},
			want: "bfcvtn2	v5.8h, v5.4s",
			wantErr: false,
		},
		{
			name: "bfcvt	h5, s3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x65, 0x40, 0x63, 0x1e}),
				address:          0,
			},
			want: "bfcvt	h5, s3",
			wantErr: false,
		},
		{
			name: "bfmlalb	v10.4s, v21.8h, v14.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaa, 0xfe, 0xce, 0x2e}),
				address:          0,
			},
			want: "bfmlalb	v10.4s, v21.8h, v14.8h",
			wantErr: false,
		},
		{
			name: "bfmlalt	v21.4s, v14.8h, v10.8h",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0xfd, 0xca, 0x6e}),
				address:          0,
			},
			want: "bfmlalt	v21.4s, v14.8h, v10.8h",
			wantErr: false,
		},
		{
			name: "bfmlalb	v14.4s, v21.8h, v10.h[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xae, 0xf2, 0xda, 0x0f}),
				address:          0,
			},
			want: "bfmlalb	v14.4s, v21.8h, v10.h[1]",
			wantErr: false,
		},
		{
			name: "bfmlalb	v14.4s, v21.8h, v10.h[2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xae, 0xf2, 0xea, 0x0f}),
				address:          0,
			},
			want: "bfmlalb	v14.4s, v21.8h, v10.h[2]",
			wantErr: false,
		},
		{
			name: "bfmlalb	v14.4s, v21.8h, v10.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xae, 0xfa, 0xfa, 0x0f}),
				address:          0,
			},
			want: "bfmlalb	v14.4s, v21.8h, v10.h[7]",
			wantErr: false,
		},
		{
			name: "bfmlalt	v21.4s, v10.8h, v14.h[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf1, 0xde, 0x4f}),
				address:          0,
			},
			want: "bfmlalt	v21.4s, v10.8h, v14.h[1]",
			wantErr: false,
		},
		{
			name: "bfmlalt	v21.4s, v10.8h, v14.h[2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf1, 0xee, 0x4f}),
				address:          0,
			},
			want: "bfmlalt	v21.4s, v10.8h, v14.h[2]",
			wantErr: false,
		},
		{
			name: "bfmlalt	v21.4s, v10.8h, v14.h[7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf9, 0xfe, 0x4f}),
				address:          0,
			},
			want: "bfmlalt	v21.4s, v10.8h, v14.h[7]",
			wantErr: false,
		},
		// llvm/test/MC/AArch64/armv8.6a-ecv.s
		{
			name: "msr	cntscale_el2, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0xe0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cntscale_el2, x1",
			wantErr: false,
		},
		{
			name: "msr	cntiscale_el2, x11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xe0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cntiscale_el2, x11",
			wantErr: false,
		},
		{
			name: "msr	cntpoff_el2, x22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd6, 0xe0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cntpoff_el2, x22",
			wantErr: false,
		},
		{
			name: "msr	cntvfrq_el2, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xe0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cntvfrq_el2, x3",
			wantErr: false,
		},
		{
			name: "msr	cntpctss_el0, x13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xad, 0xe0, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntpctss_el0, x13",
			wantErr: false,
		},
		{
			name: "msr	cntvctss_el0, x23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd7, 0xe0, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntvctss_el0, x23",
			wantErr: false,
		},
		{
			name: "mrs	x0, cntscale_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0xe0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x0, cntscale_el2",
			wantErr: false,
		},
		{
			name: "mrs	x5, cntiscale_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa5, 0xe0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x5, cntiscale_el2",
			wantErr: false,
		},
		{
			name: "mrs	x10, cntpoff_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0xe0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x10, cntpoff_el2",
			wantErr: false,
		},
		{
			name: "mrs	x15, cntvfrq_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0xe0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x15, cntvfrq_el2",
			wantErr: false,
		},
		{
			name: "mrs	x20, cntpctss_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0xe0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x20, cntpctss_el0",
			wantErr: false,
		},
		{
			name: "mrs	x30, cntvctss_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xde, 0xe0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x30, cntvctss_el0",
			wantErr: false,
		},

		// llvm/test/MC/AArch64/armv8.6a-fgt.s
		{
			name: "msr	hfgrtr_el2, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hfgrtr_el2, x0",
			wantErr: false,
		},
		{
			name: "msr	hfgwtr_el2, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa5, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hfgwtr_el2, x5",
			wantErr: false,
		},
		{
			name: "msr	hfgitr_el2, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hfgitr_el2, x10",
			wantErr: false,
		},
		{
			name: "msr	hdfgrtr_el2, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8f, 0x31, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hdfgrtr_el2, x15",
			wantErr: false,
		},
		{
			name: "msr	hdfgwtr_el2, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x31, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hdfgwtr_el2, x20",
			wantErr: false,
		},
		{
			name: "mrs	x30, hfgrtr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9e, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x30, hfgrtr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x25, hfgwtr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb9, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x25, hfgwtr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x20, hfgitr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd4, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x20, hfgitr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x15, hdfgrtr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8f, 0x31, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x15, hdfgrtr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x10, hdfgwtr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaa, 0x31, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x10, hdfgwtr_el2",
			wantErr: false,
		},

		// llvm/test/MC/AArch64/armv8.6a-simd-matmul.s
		{
			name: "smmla	v1.4s, v16.16b, v31.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0xa6, 0x9f, 0x4e}),
				address:          0,
			},
			want: "smmla	v1.4s, v16.16b, v31.16b",
			wantErr: false,
		},
		{
			name: "ummla	v1.4s, v16.16b, v31.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0xa6, 0x9f, 0x6e}),
				address:          0,
			},
			want: "ummla	v1.4s, v16.16b, v31.16b",
			wantErr: false,
		},
		{
			name: "usmmla	v1.4s, v16.16b, v31.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0xae, 0x9f, 0x4e}),
				address:          0,
			},
			want: "usmmla	v1.4s, v16.16b, v31.16b",
			wantErr: false,
		},
		{
			name: "usdot	v3.2s, v15.8b, v30.8b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x9d, 0x9e, 0x0e}),
				address:          0,
			},
			want: "usdot	v3.2s, v15.8b, v30.8b",
			wantErr: false,
		},
		{
			name: "usdot	v3.4s, v15.16b, v30.16b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x9d, 0x9e, 0x4e}),
				address:          0,
			},
			want: "usdot	v3.4s, v15.16b, v30.16b",
			wantErr: false,
		},
		{
			name: "usdot	v31.2s, v1.8b, v2.4b[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf8, 0xa2, 0x0f}),
				address:          0,
			},
			want: "usdot	v31.2s, v1.8b, v2.4b[3]",
			wantErr: false,
		},
		{
			name: "usdot	v31.4s, v1.16b, v2.4b[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf8, 0xa2, 0x4f}),
				address:          0,
			},
			want: "usdot	v31.4s, v1.16b, v2.4b[3]",
			wantErr: false,
		},
		{
			name: "sudot	v31.2s, v1.8b, v2.4b[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf8, 0x22, 0x0f}),
				address:          0,
			},
			want: "sudot	v31.2s, v1.8b, v2.4b[3]",
			wantErr: false,
		},
		{
			name: "sudot	v31.4s, v1.16b, v2.4b[3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf8, 0x22, 0x4f}),
				address:          0,
			},
			want: "sudot	v31.4s, v1.16b, v2.4b[3]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.HasPrefix(tt.want, "m") {
				fmt.Printf("want: %s\n", tt.want)
			}
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}

func Test_decompose_basic(t *testing.T) {
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
			name: "add	x2, x4, w5, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x00, 0x25, 0x8b}),
				address:          0,
			},
			want: "add	x2, x4, w5, uxtb",
			wantErr: false,
		},
		{
			name: "add	x20, sp, w19, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x23, 0x33, 0x8b}),
				address:          0,
			},
			want: "add	x20, sp, w19, uxth",
			wantErr: false,
		},
		{
			name: "add	x12, x1, w20, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x34, 0x8b}),
				address:          0,
			},
			want: "add	x12, x1, w20, uxtw",
			wantErr: false,
		},
		{
			name: "add	x20, x3, x13, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0x60, 0x2d, 0x8b}),
				address:          0,
			},
			want: "add	x20, x3, x13, uxtx",
			wantErr: false,
		},
		{
			name: "add	x17, x25, w20, sxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x31, 0x83, 0x34, 0x8b}),
				address:          0,
			},
			want: "add	x17, x25, w20, sxtb",
			wantErr: false,
		},
		{
			name: "add	x18, x13, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb2, 0xa1, 0x33, 0x8b}),
				address:          0,
			},
			want: "add	x18, x13, w19, sxth",
			wantErr: false,
		},
		{
			name: "add	sp, x2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0x8b}),
				address:          0,
			},
			want: "add	sp, x2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "add	x3, x5, x9, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xe0, 0x29, 0x8b}),
				address:          0,
			},
			want: "add	x3, x5, x9, sxtx",
			wantErr: false,
		},
		{
			name: "add	w2, w5, w7, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x00, 0x27, 0x0b}),
				address:          0,
			},
			want: "add	w2, w5, w7, uxtb",
			wantErr: false,
		},
		{
			name: "add	w21, w15, w17, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf5, 0x21, 0x31, 0x0b}),
				address:          0,
			},
			want: "add	w21, w15, w17, uxth",
			wantErr: false,
		},
		{
			name: "add	w30, w29, wzr, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x43, 0x3f, 0x0b}),
				address:          0,
			},
			want: "add	w30, w29, wzr, uxtw",
			wantErr: false,
		},
		{
			name: "add	w19, w17, w1, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x62, 0x21, 0x0b}),
				address:          0,
			},
			want: "add	w19, w17, w1, uxtx",
			wantErr: false,
		},
		{
			name: "add	w2, w5, w1, sxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x80, 0x21, 0x0b}),
				address:          0,
			},
			want: "add	w2, w5, w1, sxtb",
			wantErr: false,
		},
		{
			name: "add	w26, w17, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3a, 0xa2, 0x33, 0x0b}),
				address:          0,
			},
			want: "add	w26, w17, w19, sxth",
			wantErr: false,
		},
		{
			name: "add	w0, w2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0xc0, 0x23, 0x0b}),
				address:          0,
			},
			want: "add	w0, w2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "add	w2, w3, w5, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xe0, 0x25, 0x0b}),
				address:          0,
			},
			want: "add	w2, w3, w5, sxtx",
			wantErr: false,
		},
		{
			name: "add	x2, x3, w5, sxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x80, 0x25, 0x8b}),
				address:          0,
			},
			want: "add	x2, x3, w5, sxtb",
			wantErr: false,
		},
		{
			name: "add	x7, x11, w13, uxth #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x67, 0x31, 0x2d, 0x8b}),
				address:          0,
			},
			want: "add	x7, x11, w13, uxth #4",
			wantErr: false,
		},
		{
			name: "add	w17, w19, w23, uxtw #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x71, 0x4a, 0x37, 0x0b}),
				address:          0,
			},
			want: "add	w17, w19, w23, uxtw #2",
			wantErr: false,
		},
		{
			name: "add	w29, w23, w17, uxtx #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x66, 0x31, 0x0b}),
				address:          0,
			},
			want: "add	w29, w23, w17, uxtx #1",
			wantErr: false,
		},
		{
			name: "sub	x2, x4, w5, uxtb #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x08, 0x25, 0xcb}),
				address:          0,
			},
			want: "sub	x2, x4, w5, uxtb #2",
			wantErr: false,
		},
		{
			name: "sub	x20, sp, w19, uxth #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x33, 0x33, 0xcb}),
				address:          0,
			},
			want: "sub	x20, sp, w19, uxth #4",
			wantErr: false,
		},
		{
			name: "sub	x12, x1, w20, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x34, 0xcb}),
				address:          0,
			},
			want: "sub	x12, x1, w20, uxtw",
			wantErr: false,
		},
		{
			name: "sub	x20, x3, x13, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0x60, 0x2d, 0xcb}),
				address:          0,
			},
			want: "sub	x20, x3, x13, uxtx",
			wantErr: false,
		},
		{
			name: "sub	x17, x25, w20, sxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x31, 0x83, 0x34, 0xcb}),
				address:          0,
			},
			want: "sub	x17, x25, w20, sxtb",
			wantErr: false,
		},
		{
			name: "sub	x18, x13, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb2, 0xa1, 0x33, 0xcb}),
				address:          0,
			},
			want: "sub	x18, x13, w19, sxth",
			wantErr: false,
		},
		{
			name: "sub	sp, x2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0xcb}),
				address:          0,
			},
			want: "sub	sp, x2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "sub	x3, x5, x9, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xe0, 0x29, 0xcb}),
				address:          0,
			},
			want: "sub	x3, x5, x9, sxtx",
			wantErr: false,
		},
		{
			name: "sub	w2, w5, w7, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x00, 0x27, 0x4b}),
				address:          0,
			},
			want: "sub	w2, w5, w7, uxtb",
			wantErr: false,
		},
		{
			name: "sub	w21, w15, w17, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf5, 0x21, 0x31, 0x4b}),
				address:          0,
			},
			want: "sub	w21, w15, w17, uxth",
			wantErr: false,
		},
		{
			name: "sub	w30, w29, wzr, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x43, 0x3f, 0x4b}),
				address:          0,
			},
			want: "sub	w30, w29, wzr, uxtw",
			wantErr: false,
		},
		{
			name: "sub	w19, w17, w1, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x62, 0x21, 0x4b}),
				address:          0,
			},
			want: "sub	w19, w17, w1, uxtx",
			wantErr: false,
		},
		{
			name: "sub	w2, w5, w1, sxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x80, 0x21, 0x4b}),
				address:          0,
			},
			want: "sub	w2, w5, w1, sxtb",
			wantErr: false,
		},
		{
			name: "sub	w26, wsp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0xa3, 0x33, 0x4b}),
				address:          0,
			},
			want: "sub	w26, wsp, w19, sxth",
			wantErr: false,
		},
		{
			name: "sub	wsp, w2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0x4b}),
				address:          0,
			},
			want: "sub	wsp, w2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "sub	w2, w3, w5, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xe0, 0x25, 0x4b}),
				address:          0,
			},
			want: "sub	w2, w3, w5, sxtx",
			wantErr: false,
		},
		{
			name: "adds	x2, x4, w5, uxtb #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x08, 0x25, 0xab}),
				address:          0,
			},
			want: "adds	x2, x4, w5, uxtb #2",
			wantErr: false,
		},
		{
			name: "adds	x20, sp, w19, uxth #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x33, 0x33, 0xab}),
				address:          0,
			},
			want: "adds	x20, sp, w19, uxth #4",
			wantErr: false,
		},
		{
			name: "adds	x12, x1, w20, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x34, 0xab}),
				address:          0,
			},
			want: "adds	x12, x1, w20, uxtw",
			wantErr: false,
		},
		{
			name: "adds	x20, x3, x13, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0x60, 0x2d, 0xab}),
				address:          0,
			},
			want: "adds	x20, x3, x13, uxtx",
			wantErr: false,
		},
		{
			name: "cmn	x25, w20, sxtb #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x8f, 0x34, 0xab}),
				address:          0,
			},
			want: "cmn	x25, w20, sxtb #3",
			wantErr: false,
		},
		{
			name: "adds	x18, sp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0xa3, 0x33, 0xab}),
				address:          0,
			},
			want: "adds	x18, sp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmn	x2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0xab}),
				address:          0,
			},
			want: "cmn	x2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "adds	x3, x5, x9, sxtx #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xe8, 0x29, 0xab}),
				address:          0,
			},
			want: "adds	x3, x5, x9, sxtx #2",
			wantErr: false,
		},
		{
			name: "adds	w2, w5, w7, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x00, 0x27, 0x2b}),
				address:          0,
			},
			want: "adds	w2, w5, w7, uxtb",
			wantErr: false,
		},
		{
			name: "adds	w21, w15, w17, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf5, 0x21, 0x31, 0x2b}),
				address:          0,
			},
			want: "adds	w21, w15, w17, uxth",
			wantErr: false,
		},
		{
			name: "adds	w30, w29, wzr, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x43, 0x3f, 0x2b}),
				address:          0,
			},
			want: "adds	w30, w29, wzr, uxtw",
			wantErr: false,
		},
		{
			name: "adds	w19, w17, w1, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x62, 0x21, 0x2b}),
				address:          0,
			},
			want: "adds	w19, w17, w1, uxtx",
			wantErr: false,
		},
		{
			name: "adds	w2, w5, w1, sxtb #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x84, 0x21, 0x2b}),
				address:          0,
			},
			want: "adds	w2, w5, w1, sxtb #1",
			wantErr: false,
		},
		{
			name: "adds	w26, wsp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0xa3, 0x33, 0x2b}),
				address:          0,
			},
			want: "adds	w26, wsp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmn	w2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0x2b}),
				address:          0,
			},
			want: "cmn	w2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "adds	w2, w3, w5, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xe0, 0x25, 0x2b}),
				address:          0,
			},
			want: "adds	w2, w3, w5, sxtx",
			wantErr: false,
		},
		{
			name: "subs	x2, x4, w5, uxtb #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x08, 0x25, 0xeb}),
				address:          0,
			},
			want: "subs	x2, x4, w5, uxtb #2",
			wantErr: false,
		},
		{
			name: "subs	x20, sp, w19, uxth #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x33, 0x33, 0xeb}),
				address:          0,
			},
			want: "subs	x20, sp, w19, uxth #4",
			wantErr: false,
		},
		{
			name: "subs	x12, x1, w20, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x34, 0xeb}),
				address:          0,
			},
			want: "subs	x12, x1, w20, uxtw",
			wantErr: false,
		},
		{
			name: "subs	x20, x3, x13, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0x60, 0x2d, 0xeb}),
				address:          0,
			},
			want: "subs	x20, x3, x13, uxtx",
			wantErr: false,
		},
		{
			name: "cmp	x25, w20, sxtb #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x8f, 0x34, 0xeb}),
				address:          0,
			},
			want: "cmp	x25, w20, sxtb #3",
			wantErr: false,
		},
		{
			name: "subs	x18, sp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0xa3, 0x33, 0xeb}),
				address:          0,
			},
			want: "subs	x18, sp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmp	x2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0xeb}),
				address:          0,
			},
			want: "cmp	x2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "subs	x3, x5, x9, sxtx #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xe8, 0x29, 0xeb}),
				address:          0,
			},
			want: "subs	x3, x5, x9, sxtx #2",
			wantErr: false,
		},
		{
			name: "subs	w2, w5, w7, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x00, 0x27, 0x6b}),
				address:          0,
			},
			want: "subs	w2, w5, w7, uxtb",
			wantErr: false,
		},
		{
			name: "subs	w21, w15, w17, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf5, 0x21, 0x31, 0x6b}),
				address:          0,
			},
			want: "subs	w21, w15, w17, uxth",
			wantErr: false,
		},
		{
			name: "subs	w30, w29, wzr, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x43, 0x3f, 0x6b}),
				address:          0,
			},
			want: "subs	w30, w29, wzr, uxtw",
			wantErr: false,
		},
		{
			name: "subs	w19, w17, w1, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x62, 0x21, 0x6b}),
				address:          0,
			},
			want: "subs	w19, w17, w1, uxtx",
			wantErr: false,
		},
		{
			name: "subs	w2, w5, w1, sxtb #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x84, 0x21, 0x6b}),
				address:          0,
			},
			want: "subs	w2, w5, w1, sxtb #1",
			wantErr: false,
		},
		{
			name: "subs	w26, wsp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0xa3, 0x33, 0x6b}),
				address:          0,
			},
			want: "subs	w26, wsp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmp	w2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0x6b}),
				address:          0,
			},
			want: "cmp	w2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "subs	w2, w3, w5, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xe0, 0x25, 0x6b}),
				address:          0,
			},
			want: "subs	w2, w3, w5, sxtx",
			wantErr: false,
		},
		{
			name: "cmp	x4, w5, uxtb #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x08, 0x25, 0xeb}),
				address:          0,
			},
			want: "cmp	x4, w5, uxtb #2",
			wantErr: false,
		},
		{
			name: "cmp	sp, w19, uxth #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x33, 0xeb}),
				address:          0,
			},
			want: "cmp	sp, w19, uxth #4",
			wantErr: false,
		},
		{
			name: "cmp	x1, w20, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x40, 0x34, 0xeb}),
				address:          0,
			},
			want: "cmp	x1, w20, uxtw",
			wantErr: false,
		},
		{
			name: "cmp	x3, x13, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x60, 0x2d, 0xeb}),
				address:          0,
			},
			want: "cmp	x3, x13, uxtx",
			wantErr: false,
		},
		{
			name: "cmp	x25, w20, sxtb #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x8f, 0x34, 0xeb}),
				address:          0,
			},
			want: "cmp	x25, w20, sxtb #3",
			wantErr: false,
		},
		{
			name: "cmp	sp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa3, 0x33, 0xeb}),
				address:          0,
			},
			want: "cmp	sp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmp	x2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0xeb}),
				address:          0,
			},
			want: "cmp	x2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "cmp	x5, x9, sxtx #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0xe8, 0x29, 0xeb}),
				address:          0,
			},
			want: "cmp	x5, x9, sxtx #2",
			wantErr: false,
		},
		{
			name: "cmp	w5, w7, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x27, 0x6b}),
				address:          0,
			},
			want: "cmp	w5, w7, uxtb",
			wantErr: false,
		},
		{
			name: "cmp	w15, w17, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x21, 0x31, 0x6b}),
				address:          0,
			},
			want: "cmp	w15, w17, uxth",
			wantErr: false,
		},
		{
			name: "cmp	w29, wzr, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x43, 0x3f, 0x6b}),
				address:          0,
			},
			want: "cmp	w29, wzr, uxtw",
			wantErr: false,
		},
		{
			name: "cmp	w17, w1, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x62, 0x21, 0x6b}),
				address:          0,
			},
			want: "cmp	w17, w1, uxtx",
			wantErr: false,
		},
		{
			name: "cmp	w5, w1, sxtb #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x84, 0x21, 0x6b}),
				address:          0,
			},
			want: "cmp	w5, w1, sxtb #1",
			wantErr: false,
		},
		{
			name: "cmp	wsp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa3, 0x33, 0x6b}),
				address:          0,
			},
			want: "cmp	wsp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmp	w2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0x6b}),
				address:          0,
			},
			want: "cmp	w2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "cmp	w3, w5, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xe0, 0x25, 0x6b}),
				address:          0,
			},
			want: "cmp	w3, w5, sxtx",
			wantErr: false,
		},
		{
			name: "cmn	x4, w5, uxtb #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x08, 0x25, 0xab}),
				address:          0,
			},
			want: "cmn	x4, w5, uxtb #2",
			wantErr: false,
		},
		{
			name: "cmn	sp, w19, uxth #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x33, 0x33, 0xab}),
				address:          0,
			},
			want: "cmn	sp, w19, uxth #4",
			wantErr: false,
		},
		{
			name: "cmn	x1, w20, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x40, 0x34, 0xab}),
				address:          0,
			},
			want: "cmn	x1, w20, uxtw",
			wantErr: false,
		},
		{
			name: "cmn	x3, x13, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x60, 0x2d, 0xab}),
				address:          0,
			},
			want: "cmn	x3, x13, uxtx",
			wantErr: false,
		},
		{
			name: "cmn	x25, w20, sxtb #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x8f, 0x34, 0xab}),
				address:          0,
			},
			want: "cmn	x25, w20, sxtb #3",
			wantErr: false,
		},
		{
			name: "cmn	sp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa3, 0x33, 0xab}),
				address:          0,
			},
			want: "cmn	sp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmn	x2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0xab}),
				address:          0,
			},
			want: "cmn	x2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "cmn	x5, x9, sxtx #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0xe8, 0x29, 0xab}),
				address:          0,
			},
			want: "cmn	x5, x9, sxtx #2",
			wantErr: false,
		},
		{
			name: "cmn	w5, w7, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x27, 0x2b}),
				address:          0,
			},
			want: "cmn	w5, w7, uxtb",
			wantErr: false,
		},
		{
			name: "cmn	w15, w17, uxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x21, 0x31, 0x2b}),
				address:          0,
			},
			want: "cmn	w15, w17, uxth",
			wantErr: false,
		},
		{
			name: "cmn	w29, wzr, uxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x43, 0x3f, 0x2b}),
				address:          0,
			},
			want: "cmn	w29, wzr, uxtw",
			wantErr: false,
		},
		{
			name: "cmn	w17, w1, uxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x62, 0x21, 0x2b}),
				address:          0,
			},
			want: "cmn	w17, w1, uxtx",
			wantErr: false,
		},
		{
			name: "cmn	w5, w1, sxtb #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x84, 0x21, 0x2b}),
				address:          0,
			},
			want: "cmn	w5, w1, sxtb #1",
			wantErr: false,
		},
		{
			name: "cmn	wsp, w19, sxth",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa3, 0x33, 0x2b}),
				address:          0,
			},
			want: "cmn	wsp, w19, sxth",
			wantErr: false,
		},
		{
			name: "cmn	w2, w3, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xc0, 0x23, 0x2b}),
				address:          0,
			},
			want: "cmn	w2, w3, sxtw",
			wantErr: false,
		},
		{
			name: "cmn	w3, w5, sxtx",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xe0, 0x25, 0x2b}),
				address:          0,
			},
			want: "cmn	w3, w5, sxtx",
			wantErr: false,
		},
		{
			name: "cmp	x20, w29, uxtb #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x0e, 0x3d, 0xeb}),
				address:          0,
			},
			want: "cmp	x20, w29, uxtb #3",
			wantErr: false,
		},
		{
			name: "cmp	x12, x13, uxtx #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x71, 0x2d, 0xeb}),
				address:          0,
			},
			want: "cmp	x12, x13, uxtx #4",
			wantErr: false,
		},
		{
			name: "cmp	wsp, w1, uxtb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x21, 0x6b}),
				address:          0,
			},
			want: "cmp	wsp, w1, uxtb",
			wantErr: false,
		},
		{
			name: "cmn	wsp, wzr, sxtw",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xc3, 0x3f, 0x2b}),
				address:          0,
			},
			want: "cmn	wsp, wzr, sxtw",
			wantErr: false,
		},
		{
			name: "sub	sp, x3, x7, lsl #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x70, 0x27, 0xcb}),
				address:          0,
			},
			want: "sub	sp, x3, x7, lsl #4",
			wantErr: false,
		},
		{
			name: "add	w2, wsp, w3, lsl #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x47, 0x23, 0x0b}),
				address:          0,
			},
			want: "add	w2, wsp, w3, lsl #1",
			wantErr: false,
		},
		{
			name: "cmp	wsp, w9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x29, 0x6b}),
				address:          0,
			},
			want: "cmp	wsp, w9",
			wantErr: false,
		},
		{
			name: "cmn	wsp, w3, lsl #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x23, 0x2b}),
				address:          0,
			},
			want: "cmn	wsp, w3, lsl #4",
			wantErr: false,
		},
		{
			name: "subs	x3, sp, x9, lsl #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x6b, 0x29, 0xeb}),
				address:          0,
			},
			want: "subs	x3, sp, x9, lsl #2",
			wantErr: false,
		},
		{
			name: "add	w4, w5, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x00, 0x00, 0x11}),
				address:          0,
			},
			want: "add	w4, w5, #0",
			wantErr: false,
		},
		{
			name: "add	w2, w3, #4095",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x3f, 0x11}),
				address:          0,
			},
			want: "add	w2, w3, #4095",
			wantErr: false,
		},
		{
			name: "add	w30, w29, #1, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x07, 0x40, 0x11}),
				address:          0,
			},
			want: "add	w30, w29, #1, lsl #12",
			wantErr: false,
		},
		{
			name: "add	w13, w5, #4095, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xad, 0xfc, 0x7f, 0x11}),
				address:          0,
			},
			want: "add	w13, w5, #4095, lsl #12",
			wantErr: false,
		},
		{
			name: "add	x5, x7, #1638",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x98, 0x19, 0x91}),
				address:          0,
			},
			want: "add	x5, x7, #1638",
			wantErr: false,
		},
		{
			name: "add	w20, wsp, #801",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x87, 0x0c, 0x11}),
				address:          0,
			},
			want: "add	w20, wsp, #801",
			wantErr: false,
		},
		{
			name: "add	wsp, wsp, #1104",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x11, 0x11}),
				address:          0,
			},
			want: "add	wsp, wsp, #1104",
			wantErr: false,
		},
		{
			name: "add	wsp, w30, #4084",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xd3, 0x3f, 0x11}),
				address:          0,
			},
			want: "add	wsp, w30, #4084",
			wantErr: false,
		},
		{
			name: "add	x0, x24, #291",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x8f, 0x04, 0x91}),
				address:          0,
			},
			want: "add	x0, x24, #291",
			wantErr: false,
		},
		{
			name: "add	x3, x24, #4095, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0xff, 0x7f, 0x91}),
				address:          0,
			},
			want: "add	x3, x24, #4095, lsl #12",
			wantErr: false,
		},
		{
			name: "add	x8, sp, #1074",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0xcb, 0x10, 0x91}),
				address:          0,
			},
			want: "add	x8, sp, #1074",
			wantErr: false,
		},
		{
			name: "add	sp, x29, #3816",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0xa3, 0x3b, 0x91}),
				address:          0,
			},
			want: "add	sp, x29, #3816",
			wantErr: false,
		},
		{
			name: "sub	w0, wsp, #4077",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xb7, 0x3f, 0x51}),
				address:          0,
			},
			want: "sub	w0, wsp, #4077",
			wantErr: false,
		},
		{
			name: "sub	w4, w20, #546, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x84, 0x8a, 0x48, 0x51}),
				address:          0,
			},
			want: "sub	w4, w20, #546, lsl #12",
			wantErr: false,
		},
		{
			name: "sub	sp, sp, #288",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x83, 0x04, 0xd1}),
				address:          0,
			},
			want: "sub	sp, sp, #288",
			wantErr: false,
		},
		{
			name: "sub	wsp, w19, #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x42, 0x00, 0x51}),
				address:          0,
			},
			want: "sub	wsp, w19, #16",
			wantErr: false,
		},
		{
			name: "adds	w13, w23, #291, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x8e, 0x44, 0x31}),
				address:          0,
			},
			want: "adds	w13, w23, #291, lsl #12",
			wantErr: false,
		},
		{
			name: "cmn	w2, #4095",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xfc, 0x3f, 0x31}),
				address:          0,
			},
			want: "cmn	w2, #4095",
			wantErr: false,
		},
		{
			name: "adds	w20, wsp, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x00, 0x31}),
				address:          0,
			},
			want: "adds	w20, wsp, #0",
			wantErr: false,
		},
		{
			name: "cmn	x3, #1, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x04, 0x40, 0xb1}),
				address:          0,
			},
			want: "cmn	x3, #1, lsl #12",
			wantErr: false,
		},
		{
			name: "cmp	sp, #20, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x40, 0xf1}),
				address:          0,
			},
			want: "cmp	sp, #20, lsl #12",
			wantErr: false,
		},
		{
			name: "cmp	x30, #4095",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xff, 0x3f, 0xf1}),
				address:          0,
			},
			want: "cmp	x30, #4095",
			wantErr: false,
		},
		{
			name: "subs	x4, sp, #3822",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe4, 0xbb, 0x3b, 0xf1}),
				address:          0,
			},
			want: "subs	x4, sp, #3822",
			wantErr: false,
		},
		{
			name: "cmn	w3, #291, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x8c, 0x44, 0x31}),
				address:          0,
			},
			want: "cmn	w3, #291, lsl #12",
			wantErr: false,
		},
		{
			name: "cmn	wsp, #1365",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x57, 0x15, 0x31}),
				address:          0,
			},
			want: "cmn	wsp, #1365",
			wantErr: false,
		},
		{
			name: "cmn	sp, #1092, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x13, 0x51, 0xb1}),
				address:          0,
			},
			want: "cmn	sp, #1092, lsl #12",
			wantErr: false,
		},
		{
			name: "cmp	x4, #300, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xb0, 0x44, 0xf1}),
				address:          0,
			},
			want: "cmp	x4, #300, lsl #12",
			wantErr: false,
		},
		{
			name: "cmp	wsp, #500",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xd3, 0x07, 0x71}),
				address:          0,
			},
			want: "cmp	wsp, #500",
			wantErr: false,
		},
		{
			name: "cmp	sp, #200",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x23, 0x03, 0xf1}),
				address:          0,
			},
			want: "cmp	sp, #200",
			wantErr: false,
		},
		{
			name: "mov	sp, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x03, 0x00, 0x91}),
				address:          0,
			},
			want: "mov	sp, x30",
			wantErr: false,
		},
		{
			name: "mov	wsp, w20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x02, 0x00, 0x11}),
				address:          0,
			},
			want: "mov	wsp, w20",
			wantErr: false,
		},
		{
			name: "mov	x11, sp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xeb, 0x03, 0x00, 0x91}),
				address:          0,
			},
			want: "mov	x11, sp",
			wantErr: false,
		},
		{
			name: "mov	w24, wsp",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x03, 0x00, 0x11}),
				address:          0,
			},
			want: "mov	w24, wsp",
			wantErr: false,
		},
		{
			name: "add	w3, w5, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0x0b}),
				address:          0,
			},
			want: "add	w3, w5, w7",
			wantErr: false,
		},
		{
			name: "add	wzr, w3, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0x0b}),
				address:          0,
			},
			want: "add	wzr, w3, w5",
			wantErr: false,
		},
		{
			name: "add	w20, wzr, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0x0b}),
				address:          0,
			},
			want: "add	w20, wzr, w4",
			wantErr: false,
		},
		{
			name: "add	w4, w6, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0x0b}),
				address:          0,
			},
			want: "add	w4, w6, wzr",
			wantErr: false,
		},
		{
			name: "add	w11, w13, w15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0x0b}),
				address:          0,
			},
			want: "add	w11, w13, w15",
			wantErr: false,
		},
		{
			name: "add	w9, w3, wzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0x0b}),
				address:          0,
			},
			want: "add	w9, w3, wzr, lsl #10",
			wantErr: false,
		},
		{
			name: "add	w17, w29, w20, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0x7f, 0x14, 0x0b}),
				address:          0,
			},
			want: "add	w17, w29, w20, lsl #31",
			wantErr: false,
		},
		{
			name: "add	w17, w29, w20, lsl #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0x77, 0x14, 0x0b}),
				address:          0,
			},
			want: "add	w17, w29, w20, lsl #29",
			wantErr: false,
		},
		{
			name: "add	w21, w22, w23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0x0b}),
				address:          0,
			},
			want: "add	w21, w22, w23, lsr #0",
			wantErr: false,
		},
		{
			name: "add	w24, w25, w26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0x0b}),
				address:          0,
			},
			want: "add	w24, w25, w26, lsr #18",
			wantErr: false,
		},
		{
			name: "add	w27, w28, w29, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x7f, 0x5d, 0x0b}),
				address:          0,
			},
			want: "add	w27, w28, w29, lsr #31",
			wantErr: false,
		},
		{
			name: "add	w27, w28, w29, lsr #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x77, 0x5d, 0x0b}),
				address:          0,
			},
			want: "add	w27, w28, w29, lsr #29",
			wantErr: false,
		},
		{
			name: "add	w2, w3, w4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0x0b}),
				address:          0,
			},
			want: "add	w2, w3, w4, asr #0",
			wantErr: false,
		},
		{
			name: "add	w5, w6, w7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0x0b}),
				address:          0,
			},
			want: "add	w5, w6, w7, asr #21",
			wantErr: false,
		},
		{
			name: "add	w8, w9, w10, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x7d, 0x8a, 0x0b}),
				address:          0,
			},
			want: "add	w8, w9, w10, asr #31",
			wantErr: false,
		},
		{
			name: "add	w8, w9, w10, asr #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x75, 0x8a, 0x0b}),
				address:          0,
			},
			want: "add	w8, w9, w10, asr #29",
			wantErr: false,
		},
		{
			name: "add	x3, x5, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0x8b}),
				address:          0,
			},
			want: "add	x3, x5, x7",
			wantErr: false,
		},
		{
			name: "add	xzr, x3, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0x8b}),
				address:          0,
			},
			want: "add	xzr, x3, x5",
			wantErr: false,
		},
		{
			name: "add	x20, xzr, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0x8b}),
				address:          0,
			},
			want: "add	x20, xzr, x4",
			wantErr: false,
		},
		{
			name: "add	x4, x6, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0x8b}),
				address:          0,
			},
			want: "add	x4, x6, xzr",
			wantErr: false,
		},
		{
			name: "add	x11, x13, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0x8b}),
				address:          0,
			},
			want: "add	x11, x13, x15",
			wantErr: false,
		},
		{
			name: "add	x9, x3, xzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0x8b}),
				address:          0,
			},
			want: "add	x9, x3, xzr, lsl #10",
			wantErr: false,
		},
		{
			name: "add	x17, x29, x20, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xff, 0x14, 0x8b}),
				address:          0,
			},
			want: "add	x17, x29, x20, lsl #63",
			wantErr: false,
		},
		{
			name: "add	x17, x29, x20, lsl #58",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xeb, 0x14, 0x8b}),
				address:          0,
			},
			want: "add	x17, x29, x20, lsl #58",
			wantErr: false,
		},
		{
			name: "add	x21, x22, x23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0x8b}),
				address:          0,
			},
			want: "add	x21, x22, x23, lsr #0",
			wantErr: false,
		},
		{
			name: "add	x24, x25, x26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0x8b}),
				address:          0,
			},
			want: "add	x24, x25, x26, lsr #18",
			wantErr: false,
		},
		{
			name: "add	x27, x28, x29, lsr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0xff, 0x5d, 0x8b}),
				address:          0,
			},
			want: "add	x27, x28, x29, lsr #63",
			wantErr: false,
		},
		{
			name: "add	x17, x29, x20, lsr #58",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xeb, 0x54, 0x8b}),
				address:          0,
			},
			want: "add	x17, x29, x20, lsr #58",
			wantErr: false,
		},
		{
			name: "add	x2, x3, x4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0x8b}),
				address:          0,
			},
			want: "add	x2, x3, x4, asr #0",
			wantErr: false,
		},
		{
			name: "add	x5, x6, x7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0x8b}),
				address:          0,
			},
			want: "add	x5, x6, x7, asr #21",
			wantErr: false,
		},
		{
			name: "add	x8, x9, x10, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0xfd, 0x8a, 0x8b}),
				address:          0,
			},
			want: "add	x8, x9, x10, asr #63",
			wantErr: false,
		},
		{
			name: "add	x17, x29, x20, asr #58",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xeb, 0x94, 0x8b}),
				address:          0,
			},
			want: "add	x17, x29, x20, asr #58",
			wantErr: false,
		},
		{
			name: "adds	w3, w5, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0x2b}),
				address:          0,
			},
			want: "adds	w3, w5, w7",
			wantErr: false,
		},
		{
			name: "cmn	w3, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0x2b}),
				address:          0,
			},
			want: "cmn	w3, w5",
			wantErr: false,
		},
		{
			name: "adds	w20, wzr, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0x2b}),
				address:          0,
			},
			want: "adds	w20, wzr, w4",
			wantErr: false,
		},
		{
			name: "adds	w4, w6, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0x2b}),
				address:          0,
			},
			want: "adds	w4, w6, wzr",
			wantErr: false,
		},
		{
			name: "adds	w11, w13, w15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0x2b}),
				address:          0,
			},
			want: "adds	w11, w13, w15",
			wantErr: false,
		},
		{
			name: "adds	w9, w3, wzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0x2b}),
				address:          0,
			},
			want: "adds	w9, w3, wzr, lsl #10",
			wantErr: false,
		},
		{
			name: "adds	w17, w29, w20, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0x7f, 0x14, 0x2b}),
				address:          0,
			},
			want: "adds	w17, w29, w20, lsl #31",
			wantErr: false,
		},
		{
			name: "adds	w21, w22, w23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0x2b}),
				address:          0,
			},
			want: "adds	w21, w22, w23, lsr #0",
			wantErr: false,
		},
		{
			name: "adds	w24, w25, w26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0x2b}),
				address:          0,
			},
			want: "adds	w24, w25, w26, lsr #18",
			wantErr: false,
		},
		{
			name: "adds	w27, w28, w29, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x7f, 0x5d, 0x2b}),
				address:          0,
			},
			want: "adds	w27, w28, w29, lsr #31",
			wantErr: false,
		},
		{
			name: "adds	w2, w3, w4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0x2b}),
				address:          0,
			},
			want: "adds	w2, w3, w4, asr #0",
			wantErr: false,
		},
		{
			name: "adds	w5, w6, w7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0x2b}),
				address:          0,
			},
			want: "adds	w5, w6, w7, asr #21",
			wantErr: false,
		},
		{
			name: "adds	w8, w9, w10, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x7d, 0x8a, 0x2b}),
				address:          0,
			},
			want: "adds	w8, w9, w10, asr #31",
			wantErr: false,
		},
		{
			name: "adds	x3, x5, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0xab}),
				address:          0,
			},
			want: "adds	x3, x5, x7",
			wantErr: false,
		},
		{
			name: "cmn	x3, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0xab}),
				address:          0,
			},
			want: "cmn	x3, x5",
			wantErr: false,
		},
		{
			name: "adds	x20, xzr, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0xab}),
				address:          0,
			},
			want: "adds	x20, xzr, x4",
			wantErr: false,
		},
		{
			name: "adds	x4, x6, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0xab}),
				address:          0,
			},
			want: "adds	x4, x6, xzr",
			wantErr: false,
		},
		{
			name: "adds	x11, x13, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0xab}),
				address:          0,
			},
			want: "adds	x11, x13, x15",
			wantErr: false,
		},
		{
			name: "adds	x9, x3, xzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0xab}),
				address:          0,
			},
			want: "adds	x9, x3, xzr, lsl #10",
			wantErr: false,
		},
		{
			name: "adds	x17, x29, x20, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xff, 0x14, 0xab}),
				address:          0,
			},
			want: "adds	x17, x29, x20, lsl #63",
			wantErr: false,
		},
		{
			name: "adds	x21, x22, x23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0xab}),
				address:          0,
			},
			want: "adds	x21, x22, x23, lsr #0",
			wantErr: false,
		},
		{
			name: "adds	x24, x25, x26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0xab}),
				address:          0,
			},
			want: "adds	x24, x25, x26, lsr #18",
			wantErr: false,
		},
		{
			name: "adds	x27, x28, x29, lsr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0xff, 0x5d, 0xab}),
				address:          0,
			},
			want: "adds	x27, x28, x29, lsr #63",
			wantErr: false,
		},
		{
			name: "adds	x2, x3, x4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0xab}),
				address:          0,
			},
			want: "adds	x2, x3, x4, asr #0",
			wantErr: false,
		},
		{
			name: "adds	x5, x6, x7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0xab}),
				address:          0,
			},
			want: "adds	x5, x6, x7, asr #21",
			wantErr: false,
		},
		{
			name: "adds	x8, x9, x10, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0xfd, 0x8a, 0xab}),
				address:          0,
			},
			want: "adds	x8, x9, x10, asr #63",
			wantErr: false,
		},
		{
			name: "sub	w3, w5, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0x4b}),
				address:          0,
			},
			want: "sub	w3, w5, w7",
			wantErr: false,
		},
		{
			name: "sub	wzr, w3, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0x4b}),
				address:          0,
			},
			want: "sub	wzr, w3, w5",
			wantErr: false,
		},
		{
			name: "neg	w20, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0x4b}),
				address:          0,
			},
			want: "neg	w20, w4",
			wantErr: false,
		},
		{
			name: "sub	w4, w6, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0x4b}),
				address:          0,
			},
			want: "sub	w4, w6, wzr",
			wantErr: false,
		},
		{
			name: "sub	w11, w13, w15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0x4b}),
				address:          0,
			},
			want: "sub	w11, w13, w15",
			wantErr: false,
		},
		{
			name: "sub	w9, w3, wzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0x4b}),
				address:          0,
			},
			want: "sub	w9, w3, wzr, lsl #10",
			wantErr: false,
		},
		{
			name: "sub	w17, w29, w20, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0x7f, 0x14, 0x4b}),
				address:          0,
			},
			want: "sub	w17, w29, w20, lsl #31",
			wantErr: false,
		},
		{
			name: "sub	w21, w22, w23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0x4b}),
				address:          0,
			},
			want: "sub	w21, w22, w23, lsr #0",
			wantErr: false,
		},
		{
			name: "sub	w24, w25, w26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0x4b}),
				address:          0,
			},
			want: "sub	w24, w25, w26, lsr #18",
			wantErr: false,
		},
		{
			name: "sub	w27, w28, w29, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x7f, 0x5d, 0x4b}),
				address:          0,
			},
			want: "sub	w27, w28, w29, lsr #31",
			wantErr: false,
		},
		{
			name: "sub	w2, w3, w4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0x4b}),
				address:          0,
			},
			want: "sub	w2, w3, w4, asr #0",
			wantErr: false,
		},
		{
			name: "sub	w5, w6, w7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0x4b}),
				address:          0,
			},
			want: "sub	w5, w6, w7, asr #21",
			wantErr: false,
		},
		{
			name: "sub	w8, w9, w10, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x7d, 0x8a, 0x4b}),
				address:          0,
			},
			want: "sub	w8, w9, w10, asr #31",
			wantErr: false,
		},
		{
			name: "sub	x3, x5, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0xcb}),
				address:          0,
			},
			want: "sub	x3, x5, x7",
			wantErr: false,
		},
		{
			name: "sub	xzr, x3, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0xcb}),
				address:          0,
			},
			want: "sub	xzr, x3, x5",
			wantErr: false,
		},
		{
			name: "neg	x20, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0xcb}),
				address:          0,
			},
			want: "neg	x20, x4",
			wantErr: false,
		},
		{
			name: "sub	x4, x6, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0xcb}),
				address:          0,
			},
			want: "sub	x4, x6, xzr",
			wantErr: false,
		},
		{
			name: "sub	x11, x13, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0xcb}),
				address:          0,
			},
			want: "sub	x11, x13, x15",
			wantErr: false,
		},
		{
			name: "sub	x9, x3, xzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0xcb}),
				address:          0,
			},
			want: "sub	x9, x3, xzr, lsl #10",
			wantErr: false,
		},
		{
			name: "sub	x17, x29, x20, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xff, 0x14, 0xcb}),
				address:          0,
			},
			want: "sub	x17, x29, x20, lsl #63",
			wantErr: false,
		},
		{
			name: "sub	x21, x22, x23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0xcb}),
				address:          0,
			},
			want: "sub	x21, x22, x23, lsr #0",
			wantErr: false,
		},
		{
			name: "sub	x24, x25, x26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0xcb}),
				address:          0,
			},
			want: "sub	x24, x25, x26, lsr #18",
			wantErr: false,
		},
		{
			name: "sub	x27, x28, x29, lsr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0xff, 0x5d, 0xcb}),
				address:          0,
			},
			want: "sub	x27, x28, x29, lsr #63",
			wantErr: false,
		},
		{
			name: "sub	x2, x3, x4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0xcb}),
				address:          0,
			},
			want: "sub	x2, x3, x4, asr #0",
			wantErr: false,
		},
		{
			name: "sub	x5, x6, x7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0xcb}),
				address:          0,
			},
			want: "sub	x5, x6, x7, asr #21",
			wantErr: false,
		},
		{
			name: "sub	x8, x9, x10, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0xfd, 0x8a, 0xcb}),
				address:          0,
			},
			want: "sub	x8, x9, x10, asr #63",
			wantErr: false,
		},
		{
			name: "subs	w3, w5, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0x6b}),
				address:          0,
			},
			want: "subs	w3, w5, w7",
			wantErr: false,
		},
		{
			name: "cmp	w3, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0x6b}),
				address:          0,
			},
			want: "cmp	w3, w5",
			wantErr: false,
		},
		{
			name: "negs	w20, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0x6b}),
				address:          0,
			},
			want: "negs	w20, w4",
			wantErr: false,
		},
		{
			name: "subs	w4, w6, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0x6b}),
				address:          0,
			},
			want: "subs	w4, w6, wzr",
			wantErr: false,
		},
		{
			name: "subs	w11, w13, w15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0x6b}),
				address:          0,
			},
			want: "subs	w11, w13, w15",
			wantErr: false,
		},
		{
			name: "subs	w9, w3, wzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0x6b}),
				address:          0,
			},
			want: "subs	w9, w3, wzr, lsl #10",
			wantErr: false,
		},
		{
			name: "subs	w17, w29, w20, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0x7f, 0x14, 0x6b}),
				address:          0,
			},
			want: "subs	w17, w29, w20, lsl #31",
			wantErr: false,
		},
		{
			name: "subs	w21, w22, w23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0x6b}),
				address:          0,
			},
			want: "subs	w21, w22, w23, lsr #0",
			wantErr: false,
		},
		{
			name: "subs	w24, w25, w26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0x6b}),
				address:          0,
			},
			want: "subs	w24, w25, w26, lsr #18",
			wantErr: false,
		},
		{
			name: "subs	w27, w28, w29, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x7f, 0x5d, 0x6b}),
				address:          0,
			},
			want: "subs	w27, w28, w29, lsr #31",
			wantErr: false,
		},
		{
			name: "subs	w2, w3, w4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0x6b}),
				address:          0,
			},
			want: "subs	w2, w3, w4, asr #0",
			wantErr: false,
		},
		{
			name: "subs	w5, w6, w7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0x6b}),
				address:          0,
			},
			want: "subs	w5, w6, w7, asr #21",
			wantErr: false,
		},
		{
			name: "subs	w8, w9, w10, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x7d, 0x8a, 0x6b}),
				address:          0,
			},
			want: "subs	w8, w9, w10, asr #31",
			wantErr: false,
		},
		{
			name: "subs	x3, x5, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x07, 0xeb}),
				address:          0,
			},
			want: "subs	x3, x5, x7",
			wantErr: false,
		},
		{
			name: "cmp	x3, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x05, 0xeb}),
				address:          0,
			},
			want: "cmp	x3, x5",
			wantErr: false,
		},
		{
			name: "negs	x20, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x04, 0xeb}),
				address:          0,
			},
			want: "negs	x20, x4",
			wantErr: false,
		},
		{
			name: "subs	x4, x6, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc4, 0x00, 0x1f, 0xeb}),
				address:          0,
			},
			want: "subs	x4, x6, xzr",
			wantErr: false,
		},
		{
			name: "subs	x11, x13, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x01, 0x0f, 0xeb}),
				address:          0,
			},
			want: "subs	x11, x13, x15",
			wantErr: false,
		},
		{
			name: "subs	x9, x3, xzr, lsl #10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x28, 0x1f, 0xeb}),
				address:          0,
			},
			want: "subs	x9, x3, xzr, lsl #10",
			wantErr: false,
		},
		{
			name: "subs	x17, x29, x20, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb1, 0xff, 0x14, 0xeb}),
				address:          0,
			},
			want: "subs	x17, x29, x20, lsl #63",
			wantErr: false,
		},
		{
			name: "subs	x21, x22, x23, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x57, 0xeb}),
				address:          0,
			},
			want: "subs	x21, x22, x23, lsr #0",
			wantErr: false,
		},
		{
			name: "subs	x24, x25, x26, lsr #18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x4b, 0x5a, 0xeb}),
				address:          0,
			},
			want: "subs	x24, x25, x26, lsr #18",
			wantErr: false,
		},
		{
			name: "subs	x27, x28, x29, lsr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0xff, 0x5d, 0xeb}),
				address:          0,
			},
			want: "subs	x27, x28, x29, lsr #63",
			wantErr: false,
		},
		{
			name: "subs	x2, x3, x4, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x84, 0xeb}),
				address:          0,
			},
			want: "subs	x2, x3, x4, asr #0",
			wantErr: false,
		},
		{
			name: "subs	x5, x6, x7, asr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x54, 0x87, 0xeb}),
				address:          0,
			},
			want: "subs	x5, x6, x7, asr #21",
			wantErr: false,
		},
		{
			name: "subs	x8, x9, x10, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0xfd, 0x8a, 0xeb}),
				address:          0,
			},
			want: "subs	x8, x9, x10, asr #63",
			wantErr: false,
		},
		{
			name: "cmn	w0, w3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x03, 0x2b}),
				address:          0,
			},
			want: "cmn	w0, w3",
			wantErr: false,
		},
		{
			name: "cmn	wzr, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x04, 0x2b}),
				address:          0,
			},
			want: "cmn	wzr, w4",
			wantErr: false,
		},
		{
			name: "cmn	w5, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x1f, 0x2b}),
				address:          0,
			},
			want: "cmn	w5, wzr",
			wantErr: false,
		},
		{
			name: "cmn	wsp, w6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x26, 0x2b}),
				address:          0,
			},
			want: "cmn	wsp, w6",
			wantErr: false,
		},
		{
			name: "cmn	w6, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x00, 0x07, 0x2b}),
				address:          0,
			},
			want: "cmn	w6, w7",
			wantErr: false,
		},
		{
			name: "cmn	w8, w9, lsl #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x3d, 0x09, 0x2b}),
				address:          0,
			},
			want: "cmn	w8, w9, lsl #15",
			wantErr: false,
		},
		{
			name: "cmn	w10, w11, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x7d, 0x0b, 0x2b}),
				address:          0,
			},
			want: "cmn	w10, w11, lsl #31",
			wantErr: false,
		},
		{
			name: "cmn	w12, w13, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x4d, 0x2b}),
				address:          0,
			},
			want: "cmn	w12, w13, lsr #0",
			wantErr: false,
		},
		{
			name: "cmn	w14, w15, lsr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x55, 0x4f, 0x2b}),
				address:          0,
			},
			want: "cmn	w14, w15, lsr #21",
			wantErr: false,
		},
		{
			name: "cmn	w16, w17, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x7e, 0x51, 0x2b}),
				address:          0,
			},
			want: "cmn	w16, w17, lsr #31",
			wantErr: false,
		},
		{
			name: "cmn	w18, w19, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x02, 0x93, 0x2b}),
				address:          0,
			},
			want: "cmn	w18, w19, asr #0",
			wantErr: false,
		},
		{
			name: "cmn	w20, w21, asr #22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x5a, 0x95, 0x2b}),
				address:          0,
			},
			want: "cmn	w20, w21, asr #22",
			wantErr: false,
		},
		{
			name: "cmn	w22, w23, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x7e, 0x97, 0x2b}),
				address:          0,
			},
			want: "cmn	w22, w23, asr #31",
			wantErr: false,
		},
		{
			name: "cmn	x0, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x03, 0xab}),
				address:          0,
			},
			want: "cmn	x0, x3",
			wantErr: false,
		},
		{
			name: "cmn	xzr, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x04, 0xab}),
				address:          0,
			},
			want: "cmn	xzr, x4",
			wantErr: false,
		},
		{
			name: "cmn	x5, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x1f, 0xab}),
				address:          0,
			},
			want: "cmn	x5, xzr",
			wantErr: false,
		},
		{
			name: "cmn	sp, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x26, 0xab}),
				address:          0,
			},
			want: "cmn	sp, x6",
			wantErr: false,
		},
		{
			name: "cmn	x6, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x00, 0x07, 0xab}),
				address:          0,
			},
			want: "cmn	x6, x7",
			wantErr: false,
		},
		{
			name: "cmn	x8, x9, lsl #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x3d, 0x09, 0xab}),
				address:          0,
			},
			want: "cmn	x8, x9, lsl #15",
			wantErr: false,
		},
		{
			name: "cmn	x10, x11, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xfd, 0x0b, 0xab}),
				address:          0,
			},
			want: "cmn	x10, x11, lsl #63",
			wantErr: false,
		},
		{
			name: "cmn	x12, x13, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x4d, 0xab}),
				address:          0,
			},
			want: "cmn	x12, x13, lsr #0",
			wantErr: false,
		},
		{
			name: "cmn	x14, x15, lsr #41",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xa5, 0x4f, 0xab}),
				address:          0,
			},
			want: "cmn	x14, x15, lsr #41",
			wantErr: false,
		},
		{
			name: "cmn	x16, x17, lsr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0xfe, 0x51, 0xab}),
				address:          0,
			},
			want: "cmn	x16, x17, lsr #63",
			wantErr: false,
		},
		{
			name: "cmn	x18, x19, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x02, 0x93, 0xab}),
				address:          0,
			},
			want: "cmn	x18, x19, asr #0",
			wantErr: false,
		},
		{
			name: "cmn	x20, x21, asr #55",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xde, 0x95, 0xab}),
				address:          0,
			},
			want: "cmn	x20, x21, asr #55",
			wantErr: false,
		},
		{
			name: "cmn	x22, x23, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xfe, 0x97, 0xab}),
				address:          0,
			},
			want: "cmn	x22, x23, asr #63",
			wantErr: false,
		},
		{
			name: "cmp	w0, w3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x03, 0x6b}),
				address:          0,
			},
			want: "cmp	w0, w3",
			wantErr: false,
		},
		{
			name: "cmp	wzr, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x04, 0x6b}),
				address:          0,
			},
			want: "cmp	wzr, w4",
			wantErr: false,
		},
		{
			name: "cmp	w5, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x1f, 0x6b}),
				address:          0,
			},
			want: "cmp	w5, wzr",
			wantErr: false,
		},
		{
			name: "cmp	wsp, w6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x43, 0x26, 0x6b}),
				address:          0,
			},
			want: "cmp	wsp, w6",
			wantErr: false,
		},
		{
			name: "cmp	w6, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x00, 0x07, 0x6b}),
				address:          0,
			},
			want: "cmp	w6, w7",
			wantErr: false,
		},
		{
			name: "cmp	w8, w9, lsl #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x3d, 0x09, 0x6b}),
				address:          0,
			},
			want: "cmp	w8, w9, lsl #15",
			wantErr: false,
		},
		{
			name: "cmp	w10, w11, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x7d, 0x0b, 0x6b}),
				address:          0,
			},
			want: "cmp	w10, w11, lsl #31",
			wantErr: false,
		},
		{
			name: "cmp	w12, w13, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x4d, 0x6b}),
				address:          0,
			},
			want: "cmp	w12, w13, lsr #0",
			wantErr: false,
		},
		{
			name: "cmp	w14, w15, lsr #21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x55, 0x4f, 0x6b}),
				address:          0,
			},
			want: "cmp	w14, w15, lsr #21",
			wantErr: false,
		},
		{
			name: "cmp	w16, w17, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x7e, 0x51, 0x6b}),
				address:          0,
			},
			want: "cmp	w16, w17, lsr #31",
			wantErr: false,
		},
		{
			name: "cmp	w18, w19, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x02, 0x93, 0x6b}),
				address:          0,
			},
			want: "cmp	w18, w19, asr #0",
			wantErr: false,
		},
		{
			name: "cmp	w20, w21, asr #22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x5a, 0x95, 0x6b}),
				address:          0,
			},
			want: "cmp	w20, w21, asr #22",
			wantErr: false,
		},
		{
			name: "cmp	w22, w23, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x7e, 0x97, 0x6b}),
				address:          0,
			},
			want: "cmp	w22, w23, asr #31",
			wantErr: false,
		},
		{
			name: "cmp	x0, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x03, 0xeb}),
				address:          0,
			},
			want: "cmp	x0, x3",
			wantErr: false,
		},
		{
			name: "cmp	xzr, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x04, 0xeb}),
				address:          0,
			},
			want: "cmp	xzr, x4",
			wantErr: false,
		},
		{
			name: "cmp	x5, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x1f, 0xeb}),
				address:          0,
			},
			want: "cmp	x5, xzr",
			wantErr: false,
		},
		{
			name: "cmp	sp, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x63, 0x26, 0xeb}),
				address:          0,
			},
			want: "cmp	sp, x6",
			wantErr: false,
		},
		{
			name: "cmp	x6, x7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x00, 0x07, 0xeb}),
				address:          0,
			},
			want: "cmp	x6, x7",
			wantErr: false,
		},
		{
			name: "cmp	x8, x9, lsl #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x3d, 0x09, 0xeb}),
				address:          0,
			},
			want: "cmp	x8, x9, lsl #15",
			wantErr: false,
		},
		{
			name: "cmp	x10, x11, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xfd, 0x0b, 0xeb}),
				address:          0,
			},
			want: "cmp	x10, x11, lsl #63",
			wantErr: false,
		},
		{
			name: "cmp	x12, x13, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x4d, 0xeb}),
				address:          0,
			},
			want: "cmp	x12, x13, lsr #0",
			wantErr: false,
		},
		{
			name: "cmp	x14, x15, lsr #41",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xa5, 0x4f, 0xeb}),
				address:          0,
			},
			want: "cmp	x14, x15, lsr #41",
			wantErr: false,
		},
		{
			name: "cmp	x16, x17, lsr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0xfe, 0x51, 0xeb}),
				address:          0,
			},
			want: "cmp	x16, x17, lsr #63",
			wantErr: false,
		},
		{
			name: "cmp	x18, x19, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x02, 0x93, 0xeb}),
				address:          0,
			},
			want: "cmp	x18, x19, asr #0",
			wantErr: false,
		},
		{
			name: "cmp	x20, x21, asr #55",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xde, 0x95, 0xeb}),
				address:          0,
			},
			want: "cmp	x20, x21, asr #55",
			wantErr: false,
		},
		{
			name: "cmp	x22, x23, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xfe, 0x97, 0xeb}),
				address:          0,
			},
			want: "cmp	x22, x23, asr #63",
			wantErr: false,
		},
		{
			name: "neg	w29, w30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x03, 0x1e, 0x4b}),
				address:          0,
			},
			want: "neg	w29, w30",
			wantErr: false,
		},
		{
			name: "neg	w30, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x03, 0x1f, 0x4b}),
				address:          0,
			},
			want: "neg	w30, wzr",
			wantErr: false,
		},
		{
			name: "neg	wzr, w0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x00, 0x4b}),
				address:          0,
			},
			want: "neg	wzr, w0",
			wantErr: false,
		},
		{
			name: "neg	w28, w27",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfc, 0x03, 0x1b, 0x4b}),
				address:          0,
			},
			want: "neg	w28, w27",
			wantErr: false,
		},
		{
			name: "neg	w26, w25, lsl #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0x77, 0x19, 0x4b}),
				address:          0,
			},
			want: "neg	w26, w25, lsl #29",
			wantErr: false,
		},
		{
			name: "neg	w24, w23, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x7f, 0x17, 0x4b}),
				address:          0,
			},
			want: "neg	w24, w23, lsl #31",
			wantErr: false,
		},
		{
			name: "neg	w22, w21, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x03, 0x55, 0x4b}),
				address:          0,
			},
			want: "neg	w22, w21, lsr #0",
			wantErr: false,
		},
		{
			name: "neg	w20, w19, lsr #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x07, 0x53, 0x4b}),
				address:          0,
			},
			want: "neg	w20, w19, lsr #1",
			wantErr: false,
		},
		{
			name: "neg	w18, w17, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0x7f, 0x51, 0x4b}),
				address:          0,
			},
			want: "neg	w18, w17, lsr #31",
			wantErr: false,
		},
		{
			name: "neg	w16, w15, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf0, 0x03, 0x8f, 0x4b}),
				address:          0,
			},
			want: "neg	w16, w15, asr #0",
			wantErr: false,
		},
		{
			name: "neg	w14, w13, asr #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x33, 0x8d, 0x4b}),
				address:          0,
			},
			want: "neg	w14, w13, asr #12",
			wantErr: false,
		},
		{
			name: "neg	w12, w11, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x7f, 0x8b, 0x4b}),
				address:          0,
			},
			want: "neg	w12, w11, asr #31",
			wantErr: false,
		},
		{
			name: "neg	x29, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x03, 0x1e, 0xcb}),
				address:          0,
			},
			want: "neg	x29, x30",
			wantErr: false,
		},
		{
			name: "neg	x30, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x03, 0x1f, 0xcb}),
				address:          0,
			},
			want: "neg	x30, xzr",
			wantErr: false,
		},
		{
			name: "neg	xzr, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x00, 0xcb}),
				address:          0,
			},
			want: "neg	xzr, x0",
			wantErr: false,
		},
		{
			name: "neg	x28, x27",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfc, 0x03, 0x1b, 0xcb}),
				address:          0,
			},
			want: "neg	x28, x27",
			wantErr: false,
		},
		{
			name: "neg	x26, x25, lsl #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0x77, 0x19, 0xcb}),
				address:          0,
			},
			want: "neg	x26, x25, lsl #29",
			wantErr: false,
		},
		{
			name: "neg	x24, x23, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x7f, 0x17, 0xcb}),
				address:          0,
			},
			want: "neg	x24, x23, lsl #31",
			wantErr: false,
		},
		{
			name: "neg	x22, x21, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x03, 0x55, 0xcb}),
				address:          0,
			},
			want: "neg	x22, x21, lsr #0",
			wantErr: false,
		},
		{
			name: "neg	x20, x19, lsr #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x07, 0x53, 0xcb}),
				address:          0,
			},
			want: "neg	x20, x19, lsr #1",
			wantErr: false,
		},
		{
			name: "neg	x18, x17, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0x7f, 0x51, 0xcb}),
				address:          0,
			},
			want: "neg	x18, x17, lsr #31",
			wantErr: false,
		},
		{
			name: "neg	x16, x15, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf0, 0x03, 0x8f, 0xcb}),
				address:          0,
			},
			want: "neg	x16, x15, asr #0",
			wantErr: false,
		},
		{
			name: "neg	x14, x13, asr #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x33, 0x8d, 0xcb}),
				address:          0,
			},
			want: "neg	x14, x13, asr #12",
			wantErr: false,
		},
		{
			name: "neg	x12, x11, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x7f, 0x8b, 0xcb}),
				address:          0,
			},
			want: "neg	x12, x11, asr #31",
			wantErr: false,
		},
		{
			name: "negs	w29, w30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x03, 0x1e, 0x6b}),
				address:          0,
			},
			want: "negs	w29, w30",
			wantErr: false,
		},
		{
			name: "negs	w30, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x03, 0x1f, 0x6b}),
				address:          0,
			},
			want: "negs	w30, wzr",
			wantErr: false,
		},
		{
			name: "cmp	wzr, w0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x00, 0x6b}),
				address:          0,
			},
			want: "cmp	wzr, w0",
			wantErr: false,
		},
		{
			name: "negs	w28, w27",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfc, 0x03, 0x1b, 0x6b}),
				address:          0,
			},
			want: "negs	w28, w27",
			wantErr: false,
		},
		{
			name: "negs	w26, w25, lsl #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0x77, 0x19, 0x6b}),
				address:          0,
			},
			want: "negs	w26, w25, lsl #29",
			wantErr: false,
		},
		{
			name: "negs	w24, w23, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x7f, 0x17, 0x6b}),
				address:          0,
			},
			want: "negs	w24, w23, lsl #31",
			wantErr: false,
		},
		{
			name: "negs	w22, w21, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x03, 0x55, 0x6b}),
				address:          0,
			},
			want: "negs	w22, w21, lsr #0",
			wantErr: false,
		},
		{
			name: "negs	w20, w19, lsr #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x07, 0x53, 0x6b}),
				address:          0,
			},
			want: "negs	w20, w19, lsr #1",
			wantErr: false,
		},
		{
			name: "negs	w18, w17, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0x7f, 0x51, 0x6b}),
				address:          0,
			},
			want: "negs	w18, w17, lsr #31",
			wantErr: false,
		},
		{
			name: "negs	w16, w15, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf0, 0x03, 0x8f, 0x6b}),
				address:          0,
			},
			want: "negs	w16, w15, asr #0",
			wantErr: false,
		},
		{
			name: "negs	w14, w13, asr #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x33, 0x8d, 0x6b}),
				address:          0,
			},
			want: "negs	w14, w13, asr #12",
			wantErr: false,
		},
		{
			name: "negs	w12, w11, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x7f, 0x8b, 0x6b}),
				address:          0,
			},
			want: "negs	w12, w11, asr #31",
			wantErr: false,
		},
		{
			name: "negs	x29, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x03, 0x1e, 0xeb}),
				address:          0,
			},
			want: "negs	x29, x30",
			wantErr: false,
		},
		{
			name: "negs	x30, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x03, 0x1f, 0xeb}),
				address:          0,
			},
			want: "negs	x30, xzr",
			wantErr: false,
		},
		{
			name: "cmp	xzr, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x00, 0xeb}),
				address:          0,
			},
			want: "cmp	xzr, x0",
			wantErr: false,
		},
		{
			name: "negs	x28, x27",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfc, 0x03, 0x1b, 0xeb}),
				address:          0,
			},
			want: "negs	x28, x27",
			wantErr: false,
		},
		{
			name: "negs	x26, x25, lsl #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0x77, 0x19, 0xeb}),
				address:          0,
			},
			want: "negs	x26, x25, lsl #29",
			wantErr: false,
		},
		{
			name: "negs	x24, x23, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x7f, 0x17, 0xeb}),
				address:          0,
			},
			want: "negs	x24, x23, lsl #31",
			wantErr: false,
		},
		{
			name: "negs	x22, x21, lsr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x03, 0x55, 0xeb}),
				address:          0,
			},
			want: "negs	x22, x21, lsr #0",
			wantErr: false,
		},
		{
			name: "negs	x20, x19, lsr #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x07, 0x53, 0xeb}),
				address:          0,
			},
			want: "negs	x20, x19, lsr #1",
			wantErr: false,
		},
		{
			name: "negs	x18, x17, lsr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0x7f, 0x51, 0xeb}),
				address:          0,
			},
			want: "negs	x18, x17, lsr #31",
			wantErr: false,
		},
		{
			name: "negs	x16, x15, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf0, 0x03, 0x8f, 0xeb}),
				address:          0,
			},
			want: "negs	x16, x15, asr #0",
			wantErr: false,
		},
		{
			name: "negs	x14, x13, asr #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x33, 0x8d, 0xeb}),
				address:          0,
			},
			want: "negs	x14, x13, asr #12",
			wantErr: false,
		},
		{
			name: "negs	x12, x11, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x7f, 0x8b, 0xeb}),
				address:          0,
			},
			want: "negs	x12, x11, asr #31",
			wantErr: false,
		},
		{
			name: "adc	w29, w27, w25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0x1a}),
				address:          0,
			},
			want: "adc	w29, w27, w25",
			wantErr: false,
		},
		{
			name: "adc	wzr, w3, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0x1a}),
				address:          0,
			},
			want: "adc	wzr, w3, w4",
			wantErr: false,
		},
		{
			name: "adc	w9, wzr, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0x1a}),
				address:          0,
			},
			want: "adc	w9, wzr, w10",
			wantErr: false,
		},
		{
			name: "adc	w20, w0, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0x1a}),
				address:          0,
			},
			want: "adc	w20, w0, wzr",
			wantErr: false,
		},
		{
			name: "adc	x29, x27, x25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0x9a}),
				address:          0,
			},
			want: "adc	x29, x27, x25",
			wantErr: false,
		},
		{
			name: "adc	xzr, x3, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0x9a}),
				address:          0,
			},
			want: "adc	xzr, x3, x4",
			wantErr: false,
		},
		{
			name: "adc	x9, xzr, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0x9a}),
				address:          0,
			},
			want: "adc	x9, xzr, x10",
			wantErr: false,
		},
		{
			name: "adc	x20, x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0x9a}),
				address:          0,
			},
			want: "adc	x20, x0, xzr",
			wantErr: false,
		},
		{
			name: "adcs	w29, w27, w25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0x3a}),
				address:          0,
			},
			want: "adcs	w29, w27, w25",
			wantErr: false,
		},
		{
			name: "adcs	wzr, w3, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0x3a}),
				address:          0,
			},
			want: "adcs	wzr, w3, w4",
			wantErr: false,
		},
		{
			name: "adcs	w9, wzr, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0x3a}),
				address:          0,
			},
			want: "adcs	w9, wzr, w10",
			wantErr: false,
		},
		{
			name: "adcs	w20, w0, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0x3a}),
				address:          0,
			},
			want: "adcs	w20, w0, wzr",
			wantErr: false,
		},
		{
			name: "adcs	x29, x27, x25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0xba}),
				address:          0,
			},
			want: "adcs	x29, x27, x25",
			wantErr: false,
		},
		{
			name: "adcs	xzr, x3, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0xba}),
				address:          0,
			},
			want: "adcs	xzr, x3, x4",
			wantErr: false,
		},
		{
			name: "adcs	x9, xzr, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0xba}),
				address:          0,
			},
			want: "adcs	x9, xzr, x10",
			wantErr: false,
		},
		{
			name: "adcs	x20, x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0xba}),
				address:          0,
			},
			want: "adcs	x20, x0, xzr",
			wantErr: false,
		},
		{
			name: "sbc	w29, w27, w25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0x5a}),
				address:          0,
			},
			want: "sbc	w29, w27, w25",
			wantErr: false,
		},
		{
			name: "sbc	wzr, w3, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0x5a}),
				address:          0,
			},
			want: "sbc	wzr, w3, w4",
			wantErr: false,
		},
		{
			name: "ngc	w9, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0x5a}),
				address:          0,
			},
			want: "ngc	w9, w10",
			wantErr: false,
		},
		{
			name: "sbc	w20, w0, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0x5a}),
				address:          0,
			},
			want: "sbc	w20, w0, wzr",
			wantErr: false,
		},
		{
			name: "sbc	x29, x27, x25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0xda}),
				address:          0,
			},
			want: "sbc	x29, x27, x25",
			wantErr: false,
		},
		{
			name: "sbc	xzr, x3, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0xda}),
				address:          0,
			},
			want: "sbc	xzr, x3, x4",
			wantErr: false,
		},
		{
			name: "ngc	x9, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0xda}),
				address:          0,
			},
			want: "ngc	x9, x10",
			wantErr: false,
		},
		{
			name: "sbc	x20, x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0xda}),
				address:          0,
			},
			want: "sbc	x20, x0, xzr",
			wantErr: false,
		},
		{
			name: "sbcs	w29, w27, w25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0x7a}),
				address:          0,
			},
			want: "sbcs	w29, w27, w25",
			wantErr: false,
		},
		{
			name: "sbcs	wzr, w3, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0x7a}),
				address:          0,
			},
			want: "sbcs	wzr, w3, w4",
			wantErr: false,
		},
		{
			name: "ngcs	w9, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0x7a}),
				address:          0,
			},
			want: "ngcs	w9, w10",
			wantErr: false,
		},
		{
			name: "sbcs	w20, w0, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0x7a}),
				address:          0,
			},
			want: "sbcs	w20, w0, wzr",
			wantErr: false,
		},
		{
			name: "sbcs	x29, x27, x25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x03, 0x19, 0xfa}),
				address:          0,
			},
			want: "sbcs	x29, x27, x25",
			wantErr: false,
		},
		{
			name: "sbcs	xzr, x3, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x00, 0x04, 0xfa}),
				address:          0,
			},
			want: "sbcs	xzr, x3, x4",
			wantErr: false,
		},
		{
			name: "ngcs	x9, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x0a, 0xfa}),
				address:          0,
			},
			want: "ngcs	x9, x10",
			wantErr: false,
		},
		{
			name: "sbcs	x20, x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x1f, 0xfa}),
				address:          0,
			},
			want: "sbcs	x20, x0, xzr",
			wantErr: false,
		},
		{
			name: "ngc	w3, w12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x0c, 0x5a}),
				address:          0,
			},
			want: "ngc	w3, w12",
			wantErr: false,
		},
		{
			name: "ngc	wzr, w9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x09, 0x5a}),
				address:          0,
			},
			want: "ngc	wzr, w9",
			wantErr: false,
		},
		{
			name: "ngc	w23, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x03, 0x1f, 0x5a}),
				address:          0,
			},
			want: "ngc	w23, wzr",
			wantErr: false,
		},
		{
			name: "ngc	x29, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x03, 0x1e, 0xda}),
				address:          0,
			},
			want: "ngc	x29, x30",
			wantErr: false,
		},
		{
			name: "ngc	xzr, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x00, 0xda}),
				address:          0,
			},
			want: "ngc	xzr, x0",
			wantErr: false,
		},
		{
			name: "ngc	x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0x1f, 0xda}),
				address:          0,
			},
			want: "ngc	x0, xzr",
			wantErr: false,
		},
		{
			name: "ngcs	w3, w12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x0c, 0x7a}),
				address:          0,
			},
			want: "ngcs	w3, w12",
			wantErr: false,
		},
		{
			name: "ngcs	wzr, w9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x09, 0x7a}),
				address:          0,
			},
			want: "ngcs	wzr, w9",
			wantErr: false,
		},
		{
			name: "ngcs	w23, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x03, 0x1f, 0x7a}),
				address:          0,
			},
			want: "ngcs	w23, wzr",
			wantErr: false,
		},
		{
			name: "ngcs	x29, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x03, 0x1e, 0xfa}),
				address:          0,
			},
			want: "ngcs	x29, x30",
			wantErr: false,
		},
		{
			name: "ngcs	xzr, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x00, 0xfa}),
				address:          0,
			},
			want: "ngcs	xzr, x0",
			wantErr: false,
		},
		{
			name: "ngcs	x0, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0x1f, 0xfa}),
				address:          0,
			},
			want: "ngcs	x0, xzr",
			wantErr: false,
		},
		{
			name: "sbfx	x1, x2, #3, #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x10, 0x43, 0x93}),
				address:          0,
			},
			want: "sbfx	x1, x2, #3, #2",
			wantErr: false,
		},
		{
			name: "asr	x3, x4, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0xfc, 0x7f, 0x93}),
				address:          0,
			},
			want: "asr	x3, x4, #63",
			wantErr: false,
		},
		{
			name: "asr	wzr, wzr, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x7f, 0x1f, 0x13}),
				address:          0,
			},
			want: "asr	wzr, wzr, #31",
			wantErr: false,
		},
		{
			name: "sbfx	w12, w9, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x01, 0x00, 0x13}),
				address:          0,
			},
			want: "sbfx	w12, w9, #0, #1",
			wantErr: false,
		},
		{
			name: "ubfiz	x4, x5, #52, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x28, 0x4c, 0xd3}),
				address:          0,
			},
			want: "ubfiz	x4, x5, #52, #11",
			wantErr: false,
		},
		{
			name: "ubfx	xzr, x4, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x00, 0x40, 0xd3}),
				address:          0,
			},
			want: "ubfx	xzr, x4, #0, #1",
			wantErr: false,
		},
		{
			name: "ubfiz	x4, xzr, #1, #6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe4, 0x17, 0x7f, 0xd3}),
				address:          0,
			},
			want: "ubfiz	x4, xzr, #1, #6",
			wantErr: false,
		},
		{
			name: "lsr	x5, x6, #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0xfc, 0x4c, 0xd3}),
				address:          0,
			},
			want: "lsr	x5, x6, #12",
			wantErr: false,
		},
		{
			name: "bfi	x4, x5, #52, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x28, 0x4c, 0xb3}),
				address:          0,
			},
			want: "bfi	x4, x5, #52, #11",
			wantErr: false,
		},
		{
			name: "bfxil	xzr, x4, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x00, 0x40, 0xb3}),
				address:          0,
			},
			want: "bfxil	xzr, x4, #0, #1",
			wantErr: false,
		},
		{
			name: "bfc	x4, #1, #6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe4, 0x17, 0x7f, 0xb3}),
				address:          0,
			},
			want: "bfc	x4, #1, #6",
			wantErr: false,
		},
		{
			name: "bfxil	x5, x6, #12, #52",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0xfc, 0x4c, 0xb3}),
				address:          0,
			},
			want: "bfxil	x5, x6, #12, #52",
			wantErr: false,
		},
		{
			name: "sxtb	w1, w2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x1c, 0x00, 0x13}),
				address:          0,
			},
			want: "sxtb	w1, w2",
			wantErr: false,
		},
		{
			name: "sxtb	xzr, w3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x1c, 0x40, 0x93}),
				address:          0,
			},
			want: "sxtb	xzr, w3",
			wantErr: false,
		},
		{
			name: "sxth	w9, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x3d, 0x00, 0x13}),
				address:          0,
			},
			want: "sxth	w9, w10",
			wantErr: false,
		},
		{
			name: "sxth	x0, w1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x3c, 0x40, 0x93}),
				address:          0,
			},
			want: "sxth	x0, w1",
			wantErr: false,
		},
		{
			name: "sxtw	x3, w30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc3, 0x7f, 0x40, 0x93}),
				address:          0,
			},
			want: "sxtw	x3, w30",
			wantErr: false,
		},
		{
			name: "uxtb	w1, w2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x1c, 0x00, 0x53}),
				address:          0,
			},
			want: "uxtb	w1, w2",
			wantErr: false,
		},
		{
			name: "uxtb	wzr, w3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x1c, 0x00, 0x53}),
				address:          0,
			},
			want: "uxtb	wzr, w3",
			wantErr: false,
		},
		{
			name: "uxth	w9, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x3d, 0x00, 0x53}),
				address:          0,
			},
			want: "uxth	w9, w10",
			wantErr: false,
		},
		{
			name: "uxth	w0, w1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x3c, 0x00, 0x53}),
				address:          0,
			},
			want: "uxth	w0, w1",
			wantErr: false,
		},
		{
			name: "asr	w3, w2, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x43, 0x7c, 0x00, 0x13}),
				address:          0,
			},
			want: "asr	w3, w2, #0",
			wantErr: false,
		},
		{
			name: "asr	w9, w10, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x1f, 0x13}),
				address:          0,
			},
			want: "asr	w9, w10, #31",
			wantErr: false,
		},
		{
			name: "asr	x20, x21, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0xfe, 0x7f, 0x93}),
				address:          0,
			},
			want: "asr	x20, x21, #63",
			wantErr: false,
		},
		{
			name: "asr	w1, wzr, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x7f, 0x03, 0x13}),
				address:          0,
			},
			want: "asr	w1, wzr, #3",
			wantErr: false,
		},
		{
			name: "lsr	w3, w2, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x43, 0x7c, 0x00, 0x53}),
				address:          0,
			},
			want: "lsr	w3, w2, #0",
			wantErr: false,
		},
		{
			name: "lsr	w9, w10, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x1f, 0x53}),
				address:          0,
			},
			want: "lsr	w9, w10, #31",
			wantErr: false,
		},
		{
			name: "lsr	x20, x21, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0xfe, 0x7f, 0xd3}),
				address:          0,
			},
			want: "lsr	x20, x21, #63",
			wantErr: false,
		},
		{
			name: "lsr	wzr, wzr, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x7f, 0x03, 0x53}),
				address:          0,
			},
			want: "lsr	wzr, wzr, #3",
			wantErr: false,
		},
		{
			name: "lsr	w3, w2, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x43, 0x7c, 0x00, 0x53}),
				address:          0,
			},
			want: "lsr	w3, w2, #0",
			wantErr: false,
		},
		{
			name: "lsl	w9, w10, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x01, 0x53}),
				address:          0,
			},
			want: "lsl	w9, w10, #31",
			wantErr: false,
		},
		{
			name: "lsl	x20, x21, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x02, 0x41, 0xd3}),
				address:          0,
			},
			want: "lsl	x20, x21, #63",
			wantErr: false,
		},
		{
			name: "lsl	w1, wzr, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x73, 0x1d, 0x53}),
				address:          0,
			},
			want: "lsl	w1, wzr, #3",
			wantErr: false,
		},
		{
			name: "sbfx	w9, w10, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x00, 0x13}),
				address:          0,
			},
			want: "sbfx	w9, w10, #0, #1",
			wantErr: false,
		},
		{
			name: "sbfiz	x2, x3, #63, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x41, 0x93}),
				address:          0,
			},
			want: "sbfiz	x2, x3, #63, #1",
			wantErr: false,
		},
		{
			name: "asr	x19, x20, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x40, 0x93}),
				address:          0,
			},
			want: "asr	x19, x20, #0",
			wantErr: false,
		},
		{
			name: "sbfiz	x9, x10, #5, #59",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe9, 0x7b, 0x93}),
				address:          0,
			},
			want: "sbfiz	x9, x10, #5, #59",
			wantErr: false,
		},
		{
			name: "asr	w9, w10, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x00, 0x13}),
				address:          0,
			},
			want: "asr	w9, w10, #0",
			wantErr: false,
		},
		{
			name: "sbfiz	w11, w12, #31, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x01, 0x01, 0x13}),
				address:          0,
			},
			want: "sbfiz	w11, w12, #31, #1",
			wantErr: false,
		},
		{
			name: "sbfiz	w13, w14, #29, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x09, 0x03, 0x13}),
				address:          0,
			},
			want: "sbfiz	w13, w14, #29, #3",
			wantErr: false,
		},
		{
			name: "sbfiz	xzr, xzr, #10, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x2b, 0x76, 0x93}),
				address:          0,
			},
			want: "sbfiz	xzr, xzr, #10, #11",
			wantErr: false,
		},
		{
			name: "sbfx	w9, w10, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x00, 0x13}),
				address:          0,
			},
			want: "sbfx	w9, w10, #0, #1",
			wantErr: false,
		},
		{
			name: "asr	x2, x3, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x7f, 0x93}),
				address:          0,
			},
			want: "asr	x2, x3, #63",
			wantErr: false,
		},
		{
			name: "asr	x19, x20, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x40, 0x93}),
				address:          0,
			},
			want: "asr	x19, x20, #0",
			wantErr: false,
		},
		{
			name: "asr	x9, x10, #5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfd, 0x45, 0x93}),
				address:          0,
			},
			want: "asr	x9, x10, #5",
			wantErr: false,
		},
		{
			name: "asr	w9, w10, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x00, 0x13}),
				address:          0,
			},
			want: "asr	w9, w10, #0",
			wantErr: false,
		},
		{
			name: "asr	w11, w12, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x7d, 0x1f, 0x13}),
				address:          0,
			},
			want: "asr	w11, w12, #31",
			wantErr: false,
		},
		{
			name: "asr	w13, w14, #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x7d, 0x1d, 0x13}),
				address:          0,
			},
			want: "asr	w13, w14, #29",
			wantErr: false,
		},
		{
			name: "sbfx	xzr, xzr, #10, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x4a, 0x93}),
				address:          0,
			},
			want: "sbfx	xzr, xzr, #10, #11",
			wantErr: false,
		},
		{
			name: "bfxil	w9, w10, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x00, 0x33}),
				address:          0,
			},
			want: "bfxil	w9, w10, #0, #1",
			wantErr: false,
		},
		{
			name: "bfi	x2, x3, #63, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x41, 0xb3}),
				address:          0,
			},
			want: "bfi	x2, x3, #63, #1",
			wantErr: false,
		},
		{
			name: "bfxil	x19, x20, #0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x40, 0xb3}),
				address:          0,
			},
			want: "bfxil	x19, x20, #0, #64",
			wantErr: false,
		},
		{
			name: "bfi	x9, x10, #5, #59",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe9, 0x7b, 0xb3}),
				address:          0,
			},
			want: "bfi	x9, x10, #5, #59",
			wantErr: false,
		},
		{
			name: "bfxil	w9, w10, #0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x00, 0x33}),
				address:          0,
			},
			want: "bfxil	w9, w10, #0, #32",
			wantErr: false,
		},
		{
			name: "bfi	w11, w12, #31, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x01, 0x01, 0x33}),
				address:          0,
			},
			want: "bfi	w11, w12, #31, #1",
			wantErr: false,
		},
		{
			name: "bfi	w13, w14, #29, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x09, 0x03, 0x33}),
				address:          0,
			},
			want: "bfi	w13, w14, #29, #3",
			wantErr: false,
		},
		{
			name: "bfc	xzr, #10, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x2b, 0x76, 0xb3}),
				address:          0,
			},
			want: "bfc	xzr, #10, #11",
			wantErr: false,
		},
		{
			name: "bfxil	w9, w10, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x00, 0x33}),
				address:          0,
			},
			want: "bfxil	w9, w10, #0, #1",
			wantErr: false,
		},
		{
			name: "bfxil	x2, x3, #63, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x7f, 0xb3}),
				address:          0,
			},
			want: "bfxil	x2, x3, #63, #1",
			wantErr: false,
		},
		{
			name: "bfxil	x19, x20, #0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x40, 0xb3}),
				address:          0,
			},
			want: "bfxil	x19, x20, #0, #64",
			wantErr: false,
		},
		{
			name: "bfxil	x9, x10, #5, #59",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfd, 0x45, 0xb3}),
				address:          0,
			},
			want: "bfxil	x9, x10, #5, #59",
			wantErr: false,
		},
		{
			name: "bfxil	w9, w10, #0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x00, 0x33}),
				address:          0,
			},
			want: "bfxil	w9, w10, #0, #32",
			wantErr: false,
		},
		{
			name: "bfxil	w11, w12, #31, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x7d, 0x1f, 0x33}),
				address:          0,
			},
			want: "bfxil	w11, w12, #31, #1",
			wantErr: false,
		},
		{
			name: "bfxil	w13, w14, #29, #3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x7d, 0x1d, 0x33}),
				address:          0,
			},
			want: "bfxil	w13, w14, #29, #3",
			wantErr: false,
		},
		{
			name: "bfxil	xzr, xzr, #10, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x4a, 0xb3}),
				address:          0,
			},
			want: "bfxil	xzr, xzr, #10, #11",
			wantErr: false,
		},
		{
			name: "ubfx	w9, w10, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x00, 0x53}),
				address:          0,
			},
			want: "ubfx	w9, w10, #0, #1",
			wantErr: false,
		},
		{
			name: "lsl	x2, x3, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x41, 0xd3}),
				address:          0,
			},
			want: "lsl	x2, x3, #63",
			wantErr: false,
		},
		{
			name: "lsr	x19, x20, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x40, 0xd3}),
				address:          0,
			},
			want: "lsr	x19, x20, #0",
			wantErr: false,
		},
		{
			name: "lsl	x9, x10, #5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe9, 0x7b, 0xd3}),
				address:          0,
			},
			want: "lsl	x9, x10, #5",
			wantErr: false,
		},
		{
			name: "lsr	w9, w10, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x00, 0x53}),
				address:          0,
			},
			want: "lsr	w9, w10, #0",
			wantErr: false,
		},
		{
			name: "lsl	w11, w12, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x01, 0x01, 0x53}),
				address:          0,
			},
			want: "lsl	w11, w12, #31",
			wantErr: false,
		},
		{
			name: "lsl	w13, w14, #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x09, 0x03, 0x53}),
				address:          0,
			},
			want: "lsl	w13, w14, #29",
			wantErr: false,
		},
		{
			name: "ubfiz	xzr, xzr, #10, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x2b, 0x76, 0xd3}),
				address:          0,
			},
			want: "ubfiz	xzr, xzr, #10, #11",
			wantErr: false,
		},
		{
			name: "ubfx	w9, w10, #0, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x00, 0x53}),
				address:          0,
			},
			want: "ubfx	w9, w10, #0, #1",
			wantErr: false,
		},
		{
			name: "lsr	x2, x3, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xfc, 0x7f, 0xd3}),
				address:          0,
			},
			want: "lsr	x2, x3, #63",
			wantErr: false,
		},
		{
			name: "lsr	x19, x20, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x40, 0xd3}),
				address:          0,
			},
			want: "lsr	x19, x20, #0",
			wantErr: false,
		},
		{
			name: "lsr	x9, x10, #5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfd, 0x45, 0xd3}),
				address:          0,
			},
			want: "lsr	x9, x10, #5",
			wantErr: false,
		},
		{
			name: "lsr	w9, w10, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x7d, 0x00, 0x53}),
				address:          0,
			},
			want: "lsr	w9, w10, #0",
			wantErr: false,
		},
		{
			name: "lsr	w11, w12, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x7d, 0x1f, 0x53}),
				address:          0,
			},
			want: "lsr	w11, w12, #31",
			wantErr: false,
		},
		{
			name: "lsr	w13, w14, #29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x7d, 0x1d, 0x53}),
				address:          0,
			},
			want: "lsr	w13, w14, #29",
			wantErr: false,
		},
		{
			name: "ubfx	xzr, xzr, #10, #11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x53, 0x4a, 0xd3}),
				address:          0,
			},
			want: "ubfx	xzr, xzr, #10, #11",
			wantErr: false,
		},
		{
			name: "bfc	w3, #0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x7f, 0x00, 0x33}),
				address:          0,
			},
			want: "bfxil	w3, wzr, #0, #32",
			wantErr: false,
		},
		{
			name: "bfc	wzr, #31, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x01, 0x33}),
				address:          0,
			},
			want: "bfc	wzr, #31, #1",
			wantErr: false,
		},
		{
			name: "bfc	x0, #5, #9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x23, 0x7b, 0xb3}),
				address:          0,
			},
			want: "bfc	x0, #5, #9",
			wantErr: false,
		},
		{
			name: "bfc	xzr, #63, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x41, 0xb3}),
				address:          0,
			},
			want: "bfc	xzr, #63, #1",
			wantErr: false,
		},
		{
			name: "cbz	w5, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x05, 0x00, 0x00, 0x34}),
				address:          0,
			},
			want: "cbz	w5, #0",
			wantErr: false,
		},
		{
			name: "cbnz	x3, #-4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0xff, 0xb5}),
				address:          0,
			},
			want: "cbnz	x3, #-4",
			wantErr: false,
		},
		{
			name: "cbz	w20, #1048572",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xff, 0x7f, 0x34}),
				address:          0,
			},
			want: "cbz	w20, #1048572",
			wantErr: false,
		},
		{
			name: "cbnz	xzr, #-1048576",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x80, 0xb5}),
				address:          0,
			},
			want: "cbnz	xzr, #-1048576",
			wantErr: false,
		},
		{
			name: "b.eq	#0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x00, 0x54}),
				address:          0,
			},
			want: "b.eq	#0",
			wantErr: false,
		},
		{
			name: "b.lt	#-4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xeb, 0xff, 0xff, 0x54}),
				address:          0,
			},
			want: "b.lt	#-4",
			wantErr: false,
		},
		{
			name: "b.lo	#1048572",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0xff, 0x7f, 0x54}),
				address:          0,
			},
			want: "b.lo	#1048572",
			wantErr: false,
		},
		{
			name: "ccmp	w1, #31, #0, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0x5f, 0x7a}),
				address:          0,
			},
			want: "ccmp	w1, #31, #0, eq",
			wantErr: false,
		},
		{
			name: "ccmp	w3, #0, #15, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0x28, 0x40, 0x7a}),
				address:          0,
			},
			want: "ccmp	w3, #0, #15, hs",
			wantErr: false,
		},
		{
			name: "ccmp	wzr, #15, #13, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x2b, 0x4f, 0x7a}),
				address:          0,
			},
			want: "ccmp	wzr, #15, #13, hs",
			wantErr: false,
		},
		{
			name: "ccmp	x9, #31, #0, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd9, 0x5f, 0xfa}),
				address:          0,
			},
			want: "ccmp	x9, #31, #0, le",
			wantErr: false,
		},
		{
			name: "ccmp	x3, #0, #15, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0xc8, 0x40, 0xfa}),
				address:          0,
			},
			want: "ccmp	x3, #0, #15, gt",
			wantErr: false,
		},
		{
			name: "ccmp	xzr, #5, #7, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x1b, 0x45, 0xfa}),
				address:          0,
			},
			want: "ccmp	xzr, #5, #7, ne",
			wantErr: false,
		},
		{
			name: "ccmn	w1, #31, #0, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x08, 0x5f, 0x3a}),
				address:          0,
			},
			want: "ccmn	w1, #31, #0, eq",
			wantErr: false,
		},
		{
			name: "ccmn	w3, #0, #15, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0x28, 0x40, 0x3a}),
				address:          0,
			},
			want: "ccmn	w3, #0, #15, hs",
			wantErr: false,
		},
		{
			name: "ccmn	wzr, #15, #13, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x2b, 0x4f, 0x3a}),
				address:          0,
			},
			want: "ccmn	wzr, #15, #13, hs",
			wantErr: false,
		},
		{
			name: "ccmn	x9, #31, #0, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd9, 0x5f, 0xba}),
				address:          0,
			},
			want: "ccmn	x9, #31, #0, le",
			wantErr: false,
		},
		{
			name: "ccmn	x3, #0, #15, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0xc8, 0x40, 0xba}),
				address:          0,
			},
			want: "ccmn	x3, #0, #15, gt",
			wantErr: false,
		},
		{
			name: "ccmn	xzr, #5, #7, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x1b, 0x45, 0xba}),
				address:          0,
			},
			want: "ccmn	xzr, #5, #7, ne",
			wantErr: false,
		},
		{
			name: "ccmp	w1, wzr, #0, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0x5f, 0x7a}),
				address:          0,
			},
			want: "ccmp	w1, wzr, #0, eq",
			wantErr: false,
		},
		{
			name: "ccmp	w3, w0, #15, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0x20, 0x40, 0x7a}),
				address:          0,
			},
			want: "ccmp	w3, w0, #15, hs",
			wantErr: false,
		},
		{
			name: "ccmp	wzr, w15, #13, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x23, 0x4f, 0x7a}),
				address:          0,
			},
			want: "ccmp	wzr, w15, #13, hs",
			wantErr: false,
		},
		{
			name: "ccmp	x9, xzr, #0, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd1, 0x5f, 0xfa}),
				address:          0,
			},
			want: "ccmp	x9, xzr, #0, le",
			wantErr: false,
		},
		{
			name: "ccmp	x3, x0, #15, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0xc0, 0x40, 0xfa}),
				address:          0,
			},
			want: "ccmp	x3, x0, #15, gt",
			wantErr: false,
		},
		{
			name: "ccmp	xzr, x5, #7, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x13, 0x45, 0xfa}),
				address:          0,
			},
			want: "ccmp	xzr, x5, #7, ne",
			wantErr: false,
		},
		{
			name: "ccmn	w1, wzr, #0, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x00, 0x5f, 0x3a}),
				address:          0,
			},
			want: "ccmn	w1, wzr, #0, eq",
			wantErr: false,
		},
		{
			name: "ccmn	w3, w0, #15, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0x20, 0x40, 0x3a}),
				address:          0,
			},
			want: "ccmn	w3, w0, #15, hs",
			wantErr: false,
		},
		{
			name: "ccmn	wzr, w15, #13, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x23, 0x4f, 0x3a}),
				address:          0,
			},
			want: "ccmn	wzr, w15, #13, hs",
			wantErr: false,
		},
		{
			name: "ccmn	x9, xzr, #0, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd1, 0x5f, 0xba}),
				address:          0,
			},
			want: "ccmn	x9, xzr, #0, le",
			wantErr: false,
		},
		{
			name: "ccmn	x3, x0, #15, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0xc0, 0x40, 0xba}),
				address:          0,
			},
			want: "ccmn	x3, x0, #15, gt",
			wantErr: false,
		},
		{
			name: "ccmn	xzr, x5, #7, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x13, 0x45, 0xba}),
				address:          0,
			},
			want: "ccmn	xzr, x5, #7, ne",
			wantErr: false,
		},
		{
			name: "csel	w1, w0, w19, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x10, 0x93, 0x1a}),
				address:          0,
			},
			want: "csel	w1, w0, w19, ne",
			wantErr: false,
		},
		{
			name: "csel	wzr, w5, w9, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x89, 0x1a}),
				address:          0,
			},
			want: "csel	wzr, w5, w9, eq",
			wantErr: false,
		},
		{
			name: "csel	w9, wzr, w30, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xc3, 0x9e, 0x1a}),
				address:          0,
			},
			want: "csel	w9, wzr, w30, gt",
			wantErr: false,
		},
		{
			name: "csel	w1, w28, wzr, mi",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0x43, 0x9f, 0x1a}),
				address:          0,
			},
			want: "csel	w1, w28, wzr, mi",
			wantErr: false,
		},
		{
			name: "csel	x19, x23, x29, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xb2, 0x9d, 0x9a}),
				address:          0,
			},
			want: "csel	x19, x23, x29, lt",
			wantErr: false,
		},
		{
			name: "csel	xzr, x3, x4, ge",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xa0, 0x84, 0x9a}),
				address:          0,
			},
			want: "csel	xzr, x3, x4, ge",
			wantErr: false,
		},
		{
			name: "csel	x5, xzr, x6, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x23, 0x86, 0x9a}),
				address:          0,
			},
			want: "csel	x5, xzr, x6, hs",
			wantErr: false,
		},
		{
			name: "csel	x7, x8, xzr, lo",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0x31, 0x9f, 0x9a}),
				address:          0,
			},
			want: "csel	x7, x8, xzr, lo",
			wantErr: false,
		},
		{
			name: "csinc	w1, w0, w19, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x14, 0x93, 0x1a}),
				address:          0,
			},
			want: "csinc	w1, w0, w19, ne",
			wantErr: false,
		},
		{
			name: "csinc	wzr, w5, w9, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x04, 0x89, 0x1a}),
				address:          0,
			},
			want: "csinc	wzr, w5, w9, eq",
			wantErr: false,
		},
		{
			name: "csinc	w9, wzr, w30, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xc7, 0x9e, 0x1a}),
				address:          0,
			},
			want: "csinc	w9, wzr, w30, gt",
			wantErr: false,
		},
		{
			name: "csinc	w1, w28, wzr, mi",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0x47, 0x9f, 0x1a}),
				address:          0,
			},
			want: "csinc	w1, w28, wzr, mi",
			wantErr: false,
		},
		{
			name: "csinc	x19, x23, x29, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xb6, 0x9d, 0x9a}),
				address:          0,
			},
			want: "csinc	x19, x23, x29, lt",
			wantErr: false,
		},
		{
			name: "csinc	xzr, x3, x4, ge",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xa4, 0x84, 0x9a}),
				address:          0,
			},
			want: "csinc	xzr, x3, x4, ge",
			wantErr: false,
		},
		{
			name: "csinc	x5, xzr, x6, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x27, 0x86, 0x9a}),
				address:          0,
			},
			want: "csinc	x5, xzr, x6, hs",
			wantErr: false,
		},
		{
			name: "csinc	x7, x8, xzr, lo",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0x35, 0x9f, 0x9a}),
				address:          0,
			},
			want: "csinc	x7, x8, xzr, lo",
			wantErr: false,
		},
		{
			name: "csinv	w1, w0, w19, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x10, 0x93, 0x5a}),
				address:          0,
			},
			want: "csinv	w1, w0, w19, ne",
			wantErr: false,
		},
		{
			name: "csinv	wzr, w5, w9, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x00, 0x89, 0x5a}),
				address:          0,
			},
			want: "csinv	wzr, w5, w9, eq",
			wantErr: false,
		},
		{
			name: "csinv	w9, wzr, w30, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xc3, 0x9e, 0x5a}),
				address:          0,
			},
			want: "csinv	w9, wzr, w30, gt",
			wantErr: false,
		},
		{
			name: "csinv	w1, w28, wzr, mi",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0x43, 0x9f, 0x5a}),
				address:          0,
			},
			want: "csinv	w1, w28, wzr, mi",
			wantErr: false,
		},
		{
			name: "csinv	x19, x23, x29, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xb2, 0x9d, 0xda}),
				address:          0,
			},
			want: "csinv	x19, x23, x29, lt",
			wantErr: false,
		},
		{
			name: "csinv	xzr, x3, x4, ge",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xa0, 0x84, 0xda}),
				address:          0,
			},
			want: "csinv	xzr, x3, x4, ge",
			wantErr: false,
		},
		{
			name: "csinv	x5, xzr, x6, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x23, 0x86, 0xda}),
				address:          0,
			},
			want: "csinv	x5, xzr, x6, hs",
			wantErr: false,
		},
		{
			name: "csinv	x7, x8, xzr, lo",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0x31, 0x9f, 0xda}),
				address:          0,
			},
			want: "csinv	x7, x8, xzr, lo",
			wantErr: false,
		},
		{
			name: "csneg	w1, w0, w19, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x14, 0x93, 0x5a}),
				address:          0,
			},
			want: "csneg	w1, w0, w19, ne",
			wantErr: false,
		},
		{
			name: "csneg	wzr, w5, w9, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x04, 0x89, 0x5a}),
				address:          0,
			},
			want: "csneg	wzr, w5, w9, eq",
			wantErr: false,
		},
		{
			name: "csneg	w9, wzr, w30, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xc7, 0x9e, 0x5a}),
				address:          0,
			},
			want: "csneg	w9, wzr, w30, gt",
			wantErr: false,
		},
		{
			name: "csneg	w1, w28, wzr, mi",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0x47, 0x9f, 0x5a}),
				address:          0,
			},
			want: "csneg	w1, w28, wzr, mi",
			wantErr: false,
		},
		{
			name: "csneg	x19, x23, x29, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xb6, 0x9d, 0xda}),
				address:          0,
			},
			want: "csneg	x19, x23, x29, lt",
			wantErr: false,
		},
		{
			name: "csneg	xzr, x3, x4, ge",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xa4, 0x84, 0xda}),
				address:          0,
			},
			want: "csneg	xzr, x3, x4, ge",
			wantErr: false,
		},
		{
			name: "csneg	x5, xzr, x6, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x27, 0x86, 0xda}),
				address:          0,
			},
			want: "csneg	x5, xzr, x6, hs",
			wantErr: false,
		},
		{
			name: "csneg	x7, x8, xzr, lo",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0x35, 0x9f, 0xda}),
				address:          0,
			},
			want: "csneg	x7, x8, xzr, lo",
			wantErr: false,
		},
		{
			name: "cset	w3, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x9f, 0x1a}),
				address:          0,
			},
			want: "cset	w3, eq",
			wantErr: false,
		},
		{
			name: "cset	x9, pl",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x47, 0x9f, 0x9a}),
				address:          0,
			},
			want: "cset	x9, pl",
			wantErr: false,
		},
		{
			name: "csetm	w20, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x9f, 0x5a}),
				address:          0,
			},
			want: "csetm	w20, ne",
			wantErr: false,
		},
		{
			name: "csetm	x30, ge",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0xb3, 0x9f, 0xda}),
				address:          0,
			},
			want: "csetm	x30, ge",
			wantErr: false,
		},
		{
			name: "cinc	w3, w5, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xd4, 0x85, 0x1a}),
				address:          0,
			},
			want: "cinc	w3, w5, gt",
			wantErr: false,
		},
		{
			name: "cinc	wzr, w4, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xc4, 0x84, 0x1a}),
				address:          0,
			},
			want: "cinc	wzr, w4, le",
			wantErr: false,
		},
		{
			name: "cset	w9, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xa7, 0x9f, 0x1a}),
				address:          0,
			},
			want: "cset	w9, lt",
			wantErr: false,
		},
		{
			name: "cinc	x3, x5, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xd4, 0x85, 0x9a}),
				address:          0,
			},
			want: "cinc	x3, x5, gt",
			wantErr: false,
		},
		{
			name: "cinc	xzr, x4, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xc4, 0x84, 0x9a}),
				address:          0,
			},
			want: "cinc	xzr, x4, le",
			wantErr: false,
		},
		{
			name: "cset	x9, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xa7, 0x9f, 0x9a}),
				address:          0,
			},
			want: "cset	x9, lt",
			wantErr: false,
		},
		{
			name: "cinv	w3, w5, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xd0, 0x85, 0x5a}),
				address:          0,
			},
			want: "cinv	w3, w5, gt",
			wantErr: false,
		},
		{
			name: "cinv	wzr, w4, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xc0, 0x84, 0x5a}),
				address:          0,
			},
			want: "cinv	wzr, w4, le",
			wantErr: false,
		},
		{
			name: "csetm	w9, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xa3, 0x9f, 0x5a}),
				address:          0,
			},
			want: "csetm	w9, lt",
			wantErr: false,
		},
		{
			name: "cinv	x3, x5, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xd0, 0x85, 0xda}),
				address:          0,
			},
			want: "cinv	x3, x5, gt",
			wantErr: false,
		},
		{
			name: "cinv	xzr, x4, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xc0, 0x84, 0xda}),
				address:          0,
			},
			want: "cinv	xzr, x4, le",
			wantErr: false,
		},
		{
			name: "csetm	x9, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xa3, 0x9f, 0xda}),
				address:          0,
			},
			want: "csetm	x9, lt",
			wantErr: false,
		},
		{
			name: "cneg	w3, w5, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xd4, 0x85, 0x5a}),
				address:          0,
			},
			want: "cneg	w3, w5, gt",
			wantErr: false,
		},
		{
			name: "cneg	wzr, w4, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xc4, 0x84, 0x5a}),
				address:          0,
			},
			want: "cneg	wzr, w4, le",
			wantErr: false,
		},
		{
			name: "cneg	w9, wzr, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xa7, 0x9f, 0x5a}),
				address:          0,
			},
			want: "cneg	w9, wzr, lt",
			wantErr: false,
		},
		{
			name: "cneg	x3, x5, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xd4, 0x85, 0xda}),
				address:          0,
			},
			want: "cneg	x3, x5, gt",
			wantErr: false,
		},
		{
			name: "cneg	xzr, x4, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xc4, 0x84, 0xda}),
				address:          0,
			},
			want: "cneg	xzr, x4, le",
			wantErr: false,
		},
		{
			name: "cneg	x9, xzr, lt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xa7, 0x9f, 0xda}),
				address:          0,
			},
			want: "cneg	x9, xzr, lt",
			wantErr: false,
		},
		{
			name: "rbit	w0, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x00, 0xc0, 0x5a}),
				address:          0,
			},
			want: "rbit	w0, w7",
			wantErr: false,
		},
		{
			name: "rbit	x18, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x72, 0x00, 0xc0, 0xda}),
				address:          0,
			},
			want: "rbit	x18, x3",
			wantErr: false,
		},
		{
			name: "rev16	w17, w1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x31, 0x04, 0xc0, 0x5a}),
				address:          0,
			},
			want: "rev16	w17, w1",
			wantErr: false,
		},
		{
			name: "rev16	x5, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x45, 0x04, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev16	x5, x2",
			wantErr: false,
		},
		{
			name: "rev	w18, w0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x12, 0x08, 0xc0, 0x5a}),
				address:          0,
			},
			want: "rev	w18, w0",
			wantErr: false,
		},
		{
			name: "rev32	x20, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x34, 0x08, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev32	x20, x1",
			wantErr: false,
		},
		{
			name: "rev32	x20, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x0b, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev32	x20, xzr",
			wantErr: false,
		},
		{
			name: "rev	x22, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x56, 0x0c, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev	x22, x2",
			wantErr: false,
		},
		{
			name: "rev	x18, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf2, 0x0f, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev	x18, xzr",
			wantErr: false,
		},
		{
			name: "rev	w7, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x0b, 0xc0, 0x5a}),
				address:          0,
			},
			want: "rev	w7, wzr",
			wantErr: false,
		},
		{
			name: "clz	w24, w3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x78, 0x10, 0xc0, 0x5a}),
				address:          0,
			},
			want: "clz	w24, w3",
			wantErr: false,
		},
		{
			name: "clz	x26, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9a, 0x10, 0xc0, 0xda}),
				address:          0,
			},
			want: "clz	x26, x4",
			wantErr: false,
		},
		{
			name: "cls	w3, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x14, 0xc0, 0x5a}),
				address:          0,
			},
			want: "cls	w3, w5",
			wantErr: false,
		},
		{
			name: "cls	x20, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x14, 0xc0, 0xda}),
				address:          0,
			},
			want: "cls	x20, x5",
			wantErr: false,
		},
		{
			name: "clz	w24, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x13, 0xc0, 0x5a}),
				address:          0,
			},
			want: "clz	w24, wzr",
			wantErr: false,
		},
		{
			name: "rev	x22, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x0f, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev	x22, xzr",
			wantErr: false,
		},
		{
			name: "rev	x13, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8d, 0x0d, 0xc0, 0xda}),
				address:          0,
			},
			want: "rev	x13, x12",
			wantErr: false,
		},
		{
			name: "udiv	w0, w7, w10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x08, 0xca, 0x1a}),
				address:          0,
			},
			want: "udiv	w0, w7, w10",
			wantErr: false,
		},
		{
			name: "udiv	x9, x22, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0a, 0xc4, 0x9a}),
				address:          0,
			},
			want: "udiv	x9, x22, x4",
			wantErr: false,
		},
		{
			name: "sdiv	w12, w21, w0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0e, 0xc0, 0x1a}),
				address:          0,
			},
			want: "sdiv	w12, w21, w0",
			wantErr: false,
		},
		{
			name: "sdiv	x13, x2, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4d, 0x0c, 0xc1, 0x9a}),
				address:          0,
			},
			want: "sdiv	x13, x2, x1",
			wantErr: false,
		},
		{
			name: "lsl	w11, w12, w13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x21, 0xcd, 0x1a}),
				address:          0,
			},
			want: "lsl	w11, w12, w13",
			wantErr: false,
		},
		{
			name: "lsl	x14, x15, x16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x21, 0xd0, 0x9a}),
				address:          0,
			},
			want: "lsl	x14, x15, x16",
			wantErr: false,
		},
		{
			name: "lsr	w17, w18, w19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x51, 0x26, 0xd3, 0x1a}),
				address:          0,
			},
			want: "lsr	w17, w18, w19",
			wantErr: false,
		},
		{
			name: "lsr	x20, x21, x22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x26, 0xd6, 0x9a}),
				address:          0,
			},
			want: "lsr	x20, x21, x22",
			wantErr: false,
		},
		{
			name: "asr	w23, w24, w25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x17, 0x2b, 0xd9, 0x1a}),
				address:          0,
			},
			want: "asr	w23, w24, w25",
			wantErr: false,
		},
		{
			name: "asr	x26, x27, x28",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7a, 0x2b, 0xdc, 0x9a}),
				address:          0,
			},
			want: "asr	x26, x27, x28",
			wantErr: false,
		},
		{
			name: "ror	w0, w1, w2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x2c, 0xc2, 0x1a}),
				address:          0,
			},
			want: "ror	w0, w1, w2",
			wantErr: false,
		},
		{
			name: "ror	x3, x4, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x2c, 0xc5, 0x9a}),
				address:          0,
			},
			want: "ror	x3, x4, x5",
			wantErr: false,
		},
		{
			name: "lsl	w6, w7, w8",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x20, 0xc8, 0x1a}),
				address:          0,
			},
			want: "lsl	w6, w7, w8",
			wantErr: false,
		},
		{
			name: "lsl	x9, x10, x11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x21, 0xcb, 0x9a}),
				address:          0,
			},
			want: "lsl	x9, x10, x11",
			wantErr: false,
		},
		{
			name: "lsr	w12, w13, w14",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x25, 0xce, 0x1a}),
				address:          0,
			},
			want: "lsr	w12, w13, w14",
			wantErr: false,
		},
		{
			name: "lsr	x15, x16, x17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0f, 0x26, 0xd1, 0x9a}),
				address:          0,
			},
			want: "lsr	x15, x16, x17",
			wantErr: false,
		},
		{
			name: "asr	w18, w19, w20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x72, 0x2a, 0xd4, 0x1a}),
				address:          0,
			},
			want: "asr	w18, w19, w20",
			wantErr: false,
		},
		{
			name: "asr	x21, x22, x23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x2a, 0xd7, 0x9a}),
				address:          0,
			},
			want: "asr	x21, x22, x23",
			wantErr: false,
		},
		{
			name: "ror	w24, w25, w26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0x2f, 0xda, 0x1a}),
				address:          0,
			},
			want: "ror	w24, w25, w26",
			wantErr: false,
		},
		{
			name: "ror	x27, x28, x29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x2f, 0xdd, 0x9a}),
				address:          0,
			},
			want: "ror	x27, x28, x29",
			wantErr: false,
		},
		{
			name: "madd	w1, w3, w7, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x10, 0x07, 0x1b}),
				address:          0,
			},
			want: "madd	w1, w3, w7, w4",
			wantErr: false,
		},
		{
			name: "madd	wzr, w0, w9, w11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x2c, 0x09, 0x1b}),
				address:          0,
			},
			want: "madd	wzr, w0, w9, w11",
			wantErr: false,
		},
		{
			name: "madd	w13, wzr, w4, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x13, 0x04, 0x1b}),
				address:          0,
			},
			want: "madd	w13, wzr, w4, w4",
			wantErr: false,
		},
		{
			name: "madd	w19, w30, wzr, w29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd3, 0x77, 0x1f, 0x1b}),
				address:          0,
			},
			want: "madd	w19, w30, wzr, w29",
			wantErr: false,
		},
		{
			name: "mul	w4, w5, w6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x7c, 0x06, 0x1b}),
				address:          0,
			},
			want: "mul	w4, w5, w6",
			wantErr: false,
		},
		{
			name: "madd	x1, x3, x7, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x10, 0x07, 0x9b}),
				address:          0,
			},
			want: "madd	x1, x3, x7, x4",
			wantErr: false,
		},
		{
			name: "madd	xzr, x0, x9, x11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x2c, 0x09, 0x9b}),
				address:          0,
			},
			want: "madd	xzr, x0, x9, x11",
			wantErr: false,
		},
		{
			name: "madd	x13, xzr, x4, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x13, 0x04, 0x9b}),
				address:          0,
			},
			want: "madd	x13, xzr, x4, x4",
			wantErr: false,
		},
		{
			name: "madd	x19, x30, xzr, x29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd3, 0x77, 0x1f, 0x9b}),
				address:          0,
			},
			want: "madd	x19, x30, xzr, x29",
			wantErr: false,
		},
		{
			name: "mul	x4, x5, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x7c, 0x06, 0x9b}),
				address:          0,
			},
			want: "mul	x4, x5, x6",
			wantErr: false,
		},
		{
			name: "msub	w1, w3, w7, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x90, 0x07, 0x1b}),
				address:          0,
			},
			want: "msub	w1, w3, w7, w4",
			wantErr: false,
		},
		{
			name: "msub	wzr, w0, w9, w11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0xac, 0x09, 0x1b}),
				address:          0,
			},
			want: "msub	wzr, w0, w9, w11",
			wantErr: false,
		},
		{
			name: "msub	w13, wzr, w4, w4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x93, 0x04, 0x1b}),
				address:          0,
			},
			want: "msub	w13, wzr, w4, w4",
			wantErr: false,
		},
		{
			name: "msub	w19, w30, wzr, w29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd3, 0xf7, 0x1f, 0x1b}),
				address:          0,
			},
			want: "msub	w19, w30, wzr, w29",
			wantErr: false,
		},
		{
			name: "mneg	w4, w5, w6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xfc, 0x06, 0x1b}),
				address:          0,
			},
			want: "mneg	w4, w5, w6",
			wantErr: false,
		},
		{
			name: "msub	x1, x3, x7, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x90, 0x07, 0x9b}),
				address:          0,
			},
			want: "msub	x1, x3, x7, x4",
			wantErr: false,
		},
		{
			name: "msub	xzr, x0, x9, x11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0xac, 0x09, 0x9b}),
				address:          0,
			},
			want: "msub	xzr, x0, x9, x11",
			wantErr: false,
		},
		{
			name: "msub	x13, xzr, x4, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x93, 0x04, 0x9b}),
				address:          0,
			},
			want: "msub	x13, xzr, x4, x4",
			wantErr: false,
		},
		{
			name: "msub	x19, x30, xzr, x29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd3, 0xf7, 0x1f, 0x9b}),
				address:          0,
			},
			want: "msub	x19, x30, xzr, x29",
			wantErr: false,
		},
		{
			name: "mneg	x4, x5, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0xfc, 0x06, 0x9b}),
				address:          0,
			},
			want: "mneg	x4, x5, x6",
			wantErr: false,
		},
		{
			name: "smaddl	x3, w5, w2, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x24, 0x22, 0x9b}),
				address:          0,
			},
			want: "smaddl	x3, w5, w2, x9",
			wantErr: false,
		},
		{
			name: "smaddl	xzr, w10, w11, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x31, 0x2b, 0x9b}),
				address:          0,
			},
			want: "smaddl	xzr, w10, w11, x12",
			wantErr: false,
		},
		{
			name: "smaddl	x13, wzr, w14, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x3f, 0x2e, 0x9b}),
				address:          0,
			},
			want: "smaddl	x13, wzr, w14, x15",
			wantErr: false,
		},
		{
			name: "smaddl	x16, w17, wzr, x18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x4a, 0x3f, 0x9b}),
				address:          0,
			},
			want: "smaddl	x16, w17, wzr, x18",
			wantErr: false,
		},
		{
			name: "smull	x19, w20, w21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x7e, 0x35, 0x9b}),
				address:          0,
			},
			want: "smull	x19, w20, w21",
			wantErr: false,
		},
		{
			name: "smsubl	x3, w5, w2, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xa4, 0x22, 0x9b}),
				address:          0,
			},
			want: "smsubl	x3, w5, w2, x9",
			wantErr: false,
		},
		{
			name: "smsubl	xzr, w10, w11, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xb1, 0x2b, 0x9b}),
				address:          0,
			},
			want: "smsubl	xzr, w10, w11, x12",
			wantErr: false,
		},
		{
			name: "smsubl	x13, wzr, w14, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0xbf, 0x2e, 0x9b}),
				address:          0,
			},
			want: "smsubl	x13, wzr, w14, x15",
			wantErr: false,
		},
		{
			name: "smsubl	x16, w17, wzr, x18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0xca, 0x3f, 0x9b}),
				address:          0,
			},
			want: "smsubl	x16, w17, wzr, x18",
			wantErr: false,
		},
		{
			name: "smnegl	x19, w20, w21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x35, 0x9b}),
				address:          0,
			},
			want: "smnegl	x19, w20, w21",
			wantErr: false,
		},
		{
			name: "umaddl	x3, w5, w2, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x24, 0xa2, 0x9b}),
				address:          0,
			},
			want: "umaddl	x3, w5, w2, x9",
			wantErr: false,
		},
		{
			name: "umaddl	xzr, w10, w11, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x31, 0xab, 0x9b}),
				address:          0,
			},
			want: "umaddl	xzr, w10, w11, x12",
			wantErr: false,
		},
		{
			name: "umaddl	x13, wzr, w14, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x3f, 0xae, 0x9b}),
				address:          0,
			},
			want: "umaddl	x13, wzr, w14, x15",
			wantErr: false,
		},
		{
			name: "umaddl	x16, w17, wzr, x18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x4a, 0xbf, 0x9b}),
				address:          0,
			},
			want: "umaddl	x16, w17, wzr, x18",
			wantErr: false,
		},
		{
			name: "umull	x19, w20, w21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x7e, 0xb5, 0x9b}),
				address:          0,
			},
			want: "umull	x19, w20, w21",
			wantErr: false,
		},
		{
			name: "umsubl	x3, w5, w2, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xa4, 0xa2, 0x9b}),
				address:          0,
			},
			want: "umsubl	x3, w5, w2, x9",
			wantErr: false,
		},
		{
			name: "umsubl	xzr, w10, w11, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xb1, 0xab, 0x9b}),
				address:          0,
			},
			want: "umsubl	xzr, w10, w11, x12",
			wantErr: false,
		},
		{
			name: "umsubl	x13, wzr, w14, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0xbf, 0xae, 0x9b}),
				address:          0,
			},
			want: "umsubl	x13, wzr, w14, x15",
			wantErr: false,
		},
		{
			name: "umsubl	x16, w17, wzr, x18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0xca, 0xbf, 0x9b}),
				address:          0,
			},
			want: "umsubl	x16, w17, wzr, x18",
			wantErr: false,
		},
		{
			name: "umnegl	x19, w20, w21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0xb5, 0x9b}),
				address:          0,
			},
			want: "umnegl	x19, w20, w21",
			wantErr: false,
		},
		{
			name: "smulh	x30, x29, x28",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x7f, 0x5c, 0x9b}),
				address:          0,
			},
			want: "smulh	x30, x29, x28",
			wantErr: false,
		},
		{
			name: "smulh	xzr, x27, x26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x7f, 0x5a, 0x9b}),
				address:          0,
			},
			want: "smulh	xzr, x27, x26",
			wantErr: false,
		},
		{
			name: "smulh	x25, xzr, x24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf9, 0x7f, 0x58, 0x9b}),
				address:          0,
			},
			want: "smulh	x25, xzr, x24",
			wantErr: false,
		},
		{
			name: "smulh	x23, x22, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd7, 0x7e, 0x5f, 0x9b}),
				address:          0,
			},
			want: "smulh	x23, x22, xzr",
			wantErr: false,
		},
		{
			name: "umulh	x30, x29, x28",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbe, 0x7f, 0xdc, 0x9b}),
				address:          0,
			},
			want: "umulh	x30, x29, x28",
			wantErr: false,
		},
		{
			name: "umulh	xzr, x27, x26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x7f, 0xda, 0x9b}),
				address:          0,
			},
			want: "umulh	xzr, x27, x26",
			wantErr: false,
		},
		{
			name: "umulh	x25, xzr, x24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf9, 0x7f, 0xd8, 0x9b}),
				address:          0,
			},
			want: "umulh	x25, xzr, x24",
			wantErr: false,
		},
		{
			name: "umulh	x23, x22, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd7, 0x7e, 0xdf, 0x9b}),
				address:          0,
			},
			want: "umulh	x23, x22, xzr",
			wantErr: false,
		},
		{
			name: "mul	w3, w4, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x7c, 0x05, 0x1b}),
				address:          0,
			},
			want: "mul	w3, w4, w5",
			wantErr: false,
		},
		{
			name: "mul	wzr, w6, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x7c, 0x07, 0x1b}),
				address:          0,
			},
			want: "mul	wzr, w6, w7",
			wantErr: false,
		},
		{
			name: "mul	w8, wzr, w9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0x7f, 0x09, 0x1b}),
				address:          0,
			},
			want: "mul	w8, wzr, w9",
			wantErr: false,
		},
		{
			name: "mul	w10, w11, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x7d, 0x1f, 0x1b}),
				address:          0,
			},
			want: "mul	w10, w11, wzr",
			wantErr: false,
		},
		{
			name: "mul	x12, x13, x14",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x7d, 0x0e, 0x9b}),
				address:          0,
			},
			want: "mul	x12, x13, x14",
			wantErr: false,
		},
		{
			name: "mul	xzr, x15, x16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x7d, 0x10, 0x9b}),
				address:          0,
			},
			want: "mul	xzr, x15, x16",
			wantErr: false,
		},
		{
			name: "mul	x17, xzr, x18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0x7f, 0x12, 0x9b}),
				address:          0,
			},
			want: "mul	x17, xzr, x18",
			wantErr: false,
		},
		{
			name: "mul	x19, x20, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x7e, 0x1f, 0x9b}),
				address:          0,
			},
			want: "mul	x19, x20, xzr",
			wantErr: false,
		},
		{
			name: "mneg	w21, w22, w23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0xfe, 0x17, 0x1b}),
				address:          0,
			},
			want: "mneg	w21, w22, w23",
			wantErr: false,
		},
		{
			name: "mneg	wzr, w24, w25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0xff, 0x19, 0x1b}),
				address:          0,
			},
			want: "mneg	wzr, w24, w25",
			wantErr: false,
		},
		{
			name: "mneg	w26, wzr, w27",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0xff, 0x1b, 0x1b}),
				address:          0,
			},
			want: "mneg	w26, wzr, w27",
			wantErr: false,
		},
		{
			name: "mneg	w28, w29, wzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbc, 0xff, 0x1f, 0x1b}),
				address:          0,
			},
			want: "mneg	w28, w29, wzr",
			wantErr: false,
		},
		{
			name: "smull	x11, w13, w17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x7d, 0x31, 0x9b}),
				address:          0,
			},
			want: "smull	x11, w13, w17",
			wantErr: false,
		},
		{
			name: "umull	x11, w13, w17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x7d, 0xb1, 0x9b}),
				address:          0,
			},
			want: "umull	x11, w13, w17",
			wantErr: false,
		},
		{
			name: "smnegl	x11, w13, w17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xfd, 0x31, 0x9b}),
				address:          0,
			},
			want: "smnegl	x11, w13, w17",
			wantErr: false,
		},
		{
			name: "umnegl	x11, w13, w17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xfd, 0xb1, 0x9b}),
				address:          0,
			},
			want: "umnegl	x11, w13, w17",
			wantErr: false,
		},
		{
			name: "svc	#0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x00, 0x00, 0xd4}),
				address:          0,
			},
			want: "svc	#0",
			wantErr: false,
		},
		{
			name: "svc	#0xffff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0xff, 0x1f, 0xd4}),
				address:          0,
			},
			want: "svc	#0xffff",
			wantErr: false,
		},
		{
			name: "hvc	#0x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x22, 0x00, 0x00, 0xd4}),
				address:          0,
			},
			want: "hvc	#0x1",
			wantErr: false,
		},
		{
			name: "smc	#0x2ee0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0xdc, 0x05, 0xd4}),
				address:          0,
			},
			want: "smc	#0x2ee0",
			wantErr: false,
		},
		{
			name: "brk	#0xc",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x01, 0x20, 0xd4}),
				address:          0,
			},
			want: "brk	#0xc",
			wantErr: false,
		},
		{
			name: "hlt	#0x7b",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x0f, 0x40, 0xd4}),
				address:          0,
			},
			want: "hlt	#0x7b",
			wantErr: false,
		},
		{
			name: "dcps1	#0x2a",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x05, 0xa0, 0xd4}),
				address:          0,
			},
			want: "dcps1	#0x2a",
			wantErr: false,
		},
		{
			name: "dcps2	#0x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x22, 0x01, 0xa0, 0xd4}),
				address:          0,
			},
			want: "dcps2	#0x9",
			wantErr: false,
		},
		{
			name: "dcps3	#0x3e8",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x7d, 0xa0, 0xd4}),
				address:          0,
			},
			want: "dcps3	#0x3e8",
			wantErr: false,
		},
		{
			name: "dcps1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x00, 0xa0, 0xd4}),
				address:          0,
			},
			want:    "dcps1",
			wantErr: false,
		},
		{
			name: "dcps2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x02, 0x00, 0xa0, 0xd4}),
				address:          0,
			},
			want:    "dcps2",
			wantErr: false,
		},
		{
			name: "dcps3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x00, 0xa0, 0xd4}),
				address:          0,
			},
			want:    "dcps3",
			wantErr: false,
		},
		{
			name: "extr	w3, w5, w7, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x87, 0x13}),
				address:          0,
			},
			want: "extr	w3, w5, w7, #0",
			wantErr: false,
		},
		{
			name: "extr	w11, w13, w17, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x7d, 0x91, 0x13}),
				address:          0,
			},
			want: "extr	w11, w13, w17, #31",
			wantErr: false,
		},
		{
			name: "extr	x3, x5, x7, #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x3c, 0xc7, 0x93}),
				address:          0,
			},
			want: "extr	x3, x5, x7, #15",
			wantErr: false,
		},
		{
			name: "extr	x11, x13, x17, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xfd, 0xd1, 0x93}),
				address:          0,
			},
			want: "extr	x11, x13, x17, #63",
			wantErr: false,
		},
		{
			name: "ror	x19, x23, #24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0x62, 0xd7, 0x93}),
				address:          0,
			},
			want: "ror	x19, x23, #24",
			wantErr: false,
		},
		{
			name: "ror	x29, xzr, #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0xff, 0xdf, 0x93}),
				address:          0,
			},
			want: "ror	x29, xzr, #63",
			wantErr: false,
		},
		{
			name: "ror	w9, w13, #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x7d, 0x8d, 0x13}),
				address:          0,
			},
			want: "ror	w9, w13, #31",
			wantErr: false,
		},
		{
			name: "fcmp	s3, s5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x60, 0x20, 0x25, 0x1e}),
				address:          0,
			},
			want: "fcmp	s3, s5",
			wantErr: false,
		},
		{
			name: "fcmp	s31, #0.0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0x23, 0x20, 0x1e}),
				address:          0,
			},
			want: "fcmp	s31, #0.0",
			wantErr: false,
		},
		{
			name: "fcmpe	s29, s30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb0, 0x23, 0x3e, 0x1e}),
				address:          0,
			},
			want: "fcmpe	s29, s30",
			wantErr: false,
		},
		{
			name: "fcmpe	s15, #0.0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf8, 0x21, 0x20, 0x1e}),
				address:          0,
			},
			want: "fcmpe	s15, #0.0",
			wantErr: false,
		},
		{
			name: "fcmp	d4, d12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x20, 0x6c, 0x1e}),
				address:          0,
			},
			want: "fcmp	d4, d12",
			wantErr: false,
		},
		{
			name: "fcmp	d23, #0.0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0x22, 0x60, 0x1e}),
				address:          0,
			},
			want: "fcmp	d23, #0.0",
			wantErr: false,
		},
		{
			name: "fcmpe	d26, d22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x50, 0x23, 0x76, 0x1e}),
				address:          0,
			},
			want: "fcmpe	d26, d22",
			wantErr: false,
		},
		{
			name: "fcmpe	d29, #0.0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb8, 0x23, 0x60, 0x1e}),
				address:          0,
			},
			want: "fcmpe	d29, #0.0",
			wantErr: false,
		},
		{
			name: "fccmp	s1, s31, #0, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x04, 0x3f, 0x1e}),
				address:          0,
			},
			want: "fccmp	s1, s31, #0, eq",
			wantErr: false,
		},
		{
			name: "fccmp	s3, s0, #15, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0x24, 0x20, 0x1e}),
				address:          0,
			},
			want: "fccmp	s3, s0, #15, hs",
			wantErr: false,
		},
		{
			name: "fccmp	s31, s15, #13, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x27, 0x2f, 0x1e}),
				address:          0,
			},
			want: "fccmp	s31, s15, #13, hs",
			wantErr: false,
		},
		{
			name: "fccmp	d9, d31, #0, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0xd5, 0x7f, 0x1e}),
				address:          0,
			},
			want: "fccmp	d9, d31, #0, le",
			wantErr: false,
		},
		{
			name: "fccmp	d3, d0, #15, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6f, 0xc4, 0x60, 0x1e}),
				address:          0,
			},
			want: "fccmp	d3, d0, #15, gt",
			wantErr: false,
		},
		{
			name: "fccmp	d31, d5, #7, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x17, 0x65, 0x1e}),
				address:          0,
			},
			want: "fccmp	d31, d5, #7, ne",
			wantErr: false,
		},
		{
			name: "fccmpe	s1, s31, #0, eq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x04, 0x3f, 0x1e}),
				address:          0,
			},
			want: "fccmpe	s1, s31, #0, eq",
			wantErr: false,
		},
		{
			name: "fccmpe	s3, s0, #15, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x24, 0x20, 0x1e}),
				address:          0,
			},
			want: "fccmpe	s3, s0, #15, hs",
			wantErr: false,
		},
		{
			name: "fccmpe	s31, s15, #13, hs",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0x27, 0x2f, 0x1e}),
				address:          0,
			},
			want: "fccmpe	s31, s15, #13, hs",
			wantErr: false,
		},
		{
			name: "fccmpe	d9, d31, #0, le",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0xd5, 0x7f, 0x1e}),
				address:          0,
			},
			want: "fccmpe	d9, d31, #0, le",
			wantErr: false,
		},
		{
			name: "fccmpe	d3, d0, #15, gt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xc4, 0x60, 0x1e}),
				address:          0,
			},
			want: "fccmpe	d3, d0, #15, gt",
			wantErr: false,
		},
		{
			name: "fccmpe	d31, d5, #7, ne",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x17, 0x65, 0x1e}),
				address:          0,
			},
			want: "fccmpe	d31, d5, #7, ne",
			wantErr: false,
		},
		{
			name: "fcsel	s3, s20, s9, pl",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x5e, 0x29, 0x1e}),
				address:          0,
			},
			want: "fcsel	s3, s20, s9, pl",
			wantErr: false,
		},
		{
			name: "fcsel	d9, d10, d11, mi",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x4d, 0x6b, 0x1e}),
				address:          0,
			},
			want: "fcsel	d9, d10, d11, mi",
			wantErr: false,
		},
		{
			name: "fmov	s0, s1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x40, 0x20, 0x1e}),
				address:          0,
			},
			want: "fmov	s0, s1",
			wantErr: false,
		},
		{
			name: "fabs	s2, s3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xc0, 0x20, 0x1e}),
				address:          0,
			},
			want: "fabs	s2, s3",
			wantErr: false,
		},
		{
			name: "fneg	s4, s5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x40, 0x21, 0x1e}),
				address:          0,
			},
			want: "fneg	s4, s5",
			wantErr: false,
		},
		{
			name: "fsqrt	s6, s7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xc0, 0x21, 0x1e}),
				address:          0,
			},
			want: "fsqrt	s6, s7",
			wantErr: false,
		},
		{
			name: "fcvt	d8, s9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0xc1, 0x22, 0x1e}),
				address:          0,
			},
			want: "fcvt	d8, s9",
			wantErr: false,
		},
		{
			name: "fcvt	h10, s11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0xc1, 0x23, 0x1e}),
				address:          0,
			},
			want: "fcvt	h10, s11",
			wantErr: false,
		},
		{
			name: "frintn	s12, s13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x41, 0x24, 0x1e}),
				address:          0,
			},
			want: "frintn	s12, s13",
			wantErr: false,
		},
		{
			name: "frintp	s14, s15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0xc1, 0x24, 0x1e}),
				address:          0,
			},
			want: "frintp	s14, s15",
			wantErr: false,
		},
		{
			name: "frintm	s16, s17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x42, 0x25, 0x1e}),
				address:          0,
			},
			want: "frintm	s16, s17",
			wantErr: false,
		},
		{
			name: "frintz	s18, s19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x72, 0xc2, 0x25, 0x1e}),
				address:          0,
			},
			want: "frintz	s18, s19",
			wantErr: false,
		},
		{
			name: "frinta	s20, s21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x42, 0x26, 0x1e}),
				address:          0,
			},
			want: "frinta	s20, s21",
			wantErr: false,
		},
		{
			name: "frintx	s22, s23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x42, 0x27, 0x1e}),
				address:          0,
			},
			want: "frintx	s22, s23",
			wantErr: false,
		},
		{
			name: "frinti	s24, s25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0xc3, 0x27, 0x1e}),
				address:          0,
			},
			want: "frinti	s24, s25",
			wantErr: false,
		},
		{
			name: "fmov	d0, d1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x40, 0x60, 0x1e}),
				address:          0,
			},
			want: "fmov	d0, d1",
			wantErr: false,
		},
		{
			name: "fabs	d2, d3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0xc0, 0x60, 0x1e}),
				address:          0,
			},
			want: "fabs	d2, d3",
			wantErr: false,
		},
		{
			name: "fneg	d4, d5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x40, 0x61, 0x1e}),
				address:          0,
			},
			want: "fneg	d4, d5",
			wantErr: false,
		},
		{
			name: "fsqrt	d6, d7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0xc0, 0x61, 0x1e}),
				address:          0,
			},
			want: "fsqrt	d6, d7",
			wantErr: false,
		},
		{
			name: "fcvt	s8, d9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x41, 0x62, 0x1e}),
				address:          0,
			},
			want: "fcvt	s8, d9",
			wantErr: false,
		},
		{
			name: "fcvt	h10, d11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0xc1, 0x63, 0x1e}),
				address:          0,
			},
			want: "fcvt	h10, d11",
			wantErr: false,
		},
		{
			name: "frintn	d12, d13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x41, 0x64, 0x1e}),
				address:          0,
			},
			want: "frintn	d12, d13",
			wantErr: false,
		},
		{
			name: "frintp	d14, d15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0xc1, 0x64, 0x1e}),
				address:          0,
			},
			want: "frintp	d14, d15",
			wantErr: false,
		},
		{
			name: "frintm	d16, d17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x42, 0x65, 0x1e}),
				address:          0,
			},
			want: "frintm	d16, d17",
			wantErr: false,
		},
		{
			name: "frintz	d18, d19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x72, 0xc2, 0x65, 0x1e}),
				address:          0,
			},
			want: "frintz	d18, d19",
			wantErr: false,
		},
		{
			name: "frinta	d20, d21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x42, 0x66, 0x1e}),
				address:          0,
			},
			want: "frinta	d20, d21",
			wantErr: false,
		},
		{
			name: "frintx	d22, d23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x42, 0x67, 0x1e}),
				address:          0,
			},
			want: "frintx	d22, d23",
			wantErr: false,
		},
		{
			name: "frinti	d24, d25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x38, 0xc3, 0x67, 0x1e}),
				address:          0,
			},
			want: "frinti	d24, d25",
			wantErr: false,
		},
		{
			name: "fcvt	s26, h27",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7a, 0x43, 0xe2, 0x1e}),
				address:          0,
			},
			want: "fcvt	s26, h27",
			wantErr: false,
		},
		{
			name: "fcvt	d28, h29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbc, 0xc3, 0xe2, 0x1e}),
				address:          0,
			},
			want: "fcvt	d28, h29",
			wantErr: false,
		},
		{
			name: "fmul	s20, s19, s17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0x0a, 0x31, 0x1e}),
				address:          0,
			},
			want: "fmul	s20, s19, s17",
			wantErr: false,
		},
		{
			name: "fdiv	s1, s2, s3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x18, 0x23, 0x1e}),
				address:          0,
			},
			want: "fdiv	s1, s2, s3",
			wantErr: false,
		},
		{
			name: "fadd	s4, s5, s6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x28, 0x26, 0x1e}),
				address:          0,
			},
			want: "fadd	s4, s5, s6",
			wantErr: false,
		},
		{
			name: "fsub	s7, s8, s9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0x39, 0x29, 0x1e}),
				address:          0,
			},
			want: "fsub	s7, s8, s9",
			wantErr: false,
		},
		{
			name: "fmax	s10, s11, s12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x49, 0x2c, 0x1e}),
				address:          0,
			},
			want: "fmax	s10, s11, s12",
			wantErr: false,
		},
		{
			name: "fmin	s13, s14, s15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x59, 0x2f, 0x1e}),
				address:          0,
			},
			want: "fmin	s13, s14, s15",
			wantErr: false,
		},
		{
			name: "fmaxnm	s16, s17, s18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x6a, 0x32, 0x1e}),
				address:          0,
			},
			want: "fmaxnm	s16, s17, s18",
			wantErr: false,
		},
		{
			name: "fminnm	s19, s20, s21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x7a, 0x35, 0x1e}),
				address:          0,
			},
			want: "fminnm	s19, s20, s21",
			wantErr: false,
		},
		{
			name: "fnmul	s22, s23, s24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x8a, 0x38, 0x1e}),
				address:          0,
			},
			want: "fnmul	s22, s23, s24",
			wantErr: false,
		},
		{
			name: "fmul	d20, d19, d17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x74, 0x0a, 0x71, 0x1e}),
				address:          0,
			},
			want: "fmul	d20, d19, d17",
			wantErr: false,
		},
		{
			name: "fdiv	d1, d2, d3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0x18, 0x63, 0x1e}),
				address:          0,
			},
			want: "fdiv	d1, d2, d3",
			wantErr: false,
		},
		{
			name: "fadd	d4, d5, d6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x28, 0x66, 0x1e}),
				address:          0,
			},
			want: "fadd	d4, d5, d6",
			wantErr: false,
		},
		{
			name: "fsub	d7, d8, d9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x07, 0x39, 0x69, 0x1e}),
				address:          0,
			},
			want: "fsub	d7, d8, d9",
			wantErr: false,
		},
		{
			name: "fmax	d10, d11, d12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x49, 0x6c, 0x1e}),
				address:          0,
			},
			want: "fmax	d10, d11, d12",
			wantErr: false,
		},
		{
			name: "fmin	d13, d14, d15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x59, 0x6f, 0x1e}),
				address:          0,
			},
			want: "fmin	d13, d14, d15",
			wantErr: false,
		},
		{
			name: "fmaxnm	d16, d17, d18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x6a, 0x72, 0x1e}),
				address:          0,
			},
			want: "fmaxnm	d16, d17, d18",
			wantErr: false,
		},
		{
			name: "fminnm	d19, d20, d21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x7a, 0x75, 0x1e}),
				address:          0,
			},
			want: "fminnm	d19, d20, d21",
			wantErr: false,
		},
		{
			name: "fnmul	d22, d23, d24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x8a, 0x78, 0x1e}),
				address:          0,
			},
			want: "fnmul	d22, d23, d24",
			wantErr: false,
		},
		{
			name: "fmadd	s3, s5, s6, s31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x7c, 0x06, 0x1f}),
				address:          0,
			},
			want: "fmadd	s3, s5, s6, s31",
			wantErr: false,
		},
		{
			name: "fmadd	d3, d13, d0, d23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x5d, 0x40, 0x1f}),
				address:          0,
			},
			want: "fmadd	d3, d13, d0, d23",
			wantErr: false,
		},
		{
			name: "fmsub	s3, s5, s6, s31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x06, 0x1f}),
				address:          0,
			},
			want: "fmsub	s3, s5, s6, s31",
			wantErr: false,
		},
		{
			name: "fmsub	d3, d13, d0, d23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xdd, 0x40, 0x1f}),
				address:          0,
			},
			want: "fmsub	d3, d13, d0, d23",
			wantErr: false,
		},
		{
			name: "fnmadd	s3, s5, s6, s31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x7c, 0x26, 0x1f}),
				address:          0,
			},
			want: "fnmadd	s3, s5, s6, s31",
			wantErr: false,
		},
		{
			name: "fnmadd	d3, d13, d0, d23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x5d, 0x60, 0x1f}),
				address:          0,
			},
			want: "fnmadd	d3, d13, d0, d23",
			wantErr: false,
		},
		{
			name: "fnmsub	s3, s5, s6, s31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x26, 0x1f}),
				address:          0,
			},
			want: "fnmsub	s3, s5, s6, s31",
			wantErr: false,
		},
		{
			name: "fnmsub	d3, d13, d0, d23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xdd, 0x60, 0x1f}),
				address:          0,
			},
			want: "fnmsub	d3, d13, d0, d23",
			wantErr: false,
		},
		{
			name: "fcvtzs	w3, s5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x18, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	w3, s5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzs	wzr, s20, #13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xce, 0x18, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	wzr, s20, #13",
			wantErr: false,
		},
		{
			name: "fcvtzs	w19, s0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x80, 0x18, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	w19, s0, #32",
			wantErr: false,
		},
		{
			name: "fcvtzs	x3, s5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x18, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x3, s5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzs	x12, s30, #45",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x4f, 0x18, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x12, s30, #45",
			wantErr: false,
		},
		{
			name: "fcvtzs	x19, s0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x00, 0x18, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x19, s0, #64",
			wantErr: false,
		},
		{
			name: "fcvtzs	w3, d5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x58, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	w3, d5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzs	wzr, d20, #13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xce, 0x58, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	wzr, d20, #13",
			wantErr: false,
		},
		{
			name: "fcvtzs	w19, d0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x80, 0x58, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	w19, d0, #32",
			wantErr: false,
		},
		{
			name: "fcvtzs	x3, d5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x58, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x3, d5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzs	x12, d30, #45",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x4f, 0x58, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x12, d30, #45",
			wantErr: false,
		},
		{
			name: "fcvtzs	x19, d0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x00, 0x58, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x19, d0, #64",
			wantErr: false,
		},
		{
			name: "fcvtzu	w3, s5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x19, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	w3, s5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzu	wzr, s20, #13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xce, 0x19, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	wzr, s20, #13",
			wantErr: false,
		},
		{
			name: "fcvtzu	w19, s0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x80, 0x19, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	w19, s0, #32",
			wantErr: false,
		},
		{
			name: "fcvtzu	x3, s5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x19, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x3, s5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzu	x12, s30, #45",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x4f, 0x19, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x12, s30, #45",
			wantErr: false,
		},
		{
			name: "fcvtzu	x19, s0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x00, 0x19, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x19, s0, #64",
			wantErr: false,
		},
		{
			name: "fcvtzu	w3, d5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x59, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	w3, d5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzu	wzr, d20, #13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xce, 0x59, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	wzr, d20, #13",
			wantErr: false,
		},
		{
			name: "fcvtzu	w19, d0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x80, 0x59, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	w19, d0, #32",
			wantErr: false,
		},
		{
			name: "fcvtzu	x3, d5, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x59, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x3, d5, #1",
			wantErr: false,
		},
		{
			name: "fcvtzu	x12, d30, #45",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x4f, 0x59, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x12, d30, #45",
			wantErr: false,
		},
		{
			name: "fcvtzu	x19, d0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x00, 0x59, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x19, d0, #64",
			wantErr: false,
		},
		{
			name: "scvtf	s23, w19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x02, 0x1e}),
				address:          0,
			},
			want: "scvtf	s23, w19, #1",
			wantErr: false,
		},
		{
			name: "scvtf	s31, wzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x02, 0x1e}),
				address:          0,
			},
			want: "scvtf	s31, wzr, #20",
			wantErr: false,
		},
		{
			name: "scvtf	s14, w0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x80, 0x02, 0x1e}),
				address:          0,
			},
			want: "scvtf	s14, w0, #32",
			wantErr: false,
		},
		{
			name: "scvtf	s23, x19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x02, 0x9e}),
				address:          0,
			},
			want: "scvtf	s23, x19, #1",
			wantErr: false,
		},
		{
			name: "scvtf	s31, xzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x02, 0x9e}),
				address:          0,
			},
			want: "scvtf	s31, xzr, #20",
			wantErr: false,
		},
		{
			name: "scvtf	s14, x0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x00, 0x02, 0x9e}),
				address:          0,
			},
			want: "scvtf	s14, x0, #64",
			wantErr: false,
		},
		{
			name: "scvtf	d23, w19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x42, 0x1e}),
				address:          0,
			},
			want: "scvtf	d23, w19, #1",
			wantErr: false,
		},
		{
			name: "scvtf	d31, wzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x42, 0x1e}),
				address:          0,
			},
			want: "scvtf	d31, wzr, #20",
			wantErr: false,
		},
		{
			name: "scvtf	d14, w0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x80, 0x42, 0x1e}),
				address:          0,
			},
			want: "scvtf	d14, w0, #32",
			wantErr: false,
		},
		{
			name: "scvtf	d23, x19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x42, 0x9e}),
				address:          0,
			},
			want: "scvtf	d23, x19, #1",
			wantErr: false,
		},
		{
			name: "scvtf	d31, xzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x42, 0x9e}),
				address:          0,
			},
			want: "scvtf	d31, xzr, #20",
			wantErr: false,
		},
		{
			name: "scvtf	d14, x0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x00, 0x42, 0x9e}),
				address:          0,
			},
			want: "scvtf	d14, x0, #64",
			wantErr: false,
		},
		{
			name: "ucvtf	s23, w19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x03, 0x1e}),
				address:          0,
			},
			want: "ucvtf	s23, w19, #1",
			wantErr: false,
		},
		{
			name: "ucvtf	s31, wzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x03, 0x1e}),
				address:          0,
			},
			want: "ucvtf	s31, wzr, #20",
			wantErr: false,
		},
		{
			name: "ucvtf	s14, w0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x80, 0x03, 0x1e}),
				address:          0,
			},
			want: "ucvtf	s14, w0, #32",
			wantErr: false,
		},
		{
			name: "ucvtf	s23, x19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x03, 0x9e}),
				address:          0,
			},
			want: "ucvtf	s23, x19, #1",
			wantErr: false,
		},
		{
			name: "ucvtf	s31, xzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x03, 0x9e}),
				address:          0,
			},
			want: "ucvtf	s31, xzr, #20",
			wantErr: false,
		},
		{
			name: "ucvtf	s14, x0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x00, 0x03, 0x9e}),
				address:          0,
			},
			want: "ucvtf	s14, x0, #64",
			wantErr: false,
		},
		{
			name: "ucvtf	d23, w19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x43, 0x1e}),
				address:          0,
			},
			want: "ucvtf	d23, w19, #1",
			wantErr: false,
		},
		{
			name: "ucvtf	d31, wzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x43, 0x1e}),
				address:          0,
			},
			want: "ucvtf	d31, wzr, #20",
			wantErr: false,
		},
		{
			name: "ucvtf	d14, w0, #32",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x80, 0x43, 0x1e}),
				address:          0,
			},
			want: "ucvtf	d14, w0, #32",
			wantErr: false,
		},
		{
			name: "ucvtf	d23, x19, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfe, 0x43, 0x9e}),
				address:          0,
			},
			want: "ucvtf	d23, x19, #1",
			wantErr: false,
		},
		{
			name: "ucvtf	d31, xzr, #20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xb3, 0x43, 0x9e}),
				address:          0,
			},
			want: "ucvtf	d31, xzr, #20",
			wantErr: false,
		},
		{
			name: "ucvtf	d14, x0, #64",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0e, 0x00, 0x43, 0x9e}),
				address:          0,
			},
			want: "ucvtf	d14, x0, #64",
			wantErr: false,
		},
		{
			name: "fcvtns	w3, s31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x20, 0x1e}),
				address:          0,
			},
			want: "fcvtns	w3, s31",
			wantErr: false,
		},
		{
			name: "fcvtns	xzr, s12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x20, 0x9e}),
				address:          0,
			},
			want: "fcvtns	xzr, s12",
			wantErr: false,
		},
		{
			name: "fcvtnu	wzr, s12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x21, 0x1e}),
				address:          0,
			},
			want: "fcvtnu	wzr, s12",
			wantErr: false,
		},
		{
			name: "fcvtnu	x0, s0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x21, 0x9e}),
				address:          0,
			},
			want: "fcvtnu	x0, s0",
			wantErr: false,
		},
		{
			name: "fcvtps	wzr, s9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x01, 0x28, 0x1e}),
				address:          0,
			},
			want: "fcvtps	wzr, s9",
			wantErr: false,
		},
		{
			name: "fcvtps	x12, s20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x02, 0x28, 0x9e}),
				address:          0,
			},
			want: "fcvtps	x12, s20",
			wantErr: false,
		},
		{
			name: "fcvtpu	w30, s23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x02, 0x29, 0x1e}),
				address:          0,
			},
			want: "fcvtpu	w30, s23",
			wantErr: false,
		},
		{
			name: "fcvtpu	x29, s3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x00, 0x29, 0x9e}),
				address:          0,
			},
			want: "fcvtpu	x29, s3",
			wantErr: false,
		},
		{
			name: "fcvtms	w2, s3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x30, 0x1e}),
				address:          0,
			},
			want: "fcvtms	w2, s3",
			wantErr: false,
		},
		{
			name: "fcvtms	x4, s5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x00, 0x30, 0x9e}),
				address:          0,
			},
			want: "fcvtms	x4, s5",
			wantErr: false,
		},
		{
			name: "fcvtmu	w6, s7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x00, 0x31, 0x1e}),
				address:          0,
			},
			want: "fcvtmu	w6, s7",
			wantErr: false,
		},
		{
			name: "fcvtmu	x8, s9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x01, 0x31, 0x9e}),
				address:          0,
			},
			want: "fcvtmu	x8, s9",
			wantErr: false,
		},
		{
			name: "fcvtzs	w10, s11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x01, 0x38, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	w10, s11",
			wantErr: false,
		},
		{
			name: "fcvtzs	x12, s13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x01, 0x38, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x12, s13",
			wantErr: false,
		},
		{
			name: "fcvtzu	w14, s15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x01, 0x39, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	w14, s15",
			wantErr: false,
		},
		{
			name: "fcvtzu	x15, s16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0f, 0x02, 0x39, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x15, s16",
			wantErr: false,
		},
		{
			name: "scvtf	s17, w18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x51, 0x02, 0x22, 0x1e}),
				address:          0,
			},
			want: "scvtf	s17, w18",
			wantErr: false,
		},
		{
			name: "scvtf	s19, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x02, 0x22, 0x9e}),
				address:          0,
			},
			want: "scvtf	s19, x20",
			wantErr: false,
		},
		{
			name: "ucvtf	s21, w22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x23, 0x1e}),
				address:          0,
			},
			want: "ucvtf	s21, w22",
			wantErr: false,
		},
		{
			name: "scvtf	s23, x24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x17, 0x03, 0x22, 0x9e}),
				address:          0,
			},
			want: "scvtf	s23, x24",
			wantErr: false,
		},
		{
			name: "fcvtas	w25, s26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x59, 0x03, 0x24, 0x1e}),
				address:          0,
			},
			want: "fcvtas	w25, s26",
			wantErr: false,
		},
		{
			name: "fcvtas	x27, s28",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x03, 0x24, 0x9e}),
				address:          0,
			},
			want: "fcvtas	x27, s28",
			wantErr: false,
		},
		{
			name: "fcvtau	w29, s30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdd, 0x03, 0x25, 0x1e}),
				address:          0,
			},
			want: "fcvtau	w29, s30",
			wantErr: false,
		},
		{
			name: "fcvtau	xzr, s0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x25, 0x9e}),
				address:          0,
			},
			want: "fcvtau	xzr, s0",
			wantErr: false,
		},
		{
			name: "fcvtns	w3, d31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x60, 0x1e}),
				address:          0,
			},
			want: "fcvtns	w3, d31",
			wantErr: false,
		},
		{
			name: "fcvtns	xzr, d12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x60, 0x9e}),
				address:          0,
			},
			want: "fcvtns	xzr, d12",
			wantErr: false,
		},
		{
			name: "fcvtnu	wzr, d12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x01, 0x61, 0x1e}),
				address:          0,
			},
			want: "fcvtnu	wzr, d12",
			wantErr: false,
		},
		{
			name: "fcvtnu	x0, d0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x61, 0x9e}),
				address:          0,
			},
			want: "fcvtnu	x0, d0",
			wantErr: false,
		},
		{
			name: "fcvtps	wzr, d9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x01, 0x68, 0x1e}),
				address:          0,
			},
			want: "fcvtps	wzr, d9",
			wantErr: false,
		},
		{
			name: "fcvtps	x12, d20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x02, 0x68, 0x9e}),
				address:          0,
			},
			want: "fcvtps	x12, d20",
			wantErr: false,
		},
		{
			name: "fcvtpu	w30, d23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x02, 0x69, 0x1e}),
				address:          0,
			},
			want: "fcvtpu	w30, d23",
			wantErr: false,
		},
		{
			name: "fcvtpu	x29, d3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7d, 0x00, 0x69, 0x9e}),
				address:          0,
			},
			want: "fcvtpu	x29, d3",
			wantErr: false,
		},
		{
			name: "fcvtms	w2, d3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x00, 0x70, 0x1e}),
				address:          0,
			},
			want: "fcvtms	w2, d3",
			wantErr: false,
		},
		{
			name: "fcvtms	x4, d5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x00, 0x70, 0x9e}),
				address:          0,
			},
			want: "fcvtms	x4, d5",
			wantErr: false,
		},
		{
			name: "fcvtmu	w6, d7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x00, 0x71, 0x1e}),
				address:          0,
			},
			want: "fcvtmu	w6, d7",
			wantErr: false,
		},
		{
			name: "fcvtmu	x8, d9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x01, 0x71, 0x9e}),
				address:          0,
			},
			want: "fcvtmu	x8, d9",
			wantErr: false,
		},
		{
			name: "fcvtzs	w10, d11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x01, 0x78, 0x1e}),
				address:          0,
			},
			want: "fcvtzs	w10, d11",
			wantErr: false,
		},
		{
			name: "fcvtzs	x12, d13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x01, 0x78, 0x9e}),
				address:          0,
			},
			want: "fcvtzs	x12, d13",
			wantErr: false,
		},
		{
			name: "fcvtzu	w14, d15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x01, 0x79, 0x1e}),
				address:          0,
			},
			want: "fcvtzu	w14, d15",
			wantErr: false,
		},
		{
			name: "fcvtzu	x15, d16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0f, 0x02, 0x79, 0x9e}),
				address:          0,
			},
			want: "fcvtzu	x15, d16",
			wantErr: false,
		},
		{
			name: "scvtf	d17, w18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x51, 0x02, 0x62, 0x1e}),
				address:          0,
			},
			want: "scvtf	d17, w18",
			wantErr: false,
		},
		{
			name: "scvtf	d19, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x02, 0x62, 0x9e}),
				address:          0,
			},
			want: "scvtf	d19, x20",
			wantErr: false,
		},
		{
			name: "ucvtf	d21, w22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x02, 0x63, 0x1e}),
				address:          0,
			},
			want: "ucvtf	d21, w22",
			wantErr: false,
		},
		{
			name: "ucvtf	d23, x24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x17, 0x03, 0x63, 0x9e}),
				address:          0,
			},
			want: "ucvtf	d23, x24",
			wantErr: false,
		},
		{
			name: "fcvtas	w25, d26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x59, 0x03, 0x64, 0x1e}),
				address:          0,
			},
			want: "fcvtas	w25, d26",
			wantErr: false,
		},
		{
			name: "fcvtas	x27, d28",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9b, 0x03, 0x64, 0x9e}),
				address:          0,
			},
			want: "fcvtas	x27, d28",
			wantErr: false,
		},
		{
			name: "fcvtau	w29, d30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdd, 0x03, 0x65, 0x1e}),
				address:          0,
			},
			want: "fcvtau	w29, d30",
			wantErr: false,
		},
		{
			name: "fcvtau	xzr, d0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x00, 0x65, 0x9e}),
				address:          0,
			},
			want: "fcvtau	xzr, d0",
			wantErr: false,
		},
		{
			name: "fmov	w3, s9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x01, 0x26, 0x1e}),
				address:          0,
			},
			want: "fmov	w3, s9",
			wantErr: false,
		},
		{
			name: "fmov	s9, w3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x00, 0x27, 0x1e}),
				address:          0,
			},
			want: "fmov	s9, w3",
			wantErr: false,
		},
		{
			name: "fmov	x20, d31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x03, 0x66, 0x9e}),
				address:          0,
			},
			want: "fmov	x20, d31",
			wantErr: false,
		},
		{
			name: "fmov	d1, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0x01, 0x67, 0x9e}),
				address:          0,
			},
			want: "fmov	d1, x15",
			wantErr: false,
		},
		{
			name: "fmov	x3, v12.d[1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x01, 0xae, 0x9e}),
				address:          0,
			},
			want: "fmov	x3, v12.d[1]",
			wantErr: false,
		},
		{
			name: "fmov	v1.d[1], x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x02, 0xaf, 0x9e}),
				address:          0,
			},
			want: "fmov	v1.d[1], x19",
			wantErr: false,
		},
		{
			name: "fmov	v3.d[1], xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xaf, 0x9e}),
				address:          0,
			},
			want: "fmov	v3.d[1], xzr",
			wantErr: false,
		},
		{
			name: "fmov	s2, #0.12500000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x02, 0x10, 0x28, 0x1e}),
				address:          0,
			},
			want: "fmov	s2, #0.12500000",
			wantErr: false,
		},
		{
			name: "fmov	s3, #1.00000000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x03, 0x10, 0x2e, 0x1e}),
				address:          0,
			},
			want: "fmov	s3, #1.00000000",
			wantErr: false,
		},
		{
			name: "fmov	d30, #16.00000000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1e, 0x10, 0x66, 0x1e}),
				address:          0,
			},
			want: "fmov	d30, #16.00000000",
			wantErr: false,
		},
		{
			name: "fmov	s4, #1.06250000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x04, 0x30, 0x2e, 0x1e}),
				address:          0,
			},
			want: "fmov	s4, #1.06250000",
			wantErr: false,
		},
		{
			name: "fmov	d10, #1.93750000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0a, 0xf0, 0x6f, 0x1e}),
				address:          0,
			},
			want: "fmov	d10, #1.93750000",
			wantErr: false,
		},
		{
			name: "fmov	s12, #-1.00000000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x10, 0x3e, 0x1e}),
				address:          0,
			},
			want: "fmov	s12, #-1.00000000",
			wantErr: false,
		},
		{
			name: "fmov	d16, #8.50000000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x30, 0x64, 0x1e}),
				address:          0,
			},
			want: "fmov	d16, #8.50000000",
			wantErr: false,
		},
		{
			name: "ldr	w0, #1048572",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0xff, 0x7f, 0x18}),
				address:          0,
			},
			want: "ldr	w0, #1048572",
			wantErr: false,
		},
		{
			name: "ldr	x10, #-1048576",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0a, 0x00, 0x80, 0x58}),
				address:          0,
			},
			want: "ldr	x10, #-1048576",
			wantErr: false,
		},
		{
			name: "stxrb	w1, w2, [x3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x7c, 0x01, 0x08}),
				address:          0,
			},
			want: "stxrb	w1, w2, [x3]",
			wantErr: false,
		},
		{
			name: "stxrh	w2, w3, [x4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x7c, 0x02, 0x48}),
				address:          0,
			},
			want: "stxrh	w2, w3, [x4]",
			wantErr: false,
		},
		{
			name: "stxr	wzr, w4, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe4, 0x7f, 0x1f, 0x88}),
				address:          0,
			},
			want: "stxr	wzr, w4, [sp]",
			wantErr: false,
		},
		{
			name: "stxr	w5, x6, [x7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe6, 0x7c, 0x05, 0xc8}),
				address:          0,
			},
			want: "stxr	w5, x6, [x7]",
			wantErr: false,
		},
		{
			name: "ldxrb	w7, [x9]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x27, 0x7d, 0x5f, 0x08}),
				address:          0,
			},
			want: "ldxrb	w7, [x9]",
			wantErr: false,
		},
		{
			name: "ldxrh	wzr, [x10]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x7d, 0x5f, 0x48}),
				address:          0,
			},
			want: "ldxrh	wzr, [x10]",
			wantErr: false,
		},
		{
			name: "ldxr	w9, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x7f, 0x5f, 0x88}),
				address:          0,
			},
			want: "ldxr	w9, [sp]",
			wantErr: false,
		},
		{
			name: "ldxr	x10, [x11]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x7d, 0x5f, 0xc8}),
				address:          0,
			},
			want: "ldxr	x10, [x11]",
			wantErr: false,
		},
		{
			name: "stxp	w11, w12, w13, [x14]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x35, 0x2b, 0x88}),
				address:          0,
			},
			want: "stxp	w11, w12, w13, [x14]",
			wantErr: false,
		},
		{
			name: "stxp	wzr, x23, x14, [x15]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x39, 0x3f, 0xc8}),
				address:          0,
			},
			want: "stxp	wzr, x23, x14, [x15]",
			wantErr: false,
		},
		{
			name: "ldxp	w12, wzr, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x7f, 0x7f, 0x88}),
				address:          0,
			},
			want: "ldxp	w12, wzr, [sp]",
			wantErr: false,
		},
		{
			name: "ldxp	x13, x14, [x15]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x39, 0x7f, 0xc8}),
				address:          0,
			},
			want: "ldxp	x13, x14, [x15]",
			wantErr: false,
		},
		{
			name: "stlxrb	w14, w15, [x16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0f, 0xfe, 0x0e, 0x08}),
				address:          0,
			},
			want: "stlxrb	w14, w15, [x16]",
			wantErr: false,
		},
		{
			name: "stlxrh	w15, w16, [x17]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0xfe, 0x0f, 0x48}),
				address:          0,
			},
			want: "stlxrh	w15, w16, [x17]",
			wantErr: false,
		},
		{
			name: "stlxr	wzr, w17, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xff, 0x1f, 0x88}),
				address:          0,
			},
			want: "stlxr	wzr, w17, [sp]",
			wantErr: false,
		},
		{
			name: "stlxr	w18, x19, [x20]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xfe, 0x12, 0xc8}),
				address:          0,
			},
			want: "stlxr	w18, x19, [x20]",
			wantErr: false,
		},
		{
			name: "ldaxrb	w19, [x21]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb3, 0xfe, 0x5f, 0x08}),
				address:          0,
			},
			want: "ldaxrb	w19, [x21]",
			wantErr: false,
		},
		{
			name: "ldaxrh	w20, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xff, 0x5f, 0x48}),
				address:          0,
			},
			want: "ldaxrh	w20, [sp]",
			wantErr: false,
		},
		{
			name: "ldaxr	wzr, [x22]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0xfe, 0x5f, 0x88}),
				address:          0,
			},
			want: "ldaxr	wzr, [x22]",
			wantErr: false,
		},
		{
			name: "ldaxr	x21, [x23]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf5, 0xfe, 0x5f, 0xc8}),
				address:          0,
			},
			want: "ldaxr	x21, [x23]",
			wantErr: false,
		},
		{
			name: "stlxp	wzr, w22, w23, [x24]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x16, 0xdf, 0x3f, 0x88}),
				address:          0,
			},
			want: "stlxp	wzr, w22, w23, [x24]",
			wantErr: false,
		},
		{
			name: "stlxp	w25, x26, x27, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0xef, 0x39, 0xc8}),
				address:          0,
			},
			want: "stlxp	w25, x26, x27, [sp]",
			wantErr: false,
		},
		{
			name: "ldaxp	w26, wzr, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfa, 0xff, 0x7f, 0x88}),
				address:          0,
			},
			want: "ldaxp	w26, wzr, [sp]",
			wantErr: false,
		},
		{
			name: "ldaxp	x27, x28, [x30]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdb, 0xf3, 0x7f, 0xc8}),
				address:          0,
			},
			want: "ldaxp	x27, x28, [x30]",
			wantErr: false,
		},
		{
			name: "stlrb	w27, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfb, 0xff, 0x9f, 0x08}),
				address:          0,
			},
			want: "stlrb	w27, [sp]",
			wantErr: false,
		},
		{
			name: "stlrh	w28, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1c, 0xfc, 0x9f, 0x48}),
				address:          0,
			},
			want: "stlrh	w28, [x0]",
			wantErr: false,
		},
		{
			name: "stlr	wzr, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfc, 0x9f, 0x88}),
				address:          0,
			},
			want: "stlr	wzr, [x1]",
			wantErr: false,
		},
		{
			name: "stlr	x30, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5e, 0xfc, 0x9f, 0xc8}),
				address:          0,
			},
			want: "stlr	x30, [x2]",
			wantErr: false,
		},
		{
			name: "ldarb	w29, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0xff, 0xdf, 0x08}),
				address:          0,
			},
			want: "ldarb	w29, [sp]",
			wantErr: false,
		},
		{
			name: "ldarh	w30, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1e, 0xfc, 0xdf, 0x48}),
				address:          0,
			},
			want: "ldarh	w30, [x0]",
			wantErr: false,
		},
		{
			name: "ldar	wzr, [x1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfc, 0xdf, 0x88}),
				address:          0,
			},
			want: "ldar	wzr, [x1]",
			wantErr: false,
		},
		{
			name: "ldar	x1, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xfc, 0xdf, 0xc8}),
				address:          0,
			},
			want: "ldar	x1, [x2]",
			wantErr: false,
		},
		{
			name: "stlxp	wzr, w22, w23, [x24]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x16, 0xdf, 0x3f, 0x88}),
				address:          0,
			},
			want: "stlxp	wzr, w22, w23, [x24]",
			wantErr: false,
		},
		{
			name: "sturb	w9, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x00, 0x38}),
				address:          0,
			},
			want: "sturb	w9, [sp]",
			wantErr: false,
		},
		{
			name: "sturh	wzr, [x12, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xf1, 0x0f, 0x78}),
				address:          0,
			},
			want: "sturh	wzr, [x12, #255]",
			wantErr: false,
		},
		{
			name: "stur	w16, [x0, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x00, 0x10, 0xb8}),
				address:          0,
			},
			want: "stur	w16, [x0, #-256]",
			wantErr: false,
		},
		{
			name: "stur	x28, [x14, #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdc, 0x11, 0x00, 0xf8}),
				address:          0,
			},
			want: "stur	x28, [x14, #1]",
			wantErr: false,
		},
		{
			name: "ldurb	w1, [x20, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0xf2, 0x4f, 0x38}),
				address:          0,
			},
			want: "ldurb	w1, [x20, #255]",
			wantErr: false,
		},
		{
			name: "ldurh	w20, [x1, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x34, 0xf0, 0x4f, 0x78}),
				address:          0,
			},
			want: "ldurh	w20, [x1, #255]",
			wantErr: false,
		},
		{
			name: "ldur	w12, [sp, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xf3, 0x4f, 0xb8}),
				address:          0,
			},
			want: "ldur	w12, [sp, #255]",
			wantErr: false,
		},
		{
			name: "ldur	xzr, [x12, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xf1, 0x4f, 0xf8}),
				address:          0,
			},
			want: "ldur	xzr, [x12, #255]",
			wantErr: false,
		},
		{
			name: "ldursb	x9, [x7, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x00, 0x90, 0x38}),
				address:          0,
			},
			want: "ldursb	x9, [x7, #-256]",
			wantErr: false,
		},
		{
			name: "ldursh	x17, [x19, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x71, 0x02, 0x90, 0x78}),
				address:          0,
			},
			want: "ldursh	x17, [x19, #-256]",
			wantErr: false,
		},
		{
			name: "ldursw	x20, [x15, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x01, 0x90, 0xb8}),
				address:          0,
			},
			want: "ldursw	x20, [x15, #-256]",
			wantErr: false,
		},
		{
			name: "ldursw	x13, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4d, 0x00, 0x80, 0xb8}),
				address:          0,
			},
			want: "ldursw	x13, [x2]",
			wantErr: false,
		},
		{
			name: "prfum	pldl2keep, [sp, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x03, 0x90, 0xf8}),
				address:          0,
			},
			want: "prfum	pldl2keep, [sp, #-256]",
			wantErr: false,
		},
		{
			name: "ldursb	w19, [x1, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x00, 0xd0, 0x38}),
				address:          0,
			},
			want: "ldursb	w19, [x1, #-256]",
			wantErr: false,
		},
		{
			name: "ldursh	w15, [x21, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaf, 0x02, 0xd0, 0x78}),
				address:          0,
			},
			want: "ldursh	w15, [x21, #-256]",
			wantErr: false,
		},
		{
			name: "stur	b0, [sp, #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x13, 0x00, 0x3c}),
				address:          0,
			},
			want: "stur	b0, [sp, #1]",
			wantErr: false,
		},
		{
			name: "stur	h12, [x12, #-1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xf1, 0x1f, 0x7c}),
				address:          0,
			},
			want: "stur	h12, [x12, #-1]",
			wantErr: false,
		},
		{
			name: "stur	s15, [x0, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0f, 0xf0, 0x0f, 0xbc}),
				address:          0,
			},
			want: "stur	s15, [x0, #255]",
			wantErr: false,
		},
		{
			name: "stur	d31, [x5, #25]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x90, 0x01, 0xfc}),
				address:          0,
			},
			want: "stur	d31, [x5, #25]",
			wantErr: false,
		},
		{
			name: "stur	q9, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x00, 0x80, 0x3c}),
				address:          0,
			},
			want: "stur	q9, [x5]",
			wantErr: false,
		},
		{
			name: "ldur	b3, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x40, 0x3c}),
				address:          0,
			},
			want: "ldur	b3, [sp]",
			wantErr: false,
		},
		{
			name: "ldur	h5, [x4, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x85, 0x00, 0x50, 0x7c}),
				address:          0,
			},
			want: "ldur	h5, [x4, #-256]",
			wantErr: false,
		},
		{
			name: "ldur	s7, [x12, #-1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x87, 0xf1, 0x5f, 0xbc}),
				address:          0,
			},
			want: "ldur	s7, [x12, #-1]",
			wantErr: false,
		},
		{
			name: "ldur	d11, [x19, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6b, 0x42, 0x40, 0xfc}),
				address:          0,
			},
			want: "ldur	d11, [x19, #4]",
			wantErr: false,
		},
		{
			name: "ldur	q13, [x1, #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2d, 0x20, 0xc0, 0x3c}),
				address:          0,
			},
			want: "ldur	q13, [x1, #2]",
			wantErr: false,
		},
		{
			name: "ldr	x0, [x0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x40, 0xf9}),
				address:          0,
			},
			want: "ldr	x0, [x0]",
			wantErr: false,
		},
		{
			name: "ldr	x4, [x29]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x03, 0x40, 0xf9}),
				address:          0,
			},
			want: "ldr	x4, [x29]",
			wantErr: false,
		},
		{
			name: "ldr	x30, [x12, #32760]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9e, 0xfd, 0x7f, 0xf9}),
				address:          0,
			},
			want: "ldr	x30, [x12, #32760]",
			wantErr: false,
		},
		{
			name: "ldr	x20, [sp, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x07, 0x40, 0xf9}),
				address:          0,
			},
			want: "ldr	x20, [sp, #8]",
			wantErr: false,
		},
		{
			name: "ldr	xzr, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x40, 0xf9}),
				address:          0,
			},
			want: "ldr	xzr, [sp]",
			wantErr: false,
		},
		{
			name: "ldr	w2, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x03, 0x40, 0xb9}),
				address:          0,
			},
			want: "ldr	w2, [sp]",
			wantErr: false,
		},
		{
			name: "ldr	w17, [sp, #16380]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xff, 0x7f, 0xb9}),
				address:          0,
			},
			want: "ldr	w17, [sp, #16380]",
			wantErr: false,
		},
		{
			name: "ldr	w13, [x2, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4d, 0x04, 0x40, 0xb9}),
				address:          0,
			},
			want: "ldr	w13, [x2, #4]",
			wantErr: false,
		},
		{
			name: "ldrsw	x2, [x5, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x04, 0x80, 0xb9}),
				address:          0,
			},
			want: "ldrsw	x2, [x5, #4]",
			wantErr: false,
		},
		{
			name: "ldrsw	x23, [sp, #16380]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0xff, 0xbf, 0xb9}),
				address:          0,
			},
			want: "ldrsw	x23, [sp, #16380]",
			wantErr: false,
		},
		{
			name: "ldrh	w2, [x4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x82, 0x00, 0x40, 0x79}),
				address:          0,
			},
			want: "ldrh	w2, [x4]",
			wantErr: false,
		},
		{
			name: "ldrsh	w23, [x6, #8190]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd7, 0xfc, 0xff, 0x79}),
				address:          0,
			},
			want: "ldrsh	w23, [x6, #8190]",
			wantErr: false,
		},
		{
			name: "ldrsh	wzr, [sp, #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x07, 0xc0, 0x79}),
				address:          0,
			},
			want: "ldrsh	wzr, [sp, #2]",
			wantErr: false,
		},
		{
			name: "ldrsh	x29, [x2, #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5d, 0x04, 0x80, 0x79}),
				address:          0,
			},
			want: "ldrsh	x29, [x2, #2]",
			wantErr: false,
		},
		{
			name: "ldrb	w26, [x3, #121]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7a, 0xe4, 0x41, 0x39}),
				address:          0,
			},
			want: "ldrb	w26, [x3, #121]",
			wantErr: false,
		},
		{
			name: "ldrb	w12, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x00, 0x40, 0x39}),
				address:          0,
			},
			want: "ldrb	w12, [x2]",
			wantErr: false,
		},
		{
			name: "ldrsb	w27, [sp, #4095]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfb, 0xff, 0xff, 0x39}),
				address:          0,
			},
			want: "ldrsb	w27, [sp, #4095]",
			wantErr: false,
		},
		{
			name: "ldrsb	xzr, [x15]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x01, 0x80, 0x39}),
				address:          0,
			},
			want: "ldrsb	xzr, [x15]",
			wantErr: false,
		},
		{
			name: "str	x30, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfe, 0x03, 0x00, 0xf9}),
				address:          0,
			},
			want: "str	x30, [sp]",
			wantErr: false,
		},
		{
			name: "str	w20, [x4, #16380]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xfc, 0x3f, 0xb9}),
				address:          0,
			},
			want: "str	w20, [x4, #16380]",
			wantErr: false,
		},
		{
			name: "strh	w20, [x10, #14]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x54, 0x1d, 0x00, 0x79}),
				address:          0,
			},
			want: "strh	w20, [x10, #14]",
			wantErr: false,
		},
		{
			name: "strh	w17, [sp, #8190]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xff, 0x3f, 0x79}),
				address:          0,
			},
			want: "strh	w17, [sp, #8190]",
			wantErr: false,
		},
		{
			name: "strb	w23, [x3, #4095]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x77, 0xfc, 0x3f, 0x39}),
				address:          0,
			},
			want: "strb	w23, [x3, #4095]",
			wantErr: false,
		},
		{
			name: "strb	wzr, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x00, 0x39}),
				address:          0,
			},
			want: "strb	wzr, [x2]",
			wantErr: false,
		},
		{
			name: "prfm	pldl1keep, [sp, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x07, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pldl1keep, [sp, #8]",
			wantErr: false,
		},
		{
			name: "prfm	pldl1strm, [x3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pldl1strm, [x3]",
			wantErr: false,
		},
		{
			name: "prfm	pldl2keep, [x5, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x08, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pldl2keep, [x5, #16]",
			wantErr: false,
		},
		{
			name: "prfm	pldl2strm, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x43, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pldl2strm, [x2]",
			wantErr: false,
		},
		{
			name: "prfm	pldl3keep, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa4, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pldl3keep, [x5]",
			wantErr: false,
		},
		{
			name: "prfm	pldl3strm, [x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pldl3strm, [x6]",
			wantErr: false,
		},
		{
			name: "prfm	plil1keep, [sp, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe8, 0x07, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	plil1keep, [sp, #8]",
			wantErr: false,
		},
		{
			name: "prfm	plil1strm, [x3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	plil1strm, [x3]",
			wantErr: false,
		},
		{
			name: "prfm	plil2keep, [x5, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaa, 0x08, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	plil2keep, [x5, #16]",
			wantErr: false,
		},
		{
			name: "prfm	plil2strm, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4b, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	plil2strm, [x2]",
			wantErr: false,
		},
		{
			name: "prfm	plil3keep, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	plil3keep, [x5]",
			wantErr: false,
		},
		{
			name: "prfm	plil3strm, [x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcd, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	plil3strm, [x6]",
			wantErr: false,
		},
		{
			name: "prfm	pstl1keep, [sp, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf0, 0x07, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pstl1keep, [sp, #8]",
			wantErr: false,
		},
		{
			name: "prfm	pstl1strm, [x3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x71, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pstl1strm, [x3]",
			wantErr: false,
		},
		{
			name: "prfm	pstl2keep, [x5, #16]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb2, 0x08, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pstl2keep, [x5, #16]",
			wantErr: false,
		},
		{
			name: "prfm	pstl2strm, [x2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x53, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pstl2strm, [x2]",
			wantErr: false,
		},
		{
			name: "prfm	pstl3keep, [x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pstl3keep, [x5]",
			wantErr: false,
		},
		{
			name: "prfm	pstl3strm, [x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0x00, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	pstl3strm, [x6]",
			wantErr: false,
		},
		{
			name: "prfm	#15, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0x03, 0x80, 0xf9}),
				address:          0,
			},
			want: "prfm	#15, [sp]",
			wantErr: false,
		},
		{
			name: "ldr	b31, [sp, #4095]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xff, 0x7f, 0x3d}),
				address:          0,
			},
			want: "ldr	b31, [sp, #4095]",
			wantErr: false,
		},
		{
			name: "ldr	h20, [x2, #8190]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x54, 0xfc, 0x7f, 0x7d}),
				address:          0,
			},
			want: "ldr	h20, [x2, #8190]",
			wantErr: false,
		},
		{
			name: "ldr	s10, [x19, #16380]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0xfe, 0x7f, 0xbd}),
				address:          0,
			},
			want: "ldr	s10, [x19, #16380]",
			wantErr: false,
		},
		{
			name: "ldr	d3, [x10, #32760]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x43, 0xfd, 0x7f, 0xfd}),
				address:          0,
			},
			want: "ldr	d3, [x10, #32760]",
			wantErr: false,
		},
		{
			name: "str	q12, [sp, #65520]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xff, 0xbf, 0x3d}),
				address:          0,
			},
			want: "str	q12, [sp, #65520]",
			wantErr: false,
		},
		{
			name: "ldrb	w3, [sp, x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x6b, 0x65, 0x38}),
				address:          0,
			},
			want: "ldrb	w3, [sp, x5]",
			wantErr: false,
		},
		{
			name: "ldrb	w9, [x27, x6, lsl #0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x7b, 0x66, 0x38}),
				address:          0,
			},
			want: "ldrb	w9, [x27, x6, lsl #0]",
			wantErr: false,
		},
		{
			name: "ldrsb	w10, [x30, x7]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x6b, 0xe7, 0x38}),
				address:          0,
			},
			want: "ldrsb	w10, [x30, x7]",
			wantErr: false,
		},
		{
			name: "ldrb	w11, [x29, x3, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xeb, 0x63, 0x38}),
				address:          0,
			},
			want: "ldrb	w11, [x29, x3, sxtx]",
			wantErr: false,
		},
		{
			name: "strb	w12, [x28, xzr, sxtx #0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xfb, 0x3f, 0x38}),
				address:          0,
			},
			want: "strb	w12, [x28, xzr, sxtx #0]",
			wantErr: false,
		},
		{
			name: "ldrb	w14, [x26, w6, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4e, 0x4b, 0x66, 0x38}),
				address:          0,
			},
			want: "ldrb	w14, [x26, w6, uxtw]",
			wantErr: false,
		},
		{
			name: "ldrsb	w15, [x25, w7, uxtw #0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2f, 0x5b, 0xe7, 0x38}),
				address:          0,
			},
			want: "ldrsb	w15, [x25, w7, uxtw #0]",
			wantErr: false,
		},
		{
			name: "ldrb	w17, [x23, w9, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xca, 0x69, 0x38}),
				address:          0,
			},
			want: "ldrb	w17, [x23, w9, sxtw]",
			wantErr: false,
		},
		{
			name: "ldrsb	x18, [x22, w10, sxtw #0]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd2, 0xda, 0xaa, 0x38}),
				address:          0,
			},
			want: "ldrsb	x18, [x22, w10, sxtw #0]",
			wantErr: false,
		},
		{
			name: "ldrsh	w3, [sp, x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x6b, 0xe5, 0x78}),
				address:          0,
			},
			want: "ldrsh	w3, [sp, x5]",
			wantErr: false,
		},
		{
			name: "ldrsh	w9, [x27, x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x6b, 0xe6, 0x78}),
				address:          0,
			},
			want: "ldrsh	w9, [x27, x6]",
			wantErr: false,
		},
		{
			name: "ldrh	w10, [x30, x7, lsl #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x7b, 0x67, 0x78}),
				address:          0,
			},
			want: "ldrh	w10, [x30, x7, lsl #1]",
			wantErr: false,
		},
		{
			name: "strh	w11, [x29, x3, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xeb, 0x23, 0x78}),
				address:          0,
			},
			want: "strh	w11, [x29, x3, sxtx]",
			wantErr: false,
		},
		{
			name: "ldrh	w12, [x28, xzr, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xeb, 0x7f, 0x78}),
				address:          0,
			},
			want: "ldrh	w12, [x28, xzr, sxtx]",
			wantErr: false,
		},
		{
			name: "ldrsh	x13, [x27, x5, sxtx #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6d, 0xfb, 0xa5, 0x78}),
				address:          0,
			},
			want: "ldrsh	x13, [x27, x5, sxtx #1]",
			wantErr: false,
		},
		{
			name: "ldrh	w14, [x26, w6, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4e, 0x4b, 0x66, 0x78}),
				address:          0,
			},
			want: "ldrh	w14, [x26, w6, uxtw]",
			wantErr: false,
		},
		{
			name: "ldrh	w15, [x25, w7, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2f, 0x4b, 0x67, 0x78}),
				address:          0,
			},
			want: "ldrh	w15, [x25, w7, uxtw]",
			wantErr: false,
		},
		{
			name: "ldrsh	w16, [x24, w8, uxtw #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x5b, 0xe8, 0x78}),
				address:          0,
			},
			want: "ldrsh	w16, [x24, w8, uxtw #1]",
			wantErr: false,
		},
		{
			name: "ldrh	w17, [x23, w9, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xca, 0x69, 0x78}),
				address:          0,
			},
			want: "ldrh	w17, [x23, w9, sxtw]",
			wantErr: false,
		},
		{
			name: "ldrh	w18, [x22, w10, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd2, 0xca, 0x6a, 0x78}),
				address:          0,
			},
			want: "ldrh	w18, [x22, w10, sxtw]",
			wantErr: false,
		},
		{
			name: "strh	w19, [x21, wzr, sxtw #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb3, 0xda, 0x3f, 0x78}),
				address:          0,
			},
			want: "strh	w19, [x21, wzr, sxtw #1]",
			wantErr: false,
		},
		{
			name: "ldr	w3, [sp, x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x6b, 0x65, 0xb8}),
				address:          0,
			},
			want: "ldr	w3, [sp, x5]",
			wantErr: false,
		},
		{
			name: "ldr	s9, [x27, x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x6b, 0x66, 0xbc}),
				address:          0,
			},
			want: "ldr	s9, [x27, x6]",
			wantErr: false,
		},
		{
			name: "ldr	w10, [x30, x7, lsl #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x7b, 0x67, 0xb8}),
				address:          0,
			},
			want: "ldr	w10, [x30, x7, lsl #2]",
			wantErr: false,
		},
		{
			name: "ldr	w11, [x29, x3, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xeb, 0x63, 0xb8}),
				address:          0,
			},
			want: "ldr	w11, [x29, x3, sxtx]",
			wantErr: false,
		},
		{
			name: "str	s12, [x28, xzr, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xeb, 0x3f, 0xbc}),
				address:          0,
			},
			want: "str	s12, [x28, xzr, sxtx]",
			wantErr: false,
		},
		{
			name: "str	w13, [x27, x5, sxtx #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6d, 0xfb, 0x25, 0xb8}),
				address:          0,
			},
			want: "str	w13, [x27, x5, sxtx #2]",
			wantErr: false,
		},
		{
			name: "str	w14, [x26, w6, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4e, 0x4b, 0x26, 0xb8}),
				address:          0,
			},
			want: "str	w14, [x26, w6, uxtw]",
			wantErr: false,
		},
		{
			name: "ldr	w15, [x25, w7, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2f, 0x4b, 0x67, 0xb8}),
				address:          0,
			},
			want: "ldr	w15, [x25, w7, uxtw]",
			wantErr: false,
		},
		{
			name: "ldr	w16, [x24, w8, uxtw #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x5b, 0x68, 0xb8}),
				address:          0,
			},
			want: "ldr	w16, [x24, w8, uxtw #2]",
			wantErr: false,
		},
		{
			name: "ldrsw	x17, [x23, w9, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xca, 0xa9, 0xb8}),
				address:          0,
			},
			want: "ldrsw	x17, [x23, w9, sxtw]",
			wantErr: false,
		},
		{
			name: "ldr	w18, [x22, w10, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd2, 0xca, 0x6a, 0xb8}),
				address:          0,
			},
			want: "ldr	w18, [x22, w10, sxtw]",
			wantErr: false,
		},
		{
			name: "ldrsw	x19, [x21, wzr, sxtw #2]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb3, 0xda, 0xbf, 0xb8}),
				address:          0,
			},
			want: "ldrsw	x19, [x21, wzr, sxtw #2]",
			wantErr: false,
		},
		{
			name: "ldr	x3, [sp, x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x6b, 0x65, 0xf8}),
				address:          0,
			},
			want: "ldr	x3, [sp, x5]",
			wantErr: false,
		},
		{
			name: "str	x9, [x27, x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x6b, 0x26, 0xf8}),
				address:          0,
			},
			want: "str	x9, [x27, x6]",
			wantErr: false,
		},
		{
			name: "ldr	d10, [x30, x7, lsl #3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x7b, 0x67, 0xfc}),
				address:          0,
			},
			want: "ldr	d10, [x30, x7, lsl #3]",
			wantErr: false,
		},
		{
			name: "str	x11, [x29, x3, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xeb, 0x23, 0xf8}),
				address:          0,
			},
			want: "str	x11, [x29, x3, sxtx]",
			wantErr: false,
		},
		{
			name: "ldr	x12, [x28, xzr, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xeb, 0x7f, 0xf8}),
				address:          0,
			},
			want: "ldr	x12, [x28, xzr, sxtx]",
			wantErr: false,
		},
		{
			name: "ldr	x13, [x27, x5, sxtx #3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6d, 0xfb, 0x65, 0xf8}),
				address:          0,
			},
			want: "ldr	x13, [x27, x5, sxtx #3]",
			wantErr: false,
		},
		{
			name: "prfm	pldl1keep, [x26, w6, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x4b, 0xa6, 0xf8}),
				address:          0,
			},
			want: "prfm	pldl1keep, [x26, w6, uxtw]",
			wantErr: false,
		},
		{
			name: "ldr	x15, [x25, w7, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2f, 0x4b, 0x67, 0xf8}),
				address:          0,
			},
			want: "ldr	x15, [x25, w7, uxtw]",
			wantErr: false,
		},
		{
			name: "ldr	x16, [x24, w8, uxtw #3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x5b, 0x68, 0xf8}),
				address:          0,
			},
			want: "ldr	x16, [x24, w8, uxtw #3]",
			wantErr: false,
		},
		{
			name: "ldr	x17, [x23, w9, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xca, 0x69, 0xf8}),
				address:          0,
			},
			want: "ldr	x17, [x23, w9, sxtw]",
			wantErr: false,
		},
		{
			name: "ldr	x18, [x22, w10, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd2, 0xca, 0x6a, 0xf8}),
				address:          0,
			},
			want: "ldr	x18, [x22, w10, sxtw]",
			wantErr: false,
		},
		{
			name: "str	d19, [x21, wzr, sxtw #3]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb3, 0xda, 0x3f, 0xfc}),
				address:          0,
			},
			want: "str	d19, [x21, wzr, sxtw #3]",
			wantErr: false,
		},
		{
			name: "prfm	#6, [x0, x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x06, 0x68, 0xa5, 0xf8}),
				address:          0,
			},
			want: "prfm	#6, [x0, x5]",
			wantErr: false,
		},
		{
			name: "ldr	q3, [sp, x5]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x6b, 0xe5, 0x3c}),
				address:          0,
			},
			want: "ldr	q3, [sp, x5]",
			wantErr: false,
		},
		{
			name: "ldr	q9, [x27, x6]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x6b, 0xe6, 0x3c}),
				address:          0,
			},
			want: "ldr	q9, [x27, x6]",
			wantErr: false,
		},
		{
			name: "ldr	q10, [x30, x7, lsl #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xca, 0x7b, 0xe7, 0x3c}),
				address:          0,
			},
			want: "ldr	q10, [x30, x7, lsl #4]",
			wantErr: false,
		},
		{
			name: "str	q11, [x29, x3, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0xeb, 0xa3, 0x3c}),
				address:          0,
			},
			want: "str	q11, [x29, x3, sxtx]",
			wantErr: false,
		},
		{
			name: "str	q12, [x28, xzr, sxtx]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xeb, 0xbf, 0x3c}),
				address:          0,
			},
			want: "str	q12, [x28, xzr, sxtx]",
			wantErr: false,
		},
		{
			name: "str	q13, [x27, x5, sxtx #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6d, 0xfb, 0xa5, 0x3c}),
				address:          0,
			},
			want: "str	q13, [x27, x5, sxtx #4]",
			wantErr: false,
		},
		{
			name: "ldr	q14, [x26, w6, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4e, 0x4b, 0xe6, 0x3c}),
				address:          0,
			},
			want: "ldr	q14, [x26, w6, uxtw]",
			wantErr: false,
		},
		{
			name: "ldr	q15, [x25, w7, uxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2f, 0x4b, 0xe7, 0x3c}),
				address:          0,
			},
			want: "ldr	q15, [x25, w7, uxtw]",
			wantErr: false,
		},
		{
			name: "ldr	q16, [x24, w8, uxtw #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x5b, 0xe8, 0x3c}),
				address:          0,
			},
			want: "ldr	q16, [x24, w8, uxtw #4]",
			wantErr: false,
		},
		{
			name: "ldr	q17, [x23, w9, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xca, 0xe9, 0x3c}),
				address:          0,
			},
			want: "ldr	q17, [x23, w9, sxtw]",
			wantErr: false,
		},
		{
			name: "str	q18, [x22, w10, sxtw]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd2, 0xca, 0xaa, 0x3c}),
				address:          0,
			},
			want: "str	q18, [x22, w10, sxtw]",
			wantErr: false,
		},
		{
			name: "ldr	q19, [x21, wzr, sxtw #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb3, 0xda, 0xff, 0x3c}),
				address:          0,
			},
			want: "ldr	q19, [x21, wzr, sxtw #4]",
			wantErr: false,
		},
		{
			name: "strb	w9, [x2], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xf4, 0x0f, 0x38}),
				address:          0,
			},
			want: "strb	w9, [x2], #255",
			wantErr: false,
		},
		{
			name: "strb	w10, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x14, 0x00, 0x38}),
				address:          0,
			},
			want: "strb	w10, [x3], #1",
			wantErr: false,
		},
		{
			name: "strb	w10, [x3], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x04, 0x10, 0x38}),
				address:          0,
			},
			want: "strb	w10, [x3], #-256",
			wantErr: false,
		},
		{
			name: "strh	w9, [x2], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xf4, 0x0f, 0x78}),
				address:          0,
			},
			want: "strh	w9, [x2], #255",
			wantErr: false,
		},
		{
			name: "strh	w9, [x2], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x14, 0x00, 0x78}),
				address:          0,
			},
			want: "strh	w9, [x2], #1",
			wantErr: false,
		},
		{
			name: "strh	w10, [x3], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x04, 0x10, 0x78}),
				address:          0,
			},
			want: "strh	w10, [x3], #-256",
			wantErr: false,
		},
		{
			name: "str	w19, [sp], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xf7, 0x0f, 0xb8}),
				address:          0,
			},
			want: "str	w19, [sp], #255",
			wantErr: false,
		},
		{
			name: "str	w20, [x30], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd4, 0x17, 0x00, 0xb8}),
				address:          0,
			},
			want: "str	w20, [x30], #1",
			wantErr: false,
		},
		{
			name: "str	w21, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x95, 0x05, 0x10, 0xb8}),
				address:          0,
			},
			want: "str	w21, [x12], #-256",
			wantErr: false,
		},
		{
			name: "str	xzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0x0f, 0xf8}),
				address:          0,
			},
			want: "str	xzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "str	x2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0x00, 0xf8}),
				address:          0,
			},
			want: "str	x2, [x3], #1",
			wantErr: false,
		},
		{
			name: "str	x19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0x10, 0xf8}),
				address:          0,
			},
			want: "str	x19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldrb	w9, [x2], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xf4, 0x4f, 0x38}),
				address:          0,
			},
			want: "ldrb	w9, [x2], #255",
			wantErr: false,
		},
		{
			name: "ldrb	w10, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x14, 0x40, 0x38}),
				address:          0,
			},
			want: "ldrb	w10, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldrb	w10, [x3], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x04, 0x50, 0x38}),
				address:          0,
			},
			want: "ldrb	w10, [x3], #-256",
			wantErr: false,
		},
		{
			name: "ldrh	w9, [x2], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xf4, 0x4f, 0x78}),
				address:          0,
			},
			want: "ldrh	w9, [x2], #255",
			wantErr: false,
		},
		{
			name: "ldrh	w9, [x2], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x14, 0x40, 0x78}),
				address:          0,
			},
			want: "ldrh	w9, [x2], #1",
			wantErr: false,
		},
		{
			name: "ldrh	w10, [x3], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x04, 0x50, 0x78}),
				address:          0,
			},
			want: "ldrh	w10, [x3], #-256",
			wantErr: false,
		},
		{
			name: "ldr	w19, [sp], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xf7, 0x4f, 0xb8}),
				address:          0,
			},
			want: "ldr	w19, [sp], #255",
			wantErr: false,
		},
		{
			name: "ldr	w20, [x30], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd4, 0x17, 0x40, 0xb8}),
				address:          0,
			},
			want: "ldr	w20, [x30], #1",
			wantErr: false,
		},
		{
			name: "ldr	w21, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x95, 0x05, 0x50, 0xb8}),
				address:          0,
			},
			want: "ldr	w21, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldr	xzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0x4f, 0xf8}),
				address:          0,
			},
			want: "ldr	xzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "ldr	x2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0x40, 0xf8}),
				address:          0,
			},
			want: "ldr	x2, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldr	x19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0x50, 0xf8}),
				address:          0,
			},
			want: "ldr	x19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldrsb	xzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0x8f, 0x38}),
				address:          0,
			},
			want: "ldrsb	xzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "ldrsb	x2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0x80, 0x38}),
				address:          0,
			},
			want: "ldrsb	x2, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldrsb	x19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0x90, 0x38}),
				address:          0,
			},
			want: "ldrsb	x19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldrsh	xzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0x8f, 0x78}),
				address:          0,
			},
			want: "ldrsh	xzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "ldrsh	x2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0x80, 0x78}),
				address:          0,
			},
			want: "ldrsh	x2, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldrsh	x19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0x90, 0x78}),
				address:          0,
			},
			want: "ldrsh	x19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldrsw	xzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0x8f, 0xb8}),
				address:          0,
			},
			want: "ldrsw	xzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "ldrsw	x2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0x80, 0xb8}),
				address:          0,
			},
			want: "ldrsw	x2, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldrsw	x19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0x90, 0xb8}),
				address:          0,
			},
			want: "ldrsw	x19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldrsb	wzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0xcf, 0x38}),
				address:          0,
			},
			want: "ldrsb	wzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "ldrsb	w2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0xc0, 0x38}),
				address:          0,
			},
			want: "ldrsb	w2, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldrsb	w19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0xd0, 0x38}),
				address:          0,
			},
			want: "ldrsb	w19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "ldrsh	wzr, [x9], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xf5, 0xcf, 0x78}),
				address:          0,
			},
			want: "ldrsh	wzr, [x9], #255",
			wantErr: false,
		},
		{
			name: "ldrsh	w2, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x14, 0xc0, 0x78}),
				address:          0,
			},
			want: "ldrsh	w2, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldrsh	w19, [x12], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x05, 0xd0, 0x78}),
				address:          0,
			},
			want: "ldrsh	w19, [x12], #-256",
			wantErr: false,
		},
		{
			name: "str	b0, [x0], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xf4, 0x0f, 0x3c}),
				address:          0,
			},
			want: "str	b0, [x0], #255",
			wantErr: false,
		},
		{
			name: "str	b3, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x63, 0x14, 0x00, 0x3c}),
				address:          0,
			},
			want: "str	b3, [x3], #1",
			wantErr: false,
		},
		{
			name: "str	b5, [sp], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x07, 0x10, 0x3c}),
				address:          0,
			},
			want: "str	b5, [sp], #-256",
			wantErr: false,
		},
		{
			name: "str	h10, [x10], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0xf5, 0x0f, 0x7c}),
				address:          0,
			},
			want: "str	h10, [x10], #255",
			wantErr: false,
		},
		{
			name: "str	h13, [x23], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x16, 0x00, 0x7c}),
				address:          0,
			},
			want: "str	h13, [x23], #1",
			wantErr: false,
		},
		{
			name: "str	h15, [sp], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0x07, 0x10, 0x7c}),
				address:          0,
			},
			want: "str	h15, [sp], #-256",
			wantErr: false,
		},
		{
			name: "str	s20, [x20], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xf6, 0x0f, 0xbc}),
				address:          0,
			},
			want: "str	s20, [x20], #255",
			wantErr: false,
		},
		{
			name: "str	s23, [x23], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x16, 0x00, 0xbc}),
				address:          0,
			},
			want: "str	s23, [x23], #1",
			wantErr: false,
		},
		{
			name: "str	s25, [x0], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x04, 0x10, 0xbc}),
				address:          0,
			},
			want: "str	s25, [x0], #-256",
			wantErr: false,
		},
		{
			name: "str	d20, [x20], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xf6, 0x0f, 0xfc}),
				address:          0,
			},
			want: "str	d20, [x20], #255",
			wantErr: false,
		},
		{
			name: "str	d23, [x23], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x16, 0x00, 0xfc}),
				address:          0,
			},
			want: "str	d23, [x23], #1",
			wantErr: false,
		},
		{
			name: "str	d25, [x0], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x04, 0x10, 0xfc}),
				address:          0,
			},
			want: "str	d25, [x0], #-256",
			wantErr: false,
		},
		{
			name: "ldr	b0, [x0], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xf4, 0x4f, 0x3c}),
				address:          0,
			},
			want: "ldr	b0, [x0], #255",
			wantErr: false,
		},
		{
			name: "ldr	b3, [x3], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x63, 0x14, 0x40, 0x3c}),
				address:          0,
			},
			want: "ldr	b3, [x3], #1",
			wantErr: false,
		},
		{
			name: "ldr	b5, [sp], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x07, 0x50, 0x3c}),
				address:          0,
			},
			want: "ldr	b5, [sp], #-256",
			wantErr: false,
		},
		{
			name: "ldr	h10, [x10], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0xf5, 0x4f, 0x7c}),
				address:          0,
			},
			want: "ldr	h10, [x10], #255",
			wantErr: false,
		},
		{
			name: "ldr	h13, [x23], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x16, 0x40, 0x7c}),
				address:          0,
			},
			want: "ldr	h13, [x23], #1",
			wantErr: false,
		},
		{
			name: "ldr	h15, [sp], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0x07, 0x50, 0x7c}),
				address:          0,
			},
			want: "ldr	h15, [sp], #-256",
			wantErr: false,
		},
		{
			name: "ldr	s20, [x20], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xf6, 0x4f, 0xbc}),
				address:          0,
			},
			want: "ldr	s20, [x20], #255",
			wantErr: false,
		},
		{
			name: "ldr	s23, [x23], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x16, 0x40, 0xbc}),
				address:          0,
			},
			want: "ldr	s23, [x23], #1",
			wantErr: false,
		},
		{
			name: "ldr	s25, [x0], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x04, 0x50, 0xbc}),
				address:          0,
			},
			want: "ldr	s25, [x0], #-256",
			wantErr: false,
		},
		{
			name: "ldr	d20, [x20], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xf6, 0x4f, 0xfc}),
				address:          0,
			},
			want: "ldr	d20, [x20], #255",
			wantErr: false,
		},
		{
			name: "ldr	d23, [x23], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x16, 0x40, 0xfc}),
				address:          0,
			},
			want: "ldr	d23, [x23], #1",
			wantErr: false,
		},
		{
			name: "ldr	d25, [x0], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x04, 0x50, 0xfc}),
				address:          0,
			},
			want: "ldr	d25, [x0], #-256",
			wantErr: false,
		},
		{
			name: "ldr	q20, [x1], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x34, 0xf4, 0xcf, 0x3c}),
				address:          0,
			},
			want: "ldr	q20, [x1], #255",
			wantErr: false,
		},
		{
			name: "ldr	q23, [x9], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x37, 0x15, 0xc0, 0x3c}),
				address:          0,
			},
			want: "ldr	q23, [x9], #1",
			wantErr: false,
		},
		{
			name: "ldr	q25, [x20], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x99, 0x06, 0xd0, 0x3c}),
				address:          0,
			},
			want: "ldr	q25, [x20], #-256",
			wantErr: false,
		},
		{
			name: "str	q10, [x1], #255",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2a, 0xf4, 0x8f, 0x3c}),
				address:          0,
			},
			want: "str	q10, [x1], #255",
			wantErr: false,
		},
		{
			name: "str	q22, [sp], #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x17, 0x80, 0x3c}),
				address:          0,
			},
			want: "str	q22, [sp], #1",
			wantErr: false,
		},
		{
			name: "str	q21, [x20], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x95, 0x06, 0x90, 0x3c}),
				address:          0,
			},
			want: "str	q21, [x20], #-256",
			wantErr: false,
		},
		{
			name: "ldr	x3, [x4, #0]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x83, 0x0c, 0x40, 0xf8}),
				address:          0,
			},
			want: "ldr	x3, [x4, #0]!",
			wantErr: false,
		},
		{
			name: "ldr	xzr, [sp, #0]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x0f, 0x40, 0xf8}),
				address:          0,
			},
			want: "ldr	xzr, [sp, #0]!",
			wantErr: false,
		},
		{
			name: "strb	w9, [x2, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfc, 0x0f, 0x38}),
				address:          0,
			},
			want: "strb	w9, [x2, #255]!",
			wantErr: false,
		},
		{
			name: "strb	w10, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x1c, 0x00, 0x38}),
				address:          0,
			},
			want: "strb	w10, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "strb	w10, [x3, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x0c, 0x10, 0x38}),
				address:          0,
			},
			want: "strb	w10, [x3, #-256]!",
			wantErr: false,
		},
		{
			name: "strh	w9, [x2, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfc, 0x0f, 0x78}),
				address:          0,
			},
			want: "strh	w9, [x2, #255]!",
			wantErr: false,
		},
		{
			name: "strh	w9, [x2, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x1c, 0x00, 0x78}),
				address:          0,
			},
			want: "strh	w9, [x2, #1]!",
			wantErr: false,
		},
		{
			name: "strh	w10, [x3, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x0c, 0x10, 0x78}),
				address:          0,
			},
			want: "strh	w10, [x3, #-256]!",
			wantErr: false,
		},
		{
			name: "str	w19, [sp, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xff, 0x0f, 0xb8}),
				address:          0,
			},
			want: "str	w19, [sp, #255]!",
			wantErr: false,
		},
		{
			name: "str	w20, [x30, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd4, 0x1f, 0x00, 0xb8}),
				address:          0,
			},
			want: "str	w20, [x30, #1]!",
			wantErr: false,
		},
		{
			name: "str	w21, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x95, 0x0d, 0x10, 0xb8}),
				address:          0,
			},
			want: "str	w21, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "str	xzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0x0f, 0xf8}),
				address:          0,
			},
			want: "str	xzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "str	x2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0x00, 0xf8}),
				address:          0,
			},
			want: "str	x2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "str	x19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0x10, 0xf8}),
				address:          0,
			},
			want: "str	x19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrb	w9, [x2, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfc, 0x4f, 0x38}),
				address:          0,
			},
			want: "ldrb	w9, [x2, #255]!",
			wantErr: false,
		},
		{
			name: "ldrb	w10, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x1c, 0x40, 0x38}),
				address:          0,
			},
			want: "ldrb	w10, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldrb	w10, [x3, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x0c, 0x50, 0x38}),
				address:          0,
			},
			want: "ldrb	w10, [x3, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrh	w9, [x2, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xfc, 0x4f, 0x78}),
				address:          0,
			},
			want: "ldrh	w9, [x2, #255]!",
			wantErr: false,
		},
		{
			name: "ldrh	w9, [x2, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x1c, 0x40, 0x78}),
				address:          0,
			},
			want: "ldrh	w9, [x2, #1]!",
			wantErr: false,
		},
		{
			name: "ldrh	w10, [x3, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6a, 0x0c, 0x50, 0x78}),
				address:          0,
			},
			want: "ldrh	w10, [x3, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	w19, [sp, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf3, 0xff, 0x4f, 0xb8}),
				address:          0,
			},
			want: "ldr	w19, [sp, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	w20, [x30, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd4, 0x1f, 0x40, 0xb8}),
				address:          0,
			},
			want: "ldr	w20, [x30, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	w21, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x95, 0x0d, 0x50, 0xb8}),
				address:          0,
			},
			want: "ldr	w21, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	xzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0x4f, 0xf8}),
				address:          0,
			},
			want: "ldr	xzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	x2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0x40, 0xf8}),
				address:          0,
			},
			want: "ldr	x2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	x19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0x50, 0xf8}),
				address:          0,
			},
			want: "ldr	x19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrsb	xzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0x8f, 0x38}),
				address:          0,
			},
			want: "ldrsb	xzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "ldrsb	x2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0x80, 0x38}),
				address:          0,
			},
			want: "ldrsb	x2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldrsb	x19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0x90, 0x38}),
				address:          0,
			},
			want: "ldrsb	x19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrsh	xzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0x8f, 0x78}),
				address:          0,
			},
			want: "ldrsh	xzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "ldrsh	x2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0x80, 0x78}),
				address:          0,
			},
			want: "ldrsh	x2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldrsh	x19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0x90, 0x78}),
				address:          0,
			},
			want: "ldrsh	x19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrsw	xzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0x8f, 0xb8}),
				address:          0,
			},
			want: "ldrsw	xzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "ldrsw	x2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0x80, 0xb8}),
				address:          0,
			},
			want: "ldrsw	x2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldrsw	x19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0x90, 0xb8}),
				address:          0,
			},
			want: "ldrsw	x19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrsb	wzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0xcf, 0x38}),
				address:          0,
			},
			want: "ldrsb	wzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "ldrsb	w2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0xc0, 0x38}),
				address:          0,
			},
			want: "ldrsb	w2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldrsb	w19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0xd0, 0x38}),
				address:          0,
			},
			want: "ldrsb	w19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "ldrsh	wzr, [x9, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0xfd, 0xcf, 0x78}),
				address:          0,
			},
			want: "ldrsh	wzr, [x9, #255]!",
			wantErr: false,
		},
		{
			name: "ldrsh	w2, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x1c, 0xc0, 0x78}),
				address:          0,
			},
			want: "ldrsh	w2, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldrsh	w19, [x12, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x0d, 0xd0, 0x78}),
				address:          0,
			},
			want: "ldrsh	w19, [x12, #-256]!",
			wantErr: false,
		},
		{
			name: "str	b0, [x0, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xfc, 0x0f, 0x3c}),
				address:          0,
			},
			want: "str	b0, [x0, #255]!",
			wantErr: false,
		},
		{
			name: "str	b3, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x63, 0x1c, 0x00, 0x3c}),
				address:          0,
			},
			want: "str	b3, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "str	b5, [sp, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x0f, 0x10, 0x3c}),
				address:          0,
			},
			want: "str	b5, [sp, #-256]!",
			wantErr: false,
		},
		{
			name: "str	h10, [x10, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0xfd, 0x0f, 0x7c}),
				address:          0,
			},
			want: "str	h10, [x10, #255]!",
			wantErr: false,
		},
		{
			name: "str	h13, [x23, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x1e, 0x00, 0x7c}),
				address:          0,
			},
			want: "str	h13, [x23, #1]!",
			wantErr: false,
		},
		{
			name: "str	h15, [sp, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0x0f, 0x10, 0x7c}),
				address:          0,
			},
			want: "str	h15, [sp, #-256]!",
			wantErr: false,
		},
		{
			name: "str	s20, [x20, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xfe, 0x0f, 0xbc}),
				address:          0,
			},
			want: "str	s20, [x20, #255]!",
			wantErr: false,
		},
		{
			name: "str	s23, [x23, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x1e, 0x00, 0xbc}),
				address:          0,
			},
			want: "str	s23, [x23, #1]!",
			wantErr: false,
		},
		{
			name: "str	s25, [x0, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x0c, 0x10, 0xbc}),
				address:          0,
			},
			want: "str	s25, [x0, #-256]!",
			wantErr: false,
		},
		{
			name: "str	d20, [x20, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xfe, 0x0f, 0xfc}),
				address:          0,
			},
			want: "str	d20, [x20, #255]!",
			wantErr: false,
		},
		{
			name: "str	d23, [x23, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x1e, 0x00, 0xfc}),
				address:          0,
			},
			want: "str	d23, [x23, #1]!",
			wantErr: false,
		},
		{
			name: "str	d25, [x0, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x0c, 0x10, 0xfc}),
				address:          0,
			},
			want: "str	d25, [x0, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	b0, [x0, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0xfc, 0x4f, 0x3c}),
				address:          0,
			},
			want: "ldr	b0, [x0, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	b3, [x3, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x63, 0x1c, 0x40, 0x3c}),
				address:          0,
			},
			want: "ldr	b3, [x3, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	b5, [sp, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x0f, 0x50, 0x3c}),
				address:          0,
			},
			want: "ldr	b5, [sp, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	h10, [x10, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0xfd, 0x4f, 0x7c}),
				address:          0,
			},
			want: "ldr	h10, [x10, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	h13, [x23, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0x1e, 0x40, 0x7c}),
				address:          0,
			},
			want: "ldr	h13, [x23, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	h15, [sp, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xef, 0x0f, 0x50, 0x7c}),
				address:          0,
			},
			want: "ldr	h15, [sp, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	s20, [x20, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xfe, 0x4f, 0xbc}),
				address:          0,
			},
			want: "ldr	s20, [x20, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	s23, [x23, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x1e, 0x40, 0xbc}),
				address:          0,
			},
			want: "ldr	s23, [x23, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	s25, [x0, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x0c, 0x50, 0xbc}),
				address:          0,
			},
			want: "ldr	s25, [x0, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	d20, [x20, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0xfe, 0x4f, 0xfc}),
				address:          0,
			},
			want: "ldr	d20, [x20, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	d23, [x23, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf7, 0x1e, 0x40, 0xfc}),
				address:          0,
			},
			want: "ldr	d23, [x23, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	d25, [x0, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x19, 0x0c, 0x50, 0xfc}),
				address:          0,
			},
			want: "ldr	d25, [x0, #-256]!",
			wantErr: false,
		},
		{
			name: "ldr	q20, [x1, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x34, 0xfc, 0xcf, 0x3c}),
				address:          0,
			},
			want: "ldr	q20, [x1, #255]!",
			wantErr: false,
		},
		{
			name: "ldr	q23, [x9, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x37, 0x1d, 0xc0, 0x3c}),
				address:          0,
			},
			want: "ldr	q23, [x9, #1]!",
			wantErr: false,
		},
		{
			name: "ldr	q25, [x20, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x99, 0x0e, 0xd0, 0x3c}),
				address:          0,
			},
			want: "ldr	q25, [x20, #-256]!",
			wantErr: false,
		},
		{
			name: "str	q10, [x1, #255]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2a, 0xfc, 0x8f, 0x3c}),
				address:          0,
			},
			want: "str	q10, [x1, #255]!",
			wantErr: false,
		},
		{
			name: "str	q22, [sp, #1]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf6, 0x1f, 0x80, 0x3c}),
				address:          0,
			},
			want: "str	q22, [sp, #1]!",
			wantErr: false,
		},
		{
			name: "str	q21, [x20, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x95, 0x0e, 0x90, 0x3c}),
				address:          0,
			},
			want: "str	q21, [x20, #-256]!",
			wantErr: false,
		},
		{
			name: "sttrb	w9, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0b, 0x00, 0x38}),
				address:          0,
			},
			want: "sttrb	w9, [sp]",
			wantErr: false,
		},
		{
			name: "sttrh	wzr, [x12, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xf9, 0x0f, 0x78}),
				address:          0,
			},
			want: "sttrh	wzr, [x12, #255]",
			wantErr: false,
		},
		{
			name: "sttr	w16, [x0, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x10, 0x08, 0x10, 0xb8}),
				address:          0,
			},
			want: "sttr	w16, [x0, #-256]",
			wantErr: false,
		},
		{
			name: "sttr	x28, [x14, #1]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdc, 0x19, 0x00, 0xf8}),
				address:          0,
			},
			want: "sttr	x28, [x14, #1]",
			wantErr: false,
		},
		{
			name: "ldtrb	w1, [x20, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x81, 0xfa, 0x4f, 0x38}),
				address:          0,
			},
			want: "ldtrb	w1, [x20, #255]",
			wantErr: false,
		},
		{
			name: "ldtrh	w20, [x1, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x34, 0xf8, 0x4f, 0x78}),
				address:          0,
			},
			want: "ldtrh	w20, [x1, #255]",
			wantErr: false,
		},
		{
			name: "ldtr	w12, [sp, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xfb, 0x4f, 0xb8}),
				address:          0,
			},
			want: "ldtr	w12, [sp, #255]",
			wantErr: false,
		},
		{
			name: "ldtr	xzr, [x12, #255]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0xf9, 0x4f, 0xf8}),
				address:          0,
			},
			want: "ldtr	xzr, [x12, #255]",
			wantErr: false,
		},
		{
			name: "ldtrsb	x9, [x7, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x08, 0x90, 0x38}),
				address:          0,
			},
			want: "ldtrsb	x9, [x7, #-256]",
			wantErr: false,
		},
		{
			name: "ldtrsh	x17, [x19, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x71, 0x0a, 0x90, 0x78}),
				address:          0,
			},
			want: "ldtrsh	x17, [x19, #-256]",
			wantErr: false,
		},
		{
			name: "ldtrsw	x20, [x15, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x09, 0x90, 0xb8}),
				address:          0,
			},
			want: "ldtrsw	x20, [x15, #-256]",
			wantErr: false,
		},
		{
			name: "ldtrsb	w19, [x1, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x08, 0xd0, 0x38}),
				address:          0,
			},
			want: "ldtrsb	w19, [x1, #-256]",
			wantErr: false,
		},
		{
			name: "ldtrsh	w15, [x21, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaf, 0x0a, 0xd0, 0x78}),
				address:          0,
			},
			want: "ldtrsh	w15, [x21, #-256]",
			wantErr: false,
		},
		{
			name: "ldp	w3, w5, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x40, 0x29}),
				address:          0,
			},
			want: "ldp	w3, w5, [sp]",
			wantErr: false,
		},
		{
			name: "stp	wzr, w9, [sp, #252]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa7, 0x1f, 0x29}),
				address:          0,
			},
			want: "stp	wzr, w9, [sp, #252]",
			wantErr: false,
		},
		{
			name: "ldp	w2, wzr, [sp, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x7f, 0x60, 0x29}),
				address:          0,
			},
			want: "ldp	w2, wzr, [sp, #-256]",
			wantErr: false,
		},
		{
			name: "ldp	w9, w10, [sp, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0x40, 0x29}),
				address:          0,
			},
			want: "ldp	w9, w10, [sp, #4]",
			wantErr: false,
		},
		{
			name: "ldpsw	x9, x10, [sp, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0x40, 0x69}),
				address:          0,
			},
			want: "ldpsw	x9, x10, [sp, #4]",
			wantErr: false,
		},
		{
			name: "ldpsw	x9, x10, [x2, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x28, 0x60, 0x69}),
				address:          0,
			},
			want: "ldpsw	x9, x10, [x2, #-256]",
			wantErr: false,
		},
		{
			name: "ldpsw	x20, x30, [sp, #252]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xfb, 0x5f, 0x69}),
				address:          0,
			},
			want: "ldpsw	x20, x30, [sp, #252]",
			wantErr: false,
		},
		{
			name: "ldp	x21, x29, [x2, #504]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf4, 0x5f, 0xa9}),
				address:          0,
			},
			want: "ldp	x21, x29, [x2, #504]",
			wantErr: false,
		},
		{
			name: "ldp	x22, x23, [x3, #-512]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x76, 0x5c, 0x60, 0xa9}),
				address:          0,
			},
			want: "ldp	x22, x23, [x3, #-512]",
			wantErr: false,
		},
		{
			name: "ldp	x24, x25, [x4, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x98, 0xe4, 0x40, 0xa9}),
				address:          0,
			},
			want: "ldp	x24, x25, [x4, #8]",
			wantErr: false,
		},
		{
			name: "ldp	s29, s28, [sp, #252]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0xf3, 0x5f, 0x2d}),
				address:          0,
			},
			want: "ldp	s29, s28, [sp, #252]",
			wantErr: false,
		},
		{
			name: "stp	s27, s26, [sp, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfb, 0x6b, 0x20, 0x2d}),
				address:          0,
			},
			want: "stp	s27, s26, [sp, #-256]",
			wantErr: false,
		},
		{
			name: "ldp	s1, s2, [x3, #44]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x88, 0x45, 0x2d}),
				address:          0,
			},
			want: "ldp	s1, s2, [x3, #44]",
			wantErr: false,
		},
		{
			name: "stp	d3, d5, [x9, #504]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x95, 0x1f, 0x6d}),
				address:          0,
			},
			want: "stp	d3, d5, [x9, #504]",
			wantErr: false,
		},
		{
			name: "stp	d7, d11, [x10, #-512]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x47, 0x2d, 0x20, 0x6d}),
				address:          0,
			},
			want: "stp	d7, d11, [x10, #-512]",
			wantErr: false,
		},
		{
			name: "ldp	d2, d3, [x30, #-8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc2, 0x8f, 0x7f, 0x6d}),
				address:          0,
			},
			want: "ldp	d2, d3, [x30, #-8]",
			wantErr: false,
		},
		{
			name: "stp	q3, q5, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x00, 0xad}),
				address:          0,
			},
			want: "stp	q3, q5, [sp]",
			wantErr: false,
		},
		{
			name: "stp	q17, q19, [sp, #1008]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xcf, 0x1f, 0xad}),
				address:          0,
			},
			want: "stp	q17, q19, [sp, #1008]",
			wantErr: false,
		},
		{
			name: "ldp	q23, q29, [x1, #-1024]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x37, 0x74, 0x60, 0xad}),
				address:          0,
			},
			want: "ldp	q23, q29, [x1, #-1024]",
			wantErr: false,
		},
		{
			name: "ldp	w3, w5, [sp], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0xc0, 0x28}),
				address:          0,
			},
			want: "ldp	w3, w5, [sp], #0",
			wantErr: false,
		},
		{
			name: "stp	wzr, w9, [sp], #252",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa7, 0x9f, 0x28}),
				address:          0,
			},
			want: "stp	wzr, w9, [sp], #252",
			wantErr: false,
		},
		{
			name: "ldp	w2, wzr, [sp], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x7f, 0xe0, 0x28}),
				address:          0,
			},
			want: "ldp	w2, wzr, [sp], #-256",
			wantErr: false,
		},
		{
			name: "ldp	w9, w10, [sp], #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0xc0, 0x28}),
				address:          0,
			},
			want: "ldp	w9, w10, [sp], #4",
			wantErr: false,
		},
		{
			name: "ldpsw	x9, x10, [sp], #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0xc0, 0x68}),
				address:          0,
			},
			want: "ldpsw	x9, x10, [sp], #4",
			wantErr: false,
		},
		{
			name: "ldpsw	x9, x10, [x2], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x28, 0xe0, 0x68}),
				address:          0,
			},
			want: "ldpsw	x9, x10, [x2], #-256",
			wantErr: false,
		},
		{
			name: "ldpsw	x20, x30, [sp], #252",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xfb, 0xdf, 0x68}),
				address:          0,
			},
			want: "ldpsw	x20, x30, [sp], #252",
			wantErr: false,
		},
		{
			name: "ldp	x21, x29, [x2], #504",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf4, 0xdf, 0xa8}),
				address:          0,
			},
			want: "ldp	x21, x29, [x2], #504",
			wantErr: false,
		},
		{
			name: "ldp	x22, x23, [x3], #-512",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x76, 0x5c, 0xe0, 0xa8}),
				address:          0,
			},
			want: "ldp	x22, x23, [x3], #-512",
			wantErr: false,
		},
		{
			name: "ldp	x24, x25, [x4], #8",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x98, 0xe4, 0xc0, 0xa8}),
				address:          0,
			},
			want: "ldp	x24, x25, [x4], #8",
			wantErr: false,
		},
		{
			name: "ldp	s29, s28, [sp], #252",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0xf3, 0xdf, 0x2c}),
				address:          0,
			},
			want: "ldp	s29, s28, [sp], #252",
			wantErr: false,
		},
		{
			name: "stp	s27, s26, [sp], #-256",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfb, 0x6b, 0xa0, 0x2c}),
				address:          0,
			},
			want: "stp	s27, s26, [sp], #-256",
			wantErr: false,
		},
		{
			name: "ldp	s1, s2, [x3], #44",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x88, 0xc5, 0x2c}),
				address:          0,
			},
			want: "ldp	s1, s2, [x3], #44",
			wantErr: false,
		},
		{
			name: "stp	d3, d5, [x9], #504",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x95, 0x9f, 0x6c}),
				address:          0,
			},
			want: "stp	d3, d5, [x9], #504",
			wantErr: false,
		},
		{
			name: "stp	d7, d11, [x10], #-512",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x47, 0x2d, 0xa0, 0x6c}),
				address:          0,
			},
			want: "stp	d7, d11, [x10], #-512",
			wantErr: false,
		},
		{
			name: "ldp	d2, d3, [x30], #-8",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc2, 0x8f, 0xff, 0x6c}),
				address:          0,
			},
			want: "ldp	d2, d3, [x30], #-8",
			wantErr: false,
		},
		{
			name: "stp	q3, q5, [sp], #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x80, 0xac}),
				address:          0,
			},
			want: "stp	q3, q5, [sp], #0",
			wantErr: false,
		},
		{
			name: "stp	q17, q19, [sp], #1008",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xcf, 0x9f, 0xac}),
				address:          0,
			},
			want: "stp	q17, q19, [sp], #1008",
			wantErr: false,
		},
		{
			name: "ldp	q23, q29, [x1], #-1024",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x37, 0x74, 0xe0, 0xac}),
				address:          0,
			},
			want: "ldp	q23, q29, [x1], #-1024",
			wantErr: false,
		},
		{
			name: "ldp	w3, w5, [sp, #0]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0xc0, 0x29}),
				address:          0,
			},
			want: "ldp	w3, w5, [sp, #0]!",
			wantErr: false,
		},
		{
			name: "stp	wzr, w9, [sp, #252]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa7, 0x9f, 0x29}),
				address:          0,
			},
			want: "stp	wzr, w9, [sp, #252]!",
			wantErr: false,
		},
		{
			name: "ldp	w2, wzr, [sp, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x7f, 0xe0, 0x29}),
				address:          0,
			},
			want: "ldp	w2, wzr, [sp, #-256]!",
			wantErr: false,
		},
		{
			name: "ldp	w9, w10, [sp, #4]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0xc0, 0x29}),
				address:          0,
			},
			want: "ldp	w9, w10, [sp, #4]!",
			wantErr: false,
		},
		{
			name: "ldpsw	x9, x10, [sp, #4]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0xc0, 0x69}),
				address:          0,
			},
			want: "ldpsw	x9, x10, [sp, #4]!",
			wantErr: false,
		},
		{
			name: "ldpsw	x9, x10, [x2, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x28, 0xe0, 0x69}),
				address:          0,
			},
			want: "ldpsw	x9, x10, [x2, #-256]!",
			wantErr: false,
		},
		{
			name: "ldpsw	x20, x30, [sp, #252]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0xfb, 0xdf, 0x69}),
				address:          0,
			},
			want: "ldpsw	x20, x30, [sp, #252]!",
			wantErr: false,
		},
		{
			name: "ldp	x21, x29, [x2, #504]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf4, 0xdf, 0xa9}),
				address:          0,
			},
			want: "ldp	x21, x29, [x2, #504]!",
			wantErr: false,
		},
		{
			name: "ldp	x22, x23, [x3, #-512]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x76, 0x5c, 0xe0, 0xa9}),
				address:          0,
			},
			want: "ldp	x22, x23, [x3, #-512]!",
			wantErr: false,
		},
		{
			name: "ldp	x24, x25, [x4, #8]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x98, 0xe4, 0xc0, 0xa9}),
				address:          0,
			},
			want: "ldp	x24, x25, [x4, #8]!",
			wantErr: false,
		},
		{
			name: "ldp	s29, s28, [sp, #252]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0xf3, 0xdf, 0x2d}),
				address:          0,
			},
			want: "ldp	s29, s28, [sp, #252]!",
			wantErr: false,
		},
		{
			name: "stp	s27, s26, [sp, #-256]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfb, 0x6b, 0xa0, 0x2d}),
				address:          0,
			},
			want: "stp	s27, s26, [sp, #-256]!",
			wantErr: false,
		},
		{
			name: "ldp	s1, s2, [x3, #44]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x88, 0xc5, 0x2d}),
				address:          0,
			},
			want: "ldp	s1, s2, [x3, #44]!",
			wantErr: false,
		},
		{
			name: "stp	d3, d5, [x9, #504]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x95, 0x9f, 0x6d}),
				address:          0,
			},
			want: "stp	d3, d5, [x9, #504]!",
			wantErr: false,
		},
		{
			name: "stp	d7, d11, [x10, #-512]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x47, 0x2d, 0xa0, 0x6d}),
				address:          0,
			},
			want: "stp	d7, d11, [x10, #-512]!",
			wantErr: false,
		},
		{
			name: "ldp	d2, d3, [x30, #-8]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc2, 0x8f, 0xff, 0x6d}),
				address:          0,
			},
			want: "ldp	d2, d3, [x30, #-8]!",
			wantErr: false,
		},
		{
			name: "stp	q3, q5, [sp, #0]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x80, 0xad}),
				address:          0,
			},
			want: "stp	q3, q5, [sp, #0]!",
			wantErr: false,
		},
		{
			name: "stp	q17, q19, [sp, #1008]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xcf, 0x9f, 0xad}),
				address:          0,
			},
			want: "stp	q17, q19, [sp, #1008]!",
			wantErr: false,
		},
		{
			name: "ldp	q23, q29, [x1, #-1024]!",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x37, 0x74, 0xe0, 0xad}),
				address:          0,
			},
			want: "ldp	q23, q29, [x1, #-1024]!",
			wantErr: false,
		},
		{
			name: "ldnp	w3, w5, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x40, 0x28}),
				address:          0,
			},
			want: "ldnp	w3, w5, [sp]",
			wantErr: false,
		},
		{
			name: "stnp	wzr, w9, [sp, #252]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xa7, 0x1f, 0x28}),
				address:          0,
			},
			want: "stnp	wzr, w9, [sp, #252]",
			wantErr: false,
		},
		{
			name: "ldnp	w2, wzr, [sp, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x7f, 0x60, 0x28}),
				address:          0,
			},
			want: "ldnp	w2, wzr, [sp, #-256]",
			wantErr: false,
		},
		{
			name: "ldnp	w9, w10, [sp, #4]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xab, 0x40, 0x28}),
				address:          0,
			},
			want: "ldnp	w9, w10, [sp, #4]",
			wantErr: false,
		},
		{
			name: "ldnp	x21, x29, [x2, #504]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x55, 0xf4, 0x5f, 0xa8}),
				address:          0,
			},
			want: "ldnp	x21, x29, [x2, #504]",
			wantErr: false,
		},
		{
			name: "ldnp	x22, x23, [x3, #-512]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x76, 0x5c, 0x60, 0xa8}),
				address:          0,
			},
			want: "ldnp	x22, x23, [x3, #-512]",
			wantErr: false,
		},
		{
			name: "ldnp	x24, x25, [x4, #8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x98, 0xe4, 0x40, 0xa8}),
				address:          0,
			},
			want: "ldnp	x24, x25, [x4, #8]",
			wantErr: false,
		},
		{
			name: "ldnp	s29, s28, [sp, #252]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfd, 0xf3, 0x5f, 0x2c}),
				address:          0,
			},
			want: "ldnp	s29, s28, [sp, #252]",
			wantErr: false,
		},
		{
			name: "stnp	s27, s26, [sp, #-256]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xfb, 0x6b, 0x20, 0x2c}),
				address:          0,
			},
			want: "stnp	s27, s26, [sp, #-256]",
			wantErr: false,
		},
		{
			name: "ldnp	s1, s2, [x3, #44]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x61, 0x88, 0x45, 0x2c}),
				address:          0,
			},
			want: "ldnp	s1, s2, [x3, #44]",
			wantErr: false,
		},
		{
			name: "stnp	d3, d5, [x9, #504]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x95, 0x1f, 0x6c}),
				address:          0,
			},
			want: "stnp	d3, d5, [x9, #504]",
			wantErr: false,
		},
		{
			name: "stnp	d7, d11, [x10, #-512]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x47, 0x2d, 0x20, 0x6c}),
				address:          0,
			},
			want: "stnp	d7, d11, [x10, #-512]",
			wantErr: false,
		},
		{
			name: "ldnp	d2, d3, [x30, #-8]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc2, 0x8f, 0x7f, 0x6c}),
				address:          0,
			},
			want: "ldnp	d2, d3, [x30, #-8]",
			wantErr: false,
		},
		{
			name: "stnp	q3, q5, [sp]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x17, 0x00, 0xac}),
				address:          0,
			},
			want: "stnp	q3, q5, [sp]",
			wantErr: false,
		},
		{
			name: "stnp	q17, q19, [sp, #1008]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf1, 0xcf, 0x1f, 0xac}),
				address:          0,
			},
			want: "stnp	q17, q19, [sp, #1008]",
			wantErr: false,
		},
		{
			name: "ldnp	q23, q29, [x1, #-1024]",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x37, 0x74, 0x60, 0xac}),
				address:          0,
			},
			want: "ldnp	q23, q29, [x1, #-1024]",
			wantErr: false,
		},
		{
			name: "orr	w3, w9, #0xffff0000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x3d, 0x10, 0x32}),
				address:          0,
			},
			want: "orr	w3, w9, #0xffff0000",
			wantErr: false,
		},
		{
			name: "orr	wsp, w10, #0xe00000ff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x29, 0x03, 0x32}),
				address:          0,
			},
			want: "orr	wsp, w10, #0xe00000ff",
			wantErr: false,
		},
		{
			name: "orr	w9, w10, #0x3ff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x25, 0x00, 0x32}),
				address:          0,
			},
			want: "orr	w9, w10, #0x3ff",
			wantErr: false,
		},
		{
			name: "and	w14, w15, #0x80008000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x81, 0x01, 0x12}),
				address:          0,
			},
			want: "and	w14, w15, #0x80008000",
			wantErr: false,
		},
		{
			name: "and	w12, w13, #0xffc3ffc3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xad, 0x0a, 0x12}),
				address:          0,
			},
			want: "and	w12, w13, #0xffc3ffc3",
			wantErr: false,
		},
		{
			name: "and	w11, wzr, #0x30003",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xeb, 0x87, 0x00, 0x12}),
				address:          0,
			},
			want: "and	w11, wzr, #0x30003",
			wantErr: false,
		},
		{
			name: "eor	w3, w6, #0xe0e0e0e0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc3, 0xc8, 0x03, 0x52}),
				address:          0,
			},
			want: "eor	w3, w6, #0xe0e0e0e0",
			wantErr: false,
		},
		{
			name: "eor	wsp, wzr, #0x3030303",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xc7, 0x00, 0x52}),
				address:          0,
			},
			want: "eor	wsp, wzr, #0x3030303",
			wantErr: false,
		},
		{
			name: "eor	w16, w17, #0x81818181",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0xc6, 0x01, 0x52}),
				address:          0,
			},
			want: "eor	w16, w17, #0x81818181",
			wantErr: false,
		},
		{
			name: "tst	w18, #0xcccccccc",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xe6, 0x02, 0x72}),
				address:          0,
			},
			want: "tst	w18, #0xcccccccc",
			wantErr: false,
		},
		{
			name: "ands	w19, w20, #0x33333333",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xe6, 0x00, 0x72}),
				address:          0,
			},
			want: "ands	w19, w20, #0x33333333",
			wantErr: false,
		},
		{
			name: "ands	w21, w22, #0x99999999",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0xe6, 0x01, 0x72}),
				address:          0,
			},
			want: "ands	w21, w22, #0x99999999",
			wantErr: false,
		},
		{
			name: "tst	w3, #0xaaaaaaaa",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xf0, 0x01, 0x72}),
				address:          0,
			},
			want: "tst	w3, #0xaaaaaaaa",
			wantErr: false,
		},
		{
			name: "tst	wzr, #0x55555555",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xf3, 0x00, 0x72}),
				address:          0,
			},
			want: "tst	wzr, #0x55555555",
			wantErr: false,
		},
		{
			name: "eor	x3, x5, #0xffffffffc000000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x84, 0x66, 0xd2}),
				address:          0,
			},
			want: "eor	x3, x5, #0xffffffffc000000",
			wantErr: false,
		},
		{
			name: "and	x9, x10, #0x7fffffffffff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xb9, 0x40, 0x92}),
				address:          0,
			},
			want: "and	x9, x10, #0x7fffffffffff",
			wantErr: false,
		},
		{
			name: "orr	x11, x12, #0x8000000000000fff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8b, 0x31, 0x41, 0xb2}),
				address:          0,
			},
			want: "orr	x11, x12, #0x8000000000000fff",
			wantErr: false,
		},
		{
			name: "orr	x3, x9, #0xffff0000ffff0000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x3d, 0x10, 0xb2}),
				address:          0,
			},
			want: "orr	x3, x9, #0xffff0000ffff0000",
			wantErr: false,
		},
		{
			name: "orr	sp, x10, #0xe00000ffe00000ff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x29, 0x03, 0xb2}),
				address:          0,
			},
			want: "orr	sp, x10, #0xe00000ffe00000ff",
			wantErr: false,
		},
		{
			name: "orr	x9, x10, #0x3ff000003ff",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x25, 0x00, 0xb2}),
				address:          0,
			},
			want: "orr	x9, x10, #0x3ff000003ff",
			wantErr: false,
		},
		{
			name: "and	x14, x15, #0x8000800080008000",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x81, 0x01, 0x92}),
				address:          0,
			},
			want: "and	x14, x15, #0x8000800080008000",
			wantErr: false,
		},
		{
			name: "and	x12, x13, #0xffc3ffc3ffc3ffc3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xad, 0x0a, 0x92}),
				address:          0,
			},
			want: "and	x12, x13, #0xffc3ffc3ffc3ffc3",
			wantErr: false,
		},
		{
			name: "and	x11, xzr, #0x3000300030003",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xeb, 0x87, 0x00, 0x92}),
				address:          0,
			},
			want: "and	x11, xzr, #0x3000300030003",
			wantErr: false,
		},
		{
			name: "eor	x3, x6, #0xe0e0e0e0e0e0e0e0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc3, 0xc8, 0x03, 0xd2}),
				address:          0,
			},
			want: "eor	x3, x6, #0xe0e0e0e0e0e0e0e0",
			wantErr: false,
		},
		{
			name: "eor	sp, xzr, #0x303030303030303",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xc7, 0x00, 0xd2}),
				address:          0,
			},
			want: "eor	sp, xzr, #0x303030303030303",
			wantErr: false,
		},
		{
			name: "eor	x16, x17, #0x8181818181818181",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0xc6, 0x01, 0xd2}),
				address:          0,
			},
			want: "eor	x16, x17, #0x8181818181818181",
			wantErr: false,
		},
		{
			name: "tst	x18, #0xcccccccccccccccc",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xe6, 0x02, 0xf2}),
				address:          0,
			},
			want: "tst	x18, #0xcccccccccccccccc",
			wantErr: false,
		},
		{
			name: "ands	x19, x20, #0x3333333333333333",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0xe6, 0x00, 0xf2}),
				address:          0,
			},
			want: "ands	x19, x20, #0x3333333333333333",
			wantErr: false,
		},
		{
			name: "ands	x21, x22, #0x9999999999999999",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd5, 0xe6, 0x01, 0xf2}),
				address:          0,
			},
			want: "ands	x21, x22, #0x9999999999999999",
			wantErr: false,
		},
		{
			name: "tst	x3, #0xaaaaaaaaaaaaaaaa",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0xf0, 0x01, 0xf2}),
				address:          0,
			},
			want: "tst	x3, #0xaaaaaaaaaaaaaaaa",
			wantErr: false,
		},
		{
			name: "tst	xzr, #0x5555555555555555",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xf3, 0x00, 0xf2}),
				address:          0,
			},
			want: "tst	xzr, #0x5555555555555555",
			wantErr: false,
		},
		{
			name: "mov	w3, #983055",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x8f, 0x00, 0x32}),
				address:          0,
			},
			want: "mov	w3, #983055",
			wantErr: false,
		},
		// TODO: ADD THIS BACK IN
		// {
		// 	name: "mov	x10, #-6148914691236517206",
		// 	args: args{
		// 		instructionValue: binary.LittleEndian.Uint32([]byte{0xea, 0xf3, 0x01, 0xb2}),
		// 		address:          0,
		// 	},
		// 	want: "mov	x10, #-6148914691236517206",
		// 	wantErr: false,
		// },
		{
			name: "and	w2, w3, #0xfffffffd",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x62, 0x78, 0x1e, 0x12}),
				address:          0,
			},
			want: "and	w2, w3, #0xfffffffd",
			wantErr: false,
		},
		{
			name: "orr	w0, w1, #0xfffffffd",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x78, 0x1e, 0x32}),
				address:          0,
			},
			want: "orr	w0, w1, #0xfffffffd",
			wantErr: false,
		},
		{
			name: "eor	w16, w17, #0xfffffff9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x30, 0x76, 0x1d, 0x52}),
				address:          0,
			},
			want: "eor	w16, w17, #0xfffffff9",
			wantErr: false,
		},
		{
			name: "ands	w19, w20, #0xfffffff0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x93, 0x6e, 0x1c, 0x72}),
				address:          0,
			},
			want: "ands	w19, w20, #0xfffffff0",
			wantErr: false,
		},
		{
			name: "and	w12, w23, w21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x02, 0x15, 0x0a}),
				address:          0,
			},
			want: "and	w12, w23, w21",
			wantErr: false,
		},
		{
			name: "and	w16, w15, w1, lsl #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf0, 0x05, 0x01, 0x0a}),
				address:          0,
			},
			want: "and	w16, w15, w1, lsl #1",
			wantErr: false,
		},
		{
			name: "and	w9, w4, w10, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x7c, 0x0a, 0x0a}),
				address:          0,
			},
			want: "and	w9, w4, w10, lsl #31",
			wantErr: false,
		},
		{
			name: "and	w3, w30, w11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc3, 0x03, 0x0b, 0x0a}),
				address:          0,
			},
			want: "and	w3, w30, w11",
			wantErr: false,
		},
		{
			name: "and	x3, x5, x7, lsl #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0x07, 0x8a}),
				address:          0,
			},
			want: "and	x3, x5, x7, lsl #63",
			wantErr: false,
		},
		{
			name: "and	x5, x14, x19, asr #4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc5, 0x11, 0x93, 0x8a}),
				address:          0,
			},
			want: "and	x5, x14, x19, asr #4",
			wantErr: false,
		},
		{
			name: "and	w3, w17, w19, ror #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x7e, 0xd3, 0x0a}),
				address:          0,
			},
			want: "and	w3, w17, w19, ror #31",
			wantErr: false,
		},
		{
			name: "and	w0, w2, wzr, lsr #17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x44, 0x5f, 0x0a}),
				address:          0,
			},
			want: "and	w0, w2, wzr, lsr #17",
			wantErr: false,
		},
		{
			name: "and	w3, w30, w11, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc3, 0x03, 0x8b, 0x0a}),
				address:          0,
			},
			want: "and	w3, w30, w11, asr #0",
			wantErr: false,
		},
		{
			name: "and	xzr, x4, x26",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x00, 0x1a, 0x8a}),
				address:          0,
			},
			want: "and	xzr, x4, x26",
			wantErr: false,
		},
		{
			name: "and	w3, wzr, w20, ror #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0xd4, 0x0a}),
				address:          0,
			},
			want: "and	w3, wzr, w20, ror #0",
			wantErr: false,
		},
		{
			name: "and	x7, x20, xzr, asr #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x87, 0xfe, 0x9f, 0x8a}),
				address:          0,
			},
			want: "and	x7, x20, xzr, asr #63",
			wantErr: false,
		},
		{
			name: "bic	x13, x20, x14, lsl #47",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8d, 0xbe, 0x2e, 0x8a}),
				address:          0,
			},
			want: "bic	x13, x20, x14, lsl #47",
			wantErr: false,
		},
		{
			name: "bic	w2, w7, w9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x00, 0x29, 0x0a}),
				address:          0,
			},
			want: "bic	w2, w7, w9",
			wantErr: false,
		},
		{
			name: "orr	w2, w7, w0, asr #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe2, 0x7c, 0x80, 0x2a}),
				address:          0,
			},
			want: "orr	w2, w7, w0, asr #31",
			wantErr: false,
		},
		{
			name: "orr	x8, x9, x10, lsl #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x28, 0x31, 0x0a, 0xaa}),
				address:          0,
			},
			want: "orr	x8, x9, x10, lsl #12",
			wantErr: false,
		},
		{
			name: "orn	x3, x5, x7, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0xa7, 0xaa}),
				address:          0,
			},
			want: "orn	x3, x5, x7, asr #0",
			wantErr: false,
		},
		{
			name: "orn	w2, w5, w29",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa2, 0x00, 0x3d, 0x2a}),
				address:          0,
			},
			want: "orn	w2, w5, w29",
			wantErr: false,
		},
		{
			name: "ands	w7, wzr, w9, lsl #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe7, 0x07, 0x09, 0x6a}),
				address:          0,
			},
			want: "ands	w7, wzr, w9, lsl #1",
			wantErr: false,
		},
		{
			name: "ands	x3, x5, x20, ror #63",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0xfc, 0xd4, 0xea}),
				address:          0,
			},
			want: "ands	x3, x5, x20, ror #63",
			wantErr: false,
		},
		{
			name: "bics	w3, w5, w7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa3, 0x00, 0x27, 0x6a}),
				address:          0,
			},
			want: "bics	w3, w5, w7",
			wantErr: false,
		},
		{
			name: "bics	x3, xzr, x3, lsl #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x07, 0x23, 0xea}),
				address:          0,
			},
			want: "bics	x3, xzr, x3, lsl #1",
			wantErr: false,
		},
		{
			name: "tst	w3, w7, lsl #31",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x7c, 0x07, 0x6a}),
				address:          0,
			},
			want: "tst	w3, w7, lsl #31",
			wantErr: false,
		},
		{
			name: "tst	x2, x20, asr #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x00, 0x94, 0xea}),
				address:          0,
			},
			want: "tst	x2, x20, asr #0",
			wantErr: false,
		},
		{
			name: "mov	x3, x6",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x06, 0xaa}),
				address:          0,
			},
			want: "mov	x3, x6",
			wantErr: false,
		},
		{
			name: "mov	x3, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x1f, 0xaa}),
				address:          0,
			},
			want: "mov	x3, xzr",
			wantErr: false,
		},
		{
			name: "mov	wzr, w2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x03, 0x02, 0x2a}),
				address:          0,
			},
			want: "mov	wzr, w2",
			wantErr: false,
		},
		{
			name: "mov	w3, w5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe3, 0x03, 0x05, 0x2a}),
				address:          0,
			},
			want: "mov	w3, w5",
			wantErr: false,
		},
		{
			name: "mov	w1, #65535",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe1, 0xff, 0x9f, 0x52}),
				address:          0,
			},
			want: "mov	w1, #65535",
			wantErr: false,
		},
		{
			name: "movz	w2, #0, lsl #16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x02, 0x00, 0xa0, 0x52}),
				address:          0,
			},
			want: "movz	w2, #0, lsl #16",
			wantErr: false,
		},
		{
			name: "mov	w2, #-1235",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0x9a, 0x80, 0x12}),
				address:          0,
			},
			want: "mov	w2, #-1235",
			wantErr: false,
		},
		{
			name: "mov	x2, #5299989643264",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0x9a, 0xc0, 0xd2}),
				address:          0,
			},
			want: "mov	x2, #5299989643264",
			wantErr: false,
		},
		{
			name: "movk	xzr, #4321, lsl #48",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x1c, 0xe2, 0xf2}),
				address:          0,
			},
			want: "movk	xzr, #4321, lsl #48",
			wantErr: false,
		},
		{
			name: "adrp	x30, #4096",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1e, 0x00, 0x00, 0xb0}),
				address:          0,
			},
			want: "adrp	x30, #4096",
			wantErr: false,
		},
		{
			name: "adr	x20, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x00, 0x10}),
				address:          0,
			},
			want: "adr	x20, #0",
			wantErr: false,
		},
		{
			name: "adr	x9, #-1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xff, 0xff, 0x70}),
				address:          0,
			},
			want: "adr	x9, #-1",
			wantErr: false,
		},
		{
			name: "adr	x5, #1048575",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0xff, 0x7f, 0x70}),
				address:          0,
			},
			want: "adr	x5, #1048575",
			wantErr: false,
		},
		{
			name: "adr	x9, #1048575",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xff, 0x7f, 0x70}),
				address:          0,
			},
			want: "adr	x9, #1048575",
			wantErr: false,
		},
		{
			name: "adr	x2, #-1048576",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x02, 0x00, 0x80, 0x10}),
				address:          0,
			},
			want: "adr	x2, #-1048576",
			wantErr: false,
		},
		{
			name: "adrp	x9, #4294963200",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xff, 0x7f, 0xf0}),
				address:          0,
			},
			want: "adrp	x9, #4294963200",
			wantErr: false,
		},
		{
			name: "adrp	x20, #-4294967296",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x14, 0x00, 0x80, 0x90}),
				address:          0,
			},
			want: "adrp	x20, #-4294967296",
			wantErr: false,
		},
		{
			name: "nop",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "nop",
			wantErr: false,
		},
		{
			name: "hint	#127",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x2f, 0x03, 0xd5}),
				address:          0,
			},
			want: "hint	#127",
			wantErr: false,
		},
		{
			name: "nop",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "nop",
			wantErr: false,
		},
		{
			name: "yield",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "yield",
			wantErr: false,
		},
		{
			name: "wfe",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "wfe",
			wantErr: false,
		},
		{
			name: "wfi",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x7f, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "wfi",
			wantErr: false,
		},
		{
			name: "sev",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "sev",
			wantErr: false,
		},
		{
			name: "sevl",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "sevl",
			wantErr: false,
		},
		{
			name: "dgh",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x20, 0x03, 0xd5}),
				address:          0,
			},
			want:    "dgh",
			wantErr: false,
		},
		{
			name: "clrex",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want:    "clrex",
			wantErr: false,
		},
		{
			name: "clrex	#0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x30, 0x03, 0xd5}),
				address:          0,
			},
			want: "clrex	#0",
			wantErr: false,
		},
		{
			name: "clrex	#7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x37, 0x03, 0xd5}),
				address:          0,
			},
			want: "clrex	#7",
			wantErr: false,
		},
		{
			name: "clrex",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want:    "clrex",
			wantErr: false,
		},
		{
			name: "ssbb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x30, 0x03, 0xd5}),
				address:          0,
			},
			want:    "ssbb",
			wantErr: false,
		},
		{
			name: "pssbb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x34, 0x03, 0xd5}),
				address:          0,
			},
			want:    "pssbb",
			wantErr: false,
		},
		{
			name: "dsb	#12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3c, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	#12",
			wantErr: false,
		},
		{
			name: "dsb	sy",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	sy",
			wantErr: false,
		},
		{
			name: "dsb	oshld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x31, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	oshld",
			wantErr: false,
		},
		{
			name: "dsb	oshst",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x32, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	oshst",
			wantErr: false,
		},
		{
			name: "dsb	osh",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x33, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	osh",
			wantErr: false,
		},
		{
			name: "dsb	nshld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x35, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	nshld",
			wantErr: false,
		},
		{
			name: "dsb	nshst",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x36, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	nshst",
			wantErr: false,
		},
		{
			name: "dsb	nsh",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x37, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	nsh",
			wantErr: false,
		},
		{
			name: "dsb	ishld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x39, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	ishld",
			wantErr: false,
		},
		{
			name: "dsb	ishst",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3a, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	ishst",
			wantErr: false,
		},
		{
			name: "dsb	ish",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3b, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	ish",
			wantErr: false,
		},
		{
			name: "dsb	ld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3d, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	ld",
			wantErr: false,
		},
		{
			name: "dsb	st",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3e, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	st",
			wantErr: false,
		},
		{
			name: "dsb	sy",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want: "dsb	sy",
			wantErr: false,
		},
		{
			name: "dmb	#0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x30, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	#0",
			wantErr: false,
		},
		{
			name: "dmb	#12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3c, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	#12",
			wantErr: false,
		},
		{
			name: "dmb	sy",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	sy",
			wantErr: false,
		},
		{
			name: "dmb	oshld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x31, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	oshld",
			wantErr: false,
		},
		{
			name: "dmb	oshst",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x32, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	oshst",
			wantErr: false,
		},
		{
			name: "dmb	osh",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x33, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	osh",
			wantErr: false,
		},
		{
			name: "dmb	nshld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x35, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	nshld",
			wantErr: false,
		},
		{
			name: "dmb	nshst",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x36, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	nshst",
			wantErr: false,
		},
		{
			name: "dmb	nsh",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x37, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	nsh",
			wantErr: false,
		},
		{
			name: "dmb	ishld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x39, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	ishld",
			wantErr: false,
		},
		{
			name: "dmb	ishst",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3a, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	ishst",
			wantErr: false,
		},
		{
			name: "dmb	ish",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3b, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	ish",
			wantErr: false,
		},
		{
			name: "dmb	ld",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3d, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	ld",
			wantErr: false,
		},
		{
			name: "dmb	st",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3e, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	st",
			wantErr: false,
		},
		{
			name: "dmb	sy",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want: "dmb	sy",
			wantErr: false,
		},
		{
			name: "isb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want:    "isb",
			wantErr: false,
		},
		{
			name: "isb",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x3f, 0x03, 0xd5}),
				address:          0,
			},
			want:    "isb",
			wantErr: false,
		},
		{
			name: "isb	#12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x3c, 0x03, 0xd5}),
				address:          0,
			},
			want: "isb	#12",
			wantErr: false,
		},
		{
			name: "msr	spsel, #0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xbf, 0x40, 0x00, 0xd5}),
				address:          0,
			},
			want: "msr	spsel, #0",
			wantErr: false,
		},
		{
			name: "msr	daifset, #15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x4f, 0x03, 0xd5}),
				address:          0,
			},
			want: "msr	daifset, #15",
			wantErr: false,
		},
		{
			name: "msr	daifclr, #12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0x4c, 0x03, 0xd5}),
				address:          0,
			},
			want: "msr	daifclr, #12",
			wantErr: false,
		},
		{
			name: "sys	#7, c5, c9, #7, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0x59, 0x0f, 0xd5}),
				address:          0,
			},
			want: "sys	#7, c5, c9, #7, x5",
			wantErr: false,
		},
		{
			name: "sys	#0, c15, c15, #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5f, 0xff, 0x08, 0xd5}),
				address:          0,
			},
			want: "sys	#0, c15, c15, #2",
			wantErr: false,
		},
		{
			name: "sysl	x9, #7, c5, c9, #7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x59, 0x2f, 0xd5}),
				address:          0,
			},
			want: "sysl	x9, #7, c5, c9, #7",
			wantErr: false,
		},
		{
			name: "sysl	x1, #0, c15, c15, #2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x41, 0xff, 0x28, 0xd5}),
				address:          0,
			},
			want: "sysl	x1, #0, c15, c15, #2",
			wantErr: false,
		},
		{
			name: "ic	ialluis",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x71, 0x08, 0xd5}),
				address:          0,
			},
			want: "ic	ialluis",
			wantErr: false,
		},
		{
			name: "ic	iallu",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x75, 0x08, 0xd5}),
				address:          0,
			},
			want: "ic	iallu",
			wantErr: false,
		},
		{
			name: "ic	ivau, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x75, 0x0b, 0xd5}),
				address:          0,
			},
			want: "ic	ivau, x9",
			wantErr: false,
		},
		{
			name: "dc	zva, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x74, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	zva, x12",
			wantErr: false,
		},
		{
			name: "dc	ivac, xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x3f, 0x76, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	ivac, xzr",
			wantErr: false,
		},
		{
			name: "dc	isw, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x42, 0x76, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	isw, x2",
			wantErr: false,
		},
		{
			name: "dc	cvac, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x7a, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cvac, x9",
			wantErr: false,
		},
		{
			name: "dc	csw, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4a, 0x7a, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	csw, x10",
			wantErr: false,
		},
		{
			name: "dc	cvau, x0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x20, 0x7b, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	cvau, x0",
			wantErr: false,
		},
		{
			name: "dc	civac, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x7e, 0x0b, 0xd5}),
				address:          0,
			},
			want: "dc	civac, x3",
			wantErr: false,
		},
		{
			name: "dc	cisw, x30",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x5e, 0x7e, 0x08, 0xd5}),
				address:          0,
			},
			want: "dc	cisw, x30",
			wantErr: false,
		},
		{
			name: "at	s1e1r, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x78, 0x08, 0xd5}),
				address:          0,
			},
			want: "at	s1e1r, x19",
			wantErr: false,
		},
		{
			name: "at	s1e2r, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x78, 0x0c, 0xd5}),
				address:          0,
			},
			want: "at	s1e2r, x19",
			wantErr: false,
		},
		{
			name: "at	s1e3r, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x13, 0x78, 0x0e, 0xd5}),
				address:          0,
			},
			want: "at	s1e3r, x19",
			wantErr: false,
		},
		{
			name: "at	s1e1w, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x78, 0x08, 0xd5}),
				address:          0,
			},
			want: "at	s1e1w, x19",
			wantErr: false,
		},
		{
			name: "at	s1e2w, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x78, 0x0c, 0xd5}),
				address:          0,
			},
			want: "at	s1e2w, x19",
			wantErr: false,
		},
		{
			name: "at	s1e3w, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x78, 0x0e, 0xd5}),
				address:          0,
			},
			want: "at	s1e3w, x19",
			wantErr: false,
		},
		{
			name: "at	s1e0r, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x53, 0x78, 0x08, 0xd5}),
				address:          0,
			},
			want: "at	s1e0r, x19",
			wantErr: false,
		},
		{
			name: "at	s1e0w, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x73, 0x78, 0x08, 0xd5}),
				address:          0,
			},
			want: "at	s1e0w, x19",
			wantErr: false,
		},
		{
			name: "at	s12e1r, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x94, 0x78, 0x0c, 0xd5}),
				address:          0,
			},
			want: "at	s12e1r, x20",
			wantErr: false,
		},
		{
			name: "at	s12e1w, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb4, 0x78, 0x0c, 0xd5}),
				address:          0,
			},
			want: "at	s12e1w, x20",
			wantErr: false,
		},
		{
			name: "at	s12e0r, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xd4, 0x78, 0x0c, 0xd5}),
				address:          0,
			},
			want: "at	s12e0r, x20",
			wantErr: false,
		},
		{
			name: "at	s12e0w, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf4, 0x78, 0x0c, 0xd5}),
				address:          0,
			},
			want: "at	s12e0w, x20",
			wantErr: false,
		},
		{
			name: "tlbi	ipas2e1is, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x24, 0x80, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	ipas2e1is, x4",
			wantErr: false,
		},
		{
			name: "tlbi	ipas2le1is, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x80, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	ipas2le1is, x9",
			wantErr: false,
		},
		{
			name: "tlbi	vmalle1is",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x83, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vmalle1is",
			wantErr: false,
		},
		{
			name: "tlbi	alle2is",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x83, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	alle2is",
			wantErr: false,
		},
		{
			name: "tlbi	alle3is",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x83, 0x0e, 0xd5}),
				address:          0,
			},
			want: "tlbi	alle3is",
			wantErr: false,
		},
		{
			name: "tlbi	vae1is, x1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x21, 0x83, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vae1is, x1",
			wantErr: false,
		},
		{
			name: "tlbi	vae2is, x2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x22, 0x83, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	vae2is, x2",
			wantErr: false,
		},
		{
			name: "tlbi	vae3is, x3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x23, 0x83, 0x0e, 0xd5}),
				address:          0,
			},
			want: "tlbi	vae3is, x3",
			wantErr: false,
		},
		{
			name: "tlbi	aside1is, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x45, 0x83, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	aside1is, x5",
			wantErr: false,
		},
		{
			name: "tlbi	vaae1is, x9",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x83, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vaae1is, x9",
			wantErr: false,
		},
		{
			name: "tlbi	alle1is",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x83, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	alle1is",
			wantErr: false,
		},
		{
			name: "tlbi	vale1is, x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xaa, 0x83, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vale1is, x10",
			wantErr: false,
		},
		{
			name: "tlbi	vale2is, x11",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xab, 0x83, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	vale2is, x11",
			wantErr: false,
		},
		{
			name: "tlbi	vale3is, x13",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xad, 0x83, 0x0e, 0xd5}),
				address:          0,
			},
			want: "tlbi	vale3is, x13",
			wantErr: false,
		},
		{
			name: "tlbi	vmalls12e1is",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x83, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	vmalls12e1is",
			wantErr: false,
		},
		{
			name: "tlbi	vaale1is, x14",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xee, 0x83, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vaale1is, x14",
			wantErr: false,
		},
		{
			name: "tlbi	ipas2e1, x15",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2f, 0x84, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	ipas2e1, x15",
			wantErr: false,
		},
		{
			name: "tlbi	ipas2le1, x16",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb0, 0x84, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	ipas2le1, x16",
			wantErr: false,
		},
		{
			name: "tlbi	vmalle1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x87, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vmalle1",
			wantErr: false,
		},
		{
			name: "tlbi	alle2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x87, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	alle2",
			wantErr: false,
		},
		{
			name: "tlbi	alle3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x1f, 0x87, 0x0e, 0xd5}),
				address:          0,
			},
			want: "tlbi	alle3",
			wantErr: false,
		},
		{
			name: "tlbi	vae1, x17",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x31, 0x87, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vae1, x17",
			wantErr: false,
		},
		{
			name: "tlbi	vae2, x18",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x32, 0x87, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	vae2, x18",
			wantErr: false,
		},
		{
			name: "tlbi	vae3, x19",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x33, 0x87, 0x0e, 0xd5}),
				address:          0,
			},
			want: "tlbi	vae3, x19",
			wantErr: false,
		},
		{
			name: "tlbi	aside1, x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x54, 0x87, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	aside1, x20",
			wantErr: false,
		},
		{
			name: "tlbi	vaae1, x21",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x75, 0x87, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vaae1, x21",
			wantErr: false,
		},
		{
			name: "tlbi	alle1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x9f, 0x87, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	alle1",
			wantErr: false,
		},
		{
			name: "tlbi	vale1, x22",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb6, 0x87, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vale1, x22",
			wantErr: false,
		},
		{
			name: "tlbi	vale2, x23",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb7, 0x87, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	vale2, x23",
			wantErr: false,
		},
		{
			name: "tlbi	vale3, x24",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xb8, 0x87, 0x0e, 0xd5}),
				address:          0,
			},
			want: "tlbi	vale3, x24",
			wantErr: false,
		},
		{
			name: "tlbi	vmalls12e1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xdf, 0x87, 0x0c, 0xd5}),
				address:          0,
			},
			want: "tlbi	vmalls12e1",
			wantErr: false,
		},
		{
			name: "tlbi	vaale1, x25",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xf9, 0x87, 0x08, 0xd5}),
				address:          0,
			},
			want: "tlbi	vaale1, x25",
			wantErr: false,
		},
		{
			name: "msr	teecr32_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x00, 0x12, 0xd5}),
				address:          0,
			},
			want: "msr	teecr32_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	osdtrrx_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x00, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	osdtrrx_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	mdccint_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x02, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	mdccint_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	mdscr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x02, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	mdscr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	osdtrtx_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x03, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	osdtrtx_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgdtr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x04, 0x13, 0xd5}),
				address:          0,
			},
			want: "msr	dbgdtr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgdtrtx_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x05, 0x13, 0xd5}),
				address:          0,
			},
			want: "msr	dbgdtrtx_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	oseccr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x06, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	oseccr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgvcr32_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x07, 0x14, 0xd5}),
				address:          0,
			},
			want: "msr	dbgvcr32_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr0_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x00, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr0_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr1_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x01, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr1_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr2_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x02, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr2_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr3_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x03, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr3_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr4_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x04, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr4_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr5_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x05, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr5_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr6_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x06, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr6_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr7_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x07, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr7_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr8_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x08, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr8_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr9_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x09, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr9_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr10_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x0a, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr10_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr11_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x0b, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr11_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr12_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x0c, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr12_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr13_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x0d, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr13_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr14_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x0e, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr14_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbvr15_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x0f, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbvr15_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr0_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x00, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr0_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr1_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x01, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr1_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr2_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x02, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr2_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr3_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x03, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr3_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr4_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x04, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr4_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr5_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x05, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr5_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr6_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x06, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr6_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr7_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x07, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr7_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr8_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x08, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr8_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr9_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x09, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr9_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr10_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0a, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr10_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr11_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0b, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr11_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr12_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0c, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr12_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr13_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0d, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr13_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr14_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0e, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr14_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgbcr15_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x0f, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgbcr15_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr0_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x00, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr0_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr1_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x01, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr1_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr2_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x02, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr2_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr3_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x03, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr3_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr4_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x04, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr4_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr5_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x05, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr5_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr6_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x06, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr6_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr7_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x07, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr7_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr8_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x08, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr8_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr9_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x09, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr9_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr10_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x0a, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr10_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr11_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x0b, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr11_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr12_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x0c, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr12_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr13_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x0d, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr13_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr14_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x0e, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr14_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwvr15_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x0f, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwvr15_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr0_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x00, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr0_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr1_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x01, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr1_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr2_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x02, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr2_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr3_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x03, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr3_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr4_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x04, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr4_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr5_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x05, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr5_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr6_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x06, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr6_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr7_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x07, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr7_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr8_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x08, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr8_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr9_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x09, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr9_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr10_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x0a, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr10_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr11_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x0b, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr11_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr12_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x0c, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr12_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr13_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x0d, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr13_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr14_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x0e, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr14_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgwcr15_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x0f, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgwcr15_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	teehbr32_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x10, 0x12, 0xd5}),
				address:          0,
			},
			want: "msr	teehbr32_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	oslar_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x10, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	oslar_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	osdlr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x13, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	osdlr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgprcr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x14, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgprcr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgclaimset_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x78, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgclaimset_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	dbgclaimclr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0x79, 0x10, 0xd5}),
				address:          0,
			},
			want: "msr	dbgclaimclr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	csselr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x00, 0x1a, 0xd5}),
				address:          0,
			},
			want: "msr	csselr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	vpidr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x00, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	vpidr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	vmpidr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x00, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	vmpidr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	sctlr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x10, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	sctlr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	sctlr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x10, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	sctlr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	sctlr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x10, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	sctlr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	actlr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x10, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	actlr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	actlr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x10, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	actlr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	actlr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x10, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	actlr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	cpacr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x10, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	cpacr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	hcr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hcr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	scr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x11, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	scr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	mdcr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	mdcr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	sder32_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x11, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	sder32_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	cptr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cptr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	cptr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x11, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	cptr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	hstr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hstr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	hacr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0x11, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hacr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	mdcr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x13, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	mdcr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	ttbr0_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x20, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr0_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	ttbr0_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x20, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr0_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	ttbr0_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x20, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr0_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	ttbr1_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x20, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	ttbr1_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	tcr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x20, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	tcr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	tcr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x20, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	tcr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	tcr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x20, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	tcr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	vttbr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x21, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	vttbr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	vtcr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x21, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	vtcr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	dacr32_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x30, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	dacr32_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x40, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x40, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x40, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	elr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	elr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	elr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	elr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	elr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x40, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	elr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	sp_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x41, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	sp_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	sp_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x41, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	sp_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	sp_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x41, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	sp_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	spsel, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x42, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	spsel, x12",
			wantErr: false,
		},
		{
			name: "msr	nzcv, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x42, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	nzcv, x12",
			wantErr: false,
		},
		{
			name: "msr	daif, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x42, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	daif, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_irq, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x43, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_irq, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_abt, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x43, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_abt, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_und, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x43, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_und, x12",
			wantErr: false,
		},
		{
			name: "msr	spsr_fiq, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0x43, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	spsr_fiq, x12",
			wantErr: false,
		},
		{
			name: "msr	fpcr, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x44, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	fpcr, x12",
			wantErr: false,
		},
		{
			name: "msr	fpsr, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x44, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	fpsr, x12",
			wantErr: false,
		},
		{
			name: "msr	dspsr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x45, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	dspsr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	dlr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x45, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	dlr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	ifsr32_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x50, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	ifsr32_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	afsr0_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x51, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	afsr0_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	afsr0_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x51, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	afsr0_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	afsr0_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x51, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	afsr0_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	afsr1_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x51, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	afsr1_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	afsr1_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x51, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	afsr1_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	afsr1_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x51, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	afsr1_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	esr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x52, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	esr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	esr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x52, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	esr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	esr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x52, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	esr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	fpexc32_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x53, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	fpexc32_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	far_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x60, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	far_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	far_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x60, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	far_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	far_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x60, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	far_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	hpfar_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0x60, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	hpfar_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	par_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x74, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	par_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	pmcr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x9c, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmcr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmcntenset_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x9c, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmcntenset_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmcntenclr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x9c, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmcntenclr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmovsclr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0x9c, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmovsclr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmselr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0x9c, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmselr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmccntr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x9d, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmccntr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmxevtyper_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x9d, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmxevtyper_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmxevcntr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x9d, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmxevcntr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmuserenr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0x9e, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmuserenr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmintenset_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0x9e, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmintenset_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	pmintenclr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0x9e, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	pmintenclr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	pmovsset_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0x9e, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmovsset_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	mair_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xa2, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	mair_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	mair_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xa2, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	mair_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	mair_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xa2, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	mair_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	amair_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xa3, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	amair_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	amair_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xa3, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	amair_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	amair_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xa3, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	amair_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	vbar_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xc0, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	vbar_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	vbar_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xc0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	vbar_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	vbar_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xc0, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	vbar_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	rmr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xc0, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	rmr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	rmr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xc0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	rmr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	rmr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xc0, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	rmr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	contextidr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xd0, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	contextidr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	tpidr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xd0, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	tpidr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	tpidr_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xd0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	tpidr_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	tpidr_el3, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xd0, 0x1e, 0xd5}),
				address:          0,
			},
			want: "msr	tpidr_el3, x12",
			wantErr: false,
		},
		{
			name: "msr	tpidrro_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xd0, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	tpidrro_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	tpidr_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xd0, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	tpidr_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	cntfrq_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe0, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntfrq_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	cntvoff_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xe0, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cntvoff_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	cntkctl_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe1, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	cntkctl_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	cnthctl_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe1, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthctl_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	cntp_tval_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe2, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntp_tval_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	cnthp_tval_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe2, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthp_tval_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	cntps_tval_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe2, 0x1f, 0xd5}),
				address:          0,
			},
			want: "msr	cntps_tval_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	cntp_ctl_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xe2, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntp_ctl_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	cnthp_ctl_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xe2, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthp_ctl_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	cntps_ctl_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xe2, 0x1f, 0xd5}),
				address:          0,
			},
			want: "msr	cntps_ctl_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	cntp_cval_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xe2, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntp_cval_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	cnthp_cval_el2, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xe2, 0x1c, 0xd5}),
				address:          0,
			},
			want: "msr	cnthp_cval_el2, x12",
			wantErr: false,
		},
		{
			name: "msr	cntps_cval_el1, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xe2, 0x1f, 0xd5}),
				address:          0,
			},
			want: "msr	cntps_cval_el1, x12",
			wantErr: false,
		},
		{
			name: "msr	cntv_tval_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe3, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntv_tval_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	cntv_ctl_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xe3, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntv_ctl_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	cntv_cval_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xe3, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	cntv_cval_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr0_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr0_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr1_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr1_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr2_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr2_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr3_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr3_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr4_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr4_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr5_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr5_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr6_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr6_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr7_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xe8, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr7_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr8_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr8_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr9_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr9_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr10_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr10_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr11_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr11_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr12_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr12_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr13_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr13_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr14_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr14_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr15_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xe9, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr15_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr16_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr16_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr17_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr17_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr18_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr18_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr19_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr19_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr20_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr20_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr21_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr21_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr22_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr22_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr23_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xea, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr23_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr24_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr24_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr25_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr25_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr26_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr26_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr27_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr27_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr28_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr28_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr29_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr29_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevcntr30_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xeb, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevcntr30_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmccfiltr_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmccfiltr_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper0_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper0_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper1_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper1_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper2_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper2_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper3_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper3_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper4_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper4_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper5_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper5_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper6_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper6_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper7_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xec, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper7_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper8_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper8_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper9_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper9_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper10_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper10_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper11_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper11_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper12_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper12_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper13_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper13_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper14_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper14_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper15_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xed, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper15_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper16_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper16_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper17_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper17_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper18_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper18_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper19_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper19_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper20_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper20_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper21_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper21_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper22_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper22_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper23_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xec, 0xee, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper23_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper24_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper24_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper25_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2c, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper25_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper26_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x4c, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper26_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper27_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x6c, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper27_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper28_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x8c, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper28_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper29_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper29_el0, x12",
			wantErr: false,
		},
		{
			name: "msr	pmevtyper30_el0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xcc, 0xef, 0x1b, 0xd5}),
				address:          0,
			},
			want: "msr	pmevtyper30_el0, x12",
			wantErr: false,
		},
		{
			name: "mrs	x9, teecr32_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x00, 0x32, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, teecr32_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, osdtrrx_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x00, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, osdtrrx_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mdccsr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x01, 0x33, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mdccsr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, mdccint_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x02, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mdccint_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mdscr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x02, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mdscr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, osdtrtx_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x03, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, osdtrtx_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgdtr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x04, 0x33, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgdtr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgdtrrx_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x05, 0x33, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgdtrrx_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, oseccr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x06, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, oseccr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgvcr32_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x07, 0x34, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgvcr32_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x00, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x01, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x02, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr3_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x03, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr3_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr4_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x04, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr4_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr5_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x05, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr5_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr6_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x06, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr6_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr7_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x07, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr7_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr8_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x08, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr8_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr9_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x09, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr9_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr10_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x0a, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr10_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr11_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x0b, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr11_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr12_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x0c, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr12_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr13_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x0d, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr13_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr14_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x0e, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr14_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbvr15_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x0f, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbvr15_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x00, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x01, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x02, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr3_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x03, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr3_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr4_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x04, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr4_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr5_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x05, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr5_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr6_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x06, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr6_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr7_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x07, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr7_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr8_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x08, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr8_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr9_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x09, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr9_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr10_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x0a, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr10_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr11_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x0b, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr11_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr12_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x0c, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr12_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr13_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x0d, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr13_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr14_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x0e, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr14_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgbcr15_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x0f, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgbcr15_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x00, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x01, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x02, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr3_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x03, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr3_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr4_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x04, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr4_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr5_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x05, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr5_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr6_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x06, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr6_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr7_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x07, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr7_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr8_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x08, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr8_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr9_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x09, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr9_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr10_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0a, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr10_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr11_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0b, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr11_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr12_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0c, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr12_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr13_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0d, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr13_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr14_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0e, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr14_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwvr15_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x0f, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwvr15_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x00, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x01, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x02, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr3_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x03, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr3_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr4_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x04, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr4_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr5_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x05, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr5_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr6_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x06, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr6_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr7_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x07, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr7_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr8_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x08, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr8_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr9_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x09, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr9_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr10_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0a, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr10_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr11_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0b, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr11_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr12_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0c, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr12_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr13_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0d, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr13_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr14_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0e, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr14_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgwcr15_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x0f, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgwcr15_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mdrar_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x10, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mdrar_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, teehbr32_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x10, 0x32, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, teehbr32_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, oslsr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x11, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, oslsr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, osdlr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x13, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, osdlr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgprcr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x14, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgprcr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgclaimset_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x78, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgclaimset_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgclaimclr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x79, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgclaimclr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dbgauthstatus_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x7e, 0x30, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dbgauthstatus_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, midr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x00, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, midr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, ccsidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x00, 0x39, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ccsidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, csselr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x00, 0x3a, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, csselr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, vpidr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x00, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vpidr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, clidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x00, 0x39, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, clidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, ctr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x00, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ctr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, mpidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x00, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mpidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, vmpidr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x00, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vmpidr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, revidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x00, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, revidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, aidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x00, 0x39, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, aidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, dczid_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x00, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dczid_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_pfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_pfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_pfr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_pfr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_dfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_dfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_afr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_afr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_mmfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_mmfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_mmfr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_mmfr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_mmfr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_mmfr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_mmfr3_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x01, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_mmfr3_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_mmfr4_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_mmfr4_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_mmfr5_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x03, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_mmfr5_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_isar0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_isar0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_isar1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_isar1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_isar2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_isar2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_isar3_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_isar3_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_isar4_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_isar4_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_isar5_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x02, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_isar5_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mvfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x03, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mvfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mvfr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x03, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mvfr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mvfr2_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x03, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mvfr2_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64pfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x04, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64pfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64pfr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x04, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64pfr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64dfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x05, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64dfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64dfr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x05, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64dfr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64afr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x05, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64afr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64afr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x05, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64afr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64isar0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x06, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64isar0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64isar1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x06, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64isar1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64mmfr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x07, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64mmfr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, id_aa64mmfr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x07, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, id_aa64mmfr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, sctlr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x10, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sctlr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, sctlr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x10, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sctlr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, sctlr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x10, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sctlr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, actlr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x10, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, actlr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, actlr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x10, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, actlr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, actlr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x10, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, actlr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, cpacr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x10, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cpacr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, hcr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, hcr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, scr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x11, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, scr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, mdcr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mdcr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, sder32_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x11, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sder32_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, cptr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cptr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, cptr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x11, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cptr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, hstr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, hstr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, hacr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x11, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, hacr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, mdcr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x13, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mdcr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, ttbr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x20, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ttbr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, ttbr0_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x20, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ttbr0_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, ttbr0_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x20, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ttbr0_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, ttbr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x20, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ttbr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, tcr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x20, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tcr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, tcr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x20, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tcr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, tcr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x20, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tcr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, vttbr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x21, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vttbr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, vtcr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x21, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vtcr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, dacr32_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x30, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dacr32_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x40, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x40, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x40, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, elr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x40, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, elr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, elr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x40, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, elr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, elr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x40, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, elr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, sp_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x41, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sp_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, sp_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x41, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sp_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, sp_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x41, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, sp_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsel",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x42, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsel",
			wantErr: false,
		},
		{
			name: "mrs	x9, nzcv",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x42, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, nzcv",
			wantErr: false,
		},
		{
			name: "mrs	x9, daif",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x42, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, daif",
			wantErr: false,
		},
		{
			name: "mrs	x9, currentel",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x42, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, currentel",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_irq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x43, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_irq",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_abt",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x43, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_abt",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_und",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x43, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_und",
			wantErr: false,
		},
		{
			name: "mrs	x9, spsr_fiq",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x43, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, spsr_fiq",
			wantErr: false,
		},
		{
			name: "mrs	x9, fpcr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x44, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, fpcr",
			wantErr: false,
		},
		{
			name: "mrs	x9, fpsr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x44, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, fpsr",
			wantErr: false,
		},
		{
			name: "mrs	x9, dspsr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x45, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dspsr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, dlr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x45, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, dlr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, ifsr32_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x50, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, ifsr32_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, afsr0_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x51, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, afsr0_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, afsr0_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x51, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, afsr0_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, afsr0_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x51, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, afsr0_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, afsr1_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x51, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, afsr1_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, afsr1_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x51, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, afsr1_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, afsr1_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x51, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, afsr1_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, esr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x52, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, esr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, esr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x52, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, esr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, esr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x52, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, esr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, fpexc32_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x53, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, fpexc32_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, far_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x60, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, far_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, far_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x60, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, far_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, far_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x60, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, far_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, hpfar_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0x60, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, hpfar_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, par_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x74, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, par_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmcr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmcr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmcntenset_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmcntenset_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmcntenclr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmcntenclr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmovsclr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmovsclr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmselr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmselr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmceid0_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmceid0_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmceid1_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0x9c, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmceid1_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmccntr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x9d, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmccntr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmxevtyper_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x9d, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmxevtyper_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmxevcntr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x9d, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmxevcntr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmuserenr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0x9e, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmuserenr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmintenset_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0x9e, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmintenset_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmintenclr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0x9e, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmintenclr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmovsset_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0x9e, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmovsset_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, mair_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xa2, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mair_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, mair_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xa2, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mair_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, mair_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xa2, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, mair_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, amair_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xa3, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, amair_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, amair_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xa3, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, amair_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, amair_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xa3, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, amair_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, vbar_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xc0, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vbar_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, vbar_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xc0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vbar_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, vbar_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xc0, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, vbar_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, rvbar_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xc0, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, rvbar_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, rvbar_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xc0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, rvbar_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, rvbar_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xc0, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, rvbar_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, rmr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xc0, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, rmr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, rmr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xc0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, rmr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, rmr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xc0, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, rmr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, isr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xc1, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, isr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, contextidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xd0, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, contextidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, tpidr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xd0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tpidr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, tpidr_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xd0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tpidr_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, tpidr_el3",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xd0, 0x3e, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tpidr_el3",
			wantErr: false,
		},
		{
			name: "mrs	x9, tpidrro_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xd0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tpidrro_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, tpidr_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xd0, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, tpidr_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntfrq_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntfrq_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntpct_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntpct_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntvct_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe0, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntvct_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntvoff_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xe0, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntvoff_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntkctl_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe1, 0x38, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntkctl_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, cnthctl_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe1, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cnthctl_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntp_tval_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe2, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntp_tval_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cnthp_tval_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe2, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cnthp_tval_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntps_tval_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe2, 0x3f, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntps_tval_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntp_ctl_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe2, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntp_ctl_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cnthp_ctl_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe2, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cnthp_ctl_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntps_ctl_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe2, 0x3f, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntps_ctl_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntp_cval_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe2, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntp_cval_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cnthp_cval_el2",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe2, 0x3c, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cnthp_cval_el2",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntps_cval_el1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe2, 0x3f, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntps_cval_el1",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntv_tval_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe3, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntv_tval_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntv_ctl_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe3, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntv_ctl_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, cntv_cval_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe3, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, cntv_cval_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr0_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr0_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr1_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr1_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr2_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr2_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr3_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr3_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr4_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr4_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr5_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr5_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr6_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr6_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr7_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xe8, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr7_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr8_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr8_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr9_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr9_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr10_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr10_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr11_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr11_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr12_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr12_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr13_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr13_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr14_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr14_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr15_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xe9, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr15_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr16_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr16_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr17_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr17_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr18_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr18_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr19_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr19_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr20_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr20_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr21_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr21_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr22_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr22_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr23_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xea, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr23_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr24_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr24_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr25_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr25_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr26_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr26_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr27_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr27_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr28_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr28_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr29_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr29_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevcntr30_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xeb, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevcntr30_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmccfiltr_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmccfiltr_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper0_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper0_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper1_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper1_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper2_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper2_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper3_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper3_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper4_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper4_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper5_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper5_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper6_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper6_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper7_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xec, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper7_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper8_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper8_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper9_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper9_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper10_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper10_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper11_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper11_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper12_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper12_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper13_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper13_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper14_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper14_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper15_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xed, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper15_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper16_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper16_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper17_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper17_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper18_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper18_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper19_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper19_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper20_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper20_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper21_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper21_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper22_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper22_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper23_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe9, 0xee, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper23_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper24_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x09, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper24_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper25_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x29, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper25_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper26_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x49, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper26_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper27_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x69, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper27_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper28_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x89, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper28_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper29_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xa9, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper29_el0",
			wantErr: false,
		},
		{
			name: "mrs	x9, pmevtyper30_el0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc9, 0xef, 0x3b, 0xd5}),
				address:          0,
			},
			want: "mrs	x9, pmevtyper30_el0",
			wantErr: false,
		},
		{
			name: "mrs	x12, s3_7_c15_c1_5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xac, 0xf1, 0x3f, 0xd5}),
				address:          0,
			},
			want: "mrs	x12, s3_7_c15_c1_5",
			wantErr: false,
		},
		{
			name: "mrs	x13, s3_2_c11_c15_7",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xed, 0xbf, 0x3a, 0xd5}),
				address:          0,
			},
			want: "mrs	x13, s3_2_c11_c15_7",
			wantErr: false,
		},
		{
			name: "sysl	x14, #3, c9, c2, #1",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x2e, 0x92, 0x2b, 0xd5}),
				address:          0,
			},
			want: "sysl	x14, #3, c9, c2, #1",
			wantErr: false,
		},
		{
			name: "msr	s3_0_c15_c0_0, x12",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x0c, 0xf0, 0x18, 0xd5}),
				address:          0,
			},
			want: "msr	s3_0_c15_c0_0, x12",
			wantErr: false,
		},
		{
			name: "msr	s3_7_c11_c13_7, x5",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe5, 0xbd, 0x1f, 0xd5}),
				address:          0,
			},
			want: "msr	s3_7_c11_c13_7, x5",
			wantErr: false,
		},
		{
			name: "sys	#3, c9, c2, #1, x4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x24, 0x92, 0x0b, 0xd5}),
				address:          0,
			},
			want: "sys	#3, c9, c2, #1, x4",
			wantErr: false,
		},
		{
			name: "b	#4",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x01, 0x00, 0x00, 0x14}),
				address:          0,
			},
			want: "b	#4",
			wantErr: false,
		},
		{
			name: "bl	#0",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x00, 0x94}),
				address:          0,
			},
			want: "bl	#0",
			wantErr: false,
		},
		{
			name: "b	#134217724",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xff, 0xff, 0xff, 0x15}),
				address:          0,
			},
			want: "b	#134217724",
			wantErr: false,
		},
		{
			name: "bl	#-134217728",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x00, 0x00, 0x00, 0x96}),
				address:          0,
			},
			want: "bl	#-134217728",
			wantErr: false,
		},
		{
			name: "br	x20",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x80, 0x02, 0x1f, 0xd6}),
				address:          0,
			},
			want: "br	x20",
			wantErr: false,
		},
		{
			name: "blr	xzr",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0x3f, 0xd6}),
				address:          0,
			},
			want: "blr	xzr",
			wantErr: false,
		},
		{
			name: "ret	x10",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0x40, 0x01, 0x5f, 0xd6}),
				address:          0,
			},
			want: "ret	x10",
			wantErr: false,
		},
		{
			name: "ret",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xc0, 0x03, 0x5f, 0xd6}),
				address:          0,
			},
			want:    "ret",
			wantErr: false,
		},
		{
			name: "eret",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0x9f, 0xd6}),
				address:          0,
			},
			want:    "eret",
			wantErr: false,
		},
		{
			name: "drps",
			args: args{
				instructionValue: binary.LittleEndian.Uint32([]byte{0xe0, 0x03, 0xbf, 0xd6}),
				address:          0,
			},
			want:    "drps",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decompose(tt.args.instructionValue, tt.args.address)
			if (err != nil) != tt.wantErr {
				fmt.Printf("want: %s\n", tt.want)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				t.Errorf("disassemble() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decOut, _ := got.disassemble(true)
			hexout, _ := got.disassemble(false)
			if !reflect.DeepEqual(decOut, strings.ToLower(tt.want)) && !reflect.DeepEqual(hexout, strings.ToLower(tt.want)) {
				fmt.Printf("want: %s\n", tt.want)
				fmt.Printf("got:  %s\n", decOut)
				fmt.Printf("got:  %s (hex)\n", hexout)
				got, _ = decompose(tt.args.instructionValue, tt.args.address)
				decOut, _ := got.disassemble(true)
				t.Errorf("disassemble(dec) = %v, want %v", decOut, tt.want)
			}
		})
	}
}
