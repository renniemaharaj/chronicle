package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"chronicle/internal/campaign"
	"chronicle/internal/domain"
	"chronicle/internal/narrative"
	"chronicle/internal/speech"
)

func (s *Server) createAudioCampaign(ctx context.Context, input campaign.AudioInput) (campaign.Campaign, error) {
	data, err := s.store.Get(ctx, "projects", input.ProjectID)
	var project domain.Project
	if err == nil {
		err = json.Unmarshal(data, &project)
	}

	if err != nil {
		return campaign.Campaign{}, err
	}

	request := speech.Request{
		Provider:     input.Provider,
		VoiceName:    input.VoiceName,
		LanguageCode: input.LanguageCode,
		ModelID:      input.ModelID,
		Text:         "inspect saved narration",
	}.Normalized()
	if err := request.Validate(); err != nil {
		return campaign.Campaign{}, err
	}

	ready := []campaign.ReadyReading{}

	for _, reading := range narrative.Readings(project) {
		request.Text = reading.Text
		id := s.audioKey(request)
		exists, err := s.audio.Has(id)
		if err != nil {
			return campaign.Campaign{}, err
		}

		if !exists {
			break
		}

		seconds, err := s.audioDuration(ctx, id)
		if err != nil {
			return campaign.Campaign{}, err
		}

		ready = append(ready, campaign.ReadyReading{
			Reading:  reading,
			AudioID:  id,
			Duration: seconds,
		})
	}

	value, err := campaign.PlanAudio(project, ready, input)
	if err == nil {
		err = s.store.Put(ctx, "campaigns", value.ID, value)
	}

	return value, err
}

func (s *Server) audioDuration(ctx context.Context, id string) (float64, error) {
	file, err := s.audio.Open(id)
	if err != nil {
		return 0, err
	}

	path := file.Name()
	file.Close()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=nw=1:nk=1", path).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe could not measure saved narration: %w", err)
	}

	return strconv.ParseFloat(strings.TrimSpace(string(result)), 64)
}

func (s *Server) audioCampaign(w http.ResponseWriter, r *http.Request) {
	var input campaign.AudioInput
	if err := decode(w, r, &input); err != nil {
		fail(w, err, 400)

		return
	}

	value, err := s.createAudioCampaign(r.Context(), input)
	if err != nil {
		fail(w, err, 422)

		return
	}

	write(w, 201, value)
}
