package speech

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCachePersistsVariantsAndDeduplicatesConcurrentGeneration(t *testing.T) {
	root := t.TempDir()
	cache := NewCache(root)
	input := Request{
		Provider:     "elevenlabs",
		Text:         "Synthetic narration.",
		VoiceName:    "test-a",
		LanguageCode: "en-US",
	}.Normalized()
	id := cache.Key(input, "https://speech.example")
	var calls atomic.Int32
	var group sync.WaitGroup

	for range 8 {
		group.Add(1)
		go func() {
			defer group.Done()

			_, err := cache.Ensure(context.Background(), id, func() ([]byte, error) {
				calls.Add(1)

				return []byte("synthetic-mp3"), nil
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}

	group.Wait()
	if calls.Load() != 1 {
		t.Fatal("simultaneous requests must synthesize only once")
	}

	reopened := NewCache(root)
	file, err := reopened.Open(id)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	if string(data) != "synthetic-mp3" {
		t.Fatal("saved audio must survive cache/server recreation")
	}

	for _, field := range []string{"text", "voice", "model", "language", "provider"} {
		changed := input

		switch field {
		case "text":
			changed.Text += " Changed."
		case "voice":
			changed.VoiceName = "test-b"
		case "model":
			changed.ModelID = "another-model"
		case "language":
			changed.LanguageCode = "es"
		case "provider":
			changed.Provider = "google"
		}

		if reopened.Key(changed, "https://speech.example") == id {
			t.Fatalf("%s changes must have independent saved audio", field)
		}
	}

	if cache.Key(input, "https://other.example") == id {
		t.Fatal("deployment endpoint changes must not reuse another provider's audio")
	}
}

func TestFailedOrCancelledGenerationDoesNotPoisonCache(t *testing.T) {
	cache := NewCache(t.TempDir())
	id := cache.Key(Request{Text: "test"}, "test")
	_, err := cache.Ensure(context.Background(), id, func() ([]byte, error) {
		return nil, fmt.Errorf("provider unavailable")
	})
	if err == nil {
		t.Fatal("provider failures must be reported")
	}

	exists, err := cache.Has(id)
	if err != nil || exists {
		t.Fatal("failed generations must leave no completed audio")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = cache.Ensure(ctx, id, func() ([]byte, error) {
		t.Fatal("cancelled requests must not spend credits")

		return nil, nil
	})
	if err == nil {
		t.Fatal("cancellation must be reported")
	}

	if _, err := cache.Open("../../secrets"); err == nil {
		t.Fatal("audio identifiers must not escape the cache")
	}
}
