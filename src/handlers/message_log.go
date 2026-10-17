package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/requests"
	"github.com/alysonryan/message-service/services"
	"github.com/alysonryan/message-service/validations"
)

type MessageLogHandler struct {
	Service *services.MessageLogService
}

func NewMessageLogHandler(s *services.MessageLogService) *MessageLogHandler {
	return &MessageLogHandler{
		Service: s,
	}
}

func (h *MessageLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := &models.Message{
		MessageTypeID: req.MessageTypeID,
		ServiceOrderID: req.ServiceOrderID,
		Phone: req.Phone,
		Payload: req.Payload,
	}

	id, err := h.Service.Create(r.Context(), m)
	if err != nil {
		if errors.Is(err, validations.ErrInvalidMessage) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("failed to create message", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"messageLogId": id})
}