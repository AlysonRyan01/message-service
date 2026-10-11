package models

import "time"

type MessageLog struct {
	MessageLogID       string     `json:"messageLogId" db:"message_log_id"`
	MessageLogStatusID int16      `json:"messageLogStatusId" db:"message_log_status_id"`
	Attempts           int        `json:"attempts" db:"attempts"`
	LastError          *string    `json:"lastError,omitempty" db:"last_error"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	SentAt             *time.Time `json:"sentAt,omitempty" db:"sent_at"`
}
