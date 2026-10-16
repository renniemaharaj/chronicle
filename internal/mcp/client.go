package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client struct {
	session *sdk.ClientSession
	Tools   map[string]*sdk.Tool
	allowed map[string]bool
}

func Connect(ctx context.Context, endpoint string, allowed ...string) (*Client, error) {
	client := sdk.NewClient(&sdk.Implementation{Name: "chronicle", Version: "0.1.0"}, nil)
	timeout, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	session, err := client.Connect(timeout, &sdk.StreamableClientTransport{Endpoint: endpoint, HTTPClient: &http.Client{Timeout: 35 * time.Second}, MaxRetries: 1, DisableStandaloneSSE: true, MaxEventSize: 2 << 20}, nil)
	if err != nil {
		return nil, err
	}
	result := &Client{session: session, Tools: map[string]*sdk.Tool{}, allowed: map[string]bool{}}
	for tool, err := range session.Tools(timeout, nil) {
		if err != nil {
			session.Close()
			return nil, err
		}
		result.Tools[tool.Name] = tool
		if len(result.Tools) > 150 {
			session.Close()
			return nil, fmt.Errorf("MCP server exceeds 150-tool discovery budget")
		}
	}
	for _, name := range allowed {
		result.allowed[name] = true
	}
	return result, nil
}

func (c *Client) Call(ctx context.Context, name string, arguments map[string]any) (string, error) {
	tool, exists := c.Tools[name]
	if !exists || !c.allowed[name] {
		return "", fmt.Errorf("tool %q is not an explicitly configured research capability", name)
	}
	encoded, err := json.Marshal(tool.InputSchema)
	if err != nil {
		return "", err
	}
	var schema jsonschema.Schema
	if err := json.Unmarshal(encoded, &schema); err != nil {
		return "", err
	}
	resolved, err := schema.Resolve(nil)
	if err != nil {
		return "", err
	}
	if err := resolved.Validate(arguments); err != nil {
		return "", fmt.Errorf("tool arguments failed schema validation: %w", err)
	}
	timeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	result, err := c.session.CallTool(timeout, &sdk.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		return "", err
	}
	if result.IsError {
		return "", fmt.Errorf("MCP research tool %s returned an error", name)
	}
	parts := []string{}
	for _, content := range result.Content {
		if text, ok := content.(*sdk.TextContent); ok {
			parts = append(parts, text.Text)
		}
	}
	if result.StructuredContent != nil {
		data, err := json.Marshal(result.StructuredContent)
		if err != nil {
			return "", err
		}
		parts = append(parts, string(data))
	}
	text := strings.Join(parts, "\n")
	if len(text) == 0 || len(text) > 2<<20 {
		return "", fmt.Errorf("tool output is empty or exceeds the 2 MiB document budget")
	}
	return text, nil
}

func (c *Client) Close() error { return c.session.Close() }
