package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0xmukesh/veyra/internal/constants"
	"github.com/0xmukesh/veyra/internal/gui"
)

var (
	romFile = flag.String("rom", "", "path to the rom file")
)

func main() {
	flag.Parse()

	rom, err := os.ReadFile(*romFile)
	if err != nil {
		panic(fmt.Errorf("unable to read rom file - %w", err))
	}

	gui := gui.NewGUI(rom, constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT, constants.SCALE)
	gui.Start()
}
