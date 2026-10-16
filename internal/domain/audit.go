package domain

type ResearchRecord struct {
	ID         string     `json:"id"`
	Subject    string     `json:"subject"`
	Status     string     `json:"status"`
	CreatedAt  string     `json:"createdAt"`
	Model      string     `json:"model"`
	Config     Settings   `json:"config"`
	Requests   int        `json:"requests"`
	Plan       []Question `json:"plan"`
	Activities []Activity `json:"activities"`
}
