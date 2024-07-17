package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type result struct {
	numRows           int
	peopleCount       int
	commonName        string
	commonNameCount   int
	donationMonthFreq map[string]int
}

// processRow takes a pipe-separated line and returns the firstName, fullName, and month.
// this function was created to be somewhat compute intensive and not accurate.
func processRow(text string) (firstName, fullName, month string) {
	row := strings.Split(text, "|")

	//extract full name
	fullName = strings.Replace(strings.TrimSpace(row[7]), " ", "", -1)

	//extract first name
	name := strings.TrimSpace(row[7])
	if name != "" {
		startOfName := strings.Index(name, ", ") + 2
		if endOfName := strings.Index(name[startOfName:], " "); endOfName < 0 {
			firstName = name[startOfName:]
		} else {
			firstName = name[startOfName : startOfName+endOfName]
		}
		if strings.HasSuffix(firstName, ",") {
			firstName = strings.Replace(firstName, ",", "", -1)
		}
	}

	//extract month
	date := strings.TrimSpace(row[13])
	if len(date) == 8 {
		month = date[:2]
	} else {
		month = "--"
	}

	return firstName, fullName, month
}

func concurrent(file string, numWorkers, batchSize int) (res result) {

	res = result{donationMonthFreq: map[string]int{}}

	type processed struct {
		numRows    int
		fullNames  []string
		firstNames []string
		months     []string
	}

	//open file
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	// reader creates and returns a channel that recieves
	// batches of rows (of length batchSize) from the file
	reader := func(ctx context.Context, rowsBatch *[]string) <-chan []string {
		out := make(chan []string)

		scanner := bufio.NewScanner(f)

		go func() {
			// close channel when we are done sending all rows
			defer close(out)

			for {
				scanned := scanner.Scan()

				select {
				case <-ctx.Done():
					return
				default:
					row := scanner.Text()
					// if batch size is complete or end of file, send batch out
					if len(*rowsBatch) == batchSize || !scanned {
						out <- *rowsBatch
						// clear batch
						*rowsBatch = []string{}
					}
					*rowsBatch = append(*rowsBatch, row)
				}

				// if nothing else to scan return
				if !scanned {
					return
				}
			}
		}()

		return out
	}

	// worker takes in a read-only channel to recieve batches of rows.
	// After it processes each row-batch it sends out the processed output
	// on its channel.
	worker := func(ctx context.Context, rowBatch <-chan []string) <-chan processed {
		out := make(chan processed)

		go func() {
			defer close(out)

			p := processed{}
			for rowBatch := range rowBatch {
				for _, row := range rowBatch {
					firstName, fullName, month := processRow(row)
					p.fullNames = append(p.fullNames, fullName)
					p.firstNames = append(p.firstNames, firstName)
					p.months = append(p.months, month)
					p.numRows++
				}
			}
			out <- p
		}()

		return out
	}

	// combiner takes in multiple read-only channels that receive processed output
	// (from workers) and sends it out on it's own channel via a multiplexer.
	combiner := func(ctx context.Context, inputs ...<-chan processed) <-chan processed {
		out := make(chan processed)

		var wg sync.WaitGroup
		multiplexer := func(p <-chan processed) {
			defer wg.Done()

			for in := range p {
				select {
				case <-ctx.Done():
				case out <- in:
				}
			}
		}

		// add length of input channels to be consumed by mutiplexer
		wg.Add(len(inputs))
		for _, in := range inputs {
			go multiplexer(in)
		}

		// close channel after all inputs channels are closed
		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	// create a main context, and call cancel at the end, to ensure all our
	// goroutines exit without leaving leaks.
	// Particularly, if this function becomes part of a program with
	// a longer lifetime than this function.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// STAGE 1: start reader
	rowsBatch := []string{}
	rowsCh := reader(ctx, &rowsBatch)

	// STAGE 2: create a slice of processed output channels with size of numWorkers
	// and assign each slot with the out channel from each worker.
	workersCh := make([]<-chan processed, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workersCh[i] = worker(ctx, rowsCh)
	}

	firstNameCount := map[string]int{}
	fullNameCount := map[string]bool{}

	// STAGE 3: read from the combined channel and calculate the final result.
	// this will end once all channels from workers are closed!
	for processed := range combiner(ctx, workersCh...) {
		// add number of rows processed by worker
		res.numRows += processed.numRows

		// add months processed by worker
		for _, month := range processed.months {
			res.donationMonthFreq[month]++
		}

		// use full names to count people
		for _, fullName := range processed.fullNames {
			fullNameCount[fullName] = true
		}
		res.peopleCount = len(fullNameCount)

		// update most common first name based on processed results
		for _, firstName := range processed.firstNames {
			firstNameCount[firstName]++

			if firstNameCount[firstName] > res.commonNameCount {
				res.commonName = firstName
				res.commonNameCount = firstNameCount[firstName]
			}
		}
	}

	return res
}

// main function to run the concurrent processing
func main() {
	file := "./arquivos/sample.txt"

	numWorkers := 4
	batchSize := 100

	res := concurrent(file, numWorkers, batchSize)

	fmt.Printf("Number of rows processed: %d\n", res.numRows)
	fmt.Printf("Number of unique people: %d\n", res.peopleCount)
	fmt.Printf("Most common first name: %s (Count: %d)\n", res.commonName, res.commonNameCount)
	fmt.Println("Donation month frequencies:")
	for month, freq := range res.donationMonthFreq {
		fmt.Printf("%s: %d\n", month, freq)
	}
}
