package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL Shortener Service")
	})

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
