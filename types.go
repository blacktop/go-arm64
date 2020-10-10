package arm64

const MAX_OPERANDS = 5

type Operation uint32

const (
	ARM64_UNDEFINED Operation = iota
	ABS
	ADC
	ADCS
	ADD
	ADDG //Added for MTE
	ADDHN
	ADDHN2
	ADDP
	ADDS
	ADDV
	ADR
	ADRP
	AESD
	AESE
	AESIMC
	AESMC
	AND
	ANDS
	ASR
	AT
	AUTDA     //Added for 8.3
	AUTDB     //Added for 8.3
	AUTDZA    //Added for 8.3
	AUTDZB    //Added for 8.3
	AUTIA     //Added for 8.3
	AUTIA1716 //Added for 8.3
	AUTIASP   //Added for 8.3
	AUTIAZ    //Added for 8.3
	AUTIB     //Added for 8.3
	AUTIB1716 //Added for 8.3
	AUTIBSP   //Added for 8.3
	AUTIBZ    //Added for 8.3
	AUTIZA    //Added for 8.3
	AUTIZB    //Added for 8.3
	B
	B_AL
	B_CC
	B_CS
	B_EQ
	BFI
	BFM
	BFXIL
	B_GE
	B_GT
	B_HI
	BIC
	BICS
	BIF
	BIT
	BL
	B_LE
	BLR
	BLRAA
	BLRAAZ
	BLRAB
	BLRABZ
	B_LS
	B_LT
	B_MI
	B_NE
	B_NV
	B_PL
	BR
	BRAA
	BRAAZ
	BRAB
	BRABZ
	BRK
	BSL
	B_VC
	B_VS
	CBNZ
	CBZ
	CCMN
	CCMP
	CINC
	CINV
	CLREX
	CLS
	CLZ
	CMEQ
	CMGE
	CMGT
	CMHI
	CMHS
	CMLE
	CMLT
	CMN
	CMP
	CMPP //Added for MTE
	CMTST
	CNEG
	CNT
	CRC32B
	CRC32CB
	CRC32CH
	CRC32CW
	CRC32CX
	CRC32H
	CRC32W
	CRC32X
	CSEL
	CSET
	CSETM
	CSINC
	CSINV
	CSNEG
	DC
	DCPS1
	DCPS2
	DCPS3
	DMB
	DRPS
	DSB
	DUP
	EON
	EOR
	ERET
	ERETAA //Added for 8.3
	ERETAB //Added for 8.3
	ESB    //Added for 8.2
	EXT
	EXTR
	FABD
	FABS
	FACGE
	FACGT
	FADD
	FADDP
	FCCMP
	FCCMPE
	FCMEQ
	FCMGE
	FCMGT
	FCMLE
	FCMLT
	FCMP
	FCMPE
	FCSEL
	FCTNS
	FCTNU
	FCVT
	FCVTAS
	FCVTAU
	FCVTL
	FCVTL2
	FCVTMS
	FCVTMU
	FCVTN
	FCVTN2
	FCVTNS
	FCVTNU
	FCVTPS
	FCVTPU
	FCVTXN
	FCVTXN2
	FCVTZS
	FCVTZU
	FDIV
	FMADD
	FMAX
	FMAXNM
	FMAXNMP
	FMAXNMV
	FMAXP
	FMAXV
	FMIN
	FMINNM
	FMINNMP
	FMINNMV
	FMINP
	FMINV
	FMLA
	FMLS
	FMOV
	FMSUB
	FMUL
	FMULX
	FNEG
	FNMADD
	FNMSUB
	FNMUL
	FRECPE
	FRECPS
	FRECPX
	FRINTA
	FRINTI
	FRINTM
	FRINTN
	FRINTP
	FRINTX
	FRINTZ
	FRSQRTE
	FRSQRTS
	FSQRT
	FSUB
	GMI //Added for MTE
	HINT
	HLT
	HVC
	IC
	INS
	IRG //Added for MTE
	ISB
	LD1
	LD1R
	LD2
	LD2R
	LD3
	LD3R
	LD4
	LD4R
	LDAR
	LDARB
	LDARH
	LDAXP
	LDAXR
	LDAXRB
	LDAXRH
	LDG  //Added for MTE
	LDGM //Added for MTE
	LDNP
	LDP
	LDPSW
	LDR
	LDRAA //Added for 8.3
	LDRAB //Added for 8.3
	LDRB
	LDRH
	LDRSB
	LDRSH
	LDRSW
	LDTR
	LDTRB
	LDTRH
	LDTRSB
	LDTRSH
	LDTRSW
	LDUR
	LDURB
	LDURH
	LDURSB
	LDURSH
	LDURSW
	LDXP
	LDXR
	LDXRB
	LDXRH
	LSL
	LSR
	MADD
	MLA
	MLS
	MNEG
	MOV
	MOVI
	MOVK
	MOVN
	MOVZ
	MRS
	MSR
	MSUB
	MUL
	MVN
	MVNI
	NEG
	NEGS
	NGC
	NGCS
	NOP
	NOT
	ORN
	ORR
	PACDA     //Added for 8.3
	PACDB     //Added for 8.3
	PACDZA    //Added for 8.3
	PACDZB    //Added for 8.3
	PACIA     //Added for 8.3
	PACIA1716 //Added for 8.3
	PACIASP   //Added for 8.3
	PACIAZ    //Added for 8.3
	PACIB     //Added for 8.3
	PACIB1716 //Added for 8.3
	PACIBSP   //Added for 8.3
	PACIBZ    //Added for 8.3
	PACIZA    //Added for 8.3
	PACIZB    //Added for 8.3
	PMUL
	PMULL
	PMULL2
	PRFM
	PRFUM
	PSBCSYNC //Added for 8.2
	RADDHN
	RADDHN2
	RBIT
	RET
	RETAA //Added for 8.3
	RETAB //Added for 8.3
	REV
	REV16
	REV32
	REV64
	ROR
	RSHRN
	RSHRN2
	RSUBHN
	RSUBHN2
	SABA
	SABAL
	SABAL2
	SABD
	SABDL
	SABDL2
	SADALP
	SADDL
	SADDL2
	SADDLP
	SADDLV
	SADDW
	SADDW2
	SBC
	SBCS
	SBFIZ
	SBFM
	SBFX
	SCVTF
	SDIV
	SEV
	SEVL
	SHA1C
	SHA1H
	SHA1M
	SHA1P
	SHA1SU0
	SHA1SU1
	SHA256H
	SHA256H2
	SHA256SU0
	SHA256SU1
	SHADD
	SHL
	SHLL
	SHLL2
	SHRN
	SHRN2
	SHSUB
	SLI
	SMADDL
	SMAX
	SMAXP
	SMAXV
	SMC
	SMIN
	SMINP
	SMINV
	SMLAL
	SMLAL2
	SMLSL
	SMLSL2
	SMNEGL
	SMOV
	SMSUBL
	SMULH
	SMULL
	SMULL2
	SQABS
	SQADD
	SQDMLAL
	SQDMLAL2
	SQDMLSL
	SQDMLSL2
	SQDMULH
	SQDMULL
	SQDMULL2
	SQNEG
	SQRDMULH
	SQRSHL
	SQRSHRN
	SQRSHRN2
	SQRSHRUN
	SQRSHRUN2
	SQSHL
	SQSHLU
	SQSHRN
	SQSHRN2
	SQSHRUN
	SQSHRUN2
	SQSUB
	SQXTN
	SQXTN2
	SQXTUN
	SQXTUN2
	SRHADD
	SRI
	SRSHL
	SRSHR
	SRSRA
	SSHL
	SSHLL
	SSHLL2
	SSHR
	SSRA
	SSUBL
	SSUBL2
	SSUBW
	SSUBW2
	ST1
	ST2
	ST2G //Added for MTE
	ST3
	ST4
	STG  //Added for MTE
	STGM //Added for MTE
	STGP //Added for MTE
	STLR
	STLRB
	STLRH
	STLXP
	STLXR
	STLXRB
	STLXRH
	STNP
	STP
	STR
	STRB
	STRH
	STTR
	STTRB
	STTRH
	STUR
	STURB
	STURH
	STXP
	STXR
	STXRB
	STXRH
	STZ2G //Added for MTE
	STZG  //Added for MTE
	STZGM //Added for MTE
	SUB
	SUBG //Added for MTE
	SUBHN
	SUBHN2
	SUBP  //Added for MTE
	SUBPS //Added for MTE
	SUBS
	SUQADD
	SVC
	SXTB
	SXTH
	SXTW
	SYS
	SYSL
	TBL
	TBNZ
	TBX
	TBZ
	TLBI
	TRN1
	TRN2
	TST
	UABA
	UABAL
	UABAL2
	UABD
	UABDL
	UABDL2
	UADALP
	UADDL
	UADDL2
	UADDLP
	UADDLV
	UADDW
	UADDW2
	UBFIZ
	UBFM
	UBFX
	UCVTF
	UDIV
	UHADD
	UHSUB
	UMADDL
	UMAX
	UMAXP
	UMAXV
	UMIN
	UMINP
	UMINV
	UMLAL
	UMLAL2
	UMLSL
	UMLSL2
	UMNEGL
	UMOV
	UMSUBL
	UMULH
	UMULL
	UMULL2
	UQADD
	UQRSHL
	UQRSHRN
	UQRSHRN2
	UQSHL
	UQSHRN
	UQSHRN2
	UQSUB
	UQXTN
	UQXTN2
	URECPE
	URHADD
	URSHL
	URSHR
	URSQRTE
	URSRA
	USHL
	USHLL
	USHLL2
	USHR
	USQADD
	USRA
	USUBL
	USUBL2
	USUBW
	USUBW2
	UXTB
	UXTH
	UZP1
	UZP2
	WFE
	WFI
	XPACD   //Added for 8.3
	XPACI   //Added for 8.3
	XPACLRI //Added for 8.3
	XTN
	XTN2
	YIELD
	ZIP1
	ZIP2

	AMD64_END_TYPE //Not real instruction
)

