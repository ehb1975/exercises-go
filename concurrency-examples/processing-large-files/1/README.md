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

4.1. **Assinatura da Função**:
   ```go
   func(ctx context.Context, rowsBatch *[]string) <-chan []string
   ```
   - `ctx context.Context`: É um contexto usado para controlar a execução da função, permitindo, por exemplo, cancelar a leitura de forma graciosa.
   - `rowsBatch *[]string`: É um ponteiro para uma fatia de strings que representa um lote de linhas a serem processadas.
   - `<-chan []string`: A função retorna um canal que envia fatias de strings (lotes de linhas).

4.2. **Criação do Canal de Saída**:
   ```go
   out := make(chan []string)
   ```
   - Um canal do tipo `[]string` é criado. Este canal será usado para enviar os lotes de linhas lidas do arquivo.

4.3. **Inicialização do Scanner**:
   ```go
   scanner := bufio.NewScanner(f)
   ```
   - Um scanner é criado usando `bufio.NewScanner(f)`, onde `f` é o arquivo de onde as linhas serão lidas. Este scanner lerá o arquivo linha por linha.

4.4. **Goroutine**:
   ```go
   go func() {
       defer close(out)
       // ...
   }()
   ```
   - A função anônima é executada em uma nova goroutine para ler e processar os dados de forma assíncrona.
   - `defer close(out)`: O canal de saída será fechado quando a função anônima terminar a execução.

4.5. **Leitura e Processamento das Linhas**:
   ```go
   for {
       scanned := scanner.Scan()
       select {
       case <-ctx.Done():
           return
       default:
           row := scanner.Text()
           if len(*rowsBatch) == batchSize || !scanned {
               out <- *rowsBatch
               *rowsBatch = []string{}
           }
           *rowsBatch = append(*rowsBatch, row)
       }

       if !scanned {
           return
       }
   }
   ```
   - A função entra em um loop infinito para ler linhas do arquivo.
   - `scanned := scanner.Scan()`: `scanner.Scan()` lê a próxima linha do arquivo e retorna `true` se houver mais linhas a serem lidas.
   - `select` é usado para verificar se o contexto (`ctx`) foi cancelado (`ctx.Done()`). Se o contexto for cancelado, a goroutine retorna, encerrando a execução.
   - `row := scanner.Text()`: `scanner.Text()` retorna a linha atual lida pelo scanner.
   - `if len(*rowsBatch) == batchSize || !scanned`: Se o tamanho do lote (`rowsBatch`) atingir o tamanho máximo (`batchSize`) ou se o arquivo terminou (`!scanned`), o lote atual é enviado pelo canal `out`.
   - `out <- *rowsBatch`: O lote de linhas é enviado pelo canal `out`.
   - `*rowsBatch = []string{}`: O lote é reinicializado.
   - `*rowsBatch = append(*rowsBatch, row)`: A linha atual (`row`) é adicionada ao lote.

4.6. **Verificação de Fim de Arquivo**:
   ```go
   if !scanned {
       return
   }
   ```
   - Se não houver mais linhas a serem lidas (`!scanned`), a função retorna, encerrando a goroutine.

4.7. **Retorno do Canal de Saída**:
   ```go
   return out
   ```
   - A função retorna o canal `out`, que eventualmente conterá os lotes de linhas lidas do arquivo.

Resumindo, essa função lê linhas de um arquivo de texto e as envia em lotes através de um canal de forma assíncrona. O tamanho do lote é controlado pela variável `batchSize`, e a função pode ser cancelada a qualquer momento usando o contexto (`ctx`).


5. **Worker Function**:
   - Processes each batch received from the reader.
   - Uses `processRow` to extract data from each line.
   - Sends the processed data to another channel.

Define uma função que processa dados de forma assíncrona. Vou detalhar cada parte:

5.1. **Assinatura da Função**:
   ```go
   func(ctx context.Context, rowBatch <-chan []string) <-chan processed
   ```
   - `ctx context.Context`: É um contexto que pode ser usado para controlar a execução da função, como para cancelar a execução de forma graciosa.
   - `rowBatch <-chan []string`: É um canal de entrada que recebe lotes de strings (`[]string`). Cada lote é uma fatia de strings.
   - `<-chan processed`: A função retorna um canal de saída que envia objetos do tipo `processed`.

5.2. **Criação do Canal de Saída**:
   ```go
   out := make(chan processed)
   ```
   - Um canal do tipo `processed` é criado. Este canal será usado para enviar os resultados processados.

5.3. **Goroutine**:
   ```go
   go func() {
       defer close(out)
       // ...
   }()
   ```
   - A função anônima é executada em uma nova goroutine (uma thread leve) para processar os dados de forma assíncrona.
   - `defer close(out)`: O canal de saída será fechado quando a função anônima terminar a execução.

