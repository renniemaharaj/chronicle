package speech

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGoogleVoiceAndAudioRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Goog-User-Project") != "local-quota-project" {
			t.Error("missing quota project")
		}
		if r.URL.Path == "/voices" {
			if r.URL.Query().Get("languageCode") != "en-US" {
				t.Error("language missing")
			}
			w.Write([]byte(`{"voices":[{"name":"en-US-Test","languageCodes":["en-US"],"ssmlGender":"FEMALE"}]}`))
			return
		}
		if r.Method != "POST" || r.URL.Path != "/text:synthesize" {
			t.Errorf("unexpected speech request %s", r.URL.Path)
		}
		var body struct {
			Input struct {
				Text string `json:"text"`
			} `json:"input"`
			AudioConfig struct {
				Encoding string `json:"audioEncoding"`
			} `json:"audioConfig"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.Input.Text != "Evidence before narrative." || body.AudioConfig.Encoding != "MP3" {
			t.Error("wrong synthesis payload")
		}
		w.Write([]byte(`{"audioContent":"YXVkaW8="}`))
	}))
	defer server.Close()
	client := &Client{HTTP: server.Client(), Endpoint: server.URL, QuotaProject: "local-quota-project"}
	voices, err := client.Voices(context.Background(), "en-US")
	if err != nil || len(voices) != 1 {
		t.Fatalf("voices: %v %v", voices, err)
	}
	audio, err := client.Synthesize(context.Background(), Request{Text: "Evidence before narrative.", VoiceName: voices[0].Name, LanguageCode: "en-US"})
	if err != nil || string(audio) != "audio" {
		t.Fatalf("audio: %s %v", audio, err)
	}
	_, err = client.Synthesize(context.Background(), Request{Text: strings.Repeat("a", 5001), VoiceName: "test", LanguageCode: "en-US"})
	if err == nil {
		t.Fatal("oversized text accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Voices(ctx, "en-US"); err == nil {
		t.Fatal("cancelled request succeeded")
	}
}
