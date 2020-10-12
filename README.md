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
0x100007e58:    d503237f        pacibsp
0x100007e5c:    d102c3ff        sub     sp, sp, #0xb0
0x100007e60:    a90a7bfd        stp     x29, x30, [sp, #0xa0]
0x100007e64:    910283fd        add     x29, sp, #0xa0
0x100007e68:    b0000008        adrp    x8, #0x1000
0x100007e6c:    f9400908        ldr     x8, [x8, #0x10]
0x100007e70:    f9400108        ldr     x8, [x8, #0x0]
0x100007e74:    f81f83a8        stur    x8, [x29, #-0x8]
0x100007e78:    d2800008        mov     x8, #0x0
0x100007e7c:    9ac813e8        irg     x8, sp, x8
0x100007e80:    91810108        addg    x8, x8, #0x10, #0x0
0x100007e84:    d9a06908        st2g    x8, [x8, #0x60]
0x100007e88:    d9a04908        st2g    x8, [x8, #0x40]
0x100007e8c:    d9a02908        st2g    x8, [x8, #0x20]
0x100007e90:    d9a00908        st2g    x8, [x8, #0x0]
0x100007e94:    52800009        mov     w9, #0x0
0x100007e98:    39040109        strb    w9, [x8, #0x100]
0x100007e9c:    3943fd09        ldrb    w9, [x8, #0xff]
0x100007ea0:    d9a07bff        st2g    sp, [sp, #0x70]
0x100007ea4:    d9a05bff        st2g    sp, [sp, #0x50]
0x100007ea8:    d9a03bff        st2g    sp, [sp, #0x30]
0x100007eac:    d9a01bff        st2g    sp, [sp, #0x10]
0x100007eb0:    b0000008        adrp    x8, #0x1000
0x100007eb4:    f9400908        ldr     x8, [x8, #0x10]
0x100007eb8:    f9400108        ldr     x8, [x8, #0x0]
0x100007ebc:    f85f83aa        ldur    x10, [x29, #-0x8]
0x100007ec0:    eb0a0108        subs    x8, x8, x10
0x100007ec4:    b9000fe9        str     w9, [sp, #0xc]
0x100007ec8:    540000e1        b.ne    #0x1c
0x100007ecc:    14000001        b       #0x4
0x100007ed0:    b9400fe8        ldr     w8, [sp, #0xc]
0x100007ed4:    12001d00        and     w0, w8, #0xff
0x100007ed8:    a94a7bfd        ldp     x29, x30, [sp, #0xa0]
0x100007edc:    9102c3ff        add     sp, sp, #0xb0
0x100007ee0:    d65f0fff        retab   xzr
0x100007ee4:    9400002c        bl      #0xb0
0x100007ee8:    d4200020        brk     #0x1
```

## TODO

- [ ] fix 🐛🐛🐛
- [ ] add option for dec/hex immediates
- [ ] display opcodes like `7f 23 03 d5`
- [ ] benchmarks 🏃‍♂️💨

## Credit

This is a complete Go re-write of [Vector35/arch-arm64](https://github.com/Vector35/arch-arm64/tree/master/disassembler)

## License

MIT Copyright (c) 2020 **blacktop**
