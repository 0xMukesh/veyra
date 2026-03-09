#include "ppu.hpp"

#include <array>
#include <cstddef>
#include <cstdint>

PPU::PPU(Mirroring mirroring, std::vector<uint8_t> chr_rom)
    : chr_rom(chr_rom), palette_table({}), vram({}), mirroring(mirroring), oam_addr(0),
      oam_data({}) {}
void PPU::write_ctrl_reg(uint8_t value) { ctrl_reg.update(value); }
void PPU::write_mask_reg(uint8_t value) { mask_reg.update(value); }
void PPU::write_oam_addr(uint8_t value) { oam_addr = value; }
void PPU::write_oam_data(uint8_t value) { oam_data[oam_addr] = value; }
void PPU::write_scroll_reg(uint8_t value) { scroll_reg.write(value); }
void PPU::write_addr_reg(uint8_t value) { addr_reg.update(value); }
void PPU::write_data(uint8_t value) {
    // TODO
}
void PPU::write_oam_dma(std::array<uint8_t, 256> data) {
    for (size_t i = 0; i < data.size(); i++) {
        oam_data[i] = data[i];
        oam_addr++;
    }
}
uint8_t PPU::read_status_reg() { return status_reg.get(); }
uint8_t PPU::read_data() {
    // TODO
};
