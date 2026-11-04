package api

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"chronicle/internal/campaign"
	"chronicle/internal/config"
	"chronicle/internal/research"
	"chronicle/internal/storage"
)

func TestCampaignVideoStaysInsideExportRoot(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	store, err := storage.OpenFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	settings := config.Defaults()
	settings.Workspace = root
	handler := New(store, research.NewManager(store), settings)
	exports := filepath.Join(root, "exports")
	if err := os.MkdirAll(exports, 0700); err != nil {
		t.Fatal(err)
	}

	inside := filepath.Join(exports, "synthetic.mp4")
	outside := filepath.Join(t.TempDir(), "private.mp4")

	for _, path := range []string{inside, outside} {
		if err := os.WriteFile(path, []byte("Synthetic video transport fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	check := func(path string, status int) {
		t.Helper()

		value := campaign.Campaign{
			ID: "video-test",
			Parts: []campaign.Part{{
				Number: 1,
				Status: "recorded",
				File:   path,
			}},
		}
		if err := store.Put(ctx, "campaigns", value.ID, value); err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", "/api/campaigns/video-test/parts/1/video", nil))
		if response.Code != status {
			t.Fatalf("expected %d, got %d", status, response.Code)
		}
	}
	check(inside, 200)
	check(outside, 404)

	link := filepath.Join(exports, "outside-link.mp4")
	if err := os.Symlink(outside, link); err != nil {
		t.Skip("symlinks are unavailable")
	}

	check(link, 404)
}
