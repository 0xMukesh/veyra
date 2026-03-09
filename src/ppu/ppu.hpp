#pragma once

#include "registers/registers.hpp"
#include "rom/rom.hpp"

#include <array>
#include <cstdint>
#include <vector>

class PPU {
  public:
    PPU(Mirroring mirroring, std::vector<uint8_t> chr_rom);

    void write_ctrl_reg(uint8_t value);
    void write_mask_reg(uint8_t value);
    void write_oam_addr(uint8_t value);
    void write_oam_data(uint8_t value);
    void write_scroll_reg(uint8_t value);
    void write_addr_reg(uint8_t addr);
    void write_data(uint8_t data);
    void write_oam_dma(std::array<uint8_t, 256> data);
    uint8_t read_status_reg();
    uint8_t read_data();

  private:
    std::vector<uint8_t> chr_rom;
    std::array<uint8_t, 32> palette_table;
    std::array<uint8_t, 2048> vram;
    Mirroring mirroring;

    ControllerRegister ctrl_reg;
    MaskRegister mask_reg;
    StatusRegister status_reg;
    uint8_t oam_addr;
    std::array<uint8_t, 256> oam_data;
    ScrollRegister scroll_reg;
    AddressRegister addr_reg;
};
