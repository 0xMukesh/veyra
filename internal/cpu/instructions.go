package cpu

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
	}
}

func (c *CPU) lda(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.a = value
	c.updateZeroAndNegativeFlags(value)
}
