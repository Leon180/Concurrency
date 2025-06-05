package worker

import (
	"context"

	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job"
)

func NewSimpleWorker() job.Worker {
	return &simpleWorker{}
}

type simpleWorker struct {
	jobs    <-chan job.Ticket
	results chan<- job.Result
}

func (w *simpleWorker) ClaimWork(jobs <-chan job.Ticket, results chan<- job.Result) {
	w.jobs = jobs
	w.results = results
}

func (w *simpleWorker) StartWork(ctx context.Context) {
	go w.Work(ctx)
}

func (w *simpleWorker) Work(ctx context.Context) {
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
