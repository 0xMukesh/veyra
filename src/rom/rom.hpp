#pragma once

#include <cstdint>
#include <vector>

enum Mirroring {
    Horizontal,
    Vertical,
    FourScreen,
};

struct NesRom {
    std::vector<uint8_t> prg_rom;
    std::vector<uint8_t> chr_rom;
    uint8_t mapper_type;
    Mirroring mirroring;
};

// NES\x1a
constexpr unsigned char NES_MAGIC[4] = {0x4E, 0x45, 0x53, 0x1A};
constexpr int PRG_ROM_BANK_SIZE = 16 * 1024;
constexpr int CHR_ROM_BANK_SIZE = 8 * 1024;

NesRom parse_rom(std::vector<uint8_t> &bytes);
