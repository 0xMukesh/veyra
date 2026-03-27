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

		0x0a: NewInstruction("ASL", Accumulator, 1, c.asl),
		0x06: NewInstruction("ASL", ZeroPage, 2, c.asl),
		0x16: NewInstruction("ASL", ZeroPageX, 2, c.asl),
		0x0e: NewInstruction("ASL", Absolute, 3, c.asl),
		0x1e: NewInstruction("ASL", AbsoluteX, 3, c.asl),

		0x4a: NewInstruction("LSR", Accumulator, 1, c.lsr),
		0x46: NewInstruction("LSR", ZeroPage, 2, c.lsr),
		0x56: NewInstruction("LSR", ZeroPageX, 2, c.lsr),
		0x4e: NewInstruction("LSR", Absolute, 3, c.lsr),
		0x5e: NewInstruction("LSR", AbsoluteX, 3, c.lsr),

		0x2a: NewInstruction("ROL", Accumulator, 1, c.rol),
		0x26: NewInstruction("ROL", ZeroPage, 2, c.rol),
		0x36: NewInstruction("ROL", ZeroPageX, 2, c.rol),
		0x2e: NewInstruction("ROL", Absolute, 3, c.rol),
		0x3e: NewInstruction("ROL", AbsoluteX, 3, c.rol),

		0x6a: NewInstruction("ROR", Accumulator, 1, c.ror),
		0x66: NewInstruction("ROR", ZeroPage, 2, c.ror),
		0x76: NewInstruction("ROR", ZeroPageX, 2, c.ror),
		0x6e: NewInstruction("ROR", Absolute, 3, c.ror),
		0x7e: NewInstruction("ROR", AbsoluteX, 3, c.ror),

		0x4c: NewInstruction("JMP", Absolute, 3, c.jmp),
		0x6c: NewInstruction("JMP", Indirect, 3, c.jmp),

		0x20: NewInstruction("JSR", Absolute, 3, c.jsr),
		0x60: NewInstruction("RTS", Implicit, 1, c.rts),

		0x90: NewInstruction("BCC", Relative, 2, c.bcc),
		0xb0: NewInstruction("BCS", Relative, 2, c.bcs),
		0xf0: NewInstruction("BEQ", Relative, 2, c.beq),
		0x30: NewInstruction("BMI", Relative, 2, c.bmi),
		0xd0: NewInstruction("BNE", Relative, 2, c.bne),
		0x10: NewInstruction("BPL", Relative, 2, c.bpl),
		0x50: NewInstruction("BVC", Relative, 2, c.bvc),
		0x70: NewInstruction("BVS", Relative, 2, c.bvs),

		0x18: NewInstruction("CLC", Implicit, 1, c.clc),
		0xd8: NewInstruction("CLD", Implicit, 1, c.cld),
		0x58: NewInstruction("CLI", Implicit, 1, c.cli),
		0xb8: NewInstruction("CLV", Implicit, 1, c.clv),
		0x38: NewInstruction("SEC", Implicit, 1, c.sec),
		0xf8: NewInstruction("SED", Implicit, 1, c.sed),
		0x78: NewInstruction("SEI", Implicit, 1, c.sei),

		0x00: NewInstruction("BRK", Implicit, 1, c.brk),
		0xea: NewInstruction("NOP", Implicit, 1, c.nop),
		0x40: NewInstruction("RTI", Implicit, 1, c.rti),
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
	c.status.Set(UnusedFlag)
	c.status.Set(BreakFlag)
	c.stackPush(uint8(*c.status))
}

