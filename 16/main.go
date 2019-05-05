package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Register [4]int

type Instruction struct {
	OpcodeNum int
	A, B, C   int
}
type Sample struct {
	Before      Register
	Instruction *Instruction
	After       Register
}

type Opcode func(r Register, a, b, c int) Register

var (
	registerRegexp    = regexp.MustCompile(`^(Before|After):\s+\[(\d+), (\d+), (\d+), (\d+)\]$`)
	instructionRegexp = regexp.MustCompile(`^(\d+) (\d+) (\d+) (\d+)$`)
)

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	samples := make([]*Sample, 0)
	instructions := make([]*Instruction, 0)
	for i := 0; i < len(lines); i++ {

		if beforeRegister := ParseRegister(lines[i], "Before"); beforeRegister != nil {

			i++
			instruction := ParseInstruction(lines[i])
			if instruction == nil {
				log.Fatalln("line did not match instructions: " + lines[i])
			}

			i++
			afterRegister := ParseRegister(lines[i], "After")
			if afterRegister == nil {
				log.Fatalln("line did not match after register: " + lines[i])
			}

			sample := &Sample{
				Before:      *beforeRegister,
				Instruction: instruction,
				After:       *afterRegister,
			}
			samples = append(samples, sample)
		} else if instruction := ParseInstruction(lines[i]); instruction != nil {
			instructions = append(instructions, instruction)
		}
	}

	// an array of all opcode functions
	// (we will need to map the index position of these opcodes to the correct OpcodeNum in the samples)
	opcodes := [16]Opcode{
		addr, addi, muli, mulr, bani, banr, bori, borr,
		seti, setr, gtir, gtri, gtrr, eqir, eqri, eqrr,
	}

	// map of all OpcodeNum/OpcodeIndex combinations that were valid
	validOpcodes := make(map[int]map[int]int)

	validThreeOrMore := 0
	for _, s := range samples {
		//fmt.Printf("sample[%d]: %+v\n", i, s)
		validCount := 0
		for opcodeIndex, opcode := range opcodes {
			if IsEqual(opcode(s.Before, s.Instruction.A, s.Instruction.B, s.Instruction.C), s.After) {
				if _, ok := validOpcodes[opcodeIndex]; !ok {
					validOpcodes[opcodeIndex] = make(map[int]int)
				}
				validOpcodes[opcodeIndex][s.Instruction.OpcodeNum]++
				validCount++
			}
		}
		if validCount >= 3 {
			validThreeOrMore++
		}
	}

	fmt.Printf("part 1: %d\n", validThreeOrMore)

	// map the opcode num to the opcode index
	opcodeMap := make(map[int]int)
	for found := 0; found < 16; {
		for opcodeIndex, opcodeNums := range validOpcodes {
			notMapped := make([]int, 0)
			for opcodeNum := range opcodeNums {
				// fmt.Printf("Opcode: index=%d, num=%d\n", opcodeIndex, opcodeNum)
				if _, ok := opcodeMap[opcodeNum]; !ok {
					notMapped = append(notMapped, opcodeNum)
				}
			}
			// find the opcodeIndex that has exactly 1 corresponding opcodeNum that is not yet mapped
			if len(notMapped) == 1 {
				opcodeMap[notMapped[0]] = opcodeIndex
				found++
			}
		}
	}

	// run all the test instructions
	var register Register
	for _, instr := range instructions {
		opcode := opcodes[opcodeMap[instr.OpcodeNum]]
		register = opcode(register, instr.A, instr.B, instr.C)
	}

	fmt.Printf("part 2: %d\n", register[0])

}

func ParseRegister(line string, label string) *Register {
	matches := registerRegexp.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}
	if matches[1] != label {
		return nil
	}

	var register Register
	for i := range register {
		register[i] = mustParseInt(matches[i+2])
	}
	return &register
}

func ParseInstruction(line string) *Instruction {
	matches := instructionRegexp.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}
	return &Instruction{
		OpcodeNum: mustParseInt(matches[1]),
		A:         mustParseInt(matches[2]),
		B:         mustParseInt(matches[3]),
		C:         mustParseInt(matches[4]),
	}
}

func IsEqual(r1, r2 Register) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i := range r1 {
		if r1[i] != r2[i] {
			return false
		}
	}
	return true
}

// Addition:
// addr (add register) stores into register C the result of adding register A and register B.
func addr(r Register, a, b, c int) Register {
	r[c] = r[a] + r[b]
	return r
}

// addi (add immediate) stores into register C the result of adding register A and value B.
func addi(r Register, a, b, c int) Register {
	r[c] = r[a] + b
	return r
}

// Multiplication:
// mulr (multiply register) stores into register C the result of multiplying register A and register B.
func mulr(r Register, a, b, c int) Register {
	r[c] = r[a] * r[b]
	return r
}

// muli (multiply immediate) stores into register C the result of multiplying register A and value B.
func muli(r Register, a, b, c int) Register {
	r[c] = r[a] * b
	return r
}

// Bitwise AND:
// banr (bitwise AND register) stores into register C the result of the bitwise AND of register A and register B.
func banr(r Register, a, b, c int) Register {
	r[c] = r[a] & r[b]
	return r
}

// bani (bitwise AND immediate) stores into register C the result of the bitwise AND of register A and value B.
func bani(r Register, a, b, c int) Register {
	r[c] = r[a] & b
	return r
}

// Bitwise OR:
// borr (bitwise OR register) stores into register C the result of the bitwise OR of register A and register B.
func borr(r Register, a, b, c int) Register {
	r[c] = r[a] | r[b]
	return r
}

// bori (bitwise OR immediate) stores into register C the result of the bitwise OR of register A and value B.
func bori(r Register, a, b, c int) Register {
	r[c] = r[a] | b
	return r
}

// Assignment:
// setr (set register) copies the contents of register A into register C. (Input B is ignored.)
func setr(r Register, a, b, c int) Register {
	r[c] = r[a]
	return r
}

// seti (set immediate) stores value A into register C. (Input B is ignored.)
func seti(r Register, a, b, c int) Register {
	r[c] = a
	return r
}

// Greater-than testing:
// gtir (greater-than immediate/register) sets register C to 1 if value A is greater than register B. Otherwise, register C is set to 0.
func gtir(r Register, a, b, c int) Register {
	if a > r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

// gtri (greater-than register/immediate) sets register C to 1 if register A is greater than value B. Otherwise, register C is set to 0.
func gtri(r Register, a, b, c int) Register {
	if r[a] > b {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

// gtrr (greater-than register/register) sets register C to 1 if register A is greater than register B. Otherwise, register C is set to 0.
func gtrr(r Register, a, b, c int) Register {
	if r[a] > r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

// Equality testing:
// eqir (equal immediate/register) sets register C to 1 if value A is equal to register B. Otherwise, register C is set to 0.
func eqir(r Register, a, b, c int) Register {
	if a == r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

// eqri (equal register/immediate) sets register C to 1 if register A is equal to value B. Otherwise, register C is set to 0.
func eqri(r Register, a, b, c int) Register {
	if r[a] == b {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

// eqrr (equal register/register) sets register C to 1 if register A is equal to register B. Otherwise, register C is set to 0.
func eqrr(r Register, a, b, c int) Register {
	if r[a] == r[b] {
		r[c] = 1
	} else {
		r[c] = 0
	}
	return r
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}

func readFile(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("file path not specified")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
