package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func genCode(length int) string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)

}

func createUrl(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode := genCode(6)

	_, err := db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, requestBody.URL)
	if err != nil {
		http.Error(w, `{"error": "Failed to create short URL"}`, http.StatusInternalServerError)
		log.Printf("Database insertion error: %v", err)
		return
	}

	response := map[string]string{
		"short_url": "http://localhost:8080/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {

	shortCode := chi.URLParam(r, "shortCode")

	var originalUrl string

	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalUrl)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)

}

func main() {

	connStr := "host=localhost port=5432 user=myuser password=mypassword dbname=url_db sslmode=disable"

	var err error
	db, err = sql.Open("pgx", connStr)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS urls (
        short_code VARCHAR(10) PRIMARY KEY,
        original_url TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );`

	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating urls table: %v", err)
	}

	log.Println("Database connection successful and table checked.")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/shorten", createUrl)
	r.Get("/{shortCode:[a-zA-Z0-9]+}", handleRedirect)

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/*", fs)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
