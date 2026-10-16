package api

import (
	"context"
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/fixtures"
	"chronicle/internal/storage"
)

func TestAssistantCannotInventVerifiedMediaResources(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenFiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	projects, err := fixtures.Projects()
	if err != nil {
		t.Fatal(err)
	}
	project := projects[0]
	project.Fixture = false
	project.Sources = []domain.Source{{ID: "untrusted-source", Title: "Untrusted external reference", URL: "https://archive.example/profile", Classification: "secondary"}}
	source := &project.Sources[0]
	source.Resources = []domain.Resource{{Kind: "image", URL: "https://unverified.example/photo.jpg", Label: "Invented portrait"},
		{Kind: "link", URL: "https://unverified.example/", Label: "Official website"}}
	project.Media = []domain.MediaAsset{{ID: "photo", SourceID: source.ID, EntityID: project.SubjectID,
		Kind: "image", URL: "https://unverified.example/photo.jpg", Title: "Unsupported image", ClaimIDs: []string{project.Claims[0].ID}, UseAsAvatar: true, Retrieved: true}}
	project.SubjectLinks = []domain.SubjectLink{{ID: "site", SourceID: source.ID, EntityID: project.SubjectID,
		Kind: "website", URL: "https://unverified.example/", Label: "Unverified website", ClaimIDs: []string{project.Claims[0].ID}, Verified: true, VerificationMethod: "official_reference"}}
	server := &Server{store: store}
	validated, err := server.prepareAssistantProject(ctx, projectInput{Project: project})
	if err != nil {
		t.Fatal(err)
	}
	if validated.Media[0].Retrieved || validated.Media[0].UseAsAvatar || validated.SubjectLinks[0].Verified {
		t.Fatal("caller-provided resource lists or verification flags became evidence")
	}
	if len(validated.Sources[0].Resources) != 0 {
		t.Fatal("invented source resources were preserved")
	}
}
