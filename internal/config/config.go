package config

import "chronicle/internal/domain"

type Speech struct {
	Endpoint        string `json:"endpoint"`
	CredentialsFile string `json:"credentialsFile"`
	QuotaProject    string `json:"quotaProject"`
	ElevenLabsKey   string `json:"elevenLabsKey"`
	ElevenLabsURL   string `json:"elevenLabsUrl"`
}

type Config struct {
	Address        string          `json:"address"`
	Workspace      string          `json:"workspace"`
	StaticDir      string          `json:"staticDir"`
	DatabaseURL    string          `json:"databaseUrl"`
	AllowedOrigins []string        `json:"allowedOrigins"`
	GeocoderURL    string          `json:"geocoderUrl"`
	MapTileURL     string          `json:"mapTileUrl"`
	MapAttribution string          `json:"mapAttribution"`
	Speech         Speech          `json:"speech"`
	Research       domain.Settings `json:"research"`

	researchOverrides map[string]any
}

func Defaults() Config {
	return Config{
		Address:        "127.0.0.1:8080",
		Workspace:      ".chronicle",
		StaticDir:      "web/build",
		GeocoderURL:    "https://photon.komoot.io/api/",
		MapTileURL:     "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
		MapAttribution: `&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors`,
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5174",
		},
		Speech: Speech{
			Endpoint:      "https://texttospeech.googleapis.com/v1",
			ElevenLabsURL: "https://api.elevenlabs.io",
		},
		Research: domain.Settings{
			Provider:    "ollama",
			BaseURL:     "http://localhost:11434",
			MaxRequests: 40,
			Concurrency: 2,
		},
		researchOverrides: map[string]any{},
	}
}
