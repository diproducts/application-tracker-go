package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/diproducts/application-tracker-go/internal/repository/storage"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) SaveUser(ctx context.Context, user *models.User) (int64, error) {
	const op = "storage.postgresql.SaveUser"

	stmt, err := ur.db.PreparexContext(ctx, "INSERT INTO users(email, password, name) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, user.Email, user.HashedPassword, user.Name)
	if err != nil {
		var pqErr pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == uniqueViolationErrorCode {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserAlreadyExists)
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

	stmt, err := ur.db.PreparexContext(ctx, "SELECT * FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user dbUser
	err = stmt.Get(&user, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.User{
		ID:             user.ID,
		HashedPassword: user.Password,
		Email:          user.Email,
		Name:           user.Name,
	}, nil
}

type dbUser struct {
	ID       int64  `db:"id"`
	Password string `db:"password"`
	Email    string `db:"email"`
	Name     string `db:"name"`
}
