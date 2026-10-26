package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	results := make([]Result[R], len(in))
	if len(in) == 0 {
		return results, nil
	}

	jobs := make(chan int)
	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for i := range jobs {
			value, err := fn(ctx, in[i])
			results[i].Err = err
			if err == nil {
				results[i].Value = value
			}
		}
	}

	amount := min(workers, len(in))
	for i := 0; i < amount; i++ {
		wg.Add(1)
		go worker()
	}
	for i := range in {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return nil, ctx.Err()
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()

	return results, nil
}
