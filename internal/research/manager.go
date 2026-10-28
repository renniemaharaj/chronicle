package research

import (
	"context"
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
	return &Manager{
		Store: store,
		jobs:  map[string]*job{},
		ModelFactory: func(settings domain.Settings) providers.Model {
			return providers.Client{
				Settings: settings,
				HTTP:     &http.Client{Timeout: 90 * time.Second},
			}
		},
	}
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
	j := &job{
		run:    run,
		cancel: cancel,
		store:  m.Store,
		ctx:    ctx,
		model:  m.ModelFactory(settings),
	}
	m.jobs[run.ID] = j
	go func() {
		j.execute(settings)
		m.mu.Lock()
		delete(m.jobs, run.ID)
		m.mu.Unlock()
	}()

	return run, nil
}
