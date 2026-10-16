package backgrounds

import (
	"context"
	"log/slog"
	"time"

	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/repositories"
)

type PollingDB struct {
	repositories *repositories.Repositories
	queue        chan<- models.PendingMessage
	interval     time.Duration
}

func NewPollingDB(repos *repositories.Repositories, queue chan<- models.PendingMessage) *PollingDB {
	return &PollingDB{
		repositories: repos,
		queue:        queue,
		interval:     5 * time.Second,
	}
}

func (p *PollingDB) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *PollingDB) poll(ctx context.Context) {
	messages, err := p.repositories.MessageLog.FetchPending(ctx)
	if err != nil {
		slog.Error("failed to fetch pending messages", "error", err)
		return
	}

	for _, msg := range messages {
		select {
		case <-ctx.Done():
			return
		case p.queue <- msg:
		}
	}
}