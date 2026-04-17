package joypad

import (
	"github.com/0xmukesh/veyra/internal/helpers"
)

const (
	JoypadA      = helpers.Bitflags(1)
	JoypadB      = helpers.Bitflags(1 << 1)
	JoypadSelect = helpers.Bitflags(1 << 2)
	JoypadStart  = helpers.Bitflags(1 << 3)
	JoypadUp     = helpers.Bitflags(1 << 4)
	JoypadDown   = helpers.Bitflags(1 << 5)
	JoypadLeft   = helpers.Bitflags(1 << 6)
	JoypadRight  = helpers.Bitflags(1 << 7)
)

type Joypad struct {
	latch     bool
	state     *helpers.Bitflags
	buttonIdx int
}

func New() *Joypad {
	return &Joypad{
		latch:     false,
		state:     helpers.NewBitflags(0),
		buttonIdx: 0,
	}
}

func (j *Joypad) Write(data uint8) {
	j.latch = data&1 == 1
	if j.latch {
		j.buttonIdx = 0
	}
}

func (j *Joypad) Read() uint8 {
	if j.buttonIdx > 7 {
		return 1
	}

	res := (uint8(*j.state) & (1 << j.buttonIdx)) >> j.buttonIdx
	if !j.latch && j.buttonIdx <= 7 {
		j.buttonIdx++
	}

	return res
}

func (j *Joypad) SetButtonState(button helpers.Bitflags, state bool) {
	j.state.UpdateCond(button, state)
}
