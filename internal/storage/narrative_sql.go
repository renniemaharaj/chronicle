package storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"chronicle/internal/domain"
)

func saveNarrativeRows(ctx context.Context, tx *sql.Tx, p domain.Project) error {
	chapters, err := json.Marshal(p.Chapters)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO narratives(id,subject_id,chapters) VALUES($1,$2,$3)
		ON CONFLICT(id) DO UPDATE SET chapters=excluded.chapters,version=narratives.version+1`, "narrative-"+p.ID, p.SubjectID, chapters)
	if err != nil {
		return err
	}
	for _, record := range p.ResearchRecords {
		config, _ := json.Marshal(record.Config)
		_, err := tx.ExecContext(ctx, `INSERT INTO research_runs(id,subject_id,status,config,model,started_at) VALUES($1,$2,$3,$4,$5,$6::timestamptz)
			ON CONFLICT(id) DO UPDATE SET status=excluded.status,config=excluded.config`, record.ID, p.SubjectID, record.Status, config, record.Model, record.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}