5.4. **Inicialização da Estrutura `processed`**:
   ```go
   p := processed{}
   ```
   - Um objeto `processed` é inicializado. Este objeto irá acumular os resultados do processamento.

5.5. **Processamento dos Dados**:
   ```go
   for rowBatch := range rowBatch {
       for _, row := range rowBatch {
           firstName, fullName, month := processRow(row)
           p.fullNames = append(p.fullNames, fullName)
           p.firstNames = append(p.firstNames, firstName)
           p.months = append(p.months, month)
           p.numRows++
       }
   }
   ```
   - A função lê cada lote de strings do canal `rowBatch`.
   - Para cada lote, ela itera sobre cada string (`row`) e chama a função `processRow(row)`, que supostamente retorna `firstName`, `fullName` e `month`.
   - Esses valores são adicionados às fatias correspondentes (`p.fullNames`, `p.firstNames`, `p.months`) e o contador de linhas (`p.numRows`) é incrementado.

5.6. **Envio do Resultado Processado**:
   ```go
   out <- p
   ```
   - Após o processamento de todos os lotes, o objeto `processed` (`p`) é enviado para o canal de saída `out`.

5.7. **Retorno do Canal de Saída**:
   ```go
   return out
   ```
   - A função retorna o canal de saída `out`, que eventualmente conterá o objeto `processed` com os resultados.

Resumindo, essa função processa lotes de strings de forma assíncrona, acumulando resultados em uma estrutura `processed`, e envia esses resultados através de um canal de saída. A função `processRow` é responsável por processar cada string individualmente e extrair as informações relevantes.

6. **Combiner Function**:
   - Multiplexes the output of multiple workers.
   - Waits for all workers to finish.
   - Sends the combined result to a channel.

Define uma função que combina múltiplos canais de entrada em um único canal de saída. Vamos detalhar cada parte:

6.1. **Assinatura da Função**:
   ```go
   func(ctx context.Context, inputs ...<-chan processed) <-chan processed
   ```
   - `ctx context.Context`: É um contexto que pode ser usado para controlar a execução da função, permitindo, por exemplo, cancelar a operação de forma graciosa.
   - `inputs ...<-chan processed`: É um parâmetro variádico que recebe múltiplos canais de entrada do tipo `processed`.
   - `<-chan processed`: A função retorna um canal de saída que envia objetos do tipo `processed`.

6.2. **Criação do Canal de Saída**:
   ```go
   out := make(chan processed)
   ```
   - Um canal do tipo `processed` é criado. Este canal será usado para enviar os dados combinados dos canais de entrada.

6.3. **Criação do WaitGroup**:
   ```go
   var wg sync.WaitGroup
   ```
   - Um `WaitGroup` é criado para esperar até que todas as goroutines sejam concluídas.

6.4. **Função Multiplexadora**:
   ```go
   multiplexer := func(p <-chan processed) {
       defer wg.Done()

       for in := range p {
           select {
           case <-ctx.Done():
           case out <- in:
           }
       }
   }
   ```
   - A função `multiplexer` recebe um canal de entrada `p` do tipo `processed`.
   - `defer wg.Done()`: Quando a função `multiplexer` terminar, ela decrementar o contador do `WaitGroup`.
   - A função lê dados do canal `p` e os envia para o canal `out`.
   - `select` é usado para verificar se o contexto (`ctx`) foi cancelado (`ctx.Done()`). Se o contexto for cancelado, a função retorna. Caso contrário, o dado é enviado para o canal `out`.

6.5. **Adição de Goroutines ao WaitGroup**:
   ```go
   wg.Add(len(inputs))
   for _, in := range inputs {
       go multiplexer(in)
   }
   ```
   - `wg.Add(len(inputs))`: O contador do `WaitGroup` é incrementado pelo número de canais de entrada.
   - Para cada canal de entrada em `inputs`, uma nova goroutine é iniciada executando a função `multiplexer`.

6.6. **Goroutine para Fechar o Canal de Saída**:
   ```go
   go func() {
       wg.Wait()
       close(out)
   }()
   ```
   - Uma goroutine é iniciada para esperar até que todas as goroutines iniciadas anteriormente terminem (`wg.Wait()`).
   - Quando todas as goroutines terminarem, o canal de saída `out` é fechado (`close(out)`).

6.7. **Retorno do Canal de Saída**:
   ```go
   return out
   ```
   - A função retorna o canal de saída `out`, que eventualmente conterá os dados combinados de todos os canais de entrada.

Resumindo, essa função combina múltiplos canais de entrada em um único canal de saída de forma assíncrona. Cada canal de entrada é processado por uma goroutine que envia os dados para o canal de saída. O `WaitGroup` garante que o canal de saída só será fechado quando todos os canais de entrada tiverem sido completamente processados. O contexto (`ctx`) permite cancelar a operação de forma graciosa, interrompendo a leitura dos canais de entrada se necessário.

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