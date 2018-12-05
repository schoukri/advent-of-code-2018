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

var (
	replacer *strings.Replacer
	letters  []string
)

func init() {
	letters = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	// generate upper+lower and lower+upper letter pairs to
	// (setup replacer to replace every pair with an empty string)
	replacements := make([]string, 0)
	for _, l := range letters {
		u := strings.ToUpper(l)
		replacements = append(replacements, l+u, "", u+l, "")
	}
	replacer = strings.NewReplacer(replacements...)
}

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	// there is only 1 line in the input
	line := lines[0]

	fmt.Printf("part 1: %d\n", React(line))

	length := make([]int, len(letters))
	for i, l := range letters {
		// replace only one letter (upper and lower) with empty string
		// then run React() and see which letter produces the shortest polymer
		repl := strings.NewReplacer(l, "", strings.ToUpper(l), "")
		result := repl.Replace(line)
		length[i] = React(result)
	}

	// sort the lengths slice to get the smallest length
	sort.Ints(length)
	fmt.Printf("part 2: %d\n", length[0])

}

func React(polymer string) int {
	var res string
	for {
		// as pairs are removed, remaining letters can join together to form new pairs
		// (keep running Replace() until no more replacements are detected)
		res = replacer.Replace(polymer)
		if polymer == res {
			break
		}
		polymer = res
	}

	return len(polymer)
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
