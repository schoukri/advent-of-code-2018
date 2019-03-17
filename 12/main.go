package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	bitset "github.com/tmthrgd/go-bitset"
)

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	part := flag.Int64("part", 1, "The part of the puzzle to run.")
	naive := flag.Bool("naive", false, "Use the naive (slow) strategy.")
	flag.Parse()

	numGenerations := int64(20)
	if *part == 2 {
		numGenerations = 50000000000
	}

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	var sum int64
	if *naive {
		sum = NaiveStrategy(lines, numGenerations)
	} else {
		sum = FastStrategy(lines, numGenerations)
	}

	fmt.Printf("part %d: %d\n", *part, sum)

}

func FastStrategy(lines []string, numGenerations int64) int64 {

	line := lines[0][15:]
	bits := bitset.New(uint(len(line)))
	for i, char := range line {
		if char == '#' {
			bits.Set(uint(i))
		}
	}

	rules := make([]bitset.Bitset, 0)
	for _, line := range lines[2:] {
		if line[2] == line[9] {
			continue
		}
		rule := bitset.New(8)
		for i, char := range line[0:5] {
			if char == '#' {
				rule.Set(uint(i))
			}
		}
		rules = append(rules, rule)
	}

	offset := int64(0)
	totalOffset := int64(0)
	fastForwardMoves := int64(0)

	lastSig, lastStart := Signature(bits)

	for g := int64(1); g <= numGenerations; g++ {

		bits, offset = grow(bits)
		totalOffset += offset

		flipBits := make([]uint, 0)
		for i := uint(2); i < bits.Len()-3; i++ {
			part := bits.CloneRange(i-2, i+3)

			for _, rule := range rules {
				if part.Equal(rule) {
					flipBits = append(flipBits, i)
					break
				}
			}
		}

		for _, bit := range flipBits {
			bits.Invert(bit)
		}

		sig, start := Signature(bits)

		// fmt.Printf("SIG: g=%d, o=%d: %s\n", g, totalOffset, sig)

		if sig == lastSig {
			fmt.Printf("REPEAT: gen=%d, lastStart=%d, start=%d: %s\n", g, lastStart, start, sig)
			diffStart := start - lastStart
			fastForwardMoves = (numGenerations - g) * diffStart
			break
		}

		lastSig = sig
		lastStart = start
	}

	// fmt.Printf("FINAL: offset=%d: %s\n", totalOffset, asString(bits))

	sum := int64(0)
	for i := uint(0); i < bits.Len(); i++ {
		if bits.IsSet(i) {
			sum += (int64(i) - totalOffset) + fastForwardMoves
		}
	}

	return sum

}

func Signature(bits bitset.Bitset) (string, int64) {
	onBits := make([]uint, 0)
	for i := uint(0); i < bits.Len(); i++ {
		if bits.IsSet(i) {
			onBits = append(onBits, i)
		}
	}

	start := onBits[0]
	end := onBits[len(onBits)-1]

	b2 := bitset.New((end - start) + 1)
	for _, bit := range onBits {
		b2.Set(bit - start)
	}

	return asString(b2), int64(start)
}

func asString(bits bitset.Bitset) string {

	offset := 0
	size := int(bits.Len()) + bits.ByteLen() + 10
	result := make([]string, size)
	for i := 0; i < int(bits.Len()); i++ {
		if bits.IsSet(uint(i)) {
			result[i+offset] = "1"
		} else {
			result[i+offset] = "0"
		}
		if (i+1)%8 == 0 {
			offset += 1
			result[i+offset] = " "
		}

	}

	return "[" + strings.Join(result, "") + "]"

}

func grow(bits bitset.Bitset) (bitset.Bitset, int64) {

	offset := int64(0)

	if !bits.IsRangeClear(0, 5) {
		bitsR := bitset.New(uint(bits.Len() + 8))
		bitsR.ShiftRight(bits, 8)
		bits = bitsR
		offset = 8
	}

	if !bits.IsRangeClear(bits.Len()-5, bits.Len()) {
		bitsL := bitset.New(uint(bits.Len() + 8))
		bitsL.ShiftLeft(bits, 0)
		bits = bitsL
	}

	return bits, offset
}

func NaiveStrategy(lines []string, numGenerations int64) int64 {

	prepareGen := func(gen string) (string, int) {
		var lenPrefix int
		firstIndex := strings.Index(gen, "#")
		if firstIndex < 0 {
			return gen, len(gen)
		}
		if firstIndex < 5 {
			lenPrefix = 5 - firstIndex
			gen = strings.Repeat(".", lenPrefix) + gen
		}

		lastIndex := strings.LastIndex(gen, "#")
		if lastIndex > len(gen)-6 {
			lenSuffix := lastIndex - (len(gen) - 6)
			gen += strings.Repeat(".", lenSuffix)
		}

		return gen, lenPrefix
	}

	gen, offset := prepareGen(lines[0][15:])

	rules := make(map[string]string)
	for _, line := range lines[2:] {
		rules[line[0:5]] = string(line[9])
	}

	var newOffset int
	for g := int64(1); g <= numGenerations; g++ {
		// we only start evaluating at char index 2 (the 3rd char)
		// start the newGen with the first two chars already set
		newGen := ".."

		if g%1000 == 0 {
			fmt.Printf("%d -- %v\n", g, time.Now())
		}

		for i := 2; i < len(gen)-3; i++ {
			key := gen[i-2 : i+3]
			if _, ok := rules[key]; ok {
				//fmt.Printf("MATCHED:  i=%2d, key=[%s] (result=%s)\n", i, key, res)
				newGen += rules[key]
			} else {
				//fmt.Printf("NO MATCH: i=%2d, key=[%s]\n", i, key)
				newGen += "."
			}
		}

		gen, newOffset = prepareGen(newGen)
		offset += newOffset
		// fmt.Printf("NEW: [%02d]%s\n", g, gen)
		//	gen = newGen
	}

	sum := 0
	for i, char := range gen {
		if char == '#' {
			sum += (i - offset)
		}
	}

	return int64(sum)
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
