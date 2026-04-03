package cpu

import (
	"log/slog"
	"os"

	"github.com/0xmukesh/veyra/internal/bus"
	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/helpers"
	"github.com/0xmukesh/veyra/internal/utils"
)

const (
	CarryFlag       = helpers.Bitflags(1 << 0)
	ZeroFlag        = helpers.Bitflags(1 << 1)
	InterruptFlag   = helpers.Bitflags(1 << 2)
	DecimalModeFlag = helpers.Bitflags(1 << 3) // unused
	BreakFlag       = helpers.Bitflags(1 << 4)
	UnusedFlag      = helpers.Bitflags(1 << 5)
	OverflowFlag    = helpers.Bitflags(1 << 6)
	NegativeFlag    = helpers.Bitflags(1 << 7)
)

type CPU struct {
	pc     uint16
	sp     uint8
	a      uint8
	x      uint8
	y      uint8
	status *helpers.Bitflags
	bus    *bus.CpuBus

	halted      bool
	extraCycles int
}

func New(bus *bus.CpuBus, entrypoint uint16) *CPU {
	return &CPU{
		pc:          entrypoint,
		sp:          constants.STACK_RESET,
		a:           0,
		status:      helpers.NewBitflags(uint8(UnusedFlag) | uint8(InterruptFlag)),
		bus:         bus,
		halted:      false,
		extraCycles: 0,
	}
}

func (c *CPU) Step(trace bool) {
	c.extraCycles = 0

	opcode := c.bus.Read(c.pc)
	if trace {
		c.trace(opcode)
	}

	c.pc++
	pcAfterFetch := c.pc

	inst, ok := c.Instructions()[opcode]
	if !ok {
		slog.Error("unknown instruction", slog.String("opcode", utils.ToHexadecimalString(opcode, 2)))
		os.Exit(1)
		return
	}

	inst.handler(inst.mode)
	c.bus.Tick(inst.baseCycles + c.extraCycles)

	if c.pc == pcAfterFetch {
		c.pc += uint16(inst.bytes) - 1
	}
}

func (c *CPU) IsHalted() bool {
	return c.halted
}

func (c *CPU) resolveAddress(mode AddressingMode, addr uint16) (uint16, bool) {
	switch mode {
	case Immediate:
		return c.pc, false
	case ZeroPage:
		return uint16(c.bus.Read(addr)), false
	case ZeroPageX:
		return uint16(c.bus.Read(addr) + c.x), false
	case ZeroPageY:
		return uint16(c.bus.Read(addr) + c.y), false
	case Absolute:
		return c.bus.ReadU16(addr), false
	case AbsoluteX:
		base := c.bus.ReadU16(addr)
		effective := base + uint16(c.x)
		return effective, (base & 0xff00) != (effective & 0xff00)
	case AbsoluteY:
		base := c.bus.ReadU16(addr)
		effective := base + uint16(c.y)
		return effective, (base & 0xff00) != (effective & 0xff00)
	case Indirect:
		base := c.bus.ReadU16(addr)

		if base&0x00ff == 0x00ff {
			low := c.bus.Read(base)
			high := c.bus.Read(base & 0xff00)
			return utils.PackToLittleEndian(low, high), false
		} else {
			return c.bus.ReadU16(base), false
		}
	case IndirectX:
		base := c.bus.Read(addr)
		ptr := base + c.x
		low := c.bus.Read(uint16(ptr))
		high := c.bus.Read(uint16(ptr + 1))

		return utils.PackToLittleEndian(low, high), false
	case IndirectY:
		base := c.bus.Read(addr)
		low := c.bus.Read(uint16(base))
		high := c.bus.Read(uint16(base + 1))
		deref := utils.PackToLittleEndian(low, high)
		effective := deref + uint16(c.y)

		return effective, (deref & 0xff00) != (effective & 0xff00)
	default:
		return 0, false
	}
}

func (c *CPU) getOperandAddress(mode AddressingMode) uint16 {
	addr, _ := c.resolveAddress(mode, c.pc)
	return addr
}

func (c *CPU) readByte(mode AddressingMode) uint8 {
	addr, pageCrossed := c.resolveAddress(mode, c.pc)
	if pageCrossed {
		c.extraCycles++
	}

	return c.bus.Read(addr)
}

// register helpers
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

func (c *CPU) subFromRegisterA(data uint8) {
	// [1]: X - Y  = X + twos complement of Y
	// twos complement of Y = ones complement of Y + 1
	// ones complement of Y = ~Y
	// -Y = twos complement of Y = ~Y + 1
	// [2]: SBC = A - M - (1 - C) = A + (-M) - 1 + C = A + ~M + 1 - 1 + C = A + ~M + C
	c.addToRegisterA(uint8(-int8(data) - 1))
}

// instruction helpers
func (c *CPU) compare(mode AddressingMode, compareWith uint8) {
	value := c.readByte(mode)
	c.status.UpdateCond(CarryFlag, compareWith >= value)
	c.updateZeroAndNegativeFlags(compareWith - value)
}

func (c *CPU) branch(condition bool) {
	if condition {
		jump := int8(c.bus.Read(c.pc))
		// "+ 1" is to jump over the offset operand
		base := c.pc + 1
		dest := base + uint16(jump)

		if (base & 0xff00) != (dest & 0xff00) {
			c.extraCycles += 2
		} else {
			c.extraCycles++
		}

		c.pc = dest
	}
}

// stack
func (c *CPU) stackPush(data uint8) {
	c.bus.Write(constants.STACK_START+uint16(c.sp), data)
	c.sp--
}

func (c *CPU) stackPushU16(data uint16) {
	high := uint8(data >> 8)
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
