package conflicts

import (
	"slices"

	"chronicle/internal/domain"
)

func detectLifespans(project *domain.Project) {
	for _, event := range project.Events {
		start, end, known := event.Time.Bounds()
		if !known {
			continue
		}
		for _, entity := range project.Entities {
			if entity.Kind != "person" || !slices.Contains(event.Participants, entity.ID) {
				continue
			}
			if entity.Time.Start != nil {
				birth, _, known := entity.Time.Start.Bounds()
				if known && end < birth {
					add(project, "before-birth-"+event.ID+"-"+entity.ID, "Participation before documented birth", "lifespan", message(entity, "the entire event interval precedes the documented birth interval."), presenceSeverity(event), append(append([]string{}, event.ClaimIDs...), boundClaims(project, entity, true)...))
				}
			}
			if entity.Time.End != nil {
				_, death, known := entity.Time.End.Bounds()
				if known && start > death && event.RequiresPresence {
					add(project, "after-death-"+event.ID+"-"+entity.ID, "Physical participation after death", "lifespan", message(entity, "physical presence is required after the entire documented death interval. Posthumous publications alone are not contradictions."), "definite", append(append([]string{}, event.ClaimIDs...), boundClaims(project, entity, false)...))
				}
			}
		}
	}
}

func detectParents(project *domain.Project) {
	entities := map[string]domain.Entity{}
	for _, entity := range project.Entities {
		entities[entity.ID] = entity
	}
	for _, relation := range project.Relationships {
		if relation.Kind != "parent of" && relation.Kind != "child of" {
			continue
		}
		parent, child := entities[relation.FromID], entities[relation.ToID]
		if relation.Kind == "child of" {
			parent, child = child, parent
		}
		if parent.Time.Start == nil || child.Time.Start == nil {
			continue
		}
		if parent.Time.Start.DefinitelyAfter(*child.Time.Start) {
			ids := append(boundClaims(project, parent, true), boundClaims(project, child, true)...)
			add(project, "parent-"+relation.ID, "Parent-child chronology", "parenthood", "The parent's documented birth is later than the child's. Check identity, adoption, terminology, and date evidence.", "suspicious", append(ids, relation.ClaimIDs...))
		}
	}
}

func detectPresence(project *domain.Project) {
	for i, event := range project.Events {
		if !event.RequiresPresence || len(event.LocationIDs) == 0 {
			continue
		}
		start, end, known := event.Time.Bounds()
		if !known {
			continue
		}
		for _, other := range project.Events[i+1:] {
			if !other.RequiresPresence || len(other.LocationIDs) == 0 || slices.Equal(event.LocationIDs, other.LocationIDs) {
				continue
			}
			otherStart, otherEnd, otherKnown := other.Time.Bounds()
			if !otherKnown || start > otherEnd || otherStart > end {
				continue
			}
			shared := false
			for _, id := range event.Participants {
				if slices.Contains(other.Participants, id) {
					shared = true
				}
			}
			if shared {
				add(project, "presence-"+event.ID+"-"+other.ID, "Potentially incompatible locations", "presence", "Physical-presence intervals overlap at different recorded locations. Coarse dates, nested places, and travel may explain this; it is not a definite contradiction.", "suspicious", append(append([]string{}, event.ClaimIDs...), other.ClaimIDs...))
			}
		}
	}
}

func presenceSeverity(event domain.Event) string {
	if event.RequiresPresence {
		return "definite"
	}
	return "suspicious"
}
