package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {
	// Verify and parse arguments
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column index (1 based) on which to execute operation")
	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileNames []string, operationType string, column int, outWriter io.Writer) error {
	var opFunc statsFunc
	if len(fileNames) == 0 {
		return ErrNoFiles
	}
	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumnNumber, column)
	}

	// Validate the operation and define the opFunc accordingly
	switch operationType {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, operationType)
	}

	consolidatedData := make([]float64, 0)

	// Create the channel to receive results, errors or done status of operations
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	fileNamesCh := make(chan string)

	wg := sync.WaitGroup{}

	// Loop through all files, sending them to filesChannel.
	// When a worker gets free, it will pick up this file.
	go func() {
		defer close(fileNamesCh)
		for _, fileName := range fileNames {
			fileNamesCh <- fileName
		}
	}()

	for fileIndex := 0; fileIndex < runtime.NumCPU(); fileIndex++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileReceivedFromCh := range fileNamesCh {
				// Open the file for reading
				fileOpened, err := os.Open(fileReceivedFromCh)
				if err != nil {
					errCh <- fmt.Errorf("Cannot open file: %w", err)
					return
				}
				// Parse the CSV into a slice of float64 numbers
				data, err := csv2float(fileOpened, column)
				if err != nil {
					errCh <- err
				}
				if err := fileOpened.Close(); err != nil {
					errCh <- err
				}
				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidatedData = append(consolidatedData, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(outWriter, opFunc(consolidatedData))
			return err
		}
	}
}
