#pragma once

#include <array>
#include <cstdint>

class Bus {
public:
  Bus();

  uint8_t mem_read(uint16_t addr);
  uint16_t mem_read_u16(uint16_t addr);
  void mem_write(uint16_t addr, uint8_t data);
  void mem_write_u16(uint16_t addr, uint16_t data);

private:
  std::array<uint8_t, 65536> ram;
};
