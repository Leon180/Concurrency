package main

import (
	"context"
	"fmt"

	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job"
	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job/task"
	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job/worker"
	workerPool "github.com/Leon180/Concurrency/1_concurrency_basics/worker/job/worker_pool"
)

const (
	numJobs    = 100
	numWorkers = 5
)

func main() {
	simpleWorkerPool(numJobs, numWorkers)
}

func simpleWorker(numJobs int, numWorkers int) {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channels
	jobs := make(chan job.Ticket, numJobs)
	results := make(chan job.Result, numJobs)

	// Start workers
	workers := make([]job.Worker, numWorkers)
	for i := range workers {
		workers[i] = worker.NewSimpleWorker()
		workers[i].ClaimWork(jobs, results)
		workers[i].StartWork(ctx)
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

	// Collect results
	resultCollector := make([]job.Result, numJobs)
	count := 0
	for r := range results {
		count++
		resultCollector[r.Input.ID-1] = r
		if r.Err != nil {
			fmt.Printf("🚫 job %d failed\n", r.Input.ID)
		} else {
			fmt.Printf("✅ job %d done\n", r.Input.ID)
		}
		if count == numJobs {
			break
		}
	}
	close(results)

	// Print results
	for _, result := range resultCollector {
		if result.Err != nil {
			fmt.Printf("🚫 job %d failed: %v\n", result.Input.ID, result.Err)
			continue
		}
		fmt.Printf("✅ job %d result: %v\n", result.Input.ID, result.Output)
	}
}

func simpleWorkerPool(numJobs int, numWorkers int) {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channels
	jobs := make(chan job.Ticket, numJobs)
	results := make(chan job.Result, numJobs)

	// Start workers
	workerPool := workerPool.NewSimpleWorkerPool(numWorkers, jobs, results)
	workerPool.StartWork(ctx)

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				jobs <- job.Ticket{ID: i + 1, Task: task.SquareTask(i + 1)}
			}
		}
	}()

	workerPool.Wait(numJobs)

	// Collect results
	resultCollector := workerPool.GetResults()

	// Print results
	for _, result := range resultCollector {
		if result.Err != nil {
			fmt.Printf("🚫 job %d failed: %v\n", result.Input.ID, result.Err)
			continue
		}
		fmt.Printf("✅ job %d result: %v\n", result.Input.ID, result.Output)
	}
}
