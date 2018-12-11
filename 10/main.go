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

		// The grid is a "sparse" matrix (it will only contain our actual points and no empty cells)
		grid := make(Grid)

		// store the presence of each point in the grid by their updated X,Y coordinates
		for _, p := range points {
			// set the new location
			p.X += p.VX
			p.Y += p.VY
			if _, ok := grid[p.X]; !ok {
				grid[p.X] = make(map[int]bool)
			}
			grid[p.X][p.Y] = true
		}

		// Check each point in the grid to determine if it is a "single" point.
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
			// if one single point, there is no need to keep looking at the rest of the grid
			if singleFound {
				continue CLOCK
			}
		}

		// This gird does NOT have *any* single points.
		// It is *highly* likely it is the grid with the message.
		// Time to print it out and see.
		fmt.Println("part 1:")
		grid.Print()

		fmt.Printf("part 2: %d\n", clock)

		break CLOCK

	}
}

func (g Grid) Print() {

	// figure out the grid coordinates
	var (
		maxX = -math.MaxInt32
		maxY = -math.MaxInt32
		minX = math.MaxInt32
		minY = math.MaxInt32
	)

	for x, rows := range g {
		for y := range rows {
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
		}
	}

	// print the grid
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if g[x][y] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
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
