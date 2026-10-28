package speech

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestElevenLabsCatalogPaginationModelsAndUsage(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("xi-api-key") != "test-secret" {
			t.Error("catalog requests need the server API key")
		}

		switch r.URL.Path {
		case "/v2/voices":
			if r.URL.Query().Get("next_page_token") == "page-two" {
				fmt.Fprint(w, `{"voices":[{"voice_id":"two","name":"Second"}],"has_more":false}`)
			} else {
				fmt.Fprint(w, `{"voices":[{"voice_id":"one","name":"First"}],"has_more":true,"next_page_token":"page-two"}`)
			}
		case "/v1/models":
			fmt.Fprint(w, `[{"model_id":"speech","can_do_text_to_speech":true},{"model_id":"other","can_do_text_to_speech":false}]`)
		case "/v1/user/subscription":
			fmt.Fprint(w, `{"tier":"test","character_count":125,"character_limit":1000,"next_character_count_reset_unix":1900000000}`)
		}
	}))
	defer upstream.Close()

	client := NewElevenLabs(upstream.URL, "test-secret")
	voices, err := client.Voices(context.Background())
	if err != nil || len(voices) != 2 || voices[1].ID != "two" {
		t.Fatalf("voice discovery must follow pagination: %v", err)
	}

	models, err := client.Models(context.Background())
	if err != nil || len(models) != 1 || models[0].ID != "speech" {
		t.Fatal("only text-to-speech models should be offered")
	}

	usage, err := client.Usage(context.Background())
	if err != nil || usage.Used != 125 || usage.Limit != 1000 || usage.Reset != 1900000000 {
		t.Fatal("account usage must preserve the provider's reported values")
	}
}

func TestElevenLabsErrorsNeverExposeProviderBodiesOrKeys(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "test-secret private-account", 401)
	}))
	defer upstream.Close()

	client := NewElevenLabs(upstream.URL, "test-secret")
	_, err := client.Usage(context.Background())
	if err == nil || strings.Contains(err.Error(), "test-secret") || strings.Contains(err.Error(), "private-account") {
		t.Fatal("upstream error bodies must not leak authentication or account details")
	}

	called := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
	}))
	defer target.Close()

	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer redirect.Close()

	client = NewElevenLabs(redirect.URL, "test-secret")
	_, err = client.Usage(context.Background())
	if err == nil || called {
		t.Fatal("redirects must not forward the API key to another endpoint")
	}
}

func TestElevenLabsKeyIDErrorOffersActionableGuidance(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"detail":{"status":"api_key_id_used_as_api_key","message":"test-secret must not leak"}}`)
	}))
	defer upstream.Close()

	client := NewElevenLabs(upstream.URL, "test-secret")
	_, err := client.Usage(context.Background())
	if err == nil || !strings.Contains(err.Error(), "sk_") || strings.Contains(err.Error(), "test-secret") {
		t.Fatal("known authentication errors should explain how to fix the key without exposing upstream details")
	}
}
