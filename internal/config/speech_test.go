package config

import "testing"

func TestSpeechEnvironmentCredentials(t *testing.T) {
	t.Setenv("CHRONICLE_CONFIG", "")
	t.Setenv("ELEVENLABS_API_KEY", "synthetic-secret")
	t.Setenv("CHRONICLE_ELEVENLABS_URL", "https://speech.example")

	settings, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if settings.Speech.ElevenLabsKey != "synthetic-secret" || settings.Speech.ElevenLabsURL != "https://speech.example" {
		t.Fatal("ElevenLabs environment settings must be detected")
	}
}
