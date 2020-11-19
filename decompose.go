package arm64

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
)

func reg(arg1, arg2, arg3 int) uint32 {
	return uint32(regMap[arg1][arg2][arg3])
}

func (i *Instruction) deleteOperand(index int) {
	copy(i.operands[index:], i.operands[index+1:])                    // Shift a[i+1:] left one index.
	i.operands[len(i.operands)-1] = InstructionOperand{OpClass: NONE} // Erase last element (write zero value).
	// i.operands = i.operands[:len(i.operands)-1]                       // Truncate slice.
	i.operands[len(i.operands)-1] = InstructionOperand{OpClass: NONE}
}

func (i *Instruction) decompose_add_sub_carry() (*Instruction, error) {

	decode := AddSubWithCarry(i.raw)

	var operation = [2][2]Operation{
		{ARM64_ADC, ARM64_ADCS},
		{ARM64_SBC, ARM64_SBCS},
	}
	i.operation = operation[decode.Op()][decode.S()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[2].OpClass = REG
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))
	if decode.Rn() == 31 {
		if i.operation == ARM64_SBC {
			i.operation = ARM64_NGC
			i.deleteOperand(1)
		} else if i.operation == ARM64_SBCS {
			i.operation = ARM64_NGCS
			i.deleteOperand(1)
		}
	}
	// return decode.Opcode2() != 0
	return i, nil
}

func (i *Instruction) decompose_add_sub_extended_reg() (*Instruction, error) {

	decode := AddSubExtendedReg(i.raw)

	var operation = [2][2]Operation{{ARM64_ADD, ARM64_ADDS}, {ARM64_SUB, ARM64_SUBS}}
	var regBaseMap = [2]Register{REG_W_BASE, REG_X_BASE}
	var regBaseMap2 = [8]Register{
		REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE,
		REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE,
	}
	var decodeOptionMap = [2]uint32{2, 3}
	var shiftMap = [2][8]ShiftType{
		{
			SHIFT_UXTB, SHIFT_UXTH, SHIFT_UXTW, SHIFT_UXTX,
			SHIFT_SXTB, SHIFT_SXTH, SHIFT_SXTW, SHIFT_SXTX,
		}, {
			SHIFT_UXTB, SHIFT_UXTH, SHIFT_UXTW, SHIFT_UXTX,
			SHIFT_SXTB, SHIFT_SXTH, SHIFT_SXTW, SHIFT_SXTX,
		},
	}
	i.operation = operation[decode.Op()][decode.S()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_SP, int(regBaseMap[decode.Sf()]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_SP, int(regBaseMap[decode.Sf()]), int(decode.Rn()))
	i.operands[2].OpClass = REG
	if decode.Sf() == 0 {
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rm()))
	} else {
		i.operands[2].Reg[0] = reg(REGSET_ZR, int(regBaseMap2[decode.Option()]), int(decode.Rm()))
	}
	i.operands[2].ShiftType = shiftMap[decode.Sf()][decode.Option()]
	i.operands[2].ShiftValueUsed = 0
	//SUBS => Rn == 31
	//ADDS => Rn == 31
	//SUB  => Rd|Rn == 31
	//ADD  => Rd|Rn == 31
	if (decode.Option() == decodeOptionMap[decode.Sf()]) && ((decode.S() == 1 && decode.Rn() == 31) || (decode.S() == 0 && (decode.Rd() == 31 || decode.Rn() == 31))) {
		if decode.Imm() != 0 {
			i.operands[2].ShiftType = SHIFT_LSL
			i.operands[2].ShiftValueUsed = 1
			i.operands[2].ShiftValue = decode.Imm()
		} else {
			i.operands[2].ShiftType = SHIFT_NONE
		}
	} else if decode.Imm() != 0 {
		i.operands[2].ShiftValueUsed = 1
		i.operands[2].ShiftValue = decode.Imm()
	}
	//Now handle aliases
	if decode.Rd() == 31 {
		if i.operation == ARM64_ADDS {
			i.operation = ARM64_CMN
			i.deleteOperand(0)
		} else if i.operation == ARM64_SUBS {
			i.operation = ARM64_CMP
			i.deleteOperand(0)
		}
	}
	// 	return decode.opt != 0
	return i, nil
}

func (i *Instruction) decompose_add_sub_imm() (*Instruction, error) {
	/* C4.4.1 - Add/subtract (immediate)
	 *
	 * ADD  <Wd|WSP>, <Wn|WSP>, #<imm>{, <shift>}
	 * ADDS <Wd>,	 <Wn|WSP>, #<imm>{, <shift>}
	 * SUB  <Wd|WSP>, <Wn|WSP>, #<imm>{, <shift>}
	 * SUBS <Wd>,	 <Wn|WSP>, #<imm>{, <shift>}
	 *
	 * ADD  <Xd|SP>, <Xn|SP>, #<imm>{, <shift>}
	 * ADDS <Xd>,	<Xn|SP>, #<imm>{, <shift>}
	 * SUB  <Xd|SP>, <Xn|SP>, #<imm>{, <shift>}
	 * SUBS <Xd>,	<Xn|SP>, #<imm>{, <shift>}
	 */

	type decodeRec struct {
		operation Operation
		regType   uint32
	}
	var operationMap = [2][2]decodeRec{
		{{ARM64_ADD, REGSET_SP}, {ARM64_ADDS, REGSET_ZR}},
		{{ARM64_SUB, REGSET_SP}, {ARM64_SUBS, REGSET_ZR}}}
	var regBaseMap = [2]uint32{REG_W_BASE, REG_X_BASE}

	decode := AddSubImm(i.raw)

	i.operation = operationMap[decode.Op()][decode.S()].operation

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = uint32(regMap[operationMap[decode.Op()][decode.S()].regType][regBaseMap[decode.Sf()]][decode.Rd()])

	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = uint32(regMap[REGSET_SP][regBaseMap[decode.Sf()]][decode.Rn()])

	i.operands[2].OpClass = IMM32
	i.operands[2].Immediate = uint64(decode.Imm())
	if decode.Shift() == 1 {
		i.operands[2].ShiftValue = 12
		i.operands[2].ShiftValueUsed = 1
		i.operands[2].ShiftType = SHIFT_LSL
	} else if decode.Shift() > 1 {
		return nil, failedToDecodeInstruction
	}
	//Check for alias
	if i.operation == ARM64_SUBS && decode.Rd() == 31 {
		i.operation = ARM64_CMP
		i.deleteOperand(0)
	} else if i.operation == ARM64_ADD && i.operands[2].Immediate == 0 && decode.Shift() == 0 && (decode.Rd() == 31 || decode.Rn() == 31) {
		i.operation = ARM64_MOV
		i.operands[2].OpClass = NONE
	} else if i.operation == ARM64_ADDS && decode.Rd() == 31 {
		i.operation = ARM64_CMN
		i.deleteOperand(0)
	}

	return i, nil
}

func (i *Instruction) decompose_add_sub_imm_tags() (*Instruction, error) {
	/*
	 * ADDG <Xd|SP>, <Xn|SP>, #<uimm6>, #<uimm4>
	 * SUBG <Xd|SP>, <Xn|SP>, #<uimm6>, #<uimm4>
	 */
	decode := AddSubImmTags(i.raw)

	// ADDG: 1	0	0	1	0	0	0	1	1	0	uimm6	(0)	(0)	uimm4	Xn	Xd
	// SUBG: 1	1	0	1	0	0	0	1	1	0	uimm6	(0)	(0)	uimm4	Xn	Xd
	if ExtractBits(uint32(i.raw), 30, 1) == 0 {
		i.operation = ARM64_ADDG
	} else {
		i.operation = ARM64_SUBG
	}
	// i.operation = BF_GETI(30,1) ? ARM64_SUBG : ARM64_ADDG

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Xd()))

	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Xn()))

	// offset
	i.operands[2].OpClass = IMM64
	i.operands[2].Immediate = uint64(16 * decode.Uimm6())

	// tag_offset
	i.operands[3].OpClass = IMM32
	i.operands[3].Immediate = uint64(decode.Uimm4())

	return i, nil
}

func (i *Instruction) decompose_add_sub_shifted_reg() (*Instruction, error) {
	/* C4.5.2 Add/subtract (shifted register)
	 *
	 * ADD  <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * ADDS <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * SUB  <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * SUBS <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 *
	 * ADD  <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * ADDS <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * SUB  <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * SUBS <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 *
	 * Alias
	 * ADDS WZR, <Wn>, <Wm> {, <shift> #<amount>} -> CMN <Wn>, <Wm>{, <shift> #<amount>}
	 * ADDS XZR, <Xn>, <Xm> {, <shift> #<amount>} -> CMN <Xn>, <Xm>{, <shift> #<amount>}
	 * SUB  <Wd>, WZR, <Wm> {, <shift> #<amount>} -> NEG <Wd>, <Wm>{, <shift> #<amount>}
	 * SUB  <Xd>, XZR, <Xm> {, <shift> #<amount>} -> NEG <Xd>, <Xm>{, <shift> #<amount>}
	 * SUBS WZR, <Wn>, <Wm> {, <shift> #<amount>} -> CMP <Wn>, <Wm>{, <shift> #<amount>}
	 * SUBS XZR, <Xn>, <Xm> {, <shift> #<amount>} -> CMP <Xn>, <Xm>{, <shift> #<amount>}
	 */

	decode := AddSubShiftedReg(i.raw)

	var operation = [2][2]Operation{
		{ARM64_ADD, ARM64_SUB},
		{ARM64_ADDS, ARM64_SUBS}}
	var shift = [4]ShiftType{SHIFT_LSL, SHIFT_LSR, SHIFT_ASR, SHIFT_NONE}
	i.operation = operation[decode.S()][decode.Op()]

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))

	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))

	i.operands[2].OpClass = REG
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))

	if !(decode.Shift() == 0 && decode.Imm() == 0) {
		i.operands[2].ShiftType = shift[decode.Shift()]
		i.operands[2].ShiftValue = decode.Imm()
		i.operands[2].ShiftValueUsed = 1
	}
	//Handle aliases
	if i.operation == ARM64_ADDS && decode.Rd() == 31 {
		i.operation = ARM64_CMN
		i.deleteOperand(0)
	} else if i.operation == ARM64_SUB && decode.Rn() == 31 {
		i.operation = ARM64_NEG
		i.deleteOperand(1)
	} else if i.operation == ARM64_SUBS {
		if decode.Rd() == 31 {
			i.operation = ARM64_CMP
			i.deleteOperand(0)
		} else if decode.Rn() == 31 {
			i.operation = ARM64_NEGS
			i.deleteOperand(1)
		}
	}

	return i, nil
}

func (i *Instruction) decompose_bitfield() (*Instruction, error) {
	/* C4.4.2 Bitfield
	 *
	 * SBFM <Wd>, <Wn>, #<immr>, #<imms>
	 * BFM  <Wd>, <Wn>, #<immr>, #<imms>
	 * UBFM <Wd>, <Wn>, #<immr>, #<imms>
	 *
	 * SBFM <Xd>, <Xn>, #<immr>, #<imms>
	 * BFM  <Xd>, <Xn>, #<immr>, #<imms>
	 * UBFM <Xd>, <Xn>, #<immr>, #<imms>
	 *
	 * Alias
	 * SBFM <Wd>, <Wn>, #<shift>, #31 -> ASR <Wd>, <Wn>, #<shift>
	 * SBFM <Xd>, <Xn>, #<shift>, #63 -> ASR <Xd>, <Xn>, #<shift>
	 * SBFM <Wd>, <Wn>, #(-<lsb> MOD 32), #(<width>-1) -> SBFIZ <Wd>, <Wn>, #<lsb>, #<width>
	 * SBFM <Xd>, <Xn>, #(-<lsb> MOD 64), #(<width>-1) -> SBFIZ <Xd>, <Xn>, #<lsb>, #<width>
	 * SBFM <Wd>, <Wn>, #0, #7 -> SXTB <Wd>, <Wn>
	 * SBFM <Xd>, <Xn>, #0, #7 -> SXTB <Xd>, <Wn>
	 * SBFM <Wd>, <Wn>, #0, #15 -> SXTH <Wd>, <Wn>
	 * SBFM <Xd>, <Xn>, #0, #15 -> SXTH <Xd>, <Wn>
	 * SBFM <Xd>, <Xn>, #0, #31 -> SXTW <Xd>, <Wn>
	 *
	 * BFM <Wd>, WZR, #(-<lsb> MOD 32), #(<width>-1) -> BFC <Wd>, <Wn>, #<lsb>, #<width>
	 * BFM <Wd>, WZR, #(-<lsb> MOD 64), #(<width>-1) -> BFC <Wd>, <Wn>, #<lsb>, #<width>
	 * BFM <Wd>, <Wn>, #(-<lsb> MOD 32), #(<width>-1) -> BFI <Wd>, <Wn>, #<lsb>, #<width>
	 * BFM <Xd>, <Xn>, #(-<lsb> MOD 64), #(<width>-1) -> BFI <Xd>, <Xn>, #<lsb>, #<width>
	 * BFM <Wd>, <Wn>, #<lsb>, #(<lsb>+<width>-1) -> BFXIL <Wd>, <Wn>, #<lsb>, #<width>
	 * BFM <Xd>, <Xn>, #<lsb>, #(<lsb>+<width>-1) -> BFXIL <Xd>, <Xn>, #<lsb>, #<width>
	 *
	 * UBFM <Wd>, <Wn>, #(-<shift> MOD 32), #(31-<shift>) -> LSL <Wd>, <Wn>, #<shift>
	 * UBFM <Xd>, <Xn>, #(-<shift> MOD 64), #(63-<shift>) -> LSL <Xd>, <Xn>, #<shift>
	 * UBFM <Wd>, <Wn>, #(-<lsb> MOD 32), #(<width>-1) -> UBFIZ <Wd>, <Wn>, #<lsb>, #<width>
	 * UBFM <Xd>, <Xn>, #(-<lsb> MOD 64), #(<width>-1) -> UBFIZ <Xd>, <Xn>, #<lsb>, #<width>
	 *
	 */

	decode := Bitfield(i.raw)
	// fmt.Println(decode)
	var operation = []Operation{ARM64_SBFM, ARM64_BFM, ARM64_UBFM, ARM64_UNDEFINED}
	i.operation = operation[decode.Opc()]

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))

	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))

	i.operands[2].OpClass = IMM32
	i.operands[2].Immediate = uint64(decode.Immr())

	i.operands[3].OpClass = IMM32
	i.operands[3].Immediate = uint64(decode.Imms())

	//Handle aliases
	usebfx := bfxPreferred(decode.Sf(), decode.Opc()>>1, decode.Imms(), decode.Immr())
	if i.operation == ARM64_SBFM {
		if decode.Sf() == decode.N() && decode.Imms() == uint32(dataSize[decode.Sf()]-1) {
			i.operation = ARM64_ASR
			i.operands[3].OpClass = NONE
		} else if decode.Imms() < decode.Immr() {
			i.operation = ARM64_SBFIZ
			i.operands[2].Immediate = uint64(uint32(dataSize[decode.Sf()]) - decode.Immr())
			i.operands[3].Immediate++
			// if Register(i.operands[1].Reg[0]) == REG_WZR || Register(i.operands[1].Reg[0]) == REG_XZR {
			// 	i.operation = ARM64_BFC
			// 	i.operands[1] = i.operands[2]
			// 	i.operands[2] = i.operands[3]
			// 	i.operands[3].OpClass = NONE
			// }
		} else if usebfx > 0 {
			i.operation = ARM64_SBFX
			i.operands[3].Immediate -= i.operands[2].Immediate - 1
		} else if i.operands[2].Immediate == 0 {
			switch decode.Imms() {
			case 7:
				i.operation = ARM64_SXTB
				i.operands[1].OpClass = REG
				i.operands[1].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rn()))
				i.operands[2].OpClass = NONE
				i.operands[3].OpClass = NONE
				break
			case 15:
				i.operation = ARM64_SXTH
				i.operands[1].OpClass = REG
				i.operands[1].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rn()))
				i.operands[2].OpClass = NONE
				i.operands[3].OpClass = NONE
				break
			case 31:
				i.operation = ARM64_SXTW
				i.operands[1].OpClass = REG
				i.operands[1].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rn()))
				i.operands[2].OpClass = NONE
				i.operands[3].OpClass = NONE
				break
			default:
				break
			}
		}
	} else if i.operation == ARM64_BFM && decode.Group1() == 0x26 {
		if decode.Imms() < decode.Immr() {
			if decode.Rn() == 31 {
				i.operation = ARM64_BFC
				i.operands[1] = i.operands[2]
				i.operands[2] = i.operands[3]
				i.operands[3].OpClass = NONE
				i.operands[1].Immediate = uint64(uint32(dataSize[decode.Sf()]) - decode.Immr())
				i.operands[2].Immediate++
				// i.operands[3].OpClass = NONE
			} else {
				i.operation = ARM64_BFI
				i.operands[2].OpClass = IMM32
				i.operands[2].Immediate = uint64(uint32(dataSize[decode.Sf()]) - decode.Immr())
				i.operands[3].OpClass = IMM32
				i.operands[3].Immediate++
			}
		} else {
			i.operation = ARM64_BFXIL
			i.operands[3].OpClass = IMM32
			i.operands[3].Immediate -= i.operands[2].Immediate - 1
		}
	} else if i.operation == ARM64_UBFM {
		if decode.Imms() != uint32(dataSize[decode.Sf()]-1) && decode.Imms()+1 == decode.Immr() {
			i.operation = ARM64_LSL
			i.operands[2].OpClass = IMM32
			i.operands[2].Immediate = uint64(uint32(dataSize[decode.Sf()]) - decode.Immr())
			i.operands[3].OpClass = NONE
		} else if decode.Imms() == uint32(dataSize[decode.Sf()]-1) {
			i.operation = ARM64_LSR
			i.operands[3].OpClass = IMM32
			i.operands[3].OpClass = NONE
		} else if decode.Imms() < decode.Immr() {
			i.operation = ARM64_UBFIZ
			i.operands[2].OpClass = IMM32
			i.operands[2].Immediate = uint64(uint32(dataSize[decode.Sf()]) - decode.Immr())
			i.operands[3].OpClass = IMM32
			i.operands[3].Immediate++
		} else if usebfx > 0 {
			i.operation = ARM64_UBFX
			i.operands[3].OpClass = IMM32
			i.operands[3].Immediate -= i.operands[2].Immediate - 1
		} else if decode.Immr() == 0 {
			if decode.Imms() == 7 {
				i.operation = ARM64_UXTB
				i.operands[2].OpClass = NONE
				i.operands[3].OpClass = NONE
			} else if decode.Imms() == 15 {
				i.operation = ARM64_UXTH
				i.operands[2].OpClass = NONE
				i.operands[3].OpClass = NONE
			}
		}
	}

	return i, nil
}

func (i *Instruction) decompose_compare_branch_imm() (*Instruction, error) {
	/*
	 * C4.2.1 Compare & branch (immediate)
	 *
	 * CBZ <Wt>, <label>
	 * CBZ <Xt>, <label>
	 * CBNZ <Xt>, <label>
	 * CBNZ <Wt>, <label>
	 */
	decode := CompareBranchImm(i.raw)
	// fmt.Println(decode)
	var operation = []Operation{ARM64_CBZ, ARM64_CBNZ}
	i.operation = operation[decode.Op()]

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rt()))

	i.operands[1].OpClass = LABEL
	if decode.Imm() < 0 {
		i.operands[1].Immediate = i.address - uint64(decode.Imm()<<2)
		i.operands[1].SignedImm = 1
	} else {
		i.operands[1].Immediate = i.address + uint64(decode.Imm()<<2)
	}

	return i, nil
}

func (i *Instruction) decompose_conditional_branch() (*Instruction, error) {
	/* C4.2.2 Conditional branch (immediate)
	 *
	 * B.<cond> <label>
	 */
	decode := ConditionalBranchImm(i.raw)
	// fmt.Println(decode)
	var operation = []Operation{
		ARM64_B_EQ, ARM64_B_NE, ARM64_B_HS, ARM64_B_LO,
		ARM64_B_MI, ARM64_B_PL, ARM64_B_VS, ARM64_B_VC,
		ARM64_B_HI, ARM64_B_LS, ARM64_B_GE, ARM64_B_LT,
		ARM64_B_GT, ARM64_B_LE, ARM64_B_AL, ARM64_B_NV}
	i.operation = operation[decode.Cond()]

	i.operands[0].OpClass = LABEL
	if decode.Imm() < 0 {
		i.operands[0].Immediate = i.address - uint64(decode.Imm()<<2)
		i.operands[0].SignedImm = 1
	} else {
		i.operands[0].Immediate = i.address + uint64(decode.Imm()<<2)
	}
	if !(decode.O0() == 0 && decode.O1() == 0) {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_conditional_compare_imm() (*Instruction, error) {
	/* C4.5.4 Conditional compare (immediate)
	 *
	 * CCMN <Wn>, #<imm>, #<nzcv>, <cond>
	 * CCMN <Xn>, #<imm>, #<nzcv>, <cond>
	 * CCMP <Wn>, #<imm>, #<nzcv>, <cond>
	 * CCMP <Xn>, #<imm>, #<nzcv>, <cond>
	 */
	decode := ConditionalCompareImm(i.raw)
	var operation = [2]Operation{ARM64_CCMN, ARM64_CCMP}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[1].OpClass = IMM32
	i.operands[1].Immediate = uint64(decode.Imm())
	i.operands[2].OpClass = IMM32
	i.operands[2].Immediate = uint64(decode.Nzcv())
	i.operands[3].OpClass = CONDITION
	i.operands[3].Reg[0] = decode.Cond()

	if decode.O2() != 0 || decode.O3() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_conditional_compare_reg() (*Instruction, error) {
	/* C4.5.5 Conditional compare (register)
	 *
	 * CCMN <Wn>, <Wm>, #<nzcv>, <cond>
	 * CCMN <Xn>, <Xm>, #<nzcv>, <cond>
	 * CCMP <Wn>, <Wm>, #<nzcv>, <cond>
	 * CCMP <Xn>, <Xm>, #<nzcv>, <cond>
	 */
	decode := ConditionalCompareReg(i.raw)
	var operation = [2]Operation{ARM64_CCMN, ARM64_CCMP}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))
	i.operands[2].OpClass = IMM32
	i.operands[2].Immediate = uint64(decode.Nzcv())
	i.operands[3].OpClass = CONDITION
	i.operands[3].Reg[0] = decode.Cond()

	if decode.O2() != 0 || decode.O3() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_conditional_select() (*Instruction, error) {
	/* C4.5.6 Conditional select
	 *
	 * CSEL  <Wd>, <Wn>, <Wm>, <cond>
	 * CSEL  <Xd>, <Xn>, <Xm>, <cond>
	 * CSINC <Wd>, <Wn>, <Wm>, <cond>
	 * CSINC <Xd>, <Xn>, <Xm>, <cond>
	 * CSINV <Wd>, <Wn>, <Wm>, <cond>
	 * CSINV <Xd>, <Xn>, <Xm>, <cond>
	 * CSNEG <Wd>, <Wn>, <Wm>, <cond>
	 * CSNEG <Xd>, <Xn>, <Xm>, <cond>
	 */

	decode := ConditionalSelect(i.raw)

	var operation = [2][2]Operation{
		{ARM64_CSEL, ARM64_CSINC},
		{ARM64_CSINV, ARM64_CSNEG},
	}
	i.operation = operation[decode.Op()][decode.Op2()&1]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[2].OpClass = REG
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))
	i.operands[3].OpClass = CONDITION
	i.operands[3].Reg[0] = decode.Cond()

	if decode.Rm() != 31 && decode.Cond() < 14 && decode.Rn() != 31 && decode.Rn() == decode.Rm() {
		if i.operation == ARM64_CSINC {
			i.operation = ARM64_CINC
			i.operands[3].Reg[0] = (i.operands[3].Reg[0]) ^ 1
			i.deleteOperand(1)
		} else if i.operation == ARM64_CSINV {
			i.operation = ARM64_CINV
			i.operands[3].Reg[0] = (i.operands[3].Reg[0]) ^ 1
			i.deleteOperand(1)
		}
	}

	if decode.Rm() == 31 && decode.Rn() == 31 && decode.Cond() < 14 {
		if i.operation == ARM64_CSINC {
			i.operation = ARM64_CSET
			i.operands[1].Reg[0] = (decode.Cond()) ^ 1
			i.operands[1].OpClass = CONDITION
			i.operands[2].OpClass = NONE
			i.operands[3].OpClass = NONE
		} else if i.operation == ARM64_CSINV {
			i.operation = ARM64_CSETM
			i.operands[1].Reg[0] = (decode.Cond()) ^ 1
			i.operands[1].OpClass = CONDITION
			i.operands[2].OpClass = NONE
			i.operands[3].OpClass = NONE
		}
	}

	if i.operation == ARM64_CSNEG && decode.Cond() < 14 && decode.Rn() == decode.Rm() {
		i.operation = ARM64_CNEG
		i.operands[3].Reg[0] = (i.operands[3].Reg[0]) ^ 1
		i.deleteOperand(1)
	}

	if decode.S() != 0 || decode.Op2() > 1 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_cryptographic_2_register_sha() (*Instruction, error) {
	/* C4.6.21 Cryptographic two-register SHA
	 *
	 * SHA1H <Sd>, <Sn>
	 * SHA1SU1 <Vd>.4S, <Vn>.4S
	 * SHA256SU0 <Vd>.4S, <Vn>.4S
	 */
	decode := Cryptographic2RegSha(i.raw)
	var operation = [4]Operation{ARM64_SHA1H, ARM64_SHA1SU1, ARM64_SHA256SU0, ARM64_UNDEFINED}
	i.operation = operation[decode.Opcode()&3]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	switch decode.Opcode() {
	case 0:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_S_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_S_BASE, int(decode.Rn()))
		break
	case 1:
		fallthrough
	case 2:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
		i.operands[0].ElementSize = 4
		i.operands[0].DataSize = 4
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		i.operands[1].ElementSize = 4
		i.operands[1].DataSize = 4
		break
	default:
		return nil, failedToDecodeInstruction
	}

	if decode.Size() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_cryptographic_3_register_sha() (*Instruction, error) {
	/* C4.6.20 - Cryptographic three-register SHA
	 *
	 * SHA1C	 <Qd>, <Sn>, <Vm>.4S
	 * SHA1P	 <Qd>, <Sn>, <Vm>.4S
	 * SHA1M	 <Qd>, <Sn>, <Vm>.4S
	 * SHA256H   <Qd>, <Qn>, <Vm>.4S
	 * SHA256H2  <Qd>, <Qn>, <Vm>.4S
	 * SHA1SU0   <Vd>.4S, <Vn>.4S, <Vm>.4S
	 * SHA256SU1 <Vd>.4S, <Vn>.4S, <Vm>.4S
	 */
	decode := Cryptographic3RegSha(i.raw)

	var operation = [8]Operation{
		ARM64_SHA1C,
		ARM64_SHA1P,
		ARM64_SHA1M,
		ARM64_SHA1SU0,
		ARM64_SHA256H,
		ARM64_SHA256H2,
		ARM64_SHA256SU1,
		ARM64_UNDEFINED,
	}

	i.operation = operation[decode.Opcode()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	switch decode.Opcode() {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_Q_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_S_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		i.operands[2].ElementSize = 4
		i.operands[2].DataSize = 4
		break
	case 4:
		fallthrough
	case 5:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_Q_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_Q_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		i.operands[2].ElementSize = 4
		i.operands[2].DataSize = 4
		break
	case 3:
		fallthrough
	case 6:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		i.operands[0].ElementSize = 4
		i.operands[0].DataSize = 4
		i.operands[1].ElementSize = 4
		i.operands[1].DataSize = 4
		i.operands[2].ElementSize = 4
		i.operands[2].DataSize = 4
		break
	default:
		return nil, failedToDecodeInstruction
	}
	if decode.Size() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_cryptographic_aes() (*Instruction, error) {
	/* C4.6.19 Cryptographic AES
	 *
	 * AESE   <Vd>.16B, <Vn>.16B
	 * AESD   <Vd>.16B, <Vn>.16B
	 * AESMC  <Vd>.16B, <Vn>.16B
	 * AESIMC <Vd>.16B, <Vn>.16B
	 */
	decode := CryptographicAes(i.raw)
	var operation = [8]Operation{
		ARM64_UNDEFINED,
		ARM64_UNDEFINED,
		ARM64_UNDEFINED,
		ARM64_UNDEFINED,
		ARM64_AESE,
		ARM64_AESD,
		ARM64_AESMC,
		ARM64_AESIMC,
	}
	i.operation = operation[decode.Opcode()&7]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[0].ElementSize = 1
	i.operands[0].DataSize = 16
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[1].ElementSize = 1
	i.operands[1].DataSize = 16

	if decode.Size() != 0 || decode.Opcode() > 7 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_data_processing_1() (*Instruction, error) {
	/* C4.5.7 Data-processing (1 source)
	 *
	 * RBIT <Wd>, <Wn>
	 * RBIT <Xd>, <Xn>
	 * REV16 <Wd>, <Wn>
	 * REV16 <Xd>, <Xn>
	 * REV <Wd>, <Wn>
	 * REV <Xd>, <Xn>
	 * CLZ <Wd>, <Wn>
	 * CLZ <Xd>, <Xn>
	 * CLS <Wd>, <Wn>
	 * CLS <Xd>, <Xn>
	 */

	decode := DataProcessing1(i.raw)
	pac := PointerAuth(i.raw)

	var operation = [2][8]Operation{
		{ARM64_RBIT, ARM64_REV16, ARM64_REV, ARM64_UNDEFINED, ARM64_CLZ, ARM64_CLS, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_RBIT, ARM64_REV16, ARM64_REV32, ARM64_REV, ARM64_CLZ, ARM64_CLS, ARM64_UNDEFINED, ARM64_UNDEFINED},
	}

	var pacOperation = [2][8]Operation{
		{ARM64_PACIA, ARM64_PACIB, ARM64_PACDA, ARM64_PACDB, ARM64_AUTIA, ARM64_AUTIB, ARM64_AUTDA, ARM64_AUTDB},
		{ARM64_PACIZA, ARM64_PACIZB, ARM64_PACDZA, ARM64_PACDZB, ARM64_AUTIZA, ARM64_AUTIZB, ARM64_AUTDZA, ARM64_AUTDZB},
	}

	switch decode.Opcode2() {
	case 0:
		if decode.Opcode() > 5 {
			return i, nil
		}
		i.operation = operation[decode.Sf()][decode.Opcode()]
		i.operands[0].OpClass = REG
		i.operands[0].Reg[0] = uint32(regMap[REGSET_ZR][regSize[decode.Sf()]][decode.Rd()])
		i.operands[1].OpClass = REG
		i.operands[1].Reg[0] = uint32(regMap[REGSET_ZR][regSize[decode.Sf()]][decode.Rn()])
		if decode.S() != 0 || decode.Opcode2() != 0 || i.operation == ARM64_UNDEFINED {
			return nil, failedToDecodeInstruction
		}
	case 1:
		if (decode.Opcode() > 17) || (decode.Sf() != 1) {
			return i, nil
		}
		switch decode.Opcode() {
		case 16:
			i.operation = ARM64_XPACI
			break
		case 17:
			i.operation = ARM64_XPACD
			break
		default:
			i.operation = pacOperation[pac.Z()][pac.Group1()]
			break
		}
		i.operands[0].OpClass = REG
		i.operands[0].Reg[0] = uint32(regMap[REGSET_ZR][regSize[decode.Sf()]][pac.Rd()])
		if decode.Opcode() < 8 {
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = uint32(regMap[REGSET_SP][regSize[decode.Sf()]][pac.Rn()])
		}
		if (decode.Opcode() >= 8) && (pac.Rn() != 0x1f) {
			return nil, failedToDecodeInstruction // TODO should this be error or success?
		}
	default:
		return i, nil
	}

	return i, nil
}

func (i *Instruction) decompose_data_processing_2() (*Instruction, error) {
	/* C4.5.8 Data-processing (2 source)
	 *
	 * UDIV <Wd>, <Wn>, <Wm>
	 * UDIV <Xd>, <Xn>, <Xm>
	 * SDIV <Wd>, <Wn>, <Wm>
	 * SDIV <Wd>, <Wn>, <Wm>
	 * LSLV <Wd>, <Wn>, <Wm>
	 * LSLV <Xd>, <Xn>, <Xm>
	 * LSRV <Wd>, <Wn>, <Wm>
	 * LSRV <Xd>, <Xn>, <Xm>
	 * ASRV <Wd>, <Wn>, <Wm>
	 * ASRV <Xd>, <Xn>, <Xm>
	 * RORV <Wd>, <Wn>, <Wm>
	 * RORV <Xd>, <Xn>, <Xm>
	 * CRC32B <Wd>, <Wn>, <Wm>
	 * CRC32H <Wd>, <Wn>, <Wm>
	 * CRC32W <Wd>, <Wn>, <Wm>
	 * CRC32X <Wd>, <Wn>, <Xm>
	 * CRC32CB <Wd>, <Wn>, <Wm>
	 * CRC32CH <Wd>, <Wn>, <Wm>
	 * CRC32CW <Wd>, <Wn>, <Wm>
	 * CRC32CX <Wd>, <Wn>, <Xm>
	 * PACGA <Xd>, <Xn>, <Xm|SP>
	 * SUBP <Xd>, <Xn|SP>, <Xm|SP>
	 * SUBPS <Xd>, <Xn|SP>, <Xm|SP>
	 */
	var operation = [2][32]Operation{
		{
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UDIV, ARM64_SDIV,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_LSL, ARM64_LSR, ARM64_ASR, ARM64_ROR,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_CRC32B, ARM64_CRC32H, ARM64_CRC32W, ARM64_UNDEFINED,
			ARM64_CRC32CB, ARM64_CRC32CH, ARM64_CRC32CW, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		}, {
			ARM64_SUBP, ARM64_UNDEFINED, ARM64_UDIV, ARM64_SDIV,
			ARM64_IRG, ARM64_GMI, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_LSL, ARM64_LSR, ARM64_ASR, ARM64_ROR,
			ARM64_PACGA, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_CRC32X,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_CRC32CX,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		},
	}

	decode := DataProcessing2(i.raw)
	// fmt.Println(decode)
	if decode.Opcode() > 31 {
		return nil, failedToDisassembleOperation
	}
	if decode.S() > 0 {
		if decode.Opcode() != 0 || decode.Sf() != 1 {
			return nil, failedToDisassembleOperation
		}
		i.operation = ARM64_SUBPS
	} else {
		i.operation = operation[decode.Sf()][decode.Opcode()]
	}

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[2].OpClass = REG
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))

	switch i.operation {
	case ARM64_CRC32X:
		fallthrough
	case ARM64_CRC32CX:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rn()))
		break
	case ARM64_IRG:
		i.operands[0].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
		if decode.Rm() == 0x1f {
			i.operands[2].OpClass = NONE
		}
		fallthrough
		/* intended fall-through */
	case ARM64_SUBP:
		fallthrough
	case ARM64_SUBPS:
		i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rm()))
		break
	case ARM64_GMI:
		i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	default:
		break
	}

	// aliases
	if i.operation == ARM64_SUBPS && decode.S() == 1 && decode.Rd() == 31 {
		i.operation = ARM64_CMPP
		i.operands[0] = i.operands[1]
		i.operands[1] = i.operands[2]
		i.operands[2].OpClass = NONE
	}

	return i, nil
}

