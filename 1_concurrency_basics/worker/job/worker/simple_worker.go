package worker

import (
	"context"
	"sync"

	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job"
)

func NewSimpleWorker() job.Worker {
	return &simpleWorker{}
}

type simpleWorker struct {
	jobs    <-chan job.Ticket
	results chan<- job.Result
	wg      *sync.WaitGroup
}

func (w *simpleWorker) ClaimWork(jobs <-chan job.Ticket, results chan<- job.Result, wg *sync.WaitGroup) {
	w.jobs = jobs
	w.results = results
	w.wg = wg
}

func (w *simpleWorker) StartWork(ctx context.Context) {
	go w.Work(ctx)
}

func (w *simpleWorker) Work(ctx context.Context) {
	defer w.wg.Done() // Signal completion when work is done
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-w.jobs:
			if !ok {
				return
			}
			if ctx.Err() != nil {
				w.results <- job.Result{Input: j, Err: ctx.Err()}
				return
			}
			output, err := j.Task()
			if err != nil {
				w.results <- job.Result{Input: j, Err: err}
				continue
			}
			w.results <- job.Result{Input: j, Output: output}
		}
	}
}
