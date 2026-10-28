package speech

import (
	"context"
	"fmt"
	"net/url"
)

type ElevenVoice struct {
	ID          string            `json:"voice_id"`
	Name        string            `json:"name"`
	Category    string            `json:"category"`
	Labels      map[string]string `json:"labels"`
	Description string            `json:"description"`
}

type ElevenModel struct {
	ID            string `json:"model_id"`
	Name          string `json:"name"`
	CanSpeak      bool   `json:"can_do_text_to_speech"`
	MaxCharacters int    `json:"maximum_text_length_per_request"`
}

type Subscription struct {
	Tier   string `json:"tier"`
	Used   int64  `json:"character_count"`
	Limit  int64  `json:"character_limit"`
	Reset  int64  `json:"next_character_count_reset_unix"`
	Status string `json:"status"`
}

func (c *ElevenLabs) Voices(ctx context.Context) ([]ElevenVoice, error) {
	voices := []ElevenVoice{}
	token := ""

	for page := 0; page < 100; page++ {
		var output struct {
			Voices []ElevenVoice `json:"voices"`
			More   bool          `json:"has_more"`
			Next   string        `json:"next_page_token"`
		}
		path := "/v2/voices?page_size=100"
		if token != "" {
			path += "&next_page_token=" + url.QueryEscape(token)
		}

		if err := c.json(ctx, path, &output); err != nil {
			return nil, err
		}

		voices = append(voices, output.Voices...)
		if !output.More {
			return voices, nil
		}

		if output.Next == "" || output.Next == token {
			return nil, fmt.Errorf("ElevenLabs voice pagination could not continue")
		}

		token = output.Next
	}

	return nil, fmt.Errorf("ElevenLabs voice catalog exceeded the page limit")
}

func (c *ElevenLabs) Models(ctx context.Context) ([]ElevenModel, error) {
	var models []ElevenModel
	if err := c.json(ctx, "/v1/models", &models); err != nil {
		return nil, err
	}

	available := []ElevenModel{}

	for _, model := range models {
		if model.CanSpeak {
			available = append(available, model)
		}
	}

	return available, nil
}

func (c *ElevenLabs) Usage(ctx context.Context) (Subscription, error) {
	var usage Subscription
	err := c.json(ctx, "/v1/user/subscription", &usage)

	return usage, err
}
