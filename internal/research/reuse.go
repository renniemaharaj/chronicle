package research

import (
	"encoding/json"

	"chronicle/internal/domain"
)

func (j *job) reuseSubject(project *domain.Project, name, kind string) error {
	stored, err := j.store.List(j.ctx, "projects")
	if err != nil {
		return err
	}
	for _, data := range stored {
		var previous domain.Project
		if err := json.Unmarshal(data, &previous); err != nil {
			return err
		}
		for _, entity := range previous.Entities {
			if normalizeName(entity.Name) != normalizeName(name) || entity.Kind != kind {
				continue
			}
			if previous.SubjectID == entity.ID && !previous.Fixture {
				*project = previous
				return nil
			}
			project.SubjectID = entity.ID
			project.Entities[0].ID = entity.ID
			return nil
		}
	}
	return nil
}
