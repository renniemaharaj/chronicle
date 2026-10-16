package research

import (
	"strings"
	"testing"
)

func TestResearchQuestionsCoverHistoricalConditionsAndAftermath(t *testing.T) {
	person := historicalQuestions("A historical figure", "person", nil)
	if len(person) != 6 || !strings.Contains(person[0].Question, "technologies") || !strings.Contains(person[4].Question, "underlying causes") {
		t.Fatal("research omitted origin context or documented death investigation")
	}
	institution := historicalQuestions("An institution", "organization", nil)
	if strings.Contains(institution[0].Question, "childhood") || !strings.Contains(institution[4].Question, "dissolution") || !strings.Contains(institution[5].Question, "books") {
		t.Fatal("research imposed a human biography or omitted broad source discovery")
	}
}
