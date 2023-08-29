package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func csv2float(r io.Reader, columnIndex int) ([]float64, error) {
	// Create the CSV Reader used to read data from CSV files
	csvReader := csv.NewReader(r)
	csvReader.ReuseRecord = true
	// Adjusting for 0 based column index
	columnIndex--

	var data []float64

	// Looping through all records
	for rowIndex := 0; ; rowIndex++ {
		// Reading 1 line at a time.
		rowData, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Cant read data from file: %w", err)
		}
		// Skip header
		if rowIndex == 0 {
			continue
		}
		// Checking number of columns in CSV file
		if len(rowData) <= columnIndex {
			// File does not have that many columns
			return nil,
				fmt.Errorf("%w: File has only %d columns", ErrInvalidColumnNumber, len(rowData))
		}
		// Try to convert data read into a float number
		cellValue, err := strconv.ParseFloat(rowData[columnIndex], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNaN, err)
		}
		data = append(data, cellValue)
	}

	return data, nil
}
