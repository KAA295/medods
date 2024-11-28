package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

type RefreshToken struct {
	Token string
}

type AccessToken struct {
	Token   string
	ExpTime time.Time
}

type RefreshEntry struct {
	UserID  string
	Token   string
	Expires time.Time
}

type CustomClaims struct {
	UserID string
	Ip     string
	jwt.RegisteredClaims
}
