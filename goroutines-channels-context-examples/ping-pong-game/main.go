package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Main function
func main() {
	// Initialize a buffered channel to simulate the ping pong ball
	ball := make(chan int, 1)
	// Create a context with a 20-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Start the game by sending the ball to player 1
	ball <- 0

	// Player 1 routine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return // Exit if context is canceled
			case <-ball:
				fmt.Println("Ball is in player 1's court, sending...")
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second) // Random delay before serving
				ball <- 1 // Serve the ball to player 2
			}
		}
	}()

	// Player 2 routine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return // Exit if context is canceled
			case <-ball:
				fmt.Println("Ball is in player 2's court, sending...")
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second) // Random delay before serving
				ball <- 2 // Serve the ball back to player 1
			}
		}
	}()

	// Wait for the game to end or timeout
	select {
	case <-ctx.Done():
		fmt.Println("Game Over! Winner Is Player", <-ball)
	}
}