package main

import (
	"context"
	"encoding/json"
	"os"

	"chronicle/internal/domain"
	"chronicle/internal/narrative"
	"chronicle/internal/storage"
)

func openWorkspace(ctx context.Context) (storage.Store, error) {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		postgres, err := storage.OpenPostgres(ctx, url)
		if err != nil {
			return nil, err
		}
		if err := postgres.ApplyMigrations(ctx); err != nil {
			postgres.Close()
			return nil, err
		}
		return postgres, nil
	}
	path := os.Getenv("CHRONICLE_WORKSPACE")
	if path == "" {
		path = ".chronicle"
	}
	return storage.OpenFiles(path)
}

func prepareWorkspace(ctx context.Context, store storage.Store) error {
	runs, err := store.List(ctx, "research")
	if err != nil {
		return err
	}
	for _, data := range runs {
		var run domain.ResearchRun
		if err := json.Unmarshal(data, &run); err != nil {
			return err
		}
		if run.Status == "running" || run.Status == "pending" {
			run.Status, run.Phase = "paused", "Interrupted; ready to resume"
			if err := store.Put(ctx, "research", run.ID, run); err != nil {
				return err
			}
		}
	}
	projects, err := store.List(ctx, "projects")
	if err != nil {
		return err
	}
	for _, data := range projects {
		var project domain.Project
		if err := json.Unmarshal(data, &project); err != nil {
			return err
		}
		if project.Fixture {
			continue
		}
		narrative.EnsureStories(&project)
		if err := store.Put(ctx, "projects", project.ID, project); err != nil {
			return err
		}
	}
	return nil
}
