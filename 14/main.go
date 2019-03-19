package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

type Scoreboard struct {
	Scores []int
	ElfOne int
	ElfTwo int
	show   bool
}

func main() {

	part := flag.Int("part", 1, "The part of the puzzle to run.")
	input := flag.String("input", "681901", "The input value for the puzzle.")
	show := flag.Bool("show", false, "Show the scoreboard for each round.")
	flag.Parse()

	var recipes int
	var sequence string

	switch *part {
	case 1:
		num, err := strconv.Atoi(*input)
		if err != nil {
			log.Fatalf("error: cannot convert input '%s' to number of recipes: %v", *input, err)
		}
		recipes = num
	case 2:
		sequence = *input
	default:
		log.Fatalf("invalid part number %d\n", *part)

	}

	// initialize scoreboard
	scoreboard := Scoreboard{
		Scores: []int{3, 7},
		ElfOne: 0,
		ElfTwo: 1,
		show:   *show,
	}

	scoreboard.Show()

	if *part == 1 {
		for scoreboard.Len() < recipes+10 {
			scoreboard.Combine()
			scoreboard.Show()
		}
		finalScore := scoreboard.Sequence(recipes, 10)
		fmt.Printf("part 1: %s\n", finalScore)

	} else if *part == 2 {
		sequenceLen := len(sequence)
		for {
			added := scoreboard.Combine()
			scoreboard.Show()

			if scoreboard.Len() >= sequenceLen {
				start1 := scoreboard.Len() - sequenceLen
				if sequence == scoreboard.Sequence(start1, sequenceLen) {
					fmt.Printf("part 2: %d\n", start1)
					break
				}

				if added == 2 && scoreboard.Len() > sequenceLen {
					start2 := start1 - 1
					if sequence == scoreboard.Sequence(start2, sequenceLen) {
						fmt.Printf("part 2: %d\n", start2)
						break
					}
				}
			}
		}
	}
}

func (sc *Scoreboard) Combine() int {
	added := 1
	combined := sc.Scores[sc.ElfOne] + sc.Scores[sc.ElfTwo]
	if combined >= 10 {
		sc.Scores = append(sc.Scores, 1)
		added = 2
	}
	sc.Scores = append(sc.Scores, combined%10)
	moveOne := 1 + sc.Scores[sc.ElfOne]
	moveTwo := 1 + sc.Scores[sc.ElfTwo]
	sc.ElfOne = (sc.ElfOne + moveOne) % len(sc.Scores)
	sc.ElfTwo = (sc.ElfTwo + moveTwo) % len(sc.Scores)

	return added
}

func (sc *Scoreboard) Sequence(start, length int) string {
	var sequence string
	for _, digit := range sc.Scores[start : start+length] {
		sequence += strconv.Itoa(digit)
	}
	return sequence
}

func (sc *Scoreboard) Show() {
	if !sc.show {
		return
	}
	for i, recipe := range sc.Scores {
		if i == sc.ElfOne {
			fmt.Printf("(%d)", recipe)
		} else if i == sc.ElfTwo {
			fmt.Printf("[%d]", recipe)
		} else {
			fmt.Printf(" %d ", recipe)
		}
	}
	fmt.Println()
}

func (sc *Scoreboard) Len() int {
	return len(sc.Scores)
}
