package research

import "chronicle/internal/domain"

func existingOccurrences(project *domain.Project) []map[string]string {
	result := []map[string]string{}

	for _, event := range project.Events {
		result = append(result, map[string]string{
			"id":          event.ID,
			"title":       event.Title,
			"identityKey": event.IdentityKey,
		})
	}

	return result
}

func existingOccurrence(project domain.Project, event domain.Event) string {
	for _, candidate := range project.Events {
		if event.ID == candidate.ID {
			return candidate.ID
		}

		if event.IdentityKey != "" && event.IdentityKey == candidate.IdentityKey && event.Kind == candidate.Kind {
			return candidate.ID
		}
	}

	return ""
}
