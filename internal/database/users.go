package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`

	return m.DB.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Name,
	).Scan(
		&user.ID,
	)
}

func (m *UserModel) GetByID(userID int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, email, password FROM users WHERE id = $1`

	return m.getUser(query, ctx, userID)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT u.id, u.email, u.name, u.password FROM users u WHERE u.email = $1`

	return m.getUser(query, ctx, email)
}

func (m *UserModel) getUser(query string, ctx context.Context, args ...interface{}) (*User, error) {
	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
