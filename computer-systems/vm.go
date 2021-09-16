package vm

import "fmt"

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

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {
		position := registers[0]
		op := memory[position]

		switch op {
		case Load:
			// increment PC
			registers[0] += 3
			reg := memory[position+1]
			addr := memory[position+2]
			// load data at dataAddr into register reg
			registers[reg] = memory[addr]
		case Store:
			registers[0] += 3
			reg := memory[position+1]
			addr := memory[position+2]
			// load data at dataAddr into register reg
			memory[addr] = registers[reg]
		case Add:
			registers[0] += 3
			reg1 := memory[position+1]
			reg2 := memory[position+2]
			// add register values, store in reg1
			registers[reg1] += registers[reg2]
		case Sub:
			registers[0] += 3
			reg1 := memory[position+1]
			reg2 := memory[position+2]
			// add register values, store in reg1
			registers[reg1] -= registers[reg2]
		case Addi:
			registers[0] += 3
			reg := memory[position+1]
			val := memory[position+2]
			// add val to value stored in register
			registers[reg] += val
		case Subi:
			registers[0] += 3
			reg := memory[position+1]
			val := memory[position+2]
			// subtract val from value stored in register
			registers[reg] -= val
		case Jump:
			jumpTo := memory[position+1]
			// set PC to addr specified in arg
			registers[0] = jumpTo
		case Beqz:
			registers[0] += 3
			reg := memory[position+1]
			offset := memory[position+2]
			// move PC by offset conditional on value in reg
			if registers[reg] == 0 {
				registers[0] += offset
			}
		case Halt:
			return
		default:
			panic(fmt.Errorf("Unknown opcode: %b", op))
		}
	}
}
