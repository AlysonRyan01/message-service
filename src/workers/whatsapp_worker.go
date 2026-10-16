package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/repositories"
	"github.com/alysonryan/message-service/services"
)

type WhatsAppWorker struct {
	queue <- chan models.PendingMessage
	whatsAppService *services.WhatsAppService
	repository *repositories.MessageLogRepository
}

func NewWhatsAppWorker(queue <- chan models.PendingMessage, s *services.WhatsAppService, r *repositories.MessageLogRepository) *WhatsAppWorker {
	return &WhatsAppWorker{
		queue: queue,
		whatsAppService: s,
		repository: r,
	}
}
 
func (w *WhatsAppWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-w.queue:
			if !ok {
				return
			}
			w.process(ctx, msg)
		}
	}
}

const maxAttempts = 8

func (w *WhatsAppWorker) process(ctx context.Context, msg models.PendingMessage) {
	sendCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	sendErr := w.whatsAppService.SendMessage(sendCtx, &msg)
	cancel()

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()

	if sendErr != nil {
		slog.Error("failed to send whatsapp message",
			"message_log_id", msg.MessageLogID,
			"error", sendErr,
		)

		if msg.Attempts+1 >= maxAttempts {
			if err := w.repository.MarkFailed(ctx, msg.MessageLogID, sendErr.Error()); err != nil {
				slog.Error("failed to mark message as failed",
					"message_log_id", msg.MessageLogID,
					"error", err,
				)
			}
			return
		}

		if err := w.repository.IncrementAttempts(dbCtx, msg.MessageLogID, sendErr.Error()); err != nil {
			slog.Error("failed to register attempt",
				"message_log_id", msg.MessageLogID,
				"error", err,
			)
		}
		return
	}

	if err := w.repository.MarkSent(dbCtx, msg.MessageLogID); err != nil {
		slog.Error("failed to mark message as sent",
			"message_log_id", msg.MessageLogID,
			"error", err,
		)
	}
}