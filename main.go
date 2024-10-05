package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/lib/pq"                                // PostgreSQL driver
	"log"
	"net/http"
	"quoteservice/handler"
	"quoteservice/service"
	"quoteservice/store"
)

func main() {
	// Database connection
	dataSourceName := "postgres://quoteuser:password@localhost:5432/quotemanager?sslmode=disable"

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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
