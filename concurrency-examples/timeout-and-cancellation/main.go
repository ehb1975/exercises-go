package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a channel to receive the result
	result := make(chan int)

	// Start a goroutine to perform a long-running operation
	go func() {
		time.Sleep(3 * time.Second) // Simulate a long-running operation
		result <- 42
	}()

	// Wait for the result, but timeout if it takes too long
	select {
	case res := <-result:
		fmt.Println("Result:", res)
	case <-time.After(2 * time.Second):
		fmt.Println("Timed out!")
	}
}
