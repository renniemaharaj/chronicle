package evidence

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chronicle/internal/domain"
)

func TestCaptureHonorsSourceExclusions(t *testing.T) {
	requested := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if r.URL.Path == "/robots.txt" {
			w.Write([]byte("User-agent: Chronicle\nDisallow: /private"))
			return
		}
		requested = true
	}))
	defer server.Close()
	_, err := CaptureDocument(context.Background(), server.Client(), domain.Source{URL: server.URL + "/private"})
	if err == nil || !strings.Contains(err.Error(), "disallowed") || requested {
		t.Fatal("excluded source paths must not be retrieved")
	}
}
