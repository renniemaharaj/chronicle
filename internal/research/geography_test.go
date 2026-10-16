package research

import (
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/geography"
)

func TestGeographicReferencesMergeIntoExistingPlacesWithCanonicalProvenance(t *testing.T) {
	project := domain.EmptyProject("synthetic", "Synthetic map reference test")
	project.SubjectID = "place"
	project.Entities = []domain.Entity{{ID: "place", Name: "Synthetic place", Kind: "place", Version: 1}}
	project.Claims = []domain.Claim{{ID: "claim", EntityID: "place"}}
	project.Evidence = []domain.Evidence{{
		ClaimID: "claim", SourceID: "source", Verified: true, Mode: "quotation",
		Relation: "supporting", Excerpt: "Synthetic geographic reference 42.5, 12.25.",
	}}
	incoming := []domain.Entity{{ID: "local-place", Kind: "place", Attributes: map[string]any{
		"mapLocation": map[string]any{"latitude": 42.5, "longitude": 12.25, "kind": "region", "claimIds": []string{"local-claim"}},
	}}}
	ids := map[string]string{"local-place": "place", "local-claim": "claim"}
	source := domain.Source{ID: "source", URL: "https://example.org/geography", Publisher: "Synthetic gazetteer"}
	mergeGeography(&project, incoming, source, ids)
	reference, valid := geography.DecodeReference(project.Entities[0].Attributes["mapLocation"])
	if !valid || reference.SourceID != source.ID || reference.SourceURL != source.URL || reference.ClaimIDs[0] != "claim" {
		t.Fatalf("missing canonical map provenance: %+v", reference)
	}
	incoming[0].Attributes["mapLocation"].(map[string]any)["latitude"] = 43.5
	mergeGeography(&project, incoming, source, ids)
	reference, _ = geography.DecodeReference(project.Entities[0].Attributes["mapLocation"])
	if reference.Latitude != 42.5 {
		t.Fatal("overwrote an existing reference")
	}
}

func TestGeographyMergeOmitsUnsupportedCoordinates(t *testing.T) {
	project := domain.EmptyProject("synthetic", "Synthetic unsupported coordinates")
	project.Entities = []domain.Entity{{ID: "place", Kind: "place", Attributes: map[string]any{}}}
	incoming := []domain.Entity{{ID: "local", Kind: "place", Attributes: map[string]any{
		"mapLocation": map[string]any{"latitude": 42.5, "longitude": 12.25, "claimIds": []string{"missing"}},
	}}}
	mergeGeography(&project, incoming, domain.Source{ID: "source"}, map[string]string{"local": "place"})
	if project.Entities[0].Attributes["mapLocation"] != nil {
		t.Fatal("saved coordinates without supporting quoted evidence")
	}
}
