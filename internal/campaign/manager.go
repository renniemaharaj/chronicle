package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"chronicle/internal/domain"
	"chronicle/internal/storage"
)

type Manager struct {
	store storage.Store
	mu    sync.Mutex
}

func New(store storage.Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) Get(ctx context.Context, id string) (Campaign, error) {
	var result Campaign
	data, err := m.store.Get(ctx, "campaigns", id)
	if err == nil {
		err = json.Unmarshal(data, &result)
	}

	return result, err
}

func (m *Manager) List(ctx context.Context) ([]Campaign, error) {
	data, err := m.store.List(ctx, "campaigns")
	result := []Campaign{}

	for _, raw := range data {
		var value Campaign
		if decodeErr := json.Unmarshal(raw, &value); decodeErr != nil {
			return nil, decodeErr
		}

		result = append(result, value)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})

	return result, err
}

func (m *Manager) Create(ctx context.Context, input CreateInput) (Campaign, error) {
	data, err := m.store.Get(ctx, "projects", input.ProjectID)
	var project domain.Project
	if err == nil {
		err = json.Unmarshal(data, &project)
	}

	if err != nil {
		return Campaign{}, err
	}

	value, err := Plan(project, input.Seconds)
	if err == nil {
		err = ValidatePlatforms(input.Platforms)
		value.Platforms = input.Platforms
	}

	if err == nil {
		err = m.store.Put(ctx, "campaigns", value.ID, value)
	}

	return value, err
}

func (m *Manager) Update(ctx context.Context, input ActionInput) (Campaign, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, err := m.Get(ctx, input.ID)
	if err != nil {
		return value, err
	}

	if input.Action == "begin" {
		all, err := m.List(ctx)
		if err != nil {
			return value, err
		}

		for _, series := range all {
			for _, part := range series.Parts {
				if part.Status == "recording" {
					return value, fmt.Errorf("a recording take is already open; recover or interrupt it first")
				}
			}
		}

		data, err := m.store.Get(ctx, "projects", value.ProjectID)
		var project domain.Project
		if err != nil {
			return value, err
		}

		if err := json.Unmarshal(data, &project); err != nil {
			return value, err
		}

		if value.ProjectHash != Fingerprint(project) {
			return value, fmt.Errorf("research changed since planning; create a new series after reviewing the changes")
		}
	}

	if err := Apply(&value, input); err != nil {
		return value, err
	}

	if err := m.store.Put(ctx, "campaigns", value.ID, value); err != nil {
		return Campaign{}, err
	}

	return value, nil
}
