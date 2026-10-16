package narrative

import (
	"sort"
	"strings"

	"chronicle/internal/domain"
)

func Build(project *domain.Project) {
	EnsureStories(project)
	ordered := append([]domain.Event{}, project.Events...)
	sort.SliceStable(ordered, func(i, j int) bool {
		left, _, known := ordered[i].Time.Bounds()
		right, _, otherKnown := ordered[j].Time.Bounds()
		if !known {
			return false
		}
		if !otherKnown {
			return true
		}
		return left < right
	})
	project.Events = ordered
	project.NarrativeID = "narrative-" + project.ID
	project.Chapters = []domain.Chapter{}
	project.Coverage = []domain.Coverage{}
	for _, event := range ordered {
		phase := event.Phase
		if phase == "" {
			phase = "Documented developments"
		}
		index := -1
		for i, chapter := range project.Chapters {
			if chapter.Title == phase {
				index = i
			}
		}
		if index < 0 {
			project.Chapters = append(project.Chapters, domain.Chapter{ID: "chapter-" + event.ID, Title: phase, EventIDs: []string{}, ClaimIDs: []string{}})
			index = len(project.Chapters) - 1
		}
		chapter := &project.Chapters[index]
		chapter.EventIDs = append(chapter.EventIDs, event.ID)
		for _, claim := range project.Claims {
			if claim.EventID == event.ID && (claim.Status == "supported_limited" || claim.Status == "corroborated" || claim.Status == "disputed") {
				chapter.ClaimIDs = append(chapter.ClaimIDs, claim.ID)
				chapter.Summary = strings.TrimSpace(chapter.Summary + " " + claim.Text)
			}
		}
		if event.Time.Start != nil && event.Time.Start.Precision != "unknown" {
			start, end := event.Time.Start.Year, event.Time.Start.Year
			if event.Time.End != nil {
				end = event.Time.End.Year
			}
			project.Coverage = append(project.Coverage, domain.Coverage{StartYear: start, EndYear: end, Status: "partial", Question: event.Title, Notes: "Evidence was retrieved for this event, not for every intervening historical period."})
		}
	}
	for i := range project.Chapters {
		if project.Chapters[i].Summary == "" {
			project.Chapters[i].Summary = "No evidence-supported narration is available for this chapter."
		}
	}
	for index := 1; index < len(ordered); index++ {
		previous, current := ordered[index-1].Time.Start, ordered[index].Time.Start
		if previous != nil && current != nil && current.Year-previous.Year > 5 {
			project.Coverage = append(project.Coverage, domain.Coverage{StartYear: previous.Year + 1, EndYear: current.Year - 1, Status: "unexplored", Question: "Intervening historical interval", Notes: "No event coverage established by this run. Revisit this interval."})
		}
	}
	if len(project.Coverage) == 0 {
		project.Coverage = append(project.Coverage, domain.Coverage{Status: "unexplored", Notes: "Historical bounds remain unknown."})
	}
}
