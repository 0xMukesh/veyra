#pragma once

#include <cstdint>

// N V _ B D I Z C
namespace flags {
constexpr uint8_t CARRY = (1 << 0);
constexpr uint8_t ZERO = (1 << 1);
constexpr uint8_t INTERRUPT_DISABLE = (1 << 2);
constexpr uint8_t DECIMAL_MODE = (1 << 3);
constexpr uint8_t BREAK = (1 << 4);
constexpr uint8_t UNUSED = (1 << 5);
constexpr uint8_t OVERFLOW = (1 << 6);
constexpr uint8_t NEGATIVE = (1 << 7);
} // namespace flags

namespace memory_map {
constexpr uint16_t RAM_START = 0x0000;
constexpr uint16_t RAM_END = 0x1fff;
constexpr uint16_t STACK_START = 0x0100;
constexpr uint16_t STACK_END = 0x01ff;
constexpr uint16_t PRGROM_START = 0x8000;
} // namespace memory_map

constexpr uint16_t INTERRUPT_VECTOR = 0xfffe;
constexpr uint8_t STACK_RESET = 0xfd;
