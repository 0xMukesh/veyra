#include "rom.hpp"

#include <cstdint>
#include <cstring>
#include <stdexcept>
#include <span>
#include <vector>

NesRom parse_rom(std::vector<uint8_t> &bytes) {
    if (bytes.size() < 4 || std::memcmp(bytes.data(), NES_MAGIC, 4) != 0) {
        throw std::runtime_error("invalid iNES 1.0 header");
    }

    uint8_t ctrl_byte_one = bytes[6];
    uint8_t ctrl_byte_two = bytes[7];
    uint8_t mapper_type = (ctrl_byte_two & 0xF0) | (ctrl_byte_one >> 4);

    uint8_t ines_version = (ctrl_byte_two & 0x0f);
    if (ines_version != 0) {
        throw std::runtime_error("only iNES 1.0 formats are supported");
    }

    bool vert_flag = (ctrl_byte_one & 1) == 1;
    bool four_screen_flag = ((ctrl_byte_one >> 3) & 1) == 1;
    Mirroring mirroring;

    if (four_screen_flag) {
        mirroring = Mirroring::FourScreen;
    } else if (vert_flag) {
        mirroring = Mirroring::Vertical;
    } else {
        mirroring = Mirroring::Horizontal;
    }

    int prg_rom_size = bytes[4] * PRG_ROM_BANK_SIZE;
    int chr_rom_size = bytes[5] * CHR_ROM_BANK_SIZE;

    bool trainer_flag = ((ctrl_byte_one >> 2) & 1) == 1;
    int prg_rom_start = 16;
    if (trainer_flag) {
        prg_rom_start += 512;
    }

    int chr_rom_start = prg_rom_start + prg_rom_size;

    std::span<const uint8_t> prg_rom = {bytes.data() + prg_rom_start,
                                        static_cast<size_t>(prg_rom_size)};
    std::span<const uint8_t> chr_rom = {bytes.data() + chr_rom_start,
                                        static_cast<size_t>(chr_rom_size)};

    return NesRom{
        .prg_rom = std::vector<uint8_t>(prg_rom.begin(), prg_rom.end()),
        .chr_rom = std::vector<uint8_t>(chr_rom.begin(), chr_rom.end()),
        .mapper_type = mapper_type,
        .mirroring = mirroring,
    };
}
