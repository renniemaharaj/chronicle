package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Files struct {
	root string
	mu   sync.RWMutex
}

var safeID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,120}$`)

func OpenFiles(root string) (*Files, error) {
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
	return &Files{root: root}, nil
}

func (s *Files) path(kind, id string) (string, error) {
	if !safeID.MatchString(kind) || !safeID.MatchString(id) {
		return "", fmt.Errorf("invalid record identifier")
	}
	return filepath.Join(s.root, kind, id+".json"), nil
}

func (s *Files) Get(ctx context.Context, kind, id string) (json.RawMessage, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path, err := s.path(kind, id)
	if err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	return data, err
}

func (s *Files) Put(ctx context.Context, kind, id string, value any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := s.path(kind, id)
	if err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".pending-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if err = file.Chmod(0600); err == nil {
		_, err = file.Write(data)
	}
	if err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(file.Name(), path)
}

func (s *Files) List(ctx context.Context, kind string) ([]json.RawMessage, error) {
	if !safeID.MatchString(kind) {
		return nil, fmt.Errorf("invalid record kind")
	}
	paths, err := filepath.Glob(filepath.Join(s.root, kind, "*.json"))
	if err != nil {
		return nil, err
	}
	items := []json.RawMessage{}
	for _, path := range paths {
		id := filepath.Base(path)
		data, err := s.Get(ctx, kind, id[:len(id)-5])
		if err != nil {
			return nil, err
		}
		items = append(items, data)
	}
	return items, nil
}

func (s *Files) Close() error { return nil }
