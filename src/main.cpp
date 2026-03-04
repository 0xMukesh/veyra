#include <array>
#include <chrono>
#include <cstdint>
#include <cstdlib>
#include <ios>
#include <iostream>
#include <random>
#include <raylib.h>
#include <thread>
#include <vector>

#include "cpu/cpu.hpp"
#include "rom/rom.hpp"
#include "utils/utils.hpp"

constexpr int SCREEN_WIDTH = 32;
constexpr int SCREEN_HEIGHT = 32;
constexpr int SCALE_FACTOR = 10;

using FrameBuffer = std::array<uint8_t, SCREEN_HEIGHT * SCREEN_WIDTH * 4>;

Color get_color(uint8_t byte) {
    switch (byte) {
    case 0:
        return BLACK;
    case 1:
        return WHITE;
    case 2:
    case 9:
        return GRAY;
    case 3:
    case 10:
        return RED;
    case 4:
    case 11:
        return GREEN;
    case 5:
    case 12:
        return BLUE;
    case 6:
    case 13:
        return MAGENTA;
    case 7:
    case 14:
        return YELLOW;
    default:
        return LIME;
    }
}

bool read_screen_state(CPU &cpu, FrameBuffer &frame_buffer) {
    bool updated = false;
    int frame_idx = 0;

    for (int i = 0x0200; i < 0x600; i++) {
        Color color = get_color(cpu.mem_read(i));

        if (frame_buffer[frame_idx] != color.r || frame_buffer[frame_idx + 1] != color.g ||
            frame_buffer[frame_idx + 2] != color.b || frame_buffer[frame_idx + 3] != color.a) {
            frame_buffer[frame_idx] = color.r;
            frame_buffer[frame_idx + 1] = color.g;
            frame_buffer[frame_idx + 2] = color.b;
            frame_buffer[frame_idx + 3] = color.a;
            updated = true;
        }

        frame_idx += 4;
    }

    return updated;
}

static uint8_t last_key = 0;

void update_key_state() {
    if (IsKeyDown(KEY_W))
        last_key = 0x77;
    else if (IsKeyDown(KEY_S))
        last_key = 0x73;
    else if (IsKeyDown(KEY_A))
        last_key = 0x61;
    else if (IsKeyDown(KEY_D))
        last_key = 0x64;
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        std::cerr << "invalid usage. expected format: " << argv[0] << " <rom>" << std::endl;
        return 1;
    }

    std::vector<uint8_t> bytes = read_file(argv[1]);
    NesRom rom = parse_rom(bytes);

    CPU cpu = CPU();

    cpu.load_program(rom.prg_rom, 0x0600);
    cpu.reset();

    std::random_device rd;
    std::mt19937 rng(rd());
    std::uniform_int_distribution<int> dist(1, 15);

    FrameBuffer frame_buffer{};

    InitWindow(SCREEN_WIDTH * SCALE_FACTOR, SCREEN_HEIGHT * SCALE_FACTOR, "snake");
    SetTargetFPS(0);
    Texture2D texture = LoadTextureFromImage(GenImageColor(SCREEN_WIDTH, SCREEN_HEIGHT, BLACK));

    while (!WindowShouldClose()) {
        if (!cpu.is_halted()) {
            update_key_state();

            cpu.mem_write(0xff, last_key);
            cpu.step();
            cpu.mem_write(0xfe, static_cast<uint8_t>(dist(rng)));

            if (read_screen_state(cpu, frame_buffer)) {
                UpdateTexture(texture, frame_buffer.data());
                BeginDrawing();
                ClearBackground(BLACK);
                DrawTextureEx(texture, {0, 0}, 0.0f, SCALE_FACTOR, WHITE);
                EndDrawing();
            }

            std::this_thread::sleep_for(std::chrono::nanoseconds(70'000));
        } else {
            std::cout << "ending game..." << std::endl;
            break;
        }
    }

    UnloadTexture(texture);
    CloseWindow();
    return 0;
}
