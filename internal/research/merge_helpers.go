package research

import (
	"encoding/json"
	"strings"

	"chronicle/internal/domain"
)

func jsonMarshalProject(project domain.Project) ([]byte, error) { return json.Marshal(project) }
func jsonUnmarshalProject(data []byte, project *domain.Project) error {
	return json.Unmarshal(data, project)
}
func normalizeName(name string) string { return strings.ToLower(strings.TrimSpace(name)) }

func remap(values []string, ids map[string]string) []string {
	result := []string{}
	for _, value := range values {
		if mapped := ids[value]; mapped != "" {
			result = append(result, mapped)
		} else {
			result = append(result, value)
		}
	}
	return unique(result)
}

func unique(values []string) []string {
	result, seen := []string{}, map[string]bool{}
	for _, value := range values {
		if !seen[value] {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result
}

func hasEntity(values []domain.Entity, id string) bool {
	for _, value := range values {
		if value.ID == id {
			return true
		}
	}
	return false
}

func hasEvent(values []domain.Event, id string) bool {
	for _, value := range values {
		if value.ID == id {
			return true
		}
	}
	return false
}

func hasEvidence(values []domain.Evidence, id string) bool {
	for _, value := range values {
		if value.ID == id {
			return true
		}
	}
	return false
}

func hasRelation(values []domain.Relationship, id string) bool {
	for _, value := range values {
		if value.ID == id {
			return true
		}
	}
	return false
}
