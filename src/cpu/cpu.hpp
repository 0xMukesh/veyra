#pragma once

#include "../bus/bus.hpp"
#include "../constants.hpp"
#include <cstdint>
#include <optional>
#include <vector>

enum AddressingMode : uint8_t {
  Accumulator,
  Implied,
  Immediate,
  Relative,
  ZeroPage,
  ZeroPage_X,
  ZeroPage_Y,
  Absolute,
  Absolute_X,
  Absolute_Y,
  Indirect,
  Indirect_X, // indexed indirect
  Indirect_Y, // indirect indexed
};

class CPU {
public:
  CPU();

  void load_program(const std::vector<uint8_t> &program,
                    uint16_t start_addr = memory_map::PRGROM_START);
  void step();

  uint16_t get_pc() const { return pc; }
  uint8_t get_sp() const { return sp; }
  uint8_t get_reg_a() const { return reg_a; }
  uint8_t get_reg_x() const { return reg_x; }
  uint8_t get_reg_y() const { return reg_y; }
  uint8_t get_status() const { return status; }

  // mem utils
  uint8_t mem_read(uint16_t addr);
  uint16_t mem_read_u16(uint16_t addr);

private:
  uint16_t pc;
  uint8_t sp;
  uint8_t reg_a;
  uint8_t reg_x;
  uint8_t reg_y;
  uint8_t status;
  Bus bus;

  // opcode helpers
  struct OpCode {
    void (CPU::*handler)(AddressingMode);
    AddressingMode mode;
    uint8_t bytes;
  };

  static const std::array<OpCode, 256> &GetOpTable();

  // load operations
  void op_lda(AddressingMode);
  void op_ldx(AddressingMode);
  void op_ldy(AddressingMode);
  // store operations
  void op_sta(AddressingMode);
  void op_stx(AddressingMode);
  void op_sty(AddressingMode);
  // register transfers
  void op_tax(AddressingMode);
  void op_tay(AddressingMode);
  void op_txa(AddressingMode);
  void op_tya(AddressingMode);
  // stack operations
  void op_tsx(AddressingMode);
  void op_txs(AddressingMode);
  void op_pha(AddressingMode);
  void op_php(AddressingMode);
  void op_pla(AddressingMode);
  void op_plp(AddressingMode);
  // logical operations (performed on register A)
  void op_and(AddressingMode);
  void op_eor(AddressingMode);
  void op_ora(AddressingMode);
  void op_bit(AddressingMode);
  // arithemtic
  void op_adc(AddressingMode);
  void op_sbc(AddressingMode);
  void op_cmp(AddressingMode);
  void op_cpx(AddressingMode);
  void op_cpy(AddressingMode);
  // increment operations
  void op_inc(AddressingMode);
  void op_inx(AddressingMode);
  void op_iny(AddressingMode);
  // decrement operations
  void op_dec(AddressingMode);
  void op_dex(AddressingMode);
  void op_dey(AddressingMode);
  // shifts
  void op_asl(AddressingMode);
  void op_lsr(AddressingMode);
  void op_rol(AddressingMode);
  void op_ror(AddressingMode);
  // jumps
  void op_jmp(AddressingMode);
  void op_jsr(AddressingMode);
  void op_rts(AddressingMode);
  // branches
  void op_bcc(AddressingMode);
  void op_bcs(AddressingMode);
  void op_beq(AddressingMode);
  void op_bmi(AddressingMode);
  void op_bne(AddressingMode);
  void op_bpl(AddressingMode);
  void op_bvc(AddressingMode);
  void op_bvs(AddressingMode);
  // status flag changes
  void op_clc(AddressingMode);
  void op_cld(AddressingMode);
  void op_cli(AddressingMode);
  void op_clv(AddressingMode);
  void op_sec(AddressingMode);
  void op_sed(AddressingMode);
  void op_sei(AddressingMode);
  // system functions
  void op_brk(AddressingMode);
  void op_nop(AddressingMode);
  void op_rti(AddressingMode);

  // register utils
  void set_register_a(uint8_t value);
  void set_register_x(uint8_t value);
  void set_register_y(uint8_t value);
  void add_to_register_a(uint8_t value);

  // flag utils
  void set_flag(uint8_t mask, bool condition);
  void update_zero_and_negative_flags(uint8_t value);
  bool get_flag(uint8_t mask);

  // stack utils
  void stack_push(uint8_t data);
  void stack_push_u16(uint16_t data);
  uint8_t stack_pop();
  uint16_t stack_pop_u16();

  // op code helpers
  void branch(AddressingMode mode, bool condition);
  void compare(AddressingMode mode, uint8_t compare_with);

  // additional utils
  uint8_t fetch_next_byte();
  uint16_t get_addr(AddressingMode);
  void write_to_reg_a_or_mem(AddressingMode, std::optional<uint16_t> addr,
                             uint8_t value);
};
