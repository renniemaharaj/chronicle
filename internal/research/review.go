package research

import (
	"fmt"

	"chronicle/internal/domain"
	"chronicle/internal/evidence"
	"chronicle/internal/mcp"
)

func (j *job) review(client *mcp.Client, settings domain.Settings) error {
	run := j.snapshot()
	for _, question := range run.Plan {
		if question.Interval == "targeted review" {
			return nil
		}
	}
	if run.Config.MaxRequests-run.Requests < 5 {
		j.activity("Budget limits targeted follow-up. Coverage remains partial.", "", nil, "")
		return nil
	}
	question := domain.Question{Question: fmt.Sprintf("Find an independent primary or scholarly source that challenges or corroborates the documented origins and major turning points of %s. Investigate missing periods and competing dates.", run.Subject), Interval: "targeted review", Status: "pending"}
	index := len(run.Plan)
	j.update(func(run *domain.ResearchRun) {
		run.Phase = "Targeted evidence review"
		run.Plan = append(run.Plan, question)
	})
	batch, source, err := j.investigate(client, settings, question)
	if err != nil {
		return err
	}
	j.update(func(run *domain.ResearchRun) {
		project, mergeErr := mergeBatch(*run.Project, batch, source)
		if mergeErr != nil {
			err = mergeErr
			return
		}
		evidence.Assess(&project)
		run.Project = &project
		run.Plan[index].Status = "completed"
		run.Sources, run.Claims = len(project.Sources), len(project.Claims)
	})
	return err
}
