package vm

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {
	if len(memory) > 256 {
		panic("Out of memory")
	}

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	instructions := [9]func(byte, []byte, *[3]byte) byte{
		Load: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg, addr := memory[pc+1], memory[pc+2]
			registers[reg] = memory[addr]
			return pc + 3
		},
		Store: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg, addr := memory[pc+1], memory[pc+2]
			if addr > 7 {
				panic("Illegal write to read-only memory!")
			}
			memory[addr] = registers[reg]
			return pc + 3
		},
		Add: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg1, reg2 := memory[pc+1], memory[pc+2]
			registers[reg1] += registers[reg2]
			return pc + 3
		},
		Sub: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg1, reg2 := memory[pc+1], memory[pc+2]
			registers[reg1] -= registers[reg2]
			return pc + 3
		},
		Addi: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg, val := memory[pc+1], memory[pc+2]
			registers[reg] += val
			return pc + 3
		},
		Subi: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg, val := memory[pc+1], memory[pc+2]
			registers[reg] -= val
			return pc + 3
		},
		Jump: func(pc byte, memory []byte, registers *[3]byte) byte {
			return memory[pc+1]
		},
		Beqz: func(pc byte, memory []byte, registers *[3]byte) byte {
			reg, offset := memory[pc+1], memory[pc+2]
			// move PC by offset conditional on value in reg
			if registers[reg] == 0 {
				pc += offset
			}
			return pc + 3
		},
	}
	var pc, op byte
	// Keep looping, like a physical computer's clock
	for {
		pc = registers[0]
		op = memory[pc]
		if op == Halt {
			return
		}
		registers[0] = instructions[op](pc, memory, &registers)
	}
}
