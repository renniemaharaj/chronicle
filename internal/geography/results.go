package geography

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
)

type Result struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Label     string  `json:"label"`
	SourceURL string  `json:"sourceURL"`
	Kind      string  `json:"kind"`
}

type feature struct {
	Geometry struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Name     string `json:"name"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		OSMType  string `json:"osm_type"`
		OSMID    int64  `json:"osm_id"`
		OSMValue string `json:"osm_value"`
	} `json:"properties"`
}

func decodeResults(reader io.Reader) ([]Result, error) {
	var body struct {
		Features []feature `json:"features"`
	}
	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		return nil, err
	}

	results := []Result{}
	for _, item := range body.Features {
		coordinates := item.Geometry.Coordinates
		if len(coordinates) < 2 || math.Abs(coordinates[0]) > 180 || math.Abs(coordinates[1]) > 90 {
			continue
		}
		properties := item.Properties
		parts := []string{}
		for _, part := range []string{properties.Name, properties.City, properties.State, properties.Country} {
			if part != "" && (len(parts) == 0 || parts[len(parts)-1] != part) {
				parts = append(parts, part)
			}
		}
		kind := map[string]string{"N": "node", "W": "way", "R": "relation"}[properties.OSMType]
		if kind == "" || properties.OSMID <= 0 || len(parts) == 0 {
			continue
		}
		results = append(results, Result{
			Latitude: coordinates[1], Longitude: coordinates[0], Label: strings.Join(parts, ", "),
			SourceURL: fmt.Sprintf("https://www.openstreetmap.org/%s/%d", kind, properties.OSMID), Kind: properties.OSMValue,
		})
	}
	return results, nil
}
