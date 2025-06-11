package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var store = map[string]string{}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	fmt.Println("Received URL:", originalURL)

	sum := sha1.Sum([]byte(originalURL))
	id := strings.ToUpper(fmt.Sprintf("%x", sum)[:8])
	store[id] = originalURL
	fmt.Println("Generated ID:", id)
	fmt.Println("Current store:", store)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://localhost:8080/%s", id)
}

func handlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, ok := store[id]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	fmt.Println("Redirecting to original URL:", originalURL)

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlerPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/{id}", handlerGet) // This will handle GET requests for the shortened URLs


	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