func (i *Instruction) decompose_data_processing_3() (*Instruction, error) {

	decode := DataProcessing3(i.raw)

	var operation = [8][2]Operation{
		{ARM64_MADD, ARM64_MSUB},
		{ARM64_SMADDL, ARM64_SMSUBL},
		{ARM64_SMULH, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UMADDL, ARM64_UMSUBL},
		{ARM64_UMULH, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED},
	}

	if decode.Op31() != 0 && decode.Sf() == 0 {
		return nil, failedToDisassembleOperation
	}

	i.operation = operation[decode.Op31()][decode.O0()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	if decode.Op31() == 1 || decode.Op31() == 5 {
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rm()))
	} else {
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))
	}
	i.operands[3].OpClass = REG
	i.operands[3].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Ra()))
	if decode.Ra() == 31 {
		hasAlias := 1
		switch i.operation {
		case ARM64_MADD:
			i.operation = ARM64_MUL
			break
		case ARM64_MSUB:
			i.operation = ARM64_MNEG
			break
		case ARM64_SMADDL:
			i.operation = ARM64_SMULL
			break
		case ARM64_SMSUBL:
			i.operation = ARM64_SMNEGL
			break
		case ARM64_UMADDL:
			i.operation = ARM64_UMULL
			break
		case ARM64_UMSUBL:
			i.operation = ARM64_UMNEGL
			break
		case ARM64_UMULH:
			fallthrough
		case ARM64_SMULH:
			/*Just so we delete the extra operand*/
			break
		default:
			hasAlias = 0
		}
		if hasAlias == 1 {
			i.operands[3].OpClass = NONE
			i.operands[3].Reg[0] = uint32(REG_NONE)
		}
	}

	if i.operation == ARM64_UNDEFINED || decode.Op54() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_exception_generation() (*Instruction, error) {
	/* C4.2.3  - Exception generation
	 *
	 * SVC #<imm>
	 * HVC #<imm>
	 * SMC #<imm>
	 * BRK #<imm>
	 * HLT #<imm>
	 * DCPS1 {#<imm>}
	 * DCPS2 {#<imm>}
	 * DCPS3 {#<imm>}
	 */
	decode := ExceptionGeneration(i.raw)
	var operation = [8][4]Operation{
		{ARM64_UNDEFINED, ARM64_SVC, ARM64_HVC, ARM64_SMC},
		{ARM64_BRK, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_HLT, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_DCPS1, ARM64_DCPS2, ARM64_DCPS3},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
	}
	i.operation = operation[decode.Opc()][decode.Ll()]
	i.operands[0].OpClass = IMM32
	i.operands[0].Immediate = uint64(decode.Imm())
	if decode.Opc() == 5 && decode.Imm() == 0 {
		i.operands[0].OpClass = NONE
	}

	if i.operation == ARM64_UNDEFINED || decode.Op2() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_extract() (*Instruction, error) {
	/* C4.4.3 Extract
	 *
	 * EXTR <Wd>, <Wn>, <Wm>, #<lsb>
	 * EXTR <Xd>, <Xn>, <Xm>, #<lsb>
	 * ROR <Wd>, <Ws>, #<shift> -> EXTR <Wd>, <Ws>, <Ws>, #<shift>
	 * ROR <Xd>, <Xs>, #<shift> -> EXTR <Xd>, <Xs>, <Xs>, #<shift>
	 */

	// EXTRACT decode = *(EXTRACT*)&instructionValue
	decode := Extract(i.raw)
	i.operation = ARM64_EXTR
	if decode.Sf() != decode.N() {
		return i, nil
	}
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[2].OpClass = REG
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))
	i.operands[3].OpClass = IMM32
	i.operands[3].Immediate = uint64(decode.Imms())
	if decode.Rn() == decode.Rm() {
		i.operation = ARM64_ROR
		i.deleteOperand(2)
	}

	if decode.Sf() == 0 && decode.Imms() > 32 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_fixed_floating_conversion() (*Instruction, error) {
	/* C4.6.29 Conversion between floating-point and fixed-point
	 *
	 * SCVTF  <Sd>, <Wn>, #<fbits>
	 * SCVTF  <Dd>, <Wn>, #<fbits>
	 * SCVTF  <Sd>, <Xn>, #<fbits>
	 * SCVTF  <Dd>, <Xn>, #<fbits>
	 * UCVTF  <Sd>, <Wn>, #<fbits>
	 * UCVTF  <Dd>, <Wn>, #<fbits>
	 * UCVTF  <Sd>, <Xn>, #<fbits>
	 * UCVTF  <Dd>, <Xn>, #<fbits>
	 *
	 * FCVTZS <Wd>, <Sn>, #<fbits>
	 * FCVTZS <Xd>, <Sn>, #<fbits>
	 * FCVTZS <Wd>, <Dn>, #<fbits>
	 * FCVTZS <Xd>, <Dn>, #<fbits>
	 * FCVTZU <Wd>, <Sn>, #<fbits>
	 * FCVTZU <Xd>, <Sn>, #<fbits>
	 * FCVTZU <Wd>, <Dn>, #<fbits>
	 * FCVTZU <Xd>, <Dn>, #<fbits>
	 */

	decode := FloatingFixedConversion(i.raw)
	// fmt.Println(decode)
	var operation = [4]Operation{ARM64_FCVTZS, ARM64_FCVTZU, ARM64_SCVTF, ARM64_UCVTF}
	var sdReg = [4]uint32{REG_S_BASE, REG_D_BASE, 0, REG_H_BASE}
	i.operation = operation[decode.Opcode()&3]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = IMM32
	i.operands[2].Immediate = uint64(64 - decode.Scale())

	switch i.operation {
	case ARM64_SCVTF:
		fallthrough
	case ARM64_UCVTF:
		{
			var sdReg = [2]uint32{REG_S_BASE, REG_D_BASE}
			var regSize = [2]uint32{REG_W_BASE, REG_X_BASE}
			i.operands[0].Reg[0] = reg(REGSET_ZR, int(sdReg[decode.Type()&1]), int(decode.Rd()))
			i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
		}
		break
	case ARM64_FCVTZU:
		fallthrough
	case ARM64_FCVTZS:
		if decode.Opcode() <= 1 {
			i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
			i.operands[1].Reg[0] = reg(REGSET_ZR, int(sdReg[decode.Type()]), int(decode.Rn()))
		}
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(sdReg[decode.Type()]), int(decode.Rn()))
		break
	}

	if (decode.Sf() == 0 && (decode.Scale()>>5) == 0) || decode.Type() > 3 || decode.Opcode() > 3 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_compare() (*Instruction, error) {
	/* C4.6.22 Floating-point compare
	 *
	 * FCMP  <Sn>, <Sm>
	 * FCMP  <Dn>, <Dm>
	 * FCMPE <Sn>, <Sm>
	 * FCMPE <Dn>, <Dm>
	 * FCMP  <Sn>, #0.0
	 * FCMP  <Dn>, #0.0
	 * FCMPE <Sn>, #0.0
	 * FCMPE <Dn>, #0.0
	 */
	var operation = [2]Operation{ARM64_FCMP, ARM64_FCMPE}
	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}

	decode := FloatingCompare(i.raw)

	i.operation = operation[(decode.Opcode2()>>4)&1]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rn()))
	if ((decode.Opcode2() >> 3) & 1) == 1 {
		//zero variant
		i.operands[1].OpClass = FIMM32
		i.operands[1].Immediate = uint64(uint32(float64(0.0)))
	} else {
		i.operands[1].OpClass = REG
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rm()))
	}

	if decode.M() != 0 || decode.S() != 0 || decode.Op() != 0 || decode.Type() > 1 || decode.Opcode2()&^uint32(0x18) != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_conditional_compare() (*Instruction, error) {
	/* C4.6.23 Floating-point conditional compare
	 *
	 * FCCMP  <Sn>, <Sm>, #<nzcv>, <cond>
	 * FCCMP  <Dn>, <Dm>, #<nzcv>, <cond>
	 * FCCMPE <Sn>, <Sm>, #<nzcv>, <cond>
	 * FCCMPE <Dn>, <Dm>, #<nzcv>, <cond>
	 */
	decode := FloatingConditionalCompare(i.raw)
	var operation = [2]Operation{ARM64_FCCMP, ARM64_FCCMPE}
	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rn()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rm()))
	i.operands[2].OpClass = IMM32
	i.operands[2].Immediate = uint64(decode.Nzvb())

	i.operands[3].OpClass = CONDITION
	i.operands[3].Reg[0] = decode.Cond()

	if decode.S() != 0 || decode.M() != 0 || decode.Type() > 1 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

// TODO finish
func (i *Instruction) decompose_floating_complex_multiply_accumulate() (*Instruction, error) {
	/* C7.2.62 Floating-point Complex Multiply Accumulate
	 *
	 * FCMLA <Vd>.<T>, <Vn>.<T>, <Vm>.<Ts>[<index>], #<rotate>
	 * FCMLA <Vd>.<T>, <Vn>.<T>, <Vm>.<Ts>[<index>], #<rotate>
	 * FCMLA <Vd>.<T>, <Vn>.<T>, <Vm>.<T>, #<rotate>
	 * FCADD <Vd>.<T>, <Vn>.<T>, <Vm>.<T>, #<rotate>
	 */

	decode := FloatingComplexMultiplyAccumulate(i.raw)
	fmt.Println(decode)

	i.operation = ARM64_FCMLA
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	var rots = [4]uint32{0, 90, 180, 270}
	i.operands[2].ShiftValueUsed = 1
	i.operands[2].ShiftValue = rots[decode.Rot()]
	esize1 := uint32(1 << decode.Size())
	var dsizeMap = [2]uint32{64, 128}
	dsize1 := dsizeMap[decode.Q()] / (8 * esize1)
	var esize2 uint32
	var dsize2 uint32
	switch decode.Size() {
	case 0:
		esize2 = 2
		dsize2 = 8
		break
	case 1:
		esize2 = 4
		dsize2 = 4
		break
	case 2:
		esize2 = 8
		dsize2 = 2
		break
	case 3:
		esize2 = 16
		dsize2 = 1
		break
	}
	var elementMap = [16][3]uint32{
		{0, 1, 1},
		{0, 0, 1},
		{0, 1, 1},
		{0, 0, 1},
		{1, 0, 0},
		{0, 1, 1},
		{1, 0, 0},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
	}
	for idx := 0; idx < 3; idx++ {
		if elementMap[decode.Size()][idx] == 0 {
			i.operands[idx].ElementSize = esize2
			i.operands[idx].DataSize = dsize2
		} else {
			i.operands[idx].ElementSize = esize1
			i.operands[idx].DataSize = dsize1
		}
	}

	return i, nil
}

func (i *Instruction) decompose_floating_cselect() (*Instruction, error) {
	/* C4.6.24 Floating-point conditional select
	 *
	 * FCSEL <Sd>, <Sn>, <Sm>, <cond>
	 * FCSEL <Dd>, <Dn>, <Dm>, <cond>
	 */

	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}
	decode := FloatingConditionalSelect(i.raw)
	i.operation = ARM64_FCSEL
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rn()))
	i.operands[2].OpClass = REG
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rm()))
	i.operands[3].OpClass = CONDITION
	i.operands[3].Reg[0] = decode.Cond()

	if decode.M() != 0 || decode.S() != 0 || decode.Type() > 1 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_data_processing1() (*Instruction, error) {
	/* C4.6.25 Floating-point data-processing (1 source)
	 *
	 * FMOV   <Sd>, <Sn>
	 * FABS   <Sd>, <Sn>
	 * FNEG   <Sd>, <Sn>
	 * FSQRT  <Sd>, <Sn>
	 * FCVT   <Sd>, <Hn>
	 * FCVT   <Hd>, <Sn>
	 * FCVT   <Hd>, <Dn>
	 * FCVT   <Sd>, <Dn>
	 * FRINTN <Sd>, <Sn>
	 * FRINTP <Sd>, <Sn>
	 * FRINTM <Sd>, <Sn>
	 * FRINTZ <Sd>, <Sn>
	 * FRINTA <Sd>, <Sn>
	 * FRINTX <Sd>, <Sn>
	 * FRINTI <Sd>, <Sn>
	 * FMOV   <Dd>, <Dn>
	 * FABS   <Dd>, <Dn>
	 * FNEG   <Dd>, <Dn>
	 * FSQRT  <Dd>, <Dn>
	 * FRINTN <Dd>, <Dn>
	 * FRINTP <Dd>, <Dn>
	 * FRINTM <Dd>, <Dn>
	 * FRINTZ <Dd>, <Dn>
	 * FRINTA <Dd>, <Dn>
	 * FRINTX <Dd>, <Dn>
	 * FRINTI <Dd>, <Dn>
	 * FCVT   <Dd>, <Hn>
	 * FCVT   <Dd>, <Sn>
	 */

	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}
	decode := FloatingDataProcessing1(i.raw)
	var operation = [16]Operation{
		ARM64_FMOV, ARM64_FABS, ARM64_FNEG, ARM64_FSQRT,
		ARM64_FCVT, ARM64_FCVT, ARM64_UNDEFINED, ARM64_FCVT,
		ARM64_FRINTN, ARM64_FRINTP, ARM64_FRINTM, ARM64_FRINTZ,
		ARM64_FRINTA, ARM64_UNDEFINED, ARM64_FRINTX, ARM64_FRINTI,
	}
	i.operation = operation[decode.Opcode()&0xf]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	if (decode.Type() == 3 && (decode.Opcode() == 4 || decode.Opcode() == 5)) || i.operation == ARM64_FCVT {
		var regChoiceCVT = [4]uint32{REG_S_BASE, REG_D_BASE, math.MaxUint32, REG_H_BASE}
		regBase0 := regChoiceCVT[decode.Opcode()&3]
		regBase1 := regChoiceCVT[decode.Type()]
		if regBase0 == math.MaxUint32 || regBase1 == math.MaxUint32 {
			return nil, failedToDecodeInstruction
		}

		i.operation = ARM64_FCVT
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regBase0), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regBase1), int(decode.Rn()))
	} else {
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rn()))
	}

	if decode.M() != 0 || decode.S() != 0 || decode.Opcode() > 15 || i.operation == ARM64_UNDEFINED {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_data_processing2() (*Instruction, error) {
	/* C4.6.26  Floating-point data-processing (2 source)
	 *
	 * FMUL   <Sd>, <Sn>, <Sm>
	 * FDIV   <Sd>, <Sn>, <Sm>
	 * FADD   <Sd>, <Sn>, <Sm>
	 * FSUB   <Sd>, <Sn>, <Sm>
	 * FMAX   <Sd>, <Sn>, <Sm>
	 * FMIN   <Sd>, <Sn>, <Sm>
	 * FMAXNM <Sd>, <Sn>, <Sm>
	 * FMINNM <Sd>, <Sn>, <Sm>
	 * FNMUL  <Sd>, <Sn>, <Sm>
	 *
	 * FMUL   <Dd>, <Dn>, <Dm>
	 * FDIV   <Dd>, <Dn>, <Dm>
	 * FADD   <Dd>, <Dn>, <Dm>
	 * FSUB   <Dd>, <Dn>, <Dm>
	 * FMAX   <Dd>, <Dn>, <Dm>
	 * FMIN   <Dd>, <Dn>, <Dm>
	 * FMAXNM <Dd>, <Dn>, <Dm>
	 * FMINNM <Dd>, <Dn>, <Dm>
	 * FNMUL  <Dd>, <Dn>, <Dm>
	 */
	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}
	var operation = [16]Operation{
		ARM64_FMUL, ARM64_FDIV, ARM64_FADD, ARM64_FSUB,
		ARM64_FMAX, ARM64_FMIN, ARM64_FMAXNM, ARM64_FMINNM,
		ARM64_FNMUL, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
	}
	decode := FloatingDataProcessing2(i.raw)
	i.operation = operation[decode.Opcode()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rm()))

	if decode.M() != 0 || decode.S() != 0 || decode.Type() > 1 || decode.Opcode() > 8 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_data_processing3() (*Instruction, error) {
	/* C4.6.27 Floating-point data-processing (3 source)
	 *
	 * FMADD  <Sd>, <Sn>, <Sm>, <Sa>
	 * FMSUB  <Sd>, <Sn>, <Sm>, <Sa>
	 * FNMADD <Sd>, <Sn>, <Sm>, <Sa>
	 * FNMSUB <Sd>, <Sn>, <Sm>, <Sa>
	 * FMADD  <Dd>, <Dn>, <Dm>, <Da>
	 * FMSUB  <Dd>, <Dn>, <Dm>, <Da>
	 * FNMADD <Dd>, <Dn>, <Dm>, <Da>
	 * FNMSUB <Dd>, <Dn>, <Dm>, <Da>
	 */
	var operation = [2][2]Operation{
		{ARM64_FMADD, ARM64_FMSUB},
		{ARM64_FNMADD, ARM64_FNMSUB},
	}
	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}
	decode := FloatingDataProcessing3(i.raw)
	i.operation = operation[decode.O1()][decode.O0()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[3].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rm()))
	i.operands[3].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Ra()))

	if decode.M() != 0 || decode.S() != 0 || decode.Type() > 1 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func vFPExpandImm(imm8 uint32) uint32 {

	// bits(N) VFPExpandImm(bits(8) imm8)
	// assert N IN {16,32,64};
	// constant integer E = (if N == 16 then 5 elsif N == 32 then 8 else 11); constant integer F = N - E - 1;
	// sign = imm8<7>;
	// exp = NOT(imm8<6>):Replicate(imm8<6>,E-3):imm8<5:4>;
	// frac = imm8<3:0>:Zeros(F-4);
	// result = sign : exp : frac;
	// return result;

	var t ieee754
	var x uint32

	bit6 := (imm8 >> 6) & 1
	bit54 := (imm8 >> 4) & 3

	if bit6 != 0 {
		x = 0x1f
	} else {
		x = 0
	}

	t = t.SetFraction((imm8 & 0xf) << 19)
	t = t.SetExponent((^bit6)<<7 | x<<2 | bit54)
	t = t.SetSign(imm8 >> 7)

	return t.Value()
}

