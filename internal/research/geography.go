package research

import (
	"chronicle/internal/domain"
	"chronicle/internal/geography"
)

func mergeGeography(project *domain.Project, entities []domain.Entity, source domain.Source, ids map[string]string) {
	for _, incoming := range entities {
		if incoming.Kind != "place" {
			continue
		}
		reference, valid := geography.DecodeReference(incoming.Attributes["mapLocation"])
		if !valid {
			continue
		}
		reference.ClaimIDs = remap(reference.ClaimIDs, ids)
		reference.SourceID = source.ID
		reference.SourceURL = source.URL
		reference.Provider = source.Publisher
		reference.RetrievedAt = source.RetrievedAt
		reference.Representative = true
		entityID := ids[incoming.ID]
		if !geography.SupportedReference(*project, entityID, reference) {
			continue
		}
		for index := range project.Entities {
			entity := &project.Entities[index]
			if entity.ID != entityID {
				continue
			}
			if _, located := geography.DecodeReference(entity.Attributes["mapLocation"]); located {
				break
			}
			if reference.Label == "" {
				reference.Label = entity.Name
			}
			if entity.Attributes == nil {
				entity.Attributes = map[string]any{}
			}
			entity.Attributes["mapLocation"] = reference
			entity.Version++
			break
		}
	}
}
