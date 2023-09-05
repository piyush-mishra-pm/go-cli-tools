package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func sqlLiteTestResponse() {
	time.Sleep(time.Millisecond * 20)
	http.HandleFunc("/sql-lite", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sql-lite, %q", 1)
	})
}
