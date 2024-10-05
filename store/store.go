package store

import (
	"database/sql"
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
