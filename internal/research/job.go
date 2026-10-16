package research

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"chronicle/internal/domain"
)

func (j *job) snapshot() domain.ResearchRun {
	j.mu.Lock()
	defer j.mu.Unlock()
	data, _ := json.Marshal(j.run)
	var snapshot domain.ResearchRun
	json.Unmarshal(data, &snapshot)
	return snapshot
}

func (j *job) update(change func(*domain.ResearchRun)) {
	j.mu.Lock()
	defer j.mu.Unlock()
	change(&j.run)
	j.run.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := j.store.Put(ctx, "research", j.run.ID, j.run); err != nil {
		j.run.Error = "Checkpoint could not be persisted: " + err.Error()
		j.cancel()
	}
}

func (j *job) reserve() error {
	if err := j.ctx.Err(); err != nil {
		return err
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.run.Requests >= j.run.Config.MaxRequests {
		return fmt.Errorf("research request budget exhausted; raise the budget in Settings and resume")
	}
	j.run.Requests++
	checkpoint, cancel := context.WithTimeout(j.ctx, 5*time.Second)
	defer cancel()
	return j.store.Put(checkpoint, "research", j.run.ID, j.run)
}

func (j *job) activity(message, tool string, input any, output string) {
	encoded, _ := json.Marshal(input)
	if len(output) > 100000 {
		output = output[:100000] + "\n[output truncated in activity log]"
	}
	j.update(func(run *domain.ResearchRun) {
		run.Activities = append(run.Activities, domain.Activity{Time: time.Now().UTC().Format(time.RFC3339Nano), Message: message, Tool: tool, Input: encoded, Output: output})
	})
}

func (j *job) Complete(ctx context.Context, system, input string) (json.RawMessage, error) {
	var lastError error
	for attempt := 0; attempt < 2; attempt++ {
		if err := j.reserve(); err != nil {
			return nil, err
		}
		raw, err := j.model.Complete(ctx, system, input)
		if err == nil {
			j.activity("Structured model output received", "model", map[string]string{"system": system, "input": input}, string(raw))
			return raw, nil
		}
		lastError = err
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if attempt == 0 {
			j.activity("Model request failed; one bounded retry will be attempted", "model", nil, "")
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(250 * time.Millisecond):
			}
		}
	}
	return nil, lastError
}

func (j *job) finish(err error) {
	j.update(func(run *domain.ResearchRun) {
		if run.Status == "paused" || run.Status == "cancelled" {
			return
		}
		if err != nil {
			run.Status, run.Phase, run.Error = "failed", "Research incomplete", err.Error()
			return
		}
		run.Status, run.Phase, run.Error = "completed", "Research record assembled", ""
	})
}
