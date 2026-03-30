package gui

import (
	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/cpu"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Start(gameCode []uint8) {
	cpu := cpu.New()
	cpu.Load(gameCode, 0x0600)

	width := int32(constants.SCREEN_WIDTH * constants.SCALE)
	height := int32(constants.SCREEN_HEIGHT * constants.SCALE)

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(width, height, "test")

	if rl.WindowShouldClose() {
		return
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		text := "hi, NES :)"
		fontSize := int32(20)
		textWidth := rl.MeasureText(text, fontSize)

		rl.DrawText(
			text,
			width/2-textWidth/2,
			height/2-fontSize/2,
			fontSize,
			rl.White,
		)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
