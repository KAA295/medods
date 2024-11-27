package services

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/KAA295/medods/domain"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) generateAccessToken(userID string, ip string) (string, error) {
	claims := domain.CustomClaims{
		UserID: userID,
		Ip:     ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte("TODO ENV"))
}

func (s *AuthService) generateRefreshToken() (string, error) {
	data := make([]byte, 64)
	_, err := rand.Read(data)
	return base64.URLEncoding.EncodeToString(data), err
}

func (s *AuthService) GenerateTokens() {
}

func (s *AuthService) RefreshToken() {
}
