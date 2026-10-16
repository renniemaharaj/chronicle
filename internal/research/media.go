package research

import (
	"chronicle/internal/domain"
)

func mergeMedia(project *domain.Project, batch extraction, sourceID string, ids map[string]string) {
	for _, asset := range batch.Media {
		asset.ID = stableID("media", asset.URL+ids[asset.EntityID])
		asset.EntityID, asset.SourceID = remappedID(asset.EntityID, ids), sourceID
		asset.ClaimIDs = remap(asset.ClaimIDs, ids)
		found := false
		for _, existing := range project.Media {
			if existing.ID == asset.ID {
				found = true
			}
		}
		if !found {
			project.Media = append(project.Media, asset)
		}
	}
	for _, link := range batch.SubjectLinks {
		link.ID = stableID("link", link.URL+ids[link.EntityID])
		link.EntityID, link.SourceID = remappedID(link.EntityID, ids), sourceID
		link.ClaimIDs = remap(link.ClaimIDs, ids)
		found := false
		for _, existing := range project.SubjectLinks {
			if existing.ID == link.ID {
				found = true
			}
		}
		if !found {
			project.SubjectLinks = append(project.SubjectLinks, link)
		}
	}
}

func remappedID(id string, ids map[string]string) string {
	if ids[id] != "" {
		return ids[id]
	}
	return id
}
