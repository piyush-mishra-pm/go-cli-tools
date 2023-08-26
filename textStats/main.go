package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Flags:
	isLineMode := flag.Bool("l", false, "Count lines")
	isByteMode := flag.Bool("b", false, "Count bytes")
	flag.Parse()

	// Execute count:
	fmt.Print(count(os.Stdin, *isLineMode, *isByteMode))

	// Print suffix to result.
	if *isLineMode {
		fmt.Print(" lines")
	} else if *isByteMode {
		fmt.Print(" bytes")
	} else {
		fmt.Print(" words")
	}
}

func count(r io.Reader, isLineMode bool, isByteMode bool) int {
	scanner := bufio.NewScanner(r)

	if !isLineMode && !isByteMode {
		scanner.Split(bufio.ScanWords)
	} else if isByteMode {
		scanner.Split(bufio.ScanBytes)
	} // default isLineMode

	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
