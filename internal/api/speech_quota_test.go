package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestSpeechQuotaPreservesCompletedPartsAndCanResume(t *testing.T) {
	var calls atomic.Int32
	var exhausted atomic.Bool
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		if exhausted.Load() {
			w.WriteHeader(401)
			fmt.Fprint(w, `{"detail":{"status":"quota_exceeded","message":"private-credit-data"}}`)

			return
		}

		fmt.Fprint(w, "synthetic-audio")
	}))
	defer upstream.Close()

	handler, _ := speechTestServer(t, "synthetic-secret", upstream.URL)
	body := `{"provider":"elevenlabs","voiceName":"test-voice","languageCode":"en-US","generate":true,"text":`
	first := body + `"First synthetic part."}`
	second := body + `"Second synthetic part."}`
	if response := speechCall(handler, "POST", "/api/speech/prepare", first); response.Code != 200 {
		t.Fatal(response.Body.String())
	}

	exhausted.Store(true)

	response := speechCall(handler, "POST", "/api/speech/prepare", second)
	if response.Code != 402 || !strings.Contains(response.Body.String(), `"code":"insufficient_credits"`) || strings.Contains(response.Body.String(), "private-credit-data") {
		t.Fatalf("quota must be actionable and redact provider data: %s", response.Body.String())
	}

	// Previously generated audio stays available even while the provider rejects new work.
	response = speechCall(handler, "POST", "/api/speech/prepare", first)
	if response.Code != 200 || calls.Load() != 2 {
		t.Fatal("cached audio must bypass the exhausted provider")
	}

	exhausted.Store(false)

	response = speechCall(handler, "POST", "/api/speech/prepare", second)
	if response.Code != 200 || calls.Load() != 3 {
		t.Fatal("renewal must allow generating only the missing part")
	}
}
