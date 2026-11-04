package campaign

import (
	"strings"
	"testing"

	"chronicle/internal/narrative"
)

func TestSavedAudioPlanPreservesTextAndBoundaries(t *testing.T) {
	project := syntheticProject()
	reading := narrative.Readings(project)[0]
	value, err := PlanAudio(project, []ReadyReading{{
		Reading:  reading,
		AudioID:  "saved-audio",
		Duration: 95,
	}}, AudioInput{
		Seconds:   90,
		Platforms: []string{"facebook", "youtube"},
	})
	if err != nil {
		t.Fatal(err)
	}

	text, end := "", 0.0

	for _, passage := range value.Parts[0].Passages {
		if passage.AudioID != "saved-audio" || passage.AudioStart != end {
			t.Fatal("caption cues must cover the saved audio in order")
		}

		text += passage.Text
		end = passage.AudioEnd
	}

	if text != reading.Text || end != 95 || !strings.Contains(text, "Sources disagree") {
		t.Fatal("saved narration, qualifications and measured duration must survive exactly")
	}

	if len(value.Parts) != 1 || value.Engine != "ffmpeg" {
		t.Fatal("a complete saved passage should remain in one episode")
	}
}
