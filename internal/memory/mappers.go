package memory

type Mapper interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

type NROM struct {
	prg []uint8
	ram [0x2000]uint8
}

func NewNROM(prg []uint8) *NROM {
	return &NROM{prg: prg}
}

func (m *NROM) Read(addr uint16) uint8 {
	if addr >= 0x6000 && addr <= 0x7fff {
		return m.ram[addr-0x6000]
	}

	if addr >= 0x8000 {
		if len(m.prg) == 0x4000 {
			return m.prg[(addr-0x8000)%0x4000]
		}

		return m.prg[addr-0x8000]
	}

	return 0
}

func (m *NROM) Write(addr uint16, data uint8) {
	if addr >= 0x6000 && addr <= 0x7fff {
		m.ram[addr-0x6000] = data
	}
}
