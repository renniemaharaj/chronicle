package campaign

import (
	"chronicle/internal/domain"
	"chronicle/internal/temporal"
)

// Campaigns reference canonical records and freeze editorial passages for reproducible takes.
type Campaign struct {
	Accounts        map[string]Destination `json:"accounts,omitempty"`
	Platforms       []string               `json:"platforms,omitempty"`
	ReadyThrough    string                 `json:"readyThrough,omitempty"`
	RemainingEvents int                    `json:"remainingEvents,omitempty"`
	Engine          string                 `json:"engine,omitempty"`
	ID              string                 `json:"id"`
	ProjectID       string                 `json:"projectId"`
	ProjectHash     string                 `json:"projectHash"`
	Title           string                 `json:"title"`
	CreatedAt       string                 `json:"createdAt"`
	Revision        int                    `json:"revision"`
	Seconds         int                    `json:"seconds"`
	Parts           []Part                 `json:"parts"`
	Gaps            []string               `json:"gaps"`
	Audit           []Activity             `json:"audit"`
	Sources         []Reference            `json:"sources"`
}

type Reference struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Hash  string `json:"hash"`
}

type Passage struct {
	Visual     *FrozenVisual     `json:"visual,omitempty"`
	AudioID    string            `json:"audioId,omitempty"`
	AudioStart float64           `json:"audioStart,omitempty"`
	AudioEnd   float64           `json:"audioEnd,omitempty"`
	EventID    string            `json:"eventId,omitempty"`
	ContextID  string            `json:"contextId,omitempty"`
	Title      string            `json:"title"`
	Text       string            `json:"text"`
	Start      int               `json:"start"`
	End        int               `json:"end"`
	ClaimIDs   []string          `json:"claimIds"`
	Time       temporal.Interval `json:"time"`
}

type FrozenVisual struct {
	Cue   domain.StoryVisual `json:"cue"`
	Image domain.MediaAsset  `json:"image"`
}

type Part struct {
	Replacements []Replacement      `json:"replacements,omitempty"`
	Publications map[string]Receipt `json:"publications,omitempty"`
	Number       int                `json:"number"`
	Status       string             `json:"status"`
	Passages     []Passage          `json:"passages"`
	TakeID       string             `json:"takeId,omitempty"`
	File         string             `json:"file,omitempty"`
	Duration     float64            `json:"duration,omitempty"`
	Publication  *Receipt           `json:"publication,omitempty"`
	Metrics      []Metrics          `json:"metrics"`
	Takes        []Take             `json:"takes,omitempty"`
}

type Take struct {
	ID          string  `json:"id"`
	File        string  `json:"file"`
	Duration    float64 `json:"duration"`
	CompletedAt string  `json:"completedAt"`
}

type Receipt struct {
	File       string  `json:"file,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
	URL        string  `json:"url"`
	RecordedAt string  `json:"recordedAt"`
	Method     string  `json:"method"`
}

type Metrics struct {
	Platform   string   `json:"platform,omitempty"`
	ObservedAt string   `json:"observedAt"`
	Views      *int64   `json:"views,omitempty"`
	Reactions  *int64   `json:"reactions,omitempty"`
	Comments   *int64   `json:"comments,omitempty"`
	Shares     *int64   `json:"shares,omitempty"`
	Revenue    *float64 `json:"revenue,omitempty"`
	Currency   string   `json:"currency,omitempty"`
	Method     string   `json:"method"`
}

type Activity struct {
	At     string `json:"at"`
	Action string `json:"action"`
	Part   int    `json:"part"`
}

type CreateInput struct {
	Platforms []string `json:"platforms,omitempty"`
	ProjectID string   `json:"projectId"`
	Seconds   int      `json:"seconds" jsonschema:"Editorial duration target, 15 to 180 seconds; actual narration duration must be checked"`
}

type ActionInput struct {
	PreviousURL string       `json:"previousUrl,omitempty"`
	Destination *Destination `json:"destination,omitempty"`
	Platform    string       `json:"platform,omitempty"`
	ID          string       `json:"id"`
	Revision    int          `json:"revision"`
	Part        int          `json:"part"`
	Action      string       `json:"action" jsonschema:"begin, complete, interrupt, retake, destination, publish, replace, or metrics"`
	TakeID      string       `json:"takeId,omitempty"`
	File        string       `json:"file,omitempty"`
	Duration    float64      `json:"duration,omitempty"`
	URL         string       `json:"url,omitempty"`
	Metrics     *Metrics     `json:"metrics,omitempty"`
}
