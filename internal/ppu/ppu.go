package ppu

import (
	"fmt"

	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/helpers"
)

// PPUCTRL
const (
	CtrlBaseNameTableAddressLow       = helpers.Bitflags(1 << 0)
	CtrlBaseNameTableAddressHigh      = helpers.Bitflags(1 << 1)
	CtrlVRAMAddIncrement              = helpers.Bitflags(1 << 2)
	CtrlSpritePatternAddress          = helpers.Bitflags(1 << 3)
	CtrlBackgroundPatternTableAddress = helpers.Bitflags(1 << 4)
	CtrlSpriteSize                    = helpers.Bitflags(1 << 5)
	CtrlMasterSlaveSelect             = helpers.Bitflags(1 << 6)
	CtrlNMIEnable                     = helpers.Bitflags(1 << 7)
)

// PPUMASK
const (
	MaskGreyscale                     = helpers.Bitflags(1 << 0)
	MaskShowBackgroundLeftmost8Pixels = helpers.Bitflags(1 << 1)
	MaskShowSpritesLeftmost8Pixels    = helpers.Bitflags(1 << 2)
	MaskEnableBackgroundRendering     = helpers.Bitflags(1 << 3)
	MaskEnableSpriteRendering         = helpers.Bitflags(1 << 4)
	MaskEmphasizeRed                  = helpers.Bitflags(1 << 5)
	MaskEmphasizeGreen                = helpers.Bitflags(1 << 6)
	MaskEmphasizeBlue                 = helpers.Bitflags(1 << 7)
)

// PPUSTATUS
const (
	StatusSpriteOverflow = helpers.Bitflags(1 << 5)
	StatusSpriteZeroHit  = helpers.Bitflags(1 << 6)
	StatusVBlankFlag     = helpers.Bitflags(1 << 7)
)

// register indices
const (
	ControlRegister = iota
	MaskRegister
	StatusRegister
	OAMAddrRegister
	OAMDataRegister
	ScrollRegister
	AddrRegister
	DataRegister
	OAMDMARegister
)

type PPU struct {
	addr   uint16
	scroll uint8
	ctrl   *helpers.Bitflags
	mask   *helpers.Bitflags
	status *helpers.Bitflags
	t      uint16 // (15-bit) during rendering, specifies the starting coarse-x scroll for the next scanline and the starting y scroll for the screen. outside of rendering, used to hold scrolls and vram address before transferring to v register
	v      uint16 // (15-bit) during rendering, used for scroll position. outside of rendering, used to store current vram address
	w      bool   // (1-bit) write latch used to indicate whether it is first or second write. if it is false then it is first write else second write. cleared on reading PPUSTATUS register
	x      uint8  // (3-bit) used to store the fine-x position of the current scroll

	bus *Bus

	oamAddr uint8
	oamData [256]uint8

	internalDataBuf uint8
	scanline        uint16
	cycles          uint
	nmiLine         bool
}

func New(bus *Bus) *PPU {
	return &PPU{
		ctrl:    helpers.NewBitflags(0),
		mask:    helpers.NewBitflags(0),
		status:  helpers.NewBitflags(0),
		bus:     bus,
		oamAddr: 0,
		oamData: [256]uint8{0},
		nmiLine: false,
		w:       false,
	}
}

func (p *PPU) FetchNMIInterruptStatus() bool {
	return p.nmiLine
}

func (p *PPU) OAMData() [256]uint8 {
	return p.oamData
}

func (p *PPU) BackgroundPatternTableAddress() uint16 {
	flag := p.ctrl.Has(CtrlBackgroundPatternTableAddress)

	if flag {
		return 0x1000
	} else {
		return 0
	}
}

func (p *PPU) Tick(cycles uint) {
	p.cycles += cycles
	for p.cycles >= constants.PER_SCANLINE_CYCLE_LIFTIME {
		p.cycles -= constants.PER_SCANLINE_CYCLE_LIFTIME
		p.scanline++

		if p.scanline == constants.NMI_TRIGGER_SCANLINE {
			p.status.Set(StatusVBlankFlag)
			if p.ctrl.Has(CtrlNMIEnable) {
				p.nmiLine = true
			}
		}

		if p.scanline == 261 {
			p.status.Clear(StatusVBlankFlag)
			p.nmiLine = false
		}

		if p.scanline >= constants.NUM_SCANLINES {
			p.scanline = 0
		}
	}
}

func (p *PPU) ReadRegister(addr uint16) uint8 {
	addr &= 0x7

	switch addr {
	case StatusRegister:
		status := uint8(*p.status)
		p.status.Clear(StatusVBlankFlag)
		p.w = false
		p.nmiLine = false
		return status
	case OAMDataRegister:
		return p.oamData[p.oamAddr]
	case DataRegister:
		return p.readFromDataRegister()
	default:
		panic(fmt.Errorf("attempt to read from a write-only ppu register - 0x200%01x", addr))
	}
}

func (p *PPU) WriteRegister(addr uint16, data uint8) {
	addr &= 0x7

	switch addr {
	case ControlRegister:
		beforeVBlankNmiEnableFlag := p.ctrl.Has(CtrlNMIEnable)
		p.ctrl = helpers.NewBitflags(data)
		afterVBlankNmiEnableFlag := p.ctrl.Has(CtrlNMIEnable)
		vBlankStatus := p.status.Has(StatusVBlankFlag)

		if !beforeVBlankNmiEnableFlag && afterVBlankNmiEnableFlag && vBlankStatus {
			p.nmiLine = true
		}
	case MaskRegister:
		p.mask = helpers.NewBitflags(data)
	case OAMAddrRegister:
		p.oamAddr = data
	case OAMDataRegister:
		p.oamData[p.oamAddr] = data
		p.oamAddr++
	case ScrollRegister:
		p.updatePpuScroll(data)
	case AddrRegister:
		p.updatePpuAddr(data)
	case DataRegister:
		p.writeToDataRegister(data)
	default:
		panic(fmt.Errorf("attempt to write to a read-only ppu register - 0x200%01x", addr))
	}
}

func (p *PPU) readFromDataRegister() uint8 {
	addr := p.v
	vRamAddIncrementStatus := p.ctrl.Has(CtrlVRAMAddIncrement)

	if vRamAddIncrementStatus {
		p.v += 32
	} else {
		p.v += 1
	}

	if addr <= 0x2fff {
		result := p.internalDataBuf
		p.internalDataBuf = p.bus.Read(addr)
		return result
	} else {
		return p.bus.Read(addr)
	}
}

func (p *PPU) writeToDataRegister(data uint8) {
	addr := p.v
	vRamAddIncrementStatus := p.ctrl.Has(CtrlVRAMAddIncrement)

	if vRamAddIncrementStatus {
		p.v += 32
	} else {
		p.v += 1
	}

	p.bus.Write(addr, data)
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
