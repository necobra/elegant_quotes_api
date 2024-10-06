package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Quote struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type QuoteStore interface {
	SaveQuote(author, text string) error
	GetAllQuotes() ([]Quote, error)
	GetQuote(id int) (*Quote, error)
	GetRandomQuote() (*Quote, error)
	DeleteQuote(id int) error
}

type SQLQuoteStore struct {
	DB *sql.DB
}

func (s *SQLQuoteStore) SaveQuote(author, text string) error {
	_, err := s.DB.Exec("INSERT INTO quotes (author, text) VALUES ($1, $2)", author, text)
	return err
}

func (s *SQLQuoteStore) GetAllQuotes() ([]Quote, error) {
	rows, err := s.DB.Query("SELECT id, author, text, created_at FROM quotes")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var quotes []Quote

	for rows.Next() {
		var quote Quote

		if err := rows.Scan(&quote.ID, &quote.Author, &quote.Text, &quote.CreatedAt); err != nil {
			return nil, err
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (s *SQLQuoteStore) GetQuote(id int) (*Quote, error) {
	row, err := s.DB.Query("SELECT * FROM quotes WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var quote Quote
	row.Next()
	if err := row.Scan(&quote.ID, &quote.Author, &quote.Text, &quote.CreatedAt); err != nil {
		return nil, err
	}
	return &quote, nil
}

func (s *SQLQuoteStore) GetRandomQuote() (*Quote, error) {
	row, err := s.DB.Query(
		"SELECT * FROM quotes LIMIT 1 OFFSET FLOOR(RANDOM() * (SELECT COUNT(*) FROM quotes))",
	)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var quote Quote
	row.Next()
	if err := row.Scan(&quote.ID, &quote.Author, &quote.Text, &quote.CreatedAt); err != nil {
		return nil, err
	}
	return &quote, nil
}

func (s *SQLQuoteStore) DeleteQuote(id int) error {
	res, err := s.DB.Exec("DELETE FROM quotes WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quote with id %d not found", id)
	}

	return nil
}
