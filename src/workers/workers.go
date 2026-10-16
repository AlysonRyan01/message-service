package workers

import (
	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/repositories"
	"github.com/alysonryan/message-service/services"
)

type Workers struct {
	WhatsApp *WhatsAppWorker
}

func NewWorkersInstance(queue <-chan models.PendingMessage, services *services.Services, repositories *repositories.Repositories) *Workers {
	return &Workers{
		WhatsApp: NewWhatsAppWorker(queue, services.Whatsapp, repositories.MessageLog),
	}
}