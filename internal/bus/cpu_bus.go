package bus

import (
	"fmt"
	"log/slog"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/memory"
	"github.com/0xmukesh/veyra/internal/utils"
)

type PpuInterface interface {
	ReadRegister(addr uint16) uint8
	WriteRegister(addr uint16, data uint8)
	Tick(cycles uint)
	NMIInterruptStatus() bool
}

type CpuBus struct {
	ram    *memory.RAM
	mapper memory.Mapper
	ppu    PpuInterface

	cycles uint
}

func NewCpuBus(cartridge *cartridge.Cartridge) *CpuBus {
	return &CpuBus{
		ram:    memory.NewRam(2 * 1024),
		mapper: cartridge.Mapper(),
	}
}

func (b *CpuBus) AttachPpu(ppu PpuInterface) {
	b.ppu = ppu
}

func (b *CpuBus) Tick(cycles uint) {
	b.cycles += cycles
	b.ppu.Tick(cycles * 3)
}

func (b *CpuBus) FetchNMIInterruptStatus() bool {
	return b.ppu.NMIInterruptStatus()
}

func (b *CpuBus) Read(addr uint16) uint8 {
	switch {
	case addr <= constants.CPU_RAM_MIRRORS_END:
		return b.ram.Read(addr)
	case addr >= constants.PPU_START && addr <= constants.PPU_END:
		if b.ppu != nil {
			return b.ppu.ReadRegister(addr)
		} else {
			slog.Warn("ppu is not connected to cpu bus")
			return 0
		}
	case addr >= constants.PRGROM_START:
		return b.mapper.Read(addr)
	default:
		slog.Warn("unhandled cpu bus read", slog.String("addr", utils.ToHexadecimalString(addr, 4)))
		return 0
	}
}

func (b *CpuBus) ReadU16(addr uint16) uint16 {
	low := b.Read(addr)
	high := b.Read(addr + 1)
	return utils.PackToLittleEndian(low, high)
}

func (b *CpuBus) Write(addr uint16, data uint8) {
	switch {
	case addr <= constants.CPU_RAM_MIRRORS_END:
		addr &= constants.CPU_RAM_END
		b.ram.Write(addr, data)
	case addr >= constants.PPU_START && addr <= constants.PPU_END:
		if b.ppu != nil {
			b.ppu.WriteRegister(addr, data)
		} else {
			slog.Warn("ppu is not connected")
		}
	case addr >= constants.PRGROM_START:
		panic(fmt.Sprintf("attempt to write to cartidge ROM space - 0x%04x", addr))
	default:
		slog.Warn("ignoring memory write access on cpu bus", slog.String("address", utils.ToHexadecimalString(addr, 4)))
	}
}

func (b *CpuBus) WriteU16(addr uint16, data uint16) {
	low := uint8(data & 0x00ff)
	high := uint8(data >> 8)

	b.Write(addr, low)
	b.Write(addr+1, high)
}
