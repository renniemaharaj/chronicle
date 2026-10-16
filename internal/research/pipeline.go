package research

import (
	"fmt"
	"sync"

	"chronicle/internal/conflicts"
	"chronicle/internal/domain"
	"chronicle/internal/evidence"
	"chronicle/internal/mcp"
	"chronicle/internal/narrative"
)

func (j *job) execute(settings domain.Settings) {
	var finalError error
	defer func() { j.finish(finalError) }()
	j.update(func(run *domain.ResearchRun) {
		run.Status = "running"
		run.Error = ""
		run.Config.MaxRequests = settings.MaxRequests
	})
	client, err := mcp.Connect(j.ctx, settings.MCPURL, settings.SearchTool, settings.RetrieveTool)
	if err != nil {
		finalError = fmt.Errorf("connect to MCP research server: %w", err)
		return
	}
	defer client.Close()
	if client.Tools[settings.SearchTool] == nil || client.Tools[settings.RetrieveTool] == nil {
		finalError = fmt.Errorf("configured research tool was not discovered")
		return
	}
	j.activity("Research tools discovered and schemas inspected", "mcp", nil, "")
	if err := j.plan(); err != nil {
		finalError = err
		return
	}
	j.update(func(run *domain.ResearchRun) { run.Phase = "Researching questions" })
	questions := j.snapshot().Plan
	semaphore := make(chan struct{}, max(1, min(4, settings.Concurrency)))
	var workers sync.WaitGroup
	var errorMu sync.Mutex
	for index, question := range questions {
		if question.Status == "completed" {
			continue
		}
		workers.Add(1)
		go func(index int, question domain.Question) {
			defer workers.Done()
			select {
			case semaphore <- struct{}{}:
			case <-j.ctx.Done():
				return
			}
			defer func() { <-semaphore }()
			batch, source, err := j.investigate(client, settings, question)
			if err == nil {
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
			}
			if err != nil {
				j.activity("Question incomplete: "+err.Error(), "", nil, "")
				errorMu.Lock()
				if finalError == nil {
					finalError = err
				}
				errorMu.Unlock()
			}
		}(index, question)
	}
	workers.Wait()
	if j.ctx.Err() != nil {
		finalError = j.ctx.Err()
		return
	}
	if finalError == nil {
		finalError = j.review(client, settings)
	}
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
		audit := domain.ResearchRecord{ID: run.ID, Subject: run.Subject, Status: status, CreatedAt: run.CreatedAt, Model: run.Model, Config: run.Config, Requests: run.Requests, Plan: run.Plan, Activities: run.Activities}
		run.Project.ResearchRecords = append(run.Project.ResearchRecords, audit)
	})
	run := j.snapshot()
	if err := domain.ValidateProject(*run.Project); err != nil {
		finalError = fmt.Errorf("research output failed validation: %w", err)
		return
	}
	if err := j.store.Put(j.ctx, "projects", run.Project.ID, run.Project); err != nil {
		finalError = err
	}
	j.activity("Research record saved; coverage and unresolved questions are retained", "", nil, "")
}
