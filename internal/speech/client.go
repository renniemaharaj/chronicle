package speech

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	HTTP         *http.Client
	Endpoint     string
	QuotaProject string
}

type Voice struct {
	Name          string   `json:"name"`
	LanguageCodes []string `json:"languageCodes"`
	Gender        string   `json:"ssmlGender"`
}

type Request struct {
	Provider     string `json:"provider,omitempty"`
	ModelID      string `json:"modelId,omitempty"`
	Text         string `json:"text"`
	VoiceName    string `json:"voiceName"`
	LanguageCode string `json:"languageCode"`
}

func (c *Client) call(ctx context.Context, method, path string, input, output any) error {
	var body bytes.Buffer
	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, &body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	if c.QuotaProject != "" {
		request.Header.Set("X-Goog-User-Project", c.QuotaProject)
	}

	response, err := c.HTTP.Do(request)
	if err != nil {
		return fmt.Errorf("Google speech request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// Never echo upstream bodies which may contain authentication details.
		return fmt.Errorf("Google speech returned HTTP %d; check API access, credentials, quota and voice", response.StatusCode)
	}

	return json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(output)
}

func (c *Client) Voices(ctx context.Context, language string) ([]Voice, error) {
	var output struct {
		Voices []Voice `json:"voices"`
	}
	err := c.call(ctx, "GET", "/voices?languageCode="+url.QueryEscape(language), nil, &output)

	return output.Voices, err
}

func (c *Client) Synthesize(ctx context.Context, input Request) ([]byte, error) {
	if len(input.Text) == 0 || len(input.Text) > 5000 || input.VoiceName == "" || input.LanguageCode == "" {
		return nil, fmt.Errorf("speech needs 1–5000 text bytes, a voice and a language")
	}

	body := map[string]any{
		"input": map[string]string{"text": input.Text},
		"voice": map[string]string{
			"name":         input.VoiceName,
			"languageCode": input.LanguageCode,
		},
		"audioConfig": map[string]string{"audioEncoding": "MP3"},
	}
	var output struct {
		Audio string `json:"audioContent"`
	}
	if err := c.call(ctx, "POST", "/text:synthesize", body, &output); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(output.Audio)
}
