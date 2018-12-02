package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/arbovm/levenshtein"
)

func main() {

	var (
		filePath = flag.String("file", "input.txt", "file containing the input data")
		part     = flag.Int("part", 1, "the part number")
	)

	flag.Parse()

	if *filePath == "" {
		log.Fatal("file not specified")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	switch *part {
	case 1:
		part1(lines)
	case 2:
		part2(lines)
	default:
		log.Fatalf("invalid part number %d", *part)
	}

}

func part1(lines []string) {
	var (
		two   = 0
		three = 0
	)

	for _, line := range lines {

		count := make(map[rune]int)

		for _, letter := range line {
			count[letter]++
		}

		var found2, found3 int
		for _, v := range count {
			if v == 2 {
				found2 = 1
			} else if v == 3 {
				found3 = 1
			}
		}

		two += found2
		three += found3

	}
	fmt.Printf("answer: %d\n", two*three)
}

func part2(lines []string) {
	for _, one := range lines {
		for _, two := range lines {
			if one == two {
				continue
			}
			if levenshtein.Distance(one, two) == 1 {
				fmt.Println("one:    " + one)
				fmt.Println("two:    " + two)

				var result string
				for i, letter := range one {
					if one[i] == two[i] {
						result = result + string(letter)
					}
				}

				fmt.Println("answer: " + result)
				return

			}
		}
	}
}
