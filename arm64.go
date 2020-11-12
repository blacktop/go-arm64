package arm64

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Options Disassemble options
type Options struct {
	StartAddress int64
	DecimalImm   bool
}

// Result Disassemble instruction result
type Result struct {
	StrRepr     string
	Instruction *Instruction
	Error       error
}

// Disassemble will output the disassembly of the data of a given io.ReadSeeker
func Disassemble(r io.ReadSeeker, options Options) <-chan Result {

	out := make(chan Result)

	go func() {
		var err error
		var instrValue uint32
		for {
			addr, _ := r.Seek(0, io.SeekCurrent)

			err = binary.Read(r, binary.LittleEndian, &instrValue)

			if err == io.EOF {
				close(out)
				break
			}

			if err != nil {
				out <- Result{
					Error: fmt.Errorf("failed to read instruction: %v", err),
				}
				close(out)
				break
			}

			if options.StartAddress != 0 {
				addr += options.StartAddress
			} else {
				addr = 0
			}

			i, err := decompose(instrValue, uint64(addr))
			if err != nil {
				if err == failedToDecodeInstruction {
					out <- Result{
						StrRepr: fmt.Sprintf("%#08x:  %s\t<unknown>", uint64(addr), getOpCodeByteString(instrValue)),
						Error:   fmt.Errorf("failed to decode instruction: 0x%08x; %v", instrValue, err),
					}
				} else {
					out <- Result{
						StrRepr: fmt.Sprintf("%#08x:  %s\t💥 ERROR 💥", uint64(addr), getOpCodeByteString(instrValue)),
						Error:   fmt.Errorf("failed to decode instruction: 0x%08x; %v", instrValue, err),
					}
				}
				continue
			}

			instruction, err := i.disassemble(options.DecimalImm)
			if err != nil {
				out <- Result{
					StrRepr: fmt.Sprintf("%#08x:  %s\t<unknown>", uint64(addr), getOpCodeByteString(instrValue)),
					Error:   fmt.Errorf("failed to disassemble instruction: 0x%08x; %v", instrValue, err),
				}
				continue
			}

			out <- Result{
				StrRepr:     fmt.Sprintf("%#08x:  %s\t%s", uint64(addr), getOpCodeByteString(instrValue), instruction),
				Instruction: i,
				Error:       nil,
			}
		}
		return
	}()

	return out
}

// Instructions will output the decoded instructions of the data of a given io.ReadSeeker
func Instructions(r io.ReadSeeker, startAddr int64) (<-chan *Instruction, error) {

	out := make(chan *Instruction)

	go func() error {
		var instrValue uint32
		for {
			addr, _ := r.Seek(0, io.SeekCurrent)

			err := binary.Read(r, binary.LittleEndian, &instrValue)

			if err == io.EOF {
				break
			}

			if err != nil {
				return fmt.Errorf("failed to read instruction: %v", err)
			}

			if startAddr > 0 {
				addr += startAddr
			} else {
				addr = 0
			}

			i, err := decompose(instrValue, 0)
			if err != nil {
				return fmt.Errorf("failed to decompose instruction: 0x%08x; %v", instrValue, err)
			}

			out <- i
		}

		close(out)

		return nil
	}()

	return out, nil
}