func (i *Instruction) decompose_floating_imm() (*Instruction, error) {
	/* C4.6.28 Floating-point immediate
	 *
	 * FMOV <Sd>, #<imm>
	 * FMOV <Dd>, #<imm>
	 */
	var regChoice = [2]uint32{REG_S_BASE, REG_D_BASE}
	decode := FloatingImm(i.raw)
	i.operation = ARM64_FMOV
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = FIMM32
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.Type()&1]), int(decode.Rd()))
	i.operands[1].Immediate = uint64(vFPExpandImm(decode.Imm8())) // TODO: step this, it's wrong ⚠️

	if decode.Imm5() != 0 || decode.Type() > 1 || decode.M() != 0 || decode.S() != 0 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_integer_conversion() (*Instruction, error) {
	/* C4.6.30 Conversion between floating-point and integer
	 *
	 * FCVTNS <Wd>, <Sn>
	 * FCVTNS <Xd>, <Sn>
	 * FCVTNS <Wd>, <Dn>
	 * FCVTNS <Xd>, <Dn>
	 * FCVTNU <Wd>, <Sn>
	 * FCVTNU <Xd>, <Sn>
	 * FCVTNU <Wd>, <Dn>
	 * FCVTNU <Xd>, <Dn>
	 * FCVTAS <Wd>, <Sn>
	 * FCVTAS <Xd>, <Sn>
	 * FCVTAS <Wd>, <Dn>
	 * FCVTAS <Xd>, <Dn>
	 * FCVTAU <Wd>, <Sn>
	 * FCVTAU <Xd>, <Sn>
	 * FCVTAU <Wd>, <Dn>
	 * FCVTAU <Xd>, <Dn>
	 * FCVTPS <Wd>, <Sn>
	 * FCVTPS <Xd>, <Sn>
	 * FCVTPS <Wd>, <Dn>
	 * FCVTPS <Xd>, <Dn>
	 * FCVTPU <Wd>, <Sn>
	 * FCVTPU <Xd>, <Sn>
	 * FCVTPU <Wd>, <Dn>
	 * FCVTPU <Xd>, <Dn>
	 * FCVTMS <Wd>, <Sn>
	 * FCVTMS <Xd>, <Sn>
	 * FCVTMS <Wd>, <Dn>
	 * FCVTMS <Xd>, <Dn>
	 * FCVTMU <Wd>, <Sn>
	 * FCVTMU <Xd>, <Sn>
	 * FCVTMU <Wd>, <Dn>
	 * FCVTMU <Xd>, <Dn>
	 * FCVTZS <Wd>, <Sn>
	 * FCVTZS <Xd>, <Sn>
	 * FCVTZS <Wd>, <Dn>
	 * FCVTZS <Xd>, <Dn>
	 * FCVTZU <Wd>, <Sn>
	 * FCVTZU <Xd>, <Sn>
	 * FCVTZU <Wd>, <Dn>
	 * FCVTZU <Xd>, <Dn>
	 *
	 * SCVTF  <Sd>, <Wn>
	 * SCVTF  <Dd>, <Wn>
	 * SCVTF  <Sd>, <Xn>
	 * SCVTF  <Dd>, <Xn>
	 * UCVTF  <Sd>, <Wn>
	 * UCVTF  <Dd>, <Wn>
	 * UCVTF  <Sd>, <Xn>
	 * UCVTF  <Dd>, <Xn>
	 *
	 * FMOV   <Sd>, <Wn>
	 * FMOV   <Wd>, <Sn>
	 * FMOV   <Xd>, <Dn>
	 * FMOV   <Dd>, <Xn>
	 * FMOV   <Vd>.D[1], <Xn>
	 * FMOV   <Xd>, <Vn>.D[1]
	 */
	var operation = [2][4][8]Operation{
		{
			{
				ARM64_FCVTNS, ARM64_FCVTNU, ARM64_SCVTF, ARM64_UCVTF,
				ARM64_FCVTAS, ARM64_FCVTAU, ARM64_FMOV, ARM64_FMOV,
			}, {
				ARM64_FCVTPS, ARM64_FCVTPU, ARM64_UNDEFINED, ARM64_UNDEFINED,
				ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			}, {
				ARM64_FCVTMS, ARM64_FCVTMU, ARM64_UNDEFINED, ARM64_UNDEFINED,
				ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			}, {
				ARM64_FCVTZS, ARM64_FCVTZU, ARM64_SCVTF, ARM64_UCVTF,
				ARM64_FCVTAS, ARM64_FCVTAU, ARM64_UNDEFINED, ARM64_UNDEFINED,
			},
		}, {
			{
				ARM64_FCVTNS, ARM64_FCVTNU, ARM64_SCVTF, ARM64_UCVTF,
				ARM64_FCVTAS, ARM64_FCVTAU, ARM64_FMOV, ARM64_FMOV,
			}, {
				ARM64_FCVTPS, ARM64_FCVTPU, ARM64_UNDEFINED, ARM64_UNDEFINED,
				ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			}, {
				ARM64_FCVTMS, ARM64_FCVTMU, ARM64_UNDEFINED, ARM64_UNDEFINED,
				ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			}, {
				ARM64_FCVTZS, ARM64_FCVTZU, ARM64_UNDEFINED, ARM64_UNDEFINED,
				ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			},
		},
	}

	var srcReg = [2]uint32{REG_S_BASE, REG_D_BASE}
	var dstReg = [2]uint32{REG_W_BASE, REG_X_BASE}

	decode := FloatingIntegerConversion(i.raw)
	fmt.Println(decode)
	i.operation = operation[decode.Type()&1][decode.Rmode()][decode.Opcode()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	if decode.Type() == 2 && decode.Sf() == 1 && decode.Rmode() == 1 && decode.Opcode() >= 6 {
		i.operation = ARM64_FMOV
	}

	switch i.operation {
	case ARM64_SCVTF:
		fallthrough
	case ARM64_UCVTF:
		{
			var sdReg = [2]uint32{REG_S_BASE, REG_D_BASE}
			var wxReg = [2]uint32{REG_W_BASE, REG_X_BASE}
			i.operands[0].Reg[0] = reg(REGSET_ZR, int(sdReg[decode.Type()&1]), int(decode.Rd()))
			i.operands[1].Reg[0] = reg(REGSET_ZR, int(wxReg[decode.Sf()]), int(decode.Rn()))
		}
		break
	case ARM64_FMOV:
		if decode.Sf() == 0 {
			var swReg = [2]uint32{REG_W_BASE, REG_S_BASE}
			i.operands[0].Reg[0] = reg(REGSET_ZR, int(swReg[decode.Opcode()&1]), int(decode.Rd()))
			i.operands[1].Reg[0] = reg(REGSET_ZR, int(swReg[(^decode.Opcode()&1)]), int(decode.Rn())) // TODO: is this always correct? replaced !
		} else {
			reg1 := 1 ^ (decode.Opcode() & 1)
			reg2 := decode.Opcode() & 1
			var vxReg = [2]uint32{REG_V_BASE, REG_X_BASE}
			var dxReg = [2]uint32{REG_D_BASE, REG_X_BASE}
			if decode.Rmode() == 1 {
				i.operands[reg1].Index = 1
				i.operands[reg1].ElementSize = 8
				i.operands[0].Reg[0] = reg(REGSET_ZR, int(vxReg[reg1]), int(decode.Rd()))
				i.operands[1].Reg[0] = reg(REGSET_ZR, int(vxReg[reg2]), int(decode.Rn()))
				i.operands[reg1].Scale = (0x80000000 | 1)
			} else {
				i.operands[0].Reg[0] = reg(REGSET_ZR, int(dxReg[reg1]), int(decode.Rd()))
				i.operands[1].Reg[0] = reg(REGSET_ZR, int(dxReg[reg2]), int(decode.Rn()))
			}
		}
		break
	default:
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(dstReg[decode.Sf()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(srcReg[decode.Type()&1]), int(decode.Rn()))
		break
	}

	if decode.S() != 0 || i.operation == ARM64_UNDEFINED {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_floating_javascript_conversion() (*Instruction, error) {
	/* C7.2.99 Floating-point Javascript Convert to Signed fixed-point
	 * FJCVTZS <Wd>, <Dn>
	 */

	decode := FloatingIntegerConversion(i.raw)

	i.operation = ARM64_FJCVTZS
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = uint32(regMap[REGSET_ZR][REG_W_BASE][decode.Rd()])
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = uint32(regMap[REGSET_ZR][REG_D_BASE][decode.Rn()])
	return i, nil
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
			{ARM64_LDR, REG_W_BASE, 0},
			{ARM64_LDR, REG_X_BASE, 1},
			{ARM64_LDRSW, REG_X_BASE, 1},
			{ARM64_PRFM, REG_W_BASE, 0},
		}, {
			{ARM64_LDR, REG_S_BASE, 0},
			{ARM64_LDR, REG_D_BASE, 0},
			{ARM64_LDR, REG_Q_BASE, 0},
			{ARM64_UNDEFINED, 0, 0},
		},
	}

	op := operand[decode.V()][decode.Opc()]
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
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_mem_tags() (*Instruction, error) {

	/*
	 * STG <Xt|SP>, [<Xn|SP>], #<simm> // post-index
	 * STG <Xt|SP>, [<Xn|SP>, #<simm>]! // pre-index
	 * STG <Xt|SP>, [<Xn|SP>{, #<simm>}] // signed offset
	 *
	 * STZGM <Xt>, [<Xn|SP>]
	 *
	 * LDG <Xt>, [<Xn|SP>{, #<simm>}]
	 *
	 * STZG <Xt|SP>, [<Xn|SP>], #<simm>
	 * STZG <Xt|SP>, [<Xn|SP>, #<simm>]!
	 * STZG <Xt|SP>, [<Xn|SP>{, #<simm>}]
	 *
	 * ST2G <Xt|SP>, [<Xn|SP>], #<simm>
	 * ST2G <Xt|SP>, [<Xn|SP>, #<simm>]!
	 * ST2G <Xt|SP>, [<Xn|SP>{, #<simm>}]
	 *
	 * STGM <Xt>, [<Xn|SP>]
	 *
	 * STZ2G <Xt|SP>, [<Xn|SP>], #<simm>
	 * STZ2G <Xt|SP>, [<Xn|SP>, #<simm>]!
	 * STZ2G <Xt|SP>, [<Xn|SP>{, #<simm>}]
	 *
	 * LDGM <Xt>, [<Xn|SP>]
	 */

	decode := LdstTags(i.raw)

	var operation = [4][2][4]Operation{
		{{ARM64_STZGM, ARM64_STG, ARM64_STG, ARM64_STG},
			{ARM64_UNDEFINED, ARM64_STG, ARM64_STG, ARM64_STG},
		},
		{{ARM64_LDG, ARM64_STZG, ARM64_STZG, ARM64_STZG},
			{ARM64_LDG, ARM64_STZG, ARM64_STZG, ARM64_STZG},
		},
		{{ARM64_STGM, ARM64_ST2G, ARM64_ST2G, ARM64_ST2G},
			{ARM64_UNDEFINED, ARM64_ST2G, ARM64_ST2G, ARM64_ST2G},
		},
		{{ARM64_LDGM, ARM64_STZ2G, ARM64_STZ2G, ARM64_STZ2G},
			{ARM64_UNDEFINED, ARM64_STZ2G, ARM64_STZ2G, ARM64_STZ2G},
		},
	}
	if decode.Imm9() == 0 {
		i.operation = operation[decode.Opc()][0][decode.Op2()]
	} else {
		i.operation = operation[decode.Opc()][1][decode.Op2()]
	}
	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	i.operands[0].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))

	switch (decode.Opc() << 2) | decode.Op2() {
	case 0b0001:
		fallthrough
	case 0b0101:
		fallthrough
	case 0b1001:
		fallthrough
	case 0b1101:
		fallthrough
	case 0b0010:
		fallthrough
	case 0b0110:
		fallthrough
	case 0b1010:
		fallthrough
	case 0b1110:
		fallthrough
	case 0b0011:
		fallthrough
	case 0b0111:
		fallthrough
	case 0b1011:
		fallthrough
	case 0b1111:
		fallthrough
	case 0b0100:
		if i.operation == ARM64_LDG {
			i.operands[0].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
		} else {
			i.operands[0].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rt()))
		}

		i.operands[1].SignedImm = 1
		i.operands[1].Immediate = uint64(decode.Imm9() << 4)
		if decode.Imm9()&0x100 > 0 {
			i.operands[1].Immediate |= 0xFFFFFFFFFFFFF000
		}
		break
	default:
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
		break
	}

	switch (decode.Opc() << 2) | decode.Op2() {
	/* post-index, like: MNEMONIC <Xt|SP> [<Xn|SP>], #<simm> */
	case 0b0001:
		fallthrough
	case 0b0101:
		fallthrough
	case 0b1001:
		fallthrough
	case 0b1101:
		i.operands[1].OpClass = MEM_POST_IDX
		break
	/* signed-offset, like: MNEMONIC <Xt|SP>, [<Xn|SP>{, #<simm>}] */
	case 0b0010:
		fallthrough
	case 0b0110:
		fallthrough
	case 0b1010:
		fallthrough
	case 0b1110:
		fallthrough
	case 0b0100:
		i.operands[1].OpClass = MEM_OFFSET
		break
	/* pre-index, like: MNEMONIC <Xt|SP>, [<Xn|SP>{, #<simm>}]! */
	case 0b0011:
		fallthrough
	case 0b0111:
		fallthrough
	case 0b1011:
		fallthrough
	case 0b1111:
		i.operands[1].OpClass = MEM_PRE_IDX
		break
	/* MNEMONIC <Xt>, [<Xn|SP>] */
	case 0b0000:
		fallthrough
	case 0b1000:
		fallthrough
	case 0b1100:
		i.operands[1].OpClass = MEM_REG
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_unscaled() (*Instruction, error) {

	/*
	 * STLUR <Wt>, [<Xn|SP>{, #<simm>}]
	 * STLURB <Wt>, [<Xn|SP>{, #<simm>}]
	 * STLURH <Wt>, [<Xn|SP>{, #<simm>}]
	 *
	 * LDAPUR <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDAPURB <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDAPURH <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDAPURSB <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDAPURSH <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDAPURSW <Xt>, [<Xn|SP>{, #<simm>}]
	 */

	var operation = [4][4]Operation{
		{ARM64_STLURB, ARM64_STLURH, ARM64_STLUR, ARM64_STLUR},
		{ARM64_LDAPURB, ARM64_LDAPURH, ARM64_LDAPUR, ARM64_LDAPUR},
		{ARM64_LDAPURSB, ARM64_LDAPURSH, ARM64_LDAPURSW, ARM64_UNDEFINED},
		{ARM64_LDAPURSB, ARM64_LDAPURSH, ARM64_UNDEFINED, ARM64_UNDEFINED},
	}
	var regBase = [4][4]uint32{
		{REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE},
		{REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE},
		{REG_X_BASE, REG_X_BASE, REG_X_BASE, REG_W_BASE},
		{REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE},
	}

	decode := LdstTags(i.raw)
	// fmt.Println(decode)
	i.operation = operation[decode.Opc()][decode.Size()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regBase[decode.Opc()][decode.Size()]), int(decode.Rt()))
	i.operands[1].OpClass = MEM_OFFSET
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Immediate = uint64(decode.Imm9())
	if decode.Imm9() < 0 {
		i.operands[1].SignedImm = 1
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_exclusive() (*Instruction, error) {
	/* C4.3.6 Load/store exclusive
	 *
	 * STXRB  <Ws>, <Wt>, [<Xn|SP>{,#0}]
	 * STLXRB <Ws>, <Wt>, [<Xn|SP>{,#0}]
	 * STLLRB <Wt>, [<Xn|SP>{,#0}]
	 * STLRB  <Wt>, [<Xn|SP>{,#0}]
	 *
	 * STXRH  <Ws>, <Wt>, [<Xn|SP>{,#0}]
	 * STLXRH <Ws>, <Wt>, [<Xn|SP>{,#0}]
	 * STLLRH <Wt>, [<Xn|SP>{,#0}]
	 * STLRH  <Wt>, [<Xn|SP>{,#0}]
	 *
	 * STXR  <Ws>, <Wt>, [<Xn|SP>{,#0}]
	 * STLXR <Ws>, <Wt>, [<Xn|SP>{,#0}]
	 * STXP  <Ws>, <Wt1>, <Wt2>, [<Xn|SP>{,#0}]
	 * STLXP <Ws>, <Wt1>, <Wt2>, [<Xn|SP>{,#0}]
	 * STLR  <Wt>, [<Xn|SP>{,#0}]
	 *
	 * STXR  <Ws>, <Xt>, [<Xn|SP>{,#0}]
	 * STLXR <Ws>, <Xt>, [<Xn|SP>{,#0}]
	 * STXP  <Ws>, <Xt1>, <Xt2>, [<Xn|SP>{,#0}]
	 * STLXP <Ws>, <Xt1>, <Xt2>, [<Xn|SP>{,#0}]
	 * STLR  <Xt>, [<Xn|SP>{,#0}]
	 */

	var operation = [4][2][8]Operation{
		{
			{
				ARM64_STXRB, ARM64_STLXRB, ARM64_CASP, ARM64_CASPL,
				ARM64_STLLRB, ARM64_STLRB, ARM64_CASB, ARM64_CASLB,
			}, {
				ARM64_LDXRB, ARM64_LDAXRB, ARM64_CASPA, ARM64_CASPAL,
				ARM64_LDLARB, ARM64_LDARB, ARM64_CASAB, ARM64_CASALB,
			},
		}, {
			{
				ARM64_STXRH, ARM64_STLXRH, ARM64_CASP, ARM64_CASPL,
				ARM64_STLLRH, ARM64_STLRH, ARM64_CASH, ARM64_CASLH,
			}, {
				ARM64_LDXRH, ARM64_LDAXRH, ARM64_CASPA, ARM64_CASPAL,
				ARM64_LDLARH, ARM64_LDARH, ARM64_CASAH, ARM64_CASALH,
			},
		}, {
			{
				ARM64_STXR, ARM64_STLXR, ARM64_STXP, ARM64_STLXP,
				ARM64_STLLR, ARM64_STLR, ARM64_CAS, ARM64_CASL,
			}, {
				ARM64_LDXR, ARM64_LDAXR, ARM64_LDXP, ARM64_LDAXP,
				ARM64_LDLAR, ARM64_LDAR, ARM64_CASA, ARM64_CASAL,
			},
		}, {
			{
				ARM64_STXR, ARM64_STLXR, ARM64_STXP, ARM64_STLXP,
				ARM64_STLLR, ARM64_STLR, ARM64_CAS, ARM64_CASL,
			}, {
				ARM64_LDXR, ARM64_LDAXR, ARM64_LDXP, ARM64_LDAXP,
				ARM64_LDLAR, ARM64_LDAR, ARM64_CASA, ARM64_CASAL,
			},
		},
	}
	var regBase = []uint32{REG_W_BASE, REG_X_BASE}

	decode := LdstExclusive(i.raw)

	opcode := decode.O2()<<2 | decode.O1()<<1 | decode.O0()
	i.operation = operation[decode.Size()][decode.L()][opcode]

	idx := 0
	if i.operation >= ARM64_CASB && i.operation <= ARM64_CASL {
		var baseIdx int
		if decode.Size() == 3 || (decode.Size() == 1 && (opcode == 2 || opcode == 3)) {
			baseIdx = 1
		}
		i.operands[idx].OpClass = REG
		i.operands[idx].Reg[0] = reg(REGSET_ZR, int(regBase[baseIdx]), int(decode.Rs()))
		idx++
		if opcode == 2 || opcode == 3 {
			i.operands[idx].OpClass = REG
			i.operands[idx].Reg[0] = reg(REGSET_ZR, int(regBase[baseIdx]), int(decode.Rs()+1)%32)
			idx++
		}
		i.operands[idx].OpClass = REG
		i.operands[idx].Reg[0] = reg(REGSET_ZR, int(regBase[baseIdx]), int(decode.Rt()))
		idx++
		if opcode == 2 || opcode == 3 {
			i.operands[idx].OpClass = REG
			i.operands[idx].Reg[0] = reg(REGSET_ZR, int(regBase[baseIdx]), int(decode.Rt()+1)%32)
			idx++
		}
		i.operands[idx].OpClass = MEM_REG
		i.operands[idx].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	} else if decode.Size() < 2 {
		i.operands[idx].OpClass = REG
		i.operands[idx].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rs()))
		if opcode != 5 && decode.L() == 0 {
			idx++
		}
		i.operands[idx].OpClass = REG
		i.operands[idx].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rt()))
		idx++
		i.operands[idx].OpClass = MEM_REG
		i.operands[idx].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	} else {
		var decodeSizeIs3 int
		if decode.Size() == 3 {
			decodeSizeIs3 = 1
		}
		i.operands[idx].OpClass = REG
		i.operands[idx].Reg[0] = reg(REGSET_ZR, REG_W_BASE, int(decode.Rs()))
		if opcode != 5 && decode.L() == 0 {
			idx++
		}
		i.operands[idx].OpClass = REG
		i.operands[idx].Reg[0] = reg(REGSET_ZR, int(regBase[decodeSizeIs3]), int(decode.Rt()))
		idx++
		if opcode == 2 || opcode == 3 {
			i.operands[idx].OpClass = REG
			i.operands[idx].Reg[0] = reg(REGSET_ZR, int(regBase[decodeSizeIs3]), int(decode.Rt2()))
			idx++
		}
		i.operands[idx].OpClass = MEM_REG
		i.operands[idx].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_imm_post_idx() (*Instruction, error) {
	/* C4.3.8 Load/store register (immediate post-indexed)
	 *
	 * LDRB/STRB <Wt>, [<Xn|SP>], #<simm>	  //PI
	 * LDRSB	 <Xt>, [<Xn|SP>], #<simm>	 //64bit
	 * LDRSB	 <Wt>, [<Xn|SP>], #<simm>	 //32bit
	 * LDR/STR   <Bt>, [<Xn|SP>], #<simm>   //8bit
	 * LDR/STR   <Ht>, [<Xn|SP>], #<simm>   //16bit
	 * LDR/STR   <St>, [<Xn|SP>], #<simm>   //32bit
	 * LDR/STR   <Dt>, [<Xn|SP>], #<simm>   //64bit
	 * LDR/STR   <Qt>, [<Xn|SP>], #<simm>   //128bit
	 * LDRH/STRH <Wt>, [<Xn|SP>], #<simm>  //pi
	 * LDRSH	 <Wt>, [<Xn|SP>], #<simm>	 //32bit
	 * LDRSH	 <Xt>, [<Xn|SP>], #<simm>	 //64bit
	 * LDRSW	 <Xt>, [<Xn|SP>], #<simm>	 //pi
	 */

	decode := LdstRegPairPostIdx(i.raw)

	type opreg struct {
		operation    Operation
		registerBase uint32
	}
	var operation = [4][2][4]opreg{
		{
			{{ARM64_STRB, REG_W_BASE}, {ARM64_LDRB, REG_W_BASE}, {ARM64_LDRSB, REG_X_BASE}, {ARM64_LDRSB, REG_W_BASE}},
			{{ARM64_STR, REG_B_BASE}, {ARM64_LDR, REG_B_BASE}, {ARM64_STR, REG_Q_BASE}, {ARM64_LDR, REG_Q_BASE}},
		}, {
			{{ARM64_STRH, REG_W_BASE}, {ARM64_LDRH, REG_W_BASE}, {ARM64_LDRSH, REG_X_BASE}, {ARM64_LDRSH, REG_W_BASE}},
			{{ARM64_STR, REG_H_BASE}, {ARM64_LDR, REG_H_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		}, {
			{{ARM64_STR, REG_W_BASE}, {ARM64_LDR, REG_W_BASE}, {ARM64_LDRSW, REG_X_BASE}, {ARM64_UNDEFINED, 0}},
			{{ARM64_STR, REG_S_BASE}, {ARM64_LDR, REG_S_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		}, {
			{{ARM64_STR, REG_X_BASE}, {ARM64_LDR, REG_X_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
			{{ARM64_STR, REG_D_BASE}, {ARM64_LDR, REG_D_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		},
	}
	op := operation[decode.Size()][decode.V()][decode.Opc()]
	i.operation = op.operation
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(op.registerBase), int(decode.Rt()))

	i.operands[1].OpClass = MEM_POST_IDX
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Immediate = uint64(decode.Imm())
	if decode.Imm() < 0 {
		i.operands[1].SignedImm = 1
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_reg_imm_pre_idx() (*Instruction, error) {

	decode := LdstRegImmPreIdx(i.raw)

	type opreg struct {
		operation    Operation
		registerBase uint32
	}

	/* C4.3.9 Load/store register (immediate pre-indexed)
	 *
	 * LDRB/STRB <Wt>, [<Xn|SP>, #<simm>]!
	 * LDRSB	 <Wt>, [<Xn|SP>, #<simm>]!	 //32bit
	 * LDRSB	 <Xt>, [<Xn|SP>, #<simm>]!	 //64bit
	 * LDR/STR   <Bt>, [<Xn|SP>, #<simm>]!
	 * LDR/STR   <Ht>, [<Xn|SP>, #<simm>]!
	 * LDR/STR   <St>, [<Xn|SP>, #<simm>]!
	 * LDR/STR   <Dt>, [<Xn|SP>, #<simm>]!
	 * LDR/STR   <Qt>, [<Xn|SP>, #<simm>]!
	 * LDRH/STRH <Wt>, [<Xn|SP>, #<simm>]!
	 * LDRSH	 <Wt>, [<Xn|SP>, #<simm>]!		   /32bit
	 * LDRSH	 <Xt>, [<Xn|SP>, #<simm>]!		   //64bit
	 * LDRSW	 <Xt>, [<Xn|SP>, #<simm>]!
	 */
	var operation = [4][2][4]opreg{
		{
			{{ARM64_STRB, REG_W_BASE}, {ARM64_LDRB, REG_W_BASE}, {ARM64_LDRSB, REG_X_BASE}, {ARM64_LDRSB, REG_W_BASE}},
			{{ARM64_STR, REG_B_BASE}, {ARM64_LDR, REG_B_BASE}, {ARM64_STR, REG_Q_BASE}, {ARM64_LDR, REG_Q_BASE}},
		}, {
			{{ARM64_STRH, REG_W_BASE}, {ARM64_LDRH, REG_W_BASE}, {ARM64_LDRSH, REG_X_BASE}, {ARM64_LDRSH, REG_W_BASE}},
			{{ARM64_STR, REG_H_BASE}, {ARM64_LDR, REG_H_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		}, {
			{{ARM64_STR, REG_W_BASE}, {ARM64_LDR, REG_W_BASE}, {ARM64_LDRSW, REG_X_BASE}, {ARM64_UNDEFINED, 0}},
			{{ARM64_STR, REG_S_BASE}, {ARM64_LDR, REG_S_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		}, {
			{{ARM64_STR, REG_X_BASE}, {ARM64_LDR, REG_X_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
			{{ARM64_STR, REG_D_BASE}, {ARM64_LDR, REG_D_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		},
	}
	op := operation[decode.Size()][decode.V()][decode.Opc()]
	i.operation = op.operation

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(op.registerBase), int(decode.Rt()))

	i.operands[1].OpClass = MEM_PRE_IDX
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Immediate = uint64(decode.Imm())
	if decode.Imm() < 0 {
		i.operands[1].SignedImm = 1
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_atomic_memory_ops() (*Instruction, error) {

	var operation = [8][4][4]Operation{
		{
			{ARM64_LDADDB, ARM64_LDADDLB, ARM64_LDADDAB, ARM64_LDADDALB},
			{ARM64_LDADDH, ARM64_LDADDLH, ARM64_LDADDAH, ARM64_LDADDALH},
			{ARM64_LDADD, ARM64_LDADDL, ARM64_LDADDA, ARM64_LDADDAL},
			{ARM64_LDADD, ARM64_LDADDL, ARM64_LDADDA, ARM64_LDADDAL},
		}, {
			{ARM64_LDCLRB, ARM64_LDCLRLB, ARM64_LDCLRAB, ARM64_LDCLRALB},
			{ARM64_LDCLRH, ARM64_LDCLRLH, ARM64_LDCLRAH, ARM64_LDCLRALH},
			{ARM64_LDCLR, ARM64_LDCLRL, ARM64_LDCLRA, ARM64_LDCLRAL},
			{ARM64_LDCLR, ARM64_LDCLRL, ARM64_LDCLRA, ARM64_LDCLRAL},
		}, {
			{ARM64_LDEORB, ARM64_LDEORLB, ARM64_LDEORAB, ARM64_LDEORALB},
			{ARM64_LDEORH, ARM64_LDEORLH, ARM64_LDEORAH, ARM64_LDEORALH},
			{ARM64_LDEOR, ARM64_LDEORL, ARM64_LDEORA, ARM64_LDEORAL},
			{ARM64_LDEOR, ARM64_LDEORL, ARM64_LDEORA, ARM64_LDEORAL},
		}, {
			{ARM64_LDSETB, ARM64_LDSETLB, ARM64_LDSETAB, ARM64_LDSETALB},
			{ARM64_LDSETH, ARM64_LDSETLH, ARM64_LDSETAH, ARM64_LDSETALH},
			{ARM64_LDSET, ARM64_LDSETL, ARM64_LDSETA, ARM64_LDSETAL},
			{ARM64_LDSET, ARM64_LDSETL, ARM64_LDSETA, ARM64_LDSETAL},
		}, {
			{ARM64_LDSMAXB, ARM64_LDSMAXLB, ARM64_LDAPRB, ARM64_LDSMAXALB},
			{ARM64_LDSMAXH, ARM64_LDSMAXLH, ARM64_LDAPRH, ARM64_LDSMAXALH},
			{ARM64_LDSMAX, ARM64_LDSMAXL, ARM64_LDAPR, ARM64_LDSMAXAL},
			{ARM64_LDSMAX, ARM64_LDSMAXL, ARM64_LDAPR, ARM64_LDSMAXAL},
		}, {
			{ARM64_LDSMINB, ARM64_LDSMINLB, ARM64_LDSMINAB, ARM64_LDSMINALB},
			{ARM64_LDSMINH, ARM64_LDSMINLH, ARM64_LDSMINAH, ARM64_LDSMINALH},
			{ARM64_LDSMIN, ARM64_LDSMINL, ARM64_LDSMINA, ARM64_LDSMINAL},
			{ARM64_LDSMIN, ARM64_LDSMINL, ARM64_LDSMINA, ARM64_LDSMINAL},
		}, {
			{ARM64_LDUMAXB, ARM64_LDUMAXLB, ARM64_LDUMAXAB, ARM64_LDUMAXALB},
			{ARM64_LDUMAXH, ARM64_LDUMAXLH, ARM64_LDUMAXAH, ARM64_LDUMAXALH},
			{ARM64_LDUMAX, ARM64_LDUMAXL, ARM64_LDUMAXA, ARM64_LDUMAXAL},
			{ARM64_LDUMAX, ARM64_LDUMAXL, ARM64_LDUMAXA, ARM64_LDUMAXAL},
		}, {
			{ARM64_LDUMINB, ARM64_LDUMINLB, ARM64_LDUMINAB, ARM64_LDUMINALB},
			{ARM64_LDUMINH, ARM64_LDUMINLH, ARM64_LDUMINAH, ARM64_LDUMINALH},
			{ARM64_LDUMIN, ARM64_LDUMINL, ARM64_LDUMINA, ARM64_LDUMINAL},
			{ARM64_LDUMIN, ARM64_LDUMINL, ARM64_LDUMINA, ARM64_LDUMINAL},
		},
	}
	var regBase = [4]uint32{REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE}

	decode := LdstAtomic(i.raw)
	fmt.Println(decode)
	i.operation = operation[decode.Opc()][decode.Size()][decode.A()<<1|decode.R()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regBase[decode.Size()]), int(decode.Rs()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regBase[decode.Size()]), int(decode.Rt()))
	i.operands[2].OpClass = MEM_OFFSET
	i.operands[2].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))

	/* C6.2.10x
	* LDAPR <Wt>, [<Xn|SP> {,#0}]
	* LDAPR <Xt>, [<Xn|SP> {,#0}]
	* LDAPRB <Wt>, [<Xn|SP> {,#0}]
	* LDAPRH <Wt>, [<Xn|SP> {,#0}]
	 */
	if (i.raw & 0x38BFC000) == 0x38BFC000 {
		i.operands[0] = i.operands[1]
		i.operands[1] = i.operands[2]
		i.operands[2].OpClass = NONE
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}
	return i, nil
}

func (i *Instruction) decompose_load_store_pac() (*Instruction, error) {

	decode := LdstRegImmPac(i.raw)
	// fmt.Println(decode)
	var operation = []Operation{ARM64_LDRAA, ARM64_LDRAB}
	i.operation = operation[decode.M()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
	if decode.W() > 0 {
		i.operands[1].OpClass = MEM_PRE_IDX
	} else {
		i.operands[1].OpClass = MEM_OFFSET
	}
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	if decode.S() > 0 {
		i.operands[1].Immediate = ^uint64(0xfff)
		i.operands[1].SignedImm = 1
	} else {
		i.operands[1].Immediate = 0
	}
	i.operands[1].Immediate |= uint64(decode.Imm() << decode.Size())

	return i, nil
}

func (i *Instruction) decompose_load_store_no_allocate_pair_offset() (*Instruction, error) {
	/* C4.3.7  Load/store no-allocate pair (offset)
	 *
	 * STNP <Wt1>, <Wt2>, [<Xn|SP>{, #<imm>}]
	 * STNP <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * LDNP <Wt1>, <Wt2>, [<Xn|SP>{, #<imm>}]
	 * LDNP <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * STNP <St1>, <St2>, [<Xn|SP>{, #<imm>}]
	 * STNP <Dt1>, <Dt2>, [<Xn|SP>{, #<imm>}]
	 * STNP <Qt1>, <Qt2>, [<Xn|SP>{, #<imm>}]
	 * LDNP <St1>, <St2>, [<Xn|SP>{, #<imm>}]
	 * LDNP <Dt1>, <Dt2>, [<Xn|SP>{, #<imm>}]
	 * LDNP <Qt1>, <Qt2>, [<Xn|SP>{, #<imm>}]
	 */

	decode := LdstNoAllocPair(i.raw)
	// fmt.Println(decode)
	var operation = [2]Operation{ARM64_STNP, ARM64_LDNP}
	var regChoice = [2][4]uint32{
		{REG_W_BASE, REG_X_BASE, REG_X_BASE, REG_X_BASE},
		{REG_S_BASE, REG_D_BASE, REG_Q_BASE, REG_Q_BASE},
	}
	var immShiftBase uint32
	if decode.V() != 0 {
		immShiftBase = decode.Opc() + 2
	} else {
		immShiftBase = (decode.Opc() >> 1) + 2
	}
	i.operation = operation[decode.L()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = MEM_OFFSET
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.V()][decode.Opc()]), int(decode.Rt()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regChoice[decode.V()][decode.Opc()]), int(decode.Rt2()))
	i.operands[2].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[2].Immediate = uint64(int64(decode.Imm()) << immShiftBase)
	i.operands[2].SignedImm = 1

	if i.operation == ARM64_UNDEFINED || decode.Opc() > 2 {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_reg_imm_common() (*Instruction, error) {
	/*C4.3.14 Load/store register pair (offset)
	 *
	 * STP   <Wt1>, <Wt2>, [<Xn|SP>{, #<imm>}]
	 * STP   <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Wt1>, <Wt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * STP   <St1>, <St2>, [<Xn|SP>{, #<imm>}]
	 * STP   <Dt1>, <Dt2>, [<Xn|SP>{, #<imm>}]
	 * STP   <Qt1>, <Qt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <St1>, <St2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Dt1>, <Dt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Qt1>, <Qt2>, [<Xn|SP>{, #<imm>}]
	 * LDPSW <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * STGP  <Xt1>, <Xt2>, [<Xn|SP>], #<imm>
	 *
	 */

	decode := LdstRegPairOffset(i.raw)

	var shiftBase = []uint8{2, 3}
	var simdShiftBase = []uint8{2, 3, 4}

	i.operands[2].SignedImm = 1
	if i.operation == ARM64_LDPSW || i.operation == ARM64_STGP { // TODO this line had redundant checks in the C
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt2()))
		i.operands[2].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
		if i.operation == ARM64_STGP {
			i.operands[2].Immediate = uint64(decode.Imm() << 4)
		} else {
			i.operands[2].Immediate = uint64(decode.Imm() << shiftBase[decode.Opc()>>1])
		}
	} else if decode.V() == 0 {
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Opc()>>1]), int(decode.Rt()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Opc()>>1]), int(decode.Rt2()))
		i.operands[2].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
		i.operands[2].Immediate = uint64(decode.Imm() << shiftBase[decode.Opc()>>1])
	} else {
		if decode.Opc() == 3 {
			return nil, failedToDecodeInstruction
		}

		i.operands[0].Reg[0] = reg(REGSET_ZR, int(simdRegSize[decode.Opc()]), int(decode.Rt()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(simdRegSize[decode.Opc()]), int(decode.Rt2()))
		i.operands[2].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
		i.operands[2].Immediate = uint64(decode.Imm() << simdShiftBase[decode.Opc()])
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_reg_pair_pre_idx() (*Instruction, error) {
	/* C4.3.16 Load/store register pair (pre-indexed)
	 *
	 * STP   <Wt1>, <Wt2>, [<Xn|SP>, #<imm>]!
	 * STP   <Xt1>, <Xt2>, [<Xn|SP>, #<imm>]!
	 * LDP   <Wt1>, <Wt2>, [<Xn|SP>, #<imm>]!
	 * LDP   <Xt1>, <Xt2>, [<Xn|SP>, #<imm>]!
	 * STP   <St1>, <St2>, [<Xn|SP>, #<imm>]!
	 * STP   <Dt1>, <Dt2>, [<Xn|SP>, #<imm>]!
	 * STP   <Qt1>, <Qt2>, [<Xn|SP>, #<imm>]!
	 * LDP   <St1>, <St2>, [<Xn|SP>, #<imm>]!
	 * LDP   <Dt1>, <Dt2>, [<Xn|SP>, #<imm>]!
	 * LDP   <Qt1>, <Qt2>, [<Xn|SP>, #<imm>]!
	 * LDPSW <Xt1>, <Xt2>, [<Xn|SP>, #<imm>]!
	 */
	var operation = [4][2][2]Operation{
		{
			{ARM64_STP, ARM64_LDP},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_STGP, ARM64_LDPSW},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_STP, ARM64_LDP},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_UNDEFINED, ARM64_UNDEFINED},
			{ARM64_UNDEFINED, ARM64_UNDEFINED},
		},
	}
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = MEM_PRE_IDX

	decode := LdstRegPairOffset(i.raw)

	i.operation = operation[decode.Opc()][decode.V()][decode.L()]

	return i.decompose_load_store_reg_imm_common()
}

func (i *Instruction) decompose_load_store_reg_pair_offset() (*Instruction, error) {
	/* C4.3.14 Load/store register pair (offset)
	 *
	 * STP   <Wt1>, <Wt2>, [<Xn|SP>{, #<imm>}]
	 * STP   <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Wt1>, <Wt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 * STP   <St1>, <St2>, [<Xn|SP>{, #<imm>}]
	 * STP   <Dt1>, <Dt2>, [<Xn|SP>{, #<imm>}]
	 * STP   <Qt1>, <Qt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <St1>, <St2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Dt1>, <Dt2>, [<Xn|SP>{, #<imm>}]
	 * LDP   <Qt1>, <Qt2>, [<Xn|SP>{, #<imm>}]
	 * LDPSW <Xt1>, <Xt2>, [<Xn|SP>{, #<imm>}]
	 */

	decode := LdstRegPairOffset(i.raw)

	var operation = [4][2][2]Operation{
		{
			{ARM64_STP, ARM64_LDP},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_STGP, ARM64_LDPSW},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_STP, ARM64_LDP},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_UNDEFINED, ARM64_UNDEFINED},
			{ARM64_UNDEFINED, ARM64_UNDEFINED},
		},
	}
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = MEM_OFFSET
	i.operation = operation[decode.Opc()][decode.V()][decode.L()]
	return i.decompose_load_store_reg_imm_common()
}

func (i *Instruction) decompose_load_store_reg_pair_post_idx() (*Instruction, error) {
	/* 4.3.15 Load/store register pair (post-indexed)
	 *
	 * STP   <Wt1>, <Wt2>, [<Xn|SP>], #<imm>
	 * STP   <Xt1>, <Xt2>, [<Xn|SP>], #<imm>
	 * LDP   <Wt1>, <Wt2>, [<Xn|SP>], #<imm>
	 * LDP   <Xt1>, <Xt2>, [<Xn|SP>], #<imm>
	 * STP   <St1>, <St2>, [<Xn|SP>], #<imm>
	 * STP   <Dt1>, <Dt2>, [<Xn|SP>], #<imm>
	 * STP   <Qt1>, <Qt2>, [<Xn|SP>], #<imm>
	 * LDP   <St1>, <St2>, [<Xn|SP>], #<imm>
	 * LDP   <Dt1>, <Dt2>, [<Xn|SP>], #<imm>
	 * LDP   <Qt1>, <Qt2>, [<Xn|SP>], #<imm>
	 * LDPSW <Xt1>, <Xt2>, [<Xn|SP>], #<imm>
	 * STGP  <Xt1>, <Xt2>, [<Xn|SP>], #<imm>
	 */
	decode := LdstRegPairOffset(i.raw)
	var operation = [4][2][2]Operation{
		{
			{ARM64_STP, ARM64_LDP},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_STGP, ARM64_LDPSW},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_STP, ARM64_LDP},
			{ARM64_STP, ARM64_LDP},
		}, {
			{ARM64_UNDEFINED, ARM64_UNDEFINED},
			{ARM64_UNDEFINED, ARM64_UNDEFINED},
		},
	}
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = MEM_POST_IDX
	i.operation = operation[decode.Opc()][decode.V()][decode.L()]
	return i.decompose_load_store_reg_imm_common()
}

func (i *Instruction) decompose_load_store_reg_reg_offset() (*Instruction, error) {
	/* C4.3.10 Load/store register (register offset)
	 *
	 * STRB   <Wt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDRB   <Wt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDRSB  <Wt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDRSB  <Xt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * STR	<Bt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * STR	<Ht>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * STR	<St>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * STR	<Dt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * STR	<Qt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDR	<Bt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDR	<Ht>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDR	<St>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDR	<Dt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDR	<Qt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * STRH   <Wt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDRH   <Wt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * LDRSW  <Xt>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 * PRFM <prfop>, [<Xn|SP>, <R><m>{, <extend> {<amount>}}]
	 */
	decode := LdstRegRegOffset(i.raw)
	type opreg struct {
		operation    Operation
		registerBase uint32
		amount       [2]int32
	}
	var operation = [4][2][4]opreg{
		{
			{
				{ARM64_STRB, REG_W_BASE, [2]int32{-1, 0}},
				{ARM64_LDRB, REG_W_BASE, [2]int32{-1, 0}},
				{ARM64_LDRSB, REG_X_BASE, [2]int32{-1, 0}},
				{ARM64_LDRSB, REG_W_BASE, [2]int32{-1, 0}},
			}, {
				{ARM64_STR, REG_B_BASE, [2]int32{-1, 0}},
				{ARM64_LDR, REG_B_BASE, [2]int32{-1, 0}},
				{ARM64_STR, REG_Q_BASE, [2]int32{0, 4}},
				{ARM64_LDR, REG_Q_BASE, [2]int32{0, 4}},
			},
		}, {
			{
				{ARM64_STRH, REG_W_BASE, [2]int32{0, 1}},
				{ARM64_LDRH, REG_W_BASE, [2]int32{0, 1}},
				{ARM64_LDRSH, REG_X_BASE, [2]int32{0, 1}},
				{ARM64_LDRSH, REG_W_BASE, [2]int32{0, 1}},
			}, {
				{ARM64_STR, REG_H_BASE, [2]int32{0, 1}},
				{ARM64_LDR, REG_H_BASE, [2]int32{0, 1}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
			},
		}, {
			{
				{ARM64_STR, REG_W_BASE, [2]int32{0, 2}},
				{ARM64_LDR, REG_W_BASE, [2]int32{0, 2}},
				{ARM64_LDRSW, REG_X_BASE, [2]int32{0, 2}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
			}, {
				{ARM64_STR, REG_S_BASE, [2]int32{0, 2}},
				{ARM64_LDR, REG_S_BASE, [2]int32{0, 2}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
			},
		}, {
			{
				{ARM64_STR, REG_X_BASE, [2]int32{0, 3}},
				{ARM64_LDR, REG_X_BASE, [2]int32{0, 3}},
				{ARM64_PRFM, REG_PF_BASE, [2]int32{0, 0}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
			}, {
				{ARM64_STR, REG_D_BASE, [2]int32{0, 3}},
				{ARM64_LDR, REG_D_BASE, [2]int32{0, 3}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
				{ARM64_UNDEFINED, 0, [2]int32{0, 0}},
			},
		},
	}

	op := operation[decode.Size()][decode.V()][decode.Opc()]
	var extendRegister = []uint32{0, 0, REG_W_BASE, REG_X_BASE}
	var extendMap = []ShiftType{
		SHIFT_NONE, SHIFT_NONE, SHIFT_UXTW, SHIFT_LSL,
		SHIFT_NONE, SHIFT_NONE, SHIFT_SXTW, SHIFT_SXTX}

	if decode.Option()>>1 == 0 || decode.Option()>>1 == 2 {
		return nil, failedToDecodeInstruction
	}
	i.operation = op.operation

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(op.registerBase), int(decode.Rt()))

	i.operands[1].OpClass = MEM_EXTENDED
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Reg[1] = reg(REGSET_ZR, int(extendRegister[decode.Option()&3]), int(decode.Rm()))

	extend := extendMap[decode.Option()]
	i.operands[1].ShiftType = extend
	i.operands[1].ShiftValueUsed = 1
	i.operands[1].ShiftValue = uint32(op.amount[decode.S()])

	if i.operands[1].ShiftValue == math.MaxUint32 {
		i.operands[1].ShiftValueUsed = 0
		i.operands[1].ShiftValue = 0
		if i.operands[1].ShiftType == SHIFT_LSL {
			i.operands[1].ShiftType = SHIFT_NONE
		}
	} else if i.operation == ARM64_LDRB {
		if i.operands[1].ShiftType == SHIFT_LSL && i.operands[1].ShiftValue == 0 {
			i.operands[1].ShiftValueUsed = 1
		} else if i.operands[1].ShiftType != SHIFT_LSL && i.operands[1].ShiftValue == 0 {
			i.operands[1].ShiftValueUsed = 0
		}
	} else if i.operands[1].ShiftValue == 0 && (i.operation == ARM64_LDRSB || i.operation == ARM64_STRB) {
		i.operands[1].ShiftValueUsed = 1
	} else if i.operands[1].ShiftValue == 0 {
		if i.operands[1].ShiftType == SHIFT_LSL {
			i.operands[1].ShiftType = SHIFT_NONE
		}
		i.operands[1].ShiftValueUsed = 0
	}

	if i.operation == ARM64_UNDEFINED || extend == SHIFT_NONE {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_reg_unpriv() (*Instruction, error) {
	/* C4.3.11  Load/store register (unprivileged)
	 *
	 * STTRB  <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDTRB  <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDTRSB <Wt>, [<Xn|SP>{, #<simm>}]
	 * STTRH  <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDTRH  <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDTRSH <Wt>, [<Xn|SP>{, #<simm>}]
	 * STTR   <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDTR   <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDTRSB <Xt>, [<Xn|SP>{, #<simm>}]
	 * LDTRSH <Xt>, [<Xn|SP>{, #<simm>}]
	 * STTR   <Xt>, [<Xn|SP>{, #<simm>}]
	 * LDTR   <Xt>, [<Xn|SP>{, #<simm>}]
	 * LDTRSW <Xt>, [<Xn|SP>{, #<simm>}]
	 */

	decode := LdstRegisterUnpriv(i.raw)

	var operation = [4][4]Operation{
		{ARM64_STTRB, ARM64_LDTRB, ARM64_LDTRSB, ARM64_LDTRSB},
		{ARM64_STTRH, ARM64_LDTRH, ARM64_LDTRSH, ARM64_LDTRSH},
		{ARM64_STTR, ARM64_LDTR, ARM64_LDTRSW, ARM64_UNDEFINED},
		{ARM64_STTR, ARM64_LDTR, ARM64_UNDEFINED, ARM64_UNDEFINED},
	}
	i.operation = operation[decode.Size()][decode.Opc()]
	var regChoice int
	if decode.Opc() == 2 || decode.Size() == 3 {
		regChoice = 1
	} else {
		regChoice = 0
	}
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[regChoice]), int(decode.Rt()))
	i.operands[1].OpClass = MEM_OFFSET
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Immediate = uint64(decode.Imm())
	i.operands[1].SignedImm = 1

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_reg_unscalled_imm() (*Instruction, error) {
	/* C4.3.12 - Load/store register (unscaled immediate)
	 *
	 * LDURB/STURB <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDURSB	  <Xt>, [<Xn|SP>{, #<simm>}]
	 * LDURSB	  <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDURH/STURH <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDURSH	  <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <Wt>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <Xt>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <Bt>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <Ht>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <St>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <Dt>, [<Xn|SP>{, #<simm>}]
	 * LDUR/STUR   <Qt>, [<Xn|SP>{, #<simm>}]
	 */

	decode := LdstRegUnscaledImm(i.raw)

	type opreg struct {
		operation    Operation
		registerBase uint32
	}
	var operation = [4][2][4]opreg{
		{
			{{ARM64_STURB, REG_W_BASE}, {ARM64_LDURB, REG_W_BASE}, {ARM64_LDURSB, REG_X_BASE}, {ARM64_LDURSB, REG_W_BASE}},
			{{ARM64_STUR, REG_B_BASE}, {ARM64_LDUR, REG_B_BASE}, {ARM64_STUR, REG_Q_BASE}, {ARM64_LDUR, REG_Q_BASE}},
		}, {
			{{ARM64_STURH, REG_W_BASE}, {ARM64_LDURH, REG_W_BASE}, {ARM64_LDURSH, REG_X_BASE}, {ARM64_LDURSH, REG_W_BASE}},
			{{ARM64_STUR, REG_H_BASE}, {ARM64_LDUR, REG_H_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		}, {
			{{ARM64_STUR, REG_W_BASE}, {ARM64_LDUR, REG_W_BASE}, {ARM64_LDURSW, REG_X_BASE}, {ARM64_UNDEFINED, REG_X_BASE}},
			{{ARM64_STUR, REG_S_BASE}, {ARM64_LDUR, REG_S_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		}, {
			{{ARM64_STUR, REG_X_BASE}, {ARM64_LDUR, REG_X_BASE}, {ARM64_PRFUM, REG_PF_BASE}, {ARM64_UNDEFINED, 0}},
			{{ARM64_STUR, REG_D_BASE}, {ARM64_LDUR, REG_D_BASE}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}},
		},
	}
	op := operation[decode.Size()][decode.V()][decode.Opc()]
	i.operation = op.operation
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = MEM_OFFSET

	i.operands[0].Reg[0] = reg(REGSET_ZR, int(op.registerBase), int(decode.Rt()))
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Immediate = uint64(decode.Imm())
	i.operands[1].SignedImm = 1

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_load_store_reg_unsigned_imm() (*Instruction, error) {
	/* C4.3.13 Load/store register (unsigned immediate)
	 *
	 * STRB	<Wt>, [<Xn|SP>{, #<pimm>}]
	 * LDRB	<Wt>, [<Xn|SP>{, #<pimm>}]
	 * LDRSB   <Wt>, [<Xn|SP>{, #<pimm>}]
	 * LDRSB   <Xt>, [<Xn|SP>{, #<pimm>}]
	 * STR	 <Bt>, [<Xn|SP>{, #<pimm>}]
	 * STR	 <Ht>, [<Xn|SP>{, #<pimm>}]
	 * STR	 <St>, [<Xn|SP>{, #<pimm>}]
	 * STR	 <Dt>, [<Xn|SP>{, #<pimm>}]
	 * STR	 <Qt>, [<Xn|SP>{, #<pimm>}]
	 * LDR	 <Bt>, [<Xn|SP>{, #<pimm>}]
	 * LDR	 <Ht>, [<Xn|SP>{, #<pimm>}]
	 * LDR	 <St>, [<Xn|SP>{, #<pimm>}]
	 * LDR	 <Dt>, [<Xn|SP>{, #<pimm>}]
	 * LDR	 <Qt>, [<Xn|SP>{, #<pimm>}]
	 * STRH	<Wt>, [<Xn|SP>{, #<pimm>}]
	 * LDRH	<Wt>, [<Xn|SP>{, #<pimm>}]
	 * LDRSH   <Wt>, [<Xn|SP>{, #<pimm>}]
	 * LDRSH   <Xt>, [<Xn|SP>{, #<pimm>}]
	 * LDRSW   <Xt>, [<Xn|SP>{, #<pimm>}]
	 * PRFM <prfop>, [<Xn|SP>{, #<pimm>}]
	 */

	decode := LdstRegUnsignedImm(i.raw)

	type opreg struct {
		operation    Operation
		registerBase uint32
		amount       uint32
	}
	var operation = [4][2][4]opreg{
		{
			{
				{ARM64_STRB, REG_W_BASE, 0},
				{ARM64_LDRB, REG_W_BASE, 0},
				{ARM64_LDRSB, REG_X_BASE, 0},
				{ARM64_LDRSB, REG_W_BASE, 0},
			}, {
				{ARM64_STR, REG_B_BASE, 0},
				{ARM64_LDR, REG_B_BASE, 0},
				{ARM64_STR, REG_Q_BASE, 4},
				{ARM64_LDR, REG_Q_BASE, 4},
			},
		}, {
			{
				{ARM64_STRH, REG_W_BASE, 1},
				{ARM64_LDRH, REG_W_BASE, 1},
				{ARM64_LDRSH, REG_X_BASE, 1},
				{ARM64_LDRSH, REG_W_BASE, 1},
			}, {
				{ARM64_STR, REG_H_BASE, 1},
				{ARM64_LDR, REG_H_BASE, 1},
				{ARM64_UNDEFINED, 0, 0},
				{ARM64_UNDEFINED, 0, 0},
			},
		}, {
			{
				{ARM64_STR, REG_W_BASE, 2},
				{ARM64_LDR, REG_W_BASE, 2},
				{ARM64_LDRSW, REG_X_BASE, 2},
				{ARM64_UNDEFINED, 0, 2},
			}, {
				{ARM64_STR, REG_S_BASE, 2},
				{ARM64_LDR, REG_S_BASE, 2},
				{ARM64_UNDEFINED, 0, 0},
				{ARM64_UNDEFINED, 0, 0},
			},
		}, {
			{
				{ARM64_STR, REG_X_BASE, 3},
				{ARM64_LDR, REG_X_BASE, 3},
				{ARM64_PRFM, REG_PF_BASE, 3},
				{ARM64_UNDEFINED, 0, 0},
			}, {
				{ARM64_STR, REG_D_BASE, 3},
				{ARM64_LDR, REG_D_BASE, 3},
				{ARM64_UNDEFINED, 0, 0},
				{ARM64_UNDEFINED, 0, 0},
			},
		},
	}

	op := operation[decode.Size()][decode.V()][decode.Opc()]
	i.operation = op.operation
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = MEM_OFFSET

	i.operands[0].Reg[0] = reg(REGSET_ZR, int(op.registerBase), int(decode.Rt()))
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	i.operands[1].Immediate = uint64(decode.Imm()) << op.amount
	i.operands[1].SignedImm = 0

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func moveWidePreferred(sf, immn, imms, immr uint32) bool {
	/*
	 * boolean MoveWidePreferred(bit sf, bit immN, bits(6) imms, bits(6) immr)
	 *  integer S = UInt (imms)
	 *  integer R = UInt (immr)
	 *  integer width = if sf == '1' then 64 else 32
	 *  //element size must equal total immediate size
	 *  if sf == '1' && immN:imms != '1xxxxxx' then
	 *	  return FALSE
	 *  if sf == '0' && immN:imms != '00xxxxx' then
	 *	  return FALSE
	 *
	 *	// for MOVZ must contain no more than 16 ones
	 *	if S < 16 then
	 *	 // ones must not span halfword boundary when rotated
	 *	   return (-R MOD 16) <= (15 - S)
	 *
	 *  // for MOVN must contain no more than 16 zeros
	 *  if S >= width - 15 then
	 *   // zeros must not span halfword boundary when rotated
	 *   return (R MOD 16) <= (S - (width - 15))
	 *	return FALSE
	 */
	S := int32(imms)
	R := int32(immr)
	var width int32
	if sf == 1 {
		width = 64
	} else {
		width = 32
	}

	if sf == 1 && (((immn<<6|imms)>>6)&1) != 1 {
		return false
	}
	if sf == 0 && (((immn<<6|imms)>>5)&3) != 0 {
		return false
	}
	if S < 16 {
		return (-R % 16) <= (15 - S)
	}
	if S >= width-15 {
		return (R % 16) <= (S - (width - 15))
	}
	return false
}

func (i *Instruction) decompose_logical_imm() (*Instruction, error) {
	/* C4.4.4 Logical (immediate)
	 *
	 * AND <Wd|WSP>, <Wn>, #<imm>
	 * AND <Xd|SP>, <Xn>, #<imm>
	 * ORR <Wd|WSP>, <Wn>, #<imm>
	 * ORR <Xd|SP>, <Xn>, #<imm>
	 * EOR <Wd|WSP>, <Wn>, #<imm>
	 * EOR <Xd|SP>, <Xn>, #<imm>
	 * ANDS <Wd>, <Wn>, #<imm>
	 * ANDS <Xd>, <Xn>, #<imm>
	 */

	decode := LogicalImm(i.raw)
	// fmt.Println(decode)
	var operation = [4]Operation{ARM64_AND, ARM64_ORR, ARM64_EOR, ARM64_ANDS}
	i.operation = operation[decode.Opc()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_SP, int(regSize[decode.Sf()]), int(decode.Rd()))

	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))

	i.operands[2].OpClass = IMM64
	outBits := uint32(32)
	if decode.Sf() > 0 {
		outBits = 64
	}
	i.operands[2].Immediate = DecodeBitMasks(decode.N(), decode.Imms(), decode.Immr(), outBits)
	if i.operands[2].Immediate == 0 {
		return nil, failedToDecodeInstruction
	}

	if i.operation == ARM64_ORR && decode.Rn() == 31 && !moveWidePreferred(decode.Sf(), decode.N(), decode.Imms(), decode.Immr()) {
		i.operation = ARM64_MOV
		i.deleteOperand(1)
	}
	if i.operation == ARM64_ANDS && decode.Rd() == 31 {
		i.operation = ARM64_TST
		i.deleteOperand(0)
	}
	if (decode.Sf() == 0) && (decode.N() != 0) {
		return nil, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_logical_shifted_reg() (*Instruction, error) {
	/* C4.5.10 Logical (shifted register)
	 *
	 * AND <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * AND <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * BIC <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * BIC <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * ORR <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * ORR <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * ORN <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * ORN <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * EOR <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * EOR <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * EON <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * EON <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * ANDS <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * ANDS <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 * BICS <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
	 * BICS <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
	 *
	 * Aliases
	 * ORR <Wd>, <Wn>, <Wm>{, <shift> #<amount>} -> MOV <Wd>, <Wm>
	 * ORN <Wd>, WZR, <Wm>{, <shift> #<amount>}  -> MVN <Wd>, <Wm>{, <shift> #<amount>}
	 * ANDS WZR, <Wn>, <Wm>{, <shift> #<amount>} -> TST <Wn>, <Wm>{, <shift> #<amount>}
	 */
	// LOGICAL_SHIFTED_REG decode = *(LOGICAL_SHIFTED_REG*)&instructionValue
	decode := LogicalShiftedReg(i.raw)
	var operation = [2][4]Operation{
		{ARM64_AND, ARM64_ORR, ARM64_EOR, ARM64_ANDS},
		{ARM64_BIC, ARM64_ORN, ARM64_EON, ARM64_BICS}}
	var shiftMap = [4]ShiftType{SHIFT_LSL, SHIFT_LSR, SHIFT_ASR, SHIFT_ROR}
	i.operation = operation[decode.N()][decode.Opc()]
	// i.operands = make([]InstructionOperand, 3)
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG

	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rm()))
	i.operands[2].ShiftType = shiftMap[decode.Shift()]
	i.operands[2].ShiftValue = decode.Imm()
	i.operands[2].ShiftValueUsed = 1

	if i.operands[2].ShiftType == SHIFT_LSL && i.operands[2].ShiftValue == 0 {
		i.operands[2].ShiftType = SHIFT_NONE
	}
	if i.operation == ARM64_ORR && decode.Shift() == 0 && decode.Imm() == 0 && decode.Rn() == 31 {
		i.operation = ARM64_MOV
		i.operands[2].ShiftType = SHIFT_NONE
		i.operands[2].ShiftValue = 0
		i.deleteOperand(1)
	} else if i.operation == ARM64_ORN && decode.Rn() == 31 {
		i.operation = ARM64_MVN
		i.deleteOperand(1)
	} else if i.operation == ARM64_ANDS && decode.Rd() == 31 {
		i.operation = ARM64_TST
		i.deleteOperand(0)
	}

	return i, nil
}

func (i *Instruction) decompose_move_wide_imm() (*Instruction, error) {
	/* C4.4.5 Move wide (immediate)
	 *
	 * MOVN <Wd>, #<imm>{, LSL #<shift>}
	 * MOVN <Xd>, #<imm>{, LSL #<shift>}
	 * MOVZ <Wd>, #<imm>{, LSL #<shift>}
	 * MOVZ <Xd>, #<imm>{, LSL #<shift>}
	 * MOVK <Wd>, #<imm>{, LSL #<shift>}
	 * MOVK <Xd>, #<imm>{, LSL #<shift>}
	 */
	decode := MoveWideImm(i.raw)
	// fmt.Println(decode)
	var operation = [4]Operation{ARM64_MOVN, ARM64_UNDEFINED, ARM64_MOVZ, ARM64_MOVK}
	i.operation = operation[decode.Opc()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Sf()]), int(decode.Rd()))

	i.operands[1].OpClass = IMM32
	i.operands[1].Immediate = uint64(decode.Imm())
	if decode.Imm() < 0 {
		i.operands[1].SignedImm = 1
	}
	if decode.Hw() != 0 {
		i.operands[1].ShiftType = SHIFT_LSL
		i.operands[1].ShiftValue = decode.Hw() << 4
		i.operands[1].ShiftValueUsed = 1
	}
	if (decode.Sf() == 0 && decode.Hw()>>1 == 1) || i.operation == ARM64_UNDEFINED {
		return nil, failedToDecodeInstruction
	}

	if decode.Imm() != 0 || decode.Hw() == 0 {
		if i.operation == ARM64_MOVN && ((decode.Sf() == 0 && decode.Imm() != 0xffff) || decode.Sf() == 1) {
			i.operation = ARM64_MOV
			if decode.Sf() == 1 {
				i.operands[1].OpClass = IMM64
			} else {
				i.operands[1].OpClass = IMM32
			}
			i.operands[1].Immediate = ^(i.operands[1].Immediate << (decode.Hw() << 4))
			i.operands[1].ShiftType = SHIFT_NONE
			i.operands[1].ShiftValue = 0
		} else if i.operation == ARM64_MOVZ {
			i.operation = ARM64_MOV
			i.operands[1].OpClass = IMM64
			i.operands[1].Immediate <<= i.operands[1].ShiftValue
			i.operands[1].ShiftType = SHIFT_NONE
			i.operands[1].ShiftValue = 0
		}
	}
	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_pc_rel_addr() (*Instruction, error) {
	/* C4.4.6 PC-rel. addressing
	 *
	 * ADR <Xd>, <label>
	 * ADRP <Xd>, <label>
	 */
	decode := PcRelAddressing(i.raw)
	var operation = []Operation{ARM64_ADR, ARM64_ADRP}
	var shiftBase = []uint8{0, 12}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rd()))
	i.operands[1].OpClass = LABEL
	if decode.Immhi() < 0 {
		i.operands[1].SignedImm = 1
	}
	x := uint64(decode.Immhi())
	i.operands[1].Immediate = (((x << 2) | uint64(decode.Immlo())) << shiftBase[decode.Op()])
	if decode.Op() == 1 {
		i.operands[1].Immediate += i.address & ^uint64((1<<12)-1)
	} else {
		i.operands[1].Immediate += i.address
	}
	//printf("imm: %lx %lx %lx %lx\n", i.operands[1].Immediate, x, x<<2, ((x<<2)| decode.immlo)<< 12)
	return i, nil
}

func (i *Instruction) decompose_simd_2_reg_misc() (*Instruction, error) {
	/* C4.6.17 Advanced SIMD two-register miscellaneous
	 *
	 * REV64	 <Vd>.<T>, <Vn>.<T>
	 * REV16	 <Vd>.<T>, <Vn>.<T>
	 * SADDLP	<Vd>.<Ta>, <Vn>.<Tb>
	 * SUQADD	<Vd>.<T>, <Vn>.<T>
	 * CLS	   <Vd>.<T>, <Vn>.<T>
	 * CNT	   <Vd>.<T>, <Vn>.<T>
	 * SADALP	<Vd>.<Ta>, <Vn>.<Tb>
	 * SQABS	 <Vd>.<T>, <Vn>.<T>
	 * CMGT	  <V><d>, <V><n>, #0
	 * CMEQ	  <V><d>, <V><n>, #0
	 * CMLT	  <V><d>, <V><n>, #0
	 * ABS	   <Vd>.<T>, <Vn>.<T>
	 * XTN{2}	<Vd>.<Tb>, <Vn>.<Ta>
	 * SQXTN{2}  <Vd>.<Tb>, <Vn>.<Ta>
	 * FCVTN{2}  <Vd>.<Tb>, <Vn>.<Ta>
	 * FCVTL{2}  <Vd>.<Ta>, <Vn>.<Tb>
	 * FRINTN	<Vd>.<T>, <Vn>.<T>
	 * FRINTM	<Vd>.<T>, <Vn>.<T>
	 * FCVTAS	<Vd>.<T>, <Vn>.<T>
	 * SCVTF	 <Vd>.<T>, <Vn>.<T>
	 * FCMGT	 <Vd>.<T>, <Vn>.<T>, #0.0
	 * FCMEQ	 <Vd>.<T>, <Vn>.<T>, #0.0
	 * FCMLT	 <Vd>.<T>, <Vn>.<T>, #0.0
	 * FABS	  <Vd>.<T>, <Vn>.<T>
	 * FRINTP	<Vd>.<T>, <Vn>.<T>
	 * FRINTZ	<Vd>.<T>, <Vn>.<T>
	 * FCVTPS	<Vd>.<T>, <Vn>.<T>
	 * FCVTZS	<Vd>.<T>, <Vn>.<T>
	 * URECPE	<Vd>.<T>, <Vn>.<T>
	 * FRECPE	<Vd>.<T>, <Vn>.<T>
	 * REV32	 <Vd>.<T>, <Vn>.<T>
	 * UADDLP	<Vd>.<Ta>, <Vn>.<Tb>
	 * USQADD	<Vd>.<T>, <Vn>.<T>
	 * CLZ	   <Vd>.<T>, <Vn>.<T>
	 * UADALP	<Vd>.<Ta>, <Vn>.<Tb>
	 * SQNEG	 <Vd>.<T>, <Vn>.<T>
	 * CMGE	  <Vd>.<T>, <Vn>.<T>, #0
	 * CMLE	  <Vd>.<T>, <Vn>.<T>, #0
	 * NEG	   <Vd>.<T>, <Vn>.<T>
	 * SQXTUN{2} <Vd>.<Tb>, <Vn>.<Ta>
	 * SHLL{2}   <Vd>.<Ta>, <Vn>.<Tb>, #<shift>
	 * UQXTN{2}  <Vd>.<Tb>, <Vn>.<Ta>
	 * FCVTXN{2} <Vd>.<Tb>, <Vn>.<Ta>
	 * FRINTA	<Vd>.<T>, <Vn>.<T>
	 * FRINTX	<Vd>.<T>, <Vn>.<T>
	 * FCVTNU	<Vd>.<T>, <Vn>.<T>
	 * FCVTMU	<Vd>.<T>, <Vn>.<T>
	 * FCVTAU	<Vd>.<T>, <Vn>.<T>
	 * UCVTF	 <Vd>.<T>, <Vn>.<T>
	 * NOT	   <Vd>.<T>, <Vn>.<T>
	 * RBIT	  <Vd>.<T>, <Vn>.<T>
	 * FCMGE	 <Vd>.<T>, <Vn>.<T>, #0.0
	 * FCMLE	 <Vd>.<T>, <Vn>.<T>, #0.0
	 * FNEG	  <Vd>.<T>, <Vn>.<T>
	 * FRINTI	<Vd>.<T>, <Vn>.<T>
	 * FCVTPU	<Vd>.<T>, <Vn>.<T>
	 * FCVTZU	<Vd>.<T>, <Vn>.<T>
	 * URSQRTE   <Vd>.<T>, <Vn>.<T>
	 * FRSQRTE   <Vd>.<T>, <Vn>.<T>
	 * FSQRT	 <Vd>.<T>, <Vn>.<T>
	 *
	 * 0 - <Vd>.<T>, <Vn>.<T>
	 * 1 - <Vd>.<Ta>, <Vn>.<Tb>
	 * 2 - <Vd>.<Tb>, <Vn>.<Ta>
	 * 3 - <Vd>.<Ta>, <Vn>.<Tb>, #<shift>
	 * 4 - <Vd>.<T>, <Vn>.<T>, #0
	 * 5 - <Vd>.<T>, <Vn>.<T>, #0.0
	 * 6 - {2} <Vd>.<Tb>, <Vn>.<Ta>
	 * 7 - {2} <Vd>.<Ta>, <Vn>.<Tb>
	 * 8 - <V><d>, <V><n>, #0
	 */
	decode := Simd2RegMisc(i.raw)
	type opInfo struct {
		op      Operation
		otype   uint32
		maxSize uint32
	}
	info := opInfo{}
	if decode.U() == 0 {
		var operation1 = []opInfo{
			{ARM64_REV64, 0, 3},
			{ARM64_REV16, 0, 1},
			{ARM64_SADDLP, 1, 2},
			{ARM64_SUQADD, 0, 0},
			{ARM64_CLS, 0, 3},
			{ARM64_CNT, 0, 1},
			{ARM64_SADALP, 1, 3},
			{ARM64_SQABS, 0, 4},
			{ARM64_CMGT, 4, 4},
			{ARM64_CMEQ, 4, 4},
			{ARM64_CMLT, 4, 4},
			{ARM64_ABS, 0, 4},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_XTN, 6, 3}, //{2}
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_SQXTN, 6, 3}, //{2}
		}

		var operation2 = []opInfo{
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_FCVTN, 10, 3}, //{2}
			{ARM64_FCVTL, 12, 1}, //{2}
			{ARM64_FRINTN, 13, 7},
			{ARM64_FRINTM, 13, 7},
			{ARM64_FCVTNS, 13, 7},
			{ARM64_FCVTMS, 13, 7},
			{ARM64_FCVTAS, 13, 7},
			{ARM64_SCVTF, 13, 7},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
		}

		var operation3 = []opInfo{
			{ARM64_FCMGT, 5, 7},
			{ARM64_FCMEQ, 5, 7},
			{ARM64_FCMLT, 5, 7},
			{ARM64_FABS, 0, 7},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_FRINTP, 0, 7},
			{ARM64_FRINTZ, 0, 7},
			{ARM64_FCVTPS, 0, 7},
			{ARM64_FCVTZS, 0, 7},
			{ARM64_URECPE, 0, 8},
			{ARM64_FRECPE, 0, 7},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
		}
		// TODO: make sure COUNT_OF is the same as len()
		if decode.Opcode() < uint32(len(operation1)) && operation1[decode.Opcode()].op != ARM64_UNDEFINED {
			info = operation1[decode.Opcode()]
		} else if decode.Size() < 2 && decode.Opcode() > uint32(len(operation1)) {
			info = operation2[decode.Opcode()-uint32(len(operation1))]
		} else if decode.Size() > 1 {
			info = operation3[decode.Opcode()-12]
		} else {
			return nil, failedToDecodeInstruction
		}
	} else {
		var operation1 = []opInfo{
			{ARM64_REV32, 0, 3},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UADDLP, 1, 3},
			{ARM64_USQADD, 0, 3},
			{ARM64_CLZ, 0, 3},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UADALP, 1, 3},
			{ARM64_SQNEG, 0, 4},
			{ARM64_CMGE, 4, 4},
			{ARM64_CMLE, 4, 4},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_NEG, 0, 4},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_SQXTUN, 6, 3}, //{2}
			{ARM64_SHLL, 3, 3},   //{2}
			{ARM64_UQXTN, 6, 3},  //{2}
			{ARM64_UNDEFINED, 0, 0},
		}

		var operation2 = []opInfo{
			{ARM64_FCVTXN, 11, 5}, //{2}
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_FRINTA, 13, 0},
			{ARM64_FRINTX, 13, 7},
			{ARM64_FCVTNU, 13, 7},
			{ARM64_FCVTMU, 13, 7},
			{ARM64_FCVTAU, 13, 7},
			{ARM64_UCVTF, 13, 7},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
		}

		var operation3 = []opInfo{
			{ARM64_FCMGE, 5, 5},
			{ARM64_FCMLE, 5, 5},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_FNEG, 0, 5},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_FRINTI, 0, 5},
			{ARM64_FCVTPU, 0, 5},
			{ARM64_FCVTZU, 0, 5},
			{ARM64_URSQRTE, 0, 5},
			{ARM64_FRSQRTE, 0, 5},
			{ARM64_UNDEFINED, 0, 0},
			{ARM64_FSQRT, 0, 7},
		}

		var operation4 = []opInfo{
			{ARM64_MVN, 9, 0},
			{ARM64_RBIT, 9, 0},
		}
		if decode.Opcode() == 5 {
			info = operation4[decode.Size()&1]
		} else if decode.Opcode() < uint32(len(operation1)) && operation1[decode.Opcode()].op != ARM64_UNDEFINED {
			info = operation1[decode.Opcode()]
		} else if decode.Size() < 2 && decode.Opcode() >= 22 {
			info = operation2[decode.Opcode()-22]
		} else if decode.Size() > 1 && decode.Opcode() >= 12 {
			info = operation3[decode.Opcode()-12]
		} else {
			return nil, failedToDecodeInstruction
		}
	}
	i.operation = info.op
	/* 0 - <Vd>.<T>, <Vn>.<T>
	 * 1 - <Vd>.<Ta>, <Vn>.<Tb>
	 * 2 - <Vd>.<Tb>, <Vn>.<Ta>
	 * 3 - {2} <Vd>.<Ta>, <Vn>.<Tb>, #<shift>
	 * 4 - <Vd>.<T>, <Vn>.<T>, #0
	 * 5 - <Vd>.<T>, <Vn>.<T>, #0.0
	 * 6 - {2} <Vd>.<Tb>, <Vn>.<Ta>
	 * 7 - {2} <Vd>.<Ta>, <Vn>.<Tb>
	 * 8 - <V><d>, <V><n>, #0
	 */
	var elemSize1 uint32
	var elemSize2 uint32
	var dataSize1 uint32
	var dataSize2 uint32
	var dsizeMap = [2]uint32{64, 128}
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	switch info.otype {
	case 0:
		elemSize1 = 1 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (8 << decode.Size())
		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())
		break
	case 1:
		elemSize1 = 2 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (16 << decode.Size())
		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())
		break
	case 2:
		elemSize1 = 1 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (8 << decode.Size())
		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())
		break
	case 3:
		i.operation = Operation(uint32(i.operation) + decode.Q()) //the '2' variant is always +1
		elemSize1 = 2 << decode.Size()
		dataSize1 = 16 / (2 << decode.Size())

		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())

		i.operands[2].Immediate = 8 << decode.Size()
		i.operands[2].OpClass = IMM32
		break
	case 4:
		elemSize1 = 1 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (8 << decode.Size())
		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())
		i.operands[2].Immediate = 0
		i.operands[2].OpClass = IMM32
		break
	case 5:
		elemSize1 = 1 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (8 << decode.Size())
		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())
		i.operands[2].Immediate = 0
		i.operands[2].OpClass = IMM32
		break
	case 6:
		//good
		i.operation = Operation(uint32(i.operation) + decode.Q()) //the '2' variant is always +1
		elemSize1 = 1 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (elemSize1 << 3)

		elemSize2 = 2 << decode.Size()
		dataSize2 = 64 / (8 << decode.Size())
		break
	case 7:
		i.operation = Operation(uint32(i.operation) + decode.Q()) //the '2' variant is always +1
		elemSize1 = 2 << decode.Size()
		dataSize1 = 16 / (2 << decode.Size())

		elemSize2 = 1 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (8 << decode.Size())
		break
	case 8:
		break
	case 9:
		dataSize1 = dsizeMap[decode.Q()] / 8
		elemSize1 = 1
		dataSize2 = dsizeMap[decode.Q()] / 8
		elemSize2 = 1
		break
	case 10:
		i.operation = Operation(uint32(i.operation) + decode.Q()) //the '2' variant is always +1
		elemSize1 = 2 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (16 << decode.Size())

		elemSize2 = 4 << decode.Size()
		dataSize2 = 4 >> decode.Size()
		break
	case 11:
		i.operation = Operation(uint32(i.operation) + decode.Q()) //the '2' variant is always +1
		elemSize1 = 4
		dataSize1 = 2 << decode.Q()

		elemSize2 = 8
		dataSize2 = 2
		break
	case 12:
		i.operation = Operation(uint32(i.operation) + decode.Q()) //the '2' variant is always +1
		elemSize1 = 4 << decode.Size()
		dataSize1 = 4 >> decode.Size()

		elemSize2 = 2 << decode.Size()
		dataSize2 = dsizeMap[decode.Q()] / (16 << decode.Size())
		break
	case 13:
		elemSize1 = 4 << decode.Size()
		dataSize1 = dsizeMap[decode.Q()] / (32 << decode.Size())
		elemSize2 = elemSize1
		dataSize2 = dataSize1
		break
	}
	//element = b(1),h(2),s(4),d(8)
	//data 1,2,3,8,16
	i.operands[0].ElementSize = elemSize1
	i.operands[0].DataSize = dataSize1

	i.operands[1].ElementSize = elemSize2
	i.operands[1].DataSize = dataSize2

	switch info.maxSize {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		if decode.Size() > info.maxSize {
			return nil, failedToDecodeInstruction
		}
		break
	case 4:
		if decode.Size() == 3 && decode.Q() == 0 {
			return nil, failedToDecodeInstruction
		}
		break
	case 5:
		if decode.Size() == 0 {
			return nil, failedToDecodeInstruction
		}
		break
	case 6:
		fallthrough
	case 7:
		if decode.Size() == 1 && decode.Q() == 0 {
			return nil, failedToDecodeInstruction
		}
		break
	case 8:
		if decode.Size() == 1 {
			return nil, failedToDecodeInstruction
		}
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_3_different() (*Instruction, error) {
	/* C4.6.15 Advanced SIMD three different
	 *
	 * SADDL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SADDW{2}   <Vd>.<Ta>, <Vn>.<Ta>, <Vm>.<Tb>
	 * SSUBL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SSUBW{2}   <Vd>.<Ta>, <Vn>.<Ta>, <Vm>.<Tb>
	 * ADDHN{2}   <Vd>.<Tb>, <Vn>.<Ta>, <Vm>.<Ta>
	 * SABAL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SUBHN{2}   <Vd>.<Tb>, <Vn>.<Ta>, <Vm>.<Ta>
	 * SABDL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SMLAL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SQDMLAL{2} <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SMLSL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SQDMLSL{2} <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SMULL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * SQDMULL{2} <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * PMULL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * UADDL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * UADDW{2}   <Vd>.<Ta>, <Vn>.<Ta>, <Vm>.<Tb>
	 * USUBL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * USUBW{2}   <Vd>.<Ta>, <Vn>.<Ta>, <Vm>.<Tb>
	 * RADDHN{2}  <Vd>.<Tb>, <Vn>.<Ta>, <Vm>.<Ta>
	 * UABAL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * RSUBHN{2}  <Vd>.<Tb>, <Vn>.<Ta>, <Vm>.<Ta>
	 * UABDL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * UMLAL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * UMLSL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 * UMULL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Tb>
	 */

	var operation = [2][2][16]Operation{
		{
			{
				ARM64_SADDL,
				ARM64_SADDW,
				ARM64_SSUBL,
				ARM64_SSUBW,
				ARM64_ADDHN,
				ARM64_SABAL,
				ARM64_SUBHN,
				ARM64_SABDL,
				ARM64_SMLAL,
				ARM64_SQDMLAL,
				ARM64_SMLSL,
				ARM64_SQDMLSL,
				ARM64_SMULL,
				ARM64_SQDMULL,
				ARM64_PMULL,
				ARM64_UNDEFINED,
			}, {
				ARM64_UADDL,
				ARM64_UADDW,
				ARM64_USUBL,
				ARM64_USUBW,
				ARM64_RADDHN,
				ARM64_UABAL,
				ARM64_RSUBHN,
				ARM64_UABDL,
				ARM64_UMLAL,
				ARM64_UNDEFINED,
				ARM64_UMLSL,
				ARM64_UNDEFINED,
				ARM64_UMULL,
				ARM64_UNDEFINED,
				ARM64_UNDEFINED,
				ARM64_UNDEFINED,
			},
		}, {
			{
				ARM64_SADDL2,
				ARM64_SADDW2,
				ARM64_SSUBL2,
				ARM64_SSUBW2,
				ARM64_ADDHN2,
				ARM64_SABAL2,
				ARM64_SUBHN2,
				ARM64_SABDL2,
				ARM64_SMLAL2,
				ARM64_SQDMLAL2,
				ARM64_SMLSL2,
				ARM64_SQDMLSL2,
				ARM64_SMULL2,
				ARM64_SQDMULL2,
				ARM64_PMULL2,
				ARM64_UNDEFINED,
			}, {
				ARM64_UADDL2,
				ARM64_UADDW2,
				ARM64_USUBL2,
				ARM64_USUBW2,
				ARM64_RADDHN2,
				ARM64_UABAL2,
				ARM64_RSUBHN2,
				ARM64_UABDL2,
				ARM64_UMLAL2,
				ARM64_UNDEFINED,
				ARM64_UMLSL2,
				ARM64_UNDEFINED,
				ARM64_UMULL2,
				ARM64_UNDEFINED,
				ARM64_UNDEFINED,
				ARM64_UNDEFINED,
			},
		},
	}
	decode := Simd3Different(i.raw)
	i.operation = operation[decode.Q()][decode.U()][decode.Opcode()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	esize1 := uint32(1 << decode.Size())
	var dsizeMap = [2]uint32{64, 128}
	dsize1 := dsizeMap[decode.Q()] / (8 * esize1)
	var esize2 uint32
	var dsize2 uint32
	switch decode.Size() {
	case 0:
		esize2 = 2
		dsize2 = 8
		break
	case 1:
		esize2 = 4
		dsize2 = 4
		break
	case 2:
		esize2 = 8
		dsize2 = 2
		break
	case 3:
		esize2 = 16
		dsize2 = 1
		break
	}
	var elementMap = [16][3]uint32{
		{0, 1, 1},
		{0, 0, 1},
		{0, 1, 1},
		{0, 0, 1},
		{1, 0, 0},
		{0, 1, 1},
		{1, 0, 0},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
		{0, 1, 1},
	}
	for idx := 0; idx < 3; idx++ {
		if elementMap[decode.Opcode()][idx] == 0 {
			i.operands[idx].ElementSize = esize2
			i.operands[idx].DataSize = dsize2
		} else {
			i.operands[idx].ElementSize = esize1
			i.operands[idx].DataSize = dsize1
		}
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_3_same() (*Instruction, error) {
	/* C4.6.16 Advanced SIMD three same
	 *
	 * SHADD    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SQADD    <V><d>, <V><n>, <V><m>
	 * SRHADD   <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SHSUB    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SQSUB    <V><d>, <V><n>, <V><m>
	 * CMGT     <V><d>, <V><n>, <V><m>
	 * CMGE     <V><d>, <V><n>, <V><m>
	 * SSHL     <V><d>, <V><n>, <V><m>
	 * SQSHL    <V><d>, <V><n>, <V><m>
	 * SRSHL    <V><d>, <V><n>, <V><m>
	 * SQRSHL   <V><d>, <V><n>, <V><m>
	 * SMAX     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SMIN     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SABD     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SABA     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * ADD      <V><d>, <V><n>, <V><m>
	 * CMTST    <V><d>, <V><n>, <V><m>
	 * MLA      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * MUL      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SMAXP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SMINP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SQDMULH  <V><d>, <V><n>, <V><m>
	 * ADDP     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMAXNM   <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMLA     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FADD     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMULX    <V><d>, <V><n>, <V><m>
	 * FCMEQ    <V><d>, <V><n>, <V><m>
	 * FMAX     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FRECPS   <V><d>, <V><n>, <V><m>
	 * AND      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * BIC      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMINNM   <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMLS     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FSUB     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMIN     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FRSQRTS  <V><d>, <V><n>, <V><m>
	 * ORR      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * ORN      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UHADD    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UQADD    <V><d>, <V><n>, <V><m>
	 * URHADD   <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UHSUB    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UQSUB    <V><d>, <V><n>, <V><m>
	 * CMHI     <V><d>, <V><n>, <V><m>
	 * CMHS     <V><d>, <V><n>, <V><m>
	 * USHL     <V><d>, <V><n>, <V><m>
	 * UQRSHL   <V><d>, <V><n>, <V><m>
	 * UMAX     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UMIN     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UABD     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UABA     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SUB      <V><d>, <V><n>, <V><m>
	 * CMEQ     <V><d>, <V><n>, <V><m>
	 * MLS      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * PMUL     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UMAXP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UMINP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * SQRDMULH <V><d>, <V><n>, <V><m>
	 * FMAXNMP  <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FADDP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMUL     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FCMGE    <V><d>, <V><n>, <V><m>
	 * FMAXP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FDIV     <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * EOR      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * BSL      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FMINNMP  <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * FABD     <V><d>, <V><n>, <V><m>
	 * FCMGT    <V><d>, <V><n>, <V><m>
	 * FACGT    <V><d>, <V><n>, <V><m>
	 * FMINP    <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * BIT      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * BIF      <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 */
	decode := Simd3Same(i.raw)
	type opInfo struct {
		op     Operation
		vector uint32
	}
	var alternateEncode uint32
	if decode.U() == 0 {
		var operation1 = []opInfo{
			{ARM64_SHADD, 0},
			{ARM64_SQADD, 0},
			{ARM64_SRHADD, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SHSUB, 0},
			{ARM64_SQSUB, 0},
			{ARM64_CMGT, 0},
			{ARM64_CMGE, 0},
			{ARM64_SSHL, 0},
			{ARM64_SQSHL, 0},
			{ARM64_SRSHL, 0},
			{ARM64_SQRSHL, 0},
			{ARM64_SMAX, 0},
			{ARM64_SMIN, 0},
			{ARM64_SABD, 0},
			{ARM64_SABA, 0},
			{ARM64_ADD, 0},
			{ARM64_CMTST, 0},
			{ARM64_MLA, 0},
			{ARM64_MUL, 0},
			{ARM64_SMAXP, 0},
			{ARM64_SMINP, 0},
			{ARM64_SQDMULH, 0},
			{ARM64_ADDP, 0},
		}
		var operation2 = []opInfo{
			{ARM64_FMAXNM, 1},
			{ARM64_FMLA, 1},
			{ARM64_FADD, 1},
			{ARM64_FMULX, 1},
			{ARM64_FCMEQ, 1},
			{ARM64_UNDEFINED, 0},
			{ARM64_FMAX, 1},
			{ARM64_FRECPS, 1},
		}
		var operation3 = []opInfo{
			{ARM64_FMINNM, 0},
			{ARM64_FMLS, 0},
			{ARM64_FSUB, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_FMIN, 0},
			{ARM64_FRSQRTS, 0},
		}
		if decode.Opcode() < uint32(len(operation1)) {
			if decode.Opcode() == 3 {
				switch decode.Size() {
				case 0:
					i.operation = ARM64_AND
					break
				case 1:
					i.operation = ARM64_BIC
					break
				case 2:
					i.operation = ARM64_ORR
					break
				case 3:
					i.operation = ARM64_ORN
					break
				}
			} else {
				i.operation = operation1[decode.Opcode()].op
			}
		} else if decode.Size() < 2 {
			i.operation = operation2[decode.Opcode()-uint32(len(operation1))].op
			alternateEncode = 1
		} else {
			i.operation = operation3[decode.Opcode()-uint32(len(operation1))].op
		}
	} else {
		var operation1 = []opInfo{
			{ARM64_UHADD, 0},
			{ARM64_UQADD, 0},
			{ARM64_URHADD, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UHSUB, 0},
			{ARM64_UQSUB, 0},
			{ARM64_CMHI, 0},
			{ARM64_CMHS, 0},
			{ARM64_USHL, 0},
			{ARM64_UQSHL, 0},
			{ARM64_URSHL, 0},
			{ARM64_UQRSHL, 0},
			{ARM64_UMAX, 0},
			{ARM64_UMIN, 0},
			{ARM64_UABD, 0},
			{ARM64_UABA, 0},
			{ARM64_SUB, 0},
			{ARM64_CMEQ, 0},
			{ARM64_MLS, 0},
			{ARM64_PMUL, 0},
			{ARM64_UMAXP, 0},
			{ARM64_UMINP, 0},
			{ARM64_SQRDMULH, 0},
			{ARM64_UNDEFINED, 0},
		}

		var operation2 = []opInfo{
			{ARM64_FMAXNMP, 1},
			{ARM64_UNDEFINED, 0},
			{ARM64_FADDP, 1},
			{ARM64_FMUL, 1},
			{ARM64_FCMGE, 1},
			{ARM64_FACGE, 1},
			{ARM64_FMAXP, 1},
			{ARM64_FDIV, 0},
		}

		var operation3 = []opInfo{
			{ARM64_FMINNMP, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_FABD, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_FCMGT, 0},
			{ARM64_FACGT, 0},
			{ARM64_FMINP, 0},
			{ARM64_UNDEFINED, 0},
		}

		if decode.Opcode() < uint32(len(operation1)) {
			if decode.Opcode() == 3 {
				switch decode.Size() {
				case 0:
					i.operation = ARM64_EOR
					break
				case 1:
					i.operation = ARM64_BSL
					break
				case 2:
					i.operation = ARM64_BIT
					break
				case 3:
					i.operation = ARM64_BIF
					break
				}
			} else {
				i.operation = operation1[decode.Opcode()].op
			}
		} else if decode.Size() < 2 {
			i.operation = operation2[decode.Opcode()-uint32(len(operation1))].op
			alternateEncode = 1
		} else {
			i.operation = operation3[decode.Opcode()-uint32(len(operation1))].op
		}

	}
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	if decode.Opcode() == 3 {
		var dsizeMap = [2]uint32{8, 16}
		i.operands[0].ElementSize = 1
		i.operands[0].DataSize = dsizeMap[decode.Q()]
		i.operands[1].ElementSize = 1
		i.operands[1].DataSize = dsizeMap[decode.Q()]
		i.operands[2].ElementSize = 1
		i.operands[2].DataSize = dsizeMap[decode.Q()]
	} else {
		if alternateEncode == 1 {
			esize := uint32(32 << decode.Size())
			var dsizeMap = [2]uint32{64, 128}
			dsize := dsizeMap[decode.Q()] / (esize)
			i.operands[0].ElementSize = esize / 8
			i.operands[1].ElementSize = esize / 8
			i.operands[2].ElementSize = esize / 8
			i.operands[0].DataSize = dsize
			i.operands[1].DataSize = dsize
			i.operands[2].DataSize = dsize
		} else {
			esize := uint32(1 << decode.Size())
			var dsizeMap = [2]uint32{64, 128}
			dsize := dsizeMap[decode.Q()] / (8 * esize)
			i.operands[0].ElementSize = esize
			i.operands[1].ElementSize = esize
			i.operands[2].ElementSize = esize
			i.operands[0].DataSize = dsize
			i.operands[1].DataSize = dsize
			i.operands[2].DataSize = dsize
		}
	}
	//Aliases
	if i.operation == ARM64_ORR && decode.Rn() == decode.Rm() {
		i.operation = ARM64_MOV
		i.operands[2].OpClass = NONE
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_across_lanes() (*Instruction, error) {
	/* C4.6.1 Advanced SIMD across lanes
	 *
	 * SADDLV  <V><d>, <Vn>.<T>
	 * SMAXV   <V><d>, <Vn>.<T>
	 * SMINV   <V><d>, <Vn>.<T>
	 * ADDV    <V><d>, <Vn>.<T>
	 * UADDLV  <V><d>, <Vn>.<T>
	 * UMAXV   <V><d>, <Vn>.<T>
	 * UMINV   <V><d>, <Vn>.<T>
	 * FMAXNMV <V><d>, <Vn>.<T>
	 * FMAXV   <V><d>, <Vn>.<T>
	 * FMINNMV <V><d>, <Vn>.<T>
	 * FMINV   <V><d>, <Vn>.<T>
	 */
	decode := SimdAcrossLanes(i.raw)
	esize := uint32(1 << decode.Size())
	var dsizeMap = [2]uint32{64, 128}
	dsize := dsizeMap[decode.Q()] / (8 * esize)
	var regBaseMap = [3]uint32{REG_B_BASE, REG_H_BASE, REG_S_BASE}
	var regBaseMap2 = [3]uint32{REG_H_BASE, REG_S_BASE, REG_D_BASE}
	var reg1 uint32
	var reg2 uint32
	if decode.Size() == 3 {
		return nil, failedToDecodeInstruction
	}

	switch decode.Opcode() {
	case 3:
		if decode.U() == 0 {
			i.operation = ARM64_SADDLV
		} else {
			i.operation = ARM64_UADDLV
		}
		reg1 = reg(REGSET_ZR, int(regBaseMap2[decode.Size()]), int(decode.Rd()))
		reg2 = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		break
	case 10:
		if decode.U() == 0 {
			i.operation = ARM64_SMAXV
		} else {
			i.operation = ARM64_UMAXV
		}
		reg1 = reg(REGSET_ZR, int(regBaseMap[decode.Size()]), int(decode.Rd()))
		reg2 = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		break
	case 12:
		if decode.U() == 1 {
			if decode.Size() < 2 {
				i.operation = ARM64_FMAXNMV
			} else {
				i.operation = ARM64_FMINNMV
			}
			if decode.Q() == 0 || decode.Size() == 1 {
				return nil, failedToDecodeInstruction
			}
			reg1 = reg(REGSET_ZR, REG_S_BASE, int(decode.Rd()))
			reg2 = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
			esize = 4
			dsize = 4
		}
		break
	case 15:
		if decode.U() == 1 {
			if decode.Size() < 2 {
				i.operation = ARM64_FMAXV
			} else {
				i.operation = ARM64_FMINV
			}
			reg1 = reg(REGSET_ZR, REG_S_BASE, int(decode.Rd()))
			reg2 = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
			esize = 4
			dsize = 4
		}
		break
	case 26:
		if decode.U() == 0 {
			i.operation = ARM64_SMINV
		} else {
			i.operation = ARM64_UMINV
		}
		reg1 = reg(REGSET_ZR, int(regBaseMap[decode.Size()]), int(decode.Rd()))
		reg2 = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		break
	case 27:
		i.operation = ARM64_ADDV
		reg1 = reg(REGSET_ZR, int(regBaseMap[decode.Size()]), int(decode.Rd()))
		reg2 = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		break
	default:
		return nil, failedToDecodeInstruction
	}
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[0].Reg[0] = reg1
	i.operands[1].Reg[0] = reg2
	i.operands[1].ElementSize = esize
	i.operands[1].DataSize = dsize
	return i, nil
}

func (i *Instruction) decompose_simd_copy() (*Instruction, error) {
	/* C4.6.2 - Advanced SIMD copy
	 *
	 * DUP  <V><d>, <Vn>.<T>[<index>]
	 * DUP  <Vd>.<T>, <R><n>
	 * SMOV <Wd>, <Vn>.<Ts>[<index>]
	 * SMOV <Xd>, <Vn>.<Ts>[<index>]
	 * UMOV <Wd>, <Vn>.<Ts>[<index>]
	 * UMOV <Xd>, <Vn>.<Ts>[<index>]
	 * INS  <Vd>.<Ts>[<index>], <R><n>
	 * INS  <Vd>.<Ts>[<index1>], <Vn>.<Ts>[<index2>]
	 *
	 * Aliases:
	 * INS  <Vd>.<Ts>[<index1>], <Vn>.<Ts>[<index2>] -> MOV <Vd>.<Ts>[<index1>], <Vn>.<Ts>[<index2>]
	 * DUP  <V><d>, <Vn>.<T>[<index>] -> MOV <V><d>, <Vn>.<T>[<index>]
	 * UMOV <Wd>, <Vn>.S[<index>] -> MOV <Wd>, <Vn>.S[<index>]
	 * UMOV <Xd>, <Vn>.D[<index>] -> MOV <Xd>, <Vn>.D[<index>]
	 */

	decode := SimdCopy(i.raw)

	var elemSize1 uint32
	var size uint32
	var dsizeMap = [2]uint32{64, 128}
	var dupRegMap = [5]uint32{REG_W_BASE, REG_W_BASE, REG_W_BASE, REG_X_BASE, REG_X_BASE}
	for ; size < 4; size++ {
		if ((decode.Imm5() >> size) & 1) == 1 {
			break
		}
	}
	elemSize1 = 1 << size
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].OpClass = REG
	if decode.Op() == 0 {
		switch decode.Imm4() {
		case 0:
			i.operation = ARM64_DUP
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
			i.operands[0].ElementSize = elemSize1
			i.operands[0].DataSize = dsizeMap[decode.Q()] / (8 << size)
			i.operands[1].ElementSize = 1 << size
			i.operands[1].Scale = (0x80000000 | (decode.Imm5() >> (size + 1)))
			break
		case 1:
			i.operation = ARM64_DUP
			i.operands[1].Reg[0] = reg(REGSET_ZR, int(dupRegMap[size]), int(decode.Rn()))
			i.operands[0].ElementSize = elemSize1
			i.operands[0].DataSize = dsizeMap[decode.Q()] / (8 << size)
			break
		case 3:
			i.operation = ARM64_INS
			i.operands[1].Reg[0] = reg(REGSET_ZR, int(dupRegMap[size]), int(decode.Rn()))
			i.operands[0].ElementSize = elemSize1
			i.operands[0].Scale = 0x80000000 | (decode.Imm5() >> (size + 1))
			break
		case 5:
			i.operation = ARM64_SMOV
			i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Q()]), int(decode.Rd()))
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
			i.operands[1].ElementSize = elemSize1
			i.operands[1].Scale = 0x80000000 | (decode.Imm5() >> (size + 1))
			if (decode.Q() == 0 && (decode.Imm5()&3) == 0) || (decode.Q() == 1 && (decode.Imm5()&7) == 0) {
				return nil, failedToDecodeInstruction
			}
			break
		case 7:
			i.operation = ARM64_UMOV
			if elemSize1 == uint32(4<<decode.Q()) {
				i.operation = ARM64_MOV
			}
			i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.Q()]), int(decode.Rd()))
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
			i.operands[1].ElementSize = elemSize1
			i.operands[1].Scale = 0x80000000 | (decode.Imm5() >> (size + 1))
			/*printf("Q %d imm5 %d\n", decode.Q, decode.Imm5() )
			if ((decode.Q() == 0 && (decode.Imm5()  & 3) == 0) || (decode.Q() == 1 &&
					(((decode.Imm5()  & 15) == 0) ||
					 ((decode.Imm5()  & 1) == 1) ||
					 ((decode.Imm5()  & 3) == 2) ||
					 ((decode.Imm5()  & 7) == 4))))
				return 1
			*/
			break
		default:
			return nil, failedToDecodeInstruction
		}
	} else {
		i.operation = ARM64_INS
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		i.operands[0].ElementSize = elemSize1
		i.operands[0].Scale = 0x80000000 | (decode.Imm5() >> (size + 1))

		i.operands[1].ElementSize = elemSize1
		i.operands[1].Scale = decode.Imm4() >> size
		if (decode.Imm5() & 15) == 0 {
			return nil, failedToDecodeInstruction
		}
	}
	return i, nil
}

