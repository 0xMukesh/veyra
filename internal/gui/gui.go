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

func nmiInterruptCallback(bus *ppu.Bus, texture rl.Texture2D, buffer *[]color.RGBA, width, height, scale int) func(*ppu.PPU) {
	return func(p *ppu.PPU) {
		bank := p.BackgroundPatternTableAddress()
		chrRom := bus.ChrRom()
		numNametables := 32 * 30

		for i := range numNametables {
			tileId := bus.Read(uint16(0x2000 + i))
			tileX := i % 32
			tileY := i / 32
			tile := chrRom[(bank + uint16(tileId)*16):(bank + uint16(tileId+1)*16)]

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
						rgb = ppu.SYSTEM_PALETTE[0x0f]
					case 1:
						rgb = ppu.SYSTEM_PALETTE[0x02]
					case 2:
						rgb = ppu.SYSTEM_PALETTE[0x28]
					case 3:
						rgb = ppu.SYSTEM_PALETTE[0x16]
					}

					(*buffer)[(tileY*8+y)*width+(tileX*8+x)] = color.RGBA{
						R: rgb[0],
						G: rgb[1],
						B: rgb[2],
						A: 255,
					}
				}
			}
		}

		rl.UpdateTexture(texture, *buffer)
	}
}

func Start(width, height, scale int32, rom []byte) {
	rl.InitWindow(width*scale, height*scale, "veyra")
	rl.SetTargetFPS(60)

	cartridge, err := cartridge.New(rom)
	if err != nil {
		panic(fmt.Errorf("failed to parse cartridge file - %w", err))
	}

	buffer := make([]color.RGBA, width*height)
	img := rl.GenImageColor(int(width), int(height), rl.Black)
	texture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)

	ppuBus := ppu.NewBus(cartridge)
	ppu := ppu.New(ppuBus)

	cpuBus := cpu.NewBus(cartridge, nmiInterruptCallback(ppuBus, texture, &buffer, int(width), int(height), int(scale)))
	cpu := cpu.New(cpuBus)
	cpuBus.AttachCpu(cpu)
	cpuBus.AttachPpu(ppu)
	cpu.Reset()

	for !rl.WindowShouldClose() {
		cyclesThisFrame := uint(0)
		for cyclesThisFrame < constants.CPU_CYCLES_PER_FRAME {
			cyclesThisFrame += cpu.Step(false)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		src := rl.NewRectangle(0, 0, float32(width), float32(height))
		dst := rl.NewRectangle(0, 0, float32(width*scale), float32(height*scale))
		rl.DrawTexturePro(texture, src, dst, rl.Vector2{}, 0, rl.White)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
