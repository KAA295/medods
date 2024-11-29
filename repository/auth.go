package repository

import (
	"github.com/KAA295/medods/domain"
)

type AuthRepository interface {
	GetToken(userID string) (domain.RefreshEntry, error)
	AddToken(refreshToken domain.RefreshEntry) error
	DeleteToken(userID string) error
}
