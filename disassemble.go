package arm64

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *Instruction) disassemble(decimalImm bool) (string, error) {

	if i.operation == ARM64_UNDEFINED || i.operation == AMD64_END_TYPE {
		return "", fmt.Errorf("failed to disassemble operation")
	}

	for idx, operand := range i.operands {
		switch operand.OpClass {
		case FIMM32:
			fallthrough
		case IMM32:
			fallthrough
		case IMM64:
			fallthrough
		case LABEL:
			if err := operand.getShiftedImmediate(decimalImm); err != nil {
				return "", fmt.Errorf("failed to disassemble operation: %v", err)
			}
			i.operands[idx].strRepr = operand.strRepr
			break
		case REG:
			if err := operand.getRegister(0, decimalImm); err != nil {
				return "", fmt.Errorf("failed to disassemble operation: %v", err)
			}
			i.operands[idx].strRepr = operand.strRepr
			break
		case SYS_REG:
			i.operands[idx].strRepr = SystemReg(operand.Reg[0]).String()
			break
		case MULTI_REG:
			if err := operand.getMultiregOperand(decimalImm); err != nil {
				return "", fmt.Errorf("failed to disassemble operation: %v", err)
			}
			i.operands[idx].strRepr = operand.strRepr
			break
		case IMPLEMENTATION_SPECIFIC:
			if err := operand.getImplementationSpecific(); err != nil {
				return "", fmt.Errorf("failed to disassemble operation: %v", err)
			}
			i.operands[idx].strRepr = operand.strRepr
			break
		case MEM_REG:
			fallthrough
		case MEM_OFFSET:
			fallthrough
		case MEM_EXTENDED:
			fallthrough
		case MEM_PRE_IDX:
			fallthrough
		case MEM_POST_IDX:
			if err := operand.getMemoryOperand(decimalImm); err != nil {
				return "", fmt.Errorf("failed to disassemble operation: %v", err)
			}
			i.operands[idx].strRepr = operand.strRepr
			break
		case CONDITION:
			i.operands[idx].strRepr = Condition(operand.Reg[0]).String()
			break
		case NONE:
			break
		}

		if operand.OpClass != NONE {
			if idx == 0 {
				// i.operands[idx].strRepr = fmt.Sprintf(" %s", i.operands[idx].strRepr)
				i.operands[idx].strRepr = fmt.Sprintf("\t%s", i.operands[idx].strRepr)
			} else {
				i.operands[idx].strRepr = fmt.Sprintf(", %s", i.operands[idx].strRepr)
			}
		}
	}

	return fmt.Sprintf("%s%s%s%s%s%s",
		i.operation,
		i.operands[0],
		i.operands[1],
		i.operands[2],
		i.operands[3],
		i.operands[4]), nil
}

func (op *InstructionOperand) getShiftedImmediate(decimalImm bool) error {
	var shiftBuff string
	var immBuff string
	var sign string

	if op == nil {
		return failedToDisassembleOperand
	}

	if op.SignedImm == 1 || int64(op.Immediate) < 0 {
		sign = "-"
	}
	if op.ShiftType != SHIFT_NONE {
		if op.ShiftValueUsed != 0 || op.ShiftType != SHIFT_LSL {
			// if op.ShiftValueUsed != 0 {
			if decimalImm {
				immBuff = fmt.Sprintf(" #%d", op.ShiftValue)
			} else {
				immBuff = fmt.Sprintf(" #%#x", op.ShiftValue)
			}
		}
		shiftBuff = fmt.Sprintf(", %s%s", op.ShiftType, immBuff)
	}
	if op.OpClass == FIMM32 {
		if op.Immediate == 0 {
			shiftBuff = fmt.Sprintf("#%.1f%s", ieee754(op.Immediate).Float(), shiftBuff)
		} else {
			shiftBuff = fmt.Sprintf("#%.8f%s", ieee754(op.Immediate).Float(), shiftBuff)
		}

	} else if op.OpClass == IMM32 {
		if op.SignedImm == 1 || int32(op.Immediate) < 0 { // TODO this is gross
			if decimalImm {
				shiftBuff = fmt.Sprintf("#%d%s", int32(op.Immediate), shiftBuff)
			} else {
				shiftBuff = fmt.Sprintf("#%#x%s", int32(op.Immediate), shiftBuff)
			}
		} else {
			if decimalImm {
				shiftBuff = fmt.Sprintf("#%s%d%s", sign, op.Immediate, shiftBuff)
			} else {
				shiftBuff = fmt.Sprintf("#%s%#x%s", sign, op.Immediate, shiftBuff)
			}
		}
	} else {
		if op.SignedImm == 1 && int64(op.Immediate) < 0 { // TODO this is gross
			if decimalImm {
				shiftBuff = fmt.Sprintf("#%d%s", int64(op.Immediate), shiftBuff)
			} else {
				shiftBuff = fmt.Sprintf("#%#x%s", int64(op.Immediate), shiftBuff)
			}
		} else {
			if decimalImm {
				shiftBuff = fmt.Sprintf("#%s%d%s", sign, op.Immediate, shiftBuff)
			} else {
				shiftBuff = fmt.Sprintf("#%s%#x%s", sign, op.Immediate, shiftBuff)
			}
		}
	}

	op.strRepr = shiftBuff

	return nil
}

