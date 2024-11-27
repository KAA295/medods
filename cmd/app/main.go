package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const addr = ":8000"

func main() {
	r := chi.NewRouter()

	r.Post("/generate_tokens", nil)
	r.Post("/refresh_tokens", nil)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	httpServer.ListenAndServe()
}
