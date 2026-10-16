package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"chronicle/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct{ db *sql.DB }

func OpenPostgres(ctx context.Context, url string) (*Postgres, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (s *Postgres) Get(ctx context.Context, kind, id string) (json.RawMessage, error) {
	var data []byte
	err := s.db.QueryRowContext(ctx, "SELECT payload FROM workspace_records WHERE kind=$1 AND id=$2", kind, id).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return data, err
}

func (s *Postgres) Put(ctx context.Context, kind, id string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if kind == "projects" {
		var project domain.Project
		if err := json.Unmarshal(data, &project); err != nil {
			return err
		}
		if err := domain.ValidateProject(project); err != nil {
			return err
		}
		if err := saveProjectRows(ctx, tx, project); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO workspace_records(kind,id,payload) VALUES($1,$2,$3)
		ON CONFLICT(kind,id) DO UPDATE SET payload=excluded.payload, version=workspace_records.version+1, updated_at=now()`, kind, id, data)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Postgres) List(ctx context.Context, kind string) ([]json.RawMessage, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT payload FROM workspace_records WHERE kind=$1 ORDER BY updated_at DESC", kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []json.RawMessage{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		items = append(items, data)
	}
	return items, rows.Err()
}

func (s *Postgres) Close() error { return s.db.Close() }
