package config

import (
	"encoding/json"

	"chronicle/internal/domain"
)

// Explicit file and environment values also override saved browser settings.
func (c Config) ApplyResearch(settings domain.Settings) domain.Settings {
	encoded, _ := json.Marshal(c.researchOverrides)
	_ = json.Unmarshal(encoded, &settings)

	return settings
}

func (c Config) ManagesAPIKey() bool {
	_, configured := c.researchOverrides["apiKey"]

	return configured
}
