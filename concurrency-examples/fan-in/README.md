# FAN IN

- Fan In is used when a single function reads from multiple inputs and proceeds until all are closed. This is made possible by multiplexing the input into a single channel.

Source: https://kapoorrahul.medium.com/golang-fan-in-fan-out-concurrency-pattern-f5a29ff1f93b

This Go code provides a solution for reading lines from two files concurrently, merging the lines, and printing them to the console. Here's a detailed explanation of each part of the code:

### Package and Imports

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)
```
- `package main`: Defines the package name.
- The imports include packages for buffered I/O (`bufio`), formatted I/O (`fmt`), logging (`log`), OS functions (`os`), and synchronization primitives (`sync`).

### Function `readData`

```go
func readData(file string) <-chan string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string)
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	go func() {
		for fileScanner.Scan() {
			val := fileScanner.Text()
			out <- val
		}

		close(out)

		err := f.Close()
		if err != nil {
			fmt.Printf("Unable to close an opened file: %v\n", err.Error())
			return
		}
	}()

	return out
}
```
- This function takes a filename as input and returns a channel of strings (`<-chan string`).
- It opens the file and creates a `bufio.Scanner` to read the file line by line.
- A goroutine reads each line from the file and sends it to the `out` channel.
- The channel is closed after reading all lines from the file.
- The file is closed, and an error message is logged if the file can't be closed.

### Function `fanInMergeData`

```go
func fanInMergeData(ch1, ch2 <-chan string) chan string {
	chRes := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for val := range ch1 {
			chRes <- val
		}
		wg.Done()
	}()

	go func() {
		for val := range ch2 {
			chRes <- val
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(chRes)
	}()

	return chRes
}
```
- This function takes two input channels (`ch1` and `ch2`) and merges their data into a single output channel (`chRes`).
- A `sync.WaitGroup` is used to ensure both input channels are fully processed before closing the output channel.
- Two goroutines read data from the input channels and send it to the output channel.
- Another goroutine waits for the first two goroutines to finish (using `wg.Wait()`) and then closes the output channel.

### `main` Function

```go
func main() {
	ch1 := readData("text1.txt")
	ch2 := readData("text2.txt")

	chRes := fanInMergeData(ch1, ch2)

	for val := range chRes {
		fmt.Println(val)
	}
}
```
- The `main` function initializes two channels by reading from "text1.txt" and "text2.txt".
- The `fanInMergeData` function is called to merge the data from these two channels.
- The merged data is printed to the console line by line.

### Summary
- **Concurrency**: The code reads from files concurrently using goroutines.
- **Channel Communication**: Channels are used to communicate between goroutines and synchronize the merging process.
- **Synchronization**: A `sync.WaitGroup` ensures all data is read before closing the output channel.

This code demonstrates efficient handling of concurrent I/O operations and merging streams of data in Go.