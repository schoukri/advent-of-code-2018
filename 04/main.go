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
	"strings"
)

var re = regexp.MustCompile(`Guard \#(\d+) begins shift`)

type kv struct {
	Key   int
	Value int
}

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	minutes := make(map[int]int)
	hist := make(map[int]map[int]int)

LINE:
	for i := 0; i < len(lines); {

		line := lines[i]
		i++

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			log.Fatal("no matches for guard")
		}
		id := mustParseInt(matches[1])

		for i < len(lines) {
			sleepLine := lines[i]
			if !strings.HasSuffix(sleepLine, "falls asleep") {
				continue LINE
			}

			wakeLine := lines[i+1]
			if !strings.HasSuffix(wakeLine, "wakes up") {
				log.Fatalf("invalid wake line: %s\n", wakeLine)
			}

			sleepMin := mustParseInt(sleepLine[15:17])
			wakeMin := mustParseInt(wakeLine[15:17])
			minutes[id] += wakeMin - sleepMin

			for m := sleepMin; m < wakeMin; m++ {
				if _, ok := hist[id]; !ok {
					hist[id] = make(map[int]int)
				}
				hist[id][m]++
			}
			i += 2
		}

	}

	winner := Top(minutes)
	fmt.Printf("winner: id=%d, mins=%d\n", winner.Key, winner.Value)

	winnerHist := hist[winner.Key]

	topMin := Top(winnerHist)

	fmt.Printf("topmin: min=%d, count=%d\n", topMin.Key, topMin.Value)

	part1 := winner.Key * topMin.Key
	fmt.Printf("part 1: %d\n", part1)

	topMinPerGuard := make(map[int]int)
	mostMinPerGuard := make(map[int]int)
	for id, h := range hist {
		top := Top(h)
		topMinPerGuard[id] = top.Key
		mostMinPerGuard[id] = top.Value
	}

	p2 := Top(mostMinPerGuard)

	part2 := p2.Key * topMinPerGuard[p2.Key]
	fmt.Printf("part 2: %d\n", part2)

}

func Top(input map[int]int) kv {
	output := make([]kv, 0)
	for k, v := range input {
		output = append(output, kv{k, v})
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output[0]
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

func mustParseLine(line string) []string {

	matches := re.FindStringSubmatch(line)
	if matches == nil {
		log.Fatalf("cannot parse line: %s", line)
	}

	return matches
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
