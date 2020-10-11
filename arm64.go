package arm64

import "fmt"

const (
	REG_W_BASE  = 0
	REG_X_BASE  = 1
	REG_V_BASE  = 2
	REG_B_BASE  = 3
	REG_H_BASE  = 4
	REG_S_BASE  = 5
	REG_D_BASE  = 6
	REG_Q_BASE  = 7
	REG_PF_BASE = 8

	REGSET_SP = 0
	REGSET_ZR = 1
)

func getRegistery() {

}

func (i *Instruction) decompose_load_register_literal() (*Instruction, error) {
	/* C4.3.5 Load register (literal)
	 * LDR <Wt>, <label>
	 * LDR <Xt>, <label>
	 * LDR <St>, <label>
	 * LDR <Dt>, <label>
	 * LDR <Qt>, <label>
	 * LDRSW <Xt>, <label>
	 * PRFM <prfop>, <label>
	 */

	decode := LoadRegisterLiteral(i.raw)

	type option struct {
		operation Operation
		regBase   uint32
		signedImm uint32
	}
	var operand = [2][4]option{
		{
			{LDR, REG_W_BASE, 0},
			{LDR, REG_X_BASE, 0},
			{LDRSW, REG_X_BASE, 1},
			{PRFM, REG_W_BASE, 0},
		}, {
			{LDR, REG_S_BASE, 0},
			{LDR, REG_D_BASE, 0},
			{LDR, REG_Q_BASE, 0},
			{ARM64_UNDEFINED, 0, 0},
		},
	}
	// fmt.Println(decode)
	op := &operand[decode.V()][decode.Opc()]
	i.operation = op.operation
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = uint32(regMap[REGSET_ZR][op.regBase][decode.Rt()])

	i.operands[1].OpClass = LABEL
	i.operands[1].SignedImm = op.signedImm
	if op.signedImm > 0 {
		i.operands[1].Immediate = i.address - uint64(decode.Imm()<<2)
	}
	i.operands[1].Immediate = i.address + uint64(decode.Imm()<<2)

	if i.operation == ARM64_UNDEFINED {
		return i, fmt.Errorf("failed to decode load register literal")
	}

	return i, nil
}

