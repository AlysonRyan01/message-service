package services

import (
	"github.com/alysonryan/message-service/config"
)

type Services struct {
	Whatsapp *WhatsAppService
}

func NewServicesInstance(cfg *config.Config) *Services {
	return &Services{
		Whatsapp: NewWhatsAppService(cfg),
	}
}