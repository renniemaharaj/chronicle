package geography

import (
	"encoding/json"
	"math"
	"regexp"
	"strconv"

	"chronicle/internal/domain"
)

type Reference struct {
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	Label          string   `json:"label"`
	Kind           string   `json:"kind"`
	SourceID       string   `json:"sourceId"`
	SourceURL      string   `json:"sourceURL"`
	ClaimIDs       []string `json:"claimIds"`
	Provider       string   `json:"provider"`
	RetrievedAt    string   `json:"retrievedAt"`
	Note           string   `json:"note"`
	Representative bool     `json:"representative"`
}

func DecodeReference(value any) (Reference, bool) {
	var reference Reference
	data, err := json.Marshal(value)
	if err != nil || json.Unmarshal(data, &reference) != nil || value == nil {
		return reference, false
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(data, &fields) != nil || len(fields["latitude"]) == 0 || len(fields["longitude"]) == 0 {
		return reference, false
	}
	for _, key := range []string{"latitude", "longitude"} {
		if string(fields[key]) == "null" {
			return reference, false
		}
	}

	return reference, math.Abs(reference.Latitude) <= 90 && math.Abs(reference.Longitude) <= 180
}

var decimalCoordinate = regexp.MustCompile(`[-+]?\d+(?:\.\d+)?`)

// Require both decimal coordinates in a verified passage. Do not convert ancient
// coordinate systems or infer an unprinted position.
func SupportedReference(project domain.Project, entityID string, reference Reference) bool {
	for _, claimID := range reference.ClaimIDs {
		for _, claim := range project.Claims {
			if claim.ID != claimID || claim.EntityID != entityID {
				continue
			}
			for _, item := range project.Evidence {
				if item.ClaimID != claimID || item.SourceID != reference.SourceID || !item.Verified || item.Mode != "quotation" || item.Relation != "supporting" {
					continue
				}
				latitude, longitude := false, false
				for _, token := range decimalCoordinate.FindAllString(item.Excerpt, -1) {
					number, _ := strconv.ParseFloat(token, 64)
					latitude = latitude || number == reference.Latitude
					longitude = longitude || number == reference.Longitude
				}
				if latitude && longitude {
					return true
				}
			}
		}
	}

	return false
}

func WithoutReference(attributes map[string]any) map[string]any {
	clean := make(map[string]any, len(attributes))
	for key, value := range attributes {
		if key != "mapLocation" {
			clean[key] = value
		}
	}

	return clean
}
