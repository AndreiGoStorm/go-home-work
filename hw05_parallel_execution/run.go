package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var limitError int32

func Run(tasks []Task, n, m int) error {
	toProcess := make(chan Task)
	limitError = int32(m)
	go process(tasks, toProcess)

	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(toProcess)
		}()
	}
	wg.Wait()

	if limitError <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(toProcess <-chan Task) {
	for task := range toProcess {
		err := task()
		if err != nil {
			atomic.AddInt32(&limitError, -1)
		}
	}
}

func process(tasks []Task, toProcess chan<- Task) {
	for i := 0; i < len(tasks); i++ {
		limitError := atomic.LoadInt32(&limitError)
		if limitError <= 0 {
			break
		}
		toProcess <- tasks[i]
	}
	close(toProcess)
}