func (i *Instruction) decompose_simd_extract() (*Instruction, error) {
	decode := SimdExtract(i.raw)

	var sizeMap = []uint8{8, 16}

	i.operation = ARM64_EXT
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[3].OpClass = IMM32
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	i.operands[0].ElementSize = 1
	i.operands[1].ElementSize = 1
	i.operands[2].ElementSize = 1
	i.operands[0].DataSize = uint32(sizeMap[decode.Q()])
	i.operands[1].DataSize = uint32(sizeMap[decode.Q()])
	i.operands[2].DataSize = uint32(sizeMap[decode.Q()])
	if decode.Q() == 0 {
		i.operands[3].Immediate = uint64(decode.Imm()) & 7
	} else {
		i.operands[3].Immediate = uint64(decode.Imm())
	}

	if decode.Q() == 0 && decode.Imm() == 1 {
		return nil, failedToDecodeInstruction
	}

	return i, nil
}

func (i *Instruction) decompose_simd_load_store_multiple() (*Instruction, error) {
	/* C4.3.1 Advanced SIMD load/store multiple structures
	 *
	 * LD1/ST1 { <Vt>.<T> },								  [<Xn|SP>]
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T> },						  [<Xn|SP>]
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> },			  [<Xn|SP>]
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>]
	 * LD2/ST2 { <Vt>.<T>, <Vt2>.<T> },						  [<Xn|SP>]
	 * LD3/ST3 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> },			  [<Xn|SP>]
	 * LD4/ST4 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>]
	 */
	decode := SimdLdstMult(i.raw)
	var regCount = []byte{4, 0, 4, 0, 3, 0, 3, 1, 2, 0, 2, 0, 0, 0, 0, 0}
	var elementDataSize = [4][2]byte{{8, 16}, {4, 8}, {2, 4}, {1, 2}}
	var operation = [2][16]Operation{
		{
			ARM64_ST4, ARM64_UNDEFINED, ARM64_ST1, ARM64_UNDEFINED,
			ARM64_ST3, ARM64_UNDEFINED, ARM64_ST1, ARM64_ST1,
			ARM64_ST2, ARM64_UNDEFINED, ARM64_ST1, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		}, {
			ARM64_LD4, ARM64_UNDEFINED, ARM64_LD1, ARM64_UNDEFINED,
			ARM64_LD3, ARM64_UNDEFINED, ARM64_LD1, ARM64_LD1,
			ARM64_LD2, ARM64_UNDEFINED, ARM64_LD1, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		},
	}

	i.operation = operation[decode.L()][decode.Opcode()]
	elements := uint32(regCount[decode.Opcode()])
	i.operands[0].OpClass = MULTI_REG
	for idx := uint32(0); idx < elements; idx++ {
		i.operands[0].Reg[idx] = reg(REGSET_ZR, REG_V_BASE, int((decode.Rt())+idx)%32)
	}

	i.operands[0].DataSize = uint32(elementDataSize[decode.Size()][decode.Q()])
	i.operands[0].ElementSize = 1 << decode.Size()
	i.operands[1].OpClass = MEM_REG
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_load_store_multiple_post_idx() (*Instruction, error) {
	/* C4.3.2 Advanced SIMD load/store multiple structures (post-indexed)
	 *
	 * LD1/ST1 { <Vt>.<T> }								   [<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T> },					   [<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> },			[<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>], <Xm>
	 * LD2/ST2 { <Vt>.<T>, <Vt2>.<T> },					   [<Xn|SP>], <Xm>
	 * LD3/ST3 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> },			[<Xn|SP>], <Xm>
	 * LD4/ST4 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.<T> },								  [<Xn|SP>], <imm>
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T> },					   [<Xn|SP>], <imm>
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> },			[<Xn|SP>], <imm>
	 * LD1/ST1 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>], <imm>
	 * LD2/ST2 { <Vt>.<T>, <Vt2>.<T> },					   [<Xn|SP>], <imm>
	 * LD3/ST3 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> },			[<Xn|SP>], <imm>
	 * LD4/ST4 { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>], <imm>
	 */
	decode := SimdLdstMultPi(i.raw)
	var regCount = []byte{4, 0, 4, 0, 3, 0, 3, 1, 2, 0, 2, 0, 0, 0, 0, 0}
	var elementDataSize = [4][2]byte{{8, 16}, {4, 8}, {2, 4}, {1, 2}}
	var operation = [2][16]Operation{
		{
			ARM64_ST4, ARM64_UNDEFINED, ARM64_ST1, ARM64_UNDEFINED,
			ARM64_ST3, ARM64_UNDEFINED, ARM64_ST1, ARM64_ST1,
			ARM64_ST2, ARM64_UNDEFINED, ARM64_ST1, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		}, {
			ARM64_LD4, ARM64_UNDEFINED, ARM64_LD1, ARM64_UNDEFINED,
			ARM64_LD3, ARM64_UNDEFINED, ARM64_LD1, ARM64_LD1,
			ARM64_LD2, ARM64_UNDEFINED, ARM64_LD1, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		}}
	var imm = [4][2]byte{{8, 16}, {16, 32}, {24, 48}, {32, 64}}
	i.operation = operation[decode.L()][decode.Opcode()]
	elements := uint32(regCount[decode.Opcode()])
	if elements == 0 {
		return nil, failedToDecodeInstruction
	}
	i.operands[0].OpClass = MULTI_REG
	for idx := uint32(0); idx < elements; idx++ {
		i.operands[0].Reg[idx] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rt()+idx)%32)
	}

	i.operands[0].DataSize = uint32(elementDataSize[decode.Size()][decode.Q()])
	i.operands[0].ElementSize = 1 << decode.Size()

	i.operands[1].OpClass = MEM_POST_IDX
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))
	if decode.Rm() == 31 {
		i.operands[1].Immediate = uint64(imm[elements-1][decode.Q()])
		i.operands[1].Reg[1] = uint32(REG_NONE)
	} else {
		i.operands[1].Reg[1] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rm()))
	}

	if (i.operation == ARM64_UNDEFINED) || ((i.operation != ARM64_ST1 && i.operation != ARM64_LD1) && decode.Q() == 0 && decode.Size() == 3) {
		return nil, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_simd_load_store_single() (*Instruction, error) {
	/* C4.3.3  Advanced SIMD load/store single structure
	 *
	 * LD1/ST1 { <Vt>.B }[<index>], [<Xn|SP>]
	 * LD1/ST1 { <Vt>.H }[<index>], [<Xn|SP>]
	 * LD1/ST1 { <Vt>.S }[<index>], [<Xn|SP>]
	 * LD1/ST1 { <Vt>.D }[<index>], [<Xn|SP>]
	 * LD2/ST2 { <Vt>.B, <Vt2>.B }[<index>], [<Xn|SP>]
	 * LD2/ST2 { <Vt>.H, <Vt2>.H }[<index>], [<Xn|SP>]
	 * LD2/ST2 { <Vt>.S, <Vt2>.S }[<index>], [<Xn|SP>]
	 * LD2/ST2 { <Vt>.D, <Vt2>.D }[<index>], [<Xn|SP>]
	 * LD3/ST3 { <Vt>.B, <Vt2>.B, <Vt3>.B }[<index>], [<Xn|SP>]
	 * LD3/ST3 { <Vt>.H, <Vt2>.H, <Vt3>.H }[<index>], [<Xn|SP>]
	 * LD3/ST3 { <Vt>.S, <Vt2>.S, <Vt3>.S }[<index>], [<Xn|SP>]
	 * LD3/ST3 { <Vt>.D, <Vt2>.D, <Vt3>.D }[<index>], [<Xn|SP>]
	 * LD4/ST4 { <Vt>.B, <Vt2>.B, <Vt3>.B, <Vt4>.B }[<index>], [<Xn|SP>]
	 * LD4/ST4 { <Vt>.H, <Vt2>.H, <Vt3>.H, <Vt4>.H }[<index>], [<Xn|SP>]
	 * LD4/ST4 { <Vt>.S, <Vt2>.S, <Vt3>.S, <Vt4>.S }[<index>], [<Xn|SP>]
	 * LD4/ST4 { <Vt>.D, <Vt2>.D, <Vt3>.D, <Vt4>.D }[<index>], [<Xn|SP>]
	 * LD1R { <Vt>.<T> }, [<Xn|SP>]
	 * LD2R { <Vt>.<T>, <Vt2>.<T> }, [<Xn|SP>]
	 * LD3R { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> }, [<Xn|SP>]
	 * LD4R { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>]
	 */
	decode := SimdLdstSingle(i.raw)
	var operation = [4][8]Operation{
		{ARM64_ST1, ARM64_ST3, ARM64_ST1, ARM64_ST3, ARM64_ST1, ARM64_ST3, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_ST2, ARM64_ST4, ARM64_ST2, ARM64_ST4, ARM64_ST2, ARM64_ST4, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_LD1, ARM64_LD3, ARM64_LD1, ARM64_LD3, ARM64_LD1, ARM64_LD3, ARM64_LD1R, ARM64_LD3R},
		{ARM64_LD2, ARM64_LD4, ARM64_LD2, ARM64_LD4, ARM64_LD2, ARM64_LD4, ARM64_LD2R, ARM64_LD4R}}

	var elementMap = [4][8]byte{
		{1, 3, 1, 3, 1, 3, 0, 0},
		{2, 4, 2, 4, 2, 4, 0, 0},
		{1, 3, 1, 3, 1, 3, 1, 3},
		{2, 4, 2, 4, 2, 4, 2, 4},
	}
	i.operation = operation[(decode.L()<<1)+decode.R()][decode.Opcode()]
	i.operands[0].OpClass = MULTI_REG
	elements := uint32(elementMap[(decode.L()<<1)+decode.R()][decode.Opcode()])
	for idx := uint32(0); idx < elements; idx++ {
		i.operands[0].Reg[idx] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rt()+idx)%32)
	}

	var sizemap = [2]uint32{64, 128}
	switch decode.Opcode() >> 1 {
	case 0:
		i.operands[0].ElementSize = 1
		i.operands[0].Index = (decode.Q() << 3) | (decode.S() << 2) | (decode.Size())
		break
	case 1:
		if decode.Size() == 2 || decode.Size() == 0 {
			i.operands[0].ElementSize = 2
			i.operands[0].Index = (decode.Q() << 2) | (decode.S() << 1) | (decode.Size() >> 1)
		} else {
			return nil, failedToDecodeInstruction
		}
		break
	case 2:
		if decode.Size() == 0 {
			i.operands[0].ElementSize = 4
			i.operands[0].Index = (decode.Q() << 1) | decode.S()
		} else if decode.Size() == 1 && decode.S() == 0 {
			i.operands[0].ElementSize = 8
			i.operands[0].Index = decode.Q()
		} else {
			return nil, failedToDecodeInstruction
		}
		break
	case 3:
		i.operands[0].ElementSize = 1 << decode.Size()
		i.operands[0].DataSize = sizemap[decode.Q()] / (8 << decode.Size())
		break
	default:
		return nil, failedToDecodeInstruction
	}
	i.operands[1].OpClass = REG
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_load_store_single_post_idx() (*Instruction, error) {
	/* C4.3.4 Advanced SIMD load/store single structure (post-indexed)
	 *
	 * LD1/ST1 { <Vt>.B }[<index>], [<Xn|SP>], #1
	 * LD1/ST1 { <Vt>.H }[<index>], [<Xn|SP>], #2
	 * LD1/ST1 { <Vt>.S }[<index>], [<Xn|SP>], #4
	 * LD1/ST1 { <Vt>.D }[<index>], [<Xn|SP>], #8
	 * LD1/ST1 { <Vt>.B }[<index>], [<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.H }[<index>], [<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.S }[<index>], [<Xn|SP>], <Xm>
	 * LD1/ST1 { <Vt>.D }[<index>], [<Xn|SP>], <Xm>
	 * LD2/ST2 { <Vt>.B, <Vt2>.B }[<index>], [<Xn|SP>], #2
	 * LD2/ST2 { <Vt>.H, <Vt2>.H }[<index>], [<Xn|SP>], #4
	 * LD2/ST2 { <Vt>.S, <Vt2>.S }[<index>], [<Xn|SP>], #8
	 * LD2/ST2 { <Vt>.D, <Vt2>.D }[<index>], [<Xn|SP>], #16
	 * LD2/ST2 { <Vt>.B, <Vt2>.B }[<index>], [<Xn|SP>], <Xm>
	 * LD2/ST2 { <Vt>.H, <Vt2>.H }[<index>], [<Xn|SP>], <Xm>
	 * LD2/ST2 { <Vt>.S, <Vt2>.S }[<index>], [<Xn|SP>], <Xm>
	 * LD2/ST2 { <Vt>.D, <Vt2>.D }[<index>], [<Xn|SP>], <Xm>
	 * LD3/ST3 { <Vt>.B, <Vt2>.B, <Vt3>.B }[<index>], [<Xn|SP>], #3
	 * LD3/ST3 { <Vt>.H, <Vt2>.H, <Vt3>.H }[<index>], [<Xn|SP>], #6
	 * LD3/ST3 { <Vt>.S, <Vt2>.S, <Vt3>.S }[<index>], [<Xn|SP>], #12
	 * LD3/ST3 { <Vt>.D, <Vt2>.D, <Vt3>.D }[<index>], [<Xn|SP>], #24
	 * LD3/ST3 { <Vt>.B, <Vt2>.B, <Vt3>.B }[<index>], [<Xn|SP>], <Xm>
	 * LD3/ST3 { <Vt>.H, <Vt2>.H, <Vt3>.H }[<index>], [<Xn|SP>], <Xm>
	 * LD3/ST3 { <Vt>.S, <Vt2>.S, <Vt3>.S }[<index>], [<Xn|SP>], <Xm>
	 * LD3/ST3 { <Vt>.D, <Vt2>.D, <Vt3>.D }[<index>], [<Xn|SP>], <Xm>
	 * LD4/ST4 { <Vt>.B, <Vt2>.B, <Vt3>.B, <Vt4>.B }[<index>], [<Xn|SP>], #4
	 * LD4/ST4 { <Vt>.H, <Vt2>.H, <Vt3>.H, <Vt4>.H }[<index>], [<Xn|SP>], #8
	 * LD4/ST4 { <Vt>.S, <Vt2>.S, <Vt3>.S, <Vt4>.S }[<index>], [<Xn|SP>], #16
	 * LD4/ST4 { <Vt>.D, <Vt2>.D, <Vt3>.D, <Vt4>.D }[<index>], [<Xn|SP>], #32
	 * LD4/ST4 { <Vt>.B, <Vt2>.B, <Vt3>.B, <Vt4>.B }[<index>], [<Xn|SP>], <Xm>
	 * LD4/ST4 { <Vt>.H, <Vt2>.H, <Vt3>.H, <Vt4>.H }[<index>], [<Xn|SP>], <Xm>
	 * LD4/ST4 { <Vt>.S, <Vt2>.S, <Vt3>.S, <Vt4>.S }[<index>], [<Xn|SP>], <Xm>
	 * LD4/ST4 { <Vt>.D, <Vt2>.D, <Vt3>.D, <Vt4>.D }[<index>], [<Xn|SP>], <Xm>
	 * LD1R { <Vt>.<T> }, [<Xn|SP>], <imm>
	 * LD1R { <Vt>.<T> }, [<Xn|SP>], <Xm>
	 * LD2R { <Vt>.<T>, <Vt2>.<T> }, [<Xn|SP>], <imm>
	 * LD2R { <Vt>.<T>, <Vt2>.<T> }, [<Xn|SP>], <Xm>
	 * LD3R { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> }, [<Xn|SP>], <imm>
	 * LD3R { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T> }, [<Xn|SP>], <Xm>
	 * LD4R { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>], <imm>
	 * LD4R { <Vt>.<T>, <Vt2>.<T>, <Vt3>.<T>, <Vt4>.<T> }, [<Xn|SP>], <Xm>
	 */
	decode := SimdLdstSinglePi(i.raw)
	var immIdx uint32
	var operation = [4][8]Operation{
		{ARM64_ST1, ARM64_ST3, ARM64_ST1, ARM64_ST3, ARM64_ST1, ARM64_ST3, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_ST2, ARM64_ST4, ARM64_ST2, ARM64_ST4, ARM64_ST2, ARM64_ST4, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_LD1, ARM64_LD3, ARM64_LD1, ARM64_LD3, ARM64_LD1, ARM64_LD3, ARM64_LD1R, ARM64_LD3R},
		{ARM64_LD2, ARM64_LD4, ARM64_LD2, ARM64_LD4, ARM64_LD2, ARM64_LD4, ARM64_LD2R, ARM64_LD4R}}

	var elementMap = [4][8]byte{
		{1, 3, 1, 3, 1, 3, 0, 0},
		{2, 4, 2, 4, 2, 4, 0, 0},
		{1, 3, 1, 3, 1, 3, 1, 3},
		{2, 4, 2, 4, 2, 4, 2, 4},
	}

	var immediateValues = [4][4]byte{
		{1, 2, 4, 8},
		{2, 4, 8, 16},
		{3, 6, 12, 24},
		{4, 8, 16, 32},
	}
	i.operation = operation[(decode.L()<<1)+decode.R()][decode.Opcode()]
	i.operands[0].OpClass = MULTI_REG
	i.operands[1].OpClass = MEM_POST_IDX
	elements := uint32(elementMap[(decode.L()<<1)+decode.R()][decode.Opcode()])
	for idx := uint32(0); idx < elements; idx++ {
		i.operands[0].Reg[idx] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rt()+idx)%32)
	}
	var sizemap = [2]uint32{64, 128}
	switch decode.Opcode() >> 1 {
	case 0:
		i.operands[0].ElementSize = 1
		i.operands[0].Index = (decode.Q() << 3) | (decode.S() << 2) | (decode.Size())
		immIdx = 0
		break
	case 1:
		if decode.Size() == 2 || decode.Size() == 0 {
			i.operands[0].ElementSize = 2
			i.operands[0].Index = (decode.Q() << 2) | (decode.S() << 1) | (decode.Size() >> 1)
			immIdx = 1
		} else {
			return nil, failedToDecodeInstruction
		}
		break
	case 2:
		if decode.Size() == 0 {
			i.operands[0].ElementSize = 4
			i.operands[0].Index = (decode.Q() << 1) | decode.S()
			immIdx = 2
		} else if decode.Size() == 1 && decode.S() == 0 {
			i.operands[0].ElementSize = 8
			i.operands[0].Index = decode.Q()
			immIdx = 3
		} else {
			return nil, failedToDecodeInstruction
		}
		break
	case 3:
		i.operands[0].ElementSize = 1 << decode.Size()
		i.operands[0].DataSize = sizemap[decode.Q()] / (8 << decode.Size())
		break
	default:
		return nil, failedToDecodeInstruction
	}

	if decode.Rm() == 31 && elements != 0 {
		if decode.Opcode()>>1 == 3 {
			i.operands[1].Immediate = uint64(immediateValues[elements-1][decode.Size()])
		} else {
			i.operands[1].Immediate = uint64(immediateValues[elements-1][immIdx])
		}
		i.operands[1].Reg[1] = uint32(REG_NONE)
	} else {
		i.operands[1].Reg[1] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rm()))
	}
	i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Rn()))

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_modified_imm() (*Instruction, error) {
	/* C4.6.4 Advanced SIMD modified immediate
	 *
	 * MOVI <Vd>.<T>, #<imm8>{, LSL #0}
	 * MOVI <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * MOVI <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * MOVI <Vd>.<T>, #<imm8>, MSL #<amount>
	 * ORR  <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * ORR  <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * FMOV <Vd>.<T>, #<imm>
	 * FMOV <Vd>.2D, #<imm>
	 * MVNI <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * MVNI <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * MVNI <Vd>.<T>, #<imm8>, MSL #<amount>
	 * BIC  <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 * BIC  <Vd>.<T>, #<imm8>{, LSL #<amount>}
	 */
	decode := SimdModifiedImm(i.raw)
	type opInfo struct {
		op      Operation
		variant uint32
	}
	var operation = [2][16]opInfo{
		{
			{ARM64_MOVI, 4},
			{ARM64_ORR, 4},
			{ARM64_MOVI, 4},
			{ARM64_ORR, 4},
			{ARM64_MOVI, 4},
			{ARM64_ORR, 4},
			{ARM64_MOVI, 4},
			{ARM64_ORR, 4},
			{ARM64_MOVI, 2},
			{ARM64_ORR, 2},
			{ARM64_MOVI, 2},
			{ARM64_ORR, 2},
			{ARM64_MOVI, 6},
			{ARM64_MOVI, 6},
			{ARM64_MOVI, 1},
			{ARM64_FMOV, 5},
		}, {
			{ARM64_MVNI, 4},
			{ARM64_BIC, 4},
			{ARM64_MVNI, 4},
			{ARM64_BIC, 4},
			{ARM64_MVNI, 4},
			{ARM64_BIC, 4},
			{ARM64_MVNI, 4},
			{ARM64_BIC, 4},
			{ARM64_MVNI, 2},
			{ARM64_BIC, 2},
			{ARM64_MVNI, 2},
			{ARM64_BIC, 2},
			{ARM64_MVNI, 6},
			{ARM64_MVNI, 6},
			{ARM64_MOVI, 8},
			{ARM64_FMOV, 7},
		},
	}
	opinfo := operation[decode.Op()][decode.Cmode()]
	i.operation = opinfo.op
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = IMM32
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	var esize uint32
	var dsize uint32
	var shiftValue uint32
	shiftType := SHIFT_NONE
	immediate := uint64(decode.A()<<7 | decode.B()<<6 | decode.C()<<5 | decode.D()<<4 | decode.E()<<3 | decode.F()<<2 | decode.G()<<1 | decode.H())
	var sign = [2]int32{1, -1}
	var fvalue ieee754
	switch opinfo.variant {
	case 1:
		esize = 1
		dsize = 8 << decode.Q()
		shiftType = SHIFT_LSL
		break
	case 2:
		esize = 2
		dsize = 4 << decode.Q()
		shiftValue = 8 * ((decode.Cmode() & 2) >> 1)
		shiftType = SHIFT_LSL
		break
	case 4:
		esize = 4
		dsize = 2 << decode.Q()
		shiftValue = 8 * ((decode.Cmode() & 6) >> 1)
		shiftType = SHIFT_LSL
		break
	case 5:
		esize = 4
		dsize = 2 << decode.Q()
		i.operands[1].OpClass = FIMM32
		fvalue = fvalue.SetSign(uint32(immediate >> 7))
		fvalue = fvalue.SetExponent(uint32(immediate>>4) & 7)
		fvalue = fvalue.SetFraction(uint32(immediate & 15))
		fvalue = fvalue.SetFloat(float32(sign[fvalue.Sign()] * int32(1<<(fvalue.Exponent()-7)) * int32(1.0+(float32(fvalue.Fraction())/16))))
		immediate = uint64(fvalue.Value())
		break
	case 6:
		esize = 4
		dsize = 2 << decode.Q()
		shiftValue = 8 << (decode.Cmode() & 1)
		shiftType = SHIFT_MSL
		break
	case 7:
		esize = 8
		dsize = 2
		i.operands[1].OpClass = FIMM32
		fvalue = fvalue.SetSign(uint32(immediate & 1))
		fvalue = fvalue.SetExponent(uint32(immediate>>4) & 7)
		fvalue = fvalue.SetFraction(uint32(immediate & 15))
		fvalue = fvalue.SetFloat(float32(sign[fvalue.Sign()] * int32(1<<(fvalue.Exponent()-7)) * int32(1.0+(float32(fvalue.Fraction())/16))))
		immediate = uint64(fvalue.Value())
		break
	case 8:
		if decode.Q() == 1 {
			esize = 8
			dsize = 2
			i.operands[1].OpClass = IMM64
			shiftType = SHIFT_NONE
		} else {
			i.operands[0].Reg[0] = reg(REGSET_ZR, REG_D_BASE, int(decode.Rd()))
			i.operands[1].OpClass = IMM64
			shiftType = SHIFT_NONE
		}
		//ugh this encoding is terrible
		//Is a 64-bit immediate 'aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffffgggggggghhhhhhhh',
		// encoded in "a:b:c:d:e:f:g:h".
		// To do this we pretend that the bit is a sign bit in each byte then right shift them
		var li [8]byte
		li[7] = byte(decode.A() << 7)
		li[6] = byte(decode.B() << 7)
		li[5] = byte(decode.C() << 7)
		li[4] = byte(decode.D() << 7)
		li[3] = byte(decode.E() << 7)
		li[2] = byte(decode.F() << 7)
		li[1] = byte(decode.G() << 7)
		li[0] = byte(decode.H() << 7)
		for idx := uint32(0); idx < 8; idx++ {
			li[idx] >>= 7
		}
		immediate = binary.LittleEndian.Uint64(li[:])
		break
	}
	i.operands[0].ElementSize = esize
	i.operands[0].DataSize = dsize
	i.operands[1].Immediate = uint64(immediate)

	if shiftValue == 0 && shiftType == SHIFT_LSL {
		shiftType = SHIFT_NONE
	}
	i.operands[1].ShiftValue = shiftValue
	i.operands[1].ShiftType = shiftType
	i.operands[1].ShiftValueUsed = 1

	if decode.O2() != 0 {
		return i, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_simd_permute() (*Instruction, error) {
	/* C4.6.5 Advanced SIMD permute
	 *
	 * UZP1 <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * TRN1 <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * ZIP1 <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * UZP2 <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * TRN2 <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 * ZIP2 <Vd>.<T>, <Vn>.<T>, <Vm>.<T>
	 */
	var operation = []Operation{
		ARM64_UNDEFINED,
		ARM64_UZP1,
		ARM64_TRN1,
		ARM64_ZIP1,
		ARM64_UNDEFINED,
		ARM64_UZP2,
		ARM64_TRN2,
		ARM64_ZIP2,
	}
	decode := SimdPermute(i.raw)
	var esizeMap = [2]uint8{64, 128}
	i.operation = operation[decode.Opcode()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	i.operands[0].ElementSize = 1 << decode.Size()
	i.operands[0].DataSize = uint32(esizeMap[decode.Q()]) / (8 << decode.Size())
	i.operands[1].ElementSize = 1 << decode.Size()
	i.operands[1].DataSize = uint32(esizeMap[decode.Q()]) / (8 << decode.Size())
	i.operands[2].ElementSize = 1 << decode.Size()
	i.operands[2].DataSize = uint32(esizeMap[decode.Q()]) / (8 << decode.Size())

	if i.operation == ARM64_UNDEFINED || (decode.Size() == 3 && decode.Q() == 0) {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_scalar_2_reg_misc() (*Instruction, error) {
	/* C4.6.11 Advanced SIMD scalar two-register miscellaneous
	 *
	 * SUQADD  <V><d>, <V><n>
	 * SQABS   <V><d>, <V><n>
	 * CMGT	<V><d>, <V><n>, #0
	 * CMEQ	<V><d>, <V><n>, #0
	 * CMLT	<V><d>, <V><n>, #0
	 * ABS	 <V><d>, <V><n>
	 * SQXTN   <Vb><d>, <Va><n>
	 * FCVTNS  <V><d>, <V><n>
	 * FCVTMS  <V><d>, <V><n>
	 * FCVTAS  <V><d>, <V><n>
	 * SCVTF   <V><d>, <V><n>
	 * FCMGT   <V><d>, <V><n>, #0.0
	 * FCMEQ   <V><d>, <V><n>, #0.0
	 * FCMLT   <V><d>, <V><n>, #0.0
	 * FCVTPS  <V><d>, <V><n>
	 * FCVTZS  <V><d>, <V><n>
	 * FRECPE  <V><d>, <V><n>
	 * FRECPX  <V><d>, <V><n>
	 * FRECPX  <V><d>, <V><n>
	 * USQADD  <V><d>, <V><n>
	 * SQNEG   <V><d>, <V><n>
	 * CMGE	<V><d>, <V><n>, #0
	 * CMLE	<V><d>, <V><n>, #0
	 * NEG	 <V><d>, <V><n>
	 * SQXTUN  <Vb><d>, <Va><n>
	 * UQXTN   <Vb><d>, <Va><n>
	 * FCVTXN  <Vb><d>, <Va><n>
	 * FCVTNU  <V><d>, <V><n>
	 * FCVTMU  <V><d>, <V><n>
	 * FCVTAU  <V><d>, <V><n>
	 * UCVTF   <V><d>, <V><n>
	 * FCMGE   <V><d>, <V><n>, #0.0
	 * FCMLE   <V><d>, <V><n>, #0.0
	 * FCVTPU  <V><d>, <V><n>
	 * FCVTZU  <V><d>, <V><n>
	 * FRSQRTE <V><d>, <V><n>
	 *
	 * 0: <V><d>,  <V><n>
	 * 1: <V><d>,  <V><n>, #0.0
	 * 2: <Vb><d>, <Va><n>
	 * 3: (decode.Size() < 2)->  <V><d>, <V><n>
	 * 4: (decode.Size() < 2)->  <V><d>, <V><n>, #0.0
	 * 5: (decode.Size() < 2)->  <Vb><d>, <Va><n>
	 * 6: (decode.Size() > 1)->  <V><d>, <V><n>
	 * 7: (decode.Size() > 1)->  <V><d>, <V><n>, #0.0
	 * 8: (decode.Size() > 2)->  <Vb><d>, <Va><n>
	 */
	/*var operation = []Operation{*/
	/*0 - xx - 3  - v: 0 {ARM64_SUQADD,  0},*/
	/*0 - xx - 7  - v: 0 {ARM64_SQABS,   0},*/
	/*0 - xx - 8  - v: 1 {ARM64_CMGT,	1},*/
	/*0 - xx - 9  - v: 1 {ARM64_CMEQ,	1},*/
	/*0 - xx - 10 - v: 1 {ARM64_CMLT,	1},*/
	/*0 - xx - 11 - v: 0 {ARM64_ABS,	 0},*/
	/*0 - xx - 20 - v: 2 {ARM64_SQXTN,   2},*/
	/*0 - 0x - 26 - v: 3 {ARM64_FCVTNS,  3},*/
	/*0 - 0x - 27 - v: 3 {ARM64_FCVTMS,  3},*/
	/*0 - 0x - 28 - v: 3 {ARM64_FCVTAS,  3},*/
	/*0 - 0x - 29 - v: 3 {ARM64_SCVTF,   3},*/
	/*0 - 1x - 12 - v: 7 {ARM64_FCMGT,   7},*/
	/*0 - 1x - 13 - v: 7 {ARM64_FCMEQ,   7},*/
	/*0 - 1x - 14 - v: 7 {ARM64_FCMLT,   7},*/
	/*0 - 1x - 26 - v: 6 {ARM64_FCVTPS,  6},*/
	/*0 - 1x - 27 - v: 6 {ARM64_FCVTZS,  6},*/
	/*0 - 1x - 29 - v: 6 {ARM64_FRECPE,  6},*/
	/*0 - 1x - 31 - v: 6 {ARM64_FRECPX,  6},*/
	/*1 - xx - 3  - v: 0 {ARM64_USQADD,  0},*/
	/*1 - xx - 7  - v: 0 {ARM64_SQNEG,   0},*/
	/*1 - xx - 8  - v: 1 {ARM64_CMGE,	1},*/
	/*1 - xx - 9  - v: 1 {ARM64_CMLE,	1},*/
	/*1 - xx - 11 - v: 0 {ARM64_NEG,	 0},*/
	/*1 - xx - 18 - v: 2 {ARM64_SQXTUN,  2},*/
	/*1 - xx - 20 - v: 2 {ARM64_UQXTN,   2},*/
	/*1 - 0x - 22 - v: 5 {ARM64_FCVTXN,  5},*/
	/*1 - 0x - 26 - v: 3 {ARM64_FCVTNU,  3},*/
	/*1 - 0x - 27 - v: 3 {ARM64_FCVTMU,  3},*/
	/*1 - 0x - 28 - v: 6 {ARM64_FCVTAU,  6},*/
	/*1 - 0x - 29 - v: 6 {ARM64_UCVTF,   6},*/
	/*1 - 1x - 12 - v: 7 {ARM64_FCMGE,   7},*/
	/*1 - 1x - 13 - v: 7 {ARM64_FCMLE,   7},*/
	/*1 - 1x - 26 - v: 6 {ARM64_FCVTPU,  6},*/
	/*1 - 1x - 27 - v: 6 {ARM64_FCVTZU,  6},*/
	/*1 - 1x - 29 - v: 6 {ARM64_FRSQRTE, 6},*/
	/*};*/
	type OpInfo struct {
		op   Operation
		ovar uint32
	}
	var operation = [2][2][32]OpInfo{
		{
			{
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SUQADD, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQABS, 0},
				{ARM64_CMGT, 1}, {ARM64_CMEQ, 1}, {ARM64_CMLT, 1}, {ARM64_ABS, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_SQXTN, 2}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FCVTNS, 3}, {ARM64_FCVTMS, 3},
				{ARM64_FCVTAS, 3}, {ARM64_SCVTF, 3}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
			}, {
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SUQADD, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQABS, 0},
				{ARM64_CMGT, 1}, {ARM64_CMEQ, 1}, {ARM64_CMLT, 1}, {ARM64_ABS, 0},
				{ARM64_FCMGT, 7}, {ARM64_FCMEQ, 7}, {ARM64_FCMLT, 7}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_SQXTN, 2}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FCVTPS, 6}, {ARM64_FCVTZS, 6},
				{ARM64_UNDEFINED, 0}, {ARM64_FRECPE, 6}, {ARM64_UNDEFINED, 0}, {ARM64_FRECPX, 6},
			},
		}, {
			{
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_USQADD, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQNEG, 0},
				{ARM64_CMGE, 1}, {ARM64_CMLE, 1}, {ARM64_UNDEFINED, 0}, {ARM64_NEG, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQXTUN, 2}, {ARM64_UNDEFINED, 0},
				{ARM64_UQXTN, 2}, {ARM64_UNDEFINED, 0}, {ARM64_FCVTXN, 5}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FCVTNU, 3}, {ARM64_FCVTMU, 3},
				{ARM64_FCVTAU, 6}, {ARM64_UCVTF, 6}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
			}, {
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_USQADD, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQNEG, 0},
				{ARM64_CMGE, 1}, {ARM64_CMLE, 1}, {ARM64_UNDEFINED, 0}, {ARM64_NEG, 0},
				{ARM64_FCMGE, 7}, {ARM64_FCMLE, 7}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQXTUN, 2}, {ARM64_UNDEFINED, 0},
				{ARM64_UQXTN, 2}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FCVTPU, 6}, {ARM64_FCVTZU, 6},
				{ARM64_UNDEFINED, 0}, {ARM64_FRSQRTE, 6}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
			},
		},
	}
	decode := SimdScalar2RegisterMisc(i.raw)
	opinfo := operation[decode.U()][decode.Size()>>1][decode.Opcode()]
	var regbase = [4]uint8{REG_B_BASE, REG_H_BASE, REG_S_BASE, REG_D_BASE}
	var regbase2 = [3]uint8{REG_H_BASE, REG_S_BASE, REG_D_BASE}
	var regbase3 = [2]uint8{REG_S_BASE, REG_D_BASE}
	i.operation = opinfo.op
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	switch opinfo.ovar {
	case 0:
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rn()))
		break
	case 6:
		fallthrough
	case 3:
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase3[decode.Size()&1]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase3[decode.Size()&1]), int(decode.Rn()))
		break
	case 1:
		fallthrough
	case 4:
		if decode.Size() != 3 {
			return nil, failedToDecodeInstruction
		}
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rn()))
		i.operands[2].OpClass = IMM32
		i.operands[2].Immediate = 0
		break
	case 7:
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase3[decode.Size()&1]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase3[decode.Size()&1]), int(decode.Rn()))
		i.operands[2].OpClass = IMM32
		i.operands[2].Immediate = 0
		break
	case 2:
		fallthrough
	case 8:
		if decode.Size() == 3 {
			return nil, failedToDecodeInstruction
		}
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
		break
	case 5:
		if (decode.Size() & 1) == 0 {
			return nil, failedToDecodeInstruction
		}
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_S_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_D_BASE, int(decode.Rn()))
		break
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_scalar_3_different() (*Instruction, error) {
	/* C4.6.9 Advanced SIMD scalar three different
	 *
	 * SQDMLAL <Va><d>, <Vb><n>, <Vb><m>
	 * SQDMLSL <Va><d>, <Vb><n>, <Vb><m>
	 * SQDMULL <Va><d>, <Vb><n>, <Vb><m>
	 */
	var operation = []Operation{
		ARM64_UNDEFINED,
		ARM64_SQDMLAL,
		ARM64_UNDEFINED,
		ARM64_SQDMLSL,
		ARM64_UNDEFINED,
		ARM64_SQDMULL,
		ARM64_UNDEFINED,
		ARM64_UNDEFINED,
	}
	decode := SimdScalar3Different(i.raw)
	var regbase1 = [4]uint8{0, REG_S_BASE, REG_D_BASE, 0}
	var regbase2 = [4]uint8{0, REG_H_BASE, REG_S_BASE, 0}
	i.operation = operation[decode.Opcode()&7]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase1[decode.Size()]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rm()))

	if decode.Opcode() < 7 || decode.Size() == 0 || decode.Size() == 3 || i.operation == ARM64_UNDEFINED {
		return i, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_simd_scalar_3_same() (*Instruction, error) {
	/* C4.6.10 Advanced SIMD scalar three same
	 *
	 * SQADD	<V><d>, <V><n>, <V><m>
	 * SQSUB	<V><d>, <V><n>, <V><m>
	 * CMGT	 <V><d>, <V><n>, <V><m>
	 * CMGE	 <V><d>, <V><n>, <V><m>
	 * SSHL	 <V><d>, <V><n>, <V><m>
	 * SQSHL	<V><d>, <V><n>, <V><m>
	 * SQRSHL   <V><d>, <V><n>, <V><m>
	 * ADD	  <V><d>, <V><n>, <V><m>
	 * CMTST	<V><d>, <V><n>, <V><m>
	 * SQDMULH  <V><d>, <V><n>, <V><m>
	 * FMULX	<V><d>, <V><n>, <V><m>
	 * FCMEQ	<V><d>, <V><n>, <V><m>
	 * FRECPS   <V><d>, <V><n>, <V><m>
	 * FRSQRTS  <V><d>, <V><n>, <V><m>
	 * UQADD	<V><d>, <V><n>, <V><m>
	 * UQSUB	<V><d>, <V><n>, <V><m>
	 * CMHI	 <V><d>, <V><n>, <V><m>
	 * CMHS	 <V><d>, <V><n>, <V><m>
	 * USHL	 <V><d>, <V><n>, <V><m>
	 * UQRSHL   <V><d>, <V><n>, <V><m>
	 * URSHL	<V><d>, <V><n>, <V><m>
	 * SUB	  <V><d>, <V><n>, <V><m>
	 * CMEQ	 <V><d>, <V><n>, <V><m>
	 * CMEQ	 <V><d>, <V><n>, <V><m>
	 * SQRDMULH <V><d>, <V><n>, <V><m>
	 * FCMGE	<V><d>, <V><n>, <V><m>
	 * FACGE	<V><d>, <V><n>, <V><m>
	 * FABD	 <V><d>, <V><n>, <V><m>
	 * FCMGT	<V><d>, <V><n>, <V><m>
	 * FACGT	<V><d>, <V><n>, <V><m>
	 *
	 * Variants:
	 * 0: BHSD
	 * 1: D
	 * 2: HS
	 * 3: SD
	 */
	type OpInfo struct {
		op   Operation
		ovar uint32
	}
	/*var operation = []Operation{*/
	/*0 - xx - 0 	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 1 	{ARM64_SQADD,	 0},*/
	/*0 - xx - 2 	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 3 	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 4 	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 5 	{ARM64_SQSUB,	 0},*/
	/*0 - xx - 6 	{ARM64_CMGT,	  1},*/
	/*0 - xx - 7 	{ARM64_CMGE,	  1},*/
	/*0 - xx - 8 	{ARM64_SSHL,	  1},*/
	/*0 - xx - 9 	{ARM64_SQSHL,	 0},*/
	/*0 - xx - 10 	{ARM64_SRSHL,	 0},*/
	/*0 - xx - 11	{ARM64_SQRSHL,	0},*/
	/*0 - xx - 12	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 13	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 14	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 15	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 16	{ARM64_ADD,	   1},*/
	/*0 - xx - 17	{ARM64_CMTST,	 1},*/
	/*0 - xx - 18	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 19	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 20	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 21	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 22	{ARM64_SQDMULH,   3},*/
	/*0 - xx - 23	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 24	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 25	{ARM64_UNDEFINED, 0},*/
	/*0 - xx - 26	{ARM64_UNDEFINED, 0},*/
	/*0 - 0x - 27	{ARM64_FMULX,	 3},*/
	/*0 - 0x - 28	{ARM64_FCMEQ,	 3},*/
	/*0 - xx - 30	{ARM64_UNDEFINED, 0},*/
	/*0 - 0x - 31	{ARM64_FRECPS,	3},*/
	/*0 - 1x - 31	{ARM64_FRSQRTS,   3},*/

	/*1 - xx - 0 	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 1 	{ARM64_UQADD,	 0},*/
	/*1 - xx - 2 	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 3 	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 4 	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 5 	{ARM64_UQSUB,	 0},*/
	/*1 - xx - 6 	{ARM64_CMHI,	  1},*/
	/*1 - xx - 7 	{ARM64_CMHS,	  1},*/
	/*1 - xx - 8 	{ARM64_USHL,	  1},*/
	/*1 - xx - 9 	{ARM64_UQRSHL,	0},*/
	/*1 - xx - 10	{ARM64_URSHL,	 1},*/
	/*1 - xx - 11	{ARM64_UQRSHL,	0},*/
	/*1 - xx - 12	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 13	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 14	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 15	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 16	{ARM64_SUB,	   1},*/
	/*1 - xx - 17	{ARM64_CMEQ,	  1},*/
	/*1 - xx - 18	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 19	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 20	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 21	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 22	{ARM64_SQRDMULH,  2},*/
	/*1 - xx - 23	{ARM64_UNDEFINED, 0},*/
	/*1 - 0x - 24	{ARM64_FCMGE,	 3},*/
	/*1 - 0x - 25	{ARM64_FACGE,	 3},*/
	/*1 - 1x - 26	{ARM64_FABD,	  3},*/
	/*1 - 1x - 27	{ARM64_UNDEFINED, 0},*/
	/*1 - 1x - 28	{ARM64_FCMGT,	 3},*/
	/*1 - 1x - 29	{ARM64_FACGT,	  3},*/
	/*1 - xx - 30	{ARM64_UNDEFINED, 0},*/
	/*1 - xx - 31	{ARM64_UNDEFINED, 0},*/
	/*};*/
	decode := SimdScalar3Same(i.raw)
	var operation = [2][2][32]OpInfo{
		{
			{
				{ARM64_UNDEFINED, 0}, {ARM64_SQADD, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_SQSUB, 0}, {ARM64_CMGT, 1}, {ARM64_CMGE, 1},
				{ARM64_SSHL, 1}, {ARM64_SQSHL, 0}, {ARM64_SRSHL, 0}, {ARM64_SQRSHL, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_ADD, 1}, {ARM64_CMTST, 1}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQDMULH, 2}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FMULX, 3},
				{ARM64_FCMEQ, 3}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FRECPS, 3},
			}, {
				{ARM64_UNDEFINED, 0}, {ARM64_SQADD, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_SQSUB, 0}, {ARM64_CMGT, 1}, {ARM64_CMGE, 1},
				{ARM64_SSHL, 1}, {ARM64_SQSHL, 0}, {ARM64_SRSHL, 0}, {ARM64_SQRSHL, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_ADD, 1}, {ARM64_CMTST, 1}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQDMULH, 2}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FMULX, 3},
				{ARM64_FCMEQ, 3}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FRSQRTS, 3},
			},
		}, {
			{
				{ARM64_UNDEFINED, 0}, {ARM64_UQADD, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UQSUB, 0}, {ARM64_CMHI, 1}, {ARM64_CMHS, 1},
				{ARM64_USHL, 1}, {ARM64_UQSHL, 0}, {ARM64_URSHL, 1}, {ARM64_UQRSHL, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_SUB, 1}, {ARM64_CMEQ, 1}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQRDMULH, 2}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_FCMGE, 3}, {ARM64_FACGE, 3}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
			}, {
				{ARM64_UNDEFINED, 0}, {ARM64_UQADD, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UQSUB, 0}, {ARM64_CMHI, 1}, {ARM64_CMHS, 1},
				{ARM64_USHL, 1}, {ARM64_UQSHL, 0}, {ARM64_URSHL, 1}, {ARM64_UQRSHL, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_SUB, 1}, {ARM64_CMEQ, 1}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_SQRDMULH, 2}, {ARM64_UNDEFINED, 0},
				{ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0}, {ARM64_FABD, 3}, {ARM64_UNDEFINED, 0},
				{ARM64_FCMGT, 3}, {ARM64_FACGT, 3}, {ARM64_UNDEFINED, 0}, {ARM64_UNDEFINED, 0},
			},
		},
	}

	opinfo := operation[decode.U()][decode.Size()>>1][decode.Opcode()]
	var regbase = [4][4]uint8{
		{REG_B_BASE, REG_H_BASE, REG_S_BASE, REG_D_BASE},
		{REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE},
		{REG_S_BASE, REG_H_BASE, REG_S_BASE, REG_H_BASE},
		{REG_S_BASE, REG_D_BASE, REG_S_BASE, REG_D_BASE},
	}
	i.operation = opinfo.op
	//printf("U: %d s: %d opc: %d - str: %s\n", decode.U, decode.Size()>>1, decode.Opcode(), OperationString[opinfo.op])
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[opinfo.ovar][decode.Size()]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase[opinfo.ovar][decode.Size()]), int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, int(regbase[opinfo.ovar][decode.Size()]), int(decode.Rm()))

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_scalar_copy() (*Instruction, error) {
	/* C4.6.6 Advanced SIMD scalar copy
	 *
	 * DUP <V><d>, <Vn>.<T>[<index>]
	 */
	decode := SimdScalarCopy(i.raw)
	var size uint32
	if decode.Imm5() == 0 || decode.Imm5() == 16 {
		return nil, failedToDecodeInstruction
	}
	for ; size < 4; size++ {
		if ((decode.Imm5() >> size) & 1) == 1 {
			break
		}
	}
	i.operation = ARM64_MOV
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	var regset = [4]uint8{REG_B_BASE, REG_H_BASE, REG_S_BASE, REG_D_BASE}
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regset[size]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[1].ElementSize = 1 << size
	i.operands[1].Scale = 0x80000000 | (decode.Imm5() >> (1 + size))

	if decode.Op() != 0 || decode.Imm4() != 0 {
		return i, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_simd_scalar_indexed_element() (*Instruction, error) {
	/* C4.6.12 Advanced SIMD scalar x indexed element
	 *
	 * SQDMLAL  <Va><d>, <Vb><n>, <Vm>.<Ts>[<index>]
	 * SQDMLSL  <Va><d>, <Vb><n>, <Vm>.<Ts>[<index>]
	 * SQDMULL  <Va><d>, <Vb><n>, <Vm>.<Ts>[<index>]
	 * SQDMULH  <V><d>,  <V><n>,  <Vm>.<Ts>[<index>]
	 * SQRDMULH <V><d>,  <V><n>,  <Vm>.<Ts>[<index>]
	 * FMLA	 <V><d>,  <V><n>,  <Vm>.<Ts>[<index>]
	 * FMLS	 <V><d>,  <V><n>,  <Vm>.<Ts>[<index>]
	 * FMUL	 <V><d>,  <V><n>,  <Vm>.<Ts>[<index>]
	 * FMULX	<V><d>,  <V><n>,  <Vm>.<Ts>[<index>]
	 */
	decode := SimdScalarXIndexedElement(i.raw)
	var index uint32
	hireg := decode.M() << 4
	var regbase = [4]uint32{0, REG_S_BASE, REG_D_BASE, 0}
	var regbase2 = [4]uint32{0, REG_H_BASE, REG_S_BASE, 0}
	var regbase3 = [2]uint32{REG_S_BASE, REG_D_BASE}

	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase3[decode.Size()&1]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase3[decode.Size()&1]), int(decode.Rn()))

	if decode.Size() == 1 {
		index = decode.H()<<2 | decode.L()<<1 | decode.M()
	} else if decode.Size() == 2 {
		index = decode.H()<<1 | decode.L()
	}
	i.operands[2].ElementSize = 4 << (decode.Size() & 1)
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(hireg|decode.Rm()))
	switch decode.Opcode() {
	case 1:
		if decode.Size() < 2 {
			return nil, failedToDecodeInstruction
		}
		i.operation = ARM64_FMLA
		if (decode.Size() & 1) == 0 {
			i.operands[2].Scale = 0x80000000 | decode.H()<<1 | decode.L()
		} else if (decode.Size()&1) == 1 && decode.L() == 0 {
			i.operands[2].Scale = 0x80000000 | decode.H()
		}
		break
	case 3:
		i.operation = ARM64_SQDMLAL
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
		if decode.Size() == 1 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		} else if decode.Size() == 2 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int((decode.M()<<4)|decode.Rm()))
		}
		i.operands[2].Scale = 0x80000000 | index
		i.operands[2].ElementSize = 1 << decode.Size()
		break
	case 5:
		if decode.Size() < 2 {
			return nil, failedToDecodeInstruction
		}
		i.operation = ARM64_FMLS
		if (decode.Size() & 1) == 0 {
			i.operands[2].Scale = 0x80000000 | decode.H()<<1 | decode.L()
		} else if (decode.Size()&1) == 1 && decode.L() == 0 {
			i.operands[2].Scale = 0x80000000 | decode.H()
		}
		break
	case 7:
		i.operation = ARM64_SQDMLSL
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
		if decode.Size() == 1 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		} else if decode.Size() == 2 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int((decode.M()<<4)|decode.Rm()))
		}
		i.operands[2].Scale = 0x80000000 | index
		i.operands[2].ElementSize = 1 << decode.Size()
		break
	case 9:
		if decode.Size() < 2 {
			return nil, failedToDecodeInstruction
		}
		if decode.U() == 0 {
			i.operation = ARM64_FMUL
		} else {
			i.operation = ARM64_FMULX
		}
		if (decode.Size() & 1) == 0 {
			i.operands[2].Scale = 0x80000000 | decode.H()<<1 | decode.L()
		} else if (decode.Size()&1) == 1 && decode.L() == 0 {
			i.operands[2].Scale = 0x80000000 | decode.H()
		}
		break
	case 11:
		i.operation = ARM64_SQDMULL
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
		if decode.Size() == 1 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		} else if decode.Size() == 2 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int((decode.M()<<4)|decode.Rm()))
		}
		i.operands[2].Scale = 0x80000000 | index
		i.operands[2].ElementSize = 1 << decode.Size()
		break
	case 12:
		i.operation = ARM64_SQDMULH
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
		if decode.Size() == 1 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		} else if decode.Size() == 2 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int((decode.M()<<4)|decode.Rm()))
		}
		i.operands[2].Scale = 0x80000000 | index
		i.operands[2].ElementSize = 1 << decode.Size()
		break
	case 13:
		i.operation = ARM64_SQRDMULH
		i.operands[0].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, int(regbase2[decode.Size()]), int(decode.Rn()))
		if decode.Size() == 1 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
		} else if decode.Size() == 2 {
			i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int((decode.M()<<4)|decode.Rm()))
		}
		i.operands[2].Scale = 0x80000000 | index
		i.operands[2].ElementSize = 1 << decode.Size()
		break
	}

	return i, nil
}

