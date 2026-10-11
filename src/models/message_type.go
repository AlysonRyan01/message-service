package models

type MessageType struct {
	MessageTypeID int16  `json:"messageTypeId" db:"message_type_id"`
	Name          string `json:"name" db:"name"`
}

const (
	TypeWelcome               int16 = 1
	TypeEstimate              int16 = 2
	TypeEstimateLink          int16 = 3
	TypeEstimateApproved      int16 = 4
	TypeEstimateRejected      int16 = 5
	TypeDeviceReady           int16 = 6
	TypeRejectedReady         int16 = 7
	TypeGuarantee             int16 = 8
	TypeProductSale           int16 = 9
	TypeOwnerEstimateDecision int16 = 10
	TypeNoRepair              int16 = 11
)