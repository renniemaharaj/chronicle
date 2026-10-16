package research

import "chronicle/internal/domain"

func mergeContext(project *domain.Project, contexts []domain.HistoricalContext, ids map[string]string) {
	for _, period := range contexts {
		period.ID = stableID("context", period.Title+period.Scope)
		period.ClaimIDs = remap(period.ClaimIDs, ids)
		found := false
		for index, existing := range project.ContextPeriods {
			if existing.ID == period.ID {
				project.ContextPeriods[index].ClaimIDs = unique(append(existing.ClaimIDs, period.ClaimIDs...))
				found = true
			}
		}
		if !found {
			project.ContextPeriods = append(project.ContextPeriods, period)
		}
	}
}
