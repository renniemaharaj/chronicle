package speech

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ElevenLabs struct {
	HTTP     *http.Client
	Endpoint string
	APIKey   string
}

func NewElevenLabs(endpoint, key string) *ElevenLabs {
	return &ElevenLabs{
		HTTP: &http.Client{
			Timeout: 90 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Endpoint: strings.TrimRight(endpoint, "/"),
		APIKey:   key,
	}
}

func (c *ElevenLabs) call(ctx context.Context, method, path string, input any) ([]byte, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("configure an ElevenLabs API key in Voice & playback or ELEVENLABS_API_KEY")
	}

	var body bytes.Buffer
	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, &body)
	if err != nil {
		return nil, fmt.Errorf("invalid ElevenLabs request")
	}

	request.Header.Set("xi-api-key", c.APIKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("ElevenLabs request failed; check connectivity or try again")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, elevenLabsError(response)
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, (16<<20)+1))
	if err != nil || len(data) > 16<<20 || len(data) == 0 {
		return nil, fmt.Errorf("ElevenLabs returned empty, incomplete or oversized audio/data")
	}

	return data, nil
}

func (c *ElevenLabs) json(ctx context.Context, path string, output any) error {
	data, err := c.call(ctx, "GET", path, nil)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("ElevenLabs returned invalid account data")
	}

	return nil
}

func (c *ElevenLabs) Synthesize(ctx context.Context, input Request) ([]byte, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	body := map[string]any{
		"text":     input.Text,
		"model_id": input.ModelID,
		"voice_settings": map[string]any{
			"stability":         0.5,
			"similarity_boost":  0.75,
			"use_speaker_boost": true,
		},
	}
	// Multilingual v2 does not accept an enforced language code.
	if input.ModelID != "eleven_multilingual_v2" {
		body["language_code"] = strings.Split(input.LanguageCode, "-")[0]
	}

	return c.call(ctx, "POST", "/v1/text-to-speech/"+url.PathEscape(input.VoiceName)+"?output_format=mp3_44100_128", body)
}
