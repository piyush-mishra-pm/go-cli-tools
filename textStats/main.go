package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	isLineMode := flag.Bool("l", false, "Count lines")
	flag.Parse()
	fmt.Print(countWords(os.Stdin, *isLineMode))
	if *isLineMode {
		fmt.Print(" lines")
	} else {
		fmt.Print(" words")
	}
}

func countWords(r io.Reader, isLineMode bool) int {
	scanner := bufio.NewScanner(r)

	if !isLineMode {
		scanner.Split(bufio.ScanWords)
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}
