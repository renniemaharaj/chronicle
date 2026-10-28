package api

import (
	"net/http"
	"net/url"
	"slices"
)

func (s *Server) publicConfig(w http.ResponseWriter, _ *http.Request) {
	// Expose only presentation settings; deployment credentials stay on the server.
	write(w, http.StatusOK, map[string]string{
		"mapTileUrl":     s.config.MapTileURL,
		"mapAttribution": s.config.MapAttribution,
	})
}

func (s *Server) allowsOrigin(origin *url.URL, host string) bool {
	if origin.String() == "" {
		return true
	}

	if origin.Host == "" || (origin.Scheme != "http" && origin.Scheme != "https") || origin.User != nil || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
		return false
	}

	return origin.Host == host || slices.Contains(s.config.AllowedOrigins, origin.String())
}
