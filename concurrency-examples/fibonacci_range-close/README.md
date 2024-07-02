# FIBONACCI

## Range and Close

- A sender can close a channel to indicate that no more values will be sent. Receivers can test whether a channel has been closed by assigning a second parameter to the receive expression: after

Source: https://go.dev/tour/concurrency/4

This Go code generates the Fibonacci sequence using a goroutine and channels. Here's a detailed explanation of each part of the code:

### Package and Imports

```go
package main

import (
	"fmt"
)
```
- `package main`: Defines the package name.
- The import statement includes the `fmt` package for formatted I/O.

### Function `fibonacci`

```go
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
```
- This function generates the first `n` Fibonacci numbers and sends them to the channel `c`.
- It takes two parameters:
  - `n`: the number of Fibonacci numbers to generate.
  - `c`: a channel through which the generated Fibonacci numbers are sent.
- Inside the function:
  - Two variables, `x` and `y`, are initialized to the first two numbers in the Fibonacci sequence (0 and 1).
  - A `for` loop runs `n` times. In each iteration:
    - The current value of `x` is sent to the channel `c`.
    - The next Fibonacci number is computed by updating `x` and `y`.
  - After the loop completes, the channel `c` is closed to signal that no more values will be sent.

### `main` Function

```go
func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
```
- The `main` function coordinates the entire program.
- It creates a buffered channel `c` with a capacity of 10 (`make(chan int, 10)`).
- A new goroutine is started to run the `fibonacci` function, passing the channel's capacity (`cap(c)`) as the number of Fibonacci numbers to generate and the channel `c` for sending the numbers.
- The `for` loop iterates over the values received from the channel `c` and prints each value. This loop will automatically terminate when the channel `c` is closed.

### Summary
- **Concurrency**: The code uses a goroutine to run the `fibonacci` function concurrently.
- **Channel**: A buffered channel is used to communicate between the `fibonacci` function and the `main` function.
- **Fibonacci Generation**: The `fibonacci` function generates the Fibonacci sequence and sends each number to the channel.
- **Range Loop**: The `main` function uses a `for` loop to receive values from the channel and print them until the channel is closed.

This code demonstrates how to generate and print a sequence of Fibonacci numbers using goroutines and channels in Go.
