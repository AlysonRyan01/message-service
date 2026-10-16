package services

import (
	"github.com/alysonryan/message-service/config"
	"github.com/alysonryan/message-service/repositories"
)

type Services struct {
	cfg *config.Config
	Repositories *repositories.Repositories
	Whatsapp *WhatsAppService
}

func NewServicesInstance(cfg *config.Config, r *repositories.Repositories) *Services {
	return &Services{
		cfg: cfg,
		Repositories: r,
		Whatsapp: NewWhatsAppService(cfg),
	}
}