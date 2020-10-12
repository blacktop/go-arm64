package arm64

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Disassemble will output the disassembly of the data of a given io.ReadSeeker
func Disassemble(r io.ReadSeeker, startAddr int64) (<-chan string, error) {

	out := make(chan string)

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

			instruction, err := i.disassemble()
			if err != nil {
				return fmt.Errorf("failed to disassemble instruction: 0x%08x; %v", instrValue, err)
			}

			out <- fmt.Sprintf("0x%08x:\t%08x \t%s\n", addr, instrValue, instruction)
		}

		close(out)

		return nil
	}()

	return out, nil
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
