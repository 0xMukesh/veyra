package cpu

type AddressingMode int

const (
	Implicit AddressingMode = iota
	Accumulator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Relative
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
)

type Instruction struct {
	mnemonic string
	mode     AddressingMode
	bytes    int
	handler  func(AddressingMode)
}

func NewInstruction(mnemonic string, mode AddressingMode, bytes int, handler func(AddressingMode)) Instruction {
	return Instruction{
		mnemonic: mnemonic,
		mode:     mode,
		bytes:    bytes,
		handler:  handler,
	}
}

func (c *CPU) Instructions() map[uint8]Instruction {
	return map[uint8]Instruction{
		0xa9: NewInstruction("LDA", Immediate, 2, c.lda),
		0xa5: NewInstruction("LDA", ZeroPage, 2, c.lda),
		0xb5: NewInstruction("LDA", ZeroPageX, 2, c.lda),
		0xad: NewInstruction("LDA", Absolute, 3, c.lda),
		0xbd: NewInstruction("LDA", AbsoluteX, 3, c.lda),
		0xb9: NewInstruction("LDA", AbsoluteY, 3, c.lda),
		0xa1: NewInstruction("LDA", IndirectX, 2, c.lda),
		0xb1: NewInstruction("LDA", IndirectY, 2, c.lda),

		0xa2: NewInstruction("LDX", Immediate, 2, c.ldx),
		0xa6: NewInstruction("LDX", ZeroPage, 2, c.ldx),
		0xb6: NewInstruction("LDX", ZeroPageX, 2, c.ldx),
		0xae: NewInstruction("LDX", Absolute, 3, c.ldx),
		0xbe: NewInstruction("LDX", AbsoluteY, 3, c.ldx),

		0xa0: NewInstruction("LDY", Immediate, 2, c.ldy),
		0xa4: NewInstruction("LDY", ZeroPage, 2, c.ldy),
		0xb4: NewInstruction("LDY", ZeroPageX, 2, c.ldy),
		0xac: NewInstruction("LDY", Absolute, 3, c.ldy),
		0xbc: NewInstruction("LDY", AbsoluteX, 3, c.ldy),

		0x85: NewInstruction("STA", ZeroPage, 2, c.sta),
		0x95: NewInstruction("STA", ZeroPageX, 2, c.sta),
		0x8d: NewInstruction("STA", Absolute, 3, c.sta),
		0x9d: NewInstruction("STA", AbsoluteX, 3, c.sta),
		0x99: NewInstruction("STA", AbsoluteY, 3, c.sta),
		0x81: NewInstruction("STA", IndirectX, 2, c.sta),
		0x91: NewInstruction("STA", IndirectY, 2, c.sta),

		0x86: NewInstruction("STX", ZeroPage, 2, c.stx),
		0x96: NewInstruction("STX", ZeroPageY, 2, c.stx),
		0x8e: NewInstruction("STX", Absolute, 3, c.stx),

		0x84: NewInstruction("STY", ZeroPage, 2, c.sty),
		0x94: NewInstruction("STY", ZeroPageX, 2, c.sty),
		0x8c: NewInstruction("STY", Absolute, 3, c.sty),

		0xaa: NewInstruction("TAX", Implicit, 1, c.tax),
		0xa8: NewInstruction("TAY", Implicit, 1, c.tay),
		0x8a: NewInstruction("TXA", Implicit, 1, c.txa),
		0x98: NewInstruction("TYA", Implicit, 1, c.tya),
		0xba: NewInstruction("TSX", Implicit, 1, c.tsx),
		0x9a: NewInstruction("TXS", Implicit, 1, c.txs),

		0x48: NewInstruction("PHA", Implicit, 1, c.pha),
		0x08: NewInstruction("PHP", Implicit, 1, c.php),
		0x68: NewInstruction("PLA", Implicit, 1, c.pla),
		0x28: NewInstruction("PLP", Implicit, 1, c.plp),

		0x29: NewInstruction("AND", Immediate, 2, c.and),
		0x25: NewInstruction("AND", ZeroPage, 2, c.and),
		0x35: NewInstruction("AND", ZeroPageX, 2, c.and),
		0x2d: NewInstruction("AND", Absolute, 3, c.and),
		0x3d: NewInstruction("AND", AbsoluteX, 3, c.and),
		0x39: NewInstruction("AND", AbsoluteY, 3, c.and),
		0x21: NewInstruction("AND", IndirectX, 2, c.and),
		0x31: NewInstruction("AND", IndirectY, 2, c.and),

		0x49: NewInstruction("EOR", Immediate, 2, c.eor),
		0x45: NewInstruction("EOR", ZeroPage, 2, c.eor),
		0x55: NewInstruction("EOR", ZeroPageX, 2, c.eor),
		0x4d: NewInstruction("EOR", Absolute, 3, c.eor),
		0x5d: NewInstruction("EOR", AbsoluteX, 3, c.eor),
		0x59: NewInstruction("EOR", AbsoluteY, 3, c.eor),
		0x41: NewInstruction("EOR", IndirectX, 2, c.eor),
		0x51: NewInstruction("EOR", IndirectY, 2, c.eor),

		0x09: NewInstruction("ORA", Immediate, 2, c.ora),
		0x05: NewInstruction("ORA", ZeroPage, 2, c.ora),
		0x15: NewInstruction("ORA", ZeroPageX, 2, c.ora),
		0x0d: NewInstruction("ORA", Absolute, 3, c.ora),
		0x1d: NewInstruction("ORA", AbsoluteX, 3, c.ora),
		0x19: NewInstruction("ORA", AbsoluteY, 3, c.ora),
		0x01: NewInstruction("ORA", IndirectX, 2, c.ora),
		0x11: NewInstruction("ORA", IndirectY, 2, c.ora),

		0x24: NewInstruction("BIT", ZeroPage, 2, c.bit),
		0x2c: NewInstruction("BIT", Absolute, 3, c.bit),

		0x69: NewInstruction("ADC", Immediate, 2, c.adc),
		0x65: NewInstruction("ADC", ZeroPage, 2, c.adc),
		0x75: NewInstruction("ADC", ZeroPageX, 2, c.adc),
		0x6d: NewInstruction("ADC", Absolute, 3, c.adc),
		0x7d: NewInstruction("ADC", AbsoluteX, 3, c.adc),
		0x79: NewInstruction("ADC", AbsoluteY, 3, c.adc),
		0x61: NewInstruction("ADC", IndirectX, 2, c.adc),
		0x71: NewInstruction("ADC", IndirectY, 2, c.adc),

		0xe9: NewInstruction("SBC", Immediate, 2, c.sbc),
		0xe5: NewInstruction("SBC", ZeroPage, 2, c.sbc),
		0xf5: NewInstruction("SBC", ZeroPageX, 2, c.sbc),
		0xed: NewInstruction("SBC", Absolute, 3, c.sbc),
		0xfd: NewInstruction("SBC", AbsoluteX, 3, c.sbc),
		0xf9: NewInstruction("SBC", AbsoluteY, 3, c.sbc),
		0xe1: NewInstruction("SBC", IndirectX, 2, c.sbc),
		0xf1: NewInstruction("SBC", IndirectY, 2, c.sbc),

		0xc9: NewInstruction("CMP", Immediate, 2, c.cmp),
		0xc5: NewInstruction("CMP", ZeroPage, 2, c.cmp),
		0xd5: NewInstruction("CMP", ZeroPageX, 2, c.cmp),
		0xcd: NewInstruction("CMP", Absolute, 3, c.cmp),
		0xdd: NewInstruction("CMP", AbsoluteX, 3, c.cmp),
		0xd9: NewInstruction("CMP", AbsoluteY, 3, c.cmp),
		0xc1: NewInstruction("CMP", IndirectX, 2, c.cmp),
		0xd1: NewInstruction("CMP", IndirectY, 2, c.cmp),

		0xe0: NewInstruction("CPX", Immediate, 2, c.cpx),
		0xe4: NewInstruction("CPX", ZeroPage, 2, c.cpx),
		0xec: NewInstruction("CPX", Absolute, 3, c.cpx),

		0xc0: NewInstruction("CPY", Immediate, 2, c.cpy),
		0xc4: NewInstruction("CPY", ZeroPage, 2, c.cpy),
		0xcc: NewInstruction("CPY", Absolute, 3, c.cpy),

		0xe6: NewInstruction("INC", ZeroPage, 2, c.inc),
		0xf6: NewInstruction("INC", ZeroPageX, 2, c.inc),
		0xee: NewInstruction("INC", Absolute, 3, c.inc),
		0xfe: NewInstruction("INC", AbsoluteX, 3, c.inc),

		0xe8: NewInstruction("INX", Implicit, 1, c.inx),
		0xc8: NewInstruction("INY", Implicit, 1, c.iny),

		0xc6: NewInstruction("DEC", ZeroPage, 2, c.dec),
		0xd6: NewInstruction("DEC", ZeroPageX, 2, c.dec),
		0xce: NewInstruction("DEC", Absolute, 3, c.dec),
		0xde: NewInstruction("DEC", AbsoluteX, 3, c.dec),

		0xca: NewInstruction("DEX", Implicit, 1, c.dex),
		0x88: NewInstruction("DEY", Implicit, 1, c.dey),
	}
}

