package job

import (
	"context"
	"sync"
)

type Worker interface {
	ClaimWork(jobs <-chan Ticket, results chan<- Result, wg *sync.WaitGroup)
	StartWork(ctx context.Context)
}
