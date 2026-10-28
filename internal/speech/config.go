package speech

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"chronicle/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func New(ctx context.Context, options ...config.Speech) (*Client, error) {
	settings := config.Defaults().Speech
	if len(options) > 0 {
		settings = options[0]
	}

	credentials, err := findCredentials(ctx, settings.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("Google speech needs Application Default Credentials; see docs/speech.md")
	}

	var metadata struct {
		QuotaProjectID string `json:"quota_project_id"`
	}
	_ = json.Unmarshal(credentials.JSON, &metadata)
	if settings.QuotaProject == "" {
		settings.QuotaProject = metadata.QuotaProjectID
	}

	return &Client{
		HTTP: &http.Client{
			Transport: &oauth2.Transport{Source: credentials.TokenSource},
			Timeout:   30 * time.Second,
		},
		Endpoint:     strings.TrimRight(settings.Endpoint, "/"),
		QuotaProject: settings.QuotaProject,
	}, nil
}

func findCredentials(ctx context.Context, path string) (*google.Credentials, error) {
	const scope = "https://www.googleapis.com/auth/cloud-platform"
	if path == "" {
		return google.FindDefaultCredentials(ctx, scope)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return google.CredentialsFromJSON(ctx, data, scope)
}
