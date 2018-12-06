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
	var regionSize int

	grid := mat.NewDense(maxWidth, maxHeight, nil)
	rows, cols := grid.Dims()

	// for each cell on the grid, calculate the distance to each Point
	// and store the ID of the Point with the shortest distance.
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
			pointIDs := KeysSortedByValueAsc(pointDist)

			// check the top 2 points with the shortest distance to see if they are the same distance
			if pointDist[pointIDs[0]] == pointDist[pointIDs[1]] {
				// it's a tie! Set the pointID to a negative value to make it easy to recognize
				pointID = -1
			} else {
				pointID = pointIDs[0]
			}
			grid.Set(i, j, float64(pointID))

			// if this point is on the edge of the grid, it has infinite area
			if (i == 0 || j == 0) || (i == (rows-1) || j == (cols-1)) {
				pointIsInfinite[pointID] = true
			}
		}
	}

	// get the total area for each point
	// (points with infinite area are disqualified)
	areaPerPoint := make(map[int]int)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			id := int(grid.At(i, j))
			if pointIsInfinite[id] {
				continue
			}
			areaPerPoint[id]++
		}
	}

	// get the point id with largest area
	pointIDs := KeysSortedByValueDesc(areaPerPoint)
	areaSize := areaPerPoint[pointIDs[0]]

	fmt.Printf("part 1: %d\n", areaSize)
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

func KeysSortedByValueDesc(input map[int]int) []int {
	kvs := make([]kv, 0)
	for k, v := range input {
		kvs = append(kvs, kv{k, v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})

	kk := make([]int, len(kvs))
	for i, k := range kvs {
		kk[i] = k.Key
	}
	return kk
}

func KeysSortedByValueAsc(input map[int]int) []int {
	kvs := make([]kv, 0)
	for k, v := range input {
		kvs = append(kvs, kv{k, v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value < kvs[j].Value
	})

	kk := make([]int, len(kvs))
	for i, k := range kvs {
		kk[i] = k.Key
	}
	return kk
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
