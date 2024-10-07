package postgresql

import (
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const uniqueViolationErrorCode = pq.ErrorCode("23505")

func InitDB(dbCfg *config.Database) (*sqlx.DB, error) {
	const op = "storage.postgresql.DB"

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

	return db, nil
}
