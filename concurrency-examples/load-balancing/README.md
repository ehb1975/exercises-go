# LOAD BALANCING

- Channels can be used to distribute workloads among multiple goroutines. By creating a channel to receive tasks, and multiple goroutines to process them, you can achieve load balancing and take advantage of the concurrency of Go

- Source: https://medium.com/@varmapooja09/mastering-go-channels-part-1-7baa978a7de8

This Go code implements a simple concurrent worker pool. Workers read tasks from a channel, process them, and send the results to another channel. Here's a detailed explanation of each part of the code:

### Package and Imports

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)
```
- `package main`: Defines the package name.
- The imports include packages for formatted I/O (`fmt`), random number generation (`rand`), and time manipulation (`time`).

### Function `worker`

```go
func worker(id int, tasks <-chan int, results chan<- int) {
	// Process tasks from the input channel
	for task := range tasks {
		fmt.Printf("Worker %d started task %d\n", id, task)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate processing time
		results <- task * 2                                   // Send the result to the output channel
		fmt.Printf("Worker %d finished task %d\n", id, task)
	}
}
```
- This function represents a worker that processes tasks.
- It takes three parameters:
  - `id`: an integer representing the worker's ID.
  - `tasks`: a read-only channel from which the worker receives tasks.
  - `results`: a write-only channel to which the worker sends results.
- Inside a `for` loop, the worker reads tasks from the `tasks` channel.
- It simulates task processing by sleeping for a random duration (up to 3 seconds).
- The worker sends the result (task * 2) to the `results` channel.
- It prints messages indicating when it starts and finishes a task.

### `main` Function

```go
func main() {
	// Create channels for tasks and results
	tasks := make(chan int)
	results := make(chan int)

	// Start multiple workers to process tasks
	for i := 1; i <= 3; i++ {
		go worker(i, tasks, results)
	}

	// Send some tasks to the input channel
	for i := 1; i <= 10; i++ {
		tasks <- i
	}
	close(tasks)

	// Collect the results from the output channel
	for i := 1; i <= 10; i++ {
		result := <-results
		fmt.Printf("Result %d: %d\n", i, result)
	}
}
```
- The `main` function coordinates the entire program.
- It creates two channels: `tasks` for sending tasks to workers and `results` for receiving processed results.
- It starts three worker goroutines by calling the `worker` function with different IDs (1 to 3) in a loop.
- It sends 10 tasks (integers from 1 to 10) to the `tasks` channel and then closes the channel.
- It collects 10 results from the `results` channel and prints them.

### Summary
- **Concurrency**: The code uses goroutines to run multiple workers concurrently.
- **Channels**: Channels facilitate communication between the main function and the workers.
- **Task Processing**: Workers simulate task processing by sleeping for a random duration and then sending results to the `results` channel.
- **Synchronization**: The main function sends tasks and collects results in a synchronized manner using channels.

This code demonstrates how to create a simple concurrent worker pool in Go, where multiple workers process tasks concurrently and communicate via channels.