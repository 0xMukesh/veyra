package cpu

import (
	"testing"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/constants"
)

func TestLda(t *testing.T) {
	program := []uint8{0xa9, 0x03, 0x85, 0x07, 0xaa, 0xa8}
	rom := make([]uint8, 0x8000)
	copy(rom, program)

	cartridge := cartridge.New(rom)
	cpu := New(cartridge, constants.PRGROM_START)

	for !cpu.IsHalted() {
		cpu.Step(true)
	}
}
