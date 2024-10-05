package service

import "quoteservice/store"

type QuoteService struct {
	QuoteStore store.QuoteStore
}

func (s *QuoteService) CreateQuote(author, text string) error {
	return s.QuoteStore.SaveQuote(author, text)
}

func (s *QuoteService) ListQuotes() ([]store.Quote, error) {
	return s.QuoteStore.GetAllQuotes()
}
