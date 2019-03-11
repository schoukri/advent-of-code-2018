package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Direction rune

const (
	North Direction = '^'
	South           = 'v'
	West            = '<'
	East            = '>'
)

type Turn int

const (
	Left Turn = iota
	Right
	Straight
)

type Piece rune

const (
	Vertical     Piece = '|'
	Horizontal         = '-'
	CurveForward       = '/'
	CurveBack          = '\\'
	Intersection       = '+'
	Crash              = 'X'
	Empty              = ' '
)

type Cart struct {
	ID        int
	X, Y      int
	Direction Direction
	Turn      Turn
	Crashed   bool
}

type Pieces map[int]map[int]Piece

type Track struct {
	Pieces Pieces
	Carts  []*Cart
	MaxX   int
	MaxY   int
}

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	part := flag.Int64("part", 1, "The part of the puzzle to run.")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	track := NewTrack(lines)

	printCart := false
	if strings.Contains(*filePath, "sample") {
		printCart = true
	}

TICK:
	for tick := 1; ; tick++ {

		sort.Sort(ByPosition(track.Carts))
		track.Print(printCart)

		for _, cart := range track.Carts {
			if ok := track.Move(cart); !ok {
				// CRASH!
				if *part == 1 {
					fmt.Printf("part 1: %d,%d (tick=%d)\n", cart.X, cart.Y, tick)
					break TICK
				}
			}
		}

		if *part == 2 {
			numRemaining := track.RemoveCrashedCarts()
			if numRemaining == 0 {
				log.Fatal("no carts left on the track")
			} else if numRemaining == 1 {
				cart := track.Carts[0]
				fmt.Printf("part 2: %d,%d (tick=%d)\n", cart.X, cart.Y, tick)
				break TICK
			}
		}

	}
}

func (t *Track) Print(show bool) {
	if !show {
		return
	}

	carts := make(map[int]map[int]Direction)
	for _, cart := range t.Carts {
		if _, ok := carts[cart.X]; !ok {
			carts[cart.X] = make(map[int]Direction)
		}
		if _, ok := carts[cart.X][cart.Y]; ok {
			carts[cart.X][cart.Y] = Crash
		} else {
			carts[cart.X][cart.Y] = cart.Direction
		}
	}

	for y := 0; y <= t.MaxY; y++ {
		for x := 0; x <= t.MaxX; x++ {
			if cart, ok := carts[x][y]; ok {
				fmt.Print(string(cart))
				continue
			}
			piece, ok := t.Pieces[x][y]
			if ok {
				fmt.Print(string(piece))
			} else {
				fmt.Print(string(Empty))
			}
		}
		fmt.Println()
	}
}

func (t *Track) RemoveCrashedCarts() int {
	remaining := make([]*Cart, 0)
	for _, cart := range t.Carts {
		if cart.Crashed {
			continue
		}
		remaining = append(remaining, cart)
	}
	t.Carts = remaining
	return len(t.Carts)
}

func (t *Track) Move(c *Cart) bool {

	piece, ok := t.Pieces[c.X][c.Y]
	if !ok {
		log.Fatalf("missing track piece at %d,%d\n", c.X, c.Y)
	}

	newDirection := c.Direction

	switch piece {
	case CurveBack:
		switch c.Direction {
		case North:
			newDirection = West
		case South:
			newDirection = East
		case West:
			newDirection = North
		case East:
			newDirection = South
		}
	case CurveForward:
		switch c.Direction {
		case North:
			newDirection = East
		case South:
			newDirection = West
		case West:
			newDirection = South
		case East:
			newDirection = North
		}
	case Intersection:
		// Each time a cart has the option to turn (by arriving at any intersection),
		// it turns left the first time, goes straight the second time,
		// turns right the third time, and then repeats those directions

		var newTurn Turn
		switch c.Turn {
		case Right:
			newTurn = Left
		case Left:
			newTurn = Straight
		case Straight:
			newTurn = Right
		}
		c.Turn = newTurn

		switch c.Direction {
		case North:
			switch c.Turn {
			case Right:
				newDirection = East
			case Left:
				newDirection = West
			}
		case South:
			switch c.Turn {
			case Right:
				newDirection = West
			case Left:
				newDirection = East
			}
		case East:
			switch c.Turn {
			case Right:
				newDirection = South
			case Left:
				newDirection = North
			}
		case West:
			switch c.Turn {
			case Right:
				newDirection = North
			case Left:
				newDirection = South
			}
		}
	}

	c.Direction = newDirection

	switch c.Direction {
	case North:
		c.Y--
	case South:
		c.Y++
	case West:
		c.X--
	case East:
		c.X++
	}

	// does another cart exist at the same location?
	// If so, it's a crash!
	for _, cart := range t.Carts {
		if cart.X == c.X && cart.Y == c.Y && cart.ID != c.ID && !cart.Crashed {
			cart.Crashed = true
			c.Crashed = true
			return false
		}
	}

	return true

}

func NewTrack(lines []string) *Track {
	pieces := make(Pieces)
	carts := make([]*Cart, 0)

	var maxX, maxY, id int
	for y, line := range lines {
		if y > maxY {
			maxY = y
		}
		for x, char := range line {
			if _, ok := pieces[x]; !ok {
				pieces[x] = make(map[int]Piece)
			}

			id++

			switch char {
			case '|', '-', '/', '\\', '+':
				pieces[x][y] = Piece(char)
			case '<':
				pieces[x][y] = Horizontal
				carts = append(carts, NewCart(id, x, y, West))
			case '>':
				pieces[x][y] = Horizontal
				carts = append(carts, NewCart(id, x, y, East))
			case 'v':
				pieces[x][y] = Vertical
				carts = append(carts, NewCart(id, x, y, South))
			case '^':
				pieces[x][y] = Vertical
				carts = append(carts, NewCart(id, x, y, North))
			case ' ':
				// do nothing for space char
			default:
				log.Fatalf("unexpected char in track data at %d,%d: [%s]\n", x, y, string(char))
			}

			if x > maxX {
				maxX = x
			}

		}
	}

	return &Track{
		Pieces: pieces,
		Carts:  carts,
		MaxX:   maxX,
		MaxY:   maxY,
	}
}

// this type implements the sort interface
type ByPosition []*Cart

func (s ByPosition) Len() int {
	return len(s)
}
func (s ByPosition) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByPosition) Less(i, j int) bool {
	if s[i].Y != s[j].Y {
		return s[i].Y < s[j].Y
	}
	return s[i].X < s[j].X
}

func NewCart(id, x, y int, dir Direction) *Cart {
	return &Cart{
		ID:        id,
		X:         x,
		Y:         y,
		Direction: dir,
		Turn:      Right,
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
