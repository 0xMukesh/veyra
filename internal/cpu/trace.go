package cpu

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/0xmukesh/veyra/internal/utils"
)

func (c *CPU) trace(opcode uint8) {
	begin := c.pc

	inst, ok := c.Instructions()[opcode]
	if !ok {
		slog.Error("unknown instruction", slog.String("opcode", utils.ToHexadecimalString(opcode, 2)))
		os.Exit(1)
		return
	}

	hexDump := []string{}
	hexDump = append(hexDump, utils.ToHexadecimalString(opcode, 2))

	operandAddr, _ := c.resolveAddress(inst.mode, begin+1)
	storedValue := c.bus.Read(operandAddr)

	disassembledArgs := ""

	switch inst.bytes {
	case 1:
		switch opcode {
		case 0x0a, 0x4a, 0x2a, 0x6a:
			disassembledArgs = "A"
		}
	case 2:
		operand := c.bus.Read(begin + 1)
		hexDump = append(hexDump, utils.ToHexadecimalString(operand, 2))

		switch inst.mode {
		case Immediate:
			disassembledArgs = fmt.Sprintf("#$%02X", operand)
		case ZeroPage:
			disassembledArgs = fmt.Sprintf("$%02X = %02X", operandAddr, storedValue)
		case ZeroPageX:
			disassembledArgs = fmt.Sprintf("$%02X,X @ %02X = %02X", operand, operandAddr, storedValue)
		case ZeroPageY:
			disassembledArgs = fmt.Sprintf("$%02X,Y @ %02X = %02X", operand, operandAddr, storedValue)
		case IndirectX:
			disassembledArgs = fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", operand, operand+c.x, operandAddr, storedValue)
		case IndirectY:
			disassembledArgs = fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", operand, operandAddr-uint16(c.y), operandAddr, storedValue)
		case Relative:
			offset := int8(operand)
			target := uint16(int32(begin+2) + int32(offset))
			disassembledArgs = fmt.Sprintf("$%04X", target)
		default:
			slog.Warn("unknown addressing mode with 2 opcodes length", slog.Int("mode", int(inst.mode)))
		}
	case 3:
		operandLow := c.bus.Read(begin + 1)
		operandHigh := c.bus.Read(begin + 2)
		hexDump = append(
			hexDump,
			utils.ToHexadecimalString(operandLow, 2),
			utils.ToHexadecimalString(operandHigh, 2),
		)

		operand := utils.PackToLittleEndian(operandLow, operandHigh)

		switch inst.mode {
		case Absolute:
			if opcode == 0x4c || opcode == 0x20 {
				disassembledArgs = fmt.Sprintf("$%04X", operand)
			} else {
				disassembledArgs = fmt.Sprintf("$%04X = %02X", operandAddr, storedValue)
			}
		case AbsoluteX:
			disassembledArgs = fmt.Sprintf("$%04X,X @ %04X = %02X", operand, operandAddr, storedValue)
		case AbsoluteY:
			disassembledArgs = fmt.Sprintf("$%04X,Y @ %04X = %02X", operand, operandAddr, storedValue)
		case Indirect:
			disassembledArgs = fmt.Sprintf("($%04X) = %04X", operand, operandAddr)
		}
	}

	pcStr := fmt.Sprintf("%04X", begin)
	hexDumpStr := strings.Join(hexDump, " ")
	disassembledStr := fmt.Sprintf("%-31s", inst.mnemonic+" "+disassembledArgs)
	registersStr := fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", c.a, c.x, c.y, *c.status, c.sp)

	isUnofficial := inst.mnemonic[0] == '*'

	var line string
	if isUnofficial {
		cleanMnemonic := inst.mnemonic[1:]
		disassembledStr = fmt.Sprintf("%-31s", cleanMnemonic+" "+disassembledArgs)

		hexDumpStr = fmt.Sprintf("%-9s", hexDumpStr)
		line = fmt.Sprintf("%s  %s*%s %s", pcStr, hexDumpStr, disassembledStr, registersStr)
	} else {
		hexDumpStr = fmt.Sprintf("%-10s", hexDumpStr)
		line = fmt.Sprintf("%s  %s%s %s", pcStr, hexDumpStr, disassembledStr, registersStr)
	}

	fmt.Println(line)
}
