package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/config"
	"github.com/diproducts/application-tracker-go/internal/repository/storage"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"

	_ "github.com/lib/pq"
)

const uniqueViolationErrorCode = pq.ErrorCode("23505")

type TokenStorage struct {
	db *sqlx.DB
}

func NewTokenStorage(dbCfg *config.Database) (*TokenStorage, error) {
	const op = "storage.postgresql.NewTokenStorage"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.DBName,
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TokenStorage{db: db}, nil
}

// BlacklistToken stores user_id, token_id and expiry time in the token_blacklist table
func (ts *TokenStorage) BlacklistToken(ctx context.Context, userID int64, tokenID string, expiry time.Time) error {
	const op = "storage.postgresql.BlacklistToken"

	stmt, err := ts.db.PrepareContext(
		ctx,
		"INSERT INTO token_blacklist(user_id, token_id, expiry) VALUES(?, ?, ?)",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, userID, tokenID, expiry)
	if err != nil {
		var pqError pq.Error

		if errors.As(err, &pqError) && pqError.Code == uniqueViolationErrorCode {
			return fmt.Errorf("%s: %w", op, storage.ErrTokenAlreadyBlacklisted)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// IsBlacklisted checks if an
func (ts *TokenStorage) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	const op = "storage.postgresql.BlacklistToken"

	stmt, err := ts.db.PrepareContext(ctx, "SELECT user_id FROM token_blacklist WHERE token_id=?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var userID string
	err = stmt.QueryRowContext(ctx, tokenID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
