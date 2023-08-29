package main

import "errors"

var (
	ErrNaN                 = errors.New("Data is not numeric")
	ErrInvalidColumnNumber = errors.New("Invalid column number")
	ErrNoFiles             = errors.New("No input files")
	ErrInvalidOperation    = errors.New("Invalid operation")
)
