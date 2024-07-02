package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, tasks <-chan int, results chan<- int) {
	// Process tasks from the input channel
	for task := range tasks {
		fmt.Printf("Worker %d started task %d\n", id, task)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second) // Simulate processing time
		results <- task * 2                                   // Send the result to the output channel
		fmt.Printf("Worker %d finished task %d\n", id, task)
	}
}

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
