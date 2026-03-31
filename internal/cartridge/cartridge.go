package cartridge

import "github.com/0xmukesh/veyra/internal/memory"

type Cartridge struct {
	prg []uint8
}

func New(prg []uint8) *Cartridge {
	return &Cartridge{
		prg: prg,
	}
}

// TODO: mapper would be decided after parsing iNES header
func (c *Cartridge) Mapper() memory.Mapper {
	return memory.NewNROM(c.prg)
}
