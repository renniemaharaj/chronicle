package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func assistantCaller(t *testing.T, ctx context.Context, session *mcp.ClientSession) func(string, any, any) {
	t.Helper()

	return func(name string, args any, output any) {
		t.Helper()

		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		})
		if err != nil || result.IsError {
			t.Fatalf("%s: %v, %+v", name, err, result)
		}

		data, err := json.Marshal(result.StructuredContent)
		if err != nil {
			t.Fatal(err)
		}

		if err := json.Unmarshal(data, output); err != nil {
			t.Fatal(err)
		}
	}
}
