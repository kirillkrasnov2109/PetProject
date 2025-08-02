package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"petProject/internal/domain"
	"time"
)

type HTTPSender struct {
	client *http.Client
}

func NewHTTPSender() *HTTPSender {
	return &HTTPSender{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *HTTPSender) Send(task domain.Task) error {
	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", task.Webhook, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
