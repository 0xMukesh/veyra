#include "utils.hpp"

#include <cstdint>
#include <fstream>
#include <ios>
#include <stdexcept>
#include <vector>

std::vector<uint8_t> read_file(const std::string &path) {
    std::ifstream file(path, std::ios::binary | std::ios::ate);
    if (!file) {
        throw std::runtime_error("failed to open file");
    }

    std::streamsize size = file.tellg();
    file.seekg(0, std::ios::beg);

    std::vector<uint8_t> bytes(size);
    if (!file.read(reinterpret_cast<char *>(bytes.data()), size)) {
        throw std::runtime_error("failed to read file");
    }

    return bytes;
};
