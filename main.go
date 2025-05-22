package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

// PostgresStore represents our PostgreSQL URL storage
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgresStore instance
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// Create the urls table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (short_code TEXT PRIMARY KEY, original_url TEXT)`)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
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
func (s *PostgresStore) Put(url string) (string, error) {
	shortURL, err := generateShortURL()
	if err != nil {
		return "", err
	}
	_, err = s.db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortURL, url)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

// Get retrieves the original URL for a given short URL
func (s *PostgresStore) Get(shortURL string) (string, bool) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortURL).Scan(&originalURL)
	if err != nil {
		return "", false
	}
	return originalURL, true
}

func main() {
	// Replace with your actual PostgreSQL connection string
	connStr := "postgres://username:password@localhost:5432/urlshortener?sslmode=disable"
	store, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers for every response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method not allowed")
			return
		}

		longURL := r.FormValue("url")
		if longURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "URL is required")
			return
		}

		shortURL, err := store.Put(longURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error generating short URL")
			return
		}

		fmt.Fprintf(w, "Short URL: https://url-shortener-production-8c28.up.railway.app/%s", shortURL)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers for every response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}

		shortURL := r.URL.Path[1:] // Remove leading slash
		longURL, exists := store.Get(shortURL)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "URL not found")
			return
		}

		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	})

	// CORS middleware wraps the mux
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		mux.ServeHTTP(w, r)
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
