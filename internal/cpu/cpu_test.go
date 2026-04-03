package cpu

import (
	"fmt"
	"os"
	"testing"

	"github.com/0xmukesh/veyra/internal/bus"
	"github.com/0xmukesh/veyra/internal/cartridge"
)

func TestCpu(t *testing.T) {
	raw, err := os.ReadFile("../../roms/nestest.nes")
	if err != nil {
		panic(fmt.Errorf("failed to read nestest.nes file: %w", err))
	}

	cartridge, err := cartridge.New(raw)
	if err != nil {
		panic(fmt.Errorf("failed to initialize cartridge: %w", err))
	}

	bus := bus.NewCpuBus(cartridge)
	cpu := New(bus, 0xc000)

	for !cpu.IsHalted() {
		cpu.Step(true)
	}
}
