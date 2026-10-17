package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/repositories"
	"github.com/alysonryan/message-service/validations"
)

type MessageLogService struct {
	Repository *repositories.MessageLogRepository
}

func NewMessageLogService(r *repositories.MessageLogRepository) *MessageLogService {
	return &MessageLogService{
		Repository: r,
	}
}

func (s *MessageLogService) Create(ctx context.Context, m *models.Message) (string, error) {
	if err := validate(m); err != nil {
		return "", err
	}

	return s.Repository.Create(ctx, m)
}

func (s *MessageLogService) Delete(ctx context.Context, messageLogID string) error {
	if messageLogID == "" {
		return validations.ErrInvalidID
	}

	return s.Repository.Delete(ctx, messageLogID)
}

func (s *MessageLogService) Requeue(ctx context.Context, messageLogID string) error {
	if messageLogID == "" {
		return validations.ErrInvalidID
	}

	return s.Repository.Requeue(ctx, messageLogID)
}

func validate(m *models.Message) error {
	if strings.TrimSpace(m.Phone) == "" {
		return fmt.Errorf("%w: phone is required", validations.ErrInvalidMessage)
	}
	if strings.TrimSpace(m.Payload) == "" {
		return fmt.Errorf("%w: payload is required", validations.ErrInvalidMessage)
	}
	if m.MessageTypeID == 0 {
		return fmt.Errorf("%w: message_type_id is required", validations.ErrInvalidMessage)
	}
	return nil
}