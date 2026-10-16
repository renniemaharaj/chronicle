package domain

import "encoding/json"

type Settings struct {
	Provider     string `json:"provider"`
	BaseURL      string `json:"baseUrl"`
	Model        string `json:"model"`
	APIKey       string `json:"apiKey,omitempty"`
	HasAPIKey    bool   `json:"hasApiKey,omitempty"`
	MCPURL       string `json:"mcpUrl"`
	SearchTool   string `json:"searchTool"`
	RetrieveTool string `json:"retrieveTool"`
	MaxRequests  int    `json:"maxRequests"`
	Concurrency  int    `json:"concurrency"`
}

type Question struct {
	Question string `json:"question"`
	Interval string `json:"interval"`
	Status   string `json:"status"`
}

type Activity struct {
	Time    string          `json:"time"`
	Message string          `json:"message"`
	Tool    string          `json:"tool,omitempty"`
	Input   json.RawMessage `json:"input,omitempty"`
	Output  string          `json:"output,omitempty"`
}

type ResearchRun struct {
	ID         string     `json:"id"`
	Subject    string     `json:"subject"`
	SubjectID  string     `json:"subjectId"`
	Status     string     `json:"status"`
	Phase      string     `json:"phase"`
	Question   string     `json:"question"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
	Model      string     `json:"model"`
	Error      string     `json:"error,omitempty"`
	Requests   int        `json:"requests"`
	Sources    int        `json:"sources"`
	Claims     int        `json:"claims"`
	Conflicts  int        `json:"conflicts"`
	Plan       []Question `json:"plan"`
	Activities []Activity `json:"activities"`
	Project    *Project   `json:"project,omitempty"`
	Config     Settings   `json:"config"`
}

func EmptyProject(id, title string) Project {
	return Project{
		FormatVersion: 1, ID: id, Title: title,
		Entities: []Entity{}, Events: []Event{}, Relationships: []Relationship{},
		Claims: []Claim{}, Sources: []Source{}, Evidence: []Evidence{}, Conflicts: []Conflict{},
		Narrative: Narrative{Chapters: []Chapter{}}, Coverage: []Coverage{}, Notes: []Note{}, Layout: map[string]Position{},
	}
}