func (o Operation) String() string {
	return []string{
		"UNDEFINED",
		"abs",
		"adc",
		"adcs",
		"add",
		"addg", //Added for MTE
		"addhn",
		"addhn2",
		"addp",
		"adds",
		"addv",
		"adr",
		"adrp",
		"aesd",
		"aese",
		"aesimc",
		"aesmc",
		"and",
		"ands",
		"asr",
		"at",
		"autda",     //Added for 8.3
		"autdb",     //Added for 8.3
		"autdza",    //Added for 8.3
		"autdzb",    //Added for 8.3
		"autia",     //Added for 8.3
		"autia1716", //Added for 8.3
		"autiasp",   //Added for 8.3
		"autiaz",    //Added for 8.3
		"autib",     //Added for 8.3
		"autib1716", //Added for 8.3
		"autibsp",   //Added for 8.3
		"autibz",    //Added for 8.3
		"autiza",    //Added for 8.3
		"autizb",    //Added for 8.3
		"b",
		"b.al",
		"b.cc",
		"b.cs",
		"b.eq",
		"bfi",
		"bfm",
		"bfxil",
		"b.ge",
		"b.gt",
		"b.hi",
		"bic",
		"bics",
		"bif",
		"bit",
		"bl",
		"b.le",
		"blr",
		"blraa",
		"blraaz",
		"blrab",
		"blrabz",
		"b.ls",
		"b.lt",
		"b.mi",
		"b.ne",
		"b.nv",
		"b.pl",
		"br",
		"braa",
		"braaz",
		"brab",
		"brabz",
		"brk",
		"bsl",
		"b.vc",
		"b.vs",
		"cbnz",
		"cbz",
		"ccmn",
		"ccmp",
		"cinc",
		"cinv",
		"clrex",
		"cls",
		"clz",
		"cmeq",
		"cmge",
		"cmgt",
		"cmhi",
		"cmhs",
		"cmle",
		"cmlt",
		"cmn",
		"cmp",
		"cmpp", //Added for MTE
		"cmtst",
		"cneg",
		"cnt",
		"crc32b",
		"crc32cb",
		"crc32ch",
		"crc32cw",
		"crc32cx",
		"crc32h",
		"crc32w",
		"crc32x",
		"csel",
		"cset",
		"csetm",
		"csinc",
		"csinv",
		"csneg",
		"dc",
		"dcps1",
		"dcps2",
		"dcps3",
		"dmb",
		"drps",
		"dsb",
		"dup",
		"eon",
		"eor",
		"eret",
		"eretaa",
		"eretab",
		"esb", //Added for 8.2
		"ext",
		"extr",
		"fabd",
		"fabs",
		"facge",
		"facgt",
		"fadd",
		"faddp",
		"fccmp",
		"fccmpe",
		"fcmeq",
		"fcmge",
		"fcmgt",
		"fcmle",
		"fcmlt",
		"fcmp",
		"fcmpe",
		"fcsel",
		"fctns",
		"fctnu",
		"fcvt",
		"fcvtas",
		"fcvtau",
		"fcvtl",
		"fcvtl2",
		"fcvtms",
		"fcvtmu",
		"fcvtn",
		"fcvtn2",
		"fcvtns",
		"fcvtnu",
		"fcvtps",
		"fcvtpu",
		"fcvtxn",
		"fcvtxn2",
		"fcvtzs",
		"fcvtzu",
		"fdiv",
		"fmadd",
		"fmax",
		"fmaxnm",
		"fmaxnmp",
		"fmaxnmv",
		"fmaxp",
		"fmaxv",
		"fmin",
		"fminnm",
		"fminnmp",
		"fminnmv",
		"fminp",
		"fminv",
		"fmla",
		"fmls",
		"fmov",
		"fmsub",
		"fmul",
		"fmulx",
		"fneg",
		"fnmadd",
		"fnmsub",
		"fnmul",
		"frecpe",
		"frecps",
		"frecpx",
		"frinta",
		"frinti",
		"frintm",
		"frintn",
		"frintp",
		"frintx",
		"frintz",
		"frsqrte",
		"frsqrts",
		"fsqrt",
		"fsub",
		"gmi", //Added for MTE
		"hint",
		"hlt",
		"hvc",
		"ic",
		"ins",
		"irg", //Added for MTE
		"isb",
		"ld1",
		"ld1r",
		"ld2",
		"ld2r",
		"ld3",
		"ld3r",
		"ld4",
		"ld4r",
		"ldar",
		"ldarb",
		"ldarh",
		"ldaxp",
		"ldaxr",
		"ldaxrb",
		"ldaxrh",
		"ldg",  //Added for MTE
		"ldgm", //Added for MTE
		"ldnp",
		"ldp",
		"ldpsw",
		"ldr",
		"ldraa",
		"ldrab",
		"ldrb",
		"ldrh",
		"ldrsb",
		"ldrsh",
		"ldrsw",
		"ldtr",
		"ldtrb",
		"ldtrh",
		"ldtrsb",
		"ldtrsh",
		"ldtrsw",
		"ldur",
		"ldurb",
		"ldurh",
		"ldursb",
		"ldursh",
		"ldursw",
		"ldxp",
		"ldxr",
		"ldxrb",
		"ldxrh",
		"lsl",
		"lsr",
		"madd",
		"mla",
		"mls",
		"mneg",
		"mov",
		"movi",
		"movk",
		"movn",
		"movz",
		"mrs",
		"msr",
		"msub",
		"mul",
		"mvn",
		"mvni",
		"neg",
		"negs",
		"ngc",
		"ngcs",
		"nop",
		"not",
		"orn",
		"orr",
		"pacda",     //Added for 8.3
		"pacdb",     //Added for 8.3
		"pacdza",    //Added for 8.3
		"pacdzb",    //Added for 8.3
		"pacia",     //Added for 8.3
		"pacia1716", //Added for 8.3
		"paciasp",   //Added for 8.3
		"paciaz",    //Added for 8.3
		"pacib",     //Added for 8.3
		"pacib1716", //Added for 8.3
		"pacibsp",   //Added for 8.3
		"pacibz",    //Added for 8.3
		"paciza",    //Added for 8.3
		"pacizb",    //Added for 8.3
		"pmul",
		"pmull",
		"pmull2",
		"prfm",
		"prfum",
		"psb", //Added for 8.2
		"raddhn",
		"raddhn2",
		"rbit",
		"ret",
		"retaa", //Added for 8.3
		"retab", //Added for 8.3
		"rev",
		"rev16",
		"rev32",
		"rev64",
		"ror",
		"rshrn",
		"rshrn2",
		"rsubhn",
		"rsubhn2",
		"saba",
		"sabal",
		"sabal2",
		"sabd",
		"sabdl",
		"sabdl2",
		"sadalp",
		"saddl",
		"saddl2",
		"saddlp",
		"saddlv",
		"saddw",
		"saddw2",
		"sbc",
		"sbcs",
		"sbfiz",
		"sbfm",
		"sbfx",
		"scvtf",
		"sdiv",
		"sev",
		"sevl",
		"sha1c",
		"sha1h",
		"sha1m",
		"sha1p",
		"sha1su0",
		"sha1su1",
		"sha256h",
		"sha256h2",
		"sha256su0",
		"sha256su1",
		"shadd",
		"shl",
		"shll",
		"shll2",
		"shrn",
		"shrn2",
		"shsub",
		"sli",
		"smaddl",
		"smax",
		"smaxp",
		"smaxv",
		"smc",
		"smin",
		"sminp",
		"sminv",
		"smlal",
		"smlal2",
		"smlsl",
		"smlsl2",
		"smnegl",
		"smov",
		"smsubl",
		"smulh",
		"smull",
		"smull2",
		"sqabs",
		"sqadd",
		"sqdmlal",
		"sqdmlal2",
		"sqdmlsl",
		"sqdmlsl2",
		"sqdmulh",
		"sqdmull",
		"sqdmull2",
		"sqneg",
		"sqrdmulh",
		"sqrshl",
		"sqrshrn",
		"sqrshrn2",
		"sqrshrun",
		"sqrshrun2",
		"sqshl",
		"sqshlu",
		"sqshrn",
		"sqshrn2",
		"sqshrun",
		"sqshrun2",
		"sqsub",
		"sqxtn",
		"sqxtn2",
		"sqxtun",
		"sqxtun2",
		"srhadd",
		"sri",
		"srshl",
		"srshr",
		"srsra",
		"sshl",
		"sshll",
		"sshll2",
		"sshr",
		"ssra",
		"ssubl",
		"ssubl2",
		"ssubw",
		"ssubw2",
		"st1",
		"st2",
		"st2g", //Added for MTE
		"st3",
		"st4",
		"stg",  //Added for MTE
		"stgm", //Added for MTE
		"stgp", //Added for MTE
		"stlr",
		"stlrb",
		"stlrh",
		"stlxp",
		"stlxr",
		"stlxrb",
		"stlxrh",
		"stnp",
		"stp",
		"str",
		"strb",
		"strh",
		"sttr",
		"sttrb",
		"sttrh",
		"stur",
		"sturb",
		"sturh",
		"stxp",
		"stxr",
		"stxrb",
		"stxrh",
		"stz2g", //Added for MTE
		"stzg",  //Added for MTE
		"stzgm", //Added for MTE
		"sub",
		"subg", //Added for MTE
		"subhn",
		"subhn2",
		"subp",  //Added for MTE
		"subps", //Added for MTE
		"subs",
		"suqadd",
		"svc",
		"sxtb",
		"sxth",
		"sxtw",
		"sys",
		"sysl",
		"tbl",
		"tbnz",
		"tbx",
		"tbz",
		"tlbi",
		"trn1",
		"trn2",
		"tst",
		"uaba",
		"uabal",
		"uabal2",
		"uabd",
		"uabdl",
		"uabdl2",
		"uadalp",
		"uaddl",
		"uaddl2",
		"uaddlp",
		"uaddlv",
		"uaddw",
		"uaddw2",
		"ubfiz",
		"ubfm",
		"ubfx",
		"ucvtf",
		"udiv",
		"uhadd",
		"uhsub",
		"umaddl",
		"umax",
		"umaxp",
		"umaxv",
		"umin",
		"uminp",
		"uminv",
		"umlal",
		"umlal2",
		"umlsl",
		"umlsl2",
		"umnegl",
		"umov",
		"umsubl",
		"umulh",
		"umull",
		"umull2",
		"uqadd",
		"uqrshl",
		"uqrshrn",
		"uqrshrn2",
		"uqshl",
		"uqshrn",
		"uqshrn2",
		"uqsub",
		"uqxtn",
		"uqxtn2",
		"urecpe",
		"urhadd",
		"urshl",
		"urshr",
		"ursqrte",
		"ursra",
		"ushl",
		"ushll",
		"ushll2",
		"ushr",
		"usqadd",
		"usra",
		"usubl",
		"usubl2",
		"usubw",
		"usubw2",
		"uxtb",
		"uxth",
		"uzp1",
		"uzp2",
		"wfe",
		"wfi",
		"xpacd",
		"xpaci",
		"xpaclri",
		"xtn",
		"xtn2",
		"yield",
		"zip1",
		"zip2",
		"END_OPERATION_LIST", //NOT AN INSTRUCTION
	}[o]
}

//---------------------------------------------
//C4.4 Data processing - immediate
//---------------------------------------------

type PcRelAddressing uint32

func (i PcRelAddressing) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i PcRelAddressing) Immhi() int32 {
	return int32(ExtractBits(uint32(i), 5, 19))
}
func (i PcRelAddressing) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i PcRelAddressing) Immlo() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i PcRelAddressing) Op() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type AddSubImm uint32

func (i AddSubImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i AddSubImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i AddSubImm) Imm() uint32 {
	return ExtractBits(uint32(i), 10, 12)
}
func (i AddSubImm) Shift() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i AddSubImm) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i AddSubImm) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i AddSubImm) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i AddSubImm) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type AddSubImmTags uint32

func (i AddSubImmTags) Xd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i AddSubImmTags) Xn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i AddSubImmTags) Uimm4() uint32 {
	return ExtractBits(uint32(i), 10, 4)
}
func (i AddSubImmTags) Op3() uint32 {
	return ExtractBits(uint32(i), 14, 2)
}
func (i AddSubImmTags) Uimm6() uint32 {
	return ExtractBits(uint32(i), 16, 6)
}
func (i AddSubImmTags) Padding() uint32 {
	return ExtractBits(uint32(i), 22, 10)
}

type LogicalImm uint32

func (i LogicalImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LogicalImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LogicalImm) Imms() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i LogicalImm) Immr() uint32 {
	return ExtractBits(uint32(i), 16, 6)
}
func (i LogicalImm) N() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i LogicalImm) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 6)
}
func (i LogicalImm) Opc() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i LogicalImm) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type MoveWideImm uint32

func (i MoveWideImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i MoveWideImm) Imm() uint32 {
	return ExtractBits(uint32(i), 5, 16)
}
func (i MoveWideImm) Hw() uint32 {
	return ExtractBits(uint32(i), 21, 2)
}
func (i MoveWideImm) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 6)
}
func (i MoveWideImm) Opc() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i MoveWideImm) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type Bitfield uint32

func (i Bitfield) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Bitfield) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Bitfield) Imms() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i Bitfield) Immr() uint32 {
	return ExtractBits(uint32(i), 16, 6)
}
func (i Bitfield) N() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i Bitfield) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 6)
}
func (i Bitfield) Opc() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i Bitfield) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type Extract uint32

func (i Extract) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Extract) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Extract) Imms() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i Extract) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i Extract) O0() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i Extract) N() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i Extract) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 6)
}
func (i Extract) Op21() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i Extract) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

//--------------------------------------------------------
// C4.2  Branches, exception generating and system instructions
//--------------------------------------------------------

type UnconditionalBranch uint32

func (i UnconditionalBranch) Imm() int32 {
	return int32(ExtractBits(uint32(i), 0, 26))
}
func (i UnconditionalBranch) Opcode() uint32 {
	return ExtractBits(uint32(i), 26, 5)
}
func (i UnconditionalBranch) Op() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type CompareBranchImm uint32

func (i CompareBranchImm) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i CompareBranchImm) Imm() int32 {
	return int32(ExtractBits(uint32(i), 5, 19))
}
func (i CompareBranchImm) Op() uint32 {
	return ExtractBits(uint32(i), 24, 1)
}
func (i CompareBranchImm) Opcode() uint32 {
	return ExtractBits(uint32(i), 25, 6)
}
func (i CompareBranchImm) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type TestAndBranch uint32

func (i TestAndBranch) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i TestAndBranch) Imm() int32 {
	return int32(ExtractBits(uint32(i), 5, 14))
}
func (i TestAndBranch) B40() uint32 {
	return ExtractBits(uint32(i), 19, 5)
}
func (i TestAndBranch) Op() uint32 {
	return ExtractBits(uint32(i), 24, 1)
}
func (i TestAndBranch) Opcode() uint32 {
	return ExtractBits(uint32(i), 25, 6)
}
func (i TestAndBranch) B5() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type ConditionalBranchImm uint32

