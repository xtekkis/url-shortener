package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Link struct {
	Original  string
	Clicks    int
	CreatedAt time.Time
}

var links = make(map[string]Link)

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
	links[code] = Link{Original: original, Clicks: 0, CreatedAt: time.Now()}
	fmt.Fprintf(w, "Short URL: http://localhost:8080/r/%s\n", code)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]
	link, exists := links[code]
	if !exists {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}
	if time.Since(link.CreatedAt) > 24*time.Hour {
		http.Error(w, "Link has expired", http.StatusGone)
		return
	}
	link.Clicks++
	links[code] = link
	http.Redirect(w, r, link.Original, http.StatusFound)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%-10s %-50s %-10s %s\n", "Code", "Original URL", "Clicks", "Expires")
	fmt.Fprintf(w, "%s\n", "---------------------------------------------------------------------------------------------")
	for code, link := range links {
		expires := link.CreatedAt.Add(24 * time.Hour).Format("15:04:05")
		fmt.Fprintf(w, "%-10s %-50s %-10d %s\n", code, link.Original, link.Clicks, expires)
	}
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.HandleFunc("/stats", statsHandler)
	fmt.Println("URL Shortener running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}