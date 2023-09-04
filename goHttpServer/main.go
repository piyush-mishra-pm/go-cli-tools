package main

import (
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Hi")
	server()
}

func server() {

	http.HandleFunc("/direct-response", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	// Fib of Random Number
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	http.HandleFunc("/fibonacci", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Fib, %q", fibonacciCalculator(r1.Intn(1000_000)))
	})

	fmt.Println("Server listening @ Port 8080:")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Picks a random number between 1 to 10,000 and calculates its Fib.
func fibonacciCalculator(n int) int {
	n0 := 0
	n1 := 1

	if n <= 1 {
		return n
	}

	fib := 1
	for i := 2; i < n; i++ {
		fib = n0 + n1
		n0 = n1
		n1 = fib
	}
	return fib
}