func decompose(instructionValue uint32, address uint64) (*Instruction, error) {

	instruction := &Instruction{
		raw:       instructionValue,
		address:   address,
		operation: ARM64_UNDEFINED,
	}

	switch ExtractBits(instructionValue, 25, 4) {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		instruction.group = GROUP_UNALLOCATED
		return instruction, nil
	case 8:
		fallthrough
	case 9:
		instruction.group = GROUP_DATA_PROCESSING_IMM
		switch ExtractBits(instructionValue, 23, 3) {
		case 0:
			fallthrough
		case 1:
			// return aarch64_decompose_pc_rel_addr(instructionValue, instruction, address)
		case 2:
			// return aarch64_decompose_add_sub_imm(instructionValue, instruction)
		case 3:
			// return aarch64_decompose_add_sub_imm_tags(instructionValue, instruction)
		case 4:
			// return aarch64_decompose_logical_imm(instructionValue, instruction)
		case 5:
			// return aarch64_decompose_move_wide_imm(instructionValue, instruction)
		case 6:
			// return aarch64_decompose_bitfield(instructionValue, instruction)
		case 7:
			// return aarch64_decompose_extract(instructionValue, instruction)
		}
		break
	case 10:
		fallthrough
	case 11:
		instruction.group = GROUP_BRANCH_EXCEPTION_SYSTEM
		switch ExtractBits(instructionValue, 25, 7) {
		case 0xa:
			fallthrough
		case 0xb:
			fallthrough
		case 0x4a:
			fallthrough
		case 0x4b:
			// return aarch64_decompose_unconditional_branch(instructionValue, instruction, address)
		case 0x1a:
			fallthrough
		case 0x5a:
			// return aarch64_decompose_compare_branch_imm(instructionValue, instruction, address)
		case 0x1b:
			fallthrough
		case 0x5b:
			// return aarch64_decompose_test_branch_imm(instructionValue, instruction, address)
		case 0x2a:
			// return aarch64_decompose_conditional_branch(instructionValue, instruction, address)
		case 0x6a:
			if ExtractBits(instructionValue, 24, 1) == 0 {
				// return aarch64_decompose_exception_generation(instructionValue, instruction)
			} else if ExtractBits(instructionValue, 22, 3) == 4 {
				// return aarch64_decompose_system(instructionValue, instruction)
			}
			return instruction, nil // TODO error  ?
		case 0x6b:
			// return aarch64_decompose_unconditional_branch_reg(instructionValue, instruction)
		default:
			return instruction, nil // TODO: (error) shouldn't be able to get here
		}
		break
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 12:
		fallthrough
	case 14:
		{
			instruction.group = GROUP_LOAD_STORE

			op0 := ExtractBits(instructionValue, 28, 4)
			op1 := ExtractBits(instructionValue, 26, 1)
			op2 := ExtractBits(instructionValue, 23, 2)
			op3 := ExtractBits(instructionValue, 16, 6)
			op4 := ExtractBits(instructionValue, 10, 2)

			if (op0&0b1011) == 0 && (op1 == 1) {
				if op2 == 0 && op3 == 0 {
					// return aarch64_decompose_simd_load_store_multiple(instructionValue, instruction)
				}
				if op2 == 1 && (op3>>5) == 0 {
					// return aarch64_decompose_simd_load_store_multiple_post_idx(instructionValue, instruction)
				}
				if op2 == 2 && (op3&0x1f) == 0 {
					// return aarch64_decompose_simd_load_store_single(instructionValue, instruction)
				}
				if op2 == 3 {
					// return aarch64_decompose_simd_load_store_single_post_idx(instructionValue, instruction)
				}
			}

			if op0 == 0x0d {
				// return aarch64_decompose_load_store_mem_tags(instructionValue, instruction)
			}

			if (op0&3) == 0 && op1 == 0 && (op2>>1) == 0 {
				// return aarch64_decompose_load_store_exclusive(instructionValue, instruction)
			}
			if (op0&3) == 1 && (op2>>1) == 0 {
				return instruction.decompose_load_register_literal()
			}

			if (op0 & 3) == 2 {
				if op2 == 0 {
					// return aarch64_decompose_load_store_no_allocate_pair_offset(instructionValue, instruction)
				}
				if op2 == 1 {
					// return aarch64_decompose_load_store_reg_pair_post_idx(instructionValue, instruction)
				}
				if op2 == 2 {
					// return aarch64_decompose_load_store_reg_pair_offset(instructionValue, instruction)
				}
				if op2 == 3 {
					// return aarch64_decompose_load_store_reg_pair_pre_idx(instructionValue, instruction)
				}
			}

			if (op0 & 3) == 3 {
				if (op2 >> 1) == 0 {
					if (op3 >> 5) == 0 {
						if op4 == 0 {
							// return aarch64_decompose_load_store_reg_unscalled_imm(instructionValue, instruction)
						}
						if op4 == 1 {
							// return aarch64_decompose_load_store_imm_post_idx(instructionValue, instruction)
						}
						if op4 == 2 {
							// return aarch64_decompose_load_store_reg_unpriv(instructionValue, instruction)
						}
						if op4 == 3 {
							// return aarch64_decompose_load_store_reg_imm_pre_idx(instructionValue, instruction)
						}
					}
					if (op3 >> 5) == 1 {
						// if(op4==0) return aarch64_decompose_atomic_memory_ops(instructionValue, instruction) // TODO: remove ?
						if op4 == 2 {
							// return aarch64_decompose_load_store_reg_reg_offset(instructionValue, instruction)
						}
						if op4 == 1 || op4 == 3 {
							// return aarch64_decompose_load_store_pac(instructionValue, instruction)
						}
					}
				}
				// return aarch64_decompose_load_store_reg_unsigned_imm(instructionValue, instruction)
			}
			break
		}
	case 5:
	case 13:
		instruction.group = GROUP_DATA_PROCESSING_REG
		switch ExtractBits(instructionValue, 21, 8) {
		case 0x50:
			fallthrough
		case 0x51:
			fallthrough
		case 0x52:
			fallthrough
		case 0x53:
			fallthrough
		case 0x54:
			fallthrough
		case 0x55:
			fallthrough
		case 0x56:
			fallthrough
		case 0x57:
			// return aarch64_decompose_logical_shifted_reg(instructionValue, instruction)
		case 0x58:
			fallthrough
		case 0x5a:
			fallthrough
		case 0x5c:
			fallthrough
		case 0x5e:
			// return aarch64_decompose_add_sub_shifted_reg(instructionValue, instruction)
		case 0x59:
			fallthrough
		case 0x5b:
			fallthrough
		case 0x5d:
			fallthrough
		case 0x5f:
			// return aarch64_decompose_add_sub_extended_reg(instructionValue, instruction)
		case 0xd0:
			// return aarch64_decompose_add_sub_carry(instructionValue, instruction)
		case 0xd2:
			if ExtractBits(instructionValue, 11, 1) == 1 {
				// return aarch64_decompose_conditional_compare_imm(instructionValue, instruction)
			}
			// return aarch64_decompose_conditional_compare_reg(instructionValue, instruction)
		case 0xd4:
			// return aarch64_decompose_conditional_select(instructionValue, instruction)
		case 0xd8:
			fallthrough
		case 0xd9:
			fallthrough
		case 0xda:
			fallthrough
		case 0xdb:
			fallthrough
		case 0xdc:
			fallthrough
		case 0xdd:
			fallthrough
		case 0xde:
			fallthrough
		case 0xdf:
			// return aarch64_decompose_data_processing_3(instructionValue, instruction)
		case 0xd6:
			if ExtractBits(instructionValue, 30, 1) == 1 {
				// return aarch64_decompose_data_processing_1(instructionValue, instruction)
			}
			// return aarch64_decompose_data_processing_2(instructionValue, instruction)
		default:
			return instruction, nil // TODO: or error ?
		}
		break
	case 7:
	case 15:
		instruction.group = GROUP_DATA_PROCESSING_SIMD
		switch ExtractBits(instructionValue, 24, 8) {
		case 0x1e:
			fallthrough
		case 0x3e:
			fallthrough
		case 0x9e:
			fallthrough
		case 0xbe:
			if ExtractBits(instructionValue, 21, 1) == 0 {
				// return aarch64_decompose_fixed_floating_conversion(instructionValue, instruction)
			}
			switch ExtractBits(instructionValue, 10, 2) {
			case 0:
				if ExtractBits(instructionValue, 12, 1) == 1 {
					// return aarch64_decompose_floating_imm(instructionValue, instruction)
				} else if ExtractBits(instructionValue, 12, 2) == 2 {
					// return aarch64_decompose_floating_compare(instructionValue, instruction)
				} else if ExtractBits(instructionValue, 12, 3) == 4 {
					// return aarch64_decompose_floating_data_processing1(instructionValue, instruction)
				} else if ExtractBits(instructionValue, 12, 4) == 0 {
					// return aarch64_decompose_floating_integer_conversion(instructionValue, instruction)
				}
				break
			case 1:
				// return aarch64_decompose_floating_conditional_compare(instructionValue, instruction)
			case 2:
				// return aarch64_decompose_floating_data_processing2(instructionValue, instruction)
			case 3:
				// return aarch64_decompose_floating_cselect(instructionValue, instruction)
			}
			break
		case 0x1f:
			fallthrough
		case 0x3f:
			fallthrough
		case 0x9f:
			fallthrough
		case 0xbf:
			// return aarch64_decompose_floating_data_processing3(instructionValue, instruction)
		case 0x0e:
			fallthrough
		case 0x2e:
			fallthrough
		case 0x4e:
			fallthrough
		case 0x6e:
			if ExtractBits(instructionValue, 21, 1) == 1 {
				switch ExtractBits(instructionValue, 10, 2) {
				case 1:
				case 3:
					// return aarch64_decompose_simd_3_same(instructionValue, instruction)
				case 0:
					// return aarch64_decompose_simd_3_different(instructionValue, instruction)
				case 2:
					if ExtractBits(instructionValue, 17, 4) == 0 {
						// return aarch64_decompose_simd_2_reg_misc(instructionValue, instruction)
					} else if ExtractBits(instructionValue, 17, 4) == 8 {
						// return aarch64_decompose_simd_across_lanes(instructionValue, instruction)
					}
				}
			}
			if (instructionValue & 0x9fe08400) == 0x0e000400 {
				// return aarch64_decompose_simd_copy(instructionValue, instruction)
			}
			if (instructionValue & 0x003e0c00) == 0x00280800 {
				// return aarch64_decompose_cryptographic_aes(instructionValue, instruction)
			}
			if (instructionValue & 0xbf208c00) == 0x0e000000 {
				// return aarch64_decompose_simd_table_lookup(instructionValue, instruction)
			}
			if (instructionValue & 0xbf208c00) == 0x0e000800 {
				// return aarch64_decompose_simd_permute(instructionValue, instruction)
			}
			if (instructionValue & 0x00208400) == 0 {
				// return aarch64_decompose_simd_extract(instructionValue, instruction)
			}
			break
		case 0x0f:
			fallthrough
		case 0x2f:
			fallthrough
		case 0x4f:
			fallthrough
		case 0x6f:
			if ExtractBits(instructionValue, 10, 1) == 0 {
				// return aarch64_decompose_simd_vector_indexed_element(instructionValue, instruction)
			}
			if ExtractBits(instructionValue, 19, 5) == 0 {
				// return aarch64_decompose_simd_modified_imm(instructionValue, instruction)
			}
			// return aarch64_decompose_simd_shift_imm(instructionValue, instruction)
		case 0x5e:
			fallthrough
		case 0x7e:
			if ExtractBits(instructionValue, 21, 1) == 1 {
				switch ExtractBits(instructionValue, 10, 2) {
				case 1:
				case 3:
					// return aarch64_decompose_simd_scalar_3_same(instructionValue, instruction)
				case 0:
					// return aarch64_decompose_simd_scalar_3_different(instructionValue, instruction)
				case 2:
					if ExtractBits(instructionValue, 17, 4) == 0 {
						// return aarch64_decompose_simd_scalar_2_reg_misc(instructionValue, instruction)
					} else if ExtractBits(instructionValue, 17, 4) == 8 {
						// return aarch64_decompose_simd_scalar_pairwise(instructionValue, instruction)
					}
				}
			}
			if (instructionValue & 0xdfe08400) == 0x5e000400 {
				// return aarch64_decompose_simd_scalar_copy(instructionValue, instruction)
			} else if (instructionValue & 0xff208c00) == 0x5e000000 {
				// return aarch64_decompose_cryptographic_3_register_sha(instructionValue, instruction)
			} else if (instructionValue & 0xff3e0c00) == 0x5e280800 {
				// return aarch64_decompose_cryptographic_2_register_sha(instructionValue, instruction)
			}
			break
		case 0x5f:
			fallthrough
		case 0x7f:
			if ExtractBits(instructionValue, 10, 1) == 0 {
				// return aarch64_decompose_simd_scalar_indexed_element(instructionValue, instruction)
			} else if ExtractBits(instructionValue, 23, 1) == 0 && ExtractBits(instructionValue, 10, 1) == 1 {
				// return aarch64_decompose_simd_scalar_shift_imm(instructionValue, instruction)
			}
			break
		}
	default:
		break //should never get here
	}

	return instruction, nil
}