func (c *CPU) lda(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.a = value
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) ldx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.x = value
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) ldy(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.y = value
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) sta(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.bus.Write(addr, c.a)
}

func (c *CPU) stx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.bus.Write(addr, c.x)
}

func (c *CPU) sty(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.bus.Write(addr, c.y)
}

func (c *CPU) tax(AddressingMode) {
	c.x = c.a
	c.updateZeroAndNegativeFlags(c.x)
}

func (c *CPU) tay(AddressingMode) {
	c.y = c.a
	c.updateZeroAndNegativeFlags(c.y)
}

func (c *CPU) txa(AddressingMode) {
	c.a = c.x
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) tya(AddressingMode) {
	c.a = c.y
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) tsx(AddressingMode) {
	c.x = c.sp
	c.updateZeroAndNegativeFlags(c.x)
}

func (c *CPU) txs(AddressingMode) {
	c.sp = c.x
}

func (c *CPU) pha(AddressingMode) {
	c.stackPush(c.a)
}

func (c *CPU) php(AddressingMode) {
	c.stackPush(uint8(*c.status))
}

func (c *CPU) pla(AddressingMode) {
	c.a = c.stackPop()
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) plp(AddressingMode) {
	*c.status = ProcessorStatus(c.stackPop())
}

func (c *CPU) and(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.a = c.a & value
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) eor(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.a = c.a ^ value
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) ora(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.a = c.a | value
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) bit(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	and := c.a & value

	c.status.UpdateCond(ZeroFlag, and == 0)
	c.status.UpdateCond(OverflowFlag, value&uint8(OverflowFlag) != 0)
	c.status.UpdateCond(NegativeFlag, value&uint8(NegativeFlag) != 0)
}

