package workerPool

import (
	"context"
	"fmt"

	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job"
	"github.com/Leon180/Concurrency/1_concurrency_basics/worker/job/worker"
)

func NewSimpleWorkerPool(
	workerCnt int,
	jobs chan job.Ticket,
	results chan job.Result,
) *simpleWorkerPool {
	workers := make([]job.Worker, workerCnt)
	for i := range workers {
		workers[i] = worker.NewSimpleWorker()
	}
	return &simpleWorkerPool{
		workers:          workers,
		jobs:             jobs,
		results:          results,
		resultsCollector: make([]job.Result, 0),
	}
}

type simpleWorkerPool struct {
	workers          []job.Worker
	jobs             chan job.Ticket
	results          chan job.Result
	resultsCollector []job.Result
}

func (p *simpleWorkerPool) StartWork(ctx context.Context) {
	for i := range p.workers {
		p.workers[i].ClaimWork(p.jobs, p.results)
		p.workers[i].StartWork(ctx)
	}
	p.collectResults()
}

func (p *simpleWorkerPool) Wait(totalJobCnt int) {
	defer close(p.jobs)
	defer close(p.results)
	for {
		if len(p.resultsCollector) != totalJobCnt {
			continue
		}
		fmt.Printf("all %d jobs done\n", totalJobCnt)
		break
	}
}

func (p *simpleWorkerPool) collectResults() {
	go func() {
		for {
			r, ok := <-p.results
			if !ok {
				break
			}
			if r.Err != nil {
				fmt.Printf("🚫 job %d failed\n", r.Input.ID)
			} else {
				fmt.Printf("✅ job %d done\n", r.Input.ID)
			}
			p.resultsCollector = append(p.resultsCollector, r)
		}
	}()
}

func (p *simpleWorkerPool) GetResults() []job.Result {
	return p.resultsCollector
}
