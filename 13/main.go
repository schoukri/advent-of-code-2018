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

	// jsonBytes, _ := json.MarshalIndent(track, "", "  ")
	// fmt.Printf("Track: %s\n", string(jsonBytes))

	printCart := false
	if strings.Contains(*filePath, "sample") {
		printCart = true
	}

TICK:
	for tick := 1; ; tick++ {

		fmt.Printf("TICK: %d\n", tick)
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
			num := track.RemoveCrashedCars()
			fmt.Printf("RemoveCrashedCars: %d (tick=%d)\n", num, tick)

			if len(track.Carts) < 2 {
				cart := track.Carts[0]
				fmt.Printf("part 2: %d,%d (tick=%d)\n", cart.X, cart.Y, tick)
				fmt.Printf("cart: %#v\n", cart)
				track.Print(true)
				break TICK
			}
		}

	}
	track.Print(printCart)

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

func (t *Track) RemoveCrashedCars() int {

	sort.Sort(ByPosition(t.Carts))

	count := make(map[int]map[int]int)
	for _, cart := range t.Carts {
		if _, ok := count[cart.X]; !ok {
			count[cart.X] = make(map[int]int)
		}
		count[cart.X][cart.Y]++
	}

	fmt.Println("--------------------------")
	remaining := make([]*Cart, 0)
	removed := 0
	for _, cart := range t.Carts {
		if count[cart.X][cart.Y] == 1 {
			remaining = append(remaining, cart)
		} else {
			fmt.Printf("REMOVED DUPE CART: %#v\n", cart)
			removed++
		}
	}

	t.Carts = remaining

	return removed
}

func (t *Track) Move(c *Cart) bool {

	//fmt.Printf("Move  Start: %#v\n", c)

	piece, ok := t.Pieces[c.X][c.Y]
	if !ok {
		log.Fatalf("missing track piece at %d,%d\n", c.X, c.Y)
	}

	newTurn := c.Turn
	newDirection := c.Direction

	switch piece {
	case Vertical:
		switch c.Direction {
		case North:
			c.Y--
		case South:
			c.Y++
		}
	case Horizontal:
		switch c.Direction {
		case West:
			c.X--
		case East:
			c.X++
		}
	case CurveBack:
		switch c.Direction {
		case North:
			c.X--
			newDirection = West
		case South:
			c.X++
			newDirection = East
		case West:
			c.Y--
			newDirection = North
		case East:
			c.Y++
			newDirection = South
		}
	case CurveForward:
		switch c.Direction {
		case North:
			c.X++
			newDirection = East
		case South:
			c.X--
			newDirection = West
		case West:
			c.Y++
			newDirection = South
		case East:
			c.Y--
			newDirection = North
		}
	case Intersection:
		// Each time a cart has the option to turn (by arriving at any intersection),
		// it turns left the first time, goes straight the second time,
		// turns right the third time, and then repeats those directions

		switch c.Turn {
		case Right:
			newTurn = Left
		case Left:
			newTurn = Straight
		case Straight:
			newTurn = Right
		}

		switch c.Direction {
		case North:
			switch c.Turn {
			case Right:
				newDirection = East
				c.X++
			case Left:
				newDirection = West
				c.X--
			case Straight:
				c.Y--
			}
		case South:
			switch c.Turn {
			case Right:
				newDirection = West
				c.X--
			case Left:
				newDirection = East
				c.X++
			case Straight:
				c.Y++
			}
		case East:
			switch c.Turn {
			case Right:
				newDirection = South
				c.Y++
			case Left:
				newDirection = North
				c.Y--
			case Straight:
				c.X++
			}
		case West:
			switch c.Turn {
			case Right:
				newDirection = North
				c.Y--
			case Left:
				newDirection = South
				c.Y++
			case Straight:
				c.X--
			}
		}
	}

	c.Turn = newTurn
	c.Direction = newDirection

	//fmt.Printf("Move Finish: %#v\n", c)

	// does another cart exist at the same location?
	// If so, it's a crash!
	for _, cart := range t.Carts {
		if cart.X == c.X && cart.Y == c.Y && cart.ID != c.ID {
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
