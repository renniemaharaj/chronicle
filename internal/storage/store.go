package storage

import (
	"context"
	"encoding/json"
	"errors"
)

var ErrNotFound = errors.New("record not found")

type Store interface {
	Get(context.Context, string, string) (json.RawMessage, error)
	Put(context.Context, string, string, any) error
	List(context.Context, string) ([]json.RawMessage, error)
	Close() error
}
