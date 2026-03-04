#include "bus.hpp"
#include "../constants.hpp"

Bus::Bus() {}

uint8_t Bus::mem_read(uint16_t addr) {
  if (memory_map::RAM_START <= addr && addr <= memory_map::RAM_END) {
    addr = addr & 0b0000011111111111;
  }

  return ram[addr];
}

uint16_t Bus::mem_read_u16(uint16_t addr) {
  auto low = static_cast<uint16_t>(mem_read(addr));
  auto high = static_cast<uint16_t>(mem_read(addr + 1));

  return (high << 8) | low;
}

void Bus::mem_write(uint16_t addr, uint8_t data) { ram[addr] = data; }

void Bus::mem_write_u16(uint16_t addr, uint16_t data) {
  auto low = static_cast<uint8_t>(data & 0xff);
  auto high = static_cast<uint8_t>(data >> 8);

  mem_write(addr, low);
  mem_write(addr + 1, high);
}
