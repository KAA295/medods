package postgres

import (
	"database/sql"
	"errors"

	"github.com/KAA295/medods/domain"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (repo *AuthRepository) GetToken(userID string) (domain.RefreshEntry, error) {
	query := "SELECT token, expiration_time FROM tokens WHERE user_id = $1"
	var refreshToken domain.RefreshEntry
	err := repo.DB.QueryRow(query, userID).Scan(&refreshToken.Token, &refreshToken.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.RefreshEntry{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.RefreshEntry{}, err
	}

	return refreshToken, nil
}

func (repo *AuthRepository) AddToken(refreshToken domain.RefreshEntry) error {
	query := "INSERT INTO tokens (token, user_id, expiration_time) VALUES ($1, $2, $3)"
	_, err := repo.DB.Exec(query, refreshToken.Token, refreshToken.UserID, refreshToken.Expires)
	if err != nil {
		return err
	}
	return nil
}

func (repo *AuthRepository) DeleteToken(userID string) error {
	query := "DELETE FROM tokens WHERE user_id = $1"
	_, err := repo.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
