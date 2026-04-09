package ppu

import (
	"fmt"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/memory"
)

type Bus struct {
	cartridge    *cartridge.Cartridge
	vram         *memory.RAM
	paletteTable [32]uint8
}

func NewBus(cartridge *cartridge.Cartridge) *Bus {
	return &Bus{
		cartridge:    cartridge,
		vram:         memory.NewRam(2048),
		paletteTable: [32]uint8{},
	}
}

func (b *Bus) ChrRom() []byte {
	return b.cartridge.ChrRom()
}

func (b *Bus) Read(addr uint16) uint8 {
	addr &= 0x3fff

	switch {
	case addr <= 0x1fff:
		return b.cartridge.ChrRom()[addr]
	case addr >= 0x2000 && addr <= 0x2fff:
		return b.vram.Read(b.mirrorVRamAddr(addr))
	case addr >= 0x3000 && addr <= 0x3eff:
		panic(fmt.Errorf("ppu address 0x%04x shouldn't be used", addr))
	case addr >= 0x3f00 && addr <= 0x3fff:
		// 0x3f10, 0x3f14, 0x3f18, 0x3f1c are mirrors of 0x3f00, 0x3f04, 0x3f08, 0x3f0c
		if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c {
			addr -= 0x10
		}

		return b.paletteTable[(addr-0x3f00)&0x1f]
	default:
		panic(fmt.Errorf("unexpected ppu address 0x%04x", addr))
	}
}

func (b *Bus) Write(addr uint16, data uint8) {
	switch {
	case addr <= 0x1fff:
		panic(fmt.Errorf("attempt to write to chr rom (0x%04x)", addr))
	case addr >= 0x2000 && addr <= 0x2fff:
		b.vram.Write(b.mirrorVRamAddr(addr), data)
	case addr >= 0x3000 && addr <= 0x3eff:
		panic(fmt.Errorf("ppu address 0x%04x shouldn't be used", addr))
	case addr >= 0x3f00 && addr <= 0x3fff:
		// 0x3f10, 0x3f14, 0x3f18, 0x3f1c are mirrors of 0x3f00, 0x3f04, 0x3f08, 0x3f0c
		if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c {
			addr -= 0x10
		}

		b.paletteTable[(addr-0x3f00)&0x1f] = data
	}
}

func (b *Bus) mirrorVRamAddr(addr uint16) uint16 {
	addr &= 0x2fff
	vramIndex := addr - 0x2000
	nameTable := vramIndex / 0x400
	mirroring := b.cartridge.Mirroring()

	switch {
	case mirroring == cartridge.Vertical && (nameTable == 2 || nameTable == 3):
		return vramIndex - 0x800
	case mirroring == cartridge.Horizontal && nameTable == 2:
		return vramIndex - 0x400
	case mirroring == cartridge.Horizontal && nameTable == 1:
		return vramIndex - 0x400
	case mirroring == cartridge.Horizontal && nameTable == 3:
		return vramIndex - 0x800
	default:
		return vramIndex
	}
}
