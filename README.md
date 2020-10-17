# go-arm64 [WIP] 🚧

[![Go](https://github.com/blacktop/go-arm64/workflows/Go/badge.svg)](https://github.com/blacktop/go-arm64/actions) [![PkgGoDev](https://pkg.go.dev/badge/blacktop/go-arm64)](https://pkg.go.dev/github.com/blacktop/go-arm64) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

> Aarch64 architecture disassembler

---

## Install

```bash
$ go get github.com/blacktop/go-arm64
```

## Getting Started

```go
package main

import "github.com/blacktop/go-arm64"

func main() {
    f, err := os.Open("/path/to/hello-mte")
    if err != nil {
        panic(err)
    }

    for i := range arm64.Disassemble(f, int64(0x100007e58)) {
        if i.Error != nil{
            fmt.Println(i.StrRepr)
        }
    }
}
```

```nasm
_test:
0x100007e58:	7f 23 03 d5 	pacibsp
0x100007e5c:	ff c3 02 d1 	sub	sp, sp, #0xb0
0x100007e60:	fd 7b 0a a9 	stp	x29, x30, [sp, #0xa0]
0x100007e64:	fd 83 02 91 	add	x29, sp, #0xa0
0x100007e68:	08 00 00 b0 	adrp	x8, #0x100008000
0x100007e6c:	08 09 40 f9 	ldr	x8, [x8, #0x10]
0x100007e70:	08 01 40 f9 	ldr	x8, [x8]
0x100007e74:	a8 83 1f f8 	stur	x8, [x29, #-0x8]
0x100007e78:	08 00 80 d2 	mov	x8, #0x0
0x100007e7c:	e8 13 c8 9a 	irg	x8, sp, x8
0x100007e80:	08 01 81 91 	addg	x8, x8, #0x10, #0x0
0x100007e84:	08 69 a0 d9 	st2g	x8, [x8, #0x60]
0x100007e88:	08 49 a0 d9 	st2g	x8, [x8, #0x40]
0x100007e8c:	08 29 a0 d9 	st2g	x8, [x8, #0x20]
0x100007e90:	08 09 a0 d9 	st2g	x8, [x8]
0x100007e94:	09 00 80 52 	mov	w9, #0x0
0x100007e98:	09 01 04 39 	strb	w9, [x8, #0x100]
0x100007e9c:	09 fd 43 39 	ldrb	w9, [x8, #0xff]
0x100007ea0:	ff 7b a0 d9 	st2g	sp, [sp, #0x70]
0x100007ea4:	ff 5b a0 d9 	st2g	sp, [sp, #0x50]
0x100007ea8:	ff 3b a0 d9 	st2g	sp, [sp, #0x30]
0x100007eac:	ff 1b a0 d9 	st2g	sp, [sp, #0x10]
0x100007eb0:	08 00 00 b0 	adrp	x8, #0x100008000
0x100007eb4:	08 09 40 f9 	ldr	x8, [x8, #0x10]
0x100007eb8:	08 01 40 f9 	ldr	x8, [x8]
0x100007ebc:	aa 83 5f f8 	ldur	x10, [x29, #-0x8]
0x100007ec0:	08 01 0a eb 	subs	x8, x8, x10
0x100007ec4:	e9 0f 00 b9 	str	w9, [sp, #0xc]
0x100007ec8:	e1 00 00 54 	b.ne	#0x100007ee4
0x100007ecc:	01 00 00 14 	b	#0x100007ed0
0x100007ed0:	e8 0f 40 b9 	ldr	w8, [sp, #0xc]
0x100007ed4:	00 1d 00 12 	and	w0, w8, #0xff
0x100007ed8:	fd 7b 4a a9 	ldp	x29, x30, [sp, #0xa0]
0x100007edc:	ff c3 02 91 	add	sp, sp, #0xb0
0x100007ee0:	ff 0f 5f d6 	retab	xzr
0x100007ee4:	2c 00 00 94 	bl	#0x100007f94
0x100007ee8:	20 00 20 d4 	brk	#0x1
```

## TODO

- [ ] fix 🐛🐛🐛
- [x] add option for dec/hex immediates
- [x] display opcodes like `7f 23 03 d5`
- [ ] benchmarks 🏃‍♂️💨

## Credit

This is a complete Go re-write of [Vector35/arch-arm64](https://github.com/Vector35/arch-arm64/tree/master/disassembler)

## License

MIT Copyright (c) 2020 **blacktop**
