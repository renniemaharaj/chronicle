package evidence

import "chronicle/internal/domain"

func Assess(project *domain.Project) {
	sources := map[string]domain.Source{}
	for _, source := range project.Sources {
		sources[source.ID] = source
	}
	for index := range project.Claims {
		claim := &project.Claims[index]
		if claim.Status == "rejected_superseded" {
			continue
		}
		groups := map[string]bool{}
		support, oppose := 0, 0
		for _, item := range project.Evidence {
			if item.ClaimID != claim.ID || !item.Verified {
				continue
			}
			source := sources[item.SourceID]
			if item.Relation == "supporting" {
				support++
				if source.IndependenceReviewed && source.IndependenceGroup != "" && source.Classification != "tertiary" {
					groups[source.IndependenceGroup] = true
				}
			}
			if item.Relation == "contradicting" {
				oppose++
			}
		}
		claim.Status = "unverified"
		claim.Assessment = []string{"Passage matching verifies retrieval, not historical truth. Source independence requires review."}
		if support > 0 {
			claim.Status = "supported_limited"
			claim.Assessment = append(claim.Assessment, "At least one retrieved passage supports the assertion; corroboration or independence remains limited.")
		}
		if len(groups) >= 2 {
			claim.Status = "corroborated"
			claim.Assessment = append(claim.Assessment, "Supporting passages come from at least two explicitly reviewed independent source groups.")
		}
		if support > 0 && oppose > 0 {
			claim.Status = "disputed"
			claim.Assessment = append(claim.Assessment, "Retrieved supporting and contradicting passages are both retained.")
		}
	}
}
