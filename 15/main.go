package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

type Type rune

const (
	TypeElf    Type = 'E'
	TypeGoblin      = 'G'
	TypeWall        = '#'
	TypeOpen        = '.'
)

type Piece struct {
	X, Y      int
	Type      Type
	HitPoints int
}

type Grid [][]*Piece

type Path []*Piece
type Paths []Path

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	//part := flag.Int64("part", 1, "The part of the puzzle to run.")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	grid := NewGrid(lines)

ROUND:
	for round := 1; ; round++ {

		fmt.Printf("ROUND: %d\n", round)
		//grid.Print()

		// get a list of all units, in the order they will take a turn
		units := make([]*Piece, 0)
		for _, row := range grid {
			for _, piece := range row {
				if piece.Type == TypeElf || piece.Type == TypeGoblin {
					units = append(units, piece)
				}
			}
		}

		for _, unit := range units {

			// skip units that have died (become an empty square)
			if !(unit.Type == TypeElf || unit.Type == TypeGoblin) {
				continue
			}

			// if the unit is already in range of another unit, attack it (and don't move)
			if grid.AttackInRangeTarget(unit) {
				continue
			}

			targets := grid.Targets(unit)
			if len(targets) == 0 {
				fmt.Printf("part 1: %d\n", grid.TotalHitPoints()*(round-1))
				break ROUND
			}

			minPath := math.MaxInt32
			allPaths := make(Paths, 0)
			for _, target := range targets {
				for _, open := range grid.OpenSquares(target) {
					paths := grid.GetPaths(unit, open, minPath)
					if len(paths) > 0 {
						sort.Sort(ByShortestPath(paths))
						if len(paths[0]) < minPath {
							minPath = len(paths[0])
						}
						allPaths = append(allPaths, paths...)
					}
				}
			}

			if len(allPaths) == 0 {
				continue
			}

			// select the shortest path
			sort.Sort(ByShortestPath(allPaths))
			path := allPaths[0]

			// move to the unit to the first square in the path
			grid.Move(unit, path[0])

			// launch attack if the move put this piece in range of another
			grid.AttackInRangeTarget(unit)
		}
	}

}

func (grid Grid) GetPaths(start, end *Piece, minPath int) Paths {
	path := make(Path, 0)
	paths := make([]Path, 0)
	//minPath := math.MaxInt32
	return grid.getPaths(start, end, minPath, path, paths)
}

func (grid Grid) getPaths(start, end *Piece, minPath int, path Path, paths []Path) []Path {
	seen := make(map[Piece]bool)
	for _, piece := range path {
		seen[*piece] = true
	}
	nextRound := make([]*Piece, 0)
	for _, open := range grid.OpenSquares(start) {
		if _, ok := seen[*open]; ok {
			continue
		}
		if open.X == end.X && open.Y == end.Y {
			path = append(path, open)
			if len(path) < minPath {
				minPath = len(path)
			}
			paths = append(paths, path)
			return paths
		}
		if len(path) <= minPath {
			nextRound = append(nextRound, open)
		}
	}

	for _, open := range nextRound {
		pathCopy := make(Path, len(path))
		copy(pathCopy, path)
		pathCopy = append(pathCopy, open)
		paths = grid.getPaths(open, end, minPath, pathCopy, paths)
	}

	return paths
}

func (grid Grid) getPathsOne(start, end *Piece, path Path, paths []Path) []Path {
	seen := make(map[Piece]bool)
	for _, piece := range path {
		seen[*piece] = true
	}
	for _, open := range grid.OpenSquares(start) {
		if _, ok := seen[*open]; ok {
			continue
		}
		pathCopy := make(Path, len(path))
		copy(pathCopy, path)
		pathCopy = append(pathCopy, open)
		if open.X == end.X && open.Y == end.Y {
			paths = append(paths, pathCopy)
		} else {
			paths = grid.getPathsOne(open, end, pathCopy, paths)
		}
	}
	return paths
}

func (grid Grid) TotalHitPoints() int {
	total := 0
	for _, row := range grid {
		for _, piece := range row {
			if piece.Type == TypeElf || piece.Type == TypeGoblin {
				total += piece.HitPoints
			}
		}
	}
	return total
}

func (grid Grid) Move(source, dest *Piece) {

	if source != grid[source.Y][source.X] {
		log.Fatalf("source is not the same piece on the grid: source=%+v, grid=%+v\n", source, grid[source.Y][source.X])
	}
	if !(source.Type == TypeElf || source.Type == TypeGoblin) {
		log.Fatalf("source is not of type elf or goblin: %+v\n", source)
	}

	if dest != grid[dest.Y][dest.X] {
		log.Fatalf("dest is not the same piece on the grid: dest=%+v, grid=%+v\n", dest, grid[dest.Y][dest.X])
	}
	if dest.Type != TypeOpen {
		log.Fatalf("destination square is not open: %+v\n", dest)
	}

	// swap source and dest positions
	source.X, dest.X = dest.X, source.X
	source.Y, dest.Y = dest.Y, source.Y

	grid[source.Y][source.X] = source
	grid[dest.Y][dest.X] = dest

}

