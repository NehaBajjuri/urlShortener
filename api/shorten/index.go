package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urlMap = make(map[string]string)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/api/url" {
		handleShorten(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/url/") {
		handleRedirect(w, r)
		return
	}

	// show form for GET
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h2>URL Shortener</h2>
		<form method="post" action="/api/url">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="submit" value="Shorten">
		</form>
	`)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is missing", http.StatusBadRequest)
		return
	}

	key := generateShortKey()
	urlMap[key] = originalURL

	shortURL := fmt.Sprintf("https://<your-vercel>.vercel.app/api/url/%s", key)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h2>Shortened!</h2>
		<p>Original: %s</p>
		<p>Shortened: <a href="%s">%s</a></p>
	`, originalURL, shortURL, shortURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/url/")
	original, ok := urlMap[key]
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
