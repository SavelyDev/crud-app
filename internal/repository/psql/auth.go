package psql

import (
	"database/sql"

	"github.com/SavelyDev/crud-app/internal/domain"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (s *AuthRepo) CreateUser(user domain.User) (int, error) {
	var id int

	row := s.db.QueryRow("INSERT INTO users (name, email, password_hash, registered) values ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.PasswordHash, user.Registered)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *AuthRepo) GetUserId(email, password string) (int, error) {
	var userId int

	row := s.db.QueryRow("SELECT id FROM users WHERE email=$1 AND password_hash=$2",
		email, password)
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}
