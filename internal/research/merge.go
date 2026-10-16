package research

import (
	"fmt"

	"chronicle/internal/domain"
	"chronicle/internal/geography"
)

func mergeBatch(project domain.Project, batch extraction, source domain.Source) (domain.Project, error) {
	if len(batch.Events) > 30 || len(batch.Claims) > 50 || len(batch.Entities) > 40 {
		return project, fmt.Errorf("extraction exceeded bounded record limits")
	}
	// Clone the checkpoint so a failed merge cannot partially change the live record.
	var result domain.Project
	data, _ := jsonMarshalProject(project)
	if err := jsonUnmarshalProject(data, &result); err != nil {
		return project, err
	}
	ids := map[string]string{}
	for _, entity := range batch.Entities {
		old := entity.ID
		entity.ID = stableID("entity", entity.Kind+":"+entity.Name)
		entity.Version = 1
		entity.Attributes = geography.WithoutReference(entity.Attributes)
		for index, existing := range result.Entities {
			if normalizeName(existing.Name) == normalizeName(entity.Name) && existing.Kind == entity.Kind {
				entity.ID = existing.ID
				if existing.Time.Start == nil {
					result.Entities[index].Time.Start = entity.Time.Start
				}
				if existing.Time.End == nil {
					result.Entities[index].Time.End = entity.Time.End
				}
			}
		}
		ids[old] = entity.ID
		if !hasEntity(result.Entities, entity.ID) {
			result.Entities = append(result.Entities, entity)
		}
	}
	for _, event := range batch.Events {
		ids[event.ID] = stableID("event", event.Title+event.Kind+event.IdentityKey+inputJSON(remap(event.Participants, ids)))
	}
	for _, claim := range batch.Claims {
		ids[claim.ID] = stableID("claim", claim.Text+ids[claim.EventID])
	}
	for _, item := range batch.Evidence {
		ids[item.ID] = stableID("evidence", source.ID+ids[item.ClaimID]+item.Excerpt)
	}
	for _, event := range batch.Events {
		event.ID = ids[event.ID]
		if event.Narration != nil {
			event.Narration.ClaimIDs = remap(event.Narration.ClaimIDs, ids)
		}
		event.Participants = remap(event.Participants, ids)
		event.LocationIDs = remap(event.LocationIDs, ids)
		event.ClaimIDs = remap(event.ClaimIDs, ids)
		event.CauseIDs = remap(event.CauseIDs, ids)
		event.ConsequenceIDs = remap(event.ConsequenceIDs, ids)
		if !hasEvent(result.Events, event.ID) {
			result.Events = append(result.Events, event)
		} else {
			for i := range result.Events {
				if result.Events[i].ID == event.ID {
					result.Events[i].ClaimIDs = unique(append(result.Events[i].ClaimIDs, event.ClaimIDs...))
				}
			}
		}
	}
	for _, claim := range batch.Claims {
		claim.ID, claim.EventID, claim.EntityID = ids[claim.ID], ids[claim.EventID], ids[claim.EntityID]
		if claim.Time == nil {
			for _, event := range batch.Events {
				if ids[event.ID] == claim.EventID {
					value := event.Time
					claim.Time = &value
					break
				}
			}
		}
		claim.Status = "unverified"
		claim.EvidenceIDs = []string{}
		for _, item := range batch.Evidence {
			if ids[item.ClaimID] == claim.ID {
				claim.EvidenceIDs = append(claim.EvidenceIDs, ids[item.ID])
			}
		}
		found := false
		for index, existing := range result.Claims {
			if existing.ID == claim.ID {
				result.Claims[index].EvidenceIDs = unique(append(existing.EvidenceIDs, claim.EvidenceIDs...))
				found = true
			}
		}
		if !found {
			result.Claims = append(result.Claims, claim)
		}
	}
	for _, item := range batch.Evidence {
		item.ID, item.ClaimID, item.SourceID = ids[item.ID], ids[item.ClaimID], source.ID
		if !hasEvidence(result.Evidence, item.ID) {
			result.Evidence = append(result.Evidence, item)
		}
	}
	found := false
	for _, existing := range result.Sources {
		if existing.ID == source.ID {
			found = true
		}
	}
	if !found {
		result.Sources = append(result.Sources, source)
	}
	for _, relation := range batch.Relationships {
		relation.FromID, relation.ToID = ids[relation.FromID], ids[relation.ToID]
		relation.ID = stableID("relationship", relation.FromID+relation.Kind+relation.ToID)
		relation.ClaimIDs = remap(relation.ClaimIDs, ids)
		if !hasRelation(result.Relationships, relation.ID) {
			result.Relationships = append(result.Relationships, relation)
		}
	}
	for _, conflict := range batch.Conflicts {
		conflict.ClaimIDs = remap(conflict.ClaimIDs, ids)
		conflict.ID = stableID("conflict", conflict.Title)
		conflict.Severity = "historical"
		result.Conflicts = append(result.Conflicts, conflict)
	}
	mergeContext(&result, batch.ContextPeriods, ids)
	mergeGeography(&result, batch.Entities, source, ids)
	mergeMedia(&result, batch, source.ID, ids)
	if err := domain.ValidateProject(result); err != nil {
		return project, err
	}
	return result, nil
}
