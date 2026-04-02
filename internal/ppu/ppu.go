package ppu

import (
	"github.com/0xmukesh/veyra/internal/bus"
	"github.com/0xmukesh/veyra/internal/cartridge"
)

type PPU struct {
	addr uint16

	v uint16
	t uint16
	w bool // if w is false then high byte else low byte

	oamData [64 * 4]uint8
	bus     *bus.PpuBus
}

func New(cartridge *cartridge.Cartridge) *PPU {
	return &PPU{
		oamData: [256]uint8{0},
		bus:     bus.NewPpuBus(cartridge),
	}
}

func (p *PPU) UpdateAddrRegister(data uint8) {
	if !p.w {
		p.t = (p.t & 0xff) | ((uint16(data) & 0x3f) << 8)
	} else {
		p.t = ((p.t & 0x7f00) | uint16(data)) & 0x3fff
		p.v = p.t
	}

	p.w = !p.w
}
