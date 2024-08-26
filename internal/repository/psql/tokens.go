package psql

import (
	"database/sql"

	"github.com/SavelyDev/crud-app/internal/domain"
)

type TokensRepo struct {
	db *sql.DB
}

func NewTokensRepo(db *sql.DB) *TokensRepo {
	return &TokensRepo{db: db}
}

func (r *TokensRepo) CreateSession(refreshToken domain.RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)",
		refreshToken.UserId, refreshToken.Token, refreshToken.ExpiresAt)

	return err
}

func (r *TokensRepo) GetSession(refreshToken string) (domain.RefreshSession, error) {
	var session domain.RefreshSession

	row := r.db.QueryRow("SELECT * FROM refresh_tokens WHERE token=$1", refreshToken)
	if err := row.Scan(&session.Id, &session.UserId, &session.Token, &session.ExpiresAt); err != nil {
		return session, err
	}

	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", session.UserId)

	return session, err
}