func (i ConditionalBranchImm) Cond() uint32 {
	return ExtractBits(uint32(i), 0, 4)
}
func (i ConditionalBranchImm) O0() uint32 {
	return ExtractBits(uint32(i), 4, 1)
}
func (i ConditionalBranchImm) Imm() int32 {
	return int32(ExtractBits(uint32(i), 5, 19))
}
func (i ConditionalBranchImm) O1() uint32 {
	return ExtractBits(uint32(i), 24, 1)
}
func (i ConditionalBranchImm) Opcode() uint32 {
	return ExtractBits(uint32(i), 25, 7)
}

type ExceptionGeneration uint32

func (i ExceptionGeneration) Ll() uint32 {
	return ExtractBits(uint32(i), 0, 2)
}
func (i ExceptionGeneration) Op2() uint32 {
	return ExtractBits(uint32(i), 2, 3)
}
func (i ExceptionGeneration) Imm() uint32 {
	return ExtractBits(uint32(i), 5, 16)
}
func (i ExceptionGeneration) Opc() uint32 {
	return ExtractBits(uint32(i), 21, 3)
}
func (i ExceptionGeneration) Opcode() uint32 {
	return ExtractBits(uint32(i), 24, 8)
}

type System uint32

func (i System) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i System) Op2() uint32 {
	return ExtractBits(uint32(i), 5, 3)
}
func (i System) Crm() uint32 {
	return ExtractBits(uint32(i), 8, 4)
}
func (i System) Crn() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i System) Op1() uint32 {
	return ExtractBits(uint32(i), 16, 3)
}
func (i System) Op0() uint32 {
	return ExtractBits(uint32(i), 19, 2)
}
func (i System) L() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i System) Group1() uint32 {
	return ExtractBits(uint32(i), 22, 10)
}

type UnconditionalBranchReg uint32

func (i UnconditionalBranchReg) Op4() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i UnconditionalBranchReg) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i UnconditionalBranchReg) Op3() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i UnconditionalBranchReg) Op2() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i UnconditionalBranchReg) Opc() uint32 {
	return ExtractBits(uint32(i), 21, 4)
}
func (i UnconditionalBranchReg) Opcode() uint32 {
	return ExtractBits(uint32(i), 25, 7)
}

//--------------------------------------------------------
// C4.3 Loads and stores
//--------------------------------------------------------

type LdstTags uint32

func (i LdstTags) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstTags) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstTags) Op2() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstTags) Imm9() uint32 {
	return ExtractBits(uint32(i), 12, 9)
}
func (i LdstTags) Anon0() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstTags) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstTags) Anon1() uint32 {
	return ExtractBits(uint32(i), 24, 8)
}

type LdstExclusive uint32

func (i LdstExclusive) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstExclusive) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstExclusive) Rt2() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i LdstExclusive) O0() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i LdstExclusive) Rs() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i LdstExclusive) O1() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstExclusive) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i LdstExclusive) O2() uint32 {
	return ExtractBits(uint32(i), 23, 1)
}
func (i LdstExclusive) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 6)
}
func (i LdstExclusive) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LoadRegisterLiteral uint32

func (i LoadRegisterLiteral) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LoadRegisterLiteral) Imm() int32 {
	return int32(ExtractBits(uint32(i), 5, 19))
}
func (i LoadRegisterLiteral) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LoadRegisterLiteral) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LoadRegisterLiteral) Group2() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LoadRegisterLiteral) Opc() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstNoAllocPair uint32

func (i LdstNoAllocPair) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstNoAllocPair) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstNoAllocPair) Rt2() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i LdstNoAllocPair) Imm() int32 {
	return int32(ExtractBits(uint32(i), 15, 7))
}
func (i LdstNoAllocPair) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i LdstNoAllocPair) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 3)
}
func (i LdstNoAllocPair) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstNoAllocPair) Group2() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstNoAllocPair) Opc() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegPairPostIdx uint32

func (i LdstRegPairPostIdx) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegPairPostIdx) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegPairPostIdx) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstRegPairPostIdx) Imm() int32 {
	return int32(ExtractBits(uint32(i), 12, 9))
}
func (i LdstRegPairPostIdx) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegPairPostIdx) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegPairPostIdx) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegPairPostIdx) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegPairPostIdx) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegPairPostIdx) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegPairOffset uint32

func (i LdstRegPairOffset) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegPairOffset) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegPairOffset) Rt2() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i LdstRegPairOffset) Imm() int32 {
	return int32(ExtractBits(uint32(i), 15, 7))
}
func (i LdstRegPairOffset) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i LdstRegPairOffset) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 3)
}
func (i LdstRegPairOffset) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegPairOffset) Group2() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegPairOffset) Opc() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegPairPreIdx uint32

func (i LdstRegPairPreIdx) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegPairPreIdx) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegPairPreIdx) Rt2() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i LdstRegPairPreIdx) Imm() uint32 {
	return ExtractBits(uint32(i), 15, 7)
}
func (i LdstRegPairPreIdx) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i LdstRegPairPreIdx) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 3)
}
func (i LdstRegPairPreIdx) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegPairPreIdx) Group2() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegPairPreIdx) Opc() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegUnscaledImm uint32

func (i LdstRegUnscaledImm) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegUnscaledImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegUnscaledImm) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstRegUnscaledImm) Imm() int32 {
	return int32(ExtractBits(uint32(i), 12, 9))
}
func (i LdstRegUnscaledImm) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegUnscaledImm) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegUnscaledImm) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegUnscaledImm) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegUnscaledImm) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegUnscaledImm) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegImmPostIdx uint32

func (i LdstRegImmPostIdx) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegImmPostIdx) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegImmPostIdx) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstRegImmPostIdx) Imm() uint32 {
	return ExtractBits(uint32(i), 12, 9)
}
func (i LdstRegImmPostIdx) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegImmPostIdx) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegImmPostIdx) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegImmPostIdx) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegImmPostIdx) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegImmPostIdx) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegisterUnpriv uint32

func (i LdstRegisterUnpriv) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegisterUnpriv) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegisterUnpriv) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstRegisterUnpriv) Imm() int32 {
	return int32(ExtractBits(uint32(i), 12, 9))
}
func (i LdstRegisterUnpriv) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegisterUnpriv) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegisterUnpriv) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegisterUnpriv) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegisterUnpriv) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegisterUnpriv) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegImmPreIdx uint32

func (i LdstRegImmPreIdx) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegImmPreIdx) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegImmPreIdx) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstRegImmPreIdx) Imm() int32 {
	return int32(ExtractBits(uint32(i), 12, 9))
}
func (i LdstRegImmPreIdx) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegImmPreIdx) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegImmPreIdx) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegImmPreIdx) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegImmPreIdx) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegImmPreIdx) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegRegOffset uint32

func (i LdstRegRegOffset) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegRegOffset) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegRegOffset) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i LdstRegRegOffset) S() uint32 {
	return ExtractBits(uint32(i), 12, 1)
}
func (i LdstRegRegOffset) Option() uint32 {
	return ExtractBits(uint32(i), 13, 3)
}
func (i LdstRegRegOffset) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i LdstRegRegOffset) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegRegOffset) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegRegOffset) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegRegOffset) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegRegOffset) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegRegOffset) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegUnsignedImm uint32

func (i LdstRegUnsignedImm) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegUnsignedImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegUnsignedImm) Imm() uint32 {
	return ExtractBits(uint32(i), 10, 12)
}
func (i LdstRegUnsignedImm) Opc() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LdstRegUnsignedImm) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegUnsignedImm) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegUnsignedImm) Group2() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegUnsignedImm) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type LdstRegImmPac uint32

func (i LdstRegImmPac) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LdstRegImmPac) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LdstRegImmPac) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i LdstRegImmPac) W() uint32 {
	return ExtractBits(uint32(i), 11, 1)
}
func (i LdstRegImmPac) Imm() uint32 {
	return ExtractBits(uint32(i), 12, 9)
}
func (i LdstRegImmPac) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LdstRegImmPac) S() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i LdstRegImmPac) M() uint32 {
	return ExtractBits(uint32(i), 23, 1)
}
func (i LdstRegImmPac) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 2)
}
func (i LdstRegImmPac) V() uint32 {
	return ExtractBits(uint32(i), 26, 1)
}
func (i LdstRegImmPac) Group4() uint32 {
	return ExtractBits(uint32(i), 27, 3)
}
func (i LdstRegImmPac) Size() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdLdstMult uint32

func (i SimdLdstMult) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdLdstMult) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdLdstMult) Size() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdLdstMult) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i SimdLdstMult) Group1() uint32 {
	return ExtractBits(uint32(i), 16, 6)
}
func (i SimdLdstMult) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i SimdLdstMult) Group2() uint32 {
	return ExtractBits(uint32(i), 23, 7)
}
func (i SimdLdstMult) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdLdstMult) Group3() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdLdstMultPi uint32

func (i SimdLdstMultPi) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdLdstMultPi) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdLdstMultPi) Size() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdLdstMultPi) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i SimdLdstMultPi) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdLdstMultPi) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdLdstMultPi) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i SimdLdstMultPi) Group2() uint32 {
	return ExtractBits(uint32(i), 23, 7)
}
func (i SimdLdstMultPi) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdLdstMultPi) Group3() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdLdstSingle uint32

func (i SimdLdstSingle) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdLdstSingle) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdLdstSingle) Size() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdLdstSingle) S() uint32 {
	return ExtractBits(uint32(i), 12, 1)
}
func (i SimdLdstSingle) Opcode() uint32 {
	return ExtractBits(uint32(i), 13, 3)
}
func (i SimdLdstSingle) Group1() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdLdstSingle) R() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdLdstSingle) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i SimdLdstSingle) Group2() uint32 {
	return ExtractBits(uint32(i), 23, 7)
}
func (i SimdLdstSingle) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdLdstSingle) Group3() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdLdstSinglePi uint32

func (i SimdLdstSinglePi) Rt() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdLdstSinglePi) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdLdstSinglePi) Size() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdLdstSinglePi) S() uint32 {
	return ExtractBits(uint32(i), 12, 1)
}
func (i SimdLdstSinglePi) Opcode() uint32 {
	return ExtractBits(uint32(i), 13, 3)
}
func (i SimdLdstSinglePi) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdLdstSinglePi) R() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdLdstSinglePi) L() uint32 {
	return ExtractBits(uint32(i), 22, 1)
}
func (i SimdLdstSinglePi) Group1() uint32 {
	return ExtractBits(uint32(i), 23, 7)
}
func (i SimdLdstSinglePi) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdLdstSinglePi) Group2() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

//--------------------------------------------------------
// C4.5 Data processing - register
//--------------------------------------------------------

type LogicalShiftedReg uint32

func (i LogicalShiftedReg) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i LogicalShiftedReg) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i LogicalShiftedReg) Imm() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i LogicalShiftedReg) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i LogicalShiftedReg) N() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i LogicalShiftedReg) Shift() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i LogicalShiftedReg) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i LogicalShiftedReg) Opc() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i LogicalShiftedReg) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type AddSubShiftedReg uint32

func (i AddSubShiftedReg) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i AddSubShiftedReg) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i AddSubShiftedReg) Imm() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i AddSubShiftedReg) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i AddSubShiftedReg) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i AddSubShiftedReg) Shift() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i AddSubShiftedReg) Group2() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i AddSubShiftedReg) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i AddSubShiftedReg) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i AddSubShiftedReg) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type AddSubExtendedReg uint32

func (i AddSubExtendedReg) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i AddSubExtendedReg) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i AddSubExtendedReg) Imm() uint32 {
	return ExtractBits(uint32(i), 10, 3)
}
func (i AddSubExtendedReg) Option() uint32 {
	return ExtractBits(uint32(i), 13, 3)
}
func (i AddSubExtendedReg) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i AddSubExtendedReg) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i AddSubExtendedReg) Opt() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i AddSubExtendedReg) Group2() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i AddSubExtendedReg) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i AddSubExtendedReg) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i AddSubExtendedReg) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type AddSubWithCarry uint32

func (i AddSubWithCarry) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i AddSubWithCarry) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i AddSubWithCarry) Opcode2() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i AddSubWithCarry) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i AddSubWithCarry) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i AddSubWithCarry) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i AddSubWithCarry) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i AddSubWithCarry) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type ConditionalCompareReg uint32

func (i ConditionalCompareReg) Nzcv() uint32 {
	return ExtractBits(uint32(i), 0, 4)
}
func (i ConditionalCompareReg) O3() uint32 {
	return ExtractBits(uint32(i), 4, 1)
}
func (i ConditionalCompareReg) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i ConditionalCompareReg) O2() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i ConditionalCompareReg) Group1() uint32 {
	return ExtractBits(uint32(i), 11, 1)
}
func (i ConditionalCompareReg) Cond() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i ConditionalCompareReg) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i ConditionalCompareReg) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i ConditionalCompareReg) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i ConditionalCompareReg) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i ConditionalCompareReg) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type ConditionalCompareImm uint32

