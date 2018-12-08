package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/stevenle/topsort"
)

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	numWorkers := 5
	jobSecondsOffset := 4
	if *filePath == "sample.txt" {
		numWorkers = 2
		jobSecondsOffset = 64
	}

	re := regexp.MustCompile(`^Step (\w) must be finished before step (\w) can begin.$`)

	graph := topsort.NewGraph()
	seen := make(map[string]bool)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			log.Fatalf("cannot parse line: %s", line)
		}

		preReq := matches[1]
		job := matches[2]

		// no-op to add same node more than once
		graph.AddNode(preReq)
		graph.AddNode(job)

		seen[preReq] = true
		seen[job] = true

		err := graph.AddEdge(job, preReq)
		if err != nil {
			log.Fatalf("cannot add edge: %s -> %s\n", job, preReq)
		}

	}

	maxDeps := 0
	jobDeps := make(map[string][]string)
	nodes := make([]string, 0)
	for n := range seen {
		deps, err := graph.TopSort(n)
		if err != nil {
			log.Fatalf("cannot get deps for node %s: %+v\n", n, err)
		}
		jobDeps[n] = deps
		nodes = append(nodes, n)

		if len(deps) > maxDeps {
			maxDeps = len(deps)
		}
	}

	sort.Strings(nodes)

	added := make(map[string]bool)

	jobs := make([]string, 0)

LOOP:
	for i := 0; i < len(nodes); {
		n := nodes[i]

		if added[n] {
			i++
			continue LOOP
		}

		if !Exists(n, jobDeps[n]) {
			i++
			continue LOOP
		}

		for _, dep := range jobDeps[n] {
			if n != dep && !added[dep] {
				i++
				continue LOOP
			}
		}

		jobs = append(jobs, n)
		added[n] = true

		i = 0

	}

	fmt.Printf("part 1: %s\n", strings.Join(jobs, ""))

	workers := make([]int, numWorkers)

	// Setup a map of the number of seconds required to process each job.
	// A = 61 seconds, B = 62 seconds, C = 63 seconds, etc.
	// We could setup a map manually, but that's no fun
	// Instead, simply convert each job (letter) to its ASCII number equivalent,
	// then subtract 4 (or 64 for the sample data) to get the correct number of seconds.
	jobSeconds := make(map[string]int)
	for _, job := range jobs {
		bytes := []byte(job)
		jobSeconds[job] = int(bytes[0]) - jobSecondsOffset
		fmt.Printf("job %s seconds = %d (deps: %v)\n", job, jobSeconds[job], jobDeps[job])
	}

	done := make(map[string]bool)
	runningJob := make(map[int]string)

QUEUE:
	for {

		fmt.Printf("running jobs: %#v\n", runningJob)
		readyJobs := make([]string, 0)

	JOB:
		for job, deps := range jobDeps {

			if done[job] {
				continue JOB
			}

			notDone := make([]string, 0)

		DEP:
			for _, dep := range deps {
				for _, job := range runningJob {
					if job == dep {
						continue DEP
					}
				}
				if !done[dep] {
					notDone = append(notDone, dep)
				}
			}

			if len(notDone) == 1 {
				readyJobs = append(readyJobs, notDone[0])
			}

		}

		sort.Strings(readyJobs)
		fmt.Printf("readyJobs: %#v\n", readyJobs)

		// do the work
		for _, job := range readyJobs {
			workerIdx := GetWorker(workers, runningJob)
			if workerIdx == -1 {
				fmt.Printf("worker not available for job %s (skip for now)\n", job)
				continue
			}
			workers[workerIdx] += jobSeconds[job]
			fmt.Printf("started job %s, worker=%d, seconds=%d, total=%d\n", job, workerIdx, jobSeconds[job], workers[workerIdx])

			runningJob[workerIdx] = job
		}

		// get the lowest time for the workers that are currently running job
		doneTime := math.MaxInt32
		for workerIdx := range runningJob {
			if workers[workerIdx] < doneTime {
				doneTime = workers[workerIdx]
			}
		}

		fmt.Printf("doneTime: %d\n", doneTime)

		for workerIdx, job := range runningJob {
			if workers[workerIdx] == doneTime {
				done[job] = true
				delete(runningJob, workerIdx)
				fmt.Printf("finished running job %s, worker=%d, time=%d\n", job, workerIdx, workers[workerIdx])
			}
		}

		// make the workers wait that did not do any work for this loop
		// (update ther time to same time as the worker who did work with the max time)
		// sort.Ints(workerSeconds)
		// waitSeconds := workerSeconds[0]
		// fmt.Printf("workerSeconds: %v\n", workerSeconds)
		for i := range workers {
			if _, ok := runningJob[i]; ok {
				continue
			}
			if workers[i] < doneTime {
				fmt.Printf("changing worker i=%d seconds from %d to %d\n", i, workers[i], doneTime)
				workers[i] = doneTime
			}
		}

		if len(readyJobs) == 0 {
			fmt.Println("NO MORE JOBS TO DO")

			// finish any running jobs
			for workerIdx, job := range runningJob {
				done[job] = true
				delete(runningJob, workerIdx)
				fmt.Printf("finished running job %s, worker=%d, time=%d\n", job, workerIdx, workers[workerIdx])
			}

			break QUEUE
		}

	}
	// must be > 1006 and < 1029
	// the worker with the most time is how long it took to complete all jobs
	// (sort the workers -- the last one will have the most time)
	sort.Ints(workers)

	fmt.Printf("part 2: %d\n", workers[len(workers)-1])
	for _, job := range jobs {
		if !done[job] {
			fmt.Printf("JOB NOT DONE: %s\n", job)
		}
	}

}

// GetWorker returns the index of the worker that is least busy (the smallest number of seconds worked)
func GetWorker(workers []int, runningJobs map[int]string) int {
	idx := -1 // if we return -1, that means a worker is not avilable
	minSeconds := math.MaxInt32

	for i, seconds := range workers {
		// don't pick a worker already running a job
		if _, ok := runningJobs[i]; ok {
			continue
		}
		if seconds < minSeconds {
			minSeconds = seconds
			idx = i
		}
	}
	return idx
}

func Exists(node string, deps []string) bool {
	for _, dep := range deps {
		if node == dep {
			return true
		}
	}
	return false
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
