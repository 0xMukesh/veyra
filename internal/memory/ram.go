package memory

type RAM struct {
	data []uint8
	mask uint16
}

func NewRam(size uint16) *RAM {
	return &RAM{
		data: make([]uint8, size),
		mask: size - 1,
	}
}

func (r *RAM) Read(addr uint16) uint8 {
	return r.data[addr&r.mask]
}

func (r *RAM) Write(addr uint16, data uint8) {
	r.data[addr&r.mask] = data
}
