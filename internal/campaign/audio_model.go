package campaign

import "chronicle/internal/narrative"

type AudioInput struct {
	ProjectID    string   `json:"projectId"`
	Seconds      int      `json:"seconds"`
	Platforms    []string `json:"platforms"`
	Provider     string   `json:"provider"`
	VoiceName    string   `json:"voiceName"`
	LanguageCode string   `json:"languageCode"`
	ModelID      string   `json:"modelId"`
}

type ReadyReading struct {
	narrative.Reading
	AudioID  string
	Duration float64
}
