package main

import (
	"flag"
	"fmt"
	"log"
)

// Each fuel cell has a coordinate ranging from 1 to 300 in both the X (horizontal) and Y (vertical) direction.
// In X,Y notation, the top-left cell is 1,1, and the top-right cell is 300,1.

var (
	minX = 1
	maxX = 300
	minY = 1
	maxY = 300
)

// The power level in a given fuel cell can be found through the following process:
func Power(x, y, gridSerial int) int {

	// Find the fuel cell's rack ID, which is its X coordinate plus 10.
	rackID := x + 10

	// Begin with a power level of the rack ID times the Y coordinate.
	power := rackID * y

	// Increase the power level by the value of the grid serial number (your puzzle input).
	power += gridSerial

	// Set the power level to itself multiplied by the rack ID.
	power *= rackID

	// Keep only the hundreds digit of the power level (so 12345 becomes 3; numbers with no hundreds digit become 0).
	if power < 100 {
		power = 0
	}

	if power >= 1000 {
		power %= 1000
	}

	power /= 100

	// Subtract 5 from the power level.
	return power - 5
}

type Grid map[int]map[int]int

func main() {

	gridSerial := flag.Int("serial", 2694, "the grid serial number")
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

	var part1X, part1Y, part1Power int
	var part2X, part2Y, part2Power, part2Size int

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

			// size can range from 0 - 299
			// CX and CY can range from 1 - 300
			for size := minSize; size <= maxSize; size++ {
				// calculate total power for the row Y (for each cell of X)
				for x := CX; x <= CX+size; x++ {
					if power, ok := grid[x][CY+size]; ok {
						sizePower[size] += power
					} else {
						log.Fatalf("no X grid coord at %d,%d (y = CY=%d + size=%d)\n", x, CY+size, CY, size)
					}
				}
				// calculate total power for the column X (for each cell of Y)
				// (skip final "corner" cell because it's already been included in the row loop above)
				for y := CY; y < CY+size; y++ {
					if power, ok := grid[CX+size][y]; ok {
						sizePower[size] += power
					} else {
						log.Fatalf("no Y grid coord at %d,%d (x = CX=%d + size=%d)\n", CX+size, y, CX, size)
					}
				}

				// // add the previous size square's total power to our new row and columnn
				if size > 0 {
					sizePower[size] += sizePower[size-1]
				}

				if CX == 90 && CY == 269 && size == 15 && *gridSerial == 18 {
					fmt.Printf("gridSerial 18 example: 90,269,16 (expected power=113, actual power=%d)\n", sizePower[size])
				} else if CX == 232 && CY == 251 && size == 11 && *gridSerial == 42 {
					fmt.Printf("gridSerial 42 example: 232,251,12 (expected power=119, actual power=%d)\n", sizePower[size])
				}

				if sizePower[size] > part1Power && size == 2 {
					part1X = CX
					part1Y = CY
					part1Power = sizePower[size]
				}
				if sizePower[size] > part2Power {
					part2X = CX
					part2Y = CY
					part2Size = size
					part2Power = sizePower[size]
				}

			}
		}
	}

	fmt.Printf("part 1: %d,%d (%d)\n", part1X, part1Y, part1Power)
	fmt.Printf("part 2: %d,%d,%d (%d)\n", part2X, part2Y, part2Size+1, part2Power)
}
