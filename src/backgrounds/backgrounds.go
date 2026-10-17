package backgrounds

import (
	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/repositories"
)

type Backgrounds struct {
	PollingDB *PollingDB
}

func NewBackgroundsInstance(repos *repositories.Repositories, queue chan<- models.PendingMessage) *Backgrounds {
	return &Backgrounds{
		PollingDB: NewPollingDB(repos, queue),
	}
}