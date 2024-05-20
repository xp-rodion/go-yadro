package handlers

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(queryParams)
	if proposal := queryParams.Get(key); queryParams.Has(key) && proposal != "" {
		w.WriteHeader(http.StatusOK)
		comics, err := h.comicsService.List(proposal)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(comics)
		if err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := h.comicsService.Update()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}

func (h *HTTPHandler) Init() {
	http.HandleFunc("GET /pics", h.List)
	http.HandleFunc("POST /update", h.Update)
}
