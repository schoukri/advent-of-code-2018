package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func prepareGen(gen string) (string, int) {
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
func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	part := flag.Int64("part", 1, "The part of the puzzle to run.")
	flag.Parse()

	numGenerations := int64(20)
	if *part == 2 {
		numGenerations = 50000000000
	}

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
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

		if g%1000000 == 0 {
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
		//fmt.Printf("NEW: [%02d]%s\n", g, gen)
		//	gen = newGen
	}

	sum := 0
	for i, char := range gen {
		if char == '#' {
			sum += (i - offset)
		}
	}

	fmt.Printf("part %d: %d\n", *part, sum)
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

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
