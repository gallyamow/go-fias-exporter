package workerpool

import (
	"context"
	"sync"
)

// RunWithWorkers runs the given handler function with the given number of workers.
// The handler must wrap and return its own errors.
// The handler must respect context cancellation.
func RunWithWorkers[T any, R any](ctx context.Context, jobCh <-chan T, handler func(ctx context.Context, job T) R, workers int) <-chan R {
	resCh := make(chan R)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := range workers {
		go func(workerId int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobCh:
					if !ok {
						return
					}

					jobRes := handler(ctx, job)

					select {
					case <-ctx.Done():
						return
					case resCh <- jobRes:
					}
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	return resCh
}