func (i ConditionalCompareImm) Nzcv() uint32 {
	return ExtractBits(uint32(i), 0, 4)
}
func (i ConditionalCompareImm) O3() uint32 {
	return ExtractBits(uint32(i), 4, 1)
}
func (i ConditionalCompareImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i ConditionalCompareImm) O2() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i ConditionalCompareImm) Group1() uint32 {
	return ExtractBits(uint32(i), 11, 1)
}
func (i ConditionalCompareImm) Cond() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i ConditionalCompareImm) Imm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i ConditionalCompareImm) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i ConditionalCompareImm) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i ConditionalCompareImm) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i ConditionalCompareImm) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type ConditionalSelect uint32

func (i ConditionalSelect) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i ConditionalSelect) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i ConditionalSelect) Op2() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i ConditionalSelect) Cond() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i ConditionalSelect) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i ConditionalSelect) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i ConditionalSelect) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i ConditionalSelect) Op() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i ConditionalSelect) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type DataProcessing3 uint32

func (i DataProcessing3) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i DataProcessing3) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i DataProcessing3) Ra() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i DataProcessing3) O0() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i DataProcessing3) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i DataProcessing3) Op31() uint32 {
	return ExtractBits(uint32(i), 21, 3)
}
func (i DataProcessing3) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i DataProcessing3) Op54() uint32 {
	return ExtractBits(uint32(i), 29, 2)
}
func (i DataProcessing3) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type DataProcessing2 uint32

func (i DataProcessing2) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i DataProcessing2) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i DataProcessing2) Opcode() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i DataProcessing2) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i DataProcessing2) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i DataProcessing2) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i DataProcessing2) Group2() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i DataProcessing2) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type DataProcessing1 uint32

func (i DataProcessing1) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i DataProcessing1) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i DataProcessing1) Opcode() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i DataProcessing1) Opcode2() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i DataProcessing1) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i DataProcessing1) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i DataProcessing1) Group2() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i DataProcessing1) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

//--------------------------------------------------------
// C4.6 - Data Processing -SIMD and floating point
//--------------------------------------------------------

type FloatingFixedConversion uint32

func (i FloatingFixedConversion) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingFixedConversion) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingFixedConversion) Scale() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i FloatingFixedConversion) Opcode() uint32 {
	return ExtractBits(uint32(i), 16, 3)
}
func (i FloatingFixedConversion) Mode() uint32 {
	return ExtractBits(uint32(i), 19, 2)
}
func (i FloatingFixedConversion) Group1() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingFixedConversion) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingFixedConversion) Group2() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingFixedConversion) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingFixedConversion) Group3() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingFixedConversion) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingConditionalCompare uint32

func (i FloatingConditionalCompare) Nzvb() uint32 {
	return ExtractBits(uint32(i), 0, 4)
}
func (i FloatingConditionalCompare) Op() uint32 {
	return ExtractBits(uint32(i), 4, 1)
}
func (i FloatingConditionalCompare) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingConditionalCompare) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i FloatingConditionalCompare) Cond() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i FloatingConditionalCompare) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i FloatingConditionalCompare) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingConditionalCompare) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingConditionalCompare) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingConditionalCompare) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingConditionalCompare) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingConditionalCompare) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingDataProcessing2 uint32

func (i FloatingDataProcessing2) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingDataProcessing2) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingDataProcessing2) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i FloatingDataProcessing2) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i FloatingDataProcessing2) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i FloatingDataProcessing2) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingDataProcessing2) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingDataProcessing2) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingDataProcessing2) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingDataProcessing2) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingDataProcessing2) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingConditionalSelect uint32

func (i FloatingConditionalSelect) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingConditionalSelect) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingConditionalSelect) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i FloatingConditionalSelect) Cond() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i FloatingConditionalSelect) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i FloatingConditionalSelect) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingConditionalSelect) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingConditionalSelect) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingConditionalSelect) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingConditionalSelect) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingConditionalSelect) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingImm uint32

func (i FloatingImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingImm) Imm5() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingImm) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 3)
}
func (i FloatingImm) Imm8() uint32 {
	return ExtractBits(uint32(i), 13, 8)
}
func (i FloatingImm) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingImm) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingImm) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingImm) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingImm) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingImm) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingCompare uint32

func (i FloatingCompare) Opcode2() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingCompare) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingCompare) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 4)
}
func (i FloatingCompare) Op() uint32 {
	return ExtractBits(uint32(i), 14, 2)
}
func (i FloatingCompare) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i FloatingCompare) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingCompare) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingCompare) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingCompare) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingCompare) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingCompare) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingDataProcessing1 uint32

func (i FloatingDataProcessing1) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingDataProcessing1) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingDataProcessing1) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i FloatingDataProcessing1) Opcode() uint32 {
	return ExtractBits(uint32(i), 15, 6)
}
func (i FloatingDataProcessing1) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingDataProcessing1) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingDataProcessing1) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingDataProcessing1) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingDataProcessing1) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingDataProcessing1) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingIntegerConversion uint32

func (i FloatingIntegerConversion) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingIntegerConversion) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingIntegerConversion) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 6)
}
func (i FloatingIntegerConversion) Opcode() uint32 {
	return ExtractBits(uint32(i), 16, 3)
}
func (i FloatingIntegerConversion) Rmode() uint32 {
	return ExtractBits(uint32(i), 19, 2)
}
func (i FloatingIntegerConversion) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingIntegerConversion) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingIntegerConversion) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingIntegerConversion) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingIntegerConversion) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingIntegerConversion) Sf() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type FloatingDataProcessing3 uint32

func (i FloatingDataProcessing3) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i FloatingDataProcessing3) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i FloatingDataProcessing3) Ra() uint32 {
	return ExtractBits(uint32(i), 10, 5)
}
func (i FloatingDataProcessing3) O0() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i FloatingDataProcessing3) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i FloatingDataProcessing3) O1() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i FloatingDataProcessing3) Type() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i FloatingDataProcessing3) Group1() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i FloatingDataProcessing3) S() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i FloatingDataProcessing3) Group2() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i FloatingDataProcessing3) M() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type Simd3Same uint32

func (i Simd3Same) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Simd3Same) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Simd3Same) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i Simd3Same) Opcode() uint32 {
	return ExtractBits(uint32(i), 11, 5)
}
func (i Simd3Same) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i Simd3Same) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i Simd3Same) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i Simd3Same) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i Simd3Same) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i Simd3Same) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i Simd3Same) Group4() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type Simd3Different uint32

func (i Simd3Different) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Simd3Different) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Simd3Different) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i Simd3Different) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i Simd3Different) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i Simd3Different) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i Simd3Different) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i Simd3Different) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i Simd3Different) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i Simd3Different) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i Simd3Different) Group4() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type Simd2RegMisc uint32

func (i Simd2RegMisc) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Simd2RegMisc) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Simd2RegMisc) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i Simd2RegMisc) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 5)
}
func (i Simd2RegMisc) Group2() uint32 {
	return ExtractBits(uint32(i), 17, 5)
}
func (i Simd2RegMisc) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i Simd2RegMisc) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i Simd2RegMisc) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i Simd2RegMisc) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i Simd2RegMisc) Group4() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdAcrossLanes uint32

func (i SimdAcrossLanes) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdAcrossLanes) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdAcrossLanes) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdAcrossLanes) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 5)
}
func (i SimdAcrossLanes) Group2() uint32 {
	return ExtractBits(uint32(i), 17, 5)
}
func (i SimdAcrossLanes) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdAcrossLanes) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdAcrossLanes) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdAcrossLanes) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdAcrossLanes) Group4() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdCopy uint32

func (i SimdCopy) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdCopy) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdCopy) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdCopy) Imm4() uint32 {
	return ExtractBits(uint32(i), 11, 4)
}
func (i SimdCopy) Group2() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i SimdCopy) Imm5() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdCopy) Group3() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i SimdCopy) Op() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdCopy) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdCopy) Group4() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdVectorXIndexedElement uint32

func (i SimdVectorXIndexedElement) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdVectorXIndexedElement) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdVectorXIndexedElement) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdVectorXIndexedElement) H() uint32 {
	return ExtractBits(uint32(i), 11, 1)
}
func (i SimdVectorXIndexedElement) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i SimdVectorXIndexedElement) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 4)
}
func (i SimdVectorXIndexedElement) M() uint32 {
	return ExtractBits(uint32(i), 20, 1)
}
func (i SimdVectorXIndexedElement) L() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdVectorXIndexedElement) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdVectorXIndexedElement) Group2() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdVectorXIndexedElement) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdVectorXIndexedElement) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdVectorXIndexedElement) Group3() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdModifiedImm uint32

func (i SimdModifiedImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdModifiedImm) H() uint32 {
	return ExtractBits(uint32(i), 5, 1)
}
func (i SimdModifiedImm) G() uint32 {
	return ExtractBits(uint32(i), 6, 1)
}
func (i SimdModifiedImm) F() uint32 {
	return ExtractBits(uint32(i), 7, 1)
}
func (i SimdModifiedImm) E() uint32 {
	return ExtractBits(uint32(i), 8, 1)
}
func (i SimdModifiedImm) D() uint32 {
	return ExtractBits(uint32(i), 9, 1)
}
func (i SimdModifiedImm) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdModifiedImm) O2() uint32 {
	return ExtractBits(uint32(i), 11, 1)
}
func (i SimdModifiedImm) Cmode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i SimdModifiedImm) C() uint32 {
	return ExtractBits(uint32(i), 16, 1)
}
func (i SimdModifiedImm) B() uint32 {
	return ExtractBits(uint32(i), 17, 1)
}
func (i SimdModifiedImm) A() uint32 {
	return ExtractBits(uint32(i), 18, 1)
}
func (i SimdModifiedImm) Group2() uint32 {
	return ExtractBits(uint32(i), 19, 10)
}
func (i SimdModifiedImm) Op() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdModifiedImm) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdModifiedImm) Group3() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdShiftByImm uint32

func (i SimdShiftByImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdShiftByImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdShiftByImm) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdShiftByImm) Opcode() uint32 {
	return ExtractBits(uint32(i), 11, 5)
}
func (i SimdShiftByImm) Immb() uint32 {
	return ExtractBits(uint32(i), 16, 3)
}
func (i SimdShiftByImm) Immh() uint32 {
	return ExtractBits(uint32(i), 19, 4)
}
func (i SimdShiftByImm) Group2() uint32 {
	return ExtractBits(uint32(i), 23, 6)
}
func (i SimdShiftByImm) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdShiftByImm) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdShiftByImm) Group3() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdTableLookup uint32

func (i SimdTableLookup) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdTableLookup) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdTableLookup) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdTableLookup) Op() uint32 {
	return ExtractBits(uint32(i), 12, 1)
}
func (i SimdTableLookup) Len() uint32 {
	return ExtractBits(uint32(i), 13, 2)
}
func (i SimdTableLookup) Group2() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i SimdTableLookup) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdTableLookup) Group3() uint32 {
	return ExtractBits(uint32(i), 21, 9)
}
func (i SimdTableLookup) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdTableLookup) Group4() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdPermute uint32

func (i SimdPermute) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdPermute) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdPermute) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdPermute) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 3)
}
func (i SimdPermute) Group2() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i SimdPermute) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdPermute) Group3() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdPermute) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdPermute) Group4() uint32 {
	return ExtractBits(uint32(i), 24, 6)
}
func (i SimdPermute) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdPermute) Group5() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdExtract uint32

func (i SimdExtract) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdExtract) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdExtract) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdExtract) Imm() uint32 {
	return ExtractBits(uint32(i), 11, 4)
}
func (i SimdExtract) Group2() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i SimdExtract) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdExtract) Group3() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdExtract) Op2() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdExtract) Group4() uint32 {
	return ExtractBits(uint32(i), 24, 6)
}
func (i SimdExtract) Q() uint32 {
	return ExtractBits(uint32(i), 30, 1)
}
func (i SimdExtract) Group5() uint32 {
	return ExtractBits(uint32(i), 31, 1)
}

type SimdScalar3Same uint32

func (i SimdScalar3Same) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalar3Same) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalar3Same) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdScalar3Same) Opcode() uint32 {
	return ExtractBits(uint32(i), 11, 5)
}
func (i SimdScalar3Same) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdScalar3Same) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdScalar3Same) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdScalar3Same) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdScalar3Same) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalar3Same) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdScalar3Different uint32

