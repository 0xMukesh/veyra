package utils

import "fmt"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func PackToLittleEndian(low, high uint8) uint16 {
	return (uint16(low)) | uint16(high)<<8
}

func ToHexadecimalString[T Integer](v T, padding int) string {
	return fmt.Sprintf("0x%0*X", padding, v)
}
