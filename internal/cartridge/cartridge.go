package cartridge

import (
	"fmt"
	"slices"

	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/memory"
)

type Mirroring int

const (
	Horizontal Mirroring = iota
	Vertical
	FourScreen
)

type Cartridge struct {
	prgRom    []uint8
	chrRom    []uint8
	mapper    uint8
	mirroring Mirroring
}

func New(raw []byte) (*Cartridge, error) {
	if !slices.Equal(raw[0:4], constants.INES_MAGIC_TAG) {
		return nil, fmt.Errorf("file is not in iNES format")
	}

	ctrlByteOne := raw[6]
	ctrlByteTwo := raw[7]

	mapperType := (ctrlByteTwo & 0b1111_0000) | (ctrlByteOne >> 4)
	isFourScreen := (ctrlByteOne & (1 << 3)) != 0
	isVertical := (ctrlByteOne & 1) != 0

	var mirroring Mirroring
	if isFourScreen {
		mirroring = FourScreen
	} else if isVertical {
		mirroring = Vertical
	} else {
		mirroring = Horizontal
	}

	prgRomSize := uint(raw[4]) * constants.PRGROM_BANK_SIZE
	chrRomSize := uint(raw[5]) * constants.CHRROM_BANK_SIZE

	skipTrainer := (ctrlByteOne & (1 << 2)) == 0

	prgRomStart := uint(16)
	if !skipTrainer {
		prgRomStart += 512
	}

	chrRomStart := prgRomStart + prgRomSize

	return &Cartridge{
		prgRom:    raw[prgRomStart:(prgRomStart + prgRomSize)],
		chrRom:    raw[chrRomStart:(chrRomStart + chrRomSize)],
		mapper:    mapperType,
		mirroring: mirroring,
	}, nil
}

func (c *Cartridge) Mapper() memory.Mapper {
	switch c.mapper {
	case 0:
		return memory.NewNROM(c.prgRom)
	default:
		panic(fmt.Errorf("unimplemented mapper %d", c.mapper))
	}
}
