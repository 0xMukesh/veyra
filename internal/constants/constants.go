package constants

// cpu
const CPU_CYCLES_PER_FRAME = 29780

// ppu
const PER_SCANLINE_CYCLE_LIFTIME = 341
const NUM_SCANLINES = 262
const NMI_TRIGGER_SCANLINE = 241

// cartidge
const PRGROM_BANK_SIZE uint = 16 * 1024
const CHRROM_BANK_SIZE uint = 8 * 1024

var INES_MAGIC_TAG = []uint8{0x4e, 0x45, 0x53, 0x1a} // "NES^Z"

// gui
const SCREEN_WIDTH = 256
const SCREEN_HEIGHT = 240
const SCALE = 3
