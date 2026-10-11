package models

type PendingMessage struct {
	MessageLogID   string
	MessageTypeID  int16
	ServiceOrderID *int64
	Phone          string
	Payload        string
	Attempts       int
}