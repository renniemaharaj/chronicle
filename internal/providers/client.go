package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"chronicle/internal/domain"
)

type Model interface {
	Complete(context.Context, string, string) (json.RawMessage, error)
}

type Client struct {
	Settings domain.Settings
	HTTP     *http.Client
}

func (c Client) Complete(ctx context.Context, system, input string) (json.RawMessage, error) {
	messages := []map[string]string{{"role": "system", "content": system}, {"role": "user", "content": input}}
	endpoint := strings.TrimRight(c.Settings.BaseURL, "/")
	body := map[string]any{"model": c.Settings.Model, "messages": messages, "stream": false}
	if c.Settings.Provider == "ollama" {
		endpoint += "/api/chat"
		body["format"] = "json"
		body["options"] = map[string]any{"temperature": 0}
	} else {
		endpoint += "/chat/completions"
		body["temperature"] = 0
		body["response_format"] = map[string]string{"type": "json_object"}
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Settings.APIKey != "" {
		request.Header.Set("Authorization", "Bearer "+c.Settings.APIKey)
	}
	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("model endpoint returned HTTP %d", response.StatusCode)
	}
	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	content := result.Message.Content
	if len(result.Choices) > 0 {
		content = result.Choices[0].Message.Content
	}
	if !json.Valid([]byte(content)) {
		return nil, fmt.Errorf("model did not return a valid JSON object")
	}
	return json.RawMessage(content), nil
}
