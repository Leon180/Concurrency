package task

import (
	"fmt"
	"math/rand"
	"time"
)

func Square(n int) int {
	return n * n
}

func SquareTask(n int) func() (any, error) {
	return func(n int) func() (any, error) {
		return func() (any, error) {
			time.Sleep(time.Second * time.Duration(rand.Intn(2)))
			if rand.Intn(10) > 8 {
				return nil, fmt.Errorf("unlucky")
			}
			return Square(n), nil
		}
	}(n)
}
