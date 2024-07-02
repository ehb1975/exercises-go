package main

import "fmt"

func main() {
	// Create channels to connect the different stages of the pipeline
	nums := make(chan int)
	squares := make(chan int)

	// Start a goroutine to generate a sequence of numbers
	go func() {
		for i := 1; i <= 10; i++ {
			nums <- i
		}
		//The close built-in function closes a channel, which must be either bidirectional or send-only.
		close(nums)
	}()

	// Start a goroutine to square each number in the sequence
	go func() {
		for num := range nums {
			squares <- num * num
		}
		//The close built-in function closes a channel, which must be either bidirectional or send-only.
		close(squares)
	}()
	// Read the squared numbers from the channel and print them
	for square := range squares {
		fmt.Println(square)
	}
}
