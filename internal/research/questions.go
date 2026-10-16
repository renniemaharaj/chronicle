package research

import (
	"fmt"
	"strings"

	"chronicle/internal/domain"
)

// Six evidence questions cover a subject's whole existence without assuming a human life.
func historicalQuestions(name, kind string, proposed []domain.Question) []domain.Question {
	origins := "earliest attested origins, precursors, formation, enabling conditions and competing origin accounts"
	maturity := "development, adaptation, institutional changes, adoption and succession"
	ending := "latest documented state; dissolution, destruction, extinction or abandonment if established; successors and effects on related entities and places"
	if strings.EqualFold(kind, "person") {
		origins = "documented birth, household, childhood and early influences; separate antecedent conditions from proven causes"
		maturity = "education, training, work, maturation, decisions, setbacks and achievements"
		ending = "latest public activity if living; if deceased, establish death date, place and publicly documented cause without diagnosis speculation; distinguish immediate and underlying causes; investigate mourning, dependants, collaborators, institutions, places and later legacy"
	}
	topics := []struct{ label, question string }{
		{"Origins and historical setting", "Investigate " + origins + ". Describe the period's technologies, communication, transport, material conditions, cultural institutions and geography using dated evidence. Avoid treating an era label as an explanation."},
		{"Development and maturation", "Trace " + maturity + ". Seek lesser-known intervening periods as well as celebrated milestones."},
		{"Relationships and influences", "Document family or constituent relationships, collaborators, rivals, ownership, membership, places and contemporary surroundings as applicable. Separate documented influence from temporal association."},
		{"Accomplishments and turning points", "Examine major works, failures, changes, conflicts, decisions, adoption and consequences. Seek original records alongside scholarly interpretations and retain rival explanations."},
		{"Endings and aftermath", "Investigate " + ending + ". Follow affected subjects beyond the primary subject's end without implying that the subject was still alive or active."},
		{"Public record and media", "Identify official website references and social accounts linked by established official sites. Collect identified images, credits, licenses and external video links from retrieved sources. Seek books, archives, scholarly work, wikis and public forums; treat forum and social reports as attributed leads until independently supported. Record unavailable full texts explicitly."},
	}
	questions := make([]domain.Question, 0, len(topics))
	for index, topic := range topics {
		question := fmt.Sprintf("%s (%s): %s", name, kind, topic.question)
		if index == 2 {
			question += " Resolve event places using geographic sources with explicit coordinates; include coordinate quotations and representative regional references. Preserve ambiguous place identities as gaps."
		}
		if index < len(proposed) {
			question += " Additional research lead: " + proposed[index].Question
		}
		questions = append(questions, domain.Question{Question: question, Interval: topic.label, Status: "pending"})
	}
	return questions
}
