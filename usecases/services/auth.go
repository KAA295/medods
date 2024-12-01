package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/KAA295/medods/domain"
	"github.com/KAA295/medods/repository"
	"github.com/KAA295/medods/usecases"
)

type AuthService struct {
	authRepository repository.AuthRepository
	emailService   usecases.EmailService
}

func NewAuthService(authRepository repository.AuthRepository, emailService usecases.EmailService) *AuthService {
	return &AuthService{authRepository: authRepository, emailService: emailService}
}

func (s *AuthService) generateAccessToken(userID string, ip string) (domain.AccessToken, error) {
	expTime := time.Now().Add(time.Second * 30)
	claims := domain.CustomClaims{
		UserID: userID,
		Ip:     ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return domain.AccessToken{}, err
	}

	return domain.AccessToken{
		Token:   signedToken,
		ExpTime: expTime,
	}, nil
}

func (s *AuthService) generateRefreshToken() (domain.RefreshToken, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	return domain.RefreshToken{Token: base64.URLEncoding.EncodeToString(data)}, err
}

func (s *AuthService) GenerateTokens(userID string, ip string) (domain.Tokens, error) {
	_, err := s.authRepository.GetToken(userID)
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	accessToken, err := s.generateAccessToken(userID, ip)
	if err != nil {
		fmt.Println("1", err)
		return domain.Tokens{}, err
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		fmt.Println("2", err)
		return domain.Tokens{}, err
	}

	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken.Token), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("3", err)
		return domain.Tokens{}, err
	}
	expTime := time.Now().Add(time.Hour * 24)

	err = s.authRepository.AddToken(domain.RefreshEntry{
		UserID:  userID,
		Token:   string(encryptedToken),
		Expires: expTime,
	})
	if err != nil {
		fmt.Println("4", err)
		return domain.Tokens{}, err
	}

	return domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(userID string, ip string, accessToken string, refreshToken string) (domain.Tokens, error) {
	token, err := s.authRepository.GetToken(userID)
	if err != nil {
		return domain.Tokens{}, err
	}

	if time.Now().After(token.Expires) {
		err := s.authRepository.DeleteToken(userID)
		if err != nil {
			return domain.Tokens{}, err
		}
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	claims := &domain.CustomClaims{}

	_, err = jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return domain.Tokens{}, domain.ErrBadRequest
	}

	if time.Now().Before(claims.ExpiresAt.Time) {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	if userID != claims.UserID {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(token.Token), []byte(refreshToken))
	if err != nil {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	if claims.Ip != ip {
		s.emailService.Send("Warning, ip changed")
	}

	err = s.authRepository.DeleteToken(userID)
	if err != nil {
		return domain.Tokens{}, err
	}
	return s.GenerateTokens(userID, ip)
}
