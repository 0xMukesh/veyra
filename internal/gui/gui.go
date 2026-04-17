package gui

import (
	"fmt"
	"image/color"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/cpu"
	"github.com/0xmukesh/veyra/internal/joypad"
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
	joypad1 := joypad.New()
	joypad2 := joypad.New()

	cpuBus.AttachCpu(cpu)
	cpuBus.AttachPpu(ppu)
	cpuBus.AttachJoypads(joypad1, joypad2)
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

		joypad1.SetButtonState(joypad.JoypadA, rl.IsKeyDown(rl.KeyA))
		joypad1.SetButtonState(joypad.JoypadB, rl.IsKeyDown(rl.KeyS))
		joypad1.SetButtonState(joypad.JoypadSelect, rl.IsKeyDown(rl.KeySpace))
		joypad1.SetButtonState(joypad.JoypadStart, rl.IsKeyDown(rl.KeyEnter))
		joypad1.SetButtonState(joypad.JoypadUp, rl.IsKeyDown(rl.KeyUp))
		joypad1.SetButtonState(joypad.JoypadDown, rl.IsKeyDown(rl.KeyDown))
		joypad1.SetButtonState(joypad.JoypadLeft, rl.IsKeyDown(rl.KeyLeft))
		joypad1.SetButtonState(joypad.JoypadRight, rl.IsKeyDown(rl.KeyRight))

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
	bgPatternBank := p.BackgroundPatternTableAddress()
	chrRom := g.cartridge.ChrRom()
	numTiles := 32 * 30 // 32x30 tiles

	for i := range numTiles {
		tileIdx := g.ppuBus.Read(uint16(0x2000 + i))
		tileX := i % 32
		tileY := i / 32
		tile := chrRom[(bgPatternBank + uint16(tileIdx)*16):(bgPatternBank + uint16(tileIdx)*16 + 16)]
		bgPalette := g.getBackgroundPalette(tileX, tileY)

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
					rgb = constants.SYSTEM_PALETTE[bgPalette[0]]
				case 1:
					rgb = constants.SYSTEM_PALETTE[bgPalette[1]]
				case 2:
					rgb = constants.SYSTEM_PALETTE[bgPalette[2]]
				case 3:
					rgb = constants.SYSTEM_PALETTE[bgPalette[3]]
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

	spritePatternBank := g.ppu.SpritePatternTableAddress()
	oamData := p.OAMData()

	for i := 0; i < len(oamData); i += 4 {
		tileX := int(oamData[i+3])
		tileY := int(oamData[i])
		tileIdx := oamData[i+1]

		attr := oamData[i+2]
		flipHorizontal := (attr & (1 << 6)) != 0
		flipVertical := (attr & (1 << 7)) != 0
		paletteIdx := attr & 0b11
		spritePalette := g.getSpritePalette(paletteIdx)

		tile := chrRom[(spritePatternBank + uint16(tileIdx)*16):(spritePatternBank + uint16(tileIdx)*16 + 16)]

		for y := range 8 {
			upper := tile[y]
			lower := tile[y+8]

			for x := range 8 {
				value := ((upper>>7)&1)<<1 | ((lower >> 7) & 1)
				upper <<= 1
				lower <<= 1

				rgb := constants.SYSTEM_PALETTE[0]
				switch value {
				case 1:
					rgb = constants.SYSTEM_PALETTE[spritePalette[1]]
				case 2:
					rgb = constants.SYSTEM_PALETTE[spritePalette[2]]
				case 3:
					rgb = constants.SYSTEM_PALETTE[spritePalette[3]]
				}

				xCoord := tileX
				yCoord := tileY

				switch {
				case !flipHorizontal && !flipVertical:
					xCoord += x
					yCoord += y
				case flipHorizontal && !flipVertical:
					xCoord += 7 - x
					yCoord += y
				case !flipHorizontal && flipVertical:
					xCoord += x
					yCoord += 7 - y
				case flipHorizontal && flipVertical:
					xCoord += 7 - x
					yCoord += 7 - y
				}

				if value == 0 {
					continue
				}

				if xCoord < 0 || xCoord >= int(g.width) || yCoord < 0 || yCoord >= int(g.height) {
					continue
				}

				frameBufferIdx := (yCoord)*int(g.width) + (xCoord)
				fbValue := color.RGBA{
					R: rgb[0],
					G: rgb[1],
					B: rgb[2],
					A: 255,
				}

				g.frameBuffer[frameBufferIdx] = fbValue
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

func (g *GUI) getSpritePalette(paletteIdx uint8) [4]uint8 {
	start := 0x11 + (paletteIdx * 4)
	paletteTable := g.ppuBus.PaletteTable()

	return [4]uint8{
		paletteTable[0],
		paletteTable[start],
		paletteTable[start+1],
		paletteTable[start+2],
	}
}
