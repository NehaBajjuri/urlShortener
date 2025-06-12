package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var urlMap = make(map[string]string)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<h2>URL Shortener</h2>
			<form method="post">
				<input type="text" name="url" placeholder="Enter a URL">
				<input type="submit" value="Shorten">
			</form>
		`)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	urlMap[shortKey] = url

	shortURL := fmt.Sprintf("/api/short?key=%s", shortKey)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h2>URL Shortened!</h2>
		<p><strong>Original:</strong> %s</p>
		<p><strong>Shortened:</strong> <a href="%s">%s</a></p>
	`, url, shortURL, shortURL)
}

func generateShortKey() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, 6)
	for i := range key {
		key[i] = chars[rand.Intn(len(chars))]
	}
	return string(key)
}
