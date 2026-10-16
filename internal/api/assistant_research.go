package api

import (
	"context"
	"encoding/json"

	"chronicle/internal/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type researchInput struct {
	Subject string `json:"subject"`
}

type runControlInput struct {
	ID     string `json:"id"`
	Action string `json:"action" jsonschema:"pause, cancel, or resume"`
}

func (s *Server) addResearchTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "chronicle_start_research", Description: "Start autonomous backend research using the configured LLM and external search/retrieval MCP tools. For assistant-led research, use capture_source and save_project instead; those do not need model settings."}, func(ctx context.Context, _ *mcp.CallToolRequest, input researchInput) (*mcp.CallToolResult, domain.ResearchRun, error) {
		settings, err := s.loadSettings(ctx)
		if err != nil {
			return nil, domain.ResearchRun{}, err
		}
		run, err := s.research.Start(ctx, input.Subject, settings)
		return nil, run, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "chronicle_get_research", Description: "Inspect a backend research session's progress, sources, questions and errors.", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true}}, func(ctx context.Context, _ *mcp.CallToolRequest, input projectKey) (*mcp.CallToolResult, domain.ResearchRun, error) {
		data, err := s.store.Get(ctx, "research", input.ID)
		var run domain.ResearchRun
		if err == nil {
			err = json.Unmarshal(data, &run)
		}
		run.Config.APIKey = ""
		return nil, run, err
	})
	mcp.AddTool(server, &mcp.Tool{Name: "chronicle_control_research", Description: "Pause, cancel or resume a backend research session. Resumption preserves source and tool checkpoints."}, func(ctx context.Context, _ *mcp.CallToolRequest, input runControlInput) (*mcp.CallToolResult, domain.ResearchRun, error) {
		settings, err := s.loadSettings(ctx)
		if err != nil {
			return nil, domain.ResearchRun{}, err
		}
		run, err := s.research.Control(ctx, input.ID, input.Action, settings)
		run.Config.APIKey = ""
		return nil, run, err
	})
}
