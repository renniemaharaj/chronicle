package campaign

import (
	"encoding/json"

	"chronicle/internal/domain"
)

func freezeVisual(project domain.Project, event domain.Event) *FrozenVisual {
	if event.Visual == nil {
		return nil
	}

	for _, asset := range project.Media {
		if asset.ID != event.Visual.MediaID || asset.Kind != "image" || (asset.Generated == nil && !asset.Retrieved) {
			continue
		}

		data, err := json.Marshal(FrozenVisual{
			Cue:   *event.Visual,
			Image: asset,
		})
		if err != nil {
			return nil
		}

		var frozen FrozenVisual
		if json.Unmarshal(data, &frozen) == nil {
			return &frozen
		}
	}

	return nil
}
