package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	var (
		freq              = 0
		seenFreq          = make(map[int]bool)
		finishedFirstFreq = false
		seenFreqTwice     = false
	)

	filePath := flag.String("file", "", "file containing the input data")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("file not specified")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for !seenFreqTwice {

		_, err := file.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			value, err := strconv.Atoi(line)
			if err != nil {
				log.Fatalf("cannot parse line '%s' as int: %v", line, err)
			}

			freq += value

			if seenFreq[freq] && !seenFreqTwice {
				fmt.Printf("first freq seen twice: %d\n", freq)
				seenFreqTwice = true
			}

			seenFreq[freq] = true
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if !finishedFirstFreq {
			fmt.Printf("final freq: %d\n", freq)
			finishedFirstFreq = true
		}
	}

}
