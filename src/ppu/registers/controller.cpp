#include "registers.hpp"

#include <stdexcept>

ControllerRegister::ControllerRegister() {}
uint16_t ControllerRegister::nametable_addr() {
    switch (value & 0x11) {
    case 0:
        return 0x2000;
    case 1:
        return 0x2400;
    case 2:
        return 0x2800;
    case 3:
        return 0x2c00;
    default:
        throw std::runtime_error("value not possible");
    };
}
int ControllerRegister::vram_add_increment() {
    if (value & controller_flags::VRAM_ADD_INCREMENT) {
        return 32;
    } else {
        return 1;
    }
}
uint16_t ControllerRegister::sprite_pattern_table_addr() {
    if (value & controller_flags::SPRITE_PATTERN_ADDR) {
        return 0x1000;
    } else {
        return 0;
    }
}
uint16_t ControllerRegister::bg_pattern_table_addr() {
    if (value & controller_flags::BACKGROUND_PATTERN_ADDR) {
        return 0x1000;
    } else {
        return 0;
    }
}
int ControllerRegister::sprite_size() {
    if (value & controller_flags::SPRITE_SIZE) {
        return 16;
    } else {
        return 8;
    }
}
bool ControllerRegister::master_slave_select() {
    return value & controller_flags::MASTER_SLAVE_SELECT;
}
bool ControllerRegister::vblank_enabled() { return value & controller_flags::VBLANK_ENABLE; }
void ControllerRegister::update(uint8_t data) { value = data; }
