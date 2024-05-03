package handlers

import (
	"encoding/json"
	"net/http"
	"xkcd/internal/core/ports"
)

type HTTPHandler struct {
	comicsService ports.ComicsService
}

func NewHTTPHandler(comicsService ports.ComicsService) *HTTPHandler {
	return &HTTPHandler{
		comicsService: comicsService,
	}
}

func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	key := "search"
	w.Header().Set("Content-Type", "application/json")
	if proposal := queryParams.Get(key); queryParams.Has(key) && proposal != "" {
		w.WriteHeader(http.StatusOK)
		comics, err := h.comicsService.List(proposal)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(comics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		response, err := h.comicsService.Update()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) Init() {
	http.HandleFunc("/pics", h.List)
	http.HandleFunc("/update", h.Update)
}
