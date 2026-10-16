package storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"chronicle/internal/domain"
)

func saveProjectRows(ctx context.Context, tx *sql.Tx, p domain.Project) error {
	for _, entity := range p.Entities {
		data, _ := json.Marshal(entity)
		_, err := tx.ExecContext(ctx, `INSERT INTO entities(id,name,kind,attributes,version) VALUES($1,$2,$3,$4,$5)
			ON CONFLICT(id) DO UPDATE SET name=excluded.name,kind=excluded.kind,attributes=excluded.attributes,version=entities.version+1,updated_at=now()`, entity.ID, entity.Name, entity.Kind, data, entity.Version)
		if err != nil {
			return err
		}
	}
	for _, source := range p.Sources {
		published, _ := json.Marshal(source.PublishedAt)
		metadata, _ := json.Marshal(source)
		_, err := tx.ExecContext(ctx, `INSERT INTO sources(id,title,author,publisher,url,published_time,retrieved_at,classification,independence_group,metadata) VALUES($1,$2,$3,$4,$5,$6,NULLIF($7,'')::timestamptz,$8,$9,$10)
			ON CONFLICT(id) DO UPDATE SET title=excluded.title,author=excluded.author,publisher=excluded.publisher,url=excluded.url,classification=excluded.classification,independence_group=excluded.independence_group,metadata=excluded.metadata`, source.ID, source.Title, source.Author, source.Publisher, source.URL, published, source.RetrievedAt, source.Classification, source.IndependenceGroup, metadata)
		if err != nil {
			return err
		}
	}
	for _, event := range p.Events {
		when, _ := json.Marshal(event.Time)
		attributes, _ := json.Marshal(event)
		_, err := tx.ExecContext(ctx, `INSERT INTO events(id,name,temporal_value,attributes) VALUES($1,$2,$3,$4)
			ON CONFLICT(id) DO UPDATE SET name=excluded.name,temporal_value=excluded.temporal_value,attributes=excluded.attributes`, event.ID, event.Title, when, attributes)
		if err != nil {
			return err
		}
		for _, id := range event.Participants {
			if _, err := tx.ExecContext(ctx, `INSERT INTO event_participants(event_id,entity_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, event.ID, id); err != nil {
				return err
			}
		}
	}
	for _, claim := range p.Claims {
		data, _ := json.Marshal(claim)
		_, err := tx.ExecContext(ctx, `INSERT INTO claims(id,text,status,event_id,attributes) VALUES($1,$2,$3,NULLIF($4,''),$5)
			ON CONFLICT(id) DO UPDATE SET text=excluded.text,status=excluded.status,event_id=excluded.event_id,attributes=excluded.attributes`, claim.ID, claim.Text, claim.Status, claim.EventID, data)
		if err != nil {
			return err
		}
	}
	for _, item := range p.Evidence {
		metadata, _ := json.Marshal(item)
		_, err := tx.ExecContext(ctx, `INSERT INTO evidence(id,claim_id,source_id,relation,excerpt,location,metadata) VALUES($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT(id) DO UPDATE SET relation=excluded.relation,excerpt=excluded.excerpt,location=excluded.location,metadata=excluded.metadata`, item.ID, item.ClaimID, item.SourceID, item.Relation, item.Excerpt, item.Location, metadata)
		if err != nil {
			return err
		}
	}
	return saveRelationshipRows(ctx, tx, p)
}

func saveRelationshipRows(ctx context.Context, tx *sql.Tx, p domain.Project) error {
	for _, relation := range p.Relationships {
		from, _ := json.Marshal(relation.Time.Start)
		to, _ := json.Marshal(relation.Time.End)
		data, _ := json.Marshal(relation)
		_, err := tx.ExecContext(ctx, `INSERT INTO relationships(id,from_entity_id,to_entity_id,kind,valid_from,valid_to,attributes) VALUES($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT(id) DO UPDATE SET valid_from=excluded.valid_from,valid_to=excluded.valid_to,attributes=excluded.attributes`, relation.ID, relation.FromID, relation.ToID, relation.Kind, from, to, data)
		if err != nil {
			return err
		}
	}
	for _, conflict := range p.Conflicts {
		metadata, _ := json.Marshal(conflict)
		_, err := tx.ExecContext(ctx, `INSERT INTO conflicts(id,subject,disagreement,metadata) VALUES($1,$2,$3,$4)
			ON CONFLICT(id) DO UPDATE SET subject=excluded.subject,disagreement=excluded.disagreement,metadata=excluded.metadata`, conflict.ID, conflict.Title, conflict.Explanation, metadata)
		if err != nil {
			return err
		}
		for _, id := range conflict.ClaimIDs {
			if _, err := tx.ExecContext(ctx, `INSERT INTO conflict_claims(conflict_id,claim_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, conflict.ID, id); err != nil {
				return err
			}
		}
	}
	return saveNarrativeRows(ctx, tx, p)
}
