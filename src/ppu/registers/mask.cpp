#include "registers.hpp"

MaskRegister::MaskRegister() {}
bool MaskRegister::grayscale() { return bits & mask_flags::GREYSCALE; }
bool MaskRegister::bg_leftmost8_pixels() { return bits & mask_flags::BG_LEFTMOST8_PIXELS; }
bool MaskRegister::sprite_leftmost8_pixels() { return bits & mask_flags::SPRITES_LEFTMOST8_PIXELS; }
bool MaskRegister::background_rendering() { return bits & mask_flags::BACKGROUND_RENDERING; }
bool MaskRegister::sprite_rendering() { return bits & mask_flags::SPRITE_RENDERING; }
std::vector<Color> MaskRegister::emphasize() {
    std::vector<Color> color;

    if (bits & mask_flags::EMPHASIZE_RED) {
        color.push_back(Color::Red);
    }

    if (bits & mask_flags::EMPHASIZE_GREEN) {
        color.push_back(Color::Green);
    }

    if (bits & mask_flags::EMPHASIZE_BLUE) {
        color.push_back(Color::Blue);
    }

    return color;
}
void MaskRegister::update(uint8_t new_bits) { bits = new_bits; }