func (i SimdScalar3Different) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalar3Different) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalar3Different) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdScalar3Different) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i SimdScalar3Different) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdScalar3Different) Group2() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdScalar3Different) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdScalar3Different) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdScalar3Different) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalar3Different) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdScalar2RegisterMisc uint32

func (i SimdScalar2RegisterMisc) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalar2RegisterMisc) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalar2RegisterMisc) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdScalar2RegisterMisc) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 5)
}
func (i SimdScalar2RegisterMisc) Group2() uint32 {
	return ExtractBits(uint32(i), 17, 5)
}
func (i SimdScalar2RegisterMisc) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdScalar2RegisterMisc) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdScalar2RegisterMisc) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalar2RegisterMisc) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdScalarPairwise uint32

func (i SimdScalarPairwise) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalarPairwise) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalarPairwise) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i SimdScalarPairwise) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 5)
}
func (i SimdScalarPairwise) Group2() uint32 {
	return ExtractBits(uint32(i), 17, 5)
}
func (i SimdScalarPairwise) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdScalarPairwise) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdScalarPairwise) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalarPairwise) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdScalarCopy uint32

func (i SimdScalarCopy) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalarCopy) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalarCopy) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdScalarCopy) Imm4() uint32 {
	return ExtractBits(uint32(i), 11, 4)
}
func (i SimdScalarCopy) Group2() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i SimdScalarCopy) Imm5() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i SimdScalarCopy) Group3() uint32 {
	return ExtractBits(uint32(i), 21, 8)
}
func (i SimdScalarCopy) Op() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalarCopy) Group4() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdScalarXIndexedElement uint32

func (i SimdScalarXIndexedElement) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalarXIndexedElement) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalarXIndexedElement) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdScalarXIndexedElement) H() uint32 {
	return ExtractBits(uint32(i), 11, 1)
}
func (i SimdScalarXIndexedElement) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 4)
}
func (i SimdScalarXIndexedElement) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 4)
}
func (i SimdScalarXIndexedElement) M() uint32 {
	return ExtractBits(uint32(i), 20, 1)
}
func (i SimdScalarXIndexedElement) L() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i SimdScalarXIndexedElement) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i SimdScalarXIndexedElement) Group2() uint32 {
	return ExtractBits(uint32(i), 24, 5)
}
func (i SimdScalarXIndexedElement) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalarXIndexedElement) Group3() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type SimdScalarShiftByImm uint32

func (i SimdScalarShiftByImm) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i SimdScalarShiftByImm) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i SimdScalarShiftByImm) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 1)
}
func (i SimdScalarShiftByImm) Opcode() uint32 {
	return ExtractBits(uint32(i), 11, 5)
}
func (i SimdScalarShiftByImm) Immb() uint32 {
	return ExtractBits(uint32(i), 16, 3)
}
func (i SimdScalarShiftByImm) Immh() uint32 {
	return ExtractBits(uint32(i), 19, 4)
}
func (i SimdScalarShiftByImm) Group2() uint32 {
	return ExtractBits(uint32(i), 23, 6)
}
func (i SimdScalarShiftByImm) U() uint32 {
	return ExtractBits(uint32(i), 29, 1)
}
func (i SimdScalarShiftByImm) Group3() uint32 {
	return ExtractBits(uint32(i), 30, 2)
}

type CryptographicAes uint32

func (i CryptographicAes) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i CryptographicAes) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i CryptographicAes) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i CryptographicAes) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 5)
}
func (i CryptographicAes) Group2() uint32 {
	return ExtractBits(uint32(i), 17, 5)
}
func (i CryptographicAes) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i CryptographicAes) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 8)
}

type Cryptographic3RegSha uint32

func (i Cryptographic3RegSha) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Cryptographic3RegSha) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Cryptographic3RegSha) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i Cryptographic3RegSha) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 3)
}
func (i Cryptographic3RegSha) Group2() uint32 {
	return ExtractBits(uint32(i), 15, 1)
}
func (i Cryptographic3RegSha) Rm() uint32 {
	return ExtractBits(uint32(i), 16, 5)
}
func (i Cryptographic3RegSha) Group3() uint32 {
	return ExtractBits(uint32(i), 21, 1)
}
func (i Cryptographic3RegSha) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i Cryptographic3RegSha) Group4() uint32 {
	return ExtractBits(uint32(i), 24, 8)
}

type Cryptographic2RegSha uint32

func (i Cryptographic2RegSha) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i Cryptographic2RegSha) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i Cryptographic2RegSha) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 2)
}
func (i Cryptographic2RegSha) Opcode() uint32 {
	return ExtractBits(uint32(i), 12, 5)
}
func (i Cryptographic2RegSha) Group2() uint32 {
	return ExtractBits(uint32(i), 17, 5)
}
func (i Cryptographic2RegSha) Size() uint32 {
	return ExtractBits(uint32(i), 22, 2)
}
func (i Cryptographic2RegSha) Group3() uint32 {
	return ExtractBits(uint32(i), 24, 8)
}

type PointerAuth uint32

func (i PointerAuth) Rd() uint32 {
	return ExtractBits(uint32(i), 0, 5)
}
func (i PointerAuth) Rn() uint32 {
	return ExtractBits(uint32(i), 5, 5)
}
func (i PointerAuth) Group1() uint32 {
	return ExtractBits(uint32(i), 10, 3)
}
func (i PointerAuth) Z() uint32 {
	return ExtractBits(uint32(i), 13, 1)
}
func (i PointerAuth) Group2() uint32 {
	return ExtractBits(uint32(i), 14, 18)
}

type SystemReg uint32

