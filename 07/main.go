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

type Job struct {
	ID string
	// Started bool
	Done    bool
	Seconds int
	Start   int
	End     int
}

func NewJob(id string) *Job {
	j := &Job{ID: id}

	// Convert "A", "B", "C", etc to ASCII numbers (65, 66, 67, etc)
	// Then subtract the offset value 4 (or 64 for the sample data)
	// to get the number of seconds required to finish the job.
	bytes := []byte(j.ID)
	j.Seconds = int(bytes[0]) - jobSecondsOffset

	return j

}

type Queue struct {
	Workers []*Worker
	Done    map[string]bool
	Deps    map[string][]string
}

func NewQueue(deps map[string][]string, numWorkers int) *Queue {
	q := &Queue{
		Workers: make([]*Worker, 0),
		Done:    make(map[string]bool),
		Deps:    deps,
	}

	for i := 0; i < numWorkers; i++ {
		q.Workers = append(q.Workers, &Worker{ID: i})
	}

	return q
}

type Worker struct {
	ID   int
	Jobs []*Job
}

var (
	numWorkers       = 5
	jobSecondsOffset = 4
)

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

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

	queue := NewQueue(jobDeps, numWorkers)

CLOCK:
	for clock := 0; ; clock++ {

		// close any running jobs that have reached their end time
		runningJobs := queue.GetRunningJobs()
		for _, job := range runningJobs {
			if job.End > clock {
				// job is still running
				continue
			}
			if clock > job.End {
				log.Fatalf("current clock %d is after job end %d", clock, job.End)
			}
			fmt.Printf("[%04d] finished job %s\n", clock, job.ID)
			job.Done = true
			queue.Done[job.ID] = true
		}

		jobs := queue.GetReadyJobs()

		if len(jobs) == 0 {
			running := queue.GetRunningJobs()
			if len(running) == 0 {
				fmt.Printf("[%04d] no jobs remaining\n", clock)
				break CLOCK
			}
		}

		minJobEnd := math.MaxInt32
		for _, job := range jobs {
			worker := queue.GetWorker()
			if worker == nil {
				fmt.Printf("[%04d] skipping job %s because no worker is available\n", clock, job.ID)
				continue CLOCK
			}

			fmt.Printf("[%04d] started job %s with worker %d\n", clock, job.ID, worker.ID)

			// job.Started = true
			job.Start = clock
			job.End = clock + job.Seconds
			worker.AddJob(job)

			if job.End < minJobEnd {
				minJobEnd = job.End
			}
		}
	}

	// the total time spent is the max end time of all the last jobs
	maxEndTime := 0
	for _, worker := range queue.Workers {
		job := worker.GetLastJob()
		if job != nil {
			if job.End > maxEndTime {
				maxEndTime = job.End
			}
		}
	}

	fmt.Printf("part 2: %d\n", maxEndTime)

}

func (q *Queue) GetReadyJobs() []*Job {

	ready := make([]string, 0)

	running := make(map[string]bool)
	for _, job := range q.GetRunningJobs() {
		running[job.ID] = true
	}

	for jobStr, deps := range q.Deps {

		if q.Done[jobStr] || running[jobStr] {
			continue
		}

		notDone := 0
		for _, dep := range deps {
			if dep == jobStr {
				continue
			}
			if !q.Done[dep] {
				notDone++
				break
			}
		}

		if notDone == 0 {
			ready = append(ready, jobStr)
		}

	}

	sort.Strings(ready)

	jobs := make([]*Job, len(ready))

	for i, jobStr := range ready {
		jobs[i] = NewJob(jobStr)
	}

	return jobs
}

// GetWorker returns a worker that is the least busy (the smallest number of seconds worked)
func (q *Queue) GetWorker() *Worker {
	var worker *Worker
	minSeconds := math.MaxInt32

	for _, w := range q.Workers {
		job := w.GetLastJob()
		if job == nil {
			return w
		}

		// don't pick a worker already running a job
		if !job.Done {
			continue
		}

		if job.End < minSeconds {
			minSeconds = job.End
			worker = w
		}
	}

	return worker
}

func (w *Worker) AddJob(job *Job) {
	w.Jobs = append(w.Jobs, job)
}

func (q *Queue) GetRunningJobs() []*Job {
	jobs := make([]*Job, 0)
	for _, worker := range q.Workers {
		job := worker.GetLastJob()
		if job == nil || job.Done {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func (w *Worker) GetLastJob() *Job {
	if len(w.Jobs) == 0 {
		return nil
	}
	return w.Jobs[len(w.Jobs)-1]
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
