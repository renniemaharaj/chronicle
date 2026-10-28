package research

import "chronicle/internal/domain"

func mergeEvents(result *domain.Project, batch extraction, ids map[string]string) {
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
					previous := len(result.Events[i].ClaimIDs)
					result.Events[i].ClaimIDs = unique(append(result.Events[i].ClaimIDs, event.ClaimIDs...))
					if len(result.Events[i].ClaimIDs) > previous {
						// Rebuild the reading so added or competing evidence reaches the story.
						result.Events[i].Narration = nil
					}
				}
			}
		}
	}

}
