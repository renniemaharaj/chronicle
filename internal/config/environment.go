package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (c *Config) readEnvironment() error {
	stringsByName := map[string]*string{
		"CHRONICLE_ADDR":                 &c.Address,
		"CHRONICLE_WORKSPACE":            &c.Workspace,
		"CHRONICLE_STATIC_DIR":           &c.StaticDir,
		"DATABASE_URL":                   &c.DatabaseURL,
		"CHRONICLE_GEOCODER_URL":         &c.GeocoderURL,
		"CHRONICLE_MAP_TILE_URL":         &c.MapTileURL,
		"CHRONICLE_MAP_ATTRIBUTION":      &c.MapAttribution,
		"ELEVENLABS_API_KEY":             &c.Speech.ElevenLabsKey,
		"CHRONICLE_ELEVENLABS_URL":       &c.Speech.ElevenLabsURL,
		"CHRONICLE_SPEECH_URL":           &c.Speech.Endpoint,
		"GOOGLE_APPLICATION_CREDENTIALS": &c.Speech.CredentialsFile,
		"GOOGLE_CLOUD_QUOTA_PROJECT":     &c.Speech.QuotaProject,
	}

	for name, target := range stringsByName {
		if value, exists := os.LookupEnv(name); exists {
			*target = value
		}
	}

	if value, exists := os.LookupEnv("CHRONICLE_ALLOWED_ORIGINS"); exists {
		c.AllowedOrigins = []string{}

		for _, origin := range strings.Split(value, ",") {
			if origin = strings.TrimSpace(origin); origin != "" {
				c.AllowedOrigins = append(c.AllowedOrigins, origin)
			}
		}
	}

	researchStrings := map[string]string{
		"CHRONICLE_LLM_PROVIDER":  "provider",
		"CHRONICLE_LLM_URL":       "baseUrl",
		"CHRONICLE_LLM_MODEL":     "model",
		"CHRONICLE_LLM_API_KEY":   "apiKey",
		"CHRONICLE_MCP_URL":       "mcpUrl",
		"CHRONICLE_SEARCH_TOOL":   "searchTool",
		"CHRONICLE_RETRIEVE_TOOL": "retrieveTool",
	}

	for name, key := range researchStrings {
		if value, exists := os.LookupEnv(name); exists {
			c.researchOverrides[key] = value
		}
	}

	researchNumbers := map[string]string{
		"CHRONICLE_MAX_REQUESTS": "maxRequests",
		"CHRONICLE_CONCURRENCY":  "concurrency",
	}

	for name, key := range researchNumbers {
		if value, exists := os.LookupEnv(name); exists {
			number, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("%s must be an integer", name)
			}

			c.researchOverrides[key] = number
		}
	}

	c.Research = c.ApplyResearch(c.Research)

	return nil
}
