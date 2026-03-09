#include "registers.hpp"

AddressRegister::AddressRegister() {};
uint16_t AddressRegister::get() {
    return (static_cast<uint16_t>(std::get<0>(value)) << 8) |
           (static_cast<uint16_t>(std::get<1>(value)));
}
void AddressRegister::set(uint16_t data) {
    std::get<0>(value) = (data >> 8);
    std::get<1>(value) = (data & 0xff);
}
void AddressRegister::update(uint8_t data) {
    if (is_high) {
        std::get<0>(value) = data;
    } else {
        std::get<1>(value) = data;
    }

    uint16_t addr = get();
    if (addr > 0x3fff) {
        set(addr & 0x3fff);
    }

    is_high = !is_high;
}
void AddressRegister::reset_latch() { is_high = true; }
