package repositories

import (
	"database/sql"
)

type Repositories struct {
	MessageLog *MessageLogRepository
}

func NewRepositoriesInstance(db *sql.DB) *Repositories {
	return &Repositories{
		MessageLog: NewMessageLogRepository(db),
	}
}