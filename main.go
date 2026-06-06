package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Link struct {
	Original  string
	Clicks    int
	CreatedAt time.Time
}

var (
	links = make(map[string]Link)
	mu    sync.Mutex
)

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
	if time.Since(link.CreatedAt) > 2*time.Hour {
		http.Error(w, "Link has expired", http.StatusGone)
		return
	}
	link.Clicks++
	links[code] = link
	http.Redirect(w, r, link.Original, http.StatusFound)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Link Stats</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: Arial, sans-serif; background: #0f0f0f; color: white; padding: 40px; }
        h1 { font-size: 24px; margin-bottom: 8px; }
        p { color: #999; margin-bottom: 24px; }
        table { width: 100%%; border-collapse: collapse; }
        th { text-align: left; padding: 12px 16px; background: #1a1a1a; color: #999; font-size: 13px; }
        td { padding: 12px 16px; border-top: 1px solid #2a2a2a; font-size: 14px; }
        td a { color: #818cf8; text-decoration: none; }
        td a:hover { text-decoration: underline; }
        .clicks { color: #4ade80; }
        .expires { color: #999; }
        .back { display: inline-block; margin-bottom: 24px; color: #999; text-decoration: none; font-size: 14px; }
        .back:hover { color: #4f46e5; }
    </style>
</head>
<body>
    <a href="/" class="back">← Back</a>
    <h1>Link Stats</h1>
    <p>All shortened links and their click counts.</p>
    <table>
        <thead>
            <tr>
                <th>Code</th>
                <th>Original Link</th>
                <th>Clicks</th>
                <th>Expires At</th>
            </tr>
        </thead>
        <tbody>`)

	mu.Lock()
	for code, link := range links {
		expires := link.CreatedAt.Add(2 * time.Hour).Format("15:04:05")
		display := link.Original
		if len(display) > 100 {
    		display = display[:100] + "..."
		}
		fmt.Fprintf(w, `<tr>
            <td><a href="/r/%s" target="_blank">%s</a></td>
            <td><a href="%s" target="_blank">%s</a></td>
            <td class="clicks">%d</td>
            <td class="expires">%s</td>
        </tr>`, code, code, link.Original, display, link.Clicks, expires)
	}
	mu.Unlock()

	fmt.Fprintf(w, `</tbody></table></body></html>`)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.HandleFunc("/stats", statsHandler)
	fmt.Println("URL Shortener running on http://localhost:8080")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8080", nil)
}