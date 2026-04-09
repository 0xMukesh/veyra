package cpu

func (c *CPU) nmi() {
	c.stackPushU16(c.pc)

	flag := *c.status
	flag.Clear(BreakFlag)
	flag.Set(UnusedFlag)
	c.stackPush(uint8(flag))
	c.status.Set(InterruptDisableFlag)

	c.pc = c.bus.ReadU16(0xfffa)
}
