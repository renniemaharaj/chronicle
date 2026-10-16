package geography

import (
	"testing"

	"chronicle/internal/domain"
)

func TestReferenceRequiresExplicitCoordinates(t *testing.T) {
	for _, value := range []any{nil, map[string]any{}, map[string]any{"latitude": nil, "longitude": 0}, map[string]any{"latitude": 91, "longitude": 0}} {
		if _, valid := DecodeReference(value); valid {
			t.Fatalf("accepted missing or invalid coordinates: %v", value)
		}
	}
	if _, valid := DecodeReference(map[string]any{"latitude": 0, "longitude": 0}); !valid {
		t.Fatal("rejected explicitly supplied zero coordinates")
	}
}

func TestReferenceRequiresEntityEvidenceAndQuotedCoordinateValues(t *testing.T) {
	reference := Reference{Latitude: 42.5, Longitude: -12.25, SourceID: "source", ClaimIDs: []string{"claim"}}
	project := domain.Project{
		Claims: []domain.Claim{{ID: "claim", EntityID: "place"}},
		Evidence: []domain.Evidence{{
			ClaimID: "claim", SourceID: "source", Verified: true,
			Mode: "quotation", Relation: "supporting", Excerpt: "Synthetic coordinate reference 42.5, -12.25.",
		}},
	}
	if !SupportedReference(project, "place", reference) {
		t.Fatal("rejected sourced decimal coordinates")
	}
	if SupportedReference(project, "other-place", reference) {
		t.Fatal("accepted another entity's evidence")
	}
	reference.Longitude = 12.25
	if SupportedReference(project, "place", reference) {
		t.Fatal("accepted a coordinate with the wrong sign")
	}
	reference.Longitude = -12.25
	project.Evidence[0].Verified = false
	if SupportedReference(project, "place", reference) {
		t.Fatal("accepted an unverified passage")
	}
}
