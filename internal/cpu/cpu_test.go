package cpu

import (
	"fmt"
	"testing"
)

func TestLda(t *testing.T) {
	cpu := New()
	program := []uint8{0xa9, 0x01, 0x00}

	cpu.Load(program)
	cpu.Run()
	fmt.Println(cpu.a)
}
