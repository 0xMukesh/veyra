package cpu

import (
	"fmt"
	"log/slog"

	"github.com/0xmukesh/veyra/internal/bus"
	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/utils"
)

type CPU struct {
	pc     uint16
	sp     uint8
	a      uint8
	x      uint8
	y      uint8
	status *ProcessorStatus
	bus    *bus.Bus
}

func New() *CPU {
	return &CPU{
		pc:     0,
		sp:     constants.STACK_RESET,
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

		c.pc += 1

		inst, ok := c.Instructions()[opcode]
		if !ok {
			slog.Warn("unknown instruction", slog.String("opcode", fmt.Sprintf("0x%O2x", opcode)))
			continue
		}

		inst.handler(inst.mode)

		// PC update for relative is handled manually
		if inst.mode != Relative {
			c.pc += uint16(inst.bytes) - 1
		}

		if opcode == 0x00 {
			break
		}
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
	case Indirect:
		addr := c.bus.ReadU16(c.pc)

		if addr&0x00ff == 0x00ff {
			low := c.bus.Read(addr)
			high := c.bus.Read(addr & 0xff00)
			return utils.PackToLittleEndian(low, high)
		} else {
			return c.bus.ReadU16(addr)
		}
	case IndirectX:
		base := c.bus.Read(c.pc)
		ptr := base + c.x
		low := c.bus.Read(uint16(ptr))
		high := c.bus.Read(uint16(ptr + 1))

		return utils.PackToLittleEndian(low, high)
	case IndirectY:
		base := c.bus.Read(c.pc)
		low := c.bus.Read(uint16(base))
		high := c.bus.Read(uint16(base + 1))
		deref := utils.PackToLittleEndian(low, high)

		return uint16(c.bus.Read(deref))
	default:
		return 0
	}
}

func (c *CPU) updateZeroAndNegativeFlags(result uint8) {
	c.status.UpdateCond(ZeroFlag, result == 0)
	c.status.UpdateCond(NegativeFlag, result&uint8(NegativeFlag) != 0)
}

func (c *CPU) addToRegisterA(data uint8) {
	sum := uint16(c.a) + uint16(data)
	if c.status.Has(CarryFlag) {
		sum += 1
	}

	hadCarry := sum > 0xff
	c.status.UpdateCond(CarryFlag, hadCarry)

	result := uint8(sum)

	hadOverflow := (data^result)&(result^c.a)&(1<<7) != 0
	c.status.UpdateCond(OverflowFlag, hadOverflow)

	c.a = result
	c.updateZeroAndNegativeFlags(c.a)
}

func (c *CPU) compare(mode AddressingMode, compareWith uint8) {
	addr := c.getOperandAddress(mode)
	value := c.bus.Read(addr)

	c.status.UpdateCond(CarryFlag, compareWith >= value)
	c.updateZeroAndNegativeFlags(compareWith - value)
}

func (c *CPU) branch(condition bool) {
	if condition {
		jump := int8(c.bus.Read(c.pc))
		// "+ 1" is to jump over the offset operand
		dest := c.pc + 1 + uint16(jump)

		c.pc = dest
	}
}

func (c *CPU) stackPush(data uint8) {
	c.bus.Write(constants.STACK_START+uint16(c.sp), data)
	c.sp--
}

func (c *CPU) stackPushU16(data uint16) {
	high := uint8(data & 0xff00)
	low := uint8(data & 0x00ff)

	c.stackPush(high)
	c.stackPush(low)
}

func (c *CPU) stackPop() uint8 {
	c.sp++
	return c.bus.Read(constants.STACK_START + uint16(c.sp))
}

func (c *CPU) stackPopU16() uint16 {
	low := c.stackPop()
	high := c.stackPop()
	return utils.PackToLittleEndian(low, high)
}
