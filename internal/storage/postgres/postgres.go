package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rizzmatch/rm_auth/internal/core/models"
	"github.com/rizzmatch/rm_auth/internal/storage"
)

type Storage struct {
	db *pgx.Conn
}

func New(login, password, dbName, host string, port int, sslmode string) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", login, password, host, port, dbName, sslmode)

	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close(context.Background())
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int, error) {
	const op = "storage.postgres.SaveUser"

	stmt := `INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id;`
	// args := pgx.NamedArgs{
	// 	"userEmail":    email,
	// 	"userPassHash": passHash,
	// }

	var id int
	if err := s.db.QueryRow(context.Background(), stmt, email, passHash).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrUserExists
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	stmt := `SELECT id, email, pass_hash FROM users WHERE email = $1;`

	var user models.User
	if err := s.db.QueryRow(context.Background(), stmt, email).Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, email string) (bool, error) {
	return true, nil
}

func (s *Storage) AddPhone(ctx context.Context, id int64, phone string) error {
	const op = "storage.postgres.AddPhone"

	stmt := `INSERT INTO phones (user_id, phone) VALUES ($1, $2);`

	if _, err := s.db.Exec(context.Background(), stmt, id, phone); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) EditPhone(ctx context.Context, id int64, phone string) error {
	const op = "storage.postgres.EditPhone"

	stmt := `UPDATE phones SET phone = $1 WHERE user_id = $2;`

	if _, err := s.db.Exec(context.Background(), stmt, phone, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddEmail(ctx context.Context, id int64, email string) error {
	const op = "storage.postgres.AddEmail"

	stmt := `INSERT INTO emails (user_id, email) VALUES ($1, $2);`

	if _, err := s.db.Exec(context.Background(), stmt, id, email); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) EditEmail(ctx context.Context, id int64, email string) error {
	const op = "storage.postgres.EditEmail"

	stmt := `UPDATE emails SET email = $1 WHERE user_id = $2;`

	if _, err := s.db.Exec(context.Background(), stmt, email, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
