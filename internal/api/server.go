package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"chronicle/internal/campaign"
	"chronicle/internal/config"
	"chronicle/internal/geography"
	"chronicle/internal/research"
	"chronicle/internal/speech"
	"chronicle/internal/storage"
)

type Server struct {
	campaigns *campaign.Manager
	audio     *speech.Cache
	config    config.Config
	store     storage.Store
	research  *research.Manager
	geography *geography.Searcher
}

func New(store storage.Store, manager *research.Manager, options ...config.Config) http.Handler {
	settings := config.Defaults()
	if len(options) > 0 {
		settings = options[0]
	}

	s := &Server{
		campaigns: campaign.New(store),
		audio:     speech.NewCache(settings.Workspace),
		store:     store,
		research:  manager,
		config:    settings,
		geography: geography.NewSearcher(settings.GeocoderURL),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/config", s.publicConfig)
	mux.HandleFunc("GET /api/places/search", s.searchPlaces)
	s.speechRoutes(mux)
	s.campaignRoutes(mux)
	s.storyImageRoutes(mux)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		write(w, 200, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/projects", s.projects)
	mux.HandleFunc("PUT /api/projects/{id}", s.saveProject)
	mux.HandleFunc("GET /api/projects/{id}", s.project)
	mux.HandleFunc("GET /api/workspace/events", s.workspaceEvents)
	mux.Handle("/api/mcp", s.assistantHandler())
	mux.HandleFunc("GET /api/settings", s.settings)
	mux.HandleFunc("PUT /api/settings", s.saveSettings)
	mux.HandleFunc("GET /api/tools", s.tools)
	mux.HandleFunc("POST /api/research", s.startResearch)
	mux.HandleFunc("GET /api/research", s.runs)
	mux.HandleFunc("GET /api/research/{id}", s.run)
	mux.HandleFunc("POST /api/research/{id}/{action}", s.control)
	mux.HandleFunc("GET /api/research/{id}/events", s.events)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet {
			origin, err := url.Parse(r.Header.Get("Origin"))
			if err != nil || !s.allowsOrigin(origin, r.Host) {
				fail(w, fmt.Errorf("request origin is not the local Chronicle workspace"), 403)

				return
			}
		}

		mux.ServeHTTP(w, r)
	})
}

func write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func fail(w http.ResponseWriter, err error, status int) {
	if errors.Is(err, storage.ErrNotFound) {
		status = 404
	}

	write(w, status, map[string]string{"error": err.Error()})
}

func decode(w http.ResponseWriter, r *http.Request, value any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	return nil
}

func validateEndpoint(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || strings.TrimSpace(parsed.Host) == "" || parsed.User != nil {
		return fmt.Errorf("endpoint must be an HTTP(S) URL without embedded credentials")
	}

	return nil
}
