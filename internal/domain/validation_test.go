package domain_test

import (
	"testing"

	"chronicle/internal/domain"
	"chronicle/internal/fixtures"
)

func TestMissingProvenanceCannotEnterCanonicalProject(t *testing.T) {
	projects, err := fixtures.Projects()
	if err != nil {
		t.Fatal(err)
	}
	for _, project := range projects {
		if len(project.Evidence) == 0 {
			continue
		}
		project.Evidence[0].SourceID = "fabricated-source"
		if err := domain.ValidateProject(project); err == nil {
			t.Fatal("evidence without a canonical source must be rejected")
		}
		return
	}
	t.Fatal("test requires an evidence-bearing fixture")
}
