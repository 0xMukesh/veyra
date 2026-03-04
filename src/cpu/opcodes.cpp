#include "cpu.hpp"

const std::array<CPU::OpCode, 256> &CPU::GetOpTable() {
    static const std::array<OpCode, 256> table = [] {
        std::array<OpCode, 256> t = {};

        // load instructions
        // -- LDA
        t[0xa9] = {&CPU::op_lda, AddressingMode::Immediate, 2};
        t[0xa5] = {&CPU::op_lda, AddressingMode::ZeroPage, 2};
        t[0xb5] = {&CPU::op_lda, AddressingMode::ZeroPage_X, 2};
        t[0xad] = {&CPU::op_lda, AddressingMode::Absolute, 3};
        t[0xbd] = {&CPU::op_lda, AddressingMode::Absolute_X, 3};
        t[0xb9] = {&CPU::op_lda, AddressingMode::Absolute_Y, 3};
        t[0xa1] = {&CPU::op_lda, AddressingMode::Indirect_X, 2};
        t[0xb1] = {&CPU::op_lda, AddressingMode::Indirect_Y, 2};
        // -- LDX
        t[0xa2] = {&CPU::op_ldx, AddressingMode::Immediate, 2};
        t[0xa6] = {&CPU::op_ldx, AddressingMode::ZeroPage, 2};
        t[0xb6] = {&CPU::op_ldx, AddressingMode::ZeroPage_Y, 2};
        t[0xae] = {&CPU::op_ldx, AddressingMode::Absolute, 3};
        t[0xbe] = {&CPU::op_ldx, AddressingMode::Absolute_Y, 3};
        // -- LDY
        t[0xa0] = {&CPU::op_ldy, AddressingMode::Immediate, 2};
        t[0xa4] = {&CPU::op_ldy, AddressingMode::ZeroPage, 2};
        t[0xb4] = {&CPU::op_ldy, AddressingMode::ZeroPage_X, 2};
        t[0xac] = {&CPU::op_ldy, AddressingMode::Absolute, 3};
        t[0xbc] = {&CPU::op_ldy, AddressingMode::Absolute_X, 3};

        // store instructions
        // -- STA
        t[0x85] = {&CPU::op_sta, AddressingMode::ZeroPage, 2};
        t[0x95] = {&CPU::op_sta, AddressingMode::ZeroPage_X, 2};
        t[0x8d] = {&CPU::op_sta, AddressingMode::Absolute, 3};
        t[0x9d] = {&CPU::op_sta, AddressingMode::Absolute_X, 3};
        t[0x99] = {&CPU::op_sta, AddressingMode::Absolute_Y, 3};
        t[0x81] = {&CPU::op_sta, AddressingMode::Indirect_X, 2};
        t[0x91] = {&CPU::op_sta, AddressingMode::Indirect_Y, 2};
        // -- STX
        t[0x86] = {&CPU::op_stx, AddressingMode::ZeroPage, 2};
        t[0x96] = {&CPU::op_stx, AddressingMode::ZeroPage_Y, 2};
        t[0x8e] = {&CPU::op_stx, AddressingMode::Absolute, 3};
        // -- STY
        t[0x84] = {&CPU::op_sty, AddressingMode::ZeroPage, 2};
        t[0x94] = {&CPU::op_sty, AddressingMode::ZeroPage_X, 2};
        t[0x8c] = {&CPU::op_sty, AddressingMode::Absolute, 3};

        // register transfer instructions
        // -- TAX
        t[0xaa] = {&CPU::op_tax, AddressingMode::Implied, 1};
        // -- TAY
        t[0xa8] = {&CPU::op_tay, AddressingMode::Implied, 1};
        // -- TXA
        t[0x8a] = {&CPU::op_txa, AddressingMode::Implied, 1};
        // -- TYA
        t[0x98] = {&CPU::op_tya, AddressingMode::Implied, 1};

        // stack operations
        // -- TSX
        t[0xba] = {&CPU::op_tsx, AddressingMode::Implied, 1};
        // -- TXS
        t[0x9a] = {&CPU::op_txs, AddressingMode::Implied, 1};
        // -- PHA
        t[0x48] = {&CPU::op_pha, AddressingMode::Implied, 1};
        // -- PHP
        t[0x08] = {&CPU::op_php, AddressingMode::Implied, 1};
        // -- PLA
        t[0x68] = {&CPU::op_pla, AddressingMode::Implied, 1};
        // -- PLP
        t[0x28] = {&CPU::op_plp, AddressingMode::Implied, 1};

        // logical operations
        // -- AND
        t[0x29] = {&CPU::op_and, AddressingMode::Immediate, 2};
        t[0x25] = {&CPU::op_and, AddressingMode::ZeroPage, 2};
        t[0x35] = {&CPU::op_and, AddressingMode::ZeroPage_X, 2};
        t[0x2d] = {&CPU::op_and, AddressingMode::Absolute, 3};
        t[0x3d] = {&CPU::op_and, AddressingMode::Absolute_X, 3};
        t[0x39] = {&CPU::op_and, AddressingMode::Absolute_Y, 3};
        t[0x21] = {&CPU::op_and, AddressingMode::Indirect_X, 2};
        t[0x31] = {&CPU::op_and, AddressingMode::Indirect_Y, 2};
        // -- EOR
        t[0x49] = {&CPU::op_eor, AddressingMode::Immediate, 2};
        t[0x45] = {&CPU::op_eor, AddressingMode::ZeroPage, 2};
        t[0x55] = {&CPU::op_eor, AddressingMode::ZeroPage_X, 2};
        t[0x4d] = {&CPU::op_eor, AddressingMode::Absolute, 3};
        t[0x5d] = {&CPU::op_eor, AddressingMode::Absolute_X, 3};
        t[0x59] = {&CPU::op_eor, AddressingMode::Absolute_Y, 3};
        t[0x41] = {&CPU::op_eor, AddressingMode::Indirect_X, 2};
        t[0x51] = {&CPU::op_eor, AddressingMode::Indirect_Y, 2};
        // -- ORA
        t[0x09] = {&CPU::op_ora, AddressingMode::Immediate, 2};
        t[0x05] = {&CPU::op_ora, AddressingMode::ZeroPage, 2};
        t[0x15] = {&CPU::op_ora, AddressingMode::ZeroPage_X, 2};
        t[0x0d] = {&CPU::op_ora, AddressingMode::Absolute, 3};
        t[0x1d] = {&CPU::op_ora, AddressingMode::Absolute_X, 3};
        t[0x19] = {&CPU::op_ora, AddressingMode::Absolute_Y, 3};
        t[0x01] = {&CPU::op_ora, AddressingMode::Indirect_X, 2};
        t[0x11] = {&CPU::op_ora, AddressingMode::Indirect_Y, 2};
        // -- BIT
        t[0x24] = {&CPU::op_bit, AddressingMode::ZeroPage, 2};
        t[0x2c] = {&CPU::op_bit, AddressingMode::Absolute, 3};

        // arithemtic
        // -- ADC
        t[0x69] = {&CPU::op_adc, AddressingMode::Immediate, 2};
        t[0x65] = {&CPU::op_adc, AddressingMode::ZeroPage, 2};
        t[0x75] = {&CPU::op_adc, AddressingMode::ZeroPage_X, 2};
        t[0x6d] = {&CPU::op_adc, AddressingMode::Absolute, 3};
        t[0x7d] = {&CPU::op_adc, AddressingMode::Absolute_X, 3};
        t[0x79] = {&CPU::op_adc, AddressingMode::Absolute_Y, 3};
        t[0x61] = {&CPU::op_adc, AddressingMode::Indirect_X, 2};
        t[0x71] = {&CPU::op_adc, AddressingMode::Indirect_Y, 2};
        // -- SBC
        t[0xe9] = {&CPU::op_sbc, AddressingMode::Immediate, 2};
        t[0xe5] = {&CPU::op_sbc, AddressingMode::ZeroPage, 2};
        t[0xf5] = {&CPU::op_sbc, AddressingMode::ZeroPage_X, 2};
        t[0xed] = {&CPU::op_sbc, AddressingMode::Absolute, 3};
        t[0xfd] = {&CPU::op_sbc, AddressingMode::Absolute_X, 3};
        t[0xf9] = {&CPU::op_sbc, AddressingMode::Absolute_Y, 3};
        t[0xe1] = {&CPU::op_sbc, AddressingMode::Indirect_X, 2};
        t[0xf1] = {&CPU::op_sbc, AddressingMode::Indirect_Y, 2};
        // -- CMP
        t[0xc9] = {&CPU::op_cmp, AddressingMode::Immediate, 2};
        t[0xc5] = {&CPU::op_cmp, AddressingMode::ZeroPage, 2};
        t[0xd5] = {&CPU::op_cmp, AddressingMode::ZeroPage_X, 2};
        t[0xcd] = {&CPU::op_cmp, AddressingMode::Absolute, 3};
        t[0xdd] = {&CPU::op_cmp, AddressingMode::Absolute_X, 3};
        t[0xd9] = {&CPU::op_cmp, AddressingMode::Absolute_Y, 3};
        t[0xc1] = {&CPU::op_cmp, AddressingMode::Indirect_X, 2};
        t[0xd1] = {&CPU::op_cmp, AddressingMode::Indirect_Y, 2};
        // -- CPX
        t[0xe0] = {&CPU::op_cpx, AddressingMode::Immediate, 2};
        t[0xe4] = {&CPU::op_cpx, AddressingMode::ZeroPage, 2};
        t[0xec] = {&CPU::op_cpx, AddressingMode::Absolute, 3};
        // -- CPY
        t[0xc0] = {&CPU::op_cpy, AddressingMode::Immediate, 2};
        t[0xc4] = {&CPU::op_cpy, AddressingMode::ZeroPage, 2};
        t[0xcc] = {&CPU::op_cpy, AddressingMode::Absolute, 3};

        // increment operations
        // -- INC
        t[0xe6] = {&CPU::op_inc, AddressingMode::ZeroPage, 2};
        t[0xf6] = {&CPU::op_inc, AddressingMode::ZeroPage_X, 2};
        t[0xee] = {&CPU::op_inc, AddressingMode::Absolute, 3};
        t[0xfe] = {&CPU::op_inc, AddressingMode::Absolute_X, 3};
        // -- INX
        t[0xe8] = {&CPU::op_inx, AddressingMode::Implied, 1};
        // -- INY
        t[0xc8] = {&CPU::op_iny, AddressingMode::Implied, 1};

        // decrement operations
        // -- DEC
        t[0xc6] = {&CPU::op_dec, AddressingMode::ZeroPage, 2};
        t[0xd6] = {&CPU::op_dec, AddressingMode::ZeroPage_X, 2};
        t[0xce] = {&CPU::op_dec, AddressingMode::Absolute, 3};
        t[0xde] = {&CPU::op_dec, AddressingMode::Absolute_X, 3};
        // -- DEX
        t[0xca] = {&CPU::op_dex, AddressingMode::Implied, 1};
        // -- DEY
        t[0x88] = {&CPU::op_dey, AddressingMode::Implied, 1};

        // shifts
        // -- ASL
        t[0x0a] = {&CPU::op_asl, AddressingMode::Accumulator, 1};
        t[0x06] = {&CPU::op_asl, AddressingMode::ZeroPage, 2};
        t[0x16] = {&CPU::op_asl, AddressingMode::ZeroPage_X, 2};
        t[0x0e] = {&CPU::op_asl, AddressingMode::Absolute, 3};
        t[0x1e] = {&CPU::op_asl, AddressingMode::Absolute_X, 3};
        // -- LSR
        t[0x4a] = {&CPU::op_lsr, AddressingMode::Accumulator, 1};
        t[0x46] = {&CPU::op_lsr, AddressingMode::ZeroPage, 2};
        t[0x56] = {&CPU::op_lsr, AddressingMode::ZeroPage_X, 2};
        t[0x4e] = {&CPU::op_lsr, AddressingMode::Absolute, 3};
        t[0x5e] = {&CPU::op_lsr, AddressingMode::Absolute_X, 3};
        // -- ROL
        t[0x2a] = {&CPU::op_rol, AddressingMode::Accumulator, 1};
        t[0x26] = {&CPU::op_rol, AddressingMode::ZeroPage, 2};
        t[0x36] = {&CPU::op_rol, AddressingMode::ZeroPage_X, 2};
        t[0x2e] = {&CPU::op_rol, AddressingMode::Absolute, 3};
        t[0x3e] = {&CPU::op_rol, AddressingMode::Absolute_X, 3};
        // -- ROR
        t[0x6a] = {&CPU::op_ror, AddressingMode::Accumulator, 1};
        t[0x66] = {&CPU::op_ror, AddressingMode::ZeroPage, 2};
        t[0x76] = {&CPU::op_ror, AddressingMode::ZeroPage_X, 2};
        t[0x6e] = {&CPU::op_ror, AddressingMode::Absolute, 3};
        t[0x7e] = {&CPU::op_ror, AddressingMode::Absolute_X, 3};

        // jumps & calls
        // -- JMP
        t[0x4c] = {&CPU::op_jmp, AddressingMode::Absolute, 3};
        t[0x6c] = {&CPU::op_jmp, AddressingMode::Indirect, 3};
        // -- JSR
        t[0x20] = {&CPU::op_jsr, AddressingMode::Absolute, 3};
        // -- RTS
        t[0x60] = {&CPU::op_rts, AddressingMode::Implied, 1};

        // branches
        // -- BCC
        t[0x90] = {&CPU::op_bcc, AddressingMode::Relative, 2};
        // -- BCS
        t[0xb0] = {&CPU::op_bcs, AddressingMode::Relative, 2};
        // -- BEQ
        t[0xf0] = {&CPU::op_beq, AddressingMode::Relative, 2};
        // -- BMI
        t[0x30] = {&CPU::op_bmi, AddressingMode::Relative, 2};
        // -- BNE
        t[0xd0] = {&CPU::op_bne, AddressingMode::Relative, 2};
        // -- BPL
        t[0x10] = {&CPU::op_bpl, AddressingMode::Relative, 2};
        // -- BVC
        t[0x50] = {&CPU::op_bvc, AddressingMode::Relative, 2};
        // -- BVS
        t[0x70] = {&CPU::op_bvs, AddressingMode::Relative, 2};

        // status flag changes
        // -- CLC
        t[0x18] = {&CPU::op_clc, AddressingMode::Implied, 1};
        // -- CLD
        t[0xd8] = {&CPU::op_cld, AddressingMode::Implied, 1};
        // -- CLI
        t[0x58] = {&CPU::op_cli, AddressingMode::Implied, 1};
        // -- CLV
        t[0xb8] = {&CPU::op_clv, AddressingMode::Implied, 1};
        // -- SEC
        t[0x38] = {&CPU::op_sec, AddressingMode::Implied, 1};
        // -- SED
        t[0xf8] = {&CPU::op_sed, AddressingMode::Implied, 1};
        // -- SEI
        t[0x78] = {&CPU::op_sei, AddressingMode::Implied, 1};

        // system functions
        // -- BRK
        t[0x00] = {&CPU::op_brk, AddressingMode::Implied, 1};
        // -- NOP
        t[0xea] = {&CPU::op_nop, AddressingMode::Implied, 1};
        // -- RTI
        t[0x40] = {&CPU::op_rti, AddressingMode::Implied, 1};

        return t;
    }();

    return table;
}
