package handler

import (
	"encoding/json"
	"net/http"
	"quoteservice/service"
)

type QuoteHandler struct {
	QuoteService service.QuoteService
}

func (h *QuoteHandler) CreateQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var quote struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if quote.Author == "" || quote.Text == "" {
		http.Error(w, "Missing quote author or text", http.StatusBadRequest)
		return
	}

	if err := h.QuoteService.CreateQuote(quote.Author, quote.Text); err != nil {
		http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *QuoteHandler) ListQuotesHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := h.QuoteService.ListQuotes()

	if err != nil {
		http.Error(w, "Failed to list quotes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}
