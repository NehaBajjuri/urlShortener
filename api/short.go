package handler

import (
	"net/http"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}
	shortKey := keys[0]

	originalURL, exists := urlMap[shortKey]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
