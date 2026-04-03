package constants

const STACK_START uint16 = 0x0100
const STACK_RESET uint8 = 0xfd

const CPU_RESET_LOW_BYTE uint16 = 0xfffc
const CPU_RAM_MIRRORS_END uint16 = 0x1fff
const CPU_RAM_END uint16 = 0x07ff

const PPU_START uint16 = 0x2000
const PPU_END uint16 = 0x3fff

const PRGROM_START uint16 = 0x8000

const PRGROM_BANK_SIZE uint = 16 * 1024
const CHRROM_BANK_SIZE uint = 8 * 1024

const SCREEN_WIDTH = 32
const SCREEN_HEIGHT = 32
const SCALE = 10

// "NES^Z"
var INES_MAGIC_TAG = []uint8{0x4e, 0x45, 0x53, 0x1a}
