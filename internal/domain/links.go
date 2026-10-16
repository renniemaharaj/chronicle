package domain

import "fmt"

func validateLinks(p Project, entities, events, claims, evidence map[string]bool) error {
	for _, event := range p.Events {
		for _, id := range append(append([]string{}, event.CauseIDs...), event.ConsequenceIDs...) {
			if !events[id] {
				return fmt.Errorf("causal ordering references missing event %s", id)
			}
		}
		for _, id := range event.ClaimIDs {
			if !claims[id] {
				return fmt.Errorf("event %s has missing claim %s", event.ID, id)
			}
		}
	}
	for _, claim := range p.Claims {
		for _, id := range claim.EvidenceIDs {
			if !evidence[id] {
				return fmt.Errorf("claim %s has missing evidence %s", claim.ID, id)
			}
			for _, item := range p.Evidence {
				if item.ID == id && item.ClaimID != claim.ID {
					return fmt.Errorf("evidence belongs to another claim")
				}
			}
		}
	}
	for _, relation := range p.Relationships {
		if relation.ID == "" || relation.Kind == "" || !entities[relation.FromID] || !entities[relation.ToID] {
			return fmt.Errorf("invalid relationship %q", relation.ID)
		}
		if err := relation.Time.Validate(); err != nil {
			return err
		}
		for _, id := range relation.ClaimIDs {
			if !claims[id] {
				return fmt.Errorf("relationship has missing claim %s", id)
			}
		}
	}
	for _, conflict := range p.Conflicts {
		if conflict.ID == "" || len(conflict.ClaimIDs) < 2 {
			return fmt.Errorf("conflict requires two competing claims")
		}
		for _, id := range conflict.ClaimIDs {
			if !claims[id] {
				return fmt.Errorf("conflict references missing claim %s", id)
			}
		}
	}
	for _, chapter := range p.Chapters {
		for _, id := range chapter.EventIDs {
			if !events[id] {
				return fmt.Errorf("chapter references missing event %s", id)
			}
		}
		for _, id := range chapter.ClaimIDs {
			if !claims[id] {
				return fmt.Errorf("chapter references missing claim %s", id)
			}
		}
	}
	return nil
}