// func (op *InstructionOperand) getShiftedImmediate() error {
// 	var shiftBuff string
// 	var immBuff string
// 	var sign string

// 	if op == nil {
// 		return failedToDisassembleOperand
// 	}

// 	imm := int64(op.Immediate)

// 	if op.SignedImm == 1 && imm < 0 {
// 		sign = "-"
// 		imm = -imm
// 	}

// 	if op.ShiftType != SHIFT_NONE {
// 		if op.ShiftValueUsed != 0 {
// 			immBuff = fmt.Sprintf(" #%#x", op.ShiftValue)
// 		}
// 		shiftBuff = fmt.Sprintf(", %s%s", ShiftType(op.ShiftType), immBuff)
// 	}
// 	if op.OpClass == FIMM32 {
// 		op.strRepr = fmt.Sprintf("#%f%s", float64(op.Immediate), shiftBuff)
// 	} else if op.OpClass == IMM32 {
// 		op.strRepr = fmt.Sprintf("#%s%#x%s", sign, uint32(imm), shiftBuff)

// 	} else {
// 		op.strRepr = fmt.Sprintf("#%s%#x%s", sign, imm, shiftBuff)
// 	}

// 	return nil
// }

func (op *InstructionOperand) getRegister(registerNumber int, decimalImm bool) error {
	var scale string

	if op.Scale != 0 {
		scale = fmt.Sprintf("[%d]", 0x7fffffff&op.Scale)
	}

	if op.OpClass == SYS_REG {
		op.strRepr = fmt.Sprintf("%s", SystemReg(op.Reg[registerNumber]))
		return nil
	} else if op.OpClass != REG && op.OpClass != MULTI_REG {
		return operandIsNotRegister
	}

	if op.ShiftType != SHIFT_NONE {
		return op.getShiftedRegister(registerNumber, decimalImm)
	} else if op.ElementSize == 0 {
		op.strRepr = fmt.Sprintf("%s", Register(op.Reg[registerNumber]))
		if !decimalImm {
			if strings.HasPrefix(op.strRepr, "#") {
				i, err := strconv.Atoi(strings.TrimPrefix(op.strRepr, "#"))
				if err != nil {
					return fmt.Errorf("getRegister: failed to convert number register from str to int")
				}
				op.strRepr = fmt.Sprintf("#%#x", i)
			}
		}
		return nil
	}

	var elementSize string
	switch op.ElementSize {
	case 1:
		elementSize = "b"
		break
	case 2:
		elementSize = "h"
		break
	case 4:
		elementSize = "s"
		break
	case 8:
		elementSize = "d"
		break
	case 16:
		elementSize = "q"
		break
	default:
		return failedToDisassembleRegister
	}

	if op.DataSize != 0 {
		if registerNumber > 3 || (op.DataSize != 1 && op.DataSize != 2 && op.DataSize != 4 && op.DataSize != 8 && op.DataSize != 16) {
			return failedToDisassembleRegister
		}
		op.strRepr = fmt.Sprintf("%s.%d%s%s", Register(op.Reg[registerNumber]), op.DataSize, elementSize, scale)
	} else {
		if registerNumber > 3 {
			return failedToDisassembleRegister
		}
		op.strRepr = fmt.Sprintf("%s.%s%s", Register(op.Reg[registerNumber]), elementSize, scale)
	}

	return nil
}

func (op *InstructionOperand) getShiftedRegister(registerNumber int, decimalImm bool) error {
	var immBuff string
	var shiftBuff string

	reg := Register(op.Reg[registerNumber])
	if reg == REG_NONE {
		return failedToDisassembleRegister
	}
	if op.ShiftType != SHIFT_NONE {
		// if op.ShiftValueUsed != 0 || op.ShiftType != SHIFT_LSL {
		if op.ShiftValueUsed != 0 {
			if decimalImm {
				immBuff = fmt.Sprintf(" #%d", op.ShiftValue)
			} else {
				immBuff = fmt.Sprintf(" #%#x", op.ShiftValue)
			}
		}
		shiftBuff = fmt.Sprintf(", %s%s", ShiftType(op.ShiftType), immBuff)
	}
	op.strRepr = fmt.Sprintf("%s%s", reg, shiftBuff)
	return nil
}

