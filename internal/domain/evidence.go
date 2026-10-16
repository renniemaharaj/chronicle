package domain

import "chronicle/internal/temporal"

type Claim struct {
	ID          string             `json:"id"`
	Text        string             `json:"text"`
	Status      string             `json:"status"`
	EntityID    string             `json:"entityId,omitempty"`
	EventID     string             `json:"eventId,omitempty"`
	Time        *temporal.Interval `json:"time,omitempty"`
	EvidenceIDs []string           `json:"evidenceIds"`
	Assessment  []string           `json:"assessment"`
}

type Source struct {
	ID                   string                  `json:"id"`
	Title                string                  `json:"title"`
	Author               string                  `json:"author"`
	Publisher            string                  `json:"publisher"`
	URL                  string                  `json:"url"`
	PublishedAt          *temporal.TemporalValue `json:"publishedAt,omitempty"`
	RetrievedAt          string                  `json:"retrievedAt"`
	Classification       string                  `json:"classification"`
	IndependenceGroup    string                  `json:"independenceGroup"`
	IndependenceReviewed bool                    `json:"independenceReviewed"`
	Limitations          []string                `json:"limitations,omitempty"`
	ContentHash          string                  `json:"contentHash,omitempty"`
	Verification         string                  `json:"verification"`
}

type Evidence struct {
	ID       string `json:"id"`
	SourceID string `json:"sourceId"`
	ClaimID  string `json:"claimId"`
	Excerpt  string `json:"excerpt"`
	Location string `json:"location"`
	Relation string `json:"relation"`
	Mode     string `json:"mode"`
	Verified bool   `json:"verified"`
}

type Conflict struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	ClaimIDs     []string `json:"claimIds"`
	Kind         string   `json:"kind"`
	Explanation  string   `json:"explanation"`
	Consequences []string `json:"consequences"`
	Severity     string   `json:"severity"`
}
