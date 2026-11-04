package campaign

import (
	"context"
	"sync"
	"testing"

	"chronicle/internal/storage"
)

func TestPersistentCheckpointsAndConcurrentWriters(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	project := syntheticProject()
	if err := store.Put(ctx, "projects", project.ID, project); err != nil {
		t.Fatal(err)
	}

	manager := New(store)
	value, err := manager.Create(ctx, CreateInput{
		ProjectID: project.ID,
		Seconds:   15,
	})
	if err != nil {
		t.Fatal(err)
	}

	var wait sync.WaitGroup
	results := make(chan error, 2)

	for range 2 {
		wait.Go(func() {
			_, err := manager.Update(ctx, ActionInput{
				ID:       value.ID,
				Revision: 1,
				Part:     1,
				Action:   "begin",
			})
			results <- err
		})
	}

	wait.Wait()
	close(results)

	successes := 0

	for err := range results {
		if err == nil {
			successes++
		}
	}

	if successes != 1 {
		t.Fatal("only one concurrent take can claim the checkpoint")
	}

	reopened := New(store)
	saved, err := reopened.Get(ctx, value.ID)
	if err != nil || saved.Revision != 2 || saved.Parts[0].Status != "recording" {
		t.Fatalf("recording checkpoint did not survive a new manager: %v", err)
	}

	saved, err = reopened.Update(ctx, ActionInput{
		ID:       value.ID,
		Revision: saved.Revision,
		Part:     1,
		Action:   "interrupt",
		TakeID:   saved.Parts[0].TakeID,
	})
	if err != nil {
		t.Fatal(err)
	}

	project.Events[0].Narration.Text = "Changed synthetic narration"
	if err := store.Put(ctx, "projects", project.ID, project); err != nil {
		t.Fatal(err)
	}

	if _, err := reopened.Update(ctx, ActionInput{
		ID:       value.ID,
		Revision: saved.Revision,
		Part:     1,
		Action:   "begin",
	}); err == nil {
		t.Fatal("changed research must not silently shift a saved recording boundary")
	}
}