func (i *Instruction) decompose_simd_scalar_pairwise() (*Instruction, error) {
	/* C4.6.7 Advanced SIMD scalar pairwise
	 *
	 * ADDP	<V><d>, <Vn>.<T>
	 * FMAXNMP <V><d>, <Vn>.<T>
	 * FADDP   <V><d>, <Vn>.<T>
	 * FADDP   <V><d>, <Vn>.<T>
	 * FMAXP   <V><d>, <Vn>.<T>
	 * FMINNMP <V><d>, <Vn>.<T>
	 * FMINP   <V><d>, <Vn>.<T>
	 */
	decode := SimdScalarPairwise(i.raw)
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	var regset = [2]uint8{REG_S_BASE, REG_D_BASE}
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regset[decode.Size()&1]), int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[1].ElementSize = 4 << decode.Size()
	i.operands[1].DataSize = 2
	switch decode.Opcode() {
	case 27:
		if decode.U() == 0 {
			i.operation = ARM64_ADDP
			if decode.Size() != 3 {
				return nil, failedToDecodeInstruction
			}
			i.operands[0].Reg[0] = reg(REGSET_ZR, REG_D_BASE, int(decode.Rd()))
			i.operands[1].ElementSize = 8
			i.operands[1].DataSize = 2
		}
		break
	case 12:
		if decode.U() == 1 {
			if decode.Size() < 2 {
				i.operation = ARM64_FMAXNMP
			} else {
				i.operation = ARM64_FMINNMP
			}
		} else {
			return nil, failedToDecodeInstruction
		}
		break
	case 13:
		i.operation = ARM64_FADDP
		break
	case 15:
		if decode.U() == 1 {
			if decode.Size() < 2 {
				i.operation = ARM64_FMAXP
			} else {
				i.operation = ARM64_FMINP
			}
		} else {
			return nil, failedToDecodeInstruction
		}
		break
	default:
		return nil, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_simd_scalar_shift_imm() (*Instruction, error) {
	/* C4.6.8 Advanced SIMD scalar shift by immediate
	 *
	 * SSHR	 <V><d>, <V><n>, #<shift>
	 * SSRA	 <V><d>, <V><n>, #<shift>
	 * SRSHR	<V><d>, <V><n>, #<shift>
	 * SRSRA	<V><d>, <V><n>, #<shift>
	 * SHL	  <V><d>, <V><n>, #<shift>
	 * SQSHL	<V><d>, <V><n>, #<shift>
	 * SQSHRN   <Vb><d>, <Va><n>, #<shift>
	 * SQRSHRN  <Vb><d>, <Va><n>, #<shift>
	 * SCVTF	<V><d>, <V><n>, #<fbits>
	 * FCVTZS   <V><d>, <V><n>, #<fbits>
	 * USHR	 <V><d>, <V><n>, #<shift>
	 * USRA	 <V><d>, <V><n>, #<shift>
	 * URSHR	<V><d>, <V><n>, #<shift>
	 * URSRA	<V><d>, <V><n>, #<shift>
	 * SRI	  <V><d>, <V><n>, #<shift>
	 * SLI	  <V><d>, <V><n>, #<shift>
	 * SQSHLU   <V><d>, <V><n>, #<shift>
	 * UQSHL	<V><d>, <V><n>, #<shift>
	 * SQSHRUN  <Vb><d>, <Va><n>, #<shift>
	 * SQRSHRUN <Vb><d>, <Va><n>, #<shift>
	 * UQSHRN   <Vb><d>, <Va><n>, #<shift>
	 * UQRSHRN  <Vb><d>, <Va><n>, #<shift>
	 * UCVTF	<V><d>, <V><n>, #<fbits>
	 * FCVTZU   <V><d>, <V><n>, #<fbits>
	 */

	decode := SimdShiftByImm(i.raw)

	type shiftCalc int
	const (
		SH_3 shiftCalc = iota
		SH_HB
		SH_IH
	)
	type decodeOperation struct {
		op      Operation
		esize   uint32
		calc    shiftCalc
		regBase uint32
	}
	var operation = [2][32]decodeOperation{
		{
			{ARM64_SSHR, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SSRA, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SRSHR, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SRSRA, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SHL, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SQSHL, 8, SH_HB, 1}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SQSHRN, 8, SH_HB, 2}, {ARM64_SQRSHRN, 8, SH_HB, 2},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SCVTF, 32, SH_IH, 3}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_FCVTZS, 32, SH_IH, 3},
		}, {
			{ARM64_USHR, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_USRA, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_URSHR, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_URSRA, 8, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SRI, 8, SH_3, 1}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SLI, 8, SH_3, 1}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SQSHLU, 8, SH_HB, 1}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UQSHL, 8, SH_HB, 1}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_SQSHRUN, 8, SH_HB, 2}, {ARM64_SQRSHRUN, 8, SH_HB, 2},
			{ARM64_UQSHRN, 8, SH_HB, 2}, {ARM64_UQRSHRN, 8, SH_HB, 2},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UCVTF, 32, SH_IH, 3}, {ARM64_UNDEFINED, 0, SH_3, 0},
			{ARM64_UNDEFINED, 0, SH_3, 0}, {ARM64_FCVTZU, 32, SH_IH, 3},
		},
	}

	var regBaseMap = [4][16]Register{
		{
			REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE,
			REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE,
		}, {
			REG_NONE, REG_B_BASE, REG_H_BASE, REG_H_BASE, REG_S_BASE, REG_S_BASE, REG_S_BASE, REG_S_BASE,
			REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE,
		}, {
			REG_NONE, REG_B_BASE, REG_H_BASE, REG_H_BASE, REG_S_BASE, REG_S_BASE, REG_S_BASE, REG_S_BASE,
			REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_NONE,
		}, {
			REG_NONE, REG_NONE, REG_NONE, REG_NONE, REG_S_BASE, REG_S_BASE, REG_S_BASE, REG_S_BASE,
			REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE, REG_D_BASE,
		},
	}

	decodeOp := operation[decode.U()][decode.Opcode()]
	regBase := uint32(regBaseMap[decodeOp.regBase][decode.Immh()])
	if regBase == uint32(REG_NONE) {
		return nil, failedToDecodeInstruction
	}
	i.operation = decodeOp.op
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regBase), int(decode.Rd()))

	i.operands[1].OpClass = REG
	var regOffset uint32
	if decode.Opcode() >= 16 && decode.Opcode() <= 19 {
		regOffset = 1
	}
	i.operands[1].Reg[0] = reg(REGSET_ZR, int(regBase+regOffset), int(decode.Rn()))

	var esize uint32
	switch decodeOp.calc {
	case SH_3:
		esize = decodeOp.esize << 3
		break
	case SH_HB:
		esize = decodeOp.esize << bits.Len32(decode.Immh())
		break
	case SH_IH:
		esize = decodeOp.esize << ((decode.Immh() >> 3) & 1)
		break
	}
	if decode.Opcode() == 10 || decode.Opcode() == 14 || i.operation == ARM64_SQSHLU {
		i.operands[2].Immediate = uint64(((decode.Immh() << 3) | decode.Immb()) - (esize))
	} else {
		i.operands[2].Immediate = uint64((esize * 2) - ((decode.Immh() << 3) | decode.Immb()))
	}
	i.operands[2].OpClass = IMM32
	i.operands[2].SignedImm = decode.U()

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_shift_imm() (*Instruction, error) {
	/* C4.6.13 Advanced SIMD shift by immediate
	 *
	 * SSHR		<Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SSRA		<Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SRSHR	   <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SRSRA	   <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SHL		 <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SQSHL	   <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SHRN{2}	 <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * RSHRN{2}	<Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * SQSHRN{2}   <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * SQRSHRN{2}  <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * SSHLL{2}	<Vd>.<Ta>, <Vn>.<Tb>, #<shift>
	 * SCVTF	   <Vd>.<T>,  <Vn>.<T>,  #<fbits>
	 * FCVTZS	  <Vd>.<T>,  <Vn>.<T>,  #<fbits>
	 * USHR		<Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * USRA		<Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * URSHR	   <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * URSRA	   <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SRI		 <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SLI		 <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SQSHLU	  <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * UQSHL	   <Vd>.<T>,  <Vn>.<T>,  #<shift>
	 * SQSHRUN{2}  <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * SQRSHRUN{2} <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * UQSHRN{2}   <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * UQRSHRN{2}  <Vd>.<Tb>, <Vn>.<Ta>, #<shift>
	 * USHLL{2}	<Vd>.<Ta>, <Vn>.<Tb>, #<shift>
	 * UCVTF	   <Vd>.<T>,  <Vn>.<T>,  #<fbits>
	 * FCVTZU	  <Vd>.<T>,  <Vn>.<T>,  #<fbits>
	 *
	 * Alias
	 * SSHLL{2} <Vd>.<Ta>, <Vn>.<Tb>, #0 -> SCVTF <V><d>, <V><n>, #<fbits>
	 * USHLL{2} <Vd>.<Ta>, <Vn>.<Tb>, #0 -> UXTL{2} <Vd>.<Ta>, <Vn>.<Tb>
	 */
	type OpInfo struct {
		op   Operation
		ovar uint32
	}

	var operation = [2][32]OpInfo{
		{
			{ARM64_SSHR, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SSRA, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SRSHR, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SRSRA, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SHL, 3},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SQSHL, 3},
			{ARM64_UNDEFINED, 0},
			{ARM64_SHRN, 1},
			{ARM64_RSHRN, 1},
			{ARM64_SQSHRN, 1},
			{ARM64_SQRSHRN, 1},
			{ARM64_SSHLL, 4},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SCVTF, 2},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_FCVTZS, 2},
		}, {
			{ARM64_USHR, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_USRA, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_URSHR, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_URSRA, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SRI, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_SLI, 3},
			{ARM64_UNDEFINED, 0},
			{ARM64_SQSHLU, 3},
			{ARM64_UNDEFINED, 0},
			{ARM64_UQSHL, 3},
			{ARM64_UNDEFINED, 0},
			{ARM64_SQSHRUN, 1},
			{ARM64_SQRSHRUN, 1},
			{ARM64_UQSHRN, 1},
			{ARM64_UQRSHRN, 1},
			{ARM64_USHLL, 4},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UCVTF, 2},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_FCVTZU, 2},
		},
	}
	decode := SimdShiftByImm(i.raw)
	opinfo := operation[decode.U()][decode.Opcode()]
	//printf("opcod: %d U: %d '%s'\n", decode.Opcode(), decode.U, OperationString[opinfo.op])
	var sizemap = [2]uint8{64, 128}
	var size uint32
	for ; size < 4; size++ {
		if (decode.Immh() >> size) == 1 {
			break
		}
	}
	i.operation = opinfo.op
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = IMM32
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Immediate = uint64((16 << size) - (decode.Immh()<<3 | decode.Immb()))
	switch opinfo.ovar {
	case 0:
		i.operands[0].ElementSize = 1 << size
		i.operands[1].ElementSize = 1 << size
		i.operands[0].DataSize = uint32(sizemap[decode.Q()] / (8 << size))
		i.operands[1].DataSize = uint32(sizemap[decode.Q()] / (8 << size))
		//i.operands[2].Immediate =(decode.Immh() << 3 | decode.Immb()) - (8 << size)
		break
	case 1:
		i.operation = i.operation + Operation(decode.Q())
		i.operands[0].ElementSize = 1 << size
		i.operands[1].ElementSize = 2 << size
		i.operands[0].DataSize = uint32(sizemap[decode.Q()] / (8 << size))
		i.operands[1].DataSize = 64 / (8 << size)
		break
	case 2:
		if decode.Immh()>>2 == 1 {
			size = 0
		} else if decode.Immh()>>3 == 1 {
			size = 1
		} else {
			return nil, failedToDecodeInstruction
		}
		i.operands[0].ElementSize = 4 << size
		i.operands[1].ElementSize = 4 << size
		i.operands[0].DataSize = uint32(sizemap[decode.Q()] / (32 << size))
		i.operands[1].DataSize = uint32(sizemap[decode.Q()] / (32 << size))
		break
	case 3:
		i.operands[0].ElementSize = 1 << size
		i.operands[1].ElementSize = 1 << size
		i.operands[0].DataSize = uint32(sizemap[decode.Q()] / (8 << size))
		i.operands[1].DataSize = uint32(sizemap[decode.Q()] / (8 << size))
		i.operands[2].Immediate = uint64((decode.Immh()<<3 | decode.Immb()) - (8 << size))
		break
	case 4:
		i.operation = i.operation + Operation(decode.Q())
		i.operands[0].ElementSize = 2 << size
		i.operands[1].ElementSize = 1 << size
		i.operands[0].DataSize = 64 / (8 << size)
		i.operands[1].DataSize = uint32(sizemap[decode.Q()] / (8 << size))
		i.operands[2].Immediate = uint64((decode.Immh()<<3 | decode.Immb()) - (8 << size))
		break

	}
	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_simd_table_lookup() (*Instruction, error) {
	/* C4.6.14 Advanced SIMD table lookup
	 *
	 * TBL <Vd>.<Ta>, { <Vn>.16B }, <Vm>.<Ta>
	 * TBX <Vd>.<Ta>, { <Vn>.16B }, <Vm>.<Ta>
	 * TBL <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B }, <Vm>.<Ta>
	 * TBX <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B }, <Vm>.<Ta>
	 * TBL <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B }, <Vm>.<Ta>
	 * TBL <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B, <Vn+2>.16B }, <Vm>.<Ta>
	 * TBX <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B, <Vn+2>.16B }, <Vm>.<Ta>
	 * TBL <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B, <Vn+2>.16B, <Vn+3>.16B }, <Vm>.<Ta>
	 * TBX <Vd>.<Ta>, { <Vn>.16B, <Vn+1>.16B, <Vn+2>.16B, <Vn+3>.16B }, <Vm>.<Ta>
	 */
	decode := SimdTableLookup(i.raw)
	var operation = []Operation{
		ARM64_TBL,
		ARM64_TBX,
	}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = REG
	i.operands[1].OpClass = MULTI_REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	for idx := uint32(0); idx < decode.Len()+1; idx++ {
		i.operands[1].Reg[idx] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()+idx)%32)
	}
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	i.operands[0].ElementSize = 1
	i.operands[1].ElementSize = 1
	i.operands[2].ElementSize = 1
	i.operands[0].DataSize = 8 << decode.Q()
	i.operands[1].DataSize = 16
	i.operands[2].DataSize = 8 << decode.Q()

	return i, nil
}

