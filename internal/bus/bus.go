package bus

import (
	"log/slog"

	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/utils"
)

type Bus struct {
	cpuVRam [2 * 1024]uint8
	prgRom  []uint8
}

func New() *Bus {
	return &Bus{
		cpuVRam: [2 * 1024]uint8{0},
	}
}

func (b *Bus) Read(addr uint16) uint8 {
	switch {
	case addr <= constants.CPU_RAM_MIRRORS_END:
		return b.cpuVRam[addr&constants.CPU_RAM_END]
	case addr >= constants.PRGROM_START:
		addr -= constants.PRGROM_START
		return b.prgRom[addr]
	default:
		slog.Warn("unhandled read", slog.Uint64("addr", uint64(addr)))
		return 0
	}
}

func (b *Bus) ReadU16(addr uint16) uint16 {
	low := b.Read(addr)
	high := b.Read(addr + 1)
	return utils.PackToLittleEndian(low, high)
}

func (b *Bus) Write(addr uint16, data uint8) {
	if addr <= constants.CPU_RAM_MIRRORS_END {
		addr = addr & constants.CPU_RAM_END
		b.cpuVRam[addr] = data
	} else {
		slog.Warn("ignoring memory write access", slog.Uint64("address", uint64(addr)))
	}
}

func (b *Bus) LoadProgram(program []uint8) {
	b.prgRom = program
}
