package storage

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestAtomicWorkspaceRecordsAndPathRestrictions(t *testing.T) {
	store, err := OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	value := map[string]any{"version": 1, "claim": "A synthetic disputed assertion", "alternatives": []int{1890, 1895}}
	if err := store.Put(context.Background(), "projects", "synthetic", value); err != nil {
		t.Fatal(err)
	}
	data, err := store.Get(context.Background(), "projects", "synthetic")
	if err != nil {
		t.Fatal(err)
	}
	var restored map[string]any
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if len(restored["alternatives"].([]any)) != 2 {
		t.Fatal("competing assertions were lost")
	}
	if _, err := store.Get(context.Background(), "projects", "../settings"); err == nil {
		t.Fatal("record identifiers must not escape the local workspace")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.Put(ctx, "projects", "cancelled", value); !errors.Is(err, context.Canceled) {
		t.Fatal("cancelled requests must not write workspace records")
	}
}
