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

var (
	regSize     = []uint32{REG_W_BASE, REG_X_BASE}
	simdRegSize = []uint32{REG_S_BASE, REG_D_BASE, REG_Q_BASE}
	dataSize    = []uint8{32, 64}
)

func reg(arg1, arg2, arg3 int) uint32 {
	return uint32(regMap[arg1][arg2][arg3])
}

func (i *Instruction) deleteOperand(index int) {
	i.operands[0] = InstructionOperand{OpClass: NONE}
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
			{ARM64_LDR, REG_X_BASE, 0},
			{ARM64_LDRSW, REG_X_BASE, 1},
			{ARM64_PRFM, REG_W_BASE, 0},
		}, {
			{ARM64_LDR, REG_S_BASE, 0},
			{ARM64_LDR, REG_D_BASE, 0},
			{ARM64_LDR, REG_Q_BASE, 0},
			{ARM64_UNDEFINED, 0, 0},
		},
	}

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
	// i.operation = BF_GETI(30,1) ? ARM64_SUBG : ARM64_ADDG;

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
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_CRC32X,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_CRC32CX,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
			ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED, ARM64_UNDEFINED,
		},
	}

	decode := DataProcessing2(i.raw)

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
	if i.operation == ARM64_SUBPS && decode.S() == 1 && decode.Rd() == 0x1f {
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

	// LDST_TAGS decode = *(LDST_TAGS*)&instructionValue;
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
		fmt.Println("case 9:", ExtractBits(instructionValue, 23, 3))
		switch ExtractBits(instructionValue, 23, 3) {
		case 0:
			fallthrough
		case 1:
			// return aarch64_decompose_pc_rel_addr(instructionValue, instruction, address)
		case 2:
			return instruction.decompose_add_sub_imm()
		case 3:
			return instruction.decompose_add_sub_imm_tags()
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
				return instruction.decompose_load_store_mem_tags()
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
		fmt.Printf("case 13: 0x%x\n", ExtractBits(instructionValue, 21, 8))
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
			return instruction.decompose_data_processing_3()
		case 0xd6:
			if ExtractBits(instructionValue, 30, 1) == 1 {
				return instruction.decompose_data_processing_1()
			}
			return instruction.decompose_data_processing_2()
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
