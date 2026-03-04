cmake_minimum_required(VERSION 3.16)

project(nes-emu LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set(RAYLIB_DIR ${CMAKE_SOURCE_DIR}/external/raylib-5.5/src)

file(GLOB_RECURSE SRC_FILES src/*.cpp)

add_executable(nes-emu ${SRC_FILES})
target_include_directories(nes-emu PRIVATE
    ${CMAKE_SOURCE_DIR}/include
    ${RAYLIB_DIR}
)
target_link_libraries(nes-emu PRIVATE
    ${RAYLIB_DIR}/libraylib.a
)

target_compile_options(nes-emu PRIVATE
    -Wall -Wextra -Wpedantic
)