func (i *Instruction) decompose_simd_vector_indexed_element() (*Instruction, error) {
	/* C4.6.18 Advanced SIMD vector x indexed element
	 *
	 * SMLAL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * SQDMLAL{2} <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * SMLSL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * SQDMLSL{2} <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * MUL        <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * SMULL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * SQDMULL{2} <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * SQDMULH    <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * SQRDMULH   <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * FMLA       <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * FMLS       <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * FMUL       <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * MLA        <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * UMLAL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * MLS        <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 * UMLSL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * UMULL{2}   <Vd>.<Ta>, <Vn>.<Tb>, <Vm>.<Ts>[<index>]
	 * FMULX      <Vd>.<T>,  <Vn>.<T>,  <Vm>.<Ts>[<index>]
	 */
	type OpInfo struct {
		op   Operation
		ovar uint32
	}
	var opinfo OpInfo
	var operation = [2][16]OpInfo{
		{
			{ARM64_UNDEFINED, 0},
			{ARM64_FMLA, 3},
			{ARM64_SMLAL, 1},
			{ARM64_SQDMLAL, 1},
			{ARM64_UNDEFINED, 0},
			{ARM64_FMLS, 3},
			{ARM64_SMLSL, 1},
			{ARM64_SQDMLSL, 1},
			{ARM64_MUL, 0},
			{ARM64_FMUL, 3},
			{ARM64_SMULL, 1},
			{ARM64_SQDMULL, 1},
			{ARM64_SQDMULH, 0},
			{ARM64_SQRDMULH, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
		}, {
			{ARM64_MLA, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UMLAL, 1},
			{ARM64_UNDEFINED, 0},
			{ARM64_MLS, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UMLSL, 1},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_FMULX, 3},
			{ARM64_UMULL, 1},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
			{ARM64_UNDEFINED, 0},
		},
	}
	decode := SimdVectorXIndexedElement(i.raw)
	opinfo = operation[decode.U()][decode.Opcode()]
	i.operation = opinfo.op

	i.operands[0].OpClass = REG
	i.operands[1].OpClass = REG
	i.operands[2].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
	i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
	i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rm()))
	index := uint32(decode.H()<<2 | decode.L()<<1 | decode.M())
	if opinfo.ovar == 1 {
		if decode.Size() == 0 || decode.Size() == 3 {
			return nil, failedToDecodeInstruction
		}
		var reghi uint32
		if decode.Size() == 2 {
			reghi = decode.M() << 4
		}
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(reghi|decode.Rm()))
		//'2' variant is always the next enumeration value
		i.operation = i.operation + Operation(decode.Q())
		i.operands[0].ElementSize = 2 << decode.Size()
		i.operands[1].ElementSize = 1 << decode.Size()
		i.operands[2].ElementSize = 1 << decode.Size()
		i.operands[0].DataSize = 8 >> decode.Size()
		i.operands[1].DataSize = 8 >> (decode.Size() - decode.Q())
		i.operands[2].DataSize = 0
		i.operands[2].Scale = 0x80000000 | (index >> (decode.Size() - 1))
	} else if opinfo.ovar == 2 {
		if ((decode.Size()&1) == 1 && decode.Q() == 0) || (decode.Size() == 1 && decode.L() == 1) {
			return nil, failedToDecodeInstruction
		}
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.M()<<4|decode.Rm()))
		i.operands[0].ElementSize = 2 << (decode.Size() & 1)
		i.operands[1].ElementSize = 2 << (decode.Size() & 1)
		i.operands[2].ElementSize = 4 << (decode.Size() & 1)
		i.operands[0].DataSize = 2 << ((decode.Size() & 1) - decode.Q())
		i.operands[1].DataSize = 2 << ((decode.Size() & 1) - decode.Q())
		i.operands[2].DataSize = 0
		i.operands[2].Scale = 0x80000000 | (index >> ((decode.Size() & 1) + 1))
	} else if opinfo.ovar == 3 {
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.M()<<4|decode.Rm()))
		i.operands[0].ElementSize = 1 << decode.Size()
		i.operands[1].ElementSize = 1 << decode.Size()
		i.operands[2].ElementSize = 1 << decode.Size()
		i.operands[0].DataSize = 8 >> (decode.Size() - decode.Q())
		i.operands[1].DataSize = 8 >> (decode.Size() - decode.Q())
		i.operands[2].DataSize = 0
		i.operands[2].Scale = 0x80000000 | (index >> (decode.Size() - 1))
	} else {
		var reghi uint32
		if decode.Size() == 2 {
			reghi = decode.M() << 4
		}
		i.operands[0].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rd()))
		i.operands[1].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(decode.Rn()))
		i.operands[2].Reg[0] = reg(REGSET_ZR, REG_V_BASE, int(reghi|decode.Rm()))
		i.operands[0].ElementSize = 1 << decode.Size()
		i.operands[1].ElementSize = 1 << decode.Size()
		i.operands[2].ElementSize = 1 << decode.Size()
		i.operands[0].DataSize = 8 >> (decode.Size() - decode.Q())
		i.operands[1].DataSize = 8 >> (decode.Size() - decode.Q())
		i.operands[2].DataSize = 0
		i.operands[2].Scale = 0x80000000 | (index >> (decode.Size() - 1))
	}

	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_system_arch_hints(decode System) (*Instruction, error) {
	// fmt.Println(decode)
	switch decode.Crn() {
	case 2: //Architectural hints
		switch (decode.Crm() << 3) | decode.Op2() {
		default:
			i.operation = ARM64_HINT
			i.operands[0].OpClass = IMM32
			i.operands[0].Immediate = uint64((decode.Crm() << 3) | decode.Op2())
			break
		case 0:
			i.operation = ARM64_NOP
			break
		case 1:
			i.operation = ARM64_YIELD
			break
		case 2:
			i.operation = ARM64_WFE
			break
		case 3:
			i.operation = ARM64_WFI
			break
		case 4:
			i.operation = ARM64_SEV
			break
		case 5:
			i.operation = ARM64_SEVL
			break

		// Added in ARMv8.6
		case 6:
			i.operation = ARM64_DGH
			break

		// Added with 8.2
		case 16:
			i.operation = ARM64_ESB
			break
		case 17:
			i.operation = ARM64_PSBCSYNC
			break

		// Added for 8.3
		case 7:
			i.operation = ARM64_XPACLRI
			break

		case 8:
			i.operation = ARM64_PACIA1716
			break
		case 10:
			i.operation = ARM64_PACIB1716
			break
		case 12:
			i.operation = ARM64_AUTIA1716
			break
		case 14:
			i.operation = ARM64_AUTIB1716
			break

		case 24:
			i.operation = ARM64_PACIAZ
			break
		case 25:
			i.operation = ARM64_PACIASP
			break
		case 26:
			i.operation = ARM64_PACIBZ
			break
		case 27:
			i.operation = ARM64_PACIBSP
			break
		case 28:
			i.operation = ARM64_AUTIAZ
			break
		case 29:
			i.operation = ARM64_AUTIASP
			break
		case 30:
			i.operation = ARM64_AUTIBZ
			break
		case 31:
			i.operation = ARM64_AUTIBSP
			break
		case 32:
			i.operation = ARM64_BTI
			break
		case 34:
			i.operation = ARM64_BTI
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_TGT_C)
			break
		case 36:
			i.operation = ARM64_BTI
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_TGT_J)
			break
		case 38:
			i.operation = ARM64_BTI
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_TGT_JC)
			break
		}
		break
	case 3: //Barriers and CLREX
		switch decode.Op2() {
		case 2:
			i.operation = ARM64_CLREX
			if decode.Crm() != 0xf {
				i.operands[0].OpClass = IMM32
				i.operands[0].Immediate = uint64(decode.Crm())
			}
			break
		case 4:
			i.operation = ARM64_DSB
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_NUMBER0) + decode.Crm()
			if decode.Crm() == 0 {
				i.operation = ARM64_SSBB
				i.operands[0].OpClass = NONE
			} else if decode.Crm() == 0b100 {
				i.operation = ARM64_PSSBB
				i.operands[0].OpClass = NONE
			}
			break
		case 5:
			i.operation = ARM64_DMB
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_NUMBER0) + decode.Crm()
			break
		case 6:
			i.operation = ARM64_ISB
			if decode.Crm() != 15 {
				i.operands[0].OpClass = IMM32
				i.operands[0].Immediate = uint64(decode.Crm())
			}
			break
		default:
			return nil, failedToDecodeInstruction
		}
		break
	case 4: //PSTATE Access
		switch decode.Op2() {
		case 4:
			if decode.Op1() == 0 { // MSR PAN, <Xt>
				i.operands[0].Reg[0] = uint32(REG_PAN)
				break
			}
			if decode.Op1() != 3 {
				return nil, failedToDecodeInstruction
			}
			i.operands[0].Reg[0] = uint32(REG_TCO)
			break
		case 5:
			if decode.Op1() != 0 {
				return nil, failedToDecodeInstruction
			}
			i.operands[0].Reg[0] = uint32(REG_SPSEL)
			break
		case 6:
			if decode.Op1() != 3 {
				return nil, failedToDecodeInstruction
			}
			i.operands[0].Reg[0] = uint32(REG_DAIFSET)
			break
		case 7:
			i.operands[0].Reg[0] = uint32(REG_DAIFCLR)
			break
		default:
			return nil, failedToDecodeInstruction
		}
		i.operation = ARM64_MSR
		i.operands[0].OpClass = SYS_REG
		i.operands[1].OpClass = IMM32
		i.operands[1].Immediate = uint64(decode.Crm())
		break
	default:
		{
			var operation = [2]Operation{ARM64_SYS, ARM64_SYSL}
			var operandSet = [2][5]uint32{{0, 1, 2, 3, 4}, {1, 2, 3, 4, 0}}
			i.operation = operation[decode.L()]
			i.operands[operandSet[decode.L()][0]].OpClass = IMM32
			i.operands[operandSet[decode.L()][0]].Immediate = uint64(decode.Op1())
			i.operands[operandSet[decode.L()][1]].OpClass = SYS_REG
			i.operands[operandSet[decode.L()][1]].Reg[0] = uint32(REG_C0) + decode.Crn()
			i.operands[operandSet[decode.L()][2]].OpClass = SYS_REG
			i.operands[operandSet[decode.L()][2]].Reg[0] = uint32(REG_C0) + decode.Crm()
			i.operands[operandSet[decode.L()][3]].OpClass = IMM32
			i.operands[operandSet[decode.L()][3]].Immediate = uint64(decode.Op2())
			if decode.Rt() != 31 {
				i.operands[operandSet[decode.L()][4]].OpClass = REG
				i.operands[operandSet[decode.L()][4]].Reg[0] =
					reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			}
		}
	}
	return i, nil
}

