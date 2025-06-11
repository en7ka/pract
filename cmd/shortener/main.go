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
	if r.Method != http.MethodPost || r.URL.Path != "/" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	v, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	orig := string(v)
	sum := sha1.Sum([]byte(orig))
	id := fmt.Sprintf("%x", sum)[:8]
	store[id] = orig
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + id))
}

func handlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.URL.Path == "/" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/")
	orig, ok := store[id]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, orig, http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlerPost(w, r)
		} else {
			handlerGet(w, r)
		}
	})
	http.ListenAndServe(":8080", nil)
}
