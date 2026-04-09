package cpu

import (
	"fmt"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/memory"
	"github.com/0xmukesh/veyra/internal/ppu"
	"github.com/0xmukesh/veyra/internal/utils"
)

type Bus struct {
	ram    *memory.RAM
	mapper memory.Mapper
	cpu    *CPU
	ppu    *ppu.PPU

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
		addr &= 0x2007

		if b.ppu != nil {
			return b.ppu.ReadRegister(addr)
		} else {
			panic("ppu is not connected to cpu bus")
		}
	case addr >= 0x4000 && addr <= 0x5FFF:
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
		addr &= 0x2007

		if b.ppu != nil {
			b.ppu.WriteRegister(addr, data)
		} else {
			panic("ppu is not connected to cpu bus")
		}
	case addr == 0x4014:
		return
	case addr >= 0x4000 && addr <= 0x5FFF:
		return
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
