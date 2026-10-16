package domain

import "chronicle/internal/temporal"

type Entity struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Kind       string            `json:"kind"`
	Summary    string            `json:"summary"`
	Time       temporal.Interval `json:"time"`
	Attributes map[string]any    `json:"attributes"`
	Version    int               `json:"version"`
}

type Event struct {
	Narration        *EventNarration   `json:"narration,omitempty"`
	IdentityKey      string            `json:"identityKey,omitempty"`
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	Kind             string            `json:"kind"`
	Summary          string            `json:"summary"`
	Time             temporal.Interval `json:"time"`
	Participants     []string          `json:"participants"`
	LocationIDs      []string          `json:"locationIds"`
	ClaimIDs         []string          `json:"claimIds"`
	Importance       int               `json:"importance"`
	Phase            string            `json:"phase"`
	CauseIDs         []string          `json:"causeIds,omitempty"`
	ConsequenceIDs   []string          `json:"consequenceIds,omitempty"`
	RequiresPresence bool              `json:"requiresPresence,omitempty"`
}

type Relationship struct {
	ID       string            `json:"id"`
	FromID   string            `json:"fromId"`
	ToID     string            `json:"toId"`
	Kind     string            `json:"kind"`
	Time     temporal.Interval `json:"time"`
	ClaimIDs []string          `json:"claimIds"`
}

type Narrative struct {
	NarrativeID string    `json:"narrativeId,omitempty"`
	Chapters    []Chapter `json:"chapters"`
}

type Chapter struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	EventIDs       []string `json:"eventIds"`
	ClaimIDs       []string `json:"claimIds"`
	Interpretation bool     `json:"interpretation"`
}

type Coverage struct {
	StartYear int    `json:"startYear"`
	EndYear   int    `json:"endYear"`
	Status    string `json:"status"`
	Question  string `json:"question"`
	Notes     string `json:"notes"`
}

type Note struct {
	ID        string `json:"id"`
	TargetID  string `json:"targetId"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Project struct {
	Media        []MediaAsset  `json:"media,omitempty"`
	SubjectLinks []SubjectLink `json:"subjectLinks,omitempty"`
	Narrative
	ContextPeriods  []HistoricalContext `json:"contextPeriods,omitempty"`
	ResearchRecords []ResearchRecord    `json:"researchRecords,omitempty"`
	FormatVersion   int                 `json:"formatVersion"`
	ID              string              `json:"id"`
	SubjectID       string              `json:"subjectId"`
	Title           string              `json:"title"`
	Fixture         bool                `json:"fixture"`
	Entities        []Entity            `json:"entities"`
	Events          []Event             `json:"events"`
	Relationships   []Relationship      `json:"relationships"`
	Claims          []Claim             `json:"claims"`
	Sources         []Source            `json:"sources"`
	Evidence        []Evidence          `json:"evidence"`
	Conflicts       []Conflict          `json:"conflicts"`
	Coverage        []Coverage          `json:"coverage"`
	Notes           []Note              `json:"notes"`
	Layout          map[string]Position `json:"layout"`
}
