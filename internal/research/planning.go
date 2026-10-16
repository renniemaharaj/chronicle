package research

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"chronicle/internal/domain"
	"chronicle/internal/providers"
)

type planOutput struct {
	Name      string            `json:"name"`
	Kind      string            `json:"kind"`
	Questions []domain.Question `json:"questions"`
}

func stableID(kind, text string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(text))))
	return kind + "-" + hex.EncodeToString(hash[:8])
}

func (j *job) plan() error {
	if len(j.snapshot().Plan) > 0 {
		return nil
	}
	input := j.snapshot().Subject
	plan, err := providers.Structured[planOutput](j.ctx, j, `Identify and provisionally disambiguate the subject. Classify with an extensible entity kind.
Create four chronological research questions covering origins, development, relationships/context, and later state/legacy.
Adapt to person, family, species, place, institution, invention, event, or idea. Investigate evidence disagreements and gaps.
Use interval descriptions, not invented dates. All question statuses must be pending.`, input)
	if err != nil {
		return err
	}
	if plan.Name == "" || plan.Kind == "" || len(plan.Questions) < 3 || len(plan.Questions) > 6 {
		return fmt.Errorf("research plan requires a named subject, type, and 3–6 questions")
	}
	id := stableID("entity", plan.Kind+":"+plan.Name)
	project := domain.EmptyProject(stableID("project", input), plan.Name)
	project.SubjectID = id
	project.Entities = append(project.Entities, domain.Entity{ID: id, Name: plan.Name, Kind: plan.Kind, Summary: "Provisional subject identification. Inspect claims and evidence for documented information.", Version: 1, Attributes: map[string]any{}})
	if err := j.reuseSubject(&project, plan.Name, plan.Kind); err != nil {
		return err
	}
	id = project.SubjectID
	for i := range plan.Questions {
		plan.Questions[i].Status = "pending"
	}
	j.update(func(run *domain.ResearchRun) { run.SubjectID = id; run.Plan = plan.Questions; run.Project = &project })
	return nil
}

func inputJSON(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}
