package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
)

var store = map[string]string{}

func handlers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		v, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		orig := string(v)
		sum := sha1.Sum([]byte(orig))
		id := fmt.Sprintf("%x", sum)[:8]
		store[id] = orig
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + id))
		return
	}
}

func main() {
	http.HandleFunc("/", handlers)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
