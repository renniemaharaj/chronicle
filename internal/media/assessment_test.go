package media

import (
	"chronicle/internal/domain"
	"testing"
)

func TestOnlyEvidenceBackedOfficialReferencesEstablishAccounts(t *testing.T) {
	project := domain.Project{
		Sources: []domain.Source{
			{ID: "authority", URL: "https://archive.example/person", RetrievedAt: "2026-10-16", Resources: []domain.Resource{{Kind: "link", URL: "https://official.example/", Label: "Official website"}}},
			{ID: "website", URL: "https://official.example/about", RetrievedAt: "2026-10-16", Resources: []domain.Resource{{Kind: "link", URL: "https://social.example/person", Label: "Social account"}, {Kind: "image", URL: "https://images.example/portrait.jpg"}}},
		},
		Claims:   []domain.Claim{{ID: "identity", EntityID: "person", Status: "supported_limited"}},
		Evidence: []domain.Evidence{{ClaimID: "identity", SourceID: "authority", Verified: true, Relation: "supporting"}},
		SubjectLinks: []domain.SubjectLink{
			{ID: "social", EntityID: "person", SourceID: "website", Kind: "social", URL: "https://social.example/person", VerificationMethod: "official_site_link"},
			{ID: "official", EntityID: "person", SourceID: "authority", Kind: "website", URL: "https://official.example/", ClaimIDs: []string{"identity"}, VerificationMethod: "official_reference"},
			{ID: "fake", EntityID: "person", SourceID: "website", Kind: "social", URL: "https://fake.example/person", Verified: true, VerificationMethod: "official_site_link"},
		},
	}
	Assess(&project)
	if !project.SubjectLinks[0].Verified || !project.SubjectLinks[1].Verified || project.SubjectLinks[2].Verified {
		t.Fatal("official reference chain was not enforced")
	}
	project.Claims[0].Status = "unverified"
	Assess(&project)
	if project.SubjectLinks[0].Verified || project.SubjectLinks[1].Verified {
		t.Fatal("unverified authority established an account")
	}
}

func TestAvatarNeedsRetrievedImageAndIdentifyingEvidence(t *testing.T) {
	project := domain.Project{Sources: []domain.Source{{ID: "source", Resources: []domain.Resource{{Kind: "image", URL: "https://images.example/photo.jpg"}}}},
		Media:    []domain.MediaAsset{{ID: "photo", EntityID: "person", SourceID: "source", Kind: "image", URL: "https://images.example/photo.jpg", UseAsAvatar: true, ClaimIDs: []string{"caption"}}},
		Claims:   []domain.Claim{{ID: "caption", EntityID: "person", Status: "supported_limited"}},
		Evidence: []domain.Evidence{{ClaimID: "caption", SourceID: "source", Verified: true, Relation: "supporting"}},
	}
	Assess(&project)
	if !project.Media[0].UseAsAvatar || !project.Media[0].Retrieved {
		t.Fatal("supported portrait was lost")
	}
	project.Sources[0].Resources = nil
	Assess(&project)
	if project.Media[0].Retrieved || project.Media[0].UseAsAvatar {
		t.Fatal("invented media URL became a portrait")
	}
}

func TestOfficialReferenceLabelsRejectFanSites(t *testing.T) {
	for _, label := range []string{"Official website", "Site officiel", "Sitio oficial"} {
		if !officialLabel(label) {
			t.Errorf("official label not recognized: %s", label)
		}
	}
	for _, label := range []string{"Unofficial website", "Not the official website", "Official fan page", "Site non-officiel"} {
		if officialLabel(label) {
			t.Errorf("unverified fan reference accepted: %s", label)
		}
	}
}