func (grid Grid) OpenSquares(unit *Piece) []*Piece {
	squares := make([]*Piece, 0)
	for _, square := range grid.Adjacent(unit) {
		if square.Type == TypeOpen {
			squares = append(squares, square)
		}
	}
	return squares
}

func (grid Grid) AttackInRangeTarget(unit *Piece) bool {
	targets := make([]*Piece, 0)
	for _, adjacent := range grid.Adjacent(unit) {
		if unit.IsEnemy(adjacent) {
			targets = append(targets, adjacent)
		}
	}

	if len(targets) == 0 {
		return false
	}

	// pick the target with the fewest hit points
	sort.Sort(ByHitPoints(targets))
	target := targets[0]

	// attack!!!
	target.HitPoints -= 3

	// if the target has zero or fewer hit points, it dies (replace it with an open square)
	if target.HitPoints <= 0 {
		target.Type = TypeOpen
	}

	return true
}

func (grid Grid) Adjacent(unit *Piece) []*Piece {
	pieces := make([]*Piece, 0)

	adjacent := []*Piece{
		&Piece{X: unit.X, Y: unit.Y - 1}, // up
		&Piece{X: unit.X, Y: unit.Y + 1}, // down
		&Piece{X: unit.X - 1, Y: unit.Y}, // left
		&Piece{X: unit.X + 1, Y: unit.Y}, // right
	}

	for _, adj := range adjacent {
		if adj.Y < 0 || adj.Y > len(grid) {
			continue
		}
		if adj.X < 0 || adj.X > len(grid[adj.Y]) {
			continue
		}
		piece := grid[adj.Y][adj.X]
		if piece.Type != TypeWall {
			pieces = append(pieces, piece)
		}
	}
	return pieces
}

func (grid Grid) Targets(unit *Piece) []*Piece {
	targets := make([]*Piece, 0)
	for _, row := range grid {
		for _, piece := range row {
			if unit.IsEnemy(piece) {
				targets = append(targets, piece)
			}
		}
	}
	return targets
}

func (piece *Piece) IsEnemy(target *Piece) bool {
	return (piece.Type == TypeElf && target.Type == TypeGoblin) || (piece.Type == TypeGoblin && target.Type == TypeElf)
}

type ByHitPoints []*Piece

func (s ByHitPoints) Len() int {
	return len(s)
}
func (s ByHitPoints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByHitPoints) Less(i, j int) bool {

	// primary sort: fewest hit points
	if s[i].HitPoints != s[j].HitPoints {
		return s[i].HitPoints < s[j].HitPoints
	}

	// secondary sort: reading order
	if s[i].Y != s[j].Y {
		return s[i].Y < s[j].Y
	}
	return s[i].X < s[j].X
}

type ByShortestPath Paths

func (s ByShortestPath) Len() int {
	return len(s)
}
func (s ByShortestPath) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByShortestPath) Less(i, j int) bool {

	// primary sort: by shortest length
	if len(s[i]) != len(s[j]) {
		return len(s[i]) < len(s[j])
	}

	// secondary sort: reading order of last piece
	lastPieceI := s[i][len(s[i])-1]
	lastPieceJ := s[j][len(s[j])-1]
	if lastPieceI.Y != lastPieceJ.Y {
		return lastPieceI.Y < lastPieceJ.Y
	}
	if lastPieceI.X != lastPieceJ.X {
		return lastPieceI.X < lastPieceJ.X
	}

	// tertiary sort: reading order of first piece
	firstPieceI := s[i][0]
	firstPieceJ := s[j][0]
	if firstPieceI.Y != firstPieceJ.Y {
		return firstPieceI.Y < firstPieceJ.Y
	}
	return firstPieceI.X < firstPieceJ.X

}

func (g Grid) Print() {
	for _, row := range g {
		hits := make([]string, 0)
		for _, piece := range row {
			fmt.Print(string(piece.Type))
			if piece.Type == TypeElf || piece.Type == TypeGoblin {
				hits = append(hits, fmt.Sprintf("%s(%d)", string(piece.Type), piece.HitPoints))
			}
		}
		fmt.Println("   " + strings.Join(hits, ", "))
	}
	fmt.Println()
}

func NewGrid(lines []string) Grid {
	grid := make(Grid, len(lines))

	for y, line := range lines {
		row := make([]*Piece, len(line))
		for x, char := range line {
			typ := Type(char)
			switch typ {
			case TypeElf, TypeGoblin:
				row[x] = &Piece{x, y, typ, 200}
			case TypeWall, TypeOpen:
				row[x] = &Piece{x, y, typ, 0}
			default:
				log.Fatalf("unexpected char in grid data at %d,%d: [%s]\n", x, y, string(char))
			}
		}
		grid[y] = row
	}

	return grid
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
