package job

type Ticket struct {
	ID   int
	Task func() (any, error)
}