const (
	SYSREG_NONE SystemReg = iota
	REG_ACTLR_EL1
	REG_ACTLR_EL2
	REG_ACTLR_EL3
	REG_AFSR0_EL1
	REG_AFSR1_EL2
	REG_AFSR0_EL2
	REG_AFSR0_EL3
	REG_AFSR1_EL1
	REG_AFSR1_EL3
	REG_AIDR_EL1
	REG_ALLE1
	REG_ALLE1IS
	REG_ALLE2
	REG_ALLE2IS
	REG_ALLE3
	REG_ALLE3IS
	REG_AMAIR_EL1
	REG_AMAIR_EL2
	REG_AMAIR_EL3
	REG_ASIDE1
	REG_ASIDE1IS
	REG_CCSIDR_EL1
	REG_CISW
	REG_CIVAC
	REG_CLIDR_EL1
	REG_CNTFRQ_EL0
	REG_CNTHCTL_EL2
	REG_CNTHP_CTL_EL2
	REG_CNTHP_CVAL_EL2
	REG_CNTHP_TVAL_EL2
	REG_CNTKCTL_EL1
	REG_CNTPCT_EL0
	REG_CNTPS_CTL_EL1
	REG_CNTPS_CVAL_EL1
	REG_CNTPS_TVAL_EL1
	REG_CNTP_CTL_EL0
	REG_CNTP_CVAL_EL0
	REG_CNTP_TVAL_EL0
	REG_CNTVCT_EL0
	REG_CNTV_CTL_EL0
	REG_CNTV_CVAL_EL0
	REG_CNTV_TVAL_EL0
	REG_CONTEXTIDR_EL1
	REG_CPACR_EL1
	REG_CPTR_EL2
	REG_CPTR_EL3
	REG_CSSELR_EL1
	REG_CSW
	REG_CTR_EL0
	REG_CVAC
	REG_CVAU
	REG_DACR32_EL2
	REG_DAIFCLR
	REG_DAIFSET
	REG_DBGAUTHSTATUS_EL1
	REG_DBGCLAIMCLR_EL1
	REG_DBGCLAIMSET_EL1
	REG_DBGBCR0_EL1
	REG_DBGBCR10_EL1
	REG_DBGBCR11_EL1
	REG_DBGBCR12_EL1
	REG_DBGBCR13_EL1
	REG_DBGBCR14_EL1
	REG_DBGBCR15_EL1
	REG_DBGBCR1_EL1
	REG_DBGBCR2_EL1
	REG_DBGBCR3_EL1
	REG_DBGBCR4_EL1
	REG_DBGBCR5_EL1
	REG_DBGBCR6_EL1
	REG_DBGBCR7_EL1
	REG_DBGBCR8_EL1
	REG_DBGBCR9_EL1
	REG_DBGDTRRX_EL0
	REG_DBGDTRTX_EL0
	REG_DBGDTR_EL0
	REG_DBGPRCR_EL1
	REG_DBGVCR32_EL2
	REG_DBGBVR0_EL1
	REG_DBGBVR10_EL1
	REG_DBGBVR11_EL1
	REG_DBGBVR12_EL1
	REG_DBGBVR13_EL1
	REG_DBGBVR14_EL1
	REG_DBGBVR15_EL1
	REG_DBGBVR1_EL1
	REG_DBGBVR2_EL1
	REG_DBGBVR3_EL1
	REG_DBGBVR4_EL1
	REG_DBGBVR5_EL1
	REG_DBGBVR6_EL1
	REG_DBGBVR7_EL1
	REG_DBGBVR8_EL1
	REG_DBGBVR9_EL1
	REG_DBGWCR0_EL1
	REG_DBGWCR10_EL1
	REG_DBGWCR11_EL1
	REG_DBGWCR12_EL1
	REG_DBGWCR13_EL1
	REG_DBGWCR14_EL1
	REG_DBGWCR15_EL1
	REG_DBGWCR1_EL1
	REG_DBGWCR2_EL1
	REG_DBGWCR3_EL1
	REG_DBGWCR4_EL1
	REG_DBGWCR5_EL1
	REG_DBGWCR6_EL1
	REG_DBGWCR7_EL1
	REG_DBGWCR8_EL1
	REG_DBGWCR9_EL1
	REG_DBGWVR0_EL1
	REG_DBGWVR10_EL1
	REG_DBGWVR11_EL1
	REG_DBGWVR12_EL1
	REG_DBGWVR13_EL1
	REG_DBGWVR14_EL1
	REG_DBGWVR15_EL1
	REG_DBGWVR1_EL1
	REG_DBGWVR2_EL1
	REG_DBGWVR3_EL1
	REG_DBGWVR4_EL1
	REG_DBGWVR5_EL1
	REG_DBGWVR6_EL1
	REG_DBGWVR7_EL1
	REG_DBGWVR8_EL1
	REG_DBGWVR9_EL1
	REG_DCZID_EL0
	REG_EL1
	REG_ESR_EL1
	REG_ESR_EL2
	REG_ESR_EL3
	REG_FAR_EL1
	REG_FAR_EL2
	REG_FAR_EL3
	REG_HACR_EL2
	REG_HCR_EL2
	REG_HPFAR_EL2
	REG_HSTR_EL2
	REG_IALLU
	REG_IVAU
	REG_IALLUIS
	REG_ID_AA64DFR0_EL1
	REG_ID_AA64ISAR0_EL1
	REG_ID_AA64ISAR1_EL1
	REG_ID_AA64MMFR0_EL1
	REG_ID_AA64MMFR1_EL1
	REG_ID_AA64PFR0_EL1
	REG_ID_AA64PFR1_EL1
	REG_IPAS2E1IS
	REG_IPAS2LE1IS
	REG_IPAS2E1
	REG_IPAS2LE1
	REG_ISW
	REG_IVAC
	REG_MAIR_EL1
	REG_MAIR_EL2
	REG_MAIR_EL3
	REG_MDCCINT_EL1
	REG_MDCCSR_EL0
	REG_MDCR_EL2
	REG_MDCR_EL3
	REG_MDRAR_EL1
	REG_MDSCR_EL1
	REG_MVFR0_EL1
	REG_MVFR1_EL1
	REG_MVFR2_EL1
	REG_OSDTRRX_EL1
	REG_OSDTRTX_EL1
	REG_OSECCR_EL1
	REG_OSLAR_EL1
	REG_OSDLR_EL1
	REG_OSLSR_EL1
	REG_PAN
	REG_PAR_EL1
	REG_PMCCNTR_EL0
	REG_PMCEID0_EL0
	REG_PMCEID1_EL0
	REG_PMCNTENSET_EL0
	REG_PMCR_EL0
	REG_PMCNTENCLR_EL0
	REG_PMINTENCLR_EL1
	REG_PMINTENSET_EL1
	REG_PMOVSCLR_EL0
	REG_PMOVSSET_EL0
	REG_PMSELR_EL0
	REG_PMSWINC_EL0
	REG_PMUSERENR_EL0
	REG_PMXEVCNTR_EL0
	REG_PMXEVTYPER_EL0
	REG_RMR_EL1
	REG_RMR_EL2
	REG_RMR_EL3
	REG_RVBAR_EL1
	REG_RVBAR_EL2
	REG_RVBAR_EL3
	REG_S12E0R
	REG_S12E0W
	REG_S12E1R
	REG_S12E1W
	REG_S1E0R
	REG_S1E0W
	REG_S1E1R
	REG_S1E1W
	REG_S1E2R
	REG_S1E2W
	REG_S1E3R
	REG_S1E3W
	REG_SCR_EL3
	REG_SDER32_EL3
	REG_SCTLR_EL1
	REG_SCTLR_EL2
	REG_SCTLR_EL3
	REG_SPSEL
	REG_TCR_EL1
	REG_TCR_EL2
	REG_TCR_EL3
	REG_TPIDRRO_EL0
	REG_TPIDR_EL0
	REG_TPIDR_EL1
	REG_TPIDR_EL2
	REG_TPIDR_EL3
	REG_TTBR0_EL1
	REG_TTBR1_EL1
	REG_TTBR0_EL2
	REG_TTBR0_EL3
	REG_VAAE1
	REG_VAAE1IS
	REG_VAALE1
	REG_VAALE1IS
	REG_VAE1
	REG_VAE1IS
	REG_VAE2
	REG_VAE2IS
	REG_VAE3
	REG_VAE3IS
	REG_VALE1
	REG_VALE1IS
	REG_VALE2
	REG_VALE2IS
	REG_VALE3
	REG_VALE3IS
	REG_VBAR_EL1
	REG_VBAR_EL2
	REG_VBAR_EL3
	REG_VMALLE1
	REG_VMALLE1IS
	REG_VMALLS12E1
	REG_VMALLS12E1IS
	REG_VMPIDR_EL0
	REG_VPIDR_EL2
	REG_VTCR_EL2
	REG_VTTBR_EL2
	REG_ZVA
	REG_NUMBER0
	REG_OSHLD
	REG_OSHST
	REG_OSH
	REG_NUMBER4
	REG_NSHLD
	REG_NSHST
	REG_NSH
	REG_NUMBER8
	REG_ISHLD
	REG_ISHST
	REG_ISH
	REG_NUMBER12
	REG_LD
	REG_ST
	REG_SY
	REG_PMEVCNTR0_EL0
	REG_PMEVCNTR1_EL0
	REG_PMEVCNTR2_EL0
	REG_PMEVCNTR3_EL0
	REG_PMEVCNTR4_EL0
	REG_PMEVCNTR5_EL0
	REG_PMEVCNTR6_EL0
	REG_PMEVCNTR7_EL0
	REG_PMEVCNTR8_EL0
	REG_PMEVCNTR9_EL0
	REG_PMEVCNTR10_EL0
	REG_PMEVCNTR11_EL0
	REG_PMEVCNTR12_EL0
	REG_PMEVCNTR13_EL0
	REG_PMEVCNTR14_EL0
	REG_PMEVCNTR15_EL0
	REG_PMEVCNTR16_EL0
	REG_PMEVCNTR17_EL0
	REG_PMEVCNTR18_EL0
	REG_PMEVCNTR19_EL0
	REG_PMEVCNTR20_EL0
	REG_PMEVCNTR21_EL0
	REG_PMEVCNTR22_EL0
	REG_PMEVCNTR23_EL0
	REG_PMEVCNTR24_EL0
	REG_PMEVCNTR25_EL0
	REG_PMEVCNTR26_EL0
	REG_PMEVCNTR27_EL0
	REG_PMEVCNTR28_EL0
	REG_PMEVCNTR29_EL0
	REG_PMEVCNTR30_EL0
	REG_PMEVTYPER0_EL0
	REG_PMEVTYPER1_EL0
	REG_PMEVTYPER2_EL0
	REG_PMEVTYPER3_EL0
	REG_PMEVTYPER4_EL0
	REG_PMEVTYPER5_EL0
	REG_PMEVTYPER6_EL0
	REG_PMEVTYPER7_EL0
	REG_PMEVTYPER8_EL0
	REG_PMEVTYPER9_EL0
	REG_PMEVTYPER10_EL0
	REG_PMEVTYPER11_EL0
	REG_PMEVTYPER12_EL0
	REG_PMEVTYPER13_EL0
	REG_PMEVTYPER14_EL0
	REG_PMEVTYPER15_EL0
	REG_PMEVTYPER16_EL0
	REG_PMEVTYPER17_EL0
	REG_PMEVTYPER18_EL0
	REG_PMEVTYPER19_EL0
	REG_PMEVTYPER20_EL0
	REG_PMEVTYPER21_EL0
	REG_PMEVTYPER22_EL0
	REG_PMEVTYPER23_EL0
	REG_PMEVTYPER24_EL0
	REG_PMEVTYPER25_EL0
	REG_PMEVTYPER26_EL0
	REG_PMEVTYPER27_EL0
	REG_PMEVTYPER28_EL0
	REG_PMEVTYPER29_EL0
	REG_PMEVTYPER30_EL0
	REG_PMCCFILTR_EL0
	REG_C0
	REG_C1
	REG_C2
	REG_C3
	REG_C4
	REG_C5
	REG_C6
	REG_C7
	REG_C8
	REG_C9
	REG_C10
	REG_C11
	REG_C12
	REG_C13
	REG_C14
	REG_C15

	REG_SPSR_EL1
	REG_ELR_EL1
	REG_SP_EL0
	REG_CURRENT_EL
	REG_NZCV
	REG_FPCR
	REG_DSPSR_EL0
	REG_DAIF
	REG_FPSR
	REG_DLR_EL0
	REG_SPSR_EL2
	REG_ELR_EL2
	REG_SP_EL1
	REG_SP_EL2
	REG_SPSR_IRQ
	REG_SPSR_ABT
	REG_SPSR_UND
	REG_SPSR_FIQ
	REG_SPSR_EL3
	REG_ELR_EL3
	REG_IFSR32_EL2
	REG_FPEXC32_EL2
	REG_CNTVOFF_EL2

	REG_MIDR_EL1
	REG_MPIDR_EL1
	REG_REVIDR_EL1
	REG_ID_PFR0_EL1
	REG_ID_PFR1_EL1
	REG_ID_DFR0_EL1
	REG_ID_AFR0_EL1
	REG_ID_MMFR0_EL1
	REG_ID_MMFR1_EL1
	REG_ID_MMFR2_EL1
	REG_ID_MMFR3_EL1
	REG_ID_ISAR0_EL1
	REG_ID_ISAR1_EL1
	REG_ID_ISAR2_EL1
	REG_ID_ISAR3_EL1
	REG_ID_ISAR4_EL1
	REG_ID_ISAR5_EL1
	REG_ID_MMFR4_EL1

	REG_ICC_IAR0_EL1
	REG_ICC_EOIR0_EL1
	REG_ICC_HPPIR0_EL1
	REG_ICC_BPR0_EL1
	REG_ICC_AP0R0_EL1
	REG_ICC_AP0R1_EL1
	REG_ICC_AP0R2_EL1
	REG_ICC_AP0R3_EL1
	REG_ICC_AP1R0_EL1
	REG_ICC_AP1R1_EL1
	REG_ICC_AP1R2_EL1
	REG_ICC_AP1R3_EL1
	REG_ICC_DIR_EL1
	REG_ICC_RPR_EL1
	REG_ICC_IAR1_EL1
	REG_ICC_EOIR1_EL1
	REG_ICC_HPPIR1_EL1
	REG_ICC_BPR1_EL1
	REG_ICC_CTLR_EL1
	REG_ICC_SRE_EL1
	REG_ICC_IGRPEN0_EL1
	REG_ICC_IGRPEN1_EL1

	REG_ICC_ASGI1R_EL2
	REG_ICC_SGI0R_EL2
	REG_ICH_AP0R0_EL2
	REG_ICH_AP0R1_EL2
	REG_ICH_AP0R2_EL2
	REG_ICH_AP0R3_EL2
	REG_ICH_AP1R0_EL2
	REG_ICH_AP1R1_EL2
	REG_ICH_AP1R2_EL2
	REG_ICH_AP1R3_EL2
	REG_ICH_AP1R4_EL2
	REG_ICC_HSRE_EL2
	REG_ICH_HCR_EL2
	REG_ICH_VTR_EL2
	REG_ICH_MISR_EL2
	REG_ICH_EISR_EL2
	REG_ICH_ELRSR_EL2
	REG_ICH_VMCR_EL2

	REG_ICH_LR0_EL2
	REG_ICH_LR1_EL2
	REG_ICH_LR2_EL2
	REG_ICH_LR3_EL2
	REG_ICH_LR4_EL2
	REG_ICH_LR5_EL2
	REG_ICH_LR6_EL2
	REG_ICH_LR7_EL2
	REG_ICH_LR8_EL2
	REG_ICH_LR9_EL2
	REG_ICH_LR10_EL2
	REG_ICH_LR11_EL2
	REG_ICH_LR12_EL2
	REG_ICH_LR13_EL2
	REG_ICH_LR14_EL2
	REG_ICH_LR15_EL2

	REG_ICH_LRC0_EL2
	REG_ICH_LRC1_EL2
	REG_ICH_LRC2_EL2
	REG_ICH_LRC3_EL2
	REG_ICH_LRC4_EL2
	REG_ICH_LRC5_EL2
	REG_ICH_LRC6_EL2
	REG_ICH_LRC7_EL2
	REG_ICH_LRC8_EL2
	REG_ICH_LRC9_EL2
	REG_ICH_LRC10_EL2
	REG_ICH_LRC11_EL2
	REG_ICH_LRC12_EL2
	REG_ICH_LRC13_EL2
	REG_ICH_LRC14_EL2
	REG_ICH_LRC15_EL2

	REG_ICC_MCTLR_EL3
	REG_ICC_MSRE_EL3
	REG_ICC_MGRPEN1_EL3

	REG_TEECR32_EL1
	REG_TEEHBR32_EL1

	REG_ICC_PMR_EL1
	REG_ICC_SGI1R_EL1
	REG_ICC_SGI0R_EL1
	REG_ICC_ASGI1R_EL1
	REG_ICC_SEIEN_EL1
	REG_END_REG
)

