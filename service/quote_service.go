package service

import (
	"fmt"
	"quoteservice/store"
)

type QuoteService struct {
	QuoteStore store.QuoteStore
}

func (s *QuoteService) CreateQuote(author, text string) error {
	return s.QuoteStore.SaveQuote(author, text)
}

func (s *QuoteService) ListQuotes() ([]store.Quote, error) {
	return s.QuoteStore.GetAllQuotes()
}

func (s *QuoteService) RetrieveRandomQuote() (*store.Quote, error) {
	return s.QuoteStore.GetRandomQuote()
}

func (s *QuoteService) DeleteQuote(id int) error {
	err := s.QuoteStore.DeleteQuote(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("quote with id %d not found", id) {
			return fmt.Errorf("quote not found")
		}
		return err
	}
	return nil
}