func (i *Instruction) disassemble() (string, error) {

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
			if err := operand.getShiftedImmediate(); err != nil {
				return "", fmt.Errorf("failed to disassemble operation: %v", err)
			}
			i.operands[idx].strRepr = operand.strRepr
			break
		case REG:
			i.operands[idx].strRepr = Register(operand.Reg[0]).String()
			// if (get_register(
			// 		&instruction->operands[i],
			// 		0,
			// 		tmpOperandString,
			// 		sizeof(tmpOperandString)) != DISASM_SUCCESS)
			// 	return "", fmt.Errorf("failed to disassemble operation")
			// operand = tmpOperandString;
			break
		case SYS_REG:
			i.operands[idx].strRepr = SystemReg(operand.Reg[0]).String()
			break
		case MULTI_REG:
			// if (get_multireg_operand(
			// 			&instruction->operands[i],
			// 			tmpOperandString,
			// 			sizeof(tmpOperandString)) != DISASM_SUCCESS)
			// {
			// 	return "", fmt.Errorf("failed to disassemble operation")
			// }
			// operand = tmpOperandString;
			break
		case IMPLEMENTATION_SPECIFIC:
			// if (get_implementation_specific(
			// 		&instruction->operands[i],
			// 		tmpOperandString,
			// 		sizeof(tmpOperandString)) != DISASM_SUCCESS)
			// {
			// 	return "", fmt.Errorf("failed to disassemble operation")
			// }
			// operand = tmpOperandString;
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
			// if (get_memory_operand(&instruction->operands[i],
			// 			tmpOperandString,
			// 			sizeof(tmpOperandString)) != DISASM_SUCCESS)
			// 	return "", fmt.Errorf("failed to disassemble operation")
			// operand = tmpOperandString;
			// break
		case CONDITION:
			// if (snprintf(tmpOperandString, sizeof(tmpOperandString), "%s", get_condition((Condition)instruction->operands[i].reg[0])) < 0)
			// 	return "", fmt.Errorf("failed to disassemble operation")
			// operand = tmpOperandString;
		case NONE:
			break
		}

		if operand.OpClass != NONE {
			if idx == 0 {
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

func (op *InstructionOperand) getShiftedImmediate() error {
	var shiftBuff string
	var immBuff string
	var sign string

	if op == nil {
		return failedToDisassembleOperand
	}

	if op.SignedImm == 1 && int64(op.Immediate) < 0 {
		sign = "-"
	}
	if op.ShiftType != SHIFT_NONE {
		if op.ShiftValueUsed != 0 {
			immBuff = fmt.Sprintf(" #%#x", op.ShiftValue)
		}
		shiftBuff = fmt.Sprintf(", %s%s", op.ShiftType, immBuff)
	}
	if op.OpClass == FIMM32 {
		shiftBuff = fmt.Sprintf("#%f%s", float64(op.Immediate), shiftBuff)
	} else if op.OpClass == IMM32 {
		shiftBuff = fmt.Sprintf("#%s%#x%s", sign, uint32(op.Immediate), shiftBuff)
	}

	if op.SignedImm == 1 && int64(op.Immediate) < 0 {
		op.strRepr = fmt.Sprintf("#%s%#016x%s", sign, int64(op.Immediate), shiftBuff)
	}
	op.strRepr = fmt.Sprintf("#%s%#x%s", sign, op.Immediate, shiftBuff)

	return nil
}
