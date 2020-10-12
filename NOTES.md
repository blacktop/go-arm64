# NOTES

## `objdump`

```bash
$ xcrun objdump --macho --demangle --triple=arm64e-apple-ios -mcpu=apple-latest -mattr=+mte,+v8.5a --print-imm-hex -dis-symname=<SYMBOL> -d BINARY |  bat -l s --tabs 0 -p --theme Nord
```
