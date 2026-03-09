#include "registers.hpp"

ScrollRegister::ScrollRegister() {}
void ScrollRegister::write(uint8_t data) {
    if (is_x) {
        scroll_x = data;
    } else {
        scroll_y = data;
    }

    is_x = !is_x;
}
void ScrollRegister::reset_latch() { is_x = true; }
