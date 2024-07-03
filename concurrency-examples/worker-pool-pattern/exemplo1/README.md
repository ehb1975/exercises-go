# Worker Poll pattern

This Go code demonstrates a simple implementation of a worker pool pattern. Let's break it down step by step:

### Package and Imports

```go
package main

import (
	"fmt"
	"time"
)
```
- The `main` package is the starting point of the Go program.
- The `fmt` package is used for formatted I/O.
- The `time` package is used to introduce delays in the program.

### Worker Function

```go
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
```
- The `worker` function simulates a worker that processes jobs.
- It takes three parameters: `id` (an identifier for the worker), `jobs` (a read-only channel from which it receives jobs), and `results` (a write-only channel to which it sends results).
- It loops over the `jobs` channel, processing each job it receives.
- For each job, it prints a message indicating the start and end of the job, waits for one second to simulate job processing time, and then sends the result (job ID multiplied by 2) to the `results` channel.

### Main Function

```go
func main() {

	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
```
- The `main` function coordinates the creation of jobs and workers.
- `numJobs` is a constant representing the number of jobs to be processed.
- Two channels are created: `jobs` (with a capacity of `numJobs`) to hold the jobs, and `results` (also with a capacity of `numJobs`) to hold the results.
- Three worker goroutines are started by calling the `worker` function with different `id`s (1, 2, and 3), `jobs` channel, and `results` channel.
- The `main` function sends `numJobs` (5) jobs into the `jobs` channel.
- After sending all jobs, it closes the `jobs` channel to indicate no more jobs will be sent.
- The `main` function waits to receive `numJobs` (5) results from the `results` channel.

### Summary

- The code demonstrates concurrency in Go using goroutines and channels.
- Three workers process five jobs concurrently.
- Each worker processes a job, simulates the job processing by sleeping for one second, and then sends the processed result to the `results` channel.
- The `main` function waits for all results before terminating, ensuring that all jobs are processed.