package cpu

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/veyra/internal/constants"
)

func TestLda(t *testing.T) {
	cpu := New()
	program := []uint8{0xa9, 0x03, 0x85, 0x07, 0xaa, 0xa8, 0x00}

	cpu.Load(program, constants.PRGROM_START)
	cpu.Run()
	fmt.Println(cpu.a, cpu.x, cpu.y, cpu.bus.Read(0x07))
}
