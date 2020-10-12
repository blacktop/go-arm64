# go-arm64 [WIP] 🚧

[![Go](https://github.com/blacktop/go-arm64/workflows/Go/badge.svg)](https://github.com/blacktop/go-arm64/actions) [![PkgGoDev](https://pkg.go.dev/badge/blacktop/go-arm64)](https://pkg.go.dev/blacktop/go-arm64) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

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
    f, err := os.Open("/path/to/binary")
    if err != nil {
        panic(err)
    }

    for i := range arm64.Disassemble(f, 0) {
        if i.Error != nil{
            fmt.Println(i.StrRepr)
        }
    }
}
```

## Credit

This is a complete Go re-write of [Vector35/arch-arm64](https://github.com/Vector35/arch-arm64/tree/master/disassembler)

## License

MIT Copyright (c) 2020 **blacktop**
