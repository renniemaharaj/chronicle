package api

import (
	"context"
	"os/exec"
	"testing"

	"chronicle/internal/campaign"
	"chronicle/internal/config"
	"chronicle/internal/domain"
	"chronicle/internal/narrative"
	"chronicle/internal/speech"
	"chronicle/internal/storage"
)

func TestAudioCampaignStopsAtFirstUnsavedReading(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg is unavailable")
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe is unavailable")
	}

	ctx := context.Background()
	root := t.TempDir()
	store, err := storage.OpenFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	settings := config.Defaults()
	settings.Workspace = root
	server := &Server{
		store:  store,
		config: settings,
		audio:  speech.NewCache(root),
	}
	project := domain.Project{
		ID:    "synthetic-audio",
		Title: "Synthetic audio test, not historical evidence",
		Events: []domain.Event{
			{
				ID:    "one",
				Title: "Synthetic first reading",
			}, {
				ID:    "two",
				Title: "Synthetic second reading",
			}, {
				ID:    "three",
				Title: "Synthetic missing reading",
			},
		},
	}
	if err := store.Put(ctx, "projects", project.ID, project); err != nil {
		t.Fatal(err)
	}

	data, err := exec.Command("ffmpeg", "-v", "error", "-f", "lavfi", "-i", "sine=frequency=440:duration=1", "-f", "mp3", "pipe:1").Output()
	if err != nil {
		t.Fatal(err)
	}

	readings := narrative.Readings(project)

	for _, reading := range readings[:2] {
		request := speech.Request{
			Provider:     "elevenlabs",
			VoiceName:    "synthetic-voice",
			ModelID:      "eleven_multilingual_v2",
			LanguageCode: "en-US",
			Text:         reading.Text,
		}
		if _, err := server.audio.Ensure(ctx, server.audioKey(request), func() ([]byte, error) {
			return data, nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	value, err := server.createAudioCampaign(ctx, campaign.AudioInput{
		ProjectID:    project.ID,
		Seconds:      90,
		Platforms:    []string{"facebook", "youtube"},
		Provider:     "elevenlabs",
		VoiceName:    "synthetic-voice",
		ModelID:      "eleven_multilingual_v2",
		LanguageCode: "en-US",
	})
	if err != nil {
		t.Fatal(err)
	}

	if value.ReadyThrough != readings[1].Event.ID || value.RemainingEvents != 1 || len(value.Parts) != 1 {
		t.Fatal("audio planning must stop at the first missing reading without skipping ahead")
	}

	for _, passage := range value.Parts[0].Passages {
		if passage.AudioID == "" || passage.AudioEnd <= passage.AudioStart {
			t.Fatal("every exported passage needs a measured saved-audio range")
		}
	}
}
