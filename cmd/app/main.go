package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

	handlers "github.com/KAA295/medods/api/http"
	_ "github.com/KAA295/medods/docs"
	"github.com/KAA295/medods/repository/postgres"
	"github.com/KAA295/medods/usecases/services"
)

const addr = ":8000"

// @title Auth Service
// @version 1.0
// @description Test task for medods.

// @host localhost:8000
// @BasePath /
func main() {
	psqlInfo := "host=storage port=5432 user=postgres password=password dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	emailService := services.NewEmailService()

	authRepo := postgres.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	r := chi.NewRouter()

	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Post("/generate_tokens", authHandler.GenerateTokens)
	r.Post("/refresh_tokens", authHandler.RefreshTokens)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	httpServer.ListenAndServe()
}
