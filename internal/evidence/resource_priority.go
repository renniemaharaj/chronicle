package evidence

import (
	"net/url"
	"strings"

	"chronicle/internal/domain"
)

// Keep useful footer references even when a long article has hundreds of navigation links.
func keepResource(resources []domain.Resource, candidate domain.Resource) []domain.Resource {
	for index, existing := range resources {
		if existing.URL == candidate.URL && existing.Kind == candidate.Kind {
			if resourceScore(candidate) > resourceScore(existing) {
				resources[index] = candidate
			}
			return resources
		}
	}
	if len(resources) < 180 {
		return append(resources, candidate)
	}
	lowest := 0
	for index := range resources {
		if resourceScore(resources[index]) < resourceScore(resources[lowest]) {
			lowest = index
		}
	}
	if resourceScore(candidate) > resourceScore(resources[lowest]) {
		resources[lowest] = candidate
	}
	return resources
}

func resourceScore(resource domain.Resource) int {
	label, address := strings.ToLower(resource.Label), strings.ToLower(resource.URL)
	score := 0
	if resource.Kind == "image" {
		score += 20
	}
	if resource.Kind == "video" {
		score += 60
	}
	if strings.Contains(label, "lead image") {
		score += 45
	}
	if strings.Contains(label, "official") || strings.Contains(label, "officiel") {
		score += 100
	}
	if strings.Contains(address, "instagram.com/") || strings.Contains(address, "facebook.com/") || strings.Contains(address, "twitter.com/") {
		score += 50
	}
	if strings.Contains(address, "creativecommons.org/licenses/") {
		score += 30
	}
	if strings.Contains(address, "/wiki/special:") || strings.Contains(address, "commonswiki") || strings.Contains(label, "logo") {
		score -= 60
	}
	return score
}

func videoAddress(address string) bool {
	parsed, err := url.Parse(address)
	if err != nil || !videoHost(parsed.Hostname()) {
		return false
	}
	if strings.Contains(parsed.Hostname(), "youtu.be") {
		return strings.Trim(parsed.Path, "/") != ""
	}
	return parsed.Query().Get("v") != "" || strings.Contains(parsed.Path, "/embed/") || strings.Contains(parsed.Path, "/shorts/") || strings.Contains(parsed.Path, "/video/")
}
