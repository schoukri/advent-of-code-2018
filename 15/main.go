package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

type Type rune

const (
	TypeElf    Type = 'E'
	TypeGoblin      = 'G'
	TypeWall        = '#'
	TypeOpen        = '.'
)

type Piece struct {
	Type        Type
	HitPoints   int
	AttackPower int
}

type Node struct {
	ID    int
	X, Y  int
	Piece *Piece
}

type Graph struct {
	Nodes []*Node
	Edges [][]bool
}

type Path []*Node

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	part := flag.Int("part", 1, "The part of the puzzle to run.")
	power := flag.Int("power", 3, "The attack power for Elves (for part 2 only).")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	graph := NewGraph(lines)

	// if part 2, increase the attack power for Elves
	if *part == 2 && *power > 3 {
		for _, node := range graph.Nodes {
			if node.Piece != nil && node.Piece.Type == TypeElf {
				node.Piece.AttackPower = *power
			}
		}
	}

ROUND:
	for round := 1; ; round++ {

		fmt.Printf("ROUND: %d\n", round)

		// get a list of all units, in the order they will take a turn
		units := make([]*Node, 0)
		for _, node := range graph.Nodes {
			if node.Piece != nil {
				units = append(units, node)
			}
		}
		sort.Sort(ByPosition(units))

		for _, unit := range units {

			// skip any units that are dead
			if unit.Piece == nil {
				continue
			}

			// if the unit is already in range of another unit, attack it (and don't move)
			if graph.AttackInRangeTarget(unit, *part) {
				continue
			}

			targets := graph.Targets(unit)
			if len(targets) == 0 {
				fmt.Printf("part %d: %d\n", *part, graph.TotalHitPoints()*(round-1))
				break ROUND
			}

			allPaths := make([]Path, 0)
			for _, target := range targets {
				for _, adj := range graph.Adjacent(target) {
					if adj.Piece == nil { // is open?
						paths := graph.GetPaths(unit, adj)
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

			// move the unit to the first adjacent square in the path
			dest := path[1]
			graph.Move(unit, dest)

			// launch attack if the move put this piece in range of another
			graph.AttackInRangeTarget(dest, *part)
		}
	}
}

func (g *Graph) GetPaths(start, end *Node) []Path {
	paths := make([][]int, 0)

	queue := make([]int, 0)
	visited := make([]bool, len(g.Nodes))
	pred := make(map[int]int)

	queue = append(queue, start.ID)
	visited[start.ID] = true
	pred[start.ID] = -1

QUEUE:
	for len(queue) > 0 {
		// pop
		currID := queue[0]
		queue = queue[1:]

	EDGE:
		for nodeID, ok := range g.Edges[currID] {
			if !ok || visited[nodeID] {
				continue
			}
			if nodeID == end.ID {
				path := make([]int, 0)
				path = append(path, nodeID)
				for p := currID; p > -1; {
					path = append(path, p)
					if pNext, ok := pred[p]; ok {
						p = pNext
					}
				}
				if len(paths) > 0 {
					// compare the length of this path to the last path added to the list
					if len(path) > len(paths[len(paths)-1]) {
						break QUEUE
					}
				}
				paths = append(paths, path)
				break EDGE
			}
			// add node to queue if its open
			if g.Nodes[nodeID].Piece == nil {
				queue = append(queue, nodeID)
			}
			visited[nodeID] = true
			pred[nodeID] = currID
		}
	}
	realPaths := make([]Path, len(paths))
	for i, path := range paths {
		for j := len(path) - 1; j >= 0; j-- {
			realPaths[i] = append(realPaths[i], g.Nodes[path[j]])
		}

	}
	//fmt.Printf("paths (%d -> %d): %+v\n", start.ID, end.ID, paths)
	return realPaths
}

func (g *Graph) Move(source, dest *Node) {
	if source.Piece == nil {
		log.Fatalf("source piece is empty: %+v", source)
	}
	if dest.Piece != nil {
		log.Fatalf("dest piece is not empty: %+v", dest)
	}

	// swap pieces
	source.Piece, dest.Piece = dest.Piece, source.Piece
}

func (g *Graph) Targets(node *Node) []*Node {
	if node.Piece == nil {
		log.Fatalf("node does not have a piece: %+v", node)
	}
	targets := make([]*Node, 0)
	for _, target := range g.Nodes {
		if target.Piece == nil || node.Piece.Type == target.Piece.Type {
			continue
		}
		targets = append(targets, target)
	}
	return targets
}

func (g *Graph) Adjacent(node *Node) []*Node {
	nodes := make([]*Node, 0)
	for adjID, ok := range g.Edges[node.ID] {
		if ok {
			nodes = append(nodes, g.Nodes[adjID])
		}
	}
	return nodes
}

func (g *Graph) AttackInRangeTarget(node *Node, part int) bool {
	if node.Piece == nil {
		log.Fatalf("node piece is empty: %+v", node)
	}
	targets := make([]*Node, 0)
	for _, adj := range g.Adjacent(node) {
		if adj.Piece == nil || node.Piece.Type == adj.Piece.Type {
			continue
		}
		targets = append(targets, adj)
	}

	if len(targets) == 0 {
		return false
	}

	// pick the target with the fewest hit points
	sort.Sort(ByHitPoints(targets))
	target := targets[0]

	// attack!!!
	target.Piece.HitPoints -= node.Piece.AttackPower

	// if the target piece has zero or fewer hit points, it dies
	if target.Piece.HitPoints <= 0 {
		if part == 2 && target.Piece.Type == TypeElf {
			log.Fatalln("An elf has been killed in part 2!")
		}
		target.Piece = nil
	}

	return true
}

func (g *Graph) TotalHitPoints() int {
	total := 0
	for _, node := range g.Nodes {
		if node.Piece != nil {
			total += node.Piece.HitPoints
		}
	}
	return total
}

type ByShortestPath []Path

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

	// tertiary sort: reading order of first adjacent piece
	firstPieceI := s[i][1]
	firstPieceJ := s[j][1]
	if firstPieceI.Y != firstPieceJ.Y {
		return firstPieceI.Y < firstPieceJ.Y
	}
	return firstPieceI.X < firstPieceJ.X

}

type ByPosition []*Node

func (s ByPosition) Len() int {
	return len(s)
}
func (s ByPosition) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByPosition) Less(i, j int) bool {

	// sort by reading order
	if s[i].Y != s[j].Y {
		return s[i].Y < s[j].Y
	}
	return s[i].X < s[j].X
}

