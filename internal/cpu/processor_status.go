package cpu

type ProcessorStatus uint8

const (
	CarryFlag       ProcessorStatus = (1 << 0)
	ZeroFlag        ProcessorStatus = (1 << 1)
	InterruptFlag   ProcessorStatus = (1 << 2)
	DecimalModeFlag ProcessorStatus = (1 << 3) // unused
	BreakFlag       ProcessorStatus = (1 << 4)
	UnusedFlag      ProcessorStatus = (1 << 5)
	OverflowFlag    ProcessorStatus = (1 << 6)
	NegativeFlag    ProcessorStatus = (1 << 7)
)

func NewStatus() *ProcessorStatus {
	initialValue := UnusedFlag
	return &initialValue
}

func (p *ProcessorStatus) Set(flag ProcessorStatus) {
	*p |= flag
}

func (p *ProcessorStatus) Clear(flag ProcessorStatus) {
	*p &^= flag
}

func (p *ProcessorStatus) Has(flag ProcessorStatus) bool {
	return *p&flag != 0
}

func (p *ProcessorStatus) UpdateCond(flag ProcessorStatus, cond bool) {
	if cond {
		p.Set(flag)
	} else {
		p.Clear(flag)
	}
}
