package bus

import (
	"log/slog"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/memory"
	"github.com/0xmukesh/veyra/internal/utils"
)

type PpuBus struct {
	cartridge  *cartridge.Cartridge
	vram       *memory.RAM
	paletteRam [32]uint8
}

func NewPpuBus(cartridge *cartridge.Cartridge) *PpuBus {
	return &PpuBus{
		cartridge:  cartridge,
		vram:       memory.NewRam(2 * 1024),
		paletteRam: [32]uint8{0},
	}
}

func (b *PpuBus) Read(addr uint16) uint8 {
	addr &= 0x3fff

	switch {
	case addr < 0x2000:
		return b.cartridge.ReadChrRom(addr)
	default:
		slog.Warn("unhandled ppu bus read", slog.String("addr", utils.ToHexadecimalString(addr, 4)))
		return 0
	}
}
