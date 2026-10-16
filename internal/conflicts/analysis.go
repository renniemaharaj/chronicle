package conflicts

import (
	"fmt"
	"slices"

	"chronicle/internal/domain"
)

func Analyze(project *domain.Project) {
	detectOrdering(project)
	detectDates(project)
	detectLifespans(project)
	detectParents(project)
	detectPresence(project)
}

func add(project *domain.Project, id, title, kind, explanation, severity string, claimIDs []string) {
	unique := []string{}
	for _, id := range claimIDs {
		if id != "" && !slices.Contains(unique, id) {
			unique = append(unique, id)
		}
	}
	if len(unique) < 2 {
		return
	}
	for _, existing := range project.Conflicts {
		if existing.ID == id {
			return
		}
	}
	project.Conflicts = append(project.Conflicts, domain.Conflict{ID: id, Title: title, Kind: kind, Explanation: explanation, Severity: severity, ClaimIDs: unique, Consequences: []string{"Adopting an alternative can change ordering and interpretation; both assertions and their evidence remain available."}})
}

func markDisputed(project *domain.Project, ids ...string) {
	for index := range project.Claims {
		claim := &project.Claims[index]
		if slices.Contains(ids, claim.ID) && (claim.Status == "supported_limited" || claim.Status == "corroborated") {
			claim.Status = "disputed"
			claim.Assessment = append(claim.Assessment, "Competing asserted date intervals do not overlap. Both sources remain available.")
		}
	}
}

func detectDates(project *domain.Project) {
	for i, first := range project.Claims {
		if first.EventID == "" || first.Time == nil {
			continue
		}
		start, end, known := first.Time.Bounds()
		if !known {
			continue
		}
		for _, second := range project.Claims[i+1:] {
			if second.EventID != first.EventID || second.Time == nil {
				continue
			}
			otherStart, otherEnd, otherKnown := second.Time.Bounds()
			if otherKnown && (end < otherStart || otherEnd < start) {
				markDisputed(project, first.ID, second.ID)
				add(project, "date-"+first.ID+"-"+second.ID, "Competing event dates", "date", "The asserted date intervals do not overlap after calendar normalization. Neither assertion has been discarded.", "definite", []string{first.ID, second.ID})
			}
		}
	}
}

func boundClaims(project *domain.Project, entity domain.Entity, birth bool) []string {
	result := []string{}
	for _, claim := range project.Claims {
		if claim.EntityID == entity.ID {
			result = append(result, claim.ID)
		}
	}
	for _, event := range project.Events {
		if !slices.Contains(event.Participants, entity.ID) || event.Time.Start == nil {
			continue
		}
		bound := entity.Time.End
		if birth {
			bound = entity.Time.Start
		}
		if bound != nil && event.Time.Start.Year == bound.Year {
			result = append(result, event.ClaimIDs...)
		}
	}
	return result
}

func message(entity domain.Entity, text string) string {
	return fmt.Sprintf("%s: %s", entity.Name, text)
}
