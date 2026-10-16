package fixtures

import (
	"embed"
	"encoding/json"

	"chronicle/internal/domain"
)

//go:embed *.json
var files embed.FS

func Projects() ([]domain.Project, error) {
	entries, err := files.ReadDir(".")
	if err != nil {
		return nil, err
	}
	projects := []domain.Project{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := files.ReadFile(entry.Name())
		if err != nil {
			return nil, err
		}
		var project domain.Project
		if err := json.Unmarshal(data, &project); err != nil {
			return nil, err
		}
		if err := domain.ValidateProject(project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
