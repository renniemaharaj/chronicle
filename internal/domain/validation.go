package domain

import (
	"fmt"
	"net/url"
	"strings"
)

func ValidateProject(p Project) error {
	if p.FormatVersion != 1 || p.ID == "" || strings.TrimSpace(p.Title) == "" {
		return fmt.Errorf("project requires formatVersion 1, id, and title")
	}
	if p.Entities == nil || p.Events == nil || p.Claims == nil || p.Sources == nil || p.Evidence == nil || p.Relationships == nil || p.Conflicts == nil || p.Chapters == nil || p.Coverage == nil || p.Notes == nil || p.Layout == nil {
		return fmt.Errorf("canonical collections and layout must not be null")
	}
	entities, events, claims, sources, evidence := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, entity := range p.Entities {
		if entity.ID == "" || entity.Name == "" || entity.Kind == "" || entity.Version < 1 || entities[entity.ID] {
			return fmt.Errorf("invalid or duplicate entity %q", entity.ID)
		}
		if err := entity.Time.Validate(); err != nil {
			return fmt.Errorf("entity %s: %w", entity.ID, err)
		}
		entities[entity.ID] = true
	}
	if !entities[p.SubjectID] {
		return fmt.Errorf("subject entity does not exist")
	}
	for _, event := range p.Events {
		if event.ID == "" || event.Title == "" || event.Kind == "" || events[event.ID] || event.Importance < 0 || event.Importance > 10 {
			return fmt.Errorf("invalid or duplicate event %q", event.ID)
		}
		if err := event.Time.Validate(); err != nil {
			return fmt.Errorf("event %s: %w", event.ID, err)
		}
		for _, id := range append(append([]string{}, event.Participants...), event.LocationIDs...) {
			if !entities[id] {
				return fmt.Errorf("event %s references missing entity %s", event.ID, id)
			}
		}
		events[event.ID] = true
	}
	for _, source := range p.Sources {
		if source.ID == "" || source.Title == "" || sources[source.ID] || !oneOf(source.Classification, "primary", "secondary", "tertiary") {
			return fmt.Errorf("invalid or duplicate source %q", source.ID)
		}
		if source.URL != "" {
			parsed, err := url.Parse(source.URL)
			if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") {
				return fmt.Errorf("invalid source URL")
			}
		}
		sources[source.ID] = true
	}
	for _, claim := range p.Claims {
		if claim.ID == "" || claim.Text == "" || claims[claim.ID] || !oneOf(claim.Status, "corroborated", "supported_limited", "disputed", "unverified", "rejected_superseded") {
			return fmt.Errorf("invalid or duplicate claim %q", claim.ID)
		}
		if claim.EventID != "" && !events[claim.EventID] {
			return fmt.Errorf("claim %s references missing event", claim.ID)
		}
		if claim.EntityID != "" && !entities[claim.EntityID] {
			return fmt.Errorf("claim %s references missing entity", claim.ID)
		}
		if claim.Time != nil {
			if err := claim.Time.Validate(); err != nil {
				return err
			}
		}
		claims[claim.ID] = true
	}
	for _, item := range p.Evidence {
		if item.ID == "" || evidence[item.ID] || !claims[item.ClaimID] || !sources[item.SourceID] || strings.TrimSpace(item.Excerpt) == "" {
			return fmt.Errorf("invalid evidence provenance for %q", item.ID)
		}
		if !oneOf(item.Relation, "supporting", "contradicting", "partially_supporting", "contextual", "inconclusive") || !oneOf(item.Mode, "quotation", "paraphrase") {
			return fmt.Errorf("invalid evidence classification")
		}
		evidence[item.ID] = true
	}
	if err := validateStories(p, claims); err != nil {
		return err
	}
	if err := validateMedia(p, entities, sources, claims); err != nil {
		return err
	}
	return validateLinks(p, entities, events, claims, evidence)
}

func oneOf(value string, values ...string) bool {
	for _, allowed := range values {
		if value == allowed {
			return true
		}
	}
	return false
}
