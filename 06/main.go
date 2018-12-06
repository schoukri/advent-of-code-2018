package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

type Point struct {
	ID   int
	X, Y int
}

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	re := regexp.MustCompile(`^(\d+), (\d+)$`)

	maxWidth := 0
	maxHeight := 0

	points := make([]Point, 0)
	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			log.Fatalf("cannot parse line: %s", line)
		}

		p := Point{
			ID: i,
			X:  mustParseInt(matches[1]),
			Y:  mustParseInt(matches[2]),
		}

		if p.X > maxWidth {
			maxWidth = p.X
		}

		if p.Y > maxHeight {
			maxHeight = p.Y
		}

		points = append(points, p)

	}

	pointIsInfinite := make(map[int]bool)
	areaPerPoint := make(map[int]int)
	grid := mat.NewDense(maxWidth, maxHeight, nil)
	var regionSize int

	rows, cols := grid.Dims()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {

			var sumDist int
			pointDist := make(map[int]int)
			for _, p := range points {
				dist := ManhattanDistance(i, j, p.X, p.Y)
				pointDist[p.ID] = dist
				sumDist += dist
			}

			// for part 2
			if sumDist < 10000 {
				regionSize++
			}

			var pointID int
			// get the two points with the shortest distance
			kvs := SmallestNByValue(pointDist, 2)

			// check the 2 shortest distances to see if they are the same
			if kvs[0].Value == kvs[1].Value {
				// it's a tie! Set the pointID to a negative value to make it easy to recognize
				pointID = -1
			} else {
				pointID = kvs[0].Key
			}
			grid.Set(i, j, float64(pointID))

			// if this point is on the edge of the grid, it has infinite areaPerPoint
			if (i == 0 || j == 0) || (i == (rows-1) || j == (cols-1)) {
				pointIsInfinite[pointID] = true
			}
		}
	}

	// get the total area for each point
	// (points with infinite area are disqualified)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			id := int(grid.At(i, j))
			if pointIsInfinite[id] {

				continue
			}
			areaPerPoint[id]++
		}
	}

	pointWithBiggestArea := BiggestByValue(areaPerPoint)

	fmt.Printf("part 1: %d\n", pointWithBiggestArea.Value)
	fmt.Printf("part 2: %d\n", regionSize)

}

// The distance between two points measured along axes at right angles.
// In a plane with p1 at (x1, y1) and p2 at (x2, y2), it is |x1 - x2| + |y1 - y2|.
func ManhattanDistance(x1, y1, x2, y2 int) int {

	x := x1 - x2
	if x < 0 {
		x *= -1
	}

	y := y1 - y2
	if y < 0 {
		y *= -1
	}

	return x + y
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

type kv struct {
	Key   int
	Value int
}

func BiggestByValue(input map[int]int) kv {
	output := make([]kv, 0)
	for k, v := range input {
		output = append(output, kv{k, v})
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output[0]
}

func SmallestNByValue(input map[int]int, n int) []kv {
	output := make([]kv, 0)
	for k, v := range input {
		output = append(output, kv{k, v})
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value < output[j].Value
	})

	return output[0:n]
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
