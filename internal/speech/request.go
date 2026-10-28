package speech

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var voiceID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,120}$`)

func (r Request) Normalized() Request {
	if r.Provider == "" {
		r.Provider = "google"
	}

	if r.Provider == "elevenlabs" && r.ModelID == "" {
		r.ModelID = "eleven_multilingual_v2"
	}

	if r.Provider == "google" {
		r.ModelID = ""
	}

	return r
}

func (r Request) Validate() error {
	if r.Provider != "google" && r.Provider != "elevenlabs" {
		return fmt.Errorf("choose Google or ElevenLabs for saved audio")
	}

	if strings.TrimSpace(r.Text) == "" || len(r.Text) > 5000 || !utf8.ValidString(r.Text) || r.VoiceName == "" || r.LanguageCode == "" {
		return fmt.Errorf("speech needs 1-5000 text bytes, a voice and a language")
	}

	if r.Provider == "elevenlabs" && (!voiceID.MatchString(r.VoiceName) || !voiceID.MatchString(r.ModelID)) {
		return fmt.Errorf("invalid ElevenLabs voice or model identifier")
	}

	return nil
}
