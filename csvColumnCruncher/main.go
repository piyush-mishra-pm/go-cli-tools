package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

func run(filenames []string, operationType string, column int, outWriter io.Writer) error {
	var opFunc statsFunc
	if len(filenames) == 0 {
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
	// Loop through all files adding their data to consolidate
	for _, fname := range filenames {
		// Open the file for reading
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("Cannot open file: %w", err)
		}
		// Parse the CSV into a slice of float64 numbers
		data, err := csv2float(f, column)
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		// Append the data to consolidate
		consolidatedData = append(consolidatedData, data...)
	}
	_, err := fmt.Fprintln(outWriter, opFunc(consolidatedData))
	return err
}
