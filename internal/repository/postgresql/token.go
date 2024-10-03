package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/config"

	_ "github.com/lib/pq"
)

type TokenStorage struct {
	db *sql.DB
}

func NewTokenStorage(dbCfg *config.Database) (*TokenStorage, error) {
	const op = "repository.postgresql.NewTokenStorage"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TokenStorage{db: db}, nil
}
