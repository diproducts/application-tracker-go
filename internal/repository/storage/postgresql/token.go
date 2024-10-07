package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/repository/storage"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"

	_ "github.com/lib/pq"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// BlacklistToken stores user_id, token_id and expiry time in the token_blacklist table
func (ts *TokenRepository) BlacklistToken(ctx context.Context, userID int64, tokenID string, expiry time.Time) error {
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

// IsBlacklisted checks if the token is blacklisted.
func (ts *TokenRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
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
