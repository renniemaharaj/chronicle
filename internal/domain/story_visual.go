package domain

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

type ImageGeneration struct {
	Tool               string   `json:"tool"`
	Prompt             string   `json:"prompt"`
	CreatedAt          string   `json:"createdAt"`
	ReferenceSourceIDs []string `json:"referenceSourceIds"`
}

type StoryVisual struct {
	MediaID           string  `json:"mediaId"`
	StartSeconds      float64 `json:"startSeconds"`
	DurationSeconds   float64 `json:"durationSeconds"`
	TransitionSeconds float64 `json:"transitionSeconds"`
}

var localImageURL = regexp.MustCompile(`^/api/media/images/[a-f0-9]{64}\.(png|jpg|gif)$`)

func IsLocalImageURL(value string) bool {
	return localImageURL.MatchString(value)
}

func validateGeneratedImage(asset MediaAsset, entities, sources, claims map[string]bool) error {
	generation := asset.Generated
	_, err := time.Parse(time.RFC3339Nano, generation.CreatedAt)
	if asset.ID == "" || !entities[asset.EntityID] || asset.Kind != "image" || !IsLocalImageURL(asset.URL) || asset.Title == "" || strings.TrimSpace(asset.Caption) == "" || generation.Tool == "" || strings.TrimSpace(generation.Prompt) == "" || len(generation.Prompt) > 16000 || err != nil {
		return fmt.Errorf("generated image needs its local image, caption and generation provenance: %s", asset.ID)
	}

	if asset.SourceID != "" || asset.UseAsAvatar || asset.Retrieved || asset.Depicts == nil {
		return fmt.Errorf("generated images need a depiction period and cannot be source evidence or avatars: %s", asset.ID)
	}

	if err := asset.Depicts.Validate(); err != nil {
		return err
	}

	if asset.Depicts.Start == nil || asset.Depicts.Start.Precision == "unknown" {
		return fmt.Errorf("generated image needs an explicit depiction period")
	}

	for _, id := range generation.ReferenceSourceIDs {
		if !sources[id] {
			return fmt.Errorf("generated image references missing guidance source %s", id)
		}
	}

	for _, id := range asset.ClaimIDs {
		if !claims[id] {
			return fmt.Errorf("generated image references missing contextual claim %s", id)
		}
	}

	return nil
}

func validateStoryVisuals(project Project) error {
	assets := map[string]MediaAsset{}

	for _, asset := range project.Media {
		assets[asset.ID] = asset
	}

	for _, event := range project.Events {
		cue := event.Visual
		if cue == nil {
			continue
		}

		asset, exists := assets[cue.MediaID]
		if !exists || asset.Kind != "image" {
			return fmt.Errorf("passage %s needs an existing image", event.ID)
		}

		if asset.Depicts != nil {
			start, end, known := event.Time.Bounds()
			imageStart, imageEnd, imageKnown := asset.Depicts.Bounds()
			if known && imageKnown && (imageEnd < start || end < imageStart) {
				return fmt.Errorf("image depiction period does not overlap passage %s", event.ID)
			}
		}

		for _, value := range []float64{cue.StartSeconds, cue.DurationSeconds, cue.TransitionSeconds} {
			if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > 3600 {
				return fmt.Errorf("passage image timing must be finite seconds between 0 and 3600")
			}
		}

		if cue.DurationSeconds < 1 || cue.TransitionSeconds > 5 || cue.TransitionSeconds*2 > cue.DurationSeconds {
			return fmt.Errorf("passage image needs at least one second and two fades within its duration")
		}
	}

	return nil
}
