package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/config"
	"github.com/pressly/goose/v3"
	"log"

	_ "github.com/lib/pq"
)

func Up(cfg *config.Database, migrationsPath string) {
	db, err := initDB(cfg)
	if err != nil {
		log.Println("[ERROR] failed to connect to database:", err.Error())
		return
	}
	defer db.Close()

	err = goose.Up(db, migrationsPath)
	if err != nil {
		if errors.Is(err, goose.ErrAlreadyApplied) {
			log.Println("no migrations to apply")
			return
		}

		log.Println("[ERROR] failed to apply migrations:", err.Error())
		return
	}

	log.Println("migrations applied successfully")
}

func initDB(dbCfg *config.Database) (*sql.DB, error) {
	const op = "migrator.initDB"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.DBName,
		dbCfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