func (s SystemReg) String() string {
	return []string{
		"NONE",
		"actlr_el1",
		"actlr_el2",
		"actlr_el3",
		"afsr0_el1",
		"afsr1_el2",
		"afsr0_el2",
		"afsr0_el3",
		"afsr1_el1",
		"afsr1_el3",
		"aidr_el1",
		"alle1",
		"alle1is",
		"alle2",
		"alle2is",
		"alle3",
		"alle3is",
		"amair_el1",
		"amair_el2",
		"amair_el3",
		"aside1",
		"aside1is",
		"ccsidr_el1",
		"cisw",
		"civac",
		"clidr_el1",
		"cntfrq_el0",
		"cnthctl_el2",
		"cnthp_ctl_el2",
		"cnthp_cval_el2",
		"cnthp_tval_el2",
		"cntkctl_el1",
		"cntpct_el0",
		"cntps_ctl_el1",
		"cntps_cval_el1",
		"cntps_tval_el1",
		"cntp_ctl_el0",
		"cntp_cval_el0",
		"cntp_tval_el0",
		"cntvct_el0",
		"cntv_ctl_el0",
		"cntv_cval_el0",
		"cntv_tval_el0",
		"contextidr_el1",
		"cpacr_el1",
		"cptr_el2",
		"cptr_el3",
		"csselr_el1",
		"csw",
		"ctr_el0",
		"cvac",
		"cvau",
		"dacr32_el2",
		"daifclr",
		"daifset",
		"dbgauthstatus_el1",
		"dbgclaimclr_el1",
		"dbgclaimset_el1",
		"dbgbcr0_el1",
		"dbgbcr10_el1",
		"dbgbcr11_el1",
		"dbgbcr12_el1",
		"dbgbcr13_el1",
		"dbgbcr14_el1",
		"dbgbcr15_el1",
		"dbgbcr1_el1",
		"dbgbcr2_el1",
		"dbgbcr3_el1",
		"dbgbcr4_el1",
		"dbgbcr5_el1",
		"dbgbcr6_el1",
		"dbgbcr7_el1",
		"dbgbcr8_el1",
		"dbgbcr9_el1",
		"dbgdtrrx_el0",
		"dbgdtrtx_el0",
		"dbgdtr_el0",
		"dbgprcr_el1",
		"dbgvcr32_el2",
		"dbgbvr0_el1",
		"dbgbvr10_el1",
		"dbgbvr11_el1",
		"dbgbvr12_el1",
		"dbgbvr13_el1",
		"dbgbvr14_el1",
		"dbgbvr15_el1",
		"dbgbvr1_el1",
		"dbgbvr2_el1",
		"dbgbvr3_el1",
		"dbgbvr4_el1",
		"dbgbvr5_el1",
		"dbgbvr6_el1",
		"dbgbvr7_el1",
		"dbgbvr8_el1",
		"dbgbvr9_el1",
		"dbgwcr0_el1",
		"dbgwcr10_el1",
		"dbgwcr11_el1",
		"dbgwcr12_el1",
		"dbgwcr13_el1",
		"dbgwcr14_el1",
		"dbgwcr15_el1",
		"dbgwcr1_el1",
		"dbgwcr2_el1",
		"dbgwcr3_el1",
		"dbgwcr4_el1",
		"dbgwcr5_el1",
		"dbgwcr6_el1",
		"dbgwcr7_el1",
		"dbgwcr8_el1",
		"dbgwcr9_el1",
		"dbgwvr0_el1",
		"dbgwvr10_el1",
		"dbgwvr11_el1",
		"dbgwvr12_el1",
		"dbgwvr13_el1",
		"dbgwvr14_el1",
		"dbgwvr15_el1",
		"dbgwvr1_el1",
		"dbgwvr2_el1",
		"dbgwvr3_el1",
		"dbgwvr4_el1",
		"dbgwvr5_el1",
		"dbgwvr6_el1",
		"dbgwvr7_el1",
		"dbgwvr8_el1",
		"dbgwvr9_el1",
		"dczid_el0",
		"el1",
		"esr_el1",
		"esr_el2",
		"esr_el3",
		"far_el1",
		"far_el2",
		"far_el3",
		"hacr_el2",
		"hcr_el2",
		"hpfar_el2",
		"hstr_el2",
		"iallu",
		"ivau",
		"ialluis",
		"id_aa64dfr0_el1",
		"id_aa64isar0_el1",
		"id_aa64isar1_el1",
		"id_aa64mmfr0_el1",
		"id_aa64mmfr1_el1",
		"id_aa64pfr0_el1",
		"id_aa64pfr1_el1",
		"ipas2e1is",
		"ipas2le1is",
		"ipas2e1",
		"ipas2le1",
		"isw",
		"ivac",
		"mair_el1",
		"mair_el2",
		"mair_el3",
		"mdccint_el1",
		"mdccsr_el0",
		"mdcr_el2",
		"mdcr_el3",
		"mdrar_el1",
		"mdscr_el1",
		"mvfr0_el1",
		"mvfr1_el1",
		"mvfr2_el1",
		"osdtrrx_el1",
		"osdtrtx_el1",
		"oseccr_el1",
		"oslar_el1",
		"osdlr_el1",
		"oslsr_el1",
		"pan",
		"par_el1",
		"pmccntr_el0",
		"pmceid0_el0",
		"pmceid1_el0",
		"pmcntenset_el0",
		"pmcr_el0",
		"pmcntenclr_el0",
		"pmintenclr_el1",
		"pmintenset_el1",
		"pmovsclr_el0",
		"pmovsset_el0",
		"pmselr_el0",
		"pmswinc_el0",
		"pmuserenr_el0",
		"pmxevcntr_el0",
		"pmxevtyper_el0",
		"rmr_el1",
		"rmr_el2",
		"rmr_el3",
		"rvbar_el1",
		"rvbar_el2",
		"rvbar_el3",
		"s12e0r",
		"s12e0w",
		"s12e1r",
		"s12e1w",
		"s1e0r",
		"s1e0w",
		"s1e1r",
		"s1e1w",
		"s1e2r",
		"s1e2w",
		"s1e3r",
		"s1e3w",
		"scr_el3",
		"sder32_el3",
		"sctlr_el1",
		"sctlr_el2",
		"sctlr_el3",
		"spsel",
		"tcr_el1",
		"tcr_el2",
		"tcr_el3",
		"tpidrro_el0",
		"tpidr_el0",
		"tpidr_el1",
		"tpidr_el2",
		"tpidr_el3",
		"ttbr0_el1",
		"ttbr1_el1",
		"ttbr0_el2",
		"ttbr0_el3",
		"vaae1",
		"vaae1is",
		"vaale1",
		"vaale1is",
		"vae1",
		"vae1is",
		"vae2",
		"vae2is",
		"vae3",
		"vae3is",
		"vale1",
		"vale1is",
		"vale2",
		"vale2is",
		"vale3",
		"vale3is",
		"vbar_el1",
		"vbar_el2",
		"vbar_el3",
		"vmalle1",
		"vmalle1is",
		"vmalls12e1",
		"vmalls12e1is",
		"vmpidr_el0",
		"vpidr_el2",
		"vtcr_el2",
		"vttbr_el2",
		"zva",
		"#0x0",
		"oshld",
		"oshst",
		"osh",
		"#0x4",
		"nshld",
		"nshst",
		"nsh",
		"#0x8",
		"ishld",
		"ishst",
		"ish",
		"#0xc",
		"ld",
		"st",
		"sy",
		"pmevcntr0_el0",
		"pmevcntr1_el0",
		"pmevcntr2_el0",
		"pmevcntr3_el0",
		"pmevcntr4_el0",
		"pmevcntr5_el0",
		"pmevcntr6_el0",
		"pmevcntr7_el0",
		"pmevcntr8_el0",
		"pmevcntr9_el0",
		"pmevcntr10_el0",
		"pmevcntr11_el0",
		"pmevcntr12_el0",
		"pmevcntr13_el0",
		"pmevcntr14_el0",
		"pmevcntr15_el0",
		"pmevcntr16_el0",
		"pmevcntr17_el0",
		"pmevcntr18_el0",
		"pmevcntr19_el0",
		"pmevcntr20_el0",
		"pmevcntr21_el0",
		"pmevcntr22_el0",
		"pmevcntr23_el0",
		"pmevcntr24_el0",
		"pmevcntr25_el0",
		"pmevcntr26_el0",
		"pmevcntr27_el0",
		"pmevcntr28_el0",
		"pmevcntr29_el0",
		"pmevcntr30_el0",

		"pmevtyper0_el0",
		"pmevtyper1_el0",
		"pmevtyper2_el0",
		"pmevtyper3_el0",
		"pmevtyper4_el0",
		"pmevtyper5_el0",
		"pmevtyper6_el0",
		"pmevtyper7_el0",
		"pmevtyper8_el0",
		"pmevtyper9_el0",
		"pmevtyper10_el0",
		"pmevtyper11_el0",
		"pmevtyper12_el0",
		"pmevtyper13_el0",
		"pmevtyper14_el0",
		"pmevtyper15_el0",
		"pmevtyper16_el0",
		"pmevtyper17_el0",
		"pmevtyper18_el0",
		"pmevtyper19_el0",
		"pmevtyper20_el0",
		"pmevtyper21_el0",
		"pmevtyper22_el0",
		"pmevtyper23_el0",
		"pmevtyper24_el0",
		"pmevtyper25_el0",
		"pmevtyper26_el0",
		"pmevtyper27_el0",
		"pmevtyper28_el0",
		"pmevtyper29_el0",
		"pmevtyper30_el0",
		"pmccfiltr_el0",

		"c0",
		"c1",
		"c2",
		"c3",
		"c4",
		"c5",
		"c6",
		"c7",
		"c8",
		"c9",
		"c10",
		"c11",
		"c12",
		"c13",
		"c14",
		"c15",

		"spsr_el1",
		"elr_el1",
		"sp_el0",
		"current_el",
		"nzcv",
		"fpcr",
		"dspsr_el0",
		"daif",
		"fpsr",
		"dlr_el0",
		"spsr_el2",
		"elr_el2",
		"sp_el1",
		"sp_el2",
		"spsr_irq",
		"spsr_abt",
		"spsr_und",
		"spsr_fiq",
		"spsr_el3",
		"elr_el3",
		"ifsr32_el2",
		"fpexc32_el2",
		"cntvoff_el2",

		"midr_el1",
		"mpidr_el1",
		"revidr_el1",
		"id_pfr0_el1",
		"id_pfr1_el1",
		"id_dfr0_el1",
		"id_afr0_el1",
		"id_mmfr0_el1",
		"id_mmfr1_el1",
		"id_mmfr2_el1",
		"id_mmfr3_el1",
		"id_isar0_el1",
		"id_isar1_el1",
		"id_isar2_el1",
		"id_isar3_el1",
		"id_isar4_el1",
		"id_isar5_el1",
		"id_mmfr4_el1",

		"icc_iar0_el1",
		"icc_eoir0_el1",
		"icc_hppir0_el1",
		"icc_bpr0_el1",
		"icc_ap0r0_el1",
		"icc_ap0r1_el1",
		"icc_ap0r2_el1",
		"icc_ap0r3_el1",
		"icc_ap1r0_el1",
		"icc_ap1r1_el1",
		"icc_ap1r2_el1",
		"icc_ap1r3_el1",
		"icc_dir_el1",
		"icc_rpr_el1",
		"icc_iar1_el1",
		"icc_eoir1_el1",
		"icc_hppir1_el1",
		"icc_bpr1_el1",
		"icc_ctlr_el1",
		"icc_sre_el1",
		"icc_igrpen0_el1",
		"icc_igrpen1_el1",

		"icc_asgi1r_el2",
		"icc_sgi0r_el2",
		"ich_ap0r0_el2",
		"ich_ap0r1_el2",
		"ich_ap0r2_el2",
		"ich_ap0r3_el2",
		"ich_ap1r0_el2",
		"ich_ap1r1_el2",
		"ich_ap1r2_el2",
		"ich_ap1r3_el2",
		"ich_ap1r4_el2",
		"icc_hsre_el2",
		"ich_hcr_el2",
		"ich_vtr_el2",
		"ich_misr_el2",
		"ich_eisr_el2",
		"ich_elrsr_el2",
		"ich_vmcr_el2",

		"ich_lr0_el2",
		"ich_lr1_el2",
		"ich_lr2_el2",
		"ich_lr3_el2",
		"ich_lr4_el2",
		"ich_lr5_el2",
		"ich_lr6_el2",
		"ich_lr7_el2",
		"ich_lr8_el2",
		"ich_lr9_el2",
		"ich_lr10_el2",
		"ich_lr11_el2",
		"ich_lr12_el2",
		"ich_lr13_el2",
		"ich_lr14_el2",
		"ich_lr15_el2",

		"ich_lrc0_el2",
		"ich_lrc1_el2",
		"ich_lrc2_el2",
		"ich_lrc3_el2",
		"ich_lrc4_el2",
		"ich_lrc5_el2",
		"ich_lrc6_el2",
		"ich_lrc7_el2",
		"ich_lrc8_el2",
		"ich_lrc9_el2",
		"ich_lrc10_el2",
		"ich_lrc11_el2",
		"ich_lrc12_el2",
		"ich_lrc13_el2",
		"ich_lrc14_el2",
		"ich_lrc15_el2",

		"icc_mctlr_el3",
		"icc_msre_el3",
		"icc_mgrpen1_el3",

		"teecr32_el1",
		"teehbr32_el1",

		"icc_pmr_el1",
		"icc_sgi1r_el1",
		"icc_sgi0r_el1",
		"icc_asgi1r_el1",
		"icc_seien_el1",
		"END_REG",
	}[s]
}

// typedef union _ieee754 {
// 	uint32_t value;
// 	struct {
// 		uint32_t fraction:23;
// 		uint32_t exponent
// 		uint32_t sign
// 	};
// 	float fvalue;
// }ieee754;

type OperandClass uint32

const (
	NONE OperandClass = iota
	IMM32
	IMM64
	FIMM32
	REG
	MULTI_REG
	SYS_REG
	MEM_REG
	MEM_PRE_IDX
	MEM_POST_IDX
	MEM_OFFSET
	MEM_EXTENDED
	LABEL
	CONDITION
	IMPLEMENTATION_SPECIFIC
)

type Register uint32

