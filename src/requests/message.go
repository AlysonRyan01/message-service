package requests

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type CreateMessageRequest struct {
	MessageTypeID  int16  `json:"messageTypeId"`
	ServiceOrderID *int64 `json:"serviceOrderId,omitempty"`
	Phone          string `json:"phone"`
	Payload        string `json:"payload"`
}

var ErrInvalidRequest = errors.New("invalid request")

var phoneRegex = regexp.MustCompile(`^\d{12,13}$`)

func (r *CreateMessageRequest) Validate() error {
	r.Phone = strings.TrimSpace(r.Phone)
	r.Payload = strings.TrimSpace(r.Payload)

	if r.MessageTypeID <= 0 {
		return fmt.Errorf("%w: messageTypeId is required", ErrInvalidRequest)
	}

	if r.Phone == "" {
		return fmt.Errorf("%w: phone is required", ErrInvalidRequest)
	}
	if !phoneRegex.MatchString(r.Phone) {
		return fmt.Errorf("%w: phone must contain only digits with country code (e.g. 41999999999)", ErrInvalidRequest)
	}

	if r.Payload == "" {
		return fmt.Errorf("%w: payload is required", ErrInvalidRequest)
	}

	if r.ServiceOrderID != nil && *r.ServiceOrderID <= 0 {
		return fmt.Errorf("%w: serviceOrderId must be positive", ErrInvalidRequest)
	}

	return nil
}