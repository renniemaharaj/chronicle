package fixtures

import (
	"encoding/json"
	"reflect"
	"testing"

	"chronicle/internal/domain"
)

func TestUniversalFixturesValidateAndRoundTrip(t *testing.T) {
	projects, err := Projects()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 4 {
		t.Fatal("person, family, organization, and place fixtures are required")
	}
	for _, project := range projects {
		data, err := json.Marshal(project)
		if err != nil {
			t.Fatal(err)
		}
		var restored domain.Project
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(project, restored) {
			t.Fatal("portable project format lost records or uncertainty")
		}
	}
}
