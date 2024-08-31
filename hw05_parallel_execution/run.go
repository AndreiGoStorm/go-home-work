package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type WorkerError struct {
	limitError int32
	error      error
	sync.Once
}

func Run(tasks []Task, n, m int) error {
	toProcess := make(chan Task, len(tasks))
	go process(tasks, toProcess)

	workerError := &WorkerError{limitError: int32(m), error: nil}
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(workerError, toProcess)
		}()
	}
	wg.Wait()

	return workerError.error
}

func worker(workerError *WorkerError, toProcess <-chan Task) {
	for task := range toProcess {
		limitError := atomic.LoadInt32(&workerError.limitError)
		if limitError <= 0 {
			workerError.Do(func() {
				workerError.error = ErrErrorsLimitExceeded
			})
			return
		}
		err := task()
		if err != nil {
			atomic.AddInt32(&workerError.limitError, -1)
		}
	}
}

func process(tasks []Task, toProcess chan<- Task) {
	for i := 0; i < len(tasks); i++ {
		toProcess <- tasks[i]
	}
	close(toProcess)
}
