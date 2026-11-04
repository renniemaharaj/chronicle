package campaign

import (
	"crypto/rand"
	"fmt"
	"math"
	"strings"
	"time"
)

type Replacement struct {
	Platform   string    `json:"platform"`
	Before     Receipt   `json:"before"`
	After      Receipt   `json:"after"`
	Metrics    []Metrics `json:"metrics"`
	RecordedAt string    `json:"recordedAt"`
}

const replacementProtectionViews = int64(10)

// Replacement records observed external actions; it never deletes or uploads a platform post.
func replacePublication(value *Campaign, part *Part, input ActionInput, now string) error {
	platform := Platform(input.Platform)
	before := ReceiptFor(*part, platform)
	previous, _, err := publicationURL(input.PreviousURL, platform)
	if err != nil || before == nil || before.URL != previous {
		return fmt.Errorf("replacement requires the exact current publication receipt")
	}

	address, id, err := publicationURL(input.URL, platform)
	if err != nil || address == previous {
		return fmt.Errorf("replacement requires a different confirmed public URL")
	}

	previousFile := before.File
	if previousFile == "" {
		previousFile = part.File
	}

	if strings.TrimSpace(input.File) == "" || input.File == previousFile || input.Duration <= 0 || math.IsNaN(input.Duration) || math.IsInf(input.Duration, 0) {
		return fmt.Errorf("replacement requires a new file and measured duration")
	}

	if input.Metrics == nil || input.Metrics.Views == nil || *input.Metrics.Views >= replacementProtectionViews {
		return fmt.Errorf("replacement requires an observed view count below 10")
	}

	if err := validateMetrics(*input.Metrics); err != nil {
		return err
	}

	observed, _ := time.Parse(time.RFC3339Nano, input.Metrics.ObservedAt)
	if time.Since(observed) > 5*time.Minute {
		return fmt.Errorf("refresh performance within five minutes before replacement")
	}

	for _, candidate := range value.Parts {
		if receipt := ReceiptFor(candidate, platform); receipt != nil {
			_, existing, _ := publicationURL(receipt.URL, platform)
			if existing == id {
				return fmt.Errorf("replacement URL already belongs to an episode")
			}
		}
	}

	archived := []Metrics{}
	retained := []Metrics{}

	for _, metric := range part.Metrics {
		if Platform(metric.Platform) == platform {
			if metric.Views != nil && *metric.Views >= replacementProtectionViews {
				return fmt.Errorf("a publication observed at 10 views or more is protected")
			}

			archived = append(archived, metric)
		} else {
			retained = append(retained, metric)
		}
	}

	finalMetric := *input.Metrics
	finalMetric.Platform, finalMetric.Method = platform, "manual"
	archived = append(archived, finalMetric)
	old := *before
	if old.File == "" {
		old.File, old.Duration = part.File, part.Duration
	}

	if part.Publications == nil {
		part.Publications = map[string]Receipt{}
	}

	for target, receipt := range part.Publications {
		if receipt.File == "" {
			receipt.File, receipt.Duration = part.File, part.Duration
			part.Publications[target] = receipt
		}
	}

	after := Receipt{
		URL:        address,
		RecordedAt: now,
		Method:     "manual",
		File:       input.File,
		Duration:   input.Duration,
	}
	part.Replacements = append(part.Replacements, Replacement{
		Platform:   platform,
		Before:     old,
		After:      after,
		Metrics:    archived,
		RecordedAt: now,
	})
	part.Publications[platform] = after
	if platform == "facebook" {
		part.Publication = &after
	}

	part.Metrics = retained
	part.TakeID = "take-" + rand.Text()
	part.File, part.Duration = input.File, input.Duration
	part.Takes = append(part.Takes, Take{
		ID:          part.TakeID,
		File:        input.File,
		Duration:    input.Duration,
		CompletedAt: now,
	})

	return nil
}
