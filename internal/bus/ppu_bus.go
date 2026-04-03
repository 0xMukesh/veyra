package bus

import (
	"log/slog"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/memory"
	"github.com/0xmukesh/veyra/internal/utils"
)

type PpuBus struct {
	cartridge *cartridge.Cartridge
	vram      *memory.RAM
}

func NewPpuBus(cartridge *cartridge.Cartridge) *PpuBus {
	return &PpuBus{
		cartridge: cartridge,
		vram:      memory.NewRam(2 * 1024),
	}
}

func (b *PpuBus) Read(addr uint16) uint8 {
	addr &= 0x3fff

	switch {
	case addr < 0x2000:
		return b.cartridge.ChrRom()[addr]
	case addr >= 0x2000 && addr <= 0x2fff:
		return b.vram.Read(b.mirrorVRamAddr(addr))
	default:
		slog.Warn("unhandled ppu bus read", slog.String("addr", utils.ToHexadecimalString(addr, 4)))
		return 0
	}
}

// vertical:
//
//	[ A ] [ B ]
//	[ a ] [ b ]

// horizontal:
//
//	[ A ] [ a ]
//	[ B ] [ b ]
func (b *PpuBus) mirrorVRamAddr(addr uint16) uint16 {
	addr &= 0x2fff
	physicalVRamAddr := addr - 0x2000
	nameTableIndex := physicalVRamAddr / 0x400
	mirroring := b.cartridge.Mirroring()

	switch {
	case mirroring == cartridge.Vertical && (nameTableIndex == 2 || nameTableIndex == 3):
		return physicalVRamAddr - 0x800
	case mirroring == cartridge.Horizontal && nameTableIndex == 2:
		return physicalVRamAddr - 0x400
	case mirroring == cartridge.Horizontal && nameTableIndex == 1:
		return physicalVRamAddr - 0x400
	case mirroring == cartridge.Horizontal && nameTableIndex == 3:
		return physicalVRamAddr - 0x800
	default:
		return physicalVRamAddr
	}
}
