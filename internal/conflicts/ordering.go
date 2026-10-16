package conflicts

import "chronicle/internal/domain"

func detectOrdering(project *domain.Project) {
	events := map[string]domain.Event{}
	for _, event := range project.Events {
		events[event.ID] = event
	}
	for _, effect := range project.Events {
		effectStart, effectEnd, known := effect.Time.Bounds()
		if !known {
			continue
		}
		for _, causeID := range effect.CauseIDs {
			cause, ok := events[causeID]
			if !ok {
				continue
			}
			causeStart, _, causeKnown := cause.Time.Bounds()
			if causeKnown && causeStart > effectEnd {
				add(project, "ordering-"+cause.ID+"-"+effect.ID, "Cause follows its proposed consequence", "causal_order", "The entire documented effect interval precedes the earliest proposed cause. Review ordering and the causal assertion.", "definite", append(append([]string{}, cause.ClaimIDs...), effect.ClaimIDs...))
			}
			if len(effect.ClaimIDs) == 0 || len(cause.ClaimIDs) == 0 {
				continue
			}
			supported := false
			for _, item := range project.Evidence {
				for _, claimID := range effect.ClaimIDs {
					if item.ClaimID == claimID && item.Verified && item.Relation == "supporting" {
						supported = true
					}
				}
			}
			if !supported && causeKnown && causeStart <= effectStart {
				add(project, "causality-"+cause.ID+"-"+effect.ID, "Causal dependency lacks verified evidence", "unsupported_causality", "Chronological ordering alone does not demonstrate causation. No verified supporting passage is attached to the proposed effect.", "suspicious", append(append([]string{}, cause.ClaimIDs...), effect.ClaimIDs...))
			}
		}
	}
}
