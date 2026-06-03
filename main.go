package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

var links = make(map[string]string)

func generateCode() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	original := r.FormValue("url")
	if original == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	code := generateCode()
	links[code] = original
	fmt.Fprintf(w, "Short URL: http://localhost:8080/r/%s\n", code)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]
	original, exists := links[code]
	if !exists {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)
	fmt.Println("URL Shortener running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}