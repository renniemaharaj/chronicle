package domain

import (
	"math"
	"strings"
	"testing"

	"chronicle/internal/temporal"
)

func syntheticIllustration() MediaAsset {
	return MediaAsset{
		ID:       "illustration",
		EntityID: "subject",
		Kind:     "image",
		Title:    "Synthetic illustration",
		Caption:  "Synthetic test reconstruction, not historical evidence.",
		URL:      "/api/media/images/" + strings.Repeat("a", 64) + ".png",
		Generated: &ImageGeneration{
			Tool:      "test",
			Prompt:    "Synthetic test prompt",
			CreatedAt: "2026-01-01T00:00:00Z",
		},
		Depicts: &temporal.Interval{Start: &temporal.TemporalValue{
			Year:      100,
			Precision: "year",
			Calendar:  "julian",
		}},
	}
}

func TestGeneratedIllustrationCannotMasqueradeAsSourceEvidence(t *testing.T) {
	asset := syntheticIllustration()
	check := func(value MediaAsset) error {
		return validateMedia(Project{Media: []MediaAsset{value}}, map[string]bool{"subject": true}, map[string]bool{}, map[string]bool{})
	}
	if err := check(asset); err != nil {
		t.Fatal(err)
	}

	for _, change := range []func(*MediaAsset){
		func(a *MediaAsset) {
			a.Retrieved = true
		},
		func(a *MediaAsset) {
			a.SourceID = "fabricated-source"
		},
		func(a *MediaAsset) {
			a.UseAsAvatar = true
		},
		func(a *MediaAsset) {
			a.URL = "https://example.org/invented.png"
		},
		func(a *MediaAsset) {
			a.Depicts = nil
		},
		func(a *MediaAsset) {
			a.LicenseURL = "javascript:invalid"
		},
	} {
		value := asset
		change(&value)
		if check(value) == nil {
			t.Fatal("unsupported generated-image provenance was accepted")
		}
	}
}

func TestStoryVisualRejectsMissingImagesBadTimingAndWrongPeriod(t *testing.T) {
	asset := syntheticIllustration()
	cue := StoryVisual{
		MediaID:           asset.ID,
		StartSeconds:      4,
		DurationSeconds:   10,
		TransitionSeconds: 1,
	}
	event := Event{
		ID:     "event",
		Visual: &cue,
		Time:   *asset.Depicts,
	}
	p := Project{
		Media:  []MediaAsset{asset},
		Events: []Event{event},
	}
	if err := validateStoryVisuals(p); err != nil {
		t.Fatal(err)
	}

	for _, change := range []func(*StoryVisual){
		func(v *StoryVisual) {
			v.MediaID = "missing"
		},
		func(v *StoryVisual) {
			v.StartSeconds = -1
		},
		func(v *StoryVisual) {
			v.DurationSeconds = math.Inf(1)
		},
		func(v *StoryVisual) {
			v.TransitionSeconds = 6
		},
		func(v *StoryVisual) {
			v.DurationSeconds = 1
		},
	} {
		bad := cue
		change(&bad)

		p.Events[0].Visual = &bad
		if validateStoryVisuals(p) == nil {
			t.Fatal("invalid image cue was accepted")
		}
	}

	p.Events[0].Visual = &cue
	p.Events[0].Time.Start = &temporal.TemporalValue{
		Year:      900,
		Precision: "year",
		Calendar:  "julian",
	}
	if validateStoryVisuals(p) == nil {
		t.Fatal("a definitively different image period was accepted")
	}
}
