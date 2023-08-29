package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	wg := sync.WaitGroup{}

	// Loop through all files adding their data to consolidate
	for _, fileName := range fileNames {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()

			// Open the file for reading
			f, err := os.Open(fname)
			if err != nil {
				errCh <- fmt.Errorf("Cannot open file: %w", err)
				return
			}
			// Parse the CSV into a slice of float64 numbers
			data, err := csv2float(f, column)
			if err != nil {
				errCh <- err
			}
			if err := f.Close(); err != nil {
				errCh <- err
			}
			resCh <- data
		}(fileName)
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
