package memory

import "github.com/0xmukesh/veyra/internal/constants"

type Mapper interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

type NROM struct {
	prg []uint8
}

func NewNROM(prg []uint8) *NROM {
	if len(prg) != 0x8000 {
		panic("NROM (no mirroring) requires 32KB PRG ROM")
	}

	return &NROM{prg: prg}
}

func (m *NROM) Read(addr uint16) uint8 {
	if addr < constants.PRGROM_START {
		return 0
	}

	return m.prg[addr-constants.PRGROM_START]
}

func (m *NROM) Write(addr uint16, data uint8) {
	// read only
}
