package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	handlers "github.com/KAA295/medods/api/http"
	"github.com/KAA295/medods/repository/postgres"
	"github.com/KAA295/medods/usecases/services"
)

const addr = ":8000"

func main() {
	psqlInfo := fmt.Sprintf(
		"host=storage port=5432 user=postgres password=password dbname=postgres sslmode=disable",
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	authRepo := postgres.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := chi.NewRouter()

	r.Post("/generate_tokens", authHandler.GenerateTokens)
	r.Post("/refresh_tokens", authHandler.RefreshTokens)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	httpServer.ListenAndServe()
}
