package research

import (
	"chronicle/internal/conflicts"
	"chronicle/internal/domain"
	"chronicle/internal/narrative"
)

func (j *job) recordResearchResult(finalError error) {
	j.update(func(run *domain.ResearchRun) {
		run.Phase = "Validating chronology and evidence"
		conflicts.Analyze(run.Project)
		narrative.Build(run.Project)

		run.Conflicts = len(run.Project.Conflicts)
	})
	j.update(func(run *domain.ResearchRun) {
		status := "completed"
		if finalError != nil {
			status = "incomplete"
		}

		audit := domain.ResearchRecord{
			ID:         run.ID,
			Subject:    run.Subject,
			Status:     status,
			CreatedAt:  run.CreatedAt,
			Model:      run.Model,
			Config:     run.Config,
			Requests:   run.Requests,
			Plan:       run.Plan,
			Activities: run.Activities,
		}
		run.Project.ResearchRecords = append(run.Project.ResearchRecords, audit)
	})
}
