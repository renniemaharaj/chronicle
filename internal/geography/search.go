package geography

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Searcher struct {
	mu    sync.Mutex
	rate  sync.Mutex
	last  time.Time
	cache map[string][]Result
}

func NewSearcher() *Searcher {
	return &Searcher{cache: make(map[string][]Result)}
}

func (s *Searcher) Search(ctx context.Context, query string) ([]Result, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 || len(query) > 200 {
		return nil, fmt.Errorf("enter a place name between 2 and 200 characters")
	}

	// Serialize explicit searches, including waits, across all browser sessions.
	s.rate.Lock()
	defer s.rate.Unlock()
	key := strings.ToLower(query)
	s.mu.Lock()
	cached, found := s.cache[key]
	s.mu.Unlock()
	if found {
		return cached, nil
	}

	if delay := time.Until(s.last.Add(time.Second)); delay > 0 {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}
	}

	endpoint := os.Getenv("CHRONICLE_GEOCODER_URL")
	if endpoint == "" {
		endpoint = "https://photon.komoot.io/api/"
	}
	u, err := url.Parse(endpoint)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, fmt.Errorf("the configured map search endpoint is invalid")
	}
	params := u.Query()
	params.Set("q", query)
	params.Set("lang", "en")
	params.Set("limit", "5")
	u.RawQuery = params.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Chronicle/0.1 (local historical research workspace)")
	request.Header.Set("Accept-Language", "en")
	s.last = time.Now()
	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("map search is unavailable; try again later")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("map search returned %d; try again later", response.StatusCode)
	}

	results, err := decodeResults(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("map search returned an unreadable response")
	}

	s.mu.Lock()
	if len(s.cache) >= 256 {
		clear(s.cache)
	}
	s.cache[key] = results
	s.mu.Unlock()
	return results, nil
}
