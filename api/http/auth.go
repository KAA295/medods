package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/beevik/guid"
	"github.com/go-chi/render"

	"github.com/KAA295/medods/api/types"
	"github.com/KAA295/medods/pkg"
	"github.com/KAA295/medods/usecases"
)

type authHandler struct {
	service usecases.AuthService
}

func NewAuthHandler(authService usecases.AuthService) *authHandler {
	return &authHandler{service: authService}
}

// @Summary Generate Tokens
// @Description Generate access_token and refresh_token
// @Param json-body body types.GenerateTokensRequest true "guid"
// @Success 200 {object} types.TokensResponse "Token"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 401 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /generate_tokens [post]
func (h *authHandler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	var req types.GenerateTokensRequest
	ip := r.RemoteAddr
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "invalid request"})
		return
	}
	if req.UserID == "" {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "missing guid"})
		return
	}
	if !guid.IsGuid(req.UserID) {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "guid is not valid"})
		return
	}

	tokens, err := h.service.GenerateTokens(req.UserID, ip)
	if err != nil {
		pkg.ProcessError(w, r, pkg.ErrorResponse{Message: err.Error(), Err: err})
		return
	}

	resp := types.TokensResponse{
		AccessToken:  tokens.AccessToken.Token,
		RefreshToken: tokens.RefreshToken.Token,
	}

	render.Status(r, http.StatusOK)

	render.JSON(w, r, resp)
}

// @Summary Refresh Tokens
// @Description Refresh access_token and refresh_token
// @Param json-body body types.RefreshTokensRequest true "guid and refresh_token"
// @Param Authorization header string true "access_token with Bearer"
// @Success 200 {object} types.TokensResponse "Token"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 401 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /refresh_tokens [post]
func (h *authHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req types.RefreshTokensRequest
	ip := r.RemoteAddr
	err := json.NewDecoder(r.Body).Decode(&req) // Validate guid
	if err != nil {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "invalid request"})
		return
	}
	if req.UserID == "" {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "missing guid"})
		return
	}
	if req.RefreshToken == "" {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "missing refresh_token"})
		return
	}

	authHeader := r.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")
	if len(t) != 2 {
		pkg.Unauthorized(w, r, pkg.ErrorResponse{Message: "unauthorized"})
		return
	}

	accessToken := t[1]

	tokens, err := h.service.RefreshTokens(req.UserID, ip, accessToken, req.RefreshToken)
	if err != nil {
		pkg.ProcessError(w, r, pkg.ErrorResponse{Message: err.Error(), Err: err})
		return
	}

	cookie := &http.Cookie{
		Name:    "access_token",
		Path:    "/",
		Value:   tokens.AccessToken.Token,
		Expires: tokens.AccessToken.ExpTime,
	}

	resp := types.TokensResponse{
		AccessToken:  tokens.AccessToken.Token,
		RefreshToken: tokens.RefreshToken.Token,
	}

	http.SetCookie(w, cookie)

	render.Status(r, http.StatusOK)

	render.JSON(w, r, resp)
}
