package services

import (
	"github.com/alysonryan/message-service/config"
	"github.com/alysonryan/message-service/repositories"
)

type Services struct {
	Whatsapp *WhatsAppService
	MessageLog *MessageLogService
}

func NewServicesInstance(cfg *config.Config, r *repositories.Repositories) *Services {
	return &Services{
		Whatsapp: NewWhatsAppService(cfg),
		MessageLog: NewMessageLogService(r.MessageLog),
	}
}