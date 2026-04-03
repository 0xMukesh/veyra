package cpu

import "github.com/0xmukesh/veyra/internal/constants"

func (c *CPU) nmi() {
	c.stackPushU16(c.pc)

	flag := *c.status
	flag.Clear(BreakFlag)
	flag.Set(UnusedFlag)
	c.stackPush(uint8(flag))
	c.status.Set(InterruptDisableFlag)

	c.bus.Tick(2)
	c.pc = c.bus.ReadU16(constants.NMI_INTERRUPT_VECTOR_ADDRESS_LOW_BYTE)
}
