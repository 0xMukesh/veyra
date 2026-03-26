package utils

func PackToLittleEndian(high, low uint8) uint16 {
	return (uint16(high) << 8) | uint16(low)
}
