This Go program demonstrates a more complex implementation of a worker pool pattern. It allocates jobs to workers, where each job involves calculating the sum of the digits of a randomly generated number. Here's a detailed explanation of each part of the code:

### Package and Imports

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)
```

- **Package declaration**: `package main` indicates that this is a standalone executable program.
- **Imports**:
  - `fmt`: A package for formatted I/O.
  - `math/rand`: A package for generating random numbers.
  - `sync`: A package that provides synchronization primitives such as `WaitGroup`.
  - `time`: A package that provides functionality for measuring and displaying time.

### Types and Variables

```go
type Job struct {
	id  int
	num int
}

type Result struct {
	job   Job
	total int
}

var (
	jobs    = make(chan Job, 10)
	results = make(chan Result, 10)
)
```

- **Types**:
  - `Job`: A struct that represents a job with an `id` and a `num` (the number to process).
  - `Result`: A struct that represents the result of a job, containing the `job` and the `total` (sum of the digits).
- **Channels**:
  - `jobs`: A buffered channel for `Job` structs with a capacity of 10.
  - `results`: A buffered channel for `Result` structs with a capacity of 10.

### Sum Function

```go
func sum(number int) (total int) {
	no := number
	for no != 0 {
		digit := no % 10
		total += digit
		no /= 10
	}
	time.Sleep(2 * time.Second)
	return
}
```

- **Function declaration**: `sum` takes an integer `number` and returns the sum of its digits.
- **Logic**:
  - Extract digits of the number using modulo and integer division.
  - Add each digit to `total`.
  - Simulate processing time with `time.Sleep(2 * time.Second)`.

### Worker Function

```go
func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		output := Result{job, sum(job.num)}
		results <- output
	}
	wg.Done()
}
```

- **Function declaration**: `worker` processes jobs from the `jobs` channel and sends results to the `results` channel.
- **Parameters**: `wg` (a pointer to `sync.WaitGroup`) to synchronize the completion of all workers.
- **Logic**:
  - Loop over the `jobs` channel to process each job.
  - Calculate the sum of the digits using the `sum` function.
  - Send the result to the `results` channel.
  - Call `wg.Done()` to signal that the worker has finished.

### Create Worker Pool Function

```go
func createWorkerPool() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}
```

- **Function declaration**: `createWorkerPool` starts a pool of 10 worker goroutines.
- **Logic**:
  - Create a `sync.WaitGroup` and set the counter to 10.
  - Start 10 worker goroutines, passing the `WaitGroup` to each.
  - Wait for all workers to complete using `wg.Wait()`.
  - Close the `results` channel to signal no more results will be sent.

### Allocate Job Function

```go
func allocateJob() {
	for i := 0; i < 300; i++ {
		num := rand.Intn(999)
		job := Job{i, num}
		jobs <- job
	}
	close(jobs)
}
```

- **Function declaration**: `allocateJob` generates and sends 300 jobs to the `jobs` channel.
- **Logic**:
  - Loop 300 times to create jobs.
  - Generate a random number (0 to 998).
  - Create a `Job` struct with the loop index as `id` and the random number as `num`.
  - Send the job to the `jobs` channel.
  - Close the `jobs` channel to signal no more jobs will be sent.

### Result Function

```go
func result(done chan bool) {
	for result := range results {
		fmt.Printf("job id %d, número %d, soma dos dígitos %d\n", result.job.id, result.job.num, result.total)
	}
	done <- true
}
```

- **Function declaration**: `result` prints the results from the `results` channel.
- **Parameters**: `done` (a channel to signal completion).
- **Logic**:
  - Loop over the `results` channel to receive and print results.
  - Send `true` to the `done` channel when finished.

### Main Function

```go
func main() {
	go allocateJob()
	done := make(chan bool)
	go result(done)
	createWorkerPool()
	<-done
}
```

- **Logic**:
  - Start the `allocateJob` function as a goroutine to generate jobs concurrently.
  - Create a `done` channel to signal when results are processed.
  - Start the `result` function as a goroutine to process and print results.
  - Call `createWorkerPool` to start the worker pool and wait for completion.
  - Wait for the `result` goroutine to signal completion by receiving from the `done` channel.

### Summary

- The program creates a worker pool of 10 workers to process 300 jobs.
- Each job involves calculating the sum of the digits of a randomly generated number.
- Jobs and results are communicated via buffered channels.
- Synchronization is achieved using `sync.WaitGroup` and a `done` channel to ensure all jobs are processed before the program exits.