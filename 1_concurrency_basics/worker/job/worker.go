package job

import (
	"context"
)

type Worker interface {
	ClaimWork(jobs <-chan Ticket, results chan<- Result)
	StartWork(ctx context.Context)
}
