package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// URLStore represents our in-memory URL storage
type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

// NewURLStore creates a new URLStore instance
func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

// generateShortURL creates a random 6-character string
func generateShortURL() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:6], nil
}

// Put stores a URL and returns its short version
func (s *URLStore) Put(url string) (string, error) {
	shortURL, err := generateShortURL()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.urls[shortURL] = url
	s.mu.Unlock()

	return shortURL, nil
}

// Get retrieves the original URL for a given short URL
func (s *URLStore) Get(shortURL string) (string, bool) {
	s.mu.RLock()
	url, exists := s.urls[shortURL]
	s.mu.RUnlock()
	return url, exists
}

func main() {
	store := NewURLStore()

	// CORS middleware
	setCORS := func(w http.ResponseWriter, r *http.Request) bool {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return true
		}
		return false
	}

	// Handle URL shortening
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		longURL := r.FormValue("url")
		if longURL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		shortURL, err := store.Put(longURL)
		if err != nil {
			http.Error(w, "Error generating short URL", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Short URL: http://localhost:8080/%s", shortURL)
	})

	// Handle redirects
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w, r)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}

		shortURL := r.URL.Path[1:] // Remove leading slash
		longURL, exists := store.Get(shortURL)
		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	})

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 