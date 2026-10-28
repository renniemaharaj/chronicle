package speech

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var audioID = regexp.MustCompile(`^[a-f0-9]{64}$`)

type Cache struct {
	root string
	gate chan struct{}
}

func NewCache(workspace string) *Cache {
	return &Cache{
		root: filepath.Join(workspace, "audio"),
		gate: make(chan struct{}, 1),
	}
}

func (c *Cache) Key(input Request, endpoint string) string {
	// Version covers the fixed encoding and provider synthesis settings.
	data, _ := json.Marshal([]any{"mp3-v1-stability0.5-similarity0.75", endpoint, input.Normalized()})
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}

func (c *Cache) Open(id string) (*os.File, error) {
	if !audioID.MatchString(id) {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(filepath.Join(c.root, id+".mp3"))
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil || stat.Size() == 0 || !stat.Mode().IsRegular() {
		file.Close()

		return nil, fmt.Errorf("saved audio is empty or unreadable")
	}

	return file, nil
}

func (c *Cache) Has(id string) (bool, error) {
	file, err := c.Open(id)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, file.Close()
}

func (c *Cache) Ensure(ctx context.Context, id string, generate func() ([]byte, error)) (bool, error) {
	// Serialize generation and recheck disk so simultaneous requests spend credits once.
	select {
	case c.gate <- struct{}{}:
		defer func() {
			<-c.gate
		}()
	case <-ctx.Done():
		return false, ctx.Err()
	}

	exists, err := c.Has(id)
	if err != nil || exists {
		return exists, err
	}

	if err := ctx.Err(); err != nil {
		return false, err
	}

	data, err := generate()
	if err != nil {
		return false, err
	}

	if len(data) == 0 || len(data) > 16<<20 {
		return false, fmt.Errorf("speech returned empty or oversized audio")
	}

	return false, c.save(id, data)
}

func (c *Cache) save(id string, data []byte) error {
	if err := os.MkdirAll(c.root, 0700); err != nil {
		return err
	}

	file, err := os.CreateTemp(c.root, ".audio-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	_, err = file.Write(data)
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

	return os.Rename(file.Name(), filepath.Join(c.root, id+".mp3"))
}
