package repositories

import (
	"database/sql"
)

type Repositories struct {
	db *sql.DB
	MessageLog *MessageLogRepository
}

func NewRepositoriesInstance(db *sql.DB) *Repositories {
	return &Repositories{
		db: db,
		MessageLog: NewMessageLogRepository(db),
	}
}