package utils

func PackToLittleEndian(low, high uint8) uint16 {
	return (uint16(low) << 8) | uint16(high)
}
