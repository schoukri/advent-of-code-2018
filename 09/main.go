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
)

type Marble struct {
	Value int
	Next  *Marble
	Prev  *Marble
}

func NewMarble(value int) *Marble {
	m := new(Marble)
	m.Value = value
	m.Prev = m
	m.Next = m
	return m
}

func (m *Marble) Insert(new *Marble) *Marble {
	next := m.Next
	new.Prev = m
	new.Next = next
	m.Next = new
	next.Prev = new
	return new
}

func (m *Marble) Remove() *Marble {
	prev := m.Prev
	next := m.Next
	prev.Next = next
	next.Prev = prev
	return next
}

func (m *Marble) Move(offset int) *Marble {
	marble := m
	if offset >= 0 {
		for i := 0; i < offset; i++ {
			marble = marble.Next
		}
	} else {
		for i := 0; i > offset; i-- {
			marble = marble.Prev
		}
	}
	return marble
}

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	part := flag.Int("part", 1, "the part of the challenge to run")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	for _, line := range lines {

		re := regexp.MustCompile(`^(\d+) players; last marble is worth (\d+) points`)

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			log.Fatalf("cannot parse line: %s", line)
		}

		numPlayers := mustParseInt(matches[1])
		lastMarble := mustParseInt(matches[2])

		if *part == 2 {
			lastMarble *= 100
		}

		playerScores := make([]int, numPlayers)
		player := 0

		marble := NewMarble(0)

		for nextMarble := 1; nextMarble <= lastMarble; nextMarble++ {

			player %= numPlayers

			if nextMarble%23 == 0 {
				marble = marble.Move(-7)
				playerScores[player] += nextMarble
				playerScores[player] += marble.Value
				marble = marble.Remove()
			} else {
				marble = marble.Move(1)
				marble = marble.Insert(NewMarble(nextMarble))
			}

			player++

		}

		sort.Ints(playerScores)

		fmt.Printf("part %d: %d\n", *part, playerScores[numPlayers-1])
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
