package ppu

import (
	"github.com/0xmukesh/veyra/internal/bus"
	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/helpers"
)

// PPUCTRL
const (
	BaseNameTableAddressLow       = helpers.Bitflags(1 << 0)
	BaseNameTableAddressHigh      = helpers.Bitflags(1 << 1)
	VRAMAddIncrement              = helpers.Bitflags(1 << 2)
	SpritePatternAddress          = helpers.Bitflags(1 << 3)
	BackgroundPatternTableAddress = helpers.Bitflags(1 << 4)
	SpriteSize                    = helpers.Bitflags(1 << 5)
	MasterSlaveSelect             = helpers.Bitflags(1 << 6)
	NMIEnable                     = helpers.Bitflags(1 << 7)
)

// PPUMASK
const (
	Greyscale                     = helpers.Bitflags(1 << 0)
	ShowBackgroundLeftmost8Pixels = helpers.Bitflags(1 << 1)
	ShowSpritesLeftmost8Pixels    = helpers.Bitflags(1 << 2)
	EnableBackgroundRendering     = helpers.Bitflags(1 << 3)
	EnableSpriteRendering         = helpers.Bitflags(1 << 4)
	EmphasizeRed                  = helpers.Bitflags(1 << 5)
	EmphasizeGreen                = helpers.Bitflags(1 << 6)
	EmphasizeBlue                 = helpers.Bitflags(1 << 7)
)

// PPUSTATUS
const (
	SpriteOverflow = helpers.Bitflags(1 << 5)
	SpriteZeroHit  = helpers.Bitflags(1 << 6)
	VBlankFlag     = helpers.Bitflags(1 << 7)
)

type PPU struct {
	addr   uint16
	scroll uint8
	ctrl   *helpers.Bitflags
	mask   *helpers.Bitflags
	status *helpers.Bitflags

	t uint16 // (15-bit) during rendering, specifies the starting coarse-x scroll for the next scanline and the starting y scroll for the screen. outside of rendering, used to hold scrolls and vram address before transferring to v register
	v uint16 // (15-bit) during rendering, used for scroll position. outside of rendering, used to store current vram address
	w bool   // (1-bit) write latch used to indicate whether it is first or second write. if it is false then it is first write else second write. cleared on reading PPUSTATUS register
	x uint8  // (3-bit) used to store the fine-x position of the current scroll

	oamData         [256]uint8
	internalDataBuf uint8

	bus *bus.PpuBus
}

func New(cartridge *cartridge.Cartridge) *PPU {
	return &PPU{
		ctrl:    helpers.NewBitflags(0),
		oamData: [256]uint8{0},
		bus:     bus.NewPpuBus(cartridge),
	}
}

func (p *PPU) ReadRegister(addr uint16) uint8 {
	switch addr & 0x7 {
	case 0, 1, 3, 5, 6:
		panic("attempt to read from write-only registers")
	case 2:
		p.w = false // clear out write latch
		return uint8(*p.status)
	case 7:
		return p.readPpuData()
	}

	return 0
}

func (p *PPU) WriteRegister(addr uint16, data uint8) {
	switch addr & 0x7 {
	case 0:
		p.ctrl = helpers.NewBitflags(data)
	case 1:
		p.mask = helpers.NewBitflags(data)
	case 3:
		panic("attempt to write to a read-only register")
	case 5:
		p.updatePpuScroll(data)
	case 6:
		p.updatePpuAddr(data)
	}
}

func (p *PPU) readPpuData() uint8 {
	addr := p.v
	vRamAddIncrement := p.ctrl.Has(VRAMAddIncrement)

	if vRamAddIncrement {
		p.v += 32
	} else {
		p.v += 1
	}

	p.v &= 0x3fff

	result := p.internalDataBuf
	p.internalDataBuf = p.bus.Read(addr)

	return result
}

func (p *PPU) updatePpuAddr(data uint8) {
	if !p.w {
		// bit 14 is forced to be 0 while writing high byte
		p.t = (p.t & 0xff) | ((uint16(data) & 0b0011_1111) << 8)
	} else {
		p.t = ((p.t & 0b0111_1111_0000_0000) | uint16(data)) & 0x3fff
		p.v = p.t
	}

	p.w = !p.w
}

// layout of t register during rendering:
// bit: 14  13  12 | 11  10 | 9   8   7   6   5 | 4   3   2   1   0
// map: [  fine Y ] [nt Y X] [    coarse Y     ] [    coarse X    ]
// fine X is stored in x register
// bit 11 and 10 are written via PPUCTRL
func (p *PPU) updatePpuScroll(data uint8) {
	if !p.w {
		p.t = (p.t & 0b0111_1111_1110_0000) | (uint16(data) >> 3)
		p.x = data & 0x07
	} else {
		p.t = (p.t & 0b0000_1100_0001_1111) | ((uint16(data) & 0b1111_1000) << 2) | ((uint16(data) & 0x07) << 12)
	}

	p.w = !p.w
}
