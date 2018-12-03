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

	"gonum.org/v1/gonum/mat"
)

var claimRegexp = regexp.MustCompile(`^\#(\d+) \@ (\d+),(\d+): (\d+)x(\d+)$`)

type Claim struct {
	ID   int
	X, Y int
	W, H int
	M    *mat.Dense
}

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	maxWidth := 0
	maxHeight := 0
	claims := make([]Claim, len(lines))

	// parse the lines into Claims
	// and determine matrix dimensions
	for i, line := range lines {

		c := mustParseLine(line)

		if c.X+c.W > maxWidth {
			maxWidth = c.X + c.W
		}

		if c.Y+c.H > maxHeight {
			maxHeight = c.Y + c.H
		}

		claims[i] = c
	}

	// the "all" matrix will hold a presence indicator for all claims
	// (each claim will add 1 to every cell it occupies)
	all := mat.NewDense(maxWidth, maxHeight, nil)

	// create a matrix for each claim, setting the cells they occupy
	for idx, c := range claims {

		m := mat.NewDense(maxWidth, maxHeight, nil)
		for i := c.X; i < c.X+c.W; i++ {
			for j := c.Y; j < c.Y+c.H; j++ {
				m.Set(i, j, 1)
			}
		}

		claims[idx].M = m

		// add this claim matrix to the "all" matrix
		all.Add(all, m)

	}

	rows, cols := all.Dims()

	// PART 1: determine the number of cells that are occupied by more than 1 claim
	overlap := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if all.At(i, j) > 1 {
				overlap++
			}
		}
	}
	fmt.Printf("part 1: %d\n", overlap)

	// PART 2: find a claim that is NOT overlapped by any other claim
CLAIM:
	for _, c := range claims {

		product := mat.NewDense(maxWidth, maxHeight, nil)

		// multiply the corresponding cells of this claim matrix and the all matrix
		product.MulElem(c.M, all)

		// the product matrix will equal the claim matrix *if* the claim is NOT overlapped
		if mat.Equal(c.M, product) {
			fmt.Printf("part 2: %d\n", c.ID)
			break CLAIM
		}

	}

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

func mustParseLine(line string) Claim {

	matches := claimRegexp.FindStringSubmatch(line)
	if matches == nil {
		log.Fatalf("line not be parsed: %s", line)
	}

	claim := Claim{
		ID: mustParseInt(matches[1]),
		X:  mustParseInt(matches[2]),
		Y:  mustParseInt(matches[3]),
		W:  mustParseInt(matches[4]),
		H:  mustParseInt(matches[5]),
	}

	return claim
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
