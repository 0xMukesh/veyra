#include "cpu.hpp"
#include "../constants.hpp"
#include <cstdint>
#include <ios>
#include <iostream>
#include <optional>
#include <stdexcept>
#include <string>

CPU::CPU()
    : pc(0), sp(STACK_RESET), reg_a(0),
      reg_x(0), reg_y(0), status(flags::UNUSED) {}

void CPU::load_program(const std::vector<uint8_t> &program,
                       uint16_t start_addr) {
  for (size_t i = 0; i < program.size(); ++i) {
    bus.mem_write(start_addr + i, program[i]);
  }

  pc = start_addr;
  bus.mem_write_u16(0xFFFC, start_addr);
}

void CPU::step() {
  uint8_t code = fetch_next_byte();
  const auto &opcode = GetOpTable()[code];
  if (opcode.handler == nullptr) {
      std::cout << "missing opcode handler - " << std::hex << (int)code << std::endl;
      return;
  }

  auto pc_state = pc;
  (this->*opcode.handler)(opcode.mode);

  if (pc == pc_state) {
    pc += opcode.bytes - 1;
  }
}

void CPU::reset() {
  reg_a = 0;
  reg_x = 0;
  reg_y = 0;
  sp = STACK_RESET;
  status = 0b100100;
  pc = mem_read_u16(0xfffc);
}

// load operations
void CPU::op_lda(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  set_register_a(value);
}
void CPU::op_ldx(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  set_register_x(value);
}
void CPU::op_ldy(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  set_register_y(value);
}

// store operations
void CPU::op_sta(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  bus.mem_write(addr, reg_a);
}
void CPU::op_stx(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  bus.mem_write(addr, reg_x);
}
void CPU::op_sty(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  bus.mem_write(addr, reg_y);
}

// register transfers
void CPU::op_tax(AddressingMode) { set_register_x(reg_a); }
void CPU::op_tay(AddressingMode) { set_register_y(reg_a); }
void CPU::op_txa(AddressingMode) { set_register_a(reg_x); }
void CPU::op_tya(AddressingMode) { set_register_a(reg_y); }

// stack operations
void CPU::op_tsx(AddressingMode) { set_register_x(sp); }
void CPU::op_txs(AddressingMode) { sp = reg_x; }
void CPU::op_pha(AddressingMode) { stack_push(reg_a); }
void CPU::op_php(AddressingMode) {
  uint8_t pushed = status | flags::BREAK | flags::UNUSED;
  stack_push(pushed);
}
void CPU::op_pla(AddressingMode) {
  uint8_t data = stack_pop();
  set_register_a(data);
}
void CPU::op_plp(AddressingMode) {
  status = stack_pop();
  set_flag(flags::BREAK, false);
  set_flag(flags::UNUSED, true);
}

// arithmetic
void CPU::op_adc(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  add_to_register_a(value);
}
void CPU::op_sbc(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  add_to_register_a(static_cast<uint8_t>(~value));
}
void CPU::op_cmp(AddressingMode mode) { compare(mode, reg_a); }
void CPU::op_cpx(AddressingMode mode) { compare(mode, reg_x); }
void CPU::op_cpy(AddressingMode mode) { compare(mode, reg_y); }

// logical operations
void CPU::op_and(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  set_register_a(reg_a & value);
}
void CPU::op_eor(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  set_register_a(reg_a ^ value);
}
void CPU::op_ora(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  set_register_a(reg_a | value);
}
void CPU::op_bit(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  auto and_value = reg_a & value;

  set_flag(flags::ZERO, and_value == 0);
  set_flag(flags::NEGATIVE, (value >> 7) == 1);
  set_flag(flags::OVERFLOW, (value >> 6) == 1);
}

// increment operations
void CPU::op_inc(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  value += 1;

  bus.mem_write(addr, value);
  set_flag(flags::ZERO, value == 0);
  set_flag(flags::NEGATIVE, (value >> 7) == 1);
}
void CPU::op_inx(AddressingMode) { set_register_x(reg_x + 1); }
void CPU::op_iny(AddressingMode) { set_register_y(reg_y + 1); }

// decrement operations
void CPU::op_dec(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  uint8_t value = bus.mem_read(addr);
  value -= 1;

  bus.mem_write(addr, value);
  set_flag(flags::ZERO, value == 0);
  set_flag(flags::NEGATIVE, (value >> 7) == 1);
}
void CPU::op_dex(AddressingMode) { set_register_x(reg_x - 1); }
void CPU::op_dey(AddressingMode) { set_register_y(reg_y - 1); }

