package campaign

import (
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/narrative"
)

func TestSavedAudioFreezesVisualsWithoutChangingNarrationFingerprint(t *testing.T) {
	project := syntheticProject()
	hash := Fingerprint(project)
	project.Events[0].Visual = &domain.StoryVisual{
		MediaID:           "illustration",
		StartSeconds:      2,
		DurationSeconds:   8,
		TransitionSeconds: 1,
	}
	project.Media = []domain.MediaAsset{{
		ID:        "illustration",
		Kind:      "image",
		Title:     "Synthetic original image",
		Generated: &domain.ImageGeneration{Prompt: "Synthetic original prompt"},
	}}
	if Fingerprint(project) != hash {
		t.Fatal("presentation-only changes invalidated saved narration continuity")
	}

	reading := narrative.Readings(project)[0]
	value, err := PlanAudio(project, []ReadyReading{{
		Reading:  reading,
		AudioID:  "audio",
		Duration: 60,
	}}, AudioInput{Seconds: 90})
	if err != nil {
		t.Fatal(err)
	}

	frozen := value.Parts[0].Passages[0].Visual
	if frozen == nil || frozen.Image.Title != "Synthetic original image" || frozen.Cue.StartSeconds != 2 {
		t.Fatal("image and timing were not frozen with the saved reading")
	}

	project.Media[0].Generated.Prompt = "Changed later"
	project.Events[0].Visual.StartSeconds = 20
	if frozen.Image.Generated.Prompt != "Synthetic original prompt" || frozen.Cue.StartSeconds != 2 {
		t.Fatal("an existing reel plan followed later presentation edits")
	}
}
