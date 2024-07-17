This Go program processes a text file with rows of pipe-separated values. It extracts and counts specific information in a concurrent manner using goroutines and channels. Here's a breakdown of how it works:

### Main Components:

1. **Data Structures**:
   - `result`: A struct that stores the final processed data, including the number of rows, the number of unique people, the most common first name, and the frequency of donations by month.

2. **Functions**:
   - `processRow(text string) (firstName, fullName, month string)`: This function takes a pipe-separated line of text, extracts the first name, full name, and month from specific columns. It uses string manipulation functions to achieve this.
   - `concurrent(file string, numWorkers, batchSize int) (res result)`: This function orchestrates the concurrent processing of the file. It reads the file, processes it in batches, and combines the results.

3. **Concurrent Processing**:
   - **Reader**: Reads the file in batches and sends the batches to a channel.
   - **Worker**: Processes each batch, extracts relevant information using `processRow`, and sends the processed data to another channel.
   - **Combiner**: Combines the output from multiple workers into a single result.

### Detailed Explanation:

1. **Result Struct**:
   ```go
   type result struct {
       numRows           int
       peopleCount       int
       commonName        string
       commonNameCount   int
       donationMonthFreq map[string]int
   }
   ```

2. **processRow Function**:
   - This function splits a row of text by the pipe (`|`) character.
   - It extracts the full name by trimming spaces and removing internal spaces.
   - It extracts the first name by finding and trimming the appropriate substring.
   - It extracts the month from a date string.
   - Returns the first name, full name, and month.

3. **concurrent Function**:
   - **Initialization**:
     ```go
     res = result{donationMonthFreq: map[string]int{}}
     ```
     Initializes the result struct with an empty map for donation months.
   - **File Opening**:
     ```go
     f, err := os.Open(file)
     if err != nil {
         log.Fatal(err)
     }
     defer f.Close()
     ```
     Opens the file and defers its closure.

4. **Reader Function**:
   - Reads the file line by line.
   - Collects lines into batches.
   - Sends each batch to a channel for processing.

5. **Worker Function**:
   - Processes each batch received from the reader.
   - Uses `processRow` to extract data from each line.
   - Sends the processed data to another channel.

6. **Combiner Function**:
   - Multiplexes the output of multiple workers.
   - Waits for all workers to finish.
   - Sends the combined result to a channel.

7. **Main Function**:
   - Sets up the file path, number of workers, and batch size.
   - Calls the `concurrent` function with these parameters.
   - Prints the results.

### Example Execution:
```go
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
```
- This sets the file path, the number of worker goroutines, and the batch size.
- Calls `concurrent` to process the file.
- Prints the number of rows processed, the number of unique people, the most common first name, and the frequency of donations by month.