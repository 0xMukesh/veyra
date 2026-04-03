package gui

import (
	"fmt"
	"image/color"

	"github.com/0xmukesh/veyra/internal/cartridge"
	"github.com/0xmukesh/veyra/internal/ppu"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var BANK_IDX = 0

func Start(width, height, scale int32, rom []byte) {
	rl.InitWindow(width*scale, height*scale, "veyra")
	rl.SetTargetFPS(60)

	cartridge, err := cartridge.New(rom)
	if err != nil {
		panic(fmt.Errorf("failed to parse cartridge file - %w", err))
	}

	bankStart := BANK_IDX * 0x1000
	buffer := make([]color.RGBA, width*height)

	for tileIdx := range 256 {
		tile := cartridge.ChrRom()[(bankStart + tileIdx*16):(bankStart + (tileIdx+1)*16)]

		tileX := (tileIdx % 32) * 8
		tileY := (tileIdx / 32) * 8

		for y := range 8 {
			upper := tile[y]
			lower := tile[y+8]

			for x := range 8 {
				x = 7 - x // reversing it
				value := (1&upper)<<1 | (1 & lower)
				upper >>= 1
				lower >>= 1

				var rgb []uint8
				switch value {
				case 0:
					rgb = ppu.SYSTEM_PALETTE[0x01]
				case 1:
					rgb = ppu.SYSTEM_PALETTE[0x23]
				case 2:
					rgb = ppu.SYSTEM_PALETTE[0x27]
				case 3:
					rgb = ppu.SYSTEM_PALETTE[0x30]
				default:
					panic("can't be")
				}

				idx := (y+tileY)*int(width) + (x + tileX)
				buffer[idx] = color.RGBA{
					R: rgb[0],
					G: rgb[1],
					B: rgb[2],
					A: 255,
				}
			}
		}
	}

	img := rl.GenImageColor(int(width), int(height), rl.Black)
	texture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)

	for !rl.WindowShouldClose() {
		rl.UpdateTexture(texture, buffer)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		src := rl.NewRectangle(0, 0, float32(width), float32(height))
		dst := rl.NewRectangle(0, 0, float32(width*scale), float32(height*scale))
		rl.DrawTexturePro(texture, src, dst, rl.Vector2{}, 0, rl.White)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