type ByHitPoints []*Node

func (s ByHitPoints) Len() int {
	return len(s)
}
func (s ByHitPoints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByHitPoints) Less(i, j int) bool {

	// primary sort: fewest hit points
	if s[i].Piece.HitPoints != s[j].Piece.HitPoints {
		return s[i].Piece.HitPoints < s[j].Piece.HitPoints
	}

	// secondary sort: reading order
	if s[i].Y != s[j].Y {
		return s[i].Y < s[j].Y
	}
	return s[i].X < s[j].X
}

func NewGraph(lines []string) *Graph {

	grid := make(map[int]map[int]*Node)
	nodes := make([]*Node, 0)

	id := -1
	for y, line := range lines {
		for x, char := range line {
			if char == TypeWall {
				continue
			}

			id++
			node := &Node{
				ID: id,
				X:  x,
				Y:  y,
			}

			if _, ok := grid[x]; !ok {
				grid[x] = make(map[int]*Node)
			}

			grid[x][y] = node
			nodes = append(nodes, node)

			if char == TypeOpen {
				continue
			}
			piece := &Piece{
				Type:        Type(char),
				HitPoints:   200,
				AttackPower: 3,
			}
			node.Piece = piece
		}
	}

	// store the edges in an adjacency matrix
	edges := make([][]bool, len(nodes))
	for i := 0; i < len(edges); i++ {
		edges[i] = make([]bool, len(nodes))
	}

	for i, node := range nodes {
		if i != node.ID {
			log.Fatalf("node index=%d not equal to id=%d\n", i, node.ID)
		}
		if _, ok := grid[node.X][node.Y]; ok {
			//fmt.Printf("node exists: %+v\n", node)

			neighbors := []struct {
				x, y int
			}{
				{node.X, node.Y - 1},
				{node.X - 1, node.Y},
				{node.X + 1, node.Y},
				{node.X, node.Y + 1},
			}
			for _, n := range neighbors {
				if adj, ok := grid[n.x][n.y]; ok {
					edges[node.ID][adj.ID] = true
					edges[adj.ID][node.ID] = true
				}
			}
		}
	}

	return &Graph{
		Nodes: nodes,
		Edges: edges,
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
