package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	X, Y   int
	VX, VY int
}

type Grid map[int]map[int]bool

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	points := make([]*Point, 0)

	for _, line := range lines {

		re := regexp.MustCompile(`^position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>`)

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			log.Fatalf("cannot parse line: %s", line)
		}

		point := &Point{
			X:  mustParseInt(matches[1]),
			Y:  mustParseInt(matches[2]),
			VX: mustParseInt(matches[3]),
			VY: mustParseInt(matches[4]),
		}

		points = append(points, point)
	}

CLOCK:
	for clock := 1; ; clock++ {

		var (
			maxX = -math.MaxInt32
			maxY = -math.MaxInt32
			minX = math.MaxInt32
			minY = math.MaxInt32
		)

		// The grid is a "sparse" matrix (it will only contain our actual points and no empty cells)
		grid := make(Grid)

		// store the presence of each point in the grid by their updated X,Y coordinates
		// (keep track of the min/max X,Y coordinates so we know the grid size)
		for _, p := range points {
			// set the new location
			p.X += p.VX
			p.Y += p.VY
			if _, ok := grid[p.X]; !ok {
				grid[p.X] = make(map[int]bool)
			}
			grid[p.X][p.Y] = true

			if p.X > maxX {
				maxX = p.X
			}

			if p.Y > maxY {
				maxY = p.Y
			}

			if p.X < minX {
				minX = p.X
			}

			if p.Y < minY {
				minY = p.Y
			}
		}

		// Check each point to determine if it is a "single" point.
		// (i.e., the point does not have a neighbor in any of its 8 adjacent cells).
		// If the grid has a single point, then it can't be part of a valid letter.
		for _, p := range points {
			singleFound := true
		POINT:
			for x := p.X - 1; x <= p.X+1; x++ {
				for y := p.Y - 1; y <= p.Y+1; y++ {
					if x == p.X && y == p.Y {
						if !grid[x][y] {
							log.Fatalf("expected point in grid at X=%d, Y=%d\n", x, y)
						}
					} else if grid[x][y] {
						singleFound = false
						break POINT
					}
				}
			}
			if singleFound {
				continue CLOCK
			}
		}

		fmt.Println("part 1:")

		// print grid
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				if _, ok := grid[x]; !ok {
					fmt.Print(" ")
					continue
				}
				if grid[x][y] {
					fmt.Print("#")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}

		fmt.Printf("part 2: %d\n", clock)

		break CLOCK

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

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
