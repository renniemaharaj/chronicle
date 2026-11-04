package domain

import (
	"chronicle/internal/temporal"
	"fmt"
)

// Resources are links observed in retrieved documents, not statements of identity.
type Resource struct {
	Kind  string `json:"kind"`
	URL   string `json:"url"`
	Label string `json:"label"`
}

type MediaAsset struct {
	Generated   *ImageGeneration   `json:"generated,omitempty"`
	Depicts     *temporal.Interval `json:"depicts,omitempty"`
	ID          string             `json:"id"`
	EntityID    string             `json:"entityId"`
	SourceID    string             `json:"sourceId"`
	Kind        string             `json:"kind"`
	URL         string             `json:"url"`
	Title       string             `json:"title"`
	Caption     string             `json:"caption"`
	Creator     string             `json:"creator"`
	License     string             `json:"license"`
	LicenseURL  string             `json:"licenseUrl,omitempty"`
	ClaimIDs    []string           `json:"claimIds"`
	UseAsAvatar bool               `json:"useAsAvatar"`
	Retrieved   bool               `json:"retrieved"`
}

type SubjectLink struct {
	ID                 string   `json:"id"`
	EntityID           string   `json:"entityId"`
	SourceID           string   `json:"sourceId"`
	Kind               string   `json:"kind"`
	URL                string   `json:"url"`
	Label              string   `json:"label"`
	ClaimIDs           []string `json:"claimIds"`
	Verified           bool     `json:"verified"`
	VerificationMethod string   `json:"verificationMethod"`
	VerifiedAt         string   `json:"verifiedAt,omitempty"`
}

func validateMedia(p Project, entities, sources, claims map[string]bool) error {
	seen := map[string]bool{}
	check := func(id, entityID, sourceID, address string, ids []string) error {
		if id == "" || seen[id] || !entities[entityID] || !sources[sourceID] || !SafeWebURL(address) {
			return fmt.Errorf("invalid media/link identity or provenance: %s", id)
		}

		seen[id] = true

		for _, claimID := range ids {
			if !claims[claimID] {
				return fmt.Errorf("media/link %s references missing claim %s", id, claimID)
			}
		}

		return nil
	}

	for _, asset := range p.Media {
		if asset.Depicts != nil {
			if err := asset.Depicts.Validate(); err != nil {
				return err
			}
		}

		if asset.LicenseURL != "" && !SafeWebURL(asset.LicenseURL) {
			return fmt.Errorf("invalid media license URL")
		}

		if asset.Generated != nil {
			if err := validateGeneratedImage(asset, entities, sources, claims); err != nil {
				return err
			}

			if seen[asset.ID] {
				return fmt.Errorf("duplicate media identity: %s", asset.ID)
			}

			seen[asset.ID] = true
			continue
		}

		if err := check(asset.ID, asset.EntityID, asset.SourceID, asset.URL, asset.ClaimIDs); err != nil {
			return err
		}

		if !oneOf(asset.Kind, "image", "video") || asset.Title == "" || (asset.UseAsAvatar && (asset.Kind != "image" || len(asset.ClaimIDs) == 0)) {
			return fmt.Errorf("invalid media kind, title or avatar evidence: %s", asset.ID)
		}

	}

	for _, link := range p.SubjectLinks {
		if err := check(link.ID, link.EntityID, link.SourceID, link.URL, link.ClaimIDs); err != nil {
			return err
		}

		if !oneOf(link.Kind, "website", "social") || link.Label == "" || !oneOf(link.VerificationMethod, "unverified", "official_reference", "official_site_link") {
			return fmt.Errorf("invalid subject link type or verification method: %s", link.ID)
		}
	}

	return nil
}
