package constants

// cpu
const STACK_START uint16 = 0x0100
const STACK_RESET uint8 = 0xfd

const CPU_RAM_MIRRORS_END uint16 = 0x1fff
const CPU_RAM_END uint16 = 0x07ff

const PPU_START uint16 = 0x2000
const PPU_END uint16 = 0x3fff

const PRGROM_START uint16 = 0x8000

// ppu
const NUM_SCANLINES = 262
const PER_SCANLINE_CYCLE_LIFTIME = 341
const NMI_TRIGGER_SCANLINE = 241

// interrupt vector addresses
const NMI_INTERRUPT_VECTOR_ADDRESS_LOW_BYTE uint16 = 0xfffa

// cartidge
const PRGROM_BANK_SIZE uint = 16 * 1024
const CHRROM_BANK_SIZE uint = 8 * 1024

var INES_MAGIC_TAG = []uint8{0x4e, 0x45, 0x53, 0x1a} // "NES^Z"

// gui
const SCREEN_WIDTH = 256
const SCREEN_HEIGHT = 240
const SCALE = 3