const (
	REG_NONE Register = iota
	REG_W0
	REG_W1
	REG_W2
	REG_W3
	REG_W4
	REG_W5
	REG_W6
	REG_W7
	REG_W8
	REG_W9
	REG_W10
	REG_W11
	REG_W12
	REG_W13
	REG_W14
	REG_W15
	REG_W16
	REG_W17
	REG_W18
	REG_W19
	REG_W20
	REG_W21
	REG_W22
	REG_W23
	REG_W24
	REG_W25
	REG_W26
	REG_W27
	REG_W28
	REG_W29
	REG_W30
	REG_WZR
	REG_WSP
	REG_X0
	REG_X1
	REG_X2
	REG_X3
	REG_X4
	REG_X5
	REG_X6
	REG_X7
	REG_X8
	REG_X9
	REG_X10
	REG_X11
	REG_X12
	REG_X13
	REG_X14
	REG_X15
	REG_X16
	REG_X17
	REG_X18
	REG_X19
	REG_X20
	REG_X21
	REG_X22
	REG_X23
	REG_X24
	REG_X25
	REG_X26
	REG_X27
	REG_X28
	REG_X29
	REG_X30
	REG_XZR
	REG_SP
	REG_V0
	REG_V1
	REG_V2
	REG_V3
	REG_V4
	REG_V5
	REG_V6
	REG_V7
	REG_V8
	REG_V9
	REG_V10
	REG_V11
	REG_V12
	REG_V13
	REG_V14
	REG_V15
	REG_V16
	REG_V17
	REG_V18
	REG_V19
	REG_V20
	REG_V21
	REG_V22
	REG_V23
	REG_V24
	REG_V25
	REG_V26
	REG_V27
	REG_V28
	REG_V29
	REG_V30
	REG_VZR
	REG_V31
	REG_B0
	REG_B1
	REG_B2
	REG_B3
	REG_B4
	REG_B5
	REG_B6
	REG_B7
	REG_B8
	REG_B9
	REG_B10
	REG_B11
	REG_B12
	REG_B13
	REG_B14
	REG_B15
	REG_B16
	REG_B17
	REG_B18
	REG_B19
	REG_B20
	REG_B21
	REG_B22
	REG_B23
	REG_B24
	REG_B25
	REG_B26
	REG_B27
	REG_B28
	REG_B29
	REG_B30
	REG_BZR
	REG_B31
	REG_H0
	REG_H1
	REG_H2
	REG_H3
	REG_H4
	REG_H5
	REG_H6
	REG_H7
	REG_H8
	REG_H9
	REG_H10
	REG_H11
	REG_H12
	REG_H13
	REG_H14
	REG_H15
	REG_H16
	REG_H17
	REG_H18
	REG_H19
	REG_H20
	REG_H21
	REG_H22
	REG_H23
	REG_H24
	REG_H25
	REG_H26
	REG_H27
	REG_H28
	REG_H29
	REG_H30
	REG_HZR
	REG_H31
	REG_S0
	REG_S1
	REG_S2
	REG_S3
	REG_S4
	REG_S5
	REG_S6
	REG_S7
	REG_S8
	REG_S9
	REG_S10
	REG_S11
	REG_S12
	REG_S13
	REG_S14
	REG_S15
	REG_S16
	REG_S17
	REG_S18
	REG_S19
	REG_S20
	REG_S21
	REG_S22
	REG_S23
	REG_S24
	REG_S25
	REG_S26
	REG_S27
	REG_S28
	REG_S29
	REG_S30
	REG_SZR
	REG_S31
	REG_D0
	REG_D1
	REG_D2
	REG_D3
	REG_D4
	REG_D5
	REG_D6
	REG_D7
	REG_D8
	REG_D9
	REG_D10
	REG_D11
	REG_D12
	REG_D13
	REG_D14
	REG_D15
	REG_D16
	REG_D17
	REG_D18
	REG_D19
	REG_D20
	REG_D21
	REG_D22
	REG_D23
	REG_D24
	REG_D25
	REG_D26
	REG_D27
	REG_D28
	REG_D29
	REG_D30
	REG_DZR
	REG_D31
	REG_Q0
	REG_Q1
	REG_Q2
	REG_Q3
	REG_Q4
	REG_Q5
	REG_Q6
	REG_Q7
	REG_Q8
	REG_Q9
	REG_Q10
	REG_Q11
	REG_Q12
	REG_Q13
	REG_Q14
	REG_Q15
	REG_Q16
	REG_Q17
	REG_Q18
	REG_Q19
	REG_Q20
	REG_Q21
	REG_Q22
	REG_Q23
	REG_Q24
	REG_Q25
	REG_Q26
	REG_Q27
	REG_Q28
	REG_Q29
	REG_Q30
	REG_QZR
	REG_Q31
	REG_PF0
	REG_PF1
	REG_PF2
	REG_PF3
	REG_PF4
	REG_PF5
	REG_PF6
	REG_PF7
	REG_PF8
	REG_PF9
	REG_PF10
	REG_PF11
	REG_PF12
	REG_PF13
	REG_PF14
	REG_PF15
	REG_PF16
	REG_PF17
	REG_PF18
	REG_PF19
	REG_PF20
	REG_PF21
	REG_PF22
	REG_PF23
	REG_PF24
	REG_PF25
	REG_PF26
	REG_PF27
	REG_PF28
	REG_PF29
	REG_PF30
	REG_PF31
	REG_END
)

func (r Register) String() string {
	return []string{
		"NONE",
		"w0", "w1", "w2", "w3", "w4", "w5", "w6", "w7",
		"w8", "w9", "w10", "w11", "w12", "w13", "w14", "w15",
		"w16", "w17", "w18", "w19", "w20", "w21", "w22", "w23",
		"w24", "w25", "w26", "w27", "w28", "w29", "w30", "wzr", "wsp",
		"x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7",
		"x8", "x9", "x10", "x11", "x12", "x13", "x14", "x15",
		"x16", "x17", "x18", "x19", "x20", "x21", "x22", "x23",
		"x24", "x25", "x26", "x27", "x28", "x29", "x30", "xzr", "sp",
		"v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7",
		"v8", "v9", "v10", "v11", "v12", "v13", "v14", "v15",
		"v16", "v17", "v18", "v19", "v20", "v21", "v22", "v23",
		"v24", "v25", "v26", "v27", "v28", "v29", "v30", "v31", "v31",
		"b0", "b1", "b2", "b3", "b4", "b5", "b6", "b7",
		"b8", "b9", "b10", "b11", "b12", "b13", "b14", "b15",
		"b16", "b17", "b18", "b19", "b20", "b21", "b22", "b23",
		"b24", "b25", "b26", "b27", "b28", "b29", "b30", "b31", "b31",
		"h0", "h1", "h2", "h3", "h4", "h5", "h6", "h7",
		"h8", "h9", "h10", "h11", "h12", "h13", "h14", "h15",
		"h16", "h17", "h18", "h19", "h20", "h21", "h22", "h23",
		"h24", "h25", "h26", "h27", "h28", "h29", "h30", "h31", "h31",
		"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7",
		"s8", "s9", "s10", "s11", "s12", "s13", "s14", "s15",
		"s16", "s17", "s18", "s19", "s20", "s21", "s22", "s23",
		"s24", "s25", "s26", "s27", "s28", "s29", "s30", "s31", "s31",
		"d0", "d1", "d2", "d3", "d4", "d5", "d6", "d7",
		"d8", "d9", "d10", "d11", "d12", "d13", "d14", "d15",
		"d16", "d17", "d18", "d19", "d20", "d21", "d22", "d23",
		"d24", "d25", "d26", "d27", "d28", "d29", "d30", "d31", "d31",
		"q0", "q1", "q2", "q3", "q4", "q5", "q6", "q7",
		"q8", "q9", "q10", "q11", "q12", "q13", "q14", "q15",
		"q16", "q17", "q18", "q19", "q20", "q21", "q22", "q23",
		"q24", "q25", "q26", "q27", "q28", "q29", "q30", "q31", "q31",
		"pldl1keep", "pldl1strm", "pldl2keep", "pldl2strm",
		"pldl3keep", "pldl3strm", "#0x6", "#0x7",
		"plil1keep", "plil1strm", "plil2keep", "plil2strm",
		"plil3keep", "plil3strm", "#0xe", "#0xf",
		"pstl1keep", "pstl1strm", "pstl2keep", "pstl2strm",
		"pstl3keep", "pstl3strm",
		"#0x17", "#0x18", "#0x19", "#0x1a", "#0x1b", "#0x1c", "#0x1d", "#0x1e", "#0x1f",
	}[r]
}

type Condition uint32

const (
	COND_EQ Condition = iota
	COND_NE
	COND_CS
	COND_CC
	COND_MI
	COND_PL
	COND_VS
	COND_VC
	COND_HI
	COND_LS
	COND_GE
	COND_LT
	COND_GT
	COND_LE
	COND_AL
	COND_NV
	END_CONDITION
)

func (c Condition) String() string {
	return []string{
		"eq",
		"ne",
		"cs",
		"cc",
		"mi",
		"pl",
		"vs",
		"vc",
		"hi",
		"ls",
		"ge",
		"lt",
		"gt",
		"le",
		"al",
		"nv",
	}[c]
}

// #define INVERT_CONDITION(N) ((N)^1)

type ShiftType uint32

const (
	SHIFT_NONE ShiftType = iota
	SHIFT_LSL
	SHIFT_LSR
	SHIFT_ASR
	SHIFT_ROR
	SHIFT_UXTW
	SHIFT_SXTW
	SHIFT_SXTX
	SHIFT_UXTX
	SHIFT_SXTB
	SHIFT_SXTH
	SHIFT_UXTH
	SHIFT_UXTB
	SHIFT_MSL
	END_SHIFT
)

func (s ShiftType) String() string {
	return []string{
		"NONE", "lsl", "lsr", "asr",
		"ror", "uxtw", "sxtw", "sxtx",
		"uxtx", "sxtb", "sxth", "uxth",
		"uxtb", "msl",
	}[s]
}

type FailureCodes uint32

const (
	DISASM_SUCCESS FailureCodes = iota
	INVALID_ARGUMENTS
	FAILED_TO_DISASSEMBLE_OPERAND
	FAILED_TO_DISASSEMBLE_OPERATION
	FAILED_TO_DISASSEMBLE_REGISTER
	FAILED_TO_DECODE_INSTRUCTION
	OUTPUT_BUFFER_TOO_SMALL
	OPERAND_IS_NOT_REGISTER
	NOT_MEMORY_OPERAND
)

type Group uint32

const (
	GROUP_UNALLOCATED Group = iota
	GROUP_DATA_PROCESSING_IMM
	GROUP_BRANCH_EXCEPTION_SYSTEM
	GROUP_LOAD_STORE
	GROUP_DATA_PROCESSING_REG
	GROUP_DATA_PROCESSING_SIMD
	GROUP_DATA_PROCESSING_SIMD2
	END_GROUP
)

type InstructionOperand struct {
	OpClass        OperandClass
	Reg            [5]uint32 //registers or conditions
	Scale          uint32
	DataSize       uint32
	ElementSize    uint32
	Index          uint32
	Immediate      uint64
	ShiftType      ShiftType
	ShiftValueUsed uint32
	ShiftValue     uint32
	Extend         ShiftType
	SignedImm      uint32
}

type Instruction struct {
	group     Group
	operation Operation
	operands  [MAX_OPERANDS]InstructionOperand
}

var lsb32Mtable = [33]uint32{
	0x00000000, 0x00000001, 0x00000003,
	0x00000007, 0x0000000f, 0x0000001f,
	0x0000003f, 0x0000007f, 0x000000ff,
	0x000001ff, 0x000003ff, 0x000007ff,
	0x00000fff, 0x00001fff, 0x00003fff,
	0x00007fff, 0x0000ffff, 0x0001ffff,
	0x0003ffff, 0x0007ffff, 0x000fffff,
	0x001fffff, 0x003fffff, 0x007fffff,
	0x00ffffff, 0x01ffffff, 0x03ffffff,
	0x07ffffff, 0x0fffffff, 0x1fffffff,
	0x3fffffff, 0x7fffffff, 0xffffffff,
}

func MaskLSB32(x uint32, nbits uint8) uint32 {
	return x & lsb32Mtable[nbits]
}

func ExtractBits(x uint32, start, nbits int32) uint32 {
	return MaskLSB32(x>>start, uint8(nbits))
}

// //Given a uint32_t instructionValue decopose the instruction
// //into its components -> instruction
// uint32_t aarch64_decompose(
// 		uint32_t instructionValue,
// 		Instruction* restrict instruction,
// 		uint64_t address);

// //Get a text representation of the decomposed instruction
// //into outBuffer
// uint32_t aarch64_disassemble(
// 		Instruction* restrict instruction,
// 		char* outBuffer,
// 		uint32_t outBufferSize);

// //Get the text value of the instruction mnemonic
// const char* get_operation(const Instruction* restrict instruction);

// //Get the text value of a given register enumeration (including prefetch registers)
// //This doesn't handle vectored registers
// const char* get_register_name(Register reg);

// //Get the text value of a given system register
// const char* get_system_register_name(SystemReg reg);

// //Get the text value of a given shift type
// const char* get_shift(ShiftType shift);

// const char* get_condition(Condition cond);

// uint32_t get_implementation_specific(
// 		const InstructionOperand* restrict operand,
// 		char* outBuffer,
// 		uint32_t outBufferSize);

// uint32_t get_register_size(Register reg);