package http

import "net/http"

type AuthHandler struct{}

func NewAuthHandler()

func GenerateTokens(w http.ResponseWriter, r *http.Request) {
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
}
