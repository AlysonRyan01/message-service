package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alysonryan/message-service/config"
	"github.com/alysonryan/message-service/models"
)

type WhatsAppService struct {
	cfg *config.Config
	client *http.Client
}

type SendTextRequest struct {
	Number           string `json:"number"`
	Text             string `json:"text"`
	Delay            int    `json:"delay"`
	LinkPreview      bool   `json:"linkPreview"`
	MentionsEveryOne bool   `json:"mentionsEveryOne"`
}

func NewWhatsAppService(cfg *config.Config) *WhatsAppService {
	return &WhatsAppService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *WhatsAppService) SendMessage(ctx context.Context, message *models.PendingMessage) error {
	url := fmt.Sprintf("%s/message/sendText/%s", s.cfg.WhatsAppURL, s.cfg.WhatsAppInstance)

	payload := SendTextRequest{
		Number: message.Phone,
		Text: message.Payload,
		Delay: 0,
		LinkPreview: true,
		MentionsEveryOne: false,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.cfg.WhatsAppAPIKEY)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send whatsapp message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("whatsapp api returned status %d: %s", resp.StatusCode, body)
	}

	return nil
}