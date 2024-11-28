package repository

import "time"

type AuthRepository interface {
	GetToken(string) ([]byte, time.Time, error)
	AddToken(token []byte, userID string, expires time.Time)
	DeleteToken(string) error
}
