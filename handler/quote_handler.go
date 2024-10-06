package handler

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(err)
		return
	}

	if quote.Author == "" || quote.Text == "" {
		http.Error(w, "Missing quote author or text", http.StatusBadRequest)
		return
	}

	if err := h.QuoteService.CreateQuote(quote.Author, quote.Text); err != nil {
		http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *QuoteHandler) ListQuotesHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := h.QuoteService.ListQuotes()

	if err != nil {
		http.Error(w, "Failed to list quotes", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *QuoteHandler) RetrieveRandomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := h.QuoteService.RetrieveRandomQuote()

	if err != nil {
		http.Error(w, "Failed to retrieve random quote", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *QuoteHandler) DestroyQuoteHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		ID int `json:"id"`
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Println("Error decoding request payload:", err)
		return
	}

	if payload.ID <= 0 {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	err := h.QuoteService.DeleteQuote(payload.ID)
	if err != nil {
		if err.Error() == "quote not found" {
			http.Error(w, "Quote not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete quote", http.StatusInternalServerError)
			fmt.Println("Error deleting quote:", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
