#pragma once

#include <cstdint>
#include <tuple>
#include <vector>

// VPHB SINN
// |||| ||||
// |||| ||++- Base nametable address
// |||| ||    (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
// |||| |+--- VRAM address increment per CPU read/write of PPUDATA
// |||| |     (0: add 1, going across; 1: add 32, going down)
// |||| +---- Sprite pattern table address for 8x8 sprites
// ||||       (0: $0000; 1: $1000; ignored in 8x16 mode)
// |||+------ Background pattern table address (0: $0000; 1: $1000)
// ||+------- Sprite size (0: 8x8 pixels; 1: 8x16 pixels – see PPU OAM#Byte 1)
// |+-------- PPU master/slave select
// |          (0: read backdrop from EXT pins; 1: output color on EXT pins)
// +--------- Vblank NMI enable (0: off, 1: on)
namespace controller_flags {
constexpr uint8_t BASE_NAMETABLE_ADDR1 = (1 << 0);
constexpr uint8_t BASE_NAMETABLE_ADDR2 = (1 << 1);
constexpr uint8_t VRAM_ADD_INCREMENT = (1 << 2);
constexpr uint8_t SPRITE_PATTERN_ADDR = (1 << 3);
constexpr uint8_t BACKGROUND_PATTERN_ADDR = (1 << 4);
constexpr uint8_t SPRITE_SIZE = (1 << 5);
constexpr uint8_t MASTER_SLAVE_SELECT = (1 << 6);
constexpr uint8_t VBLANK_ENABLE = (1 << 7);
}; // namespace controller_flags

// ---- ----
// BGRs bMmG
// |||| ||||
// |||| |||+- Greyscale (0: normal color, 1: greyscale)
// |||| ||+-- 1: Show background in leftmost 8 pixels of screen, 0: Hide
// |||| |+--- 1: Show sprites in leftmost 8 pixels of screen, 0: Hide
// |||| +---- 1: Enable background rendering
// |||+------ 1: Enable sprite rendering
// ||+------- Emphasize red (green on PAL/Dendy)
// |+-------- Emphasize green (red on PAL/Dendy)
// +--------- Emphasize blue
namespace mask_flags {
constexpr uint8_t GREYSCALE = (1 << 0);
constexpr uint8_t BG_LEFTMOST8_PIXELS = (1 << 1);
constexpr uint8_t SPRITES_LEFTMOST8_PIXELS = (1 << 2);
constexpr uint8_t BACKGROUND_RENDERING = (1 << 3);
constexpr uint8_t SPRITE_RENDERING = (1 << 4);
constexpr uint8_t EMPHASIZE_RED = (1 << 5);
constexpr uint8_t EMPHASIZE_GREEN = (1 << 6);
constexpr uint8_t EMPHASIZE_BLUE = (1 << 7);
}; // namespace mask_flags

namespace status_flags {
constexpr uint8_t SPRITE_OVERFLOW = (1 << 5);
constexpr uint8_t SPRITE_ZERO_HIT = (1 << 6);
constexpr uint8_t VBLANK = (1 << 7);
}; // namespace status_flags

enum Color { Red, Green, Blue };

// $0x2000
class ControllerRegister {
  public:
    ControllerRegister();

    uint16_t nametable_addr();
    int vram_add_increment();
    uint16_t sprite_pattern_table_addr();
    uint16_t bg_pattern_table_addr();
    int sprite_size();
    bool master_slave_select();
    bool vblank_enabled();

    void update(uint8_t data);

  private:
    uint8_t value;
};

// $0x2001
class MaskRegister {
  public:
    MaskRegister();

    bool grayscale();
    bool bg_leftmost8_pixels();
    bool sprite_leftmost8_pixels();
    bool background_rendering();
    bool sprite_rendering();
    std::vector<Color> emphasize();

    void update(uint8_t new_bits);

  private:
    uint8_t bits;
};

// $0x2002
class StatusRegister {
  public:
    StatusRegister();

    uint8_t get();
    void set_sprite_overflow(bool status);
    void set_sprite_zero_hit(bool status);
    void set_vblank_flag(bool status);
    bool in_vblank();

  private:
    uint8_t value;
};

// $0x2005
class ScrollRegister {
  public:
    ScrollRegister();

    void write(uint8_t data);
    void reset_latch();

  private:
    uint8_t scroll_x;
    uint8_t scroll_y;
    bool is_x;
};

// $0x2006
class AddressRegister {
  public:
    AddressRegister();

    uint16_t get();
    void set(uint16_t data);
    void update(uint8_t data);
    void reset_latch();

  private:
    // high byte first and then low byte
    std::tuple<uint8_t, uint8_t> value;
    bool is_high;
};
