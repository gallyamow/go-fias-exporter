package workerpool_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gallyamow/go-fias-exporter/pkg/workerpool"
)

func TestRunWithWorkers(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4, 6, 8, 10}

	jobCh := make(chan int, len(jobs))
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	// handler: умножает число на 2
	handler := func(ctx context.Context, job int) int {
		return job * 2
	}

	ctx := context.Background()
	resCh := workerpool.RunWithWorkers[int, int](ctx, jobCh, handler, 3)

	results := []int{}
	for r := range resCh {
		results = append(results, r)
	}

	if len(results) != len(expected) {
		t.Fatalf("want %d, got %d", len(expected), len(results))
	}

	expectedMap := make(map[int]bool)
	for _, v := range expected {
		expectedMap[v] = true
	}

	for _, r := range results {
		if !expectedMap[r] {
			t.Errorf("got unexpected result: %d", r)
		}
	}
}

func TestRunWithWorkers_ContextCancel(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5}
	jobCh := make(chan int, len(jobs))
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	var processed int32
	handler := func(ctx context.Context, job int) int {
		atomic.AddInt32(&processed, 1)
		time.Sleep(50 * time.Millisecond)
		return job * 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	resCh := workerpool.RunWithWorkers[int, int](ctx, jobCh, handler, 2)

	count := 0
	for range resCh {
		count++
	}

	if atomic.LoadInt32(&processed) == int32(len(jobs)) {
		t.Errorf("cancel was expected")
	}

	t.Logf("handled %d before context cancelled", count)
}

func TestRunWithWorkers_EmptyJobs(t *testing.T) {
	jobCh := make(chan int)
	close(jobCh)

	handler := func(ctx context.Context, job int) int {
		t.Fatal("handler must not be called")
		return 0
	}

	ctx := context.Background()
	resCh := workerpool.RunWithWorkers[int, int](ctx, jobCh, handler, 2)

	for r := range resCh {
		t.Fatalf("got %d, want no result", r)
	}
}
