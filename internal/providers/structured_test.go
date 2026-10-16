package providers

import (
	"context"
	"encoding/json"
	"testing"
)

type fixedOutput string

func (output fixedOutput) Complete(context.Context, string, string) (json.RawMessage, error) {
	return json.RawMessage(output), nil
}

func TestStructuredOutputRejectsWrongTypesAndUnknownFields(t *testing.T) {
	type result struct {
		Name string `json:"name"`
	}
	for _, raw := range []string{`{"name":42}`, `{"name":"Synthetic","unauthorizedAction":"delete"}`} {
		if _, err := Structured[result](context.Background(), fixedOutput(raw), "extract", "untrusted document"); err == nil {
			t.Fatalf("invalid structured output was accepted: %s", raw)
		}
	}
	valid, err := Structured[result](context.Background(), fixedOutput(`{"name":"Synthetic"}`), "extract", "document")
	if err != nil || valid.Name != "Synthetic" {
		t.Fatal("valid schema-conforming output should be accepted")
	}
}
