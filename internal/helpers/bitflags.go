package helpers

type Bitflags uint8

func NewBitflags(initialValue uint8) *Bitflags {
	bf := Bitflags(initialValue)
	return &bf
}

func (bf *Bitflags) Set(flag Bitflags) {
	*bf |= flag
}

func (bf *Bitflags) Clear(flag Bitflags) {
	*bf &^= flag
}

func (bf *Bitflags) Has(flag Bitflags) bool {
	return *bf&flag != 0
}

func (bf *Bitflags) UpdateCond(flag Bitflags, cond bool) {
	if cond {
		bf.Set(flag)
	} else {
		bf.Clear(flag)
	}
}