// shifts
void CPU::op_asl(AddressingMode mode) {
  uint8_t value = 0;
  std::optional<uint16_t> addr;

  if (mode == AddressingMode::Accumulator) {
    value = reg_a;
  } else {
    addr = get_addr(mode);
    value = bus.mem_read(addr.value());
  }

  set_flag(flags::CARRY, (value >> 7) == 1);
  value <<= 1;

  write_to_reg_a_or_mem(mode, addr, value);
}
void CPU::op_lsr(AddressingMode mode) {
  uint8_t value = 0;
  std::optional<uint16_t> addr;

  if (mode == AddressingMode::Accumulator) {
    value = reg_a;
  } else {
    addr = get_addr(mode);
    value = bus.mem_read(addr.value());
  }

  set_flag(flags::CARRY, value & 1);
  value >>= 1;

  write_to_reg_a_or_mem(mode, addr, value);
}
void CPU::op_rol(AddressingMode mode) {
  uint8_t value = 0;
  std::optional<uint16_t> addr;

  if (mode == AddressingMode::Accumulator) {
    value = reg_a;
  } else {
    addr = get_addr(mode);
    value = bus.mem_read(addr.value());
  }

  uint8_t prev_carry_flag = status & flags::CARRY;
  set_flag(flags::CARRY, (value >> 7) == 1);

  value <<= 1;
  value |= prev_carry_flag & 1;

  write_to_reg_a_or_mem(mode, addr, value);
}
void CPU::op_ror(AddressingMode mode) {
  uint8_t value = 0;
  std::optional<uint16_t> addr;

  if (mode == AddressingMode::Accumulator) {
    value = reg_a;
  } else {
    addr = get_addr(mode);
    value = bus.mem_read(addr.value());
  }

  uint8_t prev_carry_flag = status & flags::CARRY;
  set_flag(flags::CARRY, (value & 1) == 1);

  value >>= 1;
  value |= (prev_carry_flag & 1) << 7;

  write_to_reg_a_or_mem(mode, addr, value);
}

// jumps
void CPU::op_jmp(AddressingMode mode) {
  uint16_t addr = get_addr(mode);
  pc = addr;
}

void CPU::op_jsr(AddressingMode mode) {
  uint16_t return_addr = pc + 1;
  uint16_t target = get_addr(mode);
  stack_push_u16(return_addr);
  pc = target;
}
void CPU::op_rts(AddressingMode) { pc = stack_pop_u16() + 1; }

// branches
void CPU::op_bcc(AddressingMode mode) { branch(mode, !get_flag(flags::CARRY)); }
void CPU::op_bcs(AddressingMode mode) { branch(mode, get_flag(flags::CARRY)); }
void CPU::op_bne(AddressingMode mode) { branch(mode, !get_flag(flags::ZERO)); }
void CPU::op_beq(AddressingMode mode) { branch(mode, get_flag(flags::ZERO)); }
void CPU::op_bpl(AddressingMode mode) { branch(mode, !get_flag(flags::NEGATIVE)); }
void CPU::op_bmi(AddressingMode mode) { branch(mode, get_flag(flags::NEGATIVE)); }
void CPU::op_bvc(AddressingMode mode) { branch(mode, !get_flag(flags::OVERFLOW)); }
void CPU::op_bvs(AddressingMode mode) { branch(mode, get_flag(flags::OVERFLOW)); }

// status flag changes
void CPU::op_clc(AddressingMode) { set_flag(flags::CARRY, false); }
void CPU::op_cld(AddressingMode) { set_flag(flags::DECIMAL_MODE, false); }
void CPU::op_cli(AddressingMode) { set_flag(flags::INTERRUPT_DISABLE, false); }
void CPU::op_clv(AddressingMode) { set_flag(flags::OVERFLOW, false); }
void CPU::op_sec(AddressingMode) { set_flag(flags::CARRY, true); }
void CPU::op_sed(AddressingMode) { set_flag(flags::DECIMAL_MODE, true); }
void CPU::op_sei(AddressingMode) { set_flag(flags::INTERRUPT_DISABLE, true); }

// system functions
void CPU::op_brk(AddressingMode) { halted = true; }
void CPU::op_nop(AddressingMode) {}
void CPU::op_rti(AddressingMode) {
  status = stack_pop();
  set_flag(flags::BREAK, false);
  set_flag(flags::UNUSED, true);
  pc = stack_pop_u16();
}

// register utils
void CPU::set_register_a(uint8_t value) {
  reg_a = value;
  update_zero_and_negative_flags(value);
}
void CPU::set_register_x(uint8_t value) {
  reg_x = value;
  update_zero_and_negative_flags(value);
}
void CPU::set_register_y(uint8_t value) {
  reg_y = value;
  update_zero_and_negative_flags(value);
}
void CPU::add_to_register_a(uint8_t value) {
  uint16_t sum = static_cast<uint16_t>(reg_a) + value;
  if (get_flag(flags::CARRY)) {
    sum += 1;
  }

  set_flag(flags::CARRY, sum > 0xff);

  auto result = static_cast<uint8_t>(sum);
  set_flag(flags::OVERFLOW, ((value ^ result) & (result ^ reg_a) & 0x80) != 0);

  set_register_a(result);
}

