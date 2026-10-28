package research

import (
	"fmt"
	"strings"
	"time"

	"chronicle/internal/domain"
)

func capturedSource(batch extraction, document, address, hostname string, contentHash [32]byte, resources []domain.Resource, truncated bool) domain.Source {
	title := batch.Title
	if title == "" || !strings.Contains(normalize(document), normalize(title)) {
		title = address
	}

	if !strings.Contains(normalize(document), normalize(batch.Author)) {
		batch.Author = ""
	}

	if !strings.Contains(normalize(document), normalize(batch.Publisher)) {
		batch.Publisher = ""
	}

	classification := batch.Classification
	if classification != "primary" && classification != "secondary" {
		classification = "tertiary"
	}

	source := domain.Source{
		ID:                stableID("source", address),
		Title:             title,
		Author:            batch.Author,
		Publisher:         batch.Publisher,
		URL:               address,
		RetrievedAt:       time.Now().UTC().Format(time.RFC3339Nano),
		Classification:    classification,
		IndependenceGroup: hostname,
		Verification:      "retrieved",
	}
	source.Resources = resources
	source.ContentHash = fmt.Sprintf("%x", contentHash)
	source.ID = stableID("source", address+"#"+source.ContentHash)
	if truncated {
		source.Limitations = []string{"Document context was truncated for bounded model processing. Remaining passages were not reviewed."}
	}

	return source
}

func normalize(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
