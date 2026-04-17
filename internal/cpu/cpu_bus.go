package cpu

import (
	"fmt"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/joypad"
	"github.com/0xmukesh/veyra/internal/memory"
	"github.com/0xmukesh/veyra/internal/ppu"
	"github.com/0xmukesh/veyra/internal/utils"
)

type Bus struct {
	ram     *memory.RAM
	mapper  memory.Mapper
	cpu     *CPU
	ppu     *ppu.PPU
	joypad1 *joypad.Joypad
	joypad2 *joypad.Joypad

	cycles            uint
	nmiRenderCallback func(ppu *ppu.PPU)
}

func NewBus(cartridge *cartridge.Cartridge, nmiRenderCallback func(ppu *ppu.PPU)) *Bus {
	return &Bus{
		ram:               memory.NewRam(2 * 1024),
		mapper:            cartridge.Mapper(),
		nmiRenderCallback: nmiRenderCallback,
	}
}

func (b *Bus) AttachPpu(ppu *ppu.PPU) {
	b.ppu = ppu
}

func (b *Bus) AttachCpu(cpu *CPU) {
	b.cpu = cpu
}

func (b *Bus) AttachJoypads(joypad1, joypad2 *joypad.Joypad) {
	b.joypad1 = joypad1
	b.joypad2 = joypad2
}

func (b *Bus) Tick(cycles uint) {
	b.cycles += cycles

	for range cycles {
		nmiBefore := b.ppu.FetchNMIInterruptStatus()
		b.ppu.Tick(3)
		nmiAfter := b.ppu.FetchNMIInterruptStatus()

		if !nmiBefore && nmiAfter {
			b.nmiRenderCallback(b.ppu)

			if b.cpu == nil {
				panic("cpu is not connected to cpu bus")
			} else {
				b.cpu.TriggerNMI()
			}
		}
	}
}

func (b *Bus) Read(addr uint16) uint8 {
	switch {
	case addr <= 0x1fff:
		return b.ram.Read(addr & 0x07ff)
	case addr >= 0x2000 && addr <= 0x3fff:
		if b.ppu != nil {
			return b.ppu.ReadRegister(addr & 0x2007)
		} else {
			panic("read attempt on ppu but it is not connected to cpu bus")
		}
	case addr >= 0x4000 && addr <= 0x4013:
		// ignore apu
		return 0x00
	case addr == 0x4014:
		panic("read attempt on oam dma register")
	case addr == 0x4015:
		// ignore apu
		return 0x00
	case addr == 0x4016:
		if b.joypad1 != nil {
			return b.joypad1.Read()
		} else {
			panic("read attempt on joypad1 but it is not connected to cpu bus")
		}
	case addr == 0x4017:
		if b.joypad2 != nil {
			return b.joypad2.Read()
		} else {
			panic("read attempt on joypad2 but it is not connected to cpu bus")
		}
	case addr >= 0x4018 && addr <= 0x401f:
		// apu and i/o functionalities which are generally disabled
		return 0x00
	case addr >= 0x6000:
		return b.mapper.Read(addr)
	default:
		panic(fmt.Errorf("unhandled cpu bus read - %s", utils.ToHexadecimalString(addr, 4)))
	}
}

func (b *Bus) Write(addr uint16, data uint8) {
	switch {
	case addr <= 0x1fff:
		b.ram.Write(addr&0x07ff, data)
	case addr >= 0x2000 && addr <= 0x3fff:
		if b.ppu != nil {
			b.ppu.WriteRegister(addr&0x2007, data)
		} else {
			panic("write attempt on ppu but it is not connected to cpu bus")
		}
	case addr >= 0x4000 && addr <= 0x4013:
	// ignore apu
	case addr == 0x4014:
		buffer := [256]uint8{}
		addr := uint16(data) << 8

		for i := range 256 {
			buffer[i] = b.Read(addr + uint16(i))
		}

		if b.ppu != nil {
			b.ppu.WriteOAMDMA(buffer)
		} else {
			panic("write attempt to oam dma but ppu is not connected to cpu bus")
		}
	case addr == 0x4015:
	// ignore apu
	case addr == 0x4016:
		if b.joypad1 != nil {
			b.joypad1.Write(data)
		} else {
			panic("write attempt on joypad1 but it is not connected to cpu bus")
		}
	case addr == 0x4017:
		if b.joypad1 != nil {
			b.joypad1.Write(data)
		} else {
			panic("write attempt on joypad2 but it is not connected to cpu bus")
		}
	case addr >= 0x4018 && addr <= 0x401f:
		// apu and i/o functionalities which are generally disabled
	case addr >= 0x6000:
		b.mapper.Write(addr, data)
	default:
		panic(fmt.Errorf("unhandled cpu bus write - %s", utils.ToHexadecimalString(addr, 4)))
	}
}

func (b *Bus) ReadU16(addr uint16) uint16 {
	low := b.Read(addr)
	high := b.Read(addr + 1)
	return utils.PackToLittleEndian(low, high)
}

func (b *Bus) WriteU16(addr uint16, data uint16) {
	low := uint8(data & 0x00ff)
	high := uint8(data >> 8)

	b.Write(addr, low)
	b.Write(addr+1, high)
}