func (c *CPU) pla(AddressingMode) {
	c.a = c.stackPop()
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) plp(AddressingMode) {
	*c.status = ProcessorStatus(c.stackPop())
	c.status.Clear(BreakFlag)
	c.status.Set(UnusedFlag)
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

func (c *CPU) asl(mode AddressingMode) {
	addr := uint16(0)
	value := c.a

	if mode != Accumulator {
		addr = c.getOperandAddress(mode)
		value = c.bus.Read(addr)
	}

	c.status.UpdateCond(CarryFlag, value>>7 != 0)
	value <<= 1

	if mode != Accumulator {
		c.bus.Write(addr, value)
	} else {
		c.a = value
	}

	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) lsr(mode AddressingMode) {
	addr := uint16(0)
	value := c.a

	if mode != Accumulator {
		addr = c.getOperandAddress(mode)
		value = c.bus.Read(addr)
	}

	c.status.UpdateCond(CarryFlag, value&1 != 0)
	value >>= 1

	if mode != Accumulator {
		c.bus.Write(addr, value)
	} else {
		c.a = value
	}

	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) rol(mode AddressingMode) {
	addr := uint16(0)
	value := c.a

	if mode != Accumulator {
		addr = c.getOperandAddress(mode)
		value = c.bus.Read(addr)
	}

	isCarrySet := c.status.Has(CarryFlag)

	c.status.UpdateCond(CarryFlag, value>>7 != 0)
	value <<= 1

	if isCarrySet {
		value |= 1
	}

	if mode != Accumulator {
		c.bus.Write(addr, value)
	} else {
		c.a = value
	}

	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) ror(mode AddressingMode) {
	addr := uint16(0)
	value := c.a

	if mode != Accumulator {
		addr = c.getOperandAddress(mode)
		value = c.bus.Read(addr)
	}

	isCarrySet := c.status.Has(CarryFlag)

	c.status.UpdateCond(CarryFlag, value&1 != 0)
	value >>= 1

	if isCarrySet {
		value |= (1 << 7)
	}

	if mode != Accumulator {
		c.bus.Write(addr, value)
	} else {
		c.a = value
	}

	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) jmp(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.pc = addr
}

func (c *CPU) jsr(mode AddressingMode) {
	// [1]: 0x8000: JSR $c000
	// 0x8003: NEXT_INSTRUCTION
	// PC = 0xc000 (target memory address)
	// 0x8003 - 1 = 0x8002 is pushed to the stack
	// [2]: at the end of the subroutine, RTS is used which pops out the return point from stack
	// RTS: PC = 0x8002 (popped from stack) + 1
	addr := c.getOperandAddress(mode)
	c.stackPushU16(c.pc - 1)
	c.pc = addr
}

func (c *CPU) rts(AddressingMode) {
	c.pc = c.stackPopU16() + 1
}

func (c *CPU) bcc(AddressingMode) {
	c.branch(!c.status.Has(CarryFlag))
}

func (c *CPU) bcs(AddressingMode) {
	c.branch(c.status.Has(CarryFlag))
}

func (c *CPU) beq(AddressingMode) {
	c.branch(c.status.Has(ZeroFlag))
}

func (c *CPU) bmi(AddressingMode) {
	c.branch(c.status.Has(NegativeFlag))
}

func (c *CPU) bne(AddressingMode) {
	c.branch(!c.status.Has(ZeroFlag))
}

func (c *CPU) bpl(AddressingMode) {
	c.branch(!c.status.Has(NegativeFlag))
}

func (c *CPU) bvc(AddressingMode) {
	c.branch(!c.status.Has(OverflowFlag))
}

func (c *CPU) bvs(AddressingMode) {
	c.branch(c.status.Has(OverflowFlag))
}

func (c *CPU) clc(AddressingMode) {
	c.status.Clear(CarryFlag)
}

func (c *CPU) cld(AddressingMode) {
	c.status.Clear(DecimalModeFlag)
}

func (c *CPU) cli(AddressingMode) {
	c.status.Clear(InterruptFlag)
}

func (c *CPU) clv(AddressingMode) {
	c.status.Clear(OverflowFlag)
}

func (c *CPU) sec(AddressingMode) {
	c.status.Set(CarryFlag)
}

func (c *CPU) sed(AddressingMode) {
	c.status.Set(DecimalModeFlag)
}

func (c *CPU) sei(AddressingMode) {
	c.status.Set(InterruptFlag)
}

func (c *CPU) brk(AddressingMode) {
	c.status.Set(BreakFlag)
}

func (c *CPU) nop(AddressingMode) {}

func (c *CPU) rti(AddressingMode) {
	*c.status = ProcessorStatus(c.stackPop())
	c.status.Clear(BreakFlag)
	c.status.Set(UnusedFlag)
	c.pc = c.stackPopU16()
}
