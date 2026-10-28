package research

import (
	"context"
	"encoding/json"
	"fmt"

	"chronicle/internal/domain"
)

func (m *Manager) Control(ctx context.Context, id, action string, settings domain.Settings) (domain.ResearchRun, error) {
	data, err := m.Store.Get(ctx, "research", id)
	if err != nil {
		return domain.ResearchRun{}, err
	}

	var run domain.ResearchRun
	if err := json.Unmarshal(data, &run); err != nil {
		return run, err
	}

	if action == "resume" {
		if run.Status == "completed" {
			return run, fmt.Errorf("completed research cannot be resumed")
		}

		config := run.Config
		config.APIKey = settings.APIKey
		config.MaxRequests = max(config.MaxRequests, settings.MaxRequests)

		return m.launch(run, config)
	}

	if action != "pause" && action != "cancel" {
		return run, fmt.Errorf("unknown research action")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	j := m.jobs[id]
	if j == nil {
		return run, fmt.Errorf("session is not running")
	}

	status := "paused"
	if action == "cancel" {
		status = "cancelled"
	}

	j.update(func(run *domain.ResearchRun) {
		run.Status = status
		run.Phase = status
	})
	j.cancel()

	return j.snapshot(), nil
}
