package research

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"chronicle/internal/domain"
	"chronicle/internal/providers"
	"chronicle/internal/storage"
)

type Manager struct {
	Store        storage.Store
	mu           sync.Mutex
	jobs         map[string]*job
	ModelFactory func(domain.Settings) providers.Model
}

type job struct {
	mu     sync.Mutex
	run    domain.ResearchRun
	cancel context.CancelFunc
	store  storage.Store
	ctx    context.Context
	model  providers.Model
}

func NewManager(store storage.Store) *Manager {
	return &Manager{Store: store, jobs: map[string]*job{}, ModelFactory: func(settings domain.Settings) providers.Model {
		return providers.Client{Settings: settings, HTTP: &http.Client{Timeout: 90 * time.Second}}
	}}
}

func (m *Manager) Start(ctx context.Context, subject string, settings domain.Settings) (domain.ResearchRun, error) {
	if subject == "" || len(subject) > 400 {
		return domain.ResearchRun{}, fmt.Errorf("subject must contain between 1 and 400 characters")
	}
	if settings.Model == "" || settings.MCPURL == "" || settings.SearchTool == "" || settings.RetrieveTool == "" {
		return domain.ResearchRun{}, fmt.Errorf("configure a model, MCP endpoint, search tool, and retrieval tool in Settings")
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	id := fmt.Sprintf("run-%d", time.Now().UnixNano())
	config := settings
	config.APIKey = ""
	run := domain.ResearchRun{ID: id, Subject: subject, Status: "pending", Phase: "Planning", CreatedAt: now, UpdatedAt: now, Model: settings.Model, Config: config, Plan: []domain.Question{}, Activities: []domain.Activity{}}
	if err := m.Store.Put(ctx, "research", id, run); err != nil {
		return run, err
	}
	started, err := m.launch(run, settings)
	if err != nil {
		run.Status, run.Phase, run.Error = "failed", "Not started", err.Error()
		m.Store.Put(ctx, "research", id, run)
	}
	return started, err
}

func (m *Manager) launch(run domain.ResearchRun, settings domain.Settings) (domain.ResearchRun, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.jobs) >= 2 {
		return run, fmt.Errorf("two research sessions are already running; pause one before starting another")
	}
	if existing := m.jobs[run.ID]; existing != nil {
		return run, fmt.Errorf("research session is already running")
	}
	run.Status = "running"
	ctx, cancel := context.WithCancel(context.Background())
	j := &job{run: run, cancel: cancel, store: m.Store, ctx: ctx, model: m.ModelFactory(settings)}
	m.jobs[run.ID] = j
	go func() {
		j.execute(settings)
		m.mu.Lock()
		delete(m.jobs, run.ID)
		m.mu.Unlock()
	}()
	return run, nil
}

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
	j.update(func(run *domain.ResearchRun) { run.Status = status; run.Phase = status })
	j.cancel()
	return j.snapshot(), nil
}
