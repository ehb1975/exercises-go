package main

import (
	"fmt"
	"time"
)

func main() {
	// Create two channels to synchronize the execution of goroutines
	ch1 := make(chan bool)
	ch2 := make(chan bool)

	// Start a goroutine to execute task1
	go func() {
		// Execute task1
		fmt.Println("Task 1 executed")
		// Signal that task1 is done by writing to channel ch1
		ch1 <- true
	}()

	// Start a goroutine to execute task2
	go func() {
		// Wait for task1 to complete by reading from channel ch1
		<-ch1
		// Execute task2
		fmt.Println("Task 2 executed")
		// Signal that task2 is done by writing to channel ch2
		ch2 <- true
	}()

	// Start a goroutine to execute task3
	go func() {
		// Wait for task2 to complete by reading from channel ch2
		<-ch2
		// Execute task3
		fmt.Println("Task 3 executed")
	}()

	// Wait for all tasks to complete
	time.Sleep(time.Second)
}
