package gui

import (
	"fmt"
	"image/color"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/cpu"
	"github.com/0xmukesh/veyra/internal/ppu"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GUI struct {
	rom                  []byte
	width, height, scale int32

	cpu       *cpu.CPU
	cpuBus    *cpu.Bus
	ppu       *ppu.PPU
	ppuBus    *ppu.Bus
	cartridge *cartridge.Cartridge

	texture     rl.Texture2D
	frameBuffer []color.RGBA
}

func NewGUI(rom []byte, width, height, scale int32) *GUI {
	return &GUI{
		rom:    rom,
		width:  width,
		height: height,
		scale:  scale,
	}
}

func (g *GUI) Start() {
	rl.InitWindow(g.width*g.scale, g.height*g.scale, "veyra")
	rl.SetTargetFPS(60)

	cartridge, err := cartridge.New(g.rom)
	if err != nil {
		panic(fmt.Errorf("failed to parse cartridge file - %w", err))
	}

	buffer := make([]color.RGBA, g.width*g.height)
	img := rl.GenImageColor(int(g.width), int(g.height), rl.Black)
	texture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)

	g.texture = texture
	g.frameBuffer = buffer

	ppuBus := ppu.NewBus(cartridge)
	ppu := ppu.New(ppuBus)

	cpuBus := cpu.NewBus(cartridge, g.nmiInterruptCallback)
	cpu := cpu.New(cpuBus)
	cpuBus.AttachCpu(cpu)
	cpuBus.AttachPpu(ppu)
	cpu.Reset()

	g.cpu = cpu
	g.cpuBus = cpuBus
	g.ppu = ppu
	g.ppuBus = ppuBus
	g.cartridge = cartridge

	for !rl.WindowShouldClose() {
		cyclesThisFrame := uint(0)
		for cyclesThisFrame < constants.CPU_CYCLES_PER_FRAME {
			cyclesThisFrame += cpu.Step(false)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		src := rl.NewRectangle(0, 0, float32(g.width), float32(g.height))
		dst := rl.NewRectangle(0, 0, float32(g.width*g.scale), float32(g.height*g.scale))
		rl.DrawTexturePro(texture, src, dst, rl.Vector2{}, 0, rl.White)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func (g *GUI) nmiInterruptCallback(p *ppu.PPU) {
	bank := p.BackgroundPatternTableAddress()
	chrRom := g.cartridge.ChrRom()
	numTiles := 32 * 30 // 32x30 tiles

	for i := range numTiles {
		tileId := g.ppuBus.Read(uint16(0x2000 + i))
		tileX := i % 32
		tileY := i / 32
		tile := chrRom[(bank + uint16(tileId)*16):(bank + uint16(tileId+1)*16)]
		palette := g.getBackgroundPalette(tileX, tileY)

		for y := range 8 {
			upper := tile[y]
			lower := tile[y+8]

			for x := range 8 {
				value := ((upper>>7)&1)<<1 | ((lower >> 7) & 1)
				upper <<= 1
				lower <<= 1

				var rgb []uint8
				switch value {
				case 0:
					rgb = constants.SYSTEM_PALETTE[palette[0]]
				case 1:
					rgb = constants.SYSTEM_PALETTE[palette[1]]
				case 2:
					rgb = constants.SYSTEM_PALETTE[palette[2]]
				case 3:
					rgb = constants.SYSTEM_PALETTE[palette[3]]
				}

				frameBufferIdx := (tileY*8+y)*int(g.width) + (tileX*8 + x)
				g.frameBuffer[frameBufferIdx] = color.RGBA{
					R: rgb[0],
					G: rgb[1],
					B: rgb[2],
					A: 255,
				}
			}
		}
	}

	rl.UpdateTexture(g.texture, g.frameBuffer)
}

func (g *GUI) getBackgroundPalette(tileX, tileY int) [4]uint8 {
	attrTableIdx := tileY/4*8 + tileX/4
	attrByte := g.ppuBus.Read(uint16(0x2000 + 960 + attrTableIdx))
	innerX := tileX % 4 / 2
	innerY := tileY % 4 / 2

	paletteIdx := 0
	switch {
	case innerX == 0 && innerY == 0:
		paletteIdx = int(attrByte) & 0b11
	case innerX == 1 && innerY == 0:
		paletteIdx = int(attrByte>>2) & 0b11
	case innerX == 0 && innerY == 1:
		paletteIdx = int(attrByte>>4) & 0b11
	case innerX == 1 && innerY == 1:
		paletteIdx = int(attrByte>>6) & 0b11
	}

	paletteStart := 1 + (paletteIdx * 4)
	paletteTable := g.ppuBus.PaletteTable()

	return [4]uint8{
		paletteTable[0],
		paletteTable[paletteStart],
		paletteTable[paletteStart+1],
		paletteTable[paletteStart+2],
	}
}
