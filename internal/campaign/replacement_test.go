package campaign

import (
	"testing"
	"time"
)

func replacementFixture() (Campaign, ActionInput) {
	views := int64(8)
	value := Campaign{
		Revision:  1,
		Platforms: []string{"facebook", "youtube"},
		Parts: []Part{{
			Number:   1,
			Status:   "published",
			File:     "original.mp4",
			Duration: 30,
			Publications: map[string]Receipt{
				"facebook": {URL: "https://www.facebook.com/reel/111"},
				"youtube":  {URL: "https://www.youtube.com/shorts/abcdefghijk"},
			},
			Metrics: []Metrics{{
				Platform: "facebook",
				Views:    &views,
			}},
		}},
	}
	input := ActionInput{
		Revision:    1,
		Part:        1,
		Action:      "replace",
		Platform:    "facebook",
		PreviousURL: "https://www.facebook.com/reel/111",
		URL:         "https://www.facebook.com/reel/222",
		File:        "illustrated.mp4",
		Duration:    30,
		Metrics: &Metrics{
			Views:      &views,
			ObservedAt: time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	return value, input
}

func TestReplacementPreservesReceiptMetricsAndOtherPlatform(t *testing.T) {
	value, input := replacementFixture()
	if err := Apply(&value, input); err != nil {
		t.Fatal(err)
	}

	part := value.Parts[0]
	if len(part.Replacements) != 1 || len(part.Replacements[0].Metrics) != 2 || len(part.Metrics) != 0 {
		t.Fatal("old observations must remain with the old receipt")
	}

	if part.Replacements[0].Before.File != "original.mp4" || part.File != "illustrated.mp4" {
		t.Fatal("both take files must remain identifiable")
	}

	if ReceiptFor(part, "youtube").File != "original.mp4" || ReceiptFor(part, "youtube").URL != "https://www.youtube.com/shorts/abcdefghijk" {
		t.Fatal("the other platform must retain its original receipt and take")
	}

	input.Revision = value.Revision
	input.Platform = "youtube"
	input.PreviousURL = "https://www.youtube.com/shorts/abcdefghijk"
	input.URL = "https://www.youtube.com/shorts/lmnopqrstuv"
	if err := Apply(&value, input); err != nil {
		t.Fatalf("both platforms may use the same illustrated file: %v", err)
	}
}

func TestReplacementRejectsProtectedUnknownStaleAndUnmatchedPosts(t *testing.T) {
	for _, scenario := range []string{"threshold", "previous-high", "unknown", "stale", "wrong-old-url", "same-url", "same-file", "duplicate-url"} {
		t.Run(scenario, func(t *testing.T) {
			value, input := replacementFixture()
			high := int64(10)

			switch scenario {
			case "threshold":
				input.Metrics.Views = &high
			case "previous-high":
				value.Parts[0].Metrics[0].Views = &high
			case "unknown":
				input.Metrics.Views = nil
			case "stale":
				input.Metrics.ObservedAt = time.Now().Add(-6 * time.Minute).Format(time.RFC3339Nano)
			case "wrong-old-url":
				input.PreviousURL = "https://www.facebook.com/reel/333"
			case "same-url":
				input.URL = input.PreviousURL
			case "same-file":
				input.File = value.Parts[0].File
			case "duplicate-url":
				value.Parts = append(value.Parts, Part{Publications: map[string]Receipt{"facebook": {URL: input.URL}}})
			}

			if Apply(&value, input) == nil {
				t.Fatal("unsafe replacement accepted")
			}

			if value.Revision != 1 || len(value.Parts[0].Replacements) != 0 {
				t.Fatal("a rejected replacement must not mutate history")
			}
		})
	}
}
