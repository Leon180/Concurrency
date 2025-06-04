package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job"
	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job/task"
	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job/worker"
)

const (
	numJobs    = 100
	numWorkers = 5
)

func main() {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channels
	jobs := make(chan job.Ticket, numJobs)
	results := make(chan job.Result, numJobs)
	var wg sync.WaitGroup

	// Start workers
	workers := make([]job.Worker, numWorkers)
	for i := range workers {
		workers[i] = worker.NewSimpleWorker()
		workers[i].ClaimWork(jobs, results, &wg)
		workers[i].StartWork(ctx)
		wg.Add(1)
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i := 0; i < numJobs; i++ {
			select {
			case <-ctx.Done():
				return
			case jobs <- job.Ticket{ID: i + 1, Task: task.SquareTask(i + 1)}:
			}
		}
	}()

	// Wait for all workers to finish and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	resultCollector := make([]job.Result, numJobs)
	for r := range results {
		resultCollector[r.Input.ID-1] = r
		if r.Err != nil {
			fmt.Printf("🚫 job %d failed\n", r.Input.ID)
			continue
		}
		fmt.Printf("✅ job %d done\n", r.Input.ID)
	}

	// Print results
	for _, result := range resultCollector {
		if result.Err != nil {
			fmt.Printf("🚫 job %d failed: %v\n", result.Input.ID, result.Err)
			continue
		}
		fmt.Printf("✅ job %d result: %v\n", result.Input.ID, result.Output)
	}
}