// flag utils
void CPU::set_flag(uint8_t mask, bool condition) {
  condition ? status |= mask : status &= ~mask;
}
void CPU::update_zero_and_negative_flags(uint8_t value) {
  set_flag(flags::ZERO, value == 0);
  set_flag(flags::NEGATIVE, (value >> 7) == 1);
}
bool CPU::get_flag(uint8_t mask) { return (status & mask) != 0; }

// stack utils
void CPU::stack_push(uint8_t data) {
  bus.mem_write(memory_map::STACK_START + static_cast<uint16_t>(sp), data);
  sp -= 1;
}
void CPU::stack_push_u16(uint16_t data) {
  uint8_t high = data >> 8;
  uint8_t low  = data & 0xff;
  stack_push(high);
  stack_push(low);
}
uint8_t CPU::stack_pop() {
  sp += 1;
  return bus.mem_read(memory_map::STACK_START + static_cast<uint16_t>(sp));
}
uint16_t CPU::stack_pop_u16() {
  auto low  = static_cast<uint16_t>(stack_pop());
  auto high = static_cast<uint16_t>(stack_pop());
  return (high << 8) | low;
}

// mem utils
uint8_t CPU::mem_read(uint16_t addr) { return bus.mem_read(addr); }
uint16_t CPU::mem_read_u16(uint16_t addr) { return bus.mem_read_u16(addr); }
void CPU::mem_write(uint16_t addr, uint8_t data) { bus.mem_write(addr, data); }

// opcode helpers
void CPU::branch(AddressingMode mode, bool condition) {
  int8_t offset = get_addr(mode);
  if (condition) {
    pc += static_cast<uint16_t>(offset);
  }
}

void CPU::compare(AddressingMode mode, uint8_t compare_with) {
  uint16_t addr = get_addr(mode);
  uint8_t data = bus.mem_read(addr);

  set_flag(flags::CARRY, data <= compare_with);
  update_zero_and_negative_flags(compare_with - data);
}

// additional utils
uint8_t CPU::fetch_next_byte() { return bus.mem_read(pc++); }

void CPU::write_to_reg_a_or_mem(AddressingMode mode,
                                std::optional<uint16_t> addr, uint8_t value) {
  if (mode == AddressingMode::Accumulator) {
    set_register_a(value);
  } else {
    if (!addr.has_value()) {
      throw std::runtime_error("failed to extract address");
    }
    bus.mem_write(addr.value(), value);
    update_zero_and_negative_flags(value);
  }
}

uint16_t CPU::get_addr(AddressingMode mode) {
  switch (mode) {
  case Immediate:
    return pc++;
  case Relative:
    return bus.mem_read(pc++);
  case ZeroPage:
    return bus.mem_read(pc++);
  case ZeroPage_X:
    return static_cast<uint8_t>(bus.mem_read(pc++) + reg_x);
  case ZeroPage_Y:
    return static_cast<uint8_t>(bus.mem_read(pc++) + reg_y);
  case Absolute: {
    uint16_t addr = bus.mem_read_u16(pc);
    pc += 2;
    return addr;
  }
  case Absolute_X: {
    uint16_t addr = bus.mem_read_u16(pc) + static_cast<uint16_t>(reg_x);
    pc += 2;
    return addr;
  }
  case Absolute_Y: {
    uint16_t addr = bus.mem_read_u16(pc) + static_cast<uint16_t>(reg_y);
    pc += 2;
    return addr;
  }
  case Indirect: {
    uint16_t ptr = bus.mem_read_u16(pc);
    pc += 2;
    uint16_t value;
    if ((ptr & 0x00FF) == 0x00FF) {
      uint8_t low  = bus.mem_read(ptr);
      uint8_t high = bus.mem_read(ptr & 0xFF00);
      value = (static_cast<uint16_t>(high) << 8) | low;
    } else {
      value = bus.mem_read_u16(ptr);
    }
    return value;
  }
  case Indirect_X: {
    uint8_t base = bus.mem_read(pc++);
    uint8_t zero_page_addr = static_cast<uint8_t>(base + reg_x);
    uint8_t low  = bus.mem_read(zero_page_addr);
    uint8_t high = bus.mem_read(static_cast<uint8_t>(zero_page_addr + 1));
    return (static_cast<uint16_t>(high) << 8) | low;
  }
  case Indirect_Y: {
    uint8_t zero_page_addr = bus.mem_read(pc++);
    uint8_t low  = bus.mem_read(zero_page_addr);
    uint8_t high = bus.mem_read(static_cast<uint8_t>(zero_page_addr + 1));
    uint16_t base_addr = (static_cast<uint16_t>(high) << 8) | low;
    return base_addr + static_cast<uint16_t>(reg_y);
  }
  default:
    throw std::runtime_error("invalid addressing mode: " +
                             std::to_string(static_cast<int>(mode)));
  }
}
