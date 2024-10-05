package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"net/http"
	"os"
	"quoteservice/handler"
	"quoteservice/service"
	"quoteservice/store"
)

func main() {
	// Database connection
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Correctly specify the file scheme
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize the store
	quoteStore := &store.SQLQuoteStore{DB: db}
	quoteService := service.QuoteService{QuoteStore: quoteStore}
	quoteHandler := handler.QuoteHandler{QuoteService: quoteService}

	http.HandleFunc("/quotes", quoteHandler.CreateQuoteHandler)
	http.HandleFunc("/quotes/list", quoteHandler.ListQuotesHandler)

	serverPort := os.Getenv("SERVER_PORT")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