func (c *CPU) adc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.addToRegisterA(value)
}

func (c *CPU) sbc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	// [1]: X - Y  = X + twos complement of Y
	// twos complement of Y = ones complement of Y + 1
	// ones complement of Y = ~Y
	// -Y = twos complement of Y = ~Y + 1
	// [2]: SBC = A - M - (1 - C) = A + (-M) - 1 + C = A + ~M + 1 - 1 + C = A + ~M + C
	c.addToRegisterA(uint8(-int8(value) - 1))
}

func (c *CPU) cmp(mode AddressingMode) {
	c.compare(mode, c.a)
}

func (c *CPU) cpx(mode AddressingMode) {
	c.compare(mode, c.x)
}

func (c *CPU) cpy(mode AddressingMode) {
	c.compare(mode, c.y)
}

func (c *CPU) inc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)
	value++

	c.bus.Write(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) inx(AddressingMode) {
	c.x++
	c.updateZeroAndNegativeFlags(c.x)
}

func (c *CPU) iny(AddressingMode) {
	c.y++
	c.updateZeroAndNegativeFlags(c.y)
}

func (c *CPU) dec(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)
	value--

	c.bus.Write(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) dex(AddressingMode) {
	c.x--
	c.updateZeroAndNegativeFlags(c.x)
}

func (c *CPU) dey(AddressingMode) {
	c.y--
	c.updateZeroAndNegativeFlags(c.y)
}
