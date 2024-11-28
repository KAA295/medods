package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/KAA295/medods/api/types"
	"github.com/KAA295/medods/pkg"
	"github.com/KAA295/medods/usecases/services"
)

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(authService services.AuthService)

func (h *authHandler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr // Real ip?
	userID := chi.URLParam(r, "guid")
	if userID == "" {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "missing guid"})
		return
	}

	tokens, err := h.service.GenerateTokens(userID, ip)
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

func (h *authHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	userID := chi.URLParam(r, "guid") // TODO: mb struct
	if userID == "" {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "missing guid"})
		return
	}
	refreshToken := chi.URLParam(r, "refresh_token")
	if refreshToken == "" {
		pkg.BadRequest(w, r, pkg.ErrorResponse{Message: "missing refresh_token"})
		return
	}
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		pkg.Unauthorized(w, r, pkg.ErrorResponse{Message: "no access_token"})
		return
	}

	accessToken := accessTokenCookie.Value

	tokens, err := h.service.RefreshToken(userID, ip, accessToken, refreshToken)
	if err != nil {
		pkg.ProcessError(w, r, pkg.ErrorResponse{Message: err.Error(), Err: err})
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
