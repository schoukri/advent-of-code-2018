package main

import (
	"flag"
	"fmt"
	"log"
)

// Each fuel cell has a coordinate ranging from 1 to 300 in both the X (horizontal) and Y (vertical) direction.
// In X,Y notation, the top-left cell is 1,1, and the top-right cell is 300,1.

// The interface lets you select any 3x3 square of fuel cells. To increase your chances of getting to your destination,
//  you decide to choose the 3x3 square with the largest total power.

var (
	minX = 1
	maxX = 300
	minY = 1
	maxY = 300
)

// The power level in a given fuel cell can be found through the following process:

// Find the fuel cell's rack ID, which is its X coordinate plus 10.
// Begin with a power level of the rack ID times the Y coordinate.
// Increase the power level by the value of the grid serial number (your puzzle input).
// Set the power level to itself multiplied by the rack ID.
// Keep only the hundreds digit of the power level (so 12345 becomes 3; numbers with no hundreds digit become 0).
// Subtract 5 from the power level.

func Power(x, y, gridSerial int) int {

	rackID := x + 10

	power := rackID * y
	power += gridSerial
	power *= rackID

	if power < 100 {
		return -5
	}

	if power >= 10000 {
		power %= 1000
	}

	for power >= 1000 {
		power /= 10
	}

	power /= 100

	return power - 5
}

type Grid map[int]map[int]int

func main() {

	gridSerial := flag.Int("serial", 2694, "the grid serial number")
	// part := flag.Int("part", 1, "the part of the challenge to run")
	flag.Parse()

	grid := make(Grid)

	totalPower := 0
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if _, ok := grid[x]; !ok {
				grid[x] = make(map[int]int)
			}
			grid[x][y] = Power(x, y, *gridSerial)
			totalPower += grid[x][y]
		}
	}

	if power := Power(122, 79, 57); power != -5 {
		log.Fatalf("Power(122, 79, 57): expected -5 (got %d)\n", power)
	}
	if power := Power(217, 196, 39); power != 0 {
		log.Fatalf("Power(217,196, 39): expected 0 (got %d)\n", power)
	}
	if power := Power(101, 153, 71); power != 4 {
		log.Fatalf("Power(101,153, 71): expected 4 (got %d)\n", power)
	}

	minSize := 0
	// maxSize := 2

	var cellX, cellY, cellSize, maxPower int

	for CX := minX; CX <= maxX; CX++ {
		for CY := minY; CY <= maxY; CY++ {
			// calculate the largest size square we can fit in the grid,
			// starting with the the current coord in the top-left corner
			// (size values are an offset and therefore start at 0)
			maxSizeX := maxX - CX
			maxSizeY := maxY - CY

			// // pick the smallest dimension for the size of the square
			maxSize := maxSizeX
			if maxSizeY < maxSizeX {
				maxSize = maxSizeY
			}

			sizePower := make(map[int]int)
			sizeCount := make(map[int]int)
			sumSizePower := make(map[int]int)
			sumSizeCount := make(map[int]int)

			for size := minSize; size <= maxSize; size++ {
				// calculate total power for the row Y (for each cell of X)
				for x := CX; x <= CX+size; x++ {
					if power, ok := grid[x][CY+size]; ok {
						sizePower[size] += power
						sizeCount[size]++

					} else {
						log.Fatalf("no X grid coord at %d,%d (y = CY=%d + size=%d)\n", x, CY+size, CY, size)
					}
				}
				// calculate total power for the column X (for each cell of Y)
				// (skip final "corner" cell because it's already been included in the row loop above)
				for y := CY; y < CY+size; y++ {
					if power, ok := grid[CX+size][y]; ok {
						sizePower[size] += power
						sizeCount[size]++
					} else {
						log.Fatalf("no Y grid coord at %d,%d (x = CX=%d + size=%d)\n", CX+size, y, CX, size)
					}
				}

				// add the new row and column power to the square's total power
				for includeSize := minSize; includeSize <= size; includeSize++ {
					sumSizePower[size] += sizePower[includeSize]
					sumSizeCount[size] += sizeCount[includeSize]
				}

			}

			thisMaxPower := 0
			for size, power := range sumSizePower {
				if power > maxPower {
					maxPower = power
					cellX = CX
					cellY = CY
					cellSize = size
				}
				if power > thisMaxPower {
					thisMaxPower = power
				}

			}

			for size := minSize; size <= maxSize; size++ {
				fmt.Printf("CELL\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%v\n",
					CX, CY, size+1,
					sizeCount[size], sumSizeCount[size],
					sizePower[size], sumSizePower[size],
					thisMaxPower, maxPower, thisMaxPower == maxPower,
				)
			}
		}
	}

	fmt.Printf("grid[1][1] power: %d\n", grid[1][1])
	fmt.Printf("totalPower: %d\n", totalPower)

	//fmt.Printf("part 1: %d,%d (%d)\n", cellX, cellY, maxPower)
	fmt.Printf("part 2: %d,%d,%d (%d)\n", cellX, cellY, cellSize+1, maxPower)
}
