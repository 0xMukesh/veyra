package utils

import "strconv"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func PackToLittleEndian(low, high uint8) uint16 {
	return (uint16(low) << 8) | uint16(high)
}

func ToHexadecimalString[T Integer](v T) string {
	return "0x" + strconv.FormatUint(uint64(v), 16)
}
