# Exercises in Go

## FAN IN
- Fan In is used when a single function reads from multiple inputs and proceeds until all are closed. This is made possible by multiplexing the input into a single channel.

## FAN OUT
- Fan out is used when multiple functions read from the same channel. The reading will stop only when the channel is closed. This characteristic is often used to distribute work amongst a group of workers to parallelize the CPU and I /O.

## Fibonacci - Range and close
- A sender can close a channel to indicate that no more values will be sent. Receivers can test whether a channel has been closed by assigning a second parameter to the receive expression: after

## Load balancing
- Channels can be used to distribute workloads among multiple goroutines. By creating a channel to receive tasks, and multiple goroutines to process them, you can achieve load balancing and take advantage of the concurrency of Go

## Ping pong game
- goroutine, channels and context

## Pipeline
- READ, PROCESS AND WRITE

## Race cars
- WaitGroup and goroutines

## Select and channels
- Handling multiple channels with select

## Tasks in order
- Execute tasks in order

## Timeout and cancellation
- Timeout and cancellation: Channels can be used to implement timeouts and cancellation of long-running operations. By using a select statement with a time. After the channel, you can wait for a specified amount of time before proceeding with an operation. Additionally, by using a cancel channel, you can signal to a goroutine that it should stop processing and return early.
