package media

import (
	"net/url"
	"regexp"
	"strings"

	"chronicle/internal/domain"
)

// Verification records source references, not a platform's paid badge or image authenticity.
func Assess(project *domain.Project) {
	sources := map[string]domain.Source{}
	for _, source := range project.Sources {
		sources[source.ID] = source
	}
	for index := range project.Media {
		asset := &project.Media[index]
		asset.Retrieved = contains(sources[asset.SourceID], asset.URL, asset.Kind)
		if asset.UseAsAvatar && (!asset.Retrieved || !supportedIdentity(*project, asset.EntityID, asset.SourceID, asset.ClaimIDs)) {
			asset.UseAsAvatar = false
		}
	}
	for index := range project.SubjectLinks {
		project.SubjectLinks[index].Verified = false
		project.SubjectLinks[index].VerifiedAt = ""
	}
	// Establish websites before evaluating social accounts linked from those sites.
	for index := range project.SubjectLinks {
		link := &project.SubjectLinks[index]
		source := sources[link.SourceID]
		if link.Kind != "website" || link.VerificationMethod != "official_reference" || host(source.URL) == host(link.URL) {
			continue
		}
		if !supportedIdentity(*project, link.EntityID, link.SourceID, link.ClaimIDs) {
			continue
		}
		for _, resource := range source.Resources {
			if equalURL(resource.URL, link.URL) && officialLabel(resource.Label) {
				link.Verified, link.VerifiedAt = true, source.RetrievedAt
			}
		}
	}
	for index := range project.SubjectLinks {
		link := &project.SubjectLinks[index]
		if link.Kind != "social" || link.VerificationMethod != "official_site_link" {
			continue
		}
		source := sources[link.SourceID]
		for _, website := range project.SubjectLinks {
			if website.Kind == "website" && website.EntityID == link.EntityID && website.Verified && host(website.URL) == host(source.URL) && contains(source, link.URL, "") {
				link.Verified, link.VerifiedAt = true, source.RetrievedAt
			}
		}
	}
}

func supportedIdentity(project domain.Project, entityID, sourceID string, ids []string) bool {
	for _, id := range ids {
		for _, claim := range project.Claims {
			if claim.ID != id || claim.EntityID != entityID || (claim.Status != "corroborated" && claim.Status != "supported_limited") {
				continue
			}
			for _, evidence := range project.Evidence {
				if evidence.ClaimID == id && evidence.SourceID == sourceID && evidence.Verified && evidence.Relation == "supporting" {
					return true
				}
			}
		}
	}
	return false
}

func contains(source domain.Source, address, kind string) bool {
	for _, resource := range source.Resources {
		if equalURL(resource.URL, address) && (kind == "" || kind == resource.Kind) {
			return true
		}
	}
	return false
}

func equalURL(left, right string) bool {
	return strings.TrimRight(left, "/") == strings.TrimRight(right, "/")
}

func host(address string) string {
	parsed, err := url.Parse(address)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
}

var officialWord = regexp.MustCompile(`(?i)(\bofficial\b|\bofficiel(?:le)?\b|\boficial\b|\boffiziell(?:e)?\b|官方网站|公式サイト)`)
var unofficialWord = regexp.MustCompile(`(?i)(unofficial|non[- ]?officiel|not.{0,12}official|fan ?(?:site|page|club))`)

func officialLabel(label string) bool {
	return officialWord.MatchString(label) && !unofficialWord.MatchString(label)
}
