#include "registers.hpp"

StatusRegister::StatusRegister() {}
uint8_t StatusRegister::get() { return value; }
void StatusRegister::set_sprite_overflow(bool status) {
    if (status) {
        value |= status_flags::SPRITE_OVERFLOW;
    } else {
        value &= ~status_flags::SPRITE_OVERFLOW;
    }
}
void StatusRegister::set_sprite_zero_hit(bool status) {
    if (status) {
        value |= status_flags::SPRITE_ZERO_HIT;
    } else {
        value &= ~status_flags::SPRITE_ZERO_HIT;
    }
}
void StatusRegister::set_vblank_flag(bool status) {
    if (status) {
        value |= status_flags::VBLANK;
    } else {
        value &= ~status_flags::VBLANK;
    }
}
bool StatusRegister::in_vblank() { return value & status_flags::VBLANK; }