func (op *InstructionOperand) getMultiregOperand(decimalImm bool) error {
	var indexBuff string
	var registers []Register
	var elementCount int

	for _, opReg := range op.Reg {
		if Register(opReg) != REG_NONE {
			if err := op.getRegister(elementCount, decimalImm); err != nil {
				return err
			}
			registers = append(registers, Register(opReg))
			elementCount++
		}
	}

	if op.Index != 0 {
		indexBuff = fmt.Sprintf("[%d]", op.Index)
	}

	switch elementCount {
	case 1:
		op.strRepr = fmt.Sprintf("{%s}%s", registers[0], indexBuff)
		break
	case 2:
		op.strRepr = fmt.Sprintf("{%s, %s}%s", registers[0], registers[1], indexBuff)
		break
	case 3:
		op.strRepr = fmt.Sprintf("{%s, %s, %s}%s", registers[0], registers[1], registers[2], indexBuff)
		break
	case 4:
		op.strRepr = fmt.Sprintf("{%s, %s, %s, %s}%s", registers[0], registers[1], registers[2], registers[3], indexBuff)
		break
	default:
		return failedToDisassembleOperand
	}

	return nil
}

func (op *InstructionOperand) getImplementationSpecific() error {
	op.strRepr = fmt.Sprintf("s%d_%d_c%d_c%d_%d", op.Reg[0], op.Reg[1], op.Reg[2], op.Reg[3], op.Reg[4])
	return nil
}

func (op *InstructionOperand) getMemoryOperand(decimalImm bool) error {
	var immBuff string
	var extendBuff string
	var paramBuff string
	var outBuffer string
	var sign string

	reg1 := Register(op.Reg[0])
	reg2 := Register(op.Reg[1])

	if op == nil {
		return failedToDisassembleOperand
	}

	imm := int64(op.Immediate)

	if op.SignedImm == 1 && imm < 0 {
		sign = "-"
		imm = -imm
	}

	switch op.OpClass {
	case MEM_REG:
		outBuffer = fmt.Sprintf("[%s]", Register(op.Reg[0]))
		break
	case MEM_PRE_IDX:
		if decimalImm {
			outBuffer = fmt.Sprintf("[%s, #%s%d]!", Register(op.Reg[0]), sign, uint64(imm))
		} else {
			outBuffer = fmt.Sprintf("[%s, #%s%#x]!", Register(op.Reg[0]), sign, uint64(imm))
		}
		break
	case MEM_POST_IDX: // [<reg>], <reg|imm>
		if Register(op.Reg[1]) != REG_NONE {
			paramBuff = fmt.Sprintf(", %s", Register(op.Reg[1]))
		} else {
			if decimalImm {
				paramBuff = fmt.Sprintf(", #%s%d", sign, uint64(imm))
			} else {
				paramBuff = fmt.Sprintf(", #%s%#x", sign, uint64(imm))
			}
		}
		outBuffer = fmt.Sprintf("[%s]%s", Register(op.Reg[0]), paramBuff)
		break
	case MEM_OFFSET: // [<reg> optional(imm)]
		if imm != 0 {
			if decimalImm {
				immBuff = fmt.Sprintf(", #%s%d", sign, uint64(imm))
			} else {
				immBuff = fmt.Sprintf(", #%s%#x", sign, uint64(imm))
			}
		}
		outBuffer = fmt.Sprintf("[%s%s]", Register(op.Reg[0]), immBuff)
		break
	case MEM_EXTENDED:
		if reg1 == REG_NONE || reg2 == REG_NONE {
			return failedToDisassembleOperand
		}
		if op.ShiftValue != 0 || op.ShiftType != SHIFT_LSL {
			if decimalImm {
				immBuff = fmt.Sprintf(" #%d", op.ShiftValue)
			} else {
				immBuff = fmt.Sprintf(" #%#x", op.ShiftValue)
			}
		}
		if op.ShiftType != SHIFT_NONE {
			extendBuff = fmt.Sprintf(", %s%s", ShiftType(op.ShiftType), immBuff)
		}
		outBuffer = fmt.Sprintf("[%s, %s%s]", reg1, reg2, extendBuff)
		break
	default:
		return notMemoryOperand
	}

	op.strRepr = outBuffer

	return nil
}
