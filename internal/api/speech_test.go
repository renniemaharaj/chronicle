package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"chronicle/internal/config"
	"chronicle/internal/storage"
)

func speechTestServer(t *testing.T, key, endpoint string) (http.Handler, storage.Store) {
	t.Helper()

	settings := config.Defaults()
	settings.Workspace = t.TempDir()
	settings.Speech.ElevenLabsKey = key
	settings.Speech.ElevenLabsURL = endpoint
	store, err := storage.OpenFiles(settings.Workspace)
	if err != nil {
		t.Fatal(err)
	}

	return New(store, nil, settings), store
}

func speechCall(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(method, path, strings.NewReader(body)))

	return response
}

func TestSavedSpeechDoesNotSpendCreditsOnReplay(t *testing.T) {
	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Header.Get("xi-api-key") != "synthetic-secret" || r.URL.Query().Get("output_format") != "mp3_44100_128" {
			t.Error("synthesis must use server credentials and the chosen encoding")
		}

		var input map[string]any
		json.NewDecoder(r.Body).Decode(&input)
		if input["model_id"] != "eleven_multilingual_v2" {
			t.Error("synthesis must send the selected model")
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte("synthetic-mp3"))
	}))
	defer upstream.Close()

	handler, _ := speechTestServer(t, "synthetic-secret", upstream.URL)
	body := `{"provider":"elevenlabs","text":"Synthetic test passage.","voiceName":"test-voice","languageCode":"en-US"`
	response := speechCall(handler, "POST", "/api/speech/prepare", body+`,"generate":false}`)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"available":false`) || calls.Load() != 0 {
		t.Fatal("checking missing audio must never call a paid provider")
	}

	response = speechCall(handler, "POST", "/api/speech/prepare", body+`,"generate":true}`)
	var saved struct {
		URL string `json:"url"`
	}
	json.Unmarshal(response.Body.Bytes(), &saved)
	if response.Code != 200 || saved.URL == "" || calls.Load() != 1 {
		t.Fatalf("generation must return a saved audio URL: %s", response.Body.String())
	}

	for range 2 {
		replay := speechCall(handler, "GET", saved.URL, "")
		if replay.Code != 200 || replay.Body.String() != "synthetic-mp3" {
			t.Fatal("saved files must replay without synthesis")
		}
	}

	replay := httptest.NewRecorder()
	request := httptest.NewRequest("GET", saved.URL, nil)
	request.Header.Set("Range", "bytes=0-8")
	handler.ServeHTTP(replay, request)
	if replay.Code != 206 || replay.Body.String() != "synthetic" {
		t.Fatal("saved audio must support range requests for streaming and seeking")
	}

	response = speechCall(handler, "POST", "/api/speech/prepare", body+`,"generate":true}`)
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"cached":true`) || calls.Load() != 1 {
		t.Fatal("preparing an existing passage must reuse its audio")
	}
}
