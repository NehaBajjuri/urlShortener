package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	urls map[string]string
}

func (res *URLShortener) Handleshorten(wr http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		wr.Header().Set("Content-Type", "text/html")
		// Render the form HTML when accessed via GET
		fmt.Fprintf(wr, `
		<h2>URL Shortener</h2>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="submit" value="Shorten">
		</form>
		`)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(wr, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(wr, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	res.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8081/short/%s", shortKey)

	//render the html res with the shortened url

	wr.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
	<h2>URL Shortener</h2>
	<p>Original URL: %s</p>
	<p>Shortened URL: <a href="%s">%s</a></p>
	<form method="post" action = "/shorten">
	 <input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
        </form>
    `, originalURL, shortenedURL, shortenedURL)
	fmt.Fprint(wr, responseHTML)

}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6
	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	fmt.Printf("Received request to redirect for short key: %s\n", shortKey)
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		fmt.Println("Key not found in URL map:", shortKey)
		return
	}

	// Redirect the user to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
func main() {
	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("/shorten", shortener.Handleshorten)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	fmt.Println("URL Shortener is running on :8081")
	http.ListenAndServe(":8081", nil)
}
