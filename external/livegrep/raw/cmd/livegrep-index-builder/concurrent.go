package main

// taskClosure is a wrapper thunk that closes over a task function. It should return a non-nil error
// if the wrapped logic encounters an error.
type taskClosure func() error

// executeTasks concurrently executes a slice of lazily-defined tasks, with a concurrency limit. It
// collects all errors from the tasks and propagates it as a slice of errors. Note that this
// function does not terminate prematurely when a task encounters an error.
func executeTasks(tasks []taskClosure, concurrency int) []error {
	var errs []error

	errC := make(chan error, len(tasks))
	semaphore := make(chan bool, concurrency)

	for _, task := range tasks {
		// Wait on the semaphore
		semaphore <- true

		go func(task taskClosure) {
			defer func() {
				// Signal on the semaphore when the task completes execution
				<-semaphore
			}()

			if err := task(); err != nil {
				errC <- err
			}
		}(task)
	}

	// Block until all tasks have completed executing
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}

	// Collect all errors
	for {
		select {
		case err := <-errC:
			errs = append(errs, err)
		default:
			return errs
		}
	}
}
