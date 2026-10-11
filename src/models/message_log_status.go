package models


type MessageLogStatus struct {
	MessageLogStatusID int16  `json:"messageLogStatusId" db:"message_log_status_id"`
	Name                string `json:"name" db:"name"`
}

const (
	StatusPending    int16 = 1
	StatusProcessing int16 = 2
	StatusSent       int16 = 3
	StatusFailed     int16 = 4
)