func (i *Instruction) decompose_system_cache_maintenance(decode System) (*Instruction, error) {
	// fmt.Println(decode)
	i.operands[1].OpClass = REG
	switch decode.Crn() {
	case 7:
		switch decode.Crm() {
		case 1: // Instruction cache maintenance instructions
			i.operation = ARM64_IC
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_IALLUIS)
			i.operands[1].OpClass = NONE
			break
		case 5: // Instruction cache maintenance instructions
			i.operation = ARM64_IC
			i.operands[0].OpClass = SYS_REG
			if decode.Op1() == 3 {
				i.operands[0].Reg[0] = uint32(REG_IVAU)
				i.operands[1].OpClass = REG
				i.operands[1].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			} else {
				i.operands[0].Reg[0] = uint32(REG_IALLU)
				i.operands[1].OpClass = NONE
			}
			break
		case 4: // Data cache zero operation
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			switch decode.Op2() {
			case 1:
				i.operands[0].Reg[0] = uint32(REG_ZVA)
			case 3:
				i.operands[0].Reg[0] = uint32(REG_GVA)
			case 4:
				i.operands[0].Reg[0] = uint32(REG_GZVA)
			}
			break
		case 6:
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			switch decode.Op2() {
			case 1:
				i.operands[0].Reg[0] = uint32(REG_IVAC)
			case 2:
				i.operands[0].Reg[0] = uint32(REG_ISW)
			case 3:
				i.operands[0].Reg[0] = uint32(REG_IGVAC)
			case 4:
				i.operands[0].Reg[0] = uint32(REG_IGSW)
			case 5:
				i.operands[0].Reg[0] = uint32(REG_IGDVAC)
			case 6:
				i.operands[0].Reg[0] = uint32(REG_IGDSW)
			}
			break
		case 10:
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_CSW)
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(1, REG_X_BASE, int(decode.Rt()))
			if decode.Op1() == 0 {
				switch decode.Op2() {
				case 2:
					i.operands[0].Reg[0] = uint32(REG_CSW)
				case 4:
					i.operands[0].Reg[0] = uint32(REG_CGSW)
				case 6:
					i.operands[0].Reg[0] = uint32(REG_CGDSW)
				default:
					return nil, failedToDecodeInstruction
				}
			} else if decode.Op1() == 3 {
				switch decode.Op2() {
				case 1:
					i.operands[0].Reg[0] = uint32(REG_CVAC)
				case 3:
					i.operands[0].Reg[0] = uint32(REG_CGVAC)
				case 5:
					i.operands[0].Reg[0] = uint32(REG_CGDVAC)
				default:
					return nil, failedToDecodeInstruction
				}
			} else {
				return nil, failedToDecodeInstruction
			}
			break
		case 11:
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(REG_CVAU)
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(1, REG_X_BASE, int(decode.Rt()))
			if decode.Op1() != 3 || decode.Op2() != 1 {
				return nil, failedToDecodeInstruction
			}
			break
		case 12: // 1100
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(1, REG_X_BASE, int(decode.Rt()))
			switch decode.Op2() {
			case 1:
				i.operands[0].Reg[0] = uint32(REG_CVAP)
			case 3:
				i.operands[0].Reg[0] = uint32(REG_CGVAP)
			case 5:
				i.operands[0].Reg[0] = uint32(REG_CGDVAP)
			default:
				return nil, failedToDecodeInstruction
			}
			break
		case 13: // 1101
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(1, REG_X_BASE, int(decode.Rt()))
			switch decode.Op2() {
			case 1:
				i.operands[0].Reg[0] = uint32(REG_CVADP)
			case 3:
				i.operands[0].Reg[0] = uint32(REG_CGVADP)
			case 5:
				i.operands[0].Reg[0] = uint32(REG_CGDVADP)
			default:
				return nil, failedToDecodeInstruction
			}
			break
		case 14: // Data cache maintenance instructions
			i.operation = ARM64_DC
			i.operands[0].OpClass = SYS_REG
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			if decode.Op1() == 0 {
				switch decode.Op2() {
				case 2:
					i.operands[0].Reg[0] = uint32(REG_CISW)
				case 4:
					i.operands[0].Reg[0] = uint32(REG_CIGSW)
				case 6:
					i.operands[0].Reg[0] = uint32(REG_CIGDSW)
				}
			} else if decode.Op1() == 3 {
				switch decode.Op2() {
				case 1:
					i.operands[0].Reg[0] = uint32(REG_CIVAC)
				case 3:
					i.operands[0].Reg[0] = uint32(REG_CIGVAC)
				case 5:
					i.operands[0].Reg[0] = uint32(REG_CIGDVAC)
				}
			} else if decode.Op1() != 3 || decode.Op2() != 1 {
				return nil, failedToDecodeInstruction
			}
			break
		case 8: // Address translation instructions
			i.operation = ARM64_AT
			i.operands[0].OpClass = SYS_REG
			i.operands[1].OpClass = REG
			i.operands[1].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			if decode.Op1() == 0 && decode.Op2() < 4 {
				var sysregs = []SystemReg{REG_S1E1R, REG_S1E1W, REG_S1E0R, REG_S1E0W}
				i.operands[0].Reg[0] = uint32(sysregs[decode.Op2()])
			} else if decode.Op1() == 4 {
				var sysregs = []SystemReg{
					REG_S1E2R, REG_S1E2W, SYSREG_NONE, SYSREG_NONE,
					REG_S12E1R, REG_S12E1W, REG_S12E0R, REG_S12E0W,
				}
				i.operands[0].Reg[0] = uint32(sysregs[decode.Op2()])
			} else if decode.Op1() == 6 && decode.Op2() < 2 {
				var sysregs = []SystemReg{REG_S1E3R, REG_S1E3W}
				i.operands[0].Reg[0] = uint32(sysregs[decode.Op2()])
			}
			break
		}
		break
	case 8: // TLB maintenance instruction
		{
			i.operation = ARM64_TLBI
			sysreg := SYSREG_NONE
			switch decode.Op1() {
			case 0:
				switch decode.Crm() {
				case 3:
					{
						var sysregs = []SystemReg{
							REG_VMALLE1IS, REG_VAE1IS, REG_ASIDE1IS, REG_VAAE1IS,
							SYSREG_NONE, REG_VALE1IS, SYSREG_NONE, REG_VAALE1IS,
						}
						sysreg = sysregs[decode.Op2()]
						if decode.Op2() == 0 {
							i.operands[1].OpClass = NONE
						} else {
							i.operands[1].OpClass = REG
						}
						break
					}
				case 7:
					{
						var sysregs = []SystemReg{
							REG_VMALLE1, REG_VAE1, REG_ASIDE1, REG_VAAE1,
							SYSREG_NONE, REG_VALE1, SYSREG_NONE, REG_VAALE1,
						}
						sysreg = sysregs[decode.Op2()]
						if decode.Op2() == 0 {
							i.operands[1].OpClass = NONE
						} else {
							i.operands[1].OpClass = REG
						}
						break
					}
				}
				break
			case 4:
				switch decode.Crm() {
				case 0:
					{
						var sysregs = []SystemReg{
							SYSREG_NONE, REG_IPAS2E1IS, SYSREG_NONE,
							SYSREG_NONE, SYSREG_NONE, REG_IPAS2LE1IS}
						sysreg = sysregs[decode.Op2()]
						break
					}
				case 3:
					{
						var sysregs = []SystemReg{
							REG_ALLE2IS, REG_VAE2IS, SYSREG_NONE, SYSREG_NONE,
							REG_ALLE1IS, REG_VALE2IS, REG_VMALLS12E1IS, SYSREG_NONE}
						sysreg = sysregs[decode.Op2()]
						if decode.Op2() == 0 || decode.Op2() == 4 || decode.Op2() == 6 {
							i.operands[1].OpClass = NONE
						}
						break
					}
				case 4:
					if decode.Op2() == 1 {
						sysreg = REG_IPAS2E1
					} else if decode.Op2() == 5 {
						sysreg = REG_IPAS2LE1
					}
					break
				case 7:
					{
						var sysregs = []SystemReg{
							REG_ALLE2, REG_VAE2, SYSREG_NONE, SYSREG_NONE,
							REG_ALLE1, REG_VALE2, REG_VMALLS12E1, SYSREG_NONE}
						sysreg = sysregs[decode.Op2()]
						if decode.Op2() == 0 || decode.Op2() == 4 || decode.Op2() == 6 {
							i.operands[1].OpClass = NONE
						}
						break
					}
				}
				break
			case 6:
				switch decode.Crm() {
				case 3:
					{
						var sysregs = []SystemReg{
							REG_ALLE3IS, REG_VAE3IS, SYSREG_NONE, SYSREG_NONE,
							SYSREG_NONE, REG_VALE3IS, SYSREG_NONE, SYSREG_NONE}
						sysreg = sysregs[decode.Op2()]
						if decode.Op2() == 0 {
							i.operands[1].OpClass = NONE
						} else {
							i.operands[1].OpClass = REG
						}
						break
					}
				case 7:
					{
						var sysregs = []SystemReg{
							REG_ALLE3, REG_VAE3, SYSREG_NONE, SYSREG_NONE,
							SYSREG_NONE, REG_VALE3, SYSREG_NONE, SYSREG_NONE}
						sysreg = sysregs[decode.Op2()]
						if decode.Op2() == 0 {
							i.operands[1].OpClass = NONE
						} else {
							i.operands[1].OpClass = REG
						}
						break
					}
				}
				break
			}
			i.operands[0].OpClass = SYS_REG
			i.operands[0].Reg[0] = uint32(sysreg)
			i.operands[1].Reg[0] = reg(1, REG_X_BASE, int(decode.Rt()))
		}
		break
	case 11:
		fallthrough
	case 12:
		fallthrough
	case 13:
		fallthrough
	case 14:
		fallthrough
	case 15:
		fallthrough
	default:
		{
			var operation = [2]Operation{ARM64_SYS, ARM64_SYSL}
			var operandSet = [2][5]uint32{{0, 1, 2, 3, 4}, {1, 2, 3, 4, 0}}
			i.operation = operation[decode.L()]
			i.operands[operandSet[decode.L()][0]].OpClass = IMM32
			i.operands[operandSet[decode.L()][0]].Immediate = uint64(decode.Op1())
			i.operands[operandSet[decode.L()][1]].OpClass = SYS_REG
			i.operands[operandSet[decode.L()][1]].Reg[0] = uint32(REG_C0) + decode.Crn()
			i.operands[operandSet[decode.L()][2]].OpClass = SYS_REG
			i.operands[operandSet[decode.L()][2]].Reg[0] = uint32(REG_C0) + decode.Crm()
			i.operands[operandSet[decode.L()][3]].OpClass = IMM32
			i.operands[operandSet[decode.L()][3]].Immediate = uint64(decode.Op2())
			if decode.Rt() != 31 {
				i.operands[operandSet[decode.L()][4]].OpClass = REG
				i.operands[operandSet[decode.L()][4]].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			}
		}
	}
	if i.operation == ARM64_UNDEFINED {
		return nil, failedToDisassembleOperation
	}

	return i, nil
}

