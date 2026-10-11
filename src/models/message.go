package models

import "time"

type Message struct {
	MessageID      string    `json:"messageId" db:"message_id"`
	MessageLogID   string    `json:"messageLogId" db:"message_log_id"`
	MessageTypeID  int16     `json:"messageTypeId" db:"message_type_id"`
	ServiceOrderID *int64    `json:"serviceOrderId,omitempty" db:"service_order_id"`
	Phone          string    `json:"phone" db:"phone"`
	Payload        string    `json:"payload" db:"payload"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}