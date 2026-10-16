package research

import (
	"context"

	"chronicle/internal/domain"
)

func (m *Manager) Shutdown(ctx context.Context) {
	m.mu.Lock()
	jobs := []*job{}
	for _, j := range m.jobs {
		jobs = append(jobs, j)
	}
	m.mu.Unlock()
	for _, j := range jobs {
		j.update(func(run *domain.ResearchRun) {
			if run.Status == "completed" || run.Status == "failed" {
				return
			}
			run.Status = "paused"
			run.Phase = "Application stopped; ready to resume"
		})
		j.cancel()
	}
}