func (i *Instruction) decompose_system_debug_and_trace_regs(decode System) (*Instruction, error) {
	// fmt.Println(decode)
	sysreg := SYSREG_NONE
	var operation = [2]Operation{ARM64_MSR, ARM64_MRS}
	var dbgreg = [4][16]SystemReg{
		{
			REG_DBGBVR0_EL1, REG_DBGBVR1_EL1, REG_DBGBVR2_EL1, REG_DBGBVR3_EL1,
			REG_DBGBVR4_EL1, REG_DBGBVR5_EL1, REG_DBGBVR6_EL1, REG_DBGBVR7_EL1,
			REG_DBGBVR8_EL1, REG_DBGBVR9_EL1, REG_DBGBVR10_EL1, REG_DBGBVR11_EL1,
			REG_DBGBVR12_EL1, REG_DBGBVR13_EL1, REG_DBGBVR14_EL1, REG_DBGBVR15_EL1,
		}, {
			REG_DBGBCR0_EL1, REG_DBGBCR1_EL1, REG_DBGBCR2_EL1, REG_DBGBCR3_EL1,
			REG_DBGBCR4_EL1, REG_DBGBCR5_EL1, REG_DBGBCR6_EL1, REG_DBGBCR7_EL1,
			REG_DBGBCR8_EL1, REG_DBGBCR9_EL1, REG_DBGBCR10_EL1, REG_DBGBCR11_EL1,
			REG_DBGBCR12_EL1, REG_DBGBCR13_EL1, REG_DBGBCR14_EL1, REG_DBGBCR15_EL1,
		}, {
			REG_DBGWVR0_EL1, REG_DBGWVR1_EL1, REG_DBGWVR2_EL1, REG_DBGWVR3_EL1,
			REG_DBGWVR4_EL1, REG_DBGWVR5_EL1, REG_DBGWVR6_EL1, REG_DBGWVR7_EL1,
			REG_DBGWVR8_EL1, REG_DBGWVR9_EL1, REG_DBGWVR10_EL1, REG_DBGWVR11_EL1,
			REG_DBGWVR12_EL1, REG_DBGWVR13_EL1, REG_DBGWVR14_EL1, REG_DBGWVR15_EL1,
		}, {
			REG_DBGWCR0_EL1, REG_DBGWCR1_EL1, REG_DBGWCR2_EL1, REG_DBGWCR3_EL1,
			REG_DBGWCR4_EL1, REG_DBGWCR5_EL1, REG_DBGWCR6_EL1, REG_DBGWCR7_EL1,
			REG_DBGWCR8_EL1, REG_DBGWCR9_EL1, REG_DBGWCR10_EL1, REG_DBGWCR11_EL1,
			REG_DBGWCR12_EL1, REG_DBGWCR13_EL1, REG_DBGWCR14_EL1, REG_DBGWCR15_EL1,
		},
	}
	switch decode.Op1() { //Table C5-5 System instruction encodings for debug System register access
	case 0:
		switch decode.Crn() {
		case 0:
			if decode.Crm() == 0 && decode.Op2() == 2 {
				sysreg = REG_OSDTRRX_EL1
			} else if decode.Crm() == 2 && decode.Op2() == 0 {
				sysreg = REG_MDCCINT_EL1
			} else if decode.Crm() == 2 && decode.Op2() == 2 {
				sysreg = REG_MDSCR_EL1
			} else if decode.Crm() == 3 && decode.Op2() == 2 {
				sysreg = REG_OSDTRTX_EL1
			} else if decode.Crm() == 6 && decode.Op2() == 2 {
				sysreg = REG_OSECCR_EL1
			} else {
				if decode.Op2() > 3 && decode.Op2() < 8 {
					sysreg = dbgreg[decode.Op2()-4][decode.Crm()]
				}
			}
			break
		case 1:
			switch decode.Crm() {
			case 0:
				if decode.Op2() == 0 {
					sysreg = REG_MDRAR_EL1
				} else if decode.Op2() == 4 {
					sysreg = REG_OSLAR_EL1
				}
				break
			case 1:
				if decode.Op2() == 4 {
					sysreg = REG_OSLSR_EL1
				}
				break
			case 3:
				if decode.Op2() == 4 {
					sysreg = REG_OSDLR_EL1
				}
				break
			case 4:
				if decode.Op2() == 4 {
					sysreg = REG_DBGPRCR_EL1
				}
				break
			}
			break
		case 7:
			if decode.Op2() != 6 {
				break
			}
			switch decode.Crm() {
			case 8:
				sysreg = REG_DBGCLAIMSET_EL1
				break
			case 9:
				sysreg = REG_DBGCLAIMCLR_EL1
				break
			case 14:
				sysreg = REG_DBGAUTHSTATUS_EL1
				break
			}
		}
		break
	case 1:
		{
			//Switch operands depending on load vs store
			op1 := ^(^decode.L()) & 1
			op2 := ^decode.L() & 1
			i.operation = operation[op1]
			i.operands[op1].OpClass = IMPLEMENTATION_SPECIFIC
			i.operands[op1].Reg[0] = decode.Op0()
			i.operands[op1].Reg[1] = decode.Op1()
			i.operands[op1].Reg[2] = decode.Crn()
			i.operands[op1].Reg[3] = decode.Crm()
			i.operands[op1].Reg[4] = decode.Op2()
			i.operands[op2].OpClass = REG
			i.operands[op2].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			return i, nil
		}
	case 2:
		if decode.Crn() == 0 && decode.Crm() == 0 {
			sysreg = REG_TEECR32_EL1
		} else if decode.Crn() == 1 && decode.Crm() == 0 {
			sysreg = REG_TEEHBR32_EL1
		}
		break
	case 3:
		if decode.Crn() != 0 || decode.Op2() != 0 {
			break
		}
		switch decode.Crm() {
		case 1:
			sysreg = REG_MDCCSR_EL0
			break
		case 4:
			sysreg = REG_DBGDTR_EL0
			break
		case 5:
			if decode.L() != 0 {
				sysreg = REG_DBGDTRRX_EL0
			} else {
				sysreg = REG_DBGDTRTX_EL0
			}
		}
		break
	case 4:
		if decode.Crn() == 0 && decode.Crm() == 7 && decode.Op2() == 0 {
			sysreg = REG_DBGVCR32_EL2
		}
		break
		//default:
		//	printf("%s\n", __FUNCTION__)
	}
	op1 := ^(^decode.L()) & 1
	op2 := ^decode.L() & 1
	i.operation = operation[op1]
	i.operands[op1].OpClass = SYS_REG
	i.operands[op1].Reg[0] = uint32(sysreg)
	i.operands[op2].OpClass = REG
	i.operands[op2].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))

	if sysreg == SYSREG_NONE {
		return nil, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_system_debug_and_trace_regs2(decode System) (*Instruction, error) {
	sysreg := SYSREG_NONE
	var operation = [2]Operation{ARM64_MSR, ARM64_MRS}
	fmt.Println(decode)
	switch decode.Crn() {
	case 0:
		if decode.Op1() == 0 {
			var sysregs = [8][8]SystemReg{
				{REG_MIDR_EL1, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, REG_MPIDR_EL1, REG_REVIDR_EL1, SYSREG_NONE},
				{REG_ID_PFR0_EL1, REG_ID_PFR1_EL1, REG_ID_DFR0_EL1, REG_ID_AFR0_EL1,
					REG_ID_MMFR0_EL1, REG_ID_MMFR1_EL1, REG_ID_MMFR2_EL1, REG_ID_MMFR3_EL1},
				{REG_ID_ISAR0_EL1, REG_ID_ISAR1_EL1, REG_ID_ISAR2_EL1, REG_ID_ISAR3_EL1,
					REG_ID_ISAR4_EL1, REG_ID_ISAR5_EL1, REG_ID_MMFR4_EL1, REG_ID_ISAR6_EL1},
				{REG_MVFR0_EL1, REG_MVFR1_EL1, REG_MVFR2_EL1, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, REG_ID_MMFR5_EL1, SYSREG_NONE},
				{REG_ID_AA64PFR0_EL1, REG_ID_AA64PFR1_EL1, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{REG_ID_AA64DFR0_EL1, REG_ID_AA64DFR1_EL1, SYSREG_NONE, SYSREG_NONE, REG_ID_AA64AFR0_EL1, REG_ID_AA64AFR1_EL1, SYSREG_NONE, SYSREG_NONE},
				{REG_ID_AA64ISAR0_EL1, REG_ID_AA64ISAR1_EL1, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{REG_ID_AA64MMFR0_EL1, REG_ID_AA64MMFR1_EL1, REG_ID_AA64MMFR2_EL1, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
			}
			sysreg = sysregs[decode.Crm()][decode.Op2()]
		} else if decode.Crm() == 0 {
			var sysregs = [8][8]SystemReg{
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{REG_CCSIDR_EL1, REG_CLIDR_EL1, SYSREG_NONE, SYSREG_NONE, REG_GMID_EL1, SYSREG_NONE, SYSREG_NONE, REG_AIDR_EL1},
				{REG_CSSELR_EL1, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, REG_CTR_EL0, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, REG_DCZID_EL0},
				{REG_VPIDR_EL2, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, REG_VMPIDR_EL2, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
			}
			sysreg = sysregs[decode.Op1()][decode.Op2()]
		}
		break
	case 1:
		switch decode.Op1() {
		case 0:
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_SCTLR_EL1
					break
				case 1:
					sysreg = REG_ACTLR_EL1
					break
				case 2:
					sysreg = REG_CPACR_EL1
					break
				case 5:
					sysreg = REG_RGSR_EL1
					break
				case 6:
					sysreg = REG_GCR_EL1
					break
				}
			}
			break
		case 5:
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_SCTLR_EL12
					break
				case 2:
					sysreg = REG_CPACR_EL12
					break
				}
			}
			break
		case 4:
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_SCTLR_EL2
					break
				case 1:
					sysreg = REG_ACTLR_EL2
					break
				}
			} else if decode.Crm() == 1 {
				var sysregs = []SystemReg{
					REG_HCR_EL2, REG_MDCR_EL2, REG_CPTR_EL2, REG_HSTR_EL2,
					REG_HFGRTR_EL2, REG_HFGWTR_EL2, REG_HFGITR_EL2, REG_HACR_EL2,
				}
				sysreg = sysregs[decode.Op2()]
			} else if decode.Crm() == 6 {
				sysreg = REG_ICC_PMR_EL1
			}
			break
		case 6:
			switch decode.Crm() {
			case 0:
				if decode.Op2() == 0 {
					sysreg = REG_SCTLR_EL3
				} else if decode.Op2() == 1 {
					sysreg = REG_ACTLR_EL3
				}
				break
			case 1:
				switch decode.Op2() {
				case 0:
					sysreg = REG_SCR_EL3
					break
				case 1:
					sysreg = REG_SDER32_EL3
					break
				case 2:
					sysreg = REG_CPTR_EL3
					break
				}
				break
			case 3:
				if decode.Op2() == 1 {
					sysreg = REG_MDCR_EL3
				}
				break
			}
			break
		}
		break
	case 2:
		switch decode.Op1() {
		case 0:
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_TTBR0_EL1
					break
				case 1:
					sysreg = REG_TTBR1_EL1
					break
				case 2:
					sysreg = REG_TCR_EL1
					break
				}
			}
			if decode.Crm() == 1 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_APIAKEYLO_EL1
					break
				case 1:
					sysreg = REG_APIAKEYHI_EL1
					break
				case 2:
					sysreg = REG_APIBKEYLO_EL1
					break
				case 3:
					sysreg = REG_APIBKEYHI_EL1
					break
				}
			}
			if decode.Crm() == 2 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_APDAKEYLO_EL1
					break
				case 1:
					sysreg = REG_APDAKEYHI_EL1
					break
				case 2:
					sysreg = REG_APDBKEYLO_EL1
					break
				case 3:
					sysreg = REG_APDBKEYHI_EL1
					break
				}
			}
			if decode.Crm() == 3 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_APGAKEYLO_EL1
					break
				case 1:
					sysreg = REG_APGAKEYHI_EL1
					break
				}
			}
			break
		case 4:
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_TTBR0_EL2
					break
				case 2:
					sysreg = REG_TCR_EL2
					break
				}
			} else if decode.Crm() == 1 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_VTTBR_EL2
					break
				case 2:
					sysreg = REG_VTCR_EL2
					break
				}
			}
			break
		case 5:
			if decode.Crm() == 0 {
				if decode.Op2() == 0 {
					sysreg = REG_TTBR0_EL12
				} else if decode.Op2() == 1 {
					sysreg = REG_TTBR1_EL12
				} else if decode.Op2() == 2 {
					sysreg = REG_TCR_EL12
				}
			} else if decode.Crm() == 2 {
				if decode.Op2() == 1 {
					sysreg = REG_TTBR1_EL12
				} else if decode.Op2() == 2 {
					sysreg = REG_TCR_EL12
				}
			}
			break
		case 6:
			if decode.Crm() == 0 {
				if decode.Op2() == 0 {
					sysreg = REG_TTBR0_EL3
				} else if decode.Op2() == 2 {
					sysreg = REG_TCR_EL3
				}
			}
			break
		}
		break
	case 3:
		if decode.Crm() == 1 && decode.Op2() == 4 {
			sysreg = REG_HDFGRTR_EL2
			break
		} else if decode.Crm() == 1 && decode.Op2() == 5 {
			sysreg = REG_HDFGWTR_EL2
			break
		}
		if decode.Op1() != 4 || decode.Crm() != 0 || decode.Op2() != 0 {
			break
		}
		sysreg = REG_DACR32_EL2
		break
	case 4:
		switch decode.Op1() {
		case 0:
			switch decode.Crm() {
			case 0:
				if decode.Op2() == 1 {
					sysreg = REG_ELR_EL1
				} else if decode.Op2() == 0 {
					sysreg = REG_SPSR_EL1
				}
				break
			case 1:
				if decode.Op2() == 0 {
					sysreg = REG_SP_EL0
				}
				break
			case 2:
				if decode.Op2() == 0 {
					sysreg = REG_SPSEL
				} else if decode.Op2() == 2 {
					sysreg = REG_CURRENT_EL
				} else if decode.Op2() == 3 {
					sysreg = REG_PAN
				}
				break
			case 6:
				sysreg = REG_ICC_PMR_EL1
				break
			}
			break
		case 3:
			if decode.Op2() == 0 {
				switch decode.Crm() {
				case 2:
					sysreg = REG_NZCV
					break
				case 4:
					sysreg = REG_FPCR
					break
				case 5:
					sysreg = REG_DSPSR_EL0
					break
				}
			} else if decode.Op2() == 1 {
				switch decode.Crm() {
				case 2:
					sysreg = REG_DAIF
					break
				case 4:
					sysreg = REG_FPSR
					break
				case 5:
					sysreg = REG_DLR_EL0
					break
				}
			} else if decode.Op2() == 7 {
				switch decode.Crm() {
				case 2:
					sysreg = REG_TCO
					break
				}
			}
			break
		case 4:
			switch decode.Crm() {
			case 0:
				if decode.Op2() == 0 {
					sysreg = REG_SPSR_EL2
				} else if decode.Op2() == 1 {
					sysreg = REG_ELR_EL2
				}
				break
			case 1:
				if decode.Op2() == 0 {
					sysreg = REG_SP_EL1
				}
				break
			case 3:
				switch decode.Op2() {
				case 0:
					sysreg = REG_SPSR_IRQ
					break
				case 1:
					sysreg = REG_SPSR_ABT
					break
				case 2:
					sysreg = REG_SPSR_UND
					break
				case 3:
					sysreg = REG_SPSR_FIQ
					break
				}
				break
			}
			break
		case 5:
			if decode.Crm() == 0 {
				if decode.Op2() == 0 {
					sysreg = REG_SPSR_EL12
				} else if decode.Op2() == 1 {
					sysreg = REG_ELR_EL12
				}
			}
			break
		case 6:
			if decode.Crm() == 0 {
				if decode.Op2() == 0 {
					sysreg = REG_SPSR_EL3
				} else if decode.Op2() == 1 {
					sysreg = REG_ELR_EL3
				}
			} else if decode.Crm() == 1 {
				if decode.Op2() == 0 {
					sysreg = REG_SP_EL2
				}
			}
			break
		}
		break
	case 5:
		{
			if decode.Crm() == 6 {
				if decode.Op2() == 0 {
					switch decode.Op1() {
					case 0:
						sysreg = REG_TFSR_EL1
					case 4:
						sysreg = REG_TFSR_EL2
					case 5:
						sysreg = REG_TFSR_EL12
					case 6:
						sysreg = REG_TFSR_EL3
					}
					break
				} else if decode.Op2() == 1 {
					switch decode.Op1() {
					case 0:
						sysreg = REG_TFSRE0_EL1
					}
					break
				}
			}

			if decode.Crm() > 3 || decode.Op2() > 1 {
				break
			}

			var sysregs = [4][4][2]SystemReg{
				{
					{SYSREG_NONE, SYSREG_NONE},
					{REG_AFSR0_EL1, REG_AFSR1_EL1},
					{REG_ESR_EL1, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE},
				}, {
					{SYSREG_NONE, REG_IFSR32_EL2},
					{REG_AFSR0_EL2, REG_AFSR1_EL2},
					{REG_ESR_EL2, SYSREG_NONE},
					{REG_FPEXC32_EL2, SYSREG_NONE},
				}, {
					{SYSREG_NONE, SYSREG_NONE},
					{REG_AFSR0_EL3, REG_AFSR1_EL3},
					{REG_ESR_EL3, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE},
				}, {
					{SYSREG_NONE, SYSREG_NONE},
					{REG_AFSR0_EL12, REG_AFSR1_EL12},
					{REG_ESR_EL12, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE},
				},
			}
			switch decode.Op1() {
			case 0:
				sysreg = sysregs[0][decode.Crm()][decode.Op2()]
				break
			case 4:
				sysreg = sysregs[1][decode.Crm()][decode.Op2()]
				break
			case 5:
				sysreg = sysregs[3][decode.Crm()][decode.Op2()]
				break
			case 6:
				sysreg = sysregs[2][decode.Crm()][decode.Op2()]
				break
			}
			break
		}
	case 6:
		if decode.Op1() == 0 && decode.Crm() == 0 && decode.Op2() == 0 {
			sysreg = REG_FAR_EL1
		} else if decode.Op1() == 4 && decode.Crm() == 0 {
			if decode.Op2() == 0 {
				sysreg = REG_FAR_EL2
			} else if decode.Op2() == 4 {
				sysreg = REG_HPFAR_EL2
			}
		} else if decode.Op1() == 6 && decode.Crm() == 0 && decode.Op2() == 0 {
			sysreg = REG_FAR_EL3
		} else if decode.Op1() == 0 && decode.Crm() == 11 && decode.Op2() == 5 {
			sysreg = REG_ICC_SGI1R_EL1
		} else if decode.Op1() == 5 && decode.Crm() == 0 && decode.Op2() == 0 {
			sysreg = REG_FAR_EL12
		}
		break
	case 7:
		if decode.Op1() == 0 && decode.Crm() == 4 && decode.Op2() == 0 {
			sysreg = REG_PAR_EL1
		}
		break
	case 9:
		{
			if decode.Op1() == 0 && decode.Crm() == 14 {
				if decode.Op2() == 1 {
					sysreg = REG_PMINTENSET_EL1
				} else if decode.Op2() == 2 {
					sysreg = REG_PMINTENCLR_EL1
				}
			} else if decode.Op1() == 3 {
				var sysregs = [3][8]SystemReg{
					{
						REG_PMCR_EL0, REG_PMCNTENSET_EL0, REG_PMCNTENCLR_EL0, REG_PMOVSCLR_EL0,
						REG_PMSWINC_EL0, REG_PMSELR_EL0, REG_PMCEID0_EL0, REG_PMCEID1_EL0,
					},
					{REG_PMCCNTR_EL0, REG_PMXEVTYPER_EL0, REG_PMXEVCNTR_EL0, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_PMUSERENR_EL0, SYSREG_NONE, SYSREG_NONE, REG_PMOVSSET_EL0, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}
				if decode.Crm() > 11 && decode.Crm() < 15 {
					sysreg = sysregs[decode.Crm()-12][decode.Op2()]
				}
			}
			break
		}
	case 10:
		{
			if decode.Op1() == 0 && decode.Crm() == 0b100 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_LORSA_EL1
				case 1:
					sysreg = REG_LOREA_EL1
				case 2:
					sysreg = REG_LORN_EL1
				case 3:
					sysreg = REG_LORC_EL1
				case 7:
					sysreg = REG_LORID_EL1
				}
			} else if decode.Op2() == 0 {
				var sysregs = [9][3]SystemReg{
					{REG_MAIR_EL1, REG_AMAIR_EL1, REG_LORSA_EL1},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_MAIR_EL2, REG_AMAIR_EL2, SYSREG_NONE},
					{REG_MAIR_EL12, REG_AMAIR_EL12, SYSREG_NONE},
					{REG_MAIR_EL3, REG_AMAIR_EL3, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}
				sysreg = sysregs[decode.Op1()][decode.Crm()-2]
			}
			break
		}
	case 12:
		if decode.Op1() == 0 {
			switch decode.Crm() {
			case 0:
				switch decode.Op2() {
				case 0:
					sysreg = REG_VBAR_EL1
					break
				case 1:
					sysreg = REG_RVBAR_EL1
					break
				case 2:
					sysreg = REG_RMR_EL1
					break
				}
				break
			case 1:
				if decode.Op2() == 0 {
					sysreg = REG_ISR_EL1
				}
				break
			case 8:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICC_IAR0_EL1
					break
				case 1:
					sysreg = REG_ICC_EOIR0_EL1
					break
				case 2:
					sysreg = REG_ICC_HPPIR0_EL1
					break
				case 3:
					sysreg = REG_ICC_BPR0_EL1
					break
				case 4:
					sysreg = REG_ICC_AP0R0_EL1
					break
				case 5:
					sysreg = REG_ICC_AP0R1_EL1
					break
				case 6:
					sysreg = REG_ICC_AP0R2_EL1
					break
				case 7:
					sysreg = REG_ICC_AP0R3_EL1
					break
				}
				break
			case 9:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICC_AP1R0_EL1
					break
				case 1:
					sysreg = REG_ICC_AP1R1_EL1
					break
				case 2:
					sysreg = REG_ICC_AP1R2_EL1
					break
				case 3:
					sysreg = REG_ICC_AP1R3_EL1
					break
				}
				break
			case 11:
				if decode.Op2() == 1 {
					sysreg = REG_ICC_DIR_EL1
				} else if decode.Op2() == 3 {
					sysreg = REG_ICC_RPR_EL1
				} else if decode.Op2() == 5 {
					sysreg = REG_ICC_SGI1R_EL1
				} else if decode.Op2() == 6 {
					sysreg = REG_ICC_ASGI1R_EL1
				} else if decode.Op2() == 7 {
					sysreg = REG_ICC_SGI0R_EL1
				}
				break
			case 12:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICC_IAR1_EL1
					break
				case 1:
					sysreg = REG_ICC_EOIR1_EL1
					break
				case 2:
					sysreg = REG_ICC_HPPIR1_EL1
					break
				case 3:
					sysreg = REG_ICC_BPR1_EL1
					break
				case 4:
					sysreg = REG_ICC_CTLR_EL1
					break
				case 5:
					sysreg = REG_ICC_SRE_EL1
					break
				case 6:
					sysreg = REG_ICC_IGRPEN0_EL1
					break
				case 7:
					sysreg = REG_ICC_IGRPEN1_EL1
					break
				}
				break
			case 13:
				if decode.Op2() == 0 {
					sysreg = REG_ICC_SEIEN_EL1
				}
				break
			default:
				break
			}
		} else if decode.Op1() == 1 && decode.Crm() == 12 {
			sysreg = REG_ICC_ASGI1R_EL2
		} else if decode.Op1() == 2 && decode.Crm() == 12 {
			sysreg = REG_ICC_SGI0R_EL2
		} else if decode.Op1() == 4 {
			switch decode.Crm() {
			case 0:
				switch decode.Op2() {
				case 0:
					sysreg = REG_VBAR_EL2
					break
				case 1:
					sysreg = REG_RVBAR_EL2
					break
				case 2:
					sysreg = REG_RMR_EL2
					break
				}
				break
			case 8:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_AP0R0_EL2
					break
				case 1:
					sysreg = REG_ICH_AP0R1_EL2
					break
				case 2:
					sysreg = REG_ICH_AP0R2_EL2
					break
				case 3:
					sysreg = REG_ICH_AP0R3_EL2
					break
				}
				break
			case 9:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_AP1R0_EL2
					break
				case 1:
					sysreg = REG_ICH_AP1R1_EL2
					break
				case 2:
					sysreg = REG_ICH_AP1R2_EL2
					break
				case 3:
					sysreg = REG_ICH_AP1R3_EL2
					break
				case 4:
					sysreg = REG_ICH_AP1R4_EL2
					break
				case 5:
					sysreg = REG_ICC_HSRE_EL2
					break
				}
				break
			case 11:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_HCR_EL2
					break
				case 1:
					sysreg = REG_ICH_VTR_EL2
					break
				case 2:
					sysreg = REG_ICH_MISR_EL2
					break
				case 3:
					sysreg = REG_ICH_EISR_EL2
					break
				case 5:
					sysreg = REG_ICH_ELRSR_EL2
					break
				case 7:
					sysreg = REG_ICH_VMCR_EL2
					break
				}
				break
			case 12:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_LR0_EL2
					break
				case 1:
					sysreg = REG_ICH_LR1_EL2
					break
				case 2:
					sysreg = REG_ICH_LR2_EL2
					break
				case 3:
					sysreg = REG_ICH_LR3_EL2
					break
				case 4:
					sysreg = REG_ICH_LR4_EL2
					break
				case 5:
					sysreg = REG_ICH_LR5_EL2
					break
				case 6:
					sysreg = REG_ICH_LR6_EL2
					break
				case 7:
					sysreg = REG_ICH_LR7_EL2
					break
				}
				break
			case 13:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_LR8_EL2
					break
				case 1:
					sysreg = REG_ICH_LR9_EL2
					break
				case 2:
					sysreg = REG_ICH_LR10_EL2
					break
				case 3:
					sysreg = REG_ICH_LR11_EL2
					break
				case 4:
					sysreg = REG_ICH_LR12_EL2
					break
				case 5:
					sysreg = REG_ICH_LR13_EL2
					break
				case 6:
					sysreg = REG_ICH_LR14_EL2
					break
				case 7:
					sysreg = REG_ICH_LR15_EL2
					break
				}
				break
			case 14:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_LRC0_EL2
					break
				case 1:
					sysreg = REG_ICH_LRC1_EL2
					break
				case 2:
					sysreg = REG_ICH_LRC2_EL2
					break
				case 3:
					sysreg = REG_ICH_LRC3_EL2
					break
				case 4:
					sysreg = REG_ICH_LRC4_EL2
					break
				case 5:
					sysreg = REG_ICH_LRC5_EL2
					break
				case 6:
					sysreg = REG_ICH_LRC6_EL2
					break
				case 7:
					sysreg = REG_ICH_LRC7_EL2
					break
				}
				break
			case 15:
				switch decode.Op2() {
				case 0:
					sysreg = REG_ICH_LRC8_EL2
					break
				case 1:
					sysreg = REG_ICH_LRC9_EL2
					break
				case 2:
					sysreg = REG_ICH_LRC10_EL2
					break
				case 3:
					sysreg = REG_ICH_LRC11_EL2
					break
				case 4:
					sysreg = REG_ICH_LRC12_EL2
					break
				case 5:
					sysreg = REG_ICH_LRC13_EL2
					break
				case 6:
					sysreg = REG_ICH_LRC14_EL2
					break
				case 7:
					sysreg = REG_ICH_LRC15_EL2
					break
				}
				break
			}
		} else if decode.Op1() == 5 {
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_VBAR_EL12
					break
				}
			}
		} else if decode.Op1() == 6 {
			if decode.Crm() == 0 {
				switch decode.Op2() {
				case 0:
					sysreg = REG_VBAR_EL3
					break
				case 1:
					sysreg = REG_RVBAR_EL3
					break
				case 2:
					sysreg = REG_RMR_EL3
					break
				}
			} else if decode.Crm() == 12 {
				switch decode.Op2() {
				case 4:
					sysreg = REG_ICC_MCTLR_EL3
					break
				case 5:
					sysreg = REG_ICC_MSRE_EL3
					break
				case 7:
					sysreg = REG_ICC_MGRPEN1_EL3
					break
				}
			}
		}
		break
	case 13:
		{
			if 4 > decode.Crm()&3 && decode.Crm() > 0 {
				var sysregs = [4][9]SystemReg{
					{REG_AMEVCNTVOFF00_EL2, REG_AMEVCNTVOFF01_EL2, REG_AMEVCNTVOFF02_EL2, REG_AMEVCNTVOFF03_EL2,
						REG_AMEVCNTVOFF04_EL2, REG_AMEVCNTVOFF05_EL2, REG_AMEVCNTVOFF06_EL2, REG_AMEVCNTVOFF07_EL2},
					{REG_AMEVCNTVOFF08_EL2, REG_AMEVCNTVOFF09_EL2, REG_AMEVCNTVOFF010_EL2, REG_AMEVCNTVOFF011_EL2,
						REG_AMEVCNTVOFF012_EL2, REG_AMEVCNTVOFF013_EL2, REG_AMEVCNTVOFF014_EL2, REG_AMEVCNTVOFF015_EL2},
					{REG_AMEVCNTVOFF10_EL2, REG_AMEVCNTVOFF11_EL2, REG_AMEVCNTVOFF12_EL2, REG_AMEVCNTVOFF13_EL2,
						REG_AMEVCNTVOFF14_EL2, REG_AMEVCNTVOFF15_EL2, REG_AMEVCNTVOFF16_EL2, REG_AMEVCNTVOFF17_EL2},
					{REG_AMEVCNTVOFF18_EL2, REG_AMEVCNTVOFF19_EL2, REG_AMEVCNTVOFF110_EL2, REG_AMEVCNTVOFF111_EL2,
						REG_AMEVCNTVOFF112_EL2, REG_AMEVCNTVOFF113_EL2, REG_AMEVCNTVOFF114_EL2, REG_AMEVCNTVOFF115_EL2},
				}
				sysreg = sysregs[decode.Crm()&3][decode.Op2()]
				break
			}
			var sysregs = [8][5]SystemReg{
				{SYSREG_NONE, REG_CONTEXTIDR_EL1, SYSREG_NONE, SYSREG_NONE, REG_TPIDR_EL1},
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, REG_TPIDR_EL0, REG_TPIDRRO_EL0, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, REG_TPIDR_EL2, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, REG_CONTEXTIDR_EL12, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, REG_TPIDR_EL3, SYSREG_NONE, SYSREG_NONE},
				{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
			}
			sysreg = sysregs[decode.Op1()][decode.Op2()]
			break
		}
	case 14:
		{
			if decode.Op1() == 3 {
				reg := ((decode.Crm() & 3) << 3) | decode.Op2()
				if (decode.Crm() >= 8 && decode.Crm() <= 10 && decode.Op2() <= 7) || (decode.Crm() == 11 && decode.Op2() <= 6) {
					sysreg = REG_PMEVCNTR0_EL0 + SystemReg(reg)
					break
				} else if (decode.Crm() >= 12 && decode.Crm() <= 14 && decode.Op2() <= 7) || (decode.Crm() == 15 && decode.Op2() <= 6) {
					sysreg = REG_PMEVTYPER0_EL0 + SystemReg(reg)
					break
				} else if decode.Crm() == 15 && decode.Op2() == 7 {
					sysreg = REG_PMCCFILTR_EL0
					break
				} else if decode.Crm() == 0 && decode.Op2() == 5 {
					sysreg = REG_CNTPCTSS_EL0
					break
				} else if decode.Crm() == 0 && decode.Op2() == 6 {
					sysreg = REG_CNTVCTSS_EL0
					break
				}
			} else if decode.Op1() == 4 && decode.Crm() == 0 && decode.Op2() == 3 {
				sysreg = REG_CNTVOFF_EL2
				break
			} else if decode.Op2() > 2 && decode.Op1() == 4 {
				switch decode.Op2() {
				case 4:
					sysreg = REG_CNTSCALE_EL2
				case 5:
					sysreg = REG_CNTISCALE_EL2
				case 6:
					sysreg = REG_CNTPOFF_EL2
				case 7:
					sysreg = REG_CNTVFRQ_EL2
				}
				break
			}
			var sysregs = [8][4][3]SystemReg{
				{
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_CNTKCTL_EL1, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}, {
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}, {
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}, {
					{REG_CNTFRQ_EL0, REG_CNTPCT_EL0, REG_CNTVCT_EL0},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_CNTP_TVAL_EL0, REG_CNTP_CTL_EL0, REG_CNTP_CVAL_EL0},
					{REG_CNTV_TVAL_EL0, REG_CNTV_CTL_EL0, REG_CNTV_CVAL_EL0},
				}, {
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_CNTHCTL_EL2, SYSREG_NONE, SYSREG_NONE},
					{REG_CNTHP_TVAL_EL2, REG_CNTHP_CTL_EL2, REG_CNTHP_CVAL_EL2},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}, {
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_CNTKCTL_EL12, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, REG_CNTP_CTL_EL02, REG_CNTP_CVAL_EL02},
					{SYSREG_NONE, REG_CNTV_CTL_EL02, REG_CNTV_CVAL_EL02},
				}, {
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				}, {
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
					{REG_CNTPS_TVAL_EL1, REG_CNTPS_CTL_EL1, REG_CNTPS_CVAL_EL1},
					{SYSREG_NONE, SYSREG_NONE, SYSREG_NONE},
				},
			}
			if decode.Op1() > 7 || decode.Crm() > 3 || decode.Op2() > 2 {
				return nil, failedToDecodeInstruction
			}
			sysreg = sysregs[decode.Op1()][decode.Crm()][decode.Op2()]
		}
		break
	case 11:
		fallthrough
	case 15:
		{
			//Switch operands depending on load vs store
			op1 := ^(^decode.L()) & 1
			op2 := ^decode.L() & 1
			i.operation = operation[op1]
			i.operands[op1].OpClass = IMPLEMENTATION_SPECIFIC
			i.operands[op1].Reg[0] = decode.Op0()
			i.operands[op1].Reg[1] = decode.Op1()
			i.operands[op1].Reg[2] = decode.Crn()
			i.operands[op1].Reg[3] = decode.Crm()
			i.operands[op1].Reg[4] = decode.Op2()
			i.operands[op2].OpClass = REG
			i.operands[op2].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))
			return i, nil
		}
	}
	op1 := ^(^decode.L()) & 1
	op2 := ^decode.L() & 1
	i.operation = operation[op1]
	i.operands[op1].OpClass = SYS_REG
	i.operands[op1].Reg[0] = uint32(sysreg)
	i.operands[op2].OpClass = REG
	i.operands[op2].Reg[0] = reg(REGSET_ZR, REG_X_BASE, int(decode.Rt()))

	if sysreg == SYSREG_NONE {
		return nil, failedToDecodeInstruction
	}
	return i, nil
}

func (i *Instruction) decompose_system() (*Instruction, error) {
	decode := System(i.raw)
	switch decode.Op0() {
	case 0: //C5.2.3 - Architectural hints, barriers and CLREX, PSTATE Access
		return i.decompose_system_arch_hints(decode)
	case 1: //C5.2.4 - Cache maintenance, TLB maintenance, and address translation instructions
		return i.decompose_system_cache_maintenance(decode)
	case 2: //C5.2.5 - Moves to and from debug and trace system registers
		return i.decompose_system_debug_and_trace_regs(decode)
	case 3: //C5.2.6 - Moves to and from non-debug System registers and special purpose registers
		return i.decompose_system_debug_and_trace_regs2(decode)
	}
	return nil, failedToDecodeInstruction
}

func (i *Instruction) decompose_test_branch_imm() (*Instruction, error) {
	/* C4.2.5 Test & branch (immediate)
	 *
	 * TBZ <R><t>, #<imm>, <label>
	 * TBNZ <R><t>, #<imm>, <label>
	 */
	decode := TestAndBranch(i.raw)
	var operation = [2]Operation{ARM64_TBZ, ARM64_TBNZ}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = reg(REGSET_ZR, int(regSize[decode.B5()]), int(decode.Rt()))
	i.operands[1].OpClass = IMM32
	i.operands[1].Immediate = uint64(decode.B5()<<5 | decode.B40())

	i.operands[2].OpClass = LABEL
	i.operands[2].Immediate = i.address + uint64(decode.Imm()<<2)

	return i, nil
}

func (i *Instruction) decompose_unconditional_branch() (*Instruction, error) {
	/*
	 * B <label>
	 * BL <label>
	 */
	decode := UnconditionalBranch(i.raw)
	// fmt.Println(decode)
	var operation = []Operation{ARM64_B, ARM64_BL}
	i.operation = operation[decode.Op()]
	i.operands[0].OpClass = LABEL
	i.operands[0].Immediate = i.address + uint64(decode.Imm()<<2)
	if decode.Imm() < 0 {
		i.operands[0].SignedImm = 1
	}
	return i, nil
}

func (i *Instruction) decompose_unconditional_branch_reg() (*Instruction, error) {
	/* C4.2.7 Unconditional branch (register)
	 *
	 * BR <Xn>
	 * BLR <Xn>
	 * RET {<Xn>}
	 * ERET
	 * DRPS
	 */
	decode := UnconditionalBranchReg(i.raw)
	var operations = [10][4]Operation{
		{ARM64_BR, ARM64_UNDEFINED, ARM64_BRAAZ, ARM64_BRABZ},
		{ARM64_BLR, ARM64_UNDEFINED, ARM64_BLRAAZ, ARM64_BLRABZ},
		{ARM64_RET, ARM64_UNDEFINED, ARM64_RETAA, ARM64_RETAB},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},

		{ARM64_ERET, ARM64_UNDEFINED, ARM64_ERETAA, ARM64_ERETAB},
		{ARM64_DRPS, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED},

		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_BRAA, ARM64_BRAB},
		{ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_BLRAA, ARM64_BLRAB},
	}

	if decode.Opc() > 9 {
		return nil, failedToDecodeInstruction
	}
	if decode.Op3() > 3 {
		return nil, failedToDecodeInstruction
	}
	if decode.Opc() < 8 && decode.Op3() != 0 && decode.Op4() != 0x1f {
		return nil, failedToDecodeInstruction
	}
	if decode.Op3() == 0 && decode.Op4() != 0 {
		return nil, failedToDecodeInstruction
	}
	if decode.Op2() != 0x1f {
		return nil, failedToDecodeInstruction
	}

	i.operation = operations[decode.Opc()][decode.Op3()]
	r := reg(1, REG_X_BASE, int(decode.Rn()))

	switch decode.Opc() {
	case 2: // RET
		if decode.Op3() == 0 {
			if r == uint32(REG_X30) {
				return i, nil
			}
			break
		}
		fallthrough
	case 4: // ERET
		fallthrough
	case 5: // DRPS
		if decode.Rn() != 0x1f {
			return nil, failedToDecodeInstruction
		}
		return i, nil
	case 8:
		fallthrough
	case 9:
		i.operands[1].OpClass = REG
		i.operands[1].Reg[0] = reg(REGSET_SP, REG_X_BASE, int(decode.Op4()))
		break
	default:
		break
	}

	i.operands[0].OpClass = REG
	i.operands[0].Reg[0] = r

	return i, nil
}

func decompose(instructionValue uint32, address uint64) (*Instruction, error) {

	instruction := &Instruction{
		raw:       instructionValue,
		address:   address,
		operation: ARM64_UNDEFINED,
	}
	// fmt.Printf("main SWITCH: %d\n", ExtractBits(instructionValue, 25, 4))
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
			return instruction.decompose_pc_rel_addr()
		case 2:
			return instruction.decompose_add_sub_imm()
		case 3:
			return instruction.decompose_add_sub_imm_tags()
		case 4:
			return instruction.decompose_logical_imm()
		case 5:
			return instruction.decompose_move_wide_imm()
		case 6:
			return instruction.decompose_bitfield()
		case 7:
			return instruction.decompose_extract()
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
			return instruction.decompose_unconditional_branch()
		case 0x1a:
			fallthrough
		case 0x5a:
			return instruction.decompose_compare_branch_imm()
		case 0x1b:
			fallthrough
		case 0x5b:
			return instruction.decompose_test_branch_imm()
		case 0x2a:
			return instruction.decompose_conditional_branch()
		case 0x6a:
			if ExtractBits(instructionValue, 24, 1) == 0 {
				return instruction.decompose_exception_generation()
			} else if ExtractBits(instructionValue, 22, 3) == 4 {
				return instruction.decompose_system()
			}
			return instruction, nil // TODO error  ?
		case 0x6b:
			return instruction.decompose_unconditional_branch_reg()
		default:
			return nil, failedToDecodeInstruction
		}
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
					return instruction.decompose_simd_load_store_multiple()
				}
				if op2 == 1 && (op3>>5) == 0 {
					return instruction.decompose_simd_load_store_multiple_post_idx()
				}
				if op2 == 2 && (op3&0x1f) == 0 {
					return instruction.decompose_simd_load_store_single()
				}
				if op2 == 3 {
					return instruction.decompose_simd_load_store_single_post_idx()
				}
			}

			if ExtractBits(instructionValue, 24, 6) == 25 {
				if op0 == 13 && ExtractBits(instructionValue, 21, 1) != 0 {
					return instruction.decompose_load_store_mem_tags()
				}
				return instruction.decompose_load_store_unscaled()
			}

			if (op0&3) == 0 && op1 == 0 && (op2>>1) == 0 {
				return instruction.decompose_load_store_exclusive()
			}
			if (op0&3) == 1 && (op2>>1) == 0 {
				return instruction.decompose_load_register_literal()
			}

			if (op0 & 3) == 2 {
				if op2 == 0 {
					return instruction.decompose_load_store_no_allocate_pair_offset()
				}
				if op2 == 1 {
					return instruction.decompose_load_store_reg_pair_post_idx()
				}
				if op2 == 2 {
					return instruction.decompose_load_store_reg_pair_offset()
				}
				if op2 == 3 {
					return instruction.decompose_load_store_reg_pair_pre_idx()
				}
			}

			if (op0 & 3) == 3 {
				if (op2 >> 1) == 0 {
					if (op3 >> 5) == 0 {
						if op4 == 0 {
							return instruction.decompose_load_store_reg_unscalled_imm()
						}
						if op4 == 1 {
							return instruction.decompose_load_store_imm_post_idx()
						}
						if op4 == 2 {
							return instruction.decompose_load_store_reg_unpriv()
						}
						if op4 == 3 {
							return instruction.decompose_load_store_reg_imm_pre_idx()
						}
					}
					if (op3 >> 5) == 1 {
						if op4 == 0 {
							// if ExtractBits(instructionValue, 24, 6) == 56 {
							return instruction.decompose_atomic_memory_ops()
						}
						if op4 == 2 {
							return instruction.decompose_load_store_reg_reg_offset()
						}
						if op4 == 1 || op4 == 3 {
							return instruction.decompose_load_store_pac()
						}
					}
				}
				return instruction.decompose_load_store_reg_unsigned_imm()
			}
			break
		}
	case 5:
		fallthrough
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
			return instruction.decompose_logical_shifted_reg()
		case 0x58:
			fallthrough
		case 0x5a:
			fallthrough
		case 0x5c:
			fallthrough
		case 0x5e:
			return instruction.decompose_add_sub_shifted_reg()
		case 0x59:
			fallthrough
		case 0x5b:
			fallthrough
		case 0x5d:
			fallthrough
		case 0x5f:
			return instruction.decompose_add_sub_extended_reg()
		case 0xd0:
			return instruction.decompose_add_sub_carry()
		case 0xd2:
			if ExtractBits(instructionValue, 11, 1) == 1 {
				return instruction.decompose_conditional_compare_imm()
			}
			return instruction.decompose_conditional_compare_reg()
		case 0xd4:
			return instruction.decompose_conditional_select()
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
			return instruction.decompose_data_processing_3()
		case 0xd6:
			if ExtractBits(instructionValue, 30, 1) == 1 {
				return instruction.decompose_data_processing_1()
			}
			return instruction.decompose_data_processing_2()
		default:
			return nil, failedToDecodeInstruction
		}
	case 7:
		switch ExtractBits(instructionValue, 24, 5) {
		case 14:
			fallthrough
		case 15:
			return instruction.decompose_floating_complex_multiply_accumulate()
		default:
			return nil, failedToDecodeInstruction
		}
	case 15:
		instruction.group = GROUP_DATA_PROCESSING_SIMD
		fmt.Printf("case 15: %#x\n", ExtractBits(instructionValue, 24, 8))
		switch ExtractBits(instructionValue, 24, 8) {
		case 0x1e:
			fallthrough
		case 0x3e:
			fallthrough
		case 0x9e:
			fallthrough
		case 0xbe:
			if ExtractBits(instructionValue, 21, 1) == 0 {
				return instruction.decompose_fixed_floating_conversion()
			}
			fmt.Printf("switch: %#x\n", ExtractBits(instructionValue, 10, 2))
			switch ExtractBits(instructionValue, 10, 2) {
			case 0:
				if (instructionValue & 0x1E7E0000) == 0x1E7E0000 {
					return instruction.decompose_floating_javascript_conversion()
				} else if ExtractBits(instructionValue, 12, 1) == 1 {
					return instruction.decompose_floating_imm()
				} else if ExtractBits(instructionValue, 12, 2) == 2 {
					return instruction.decompose_floating_compare()
				} else if ExtractBits(instructionValue, 12, 3) == 4 {
					return instruction.decompose_floating_data_processing1()
				} else if ExtractBits(instructionValue, 12, 4) == 0 {
					return instruction.decompose_floating_integer_conversion()
				}
				break
			case 1:
				return instruction.decompose_floating_conditional_compare()
			case 2:
				return instruction.decompose_floating_data_processing2()
			case 3:
				return instruction.decompose_floating_cselect()
			}
			break
		case 0x1f:
			fallthrough
		case 0x3f:
			fallthrough
		case 0x9f:
			fallthrough
		case 0xbf:
			return instruction.decompose_floating_data_processing3()
		case 0x0e:
			fallthrough
		case 0x2e:
			fallthrough
		case 0x4e:
			fallthrough
		case 0x6e:
			// fmt.Printf("case 0x6e: %d\n", ExtractBits(instructionValue, 21, 1))
			if ExtractBits(instructionValue, 21, 1) == 1 {
				// fmt.Printf("switch: %d\n", ExtractBits(instructionValue, 10, 2))
				switch ExtractBits(instructionValue, 10, 2) {
				case 1:
					fallthrough
				case 3:
					return instruction.decompose_simd_3_same()
				case 0:
					return instruction.decompose_simd_3_different()
				case 2:
					if ExtractBits(instructionValue, 17, 4) == 0 {
						return instruction.decompose_simd_2_reg_misc()
					} else if ExtractBits(instructionValue, 17, 4) == 8 {
						return instruction.decompose_simd_across_lanes()
					}
				}
			}
			// else if ExtractBits(instructionValue, 24, 5) == 15 || {

			// } else if ExtractBits(instructionValue, 24, 5) == 14 {

			// }
			if (instructionValue & 0x9fe08400) == 0x0e000400 {
				return instruction.decompose_simd_copy()
			}
			if (instructionValue & 0x003e0c00) == 0x00280800 {
				return instruction.decompose_cryptographic_aes()
			}
			if (instructionValue & 0xbf208c00) == 0x0e000000 {
				return instruction.decompose_simd_table_lookup()
			}
			if (instructionValue & 0xbf208c00) == 0x0e000800 {
				return instruction.decompose_simd_permute()
			}
			if (instructionValue & 0xbfe08400) == 0x2e000000 {
				return instruction.decompose_simd_extract()
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
				return instruction.decompose_simd_vector_indexed_element()
			}
			if ExtractBits(instructionValue, 19, 5) == 0 {
				return instruction.decompose_simd_modified_imm()
			}
			return instruction.decompose_simd_shift_imm()
		case 0x5e:
			fallthrough
		case 0x7e:
			if ExtractBits(instructionValue, 21, 1) == 1 {
				switch ExtractBits(instructionValue, 10, 2) {
				case 1:
					fallthrough
				case 3:
					return instruction.decompose_simd_scalar_3_same()
				case 0:
					return instruction.decompose_simd_scalar_3_different()
				case 2:
					if ExtractBits(instructionValue, 17, 4) == 0 {
						return instruction.decompose_simd_scalar_2_reg_misc()
					} else if ExtractBits(instructionValue, 17, 4) == 8 {
						return instruction.decompose_simd_scalar_pairwise()
					}
				}
			}
			if (instructionValue & 0xdfe08400) == 0x5e000400 {
				return instruction.decompose_simd_scalar_copy()
			} else if (instructionValue & 0xff208c00) == 0x5e000000 {
				return instruction.decompose_cryptographic_3_register_sha()
			} else if (instructionValue & 0xff3e0c00) == 0x5e280800 {
				return instruction.decompose_cryptographic_2_register_sha()
			}
			break
		case 0x5f:
			fallthrough
		case 0x7f:
			if ExtractBits(instructionValue, 10, 1) == 0 {
				return instruction.decompose_simd_scalar_indexed_element()
			} else if ExtractBits(instructionValue, 23, 1) == 0 && ExtractBits(instructionValue, 10, 1) == 1 {
				return instruction.decompose_simd_scalar_shift_imm()
			}
			break
		}
	default:
		return nil, failedToDecodeInstruction
	}

	return instruction, nil
}
