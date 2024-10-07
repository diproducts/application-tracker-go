package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) SaveUser(ctx context.Context, user *models.User) (int64, error) {
	const op = "storage.postgresql.SaveUser"

	stmt, err := ur.db.PrepareContext(ctx, "INSERT INTO users(email, password, name) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, user.Email, user.HashedPassword, user.Name)
	if err != nil {
		var pqErr pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == uniqueViolationErrorCode {
			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (ur *UserRepository) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgresql.User"

	stmt, err := ur.db.PreparexContext(ctx, "SELECT email, password, name FROM users WHERE email=?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User

	err = stmt.Get(&user, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
