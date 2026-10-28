package research

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chronicle/internal/domain"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type modelFunc func(context.Context, string, string) (json.RawMessage, error)

func (f modelFunc) Complete(ctx context.Context, system, input string) (json.RawMessage, error) {
	return f(ctx, system, input)
}

const testDocument = "Meridian Archive. The synthetic observatory was established in 1800. Section: Foundation."
const testURL = "https://archive.example/meridian"

func testMCP(t *testing.T) *httptest.Server {
	t.Helper()

	server := sdk.NewServer(&sdk.Implementation{
		Name:    "synthetic-archive",
		Version: "1",
	}, nil)

	for _, name := range []string{"search", "retrieve"} {
		server.AddTool(&sdk.Tool{
			Name:        name,
			InputSchema: map[string]any{"type": "object"},
		}, func(ctx context.Context, request *sdk.CallToolRequest) (*sdk.CallToolResult, error) {
			text := testURL
			if request.Params.Name == "retrieve" {
				text = testDocument
			}

			return &sdk.CallToolResult{Content: []sdk.Content{&sdk.TextContent{Text: text}}}, nil
		})
	}

	handler := sdk.NewStreamableHTTPHandler(func(*http.Request) *sdk.Server {
		return server
	}, nil)

	return httptest.NewServer(handler)
}

func syntheticModel(ctx context.Context, system, input string) (json.RawMessage, error) {
	var value any

	switch {
	case strings.Contains(system, "Identify and provisionally"):
		value = planOutput{
			Name: "Meridian Observatory",
			Kind: "organization",
			Questions: []domain.Question{{
				Question: "Origins",
				Interval: "earliest records",
				Status:   "pending",
			}, {
				Question: "Development",
				Interval: "early period",
				Status:   "pending",
			}, {
				Question: "Relationships",
				Interval: "middle period",
				Status:   "pending",
			}, {
				Question: "Legacy",
				Interval: "later period",
				Status:   "pending",
			}},
		}
	case strings.Contains(system, "Create JSON arguments"):
		value = argumentsOutput{Arguments: map[string]any{"query": "Meridian Observatory"}}
	case strings.Contains(system, "Select one relevant"):
		value = selectionOutput{
			URL:       testURL,
			Arguments: map[string]any{"url": testURL},
		}
	default:
		var extracted extraction
		json.Unmarshal([]byte(testBatch), &extracted)
		value = extracted
	}

	return json.Marshal(value)
}
