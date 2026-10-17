package handlers

import "github.com/alysonryan/message-service/services"

type Handlers struct {
	messageLog *MessageLogHandler
}

func NewHandlersInstance(s services.Services) *Handlers {
	return &Handlers{
		messageLog: NewMessageLogHandler(s.MessageLog),
	}
}