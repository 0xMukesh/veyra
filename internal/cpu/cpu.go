package cpu

import (
	"fmt"
	"log/slog"

	"github.com/0xmukesh/veyra/internal/bus"
	"github.com/0xmukesh/veyra/internal/utils"
)

type CPU struct {
	pc     uint16
	sp     uint8
	a      uint8
	x      uint8
	y      uint8
	status ProcessorStatus
	bus    *bus.Bus
}

type AddressingMode int

const (
	Immediate AddressingMode = iota
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

func New() *CPU {
	return &CPU{
		pc:     0,
		a:      0,
		status: NewStatus(),
		bus:    bus.New(),
	}
}

func (c *CPU) Load(program []uint8) {
	c.bus.LoadProgram(program)
	c.pc = 0x8000
}

func (c *CPU) Run() {
	for {
		opcode := c.bus.Read(c.pc)
		if opcode == 0x00 {
			break
		}

		c.pc += 1

		inst, ok := c.Instructions()[opcode]
		if !ok {
			slog.Warn("unknown instruction", slog.String("opcode", fmt.Sprintf("0x%O2x", opcode)))
			continue
		}

		inst.handler(inst.mode)
		c.pc += uint16(inst.bytes) - 1
	}
}

func (c *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case Immediate:
		return c.pc
	case ZeroPage:
		return uint16(c.bus.Read(c.pc))
	case ZeroPageX:
		return uint16(c.bus.Read(c.pc) + c.x)
	case ZeroPageY:
		return uint16(c.bus.Read(c.pc) + c.y)
	case Absolute:
		return c.bus.ReadU16(c.pc)
	case AbsoluteX:
		return c.bus.ReadU16(c.pc) + uint16(c.x)
	case AbsoluteY:
		return c.bus.ReadU16(c.pc) + uint16(c.y)
	case IndirectX:
		base := c.bus.Read(c.pc)
		ptr := base + c.x
		low := c.bus.Read(uint16(ptr))
		high := c.bus.Read(uint16(ptr + 1))

		return utils.PackToLittleEndian(high, low)
	case IndirectY:
		base := c.bus.Read(c.pc)
		low := c.bus.Read(uint16(base))
		high := c.bus.Read(uint16(base + 1))
		deref := utils.PackToLittleEndian(high, low)

		return uint16(c.bus.Read(deref))
	default:
		return 0
	}
}

func (c *CPU) updateZeroAndNegativeFlags(result uint8) {
	if result == 0 {
		c.status.Set(ZeroFlag)
	} else {
		c.status.Clear(ZeroFlag)
	}

	if result&(1<<7) != 0 {
		c.status.Set(NegativeFlag)
	} else {
		c.status.Clear(NegativeFlag)
	}
}
