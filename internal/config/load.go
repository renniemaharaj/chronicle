package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Load applies explicit deployment settings after defaults; environment wins.
func Load() (Config, error) {
	settings := Defaults()
	if path := os.Getenv("CHRONICLE_CONFIG"); path != "" {
		if err := settings.readFile(path); err != nil {
			return Config{}, err
		}
	}

	if err := settings.readEnvironment(); err != nil {
		return Config{}, err
	}

	if err := settings.Validate(); err != nil {
		return Config{}, err
	}

	return settings, nil
}

func (c *Config) readFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read Chronicle configuration: %w", err)
	}

	if trimmed := bytes.TrimSpace(data); len(trimmed) == 0 || trimmed[0] != '{' {
		return fmt.Errorf("Chronicle configuration must be a JSON object")
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf("decode Chronicle configuration: %w", err)
	}

	if err := decoder.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("Chronicle configuration must contain one JSON object")
	}

	var explicit struct {
		Research map[string]any `json:"research"`
	}
	if err := json.Unmarshal(data, &explicit); err != nil {
		return err
	}

	for key, value := range explicit.Research {
		c.researchOverrides[key] = value
	}

	return nil
}
