package db

import (
	"database/sql"
	"fmt"

	"github.com/alysonryan/message-service/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Init(cfg *config.Config) (*sql.DB, error) {
	connectionString := cfg.ConnectionString()
	
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}