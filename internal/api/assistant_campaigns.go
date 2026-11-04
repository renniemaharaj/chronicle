package api

import (
	"context"

	"chronicle/internal/campaign"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type campaignList struct {
	Campaigns []campaign.Campaign `json:"campaigns"`
}

func (s *Server) addCampaignTools(server *mcp.Server) {
	readOnly := &mcp.ToolAnnotations{ReadOnlyHint: true}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "chronicle_list_campaigns",
		Description: "List saved reel series, recording checkpoints, publication receipts and manual performance observations.",
		Annotations: readOnly,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, campaignList, error) {
		values, err := s.campaigns.List(ctx)

		return nil, campaignList{Campaigns: values}, err
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "chronicle_get_campaign",
		Description: "Read a reel series before resuming. Inspect the last completed take and first unpublished part. Offsets are UTF-8 bytes in frozen evidence-qualified passages.",
		Annotations: readOnly,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input projectKey) (*mcp.CallToolResult, campaign.Campaign, error) {
		value, err := s.campaigns.Get(ctx, input.ID)

		return nil, value, err
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "chronicle_create_campaign",
		Description: "Plan a vertical reel series from a saved research project. Freezes the chronological event narration and historical context without omitting passages. Duration is an estimate at 140 words per minute, not measured media length. Does not publish.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input campaign.CreateInput) (*mcp.CallToolResult, campaign.Campaign, error) {
		value, err := s.campaigns.Create(ctx, input)

		return nil, value, err
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "chronicle_update_campaign",
		Description: "Checkpoint a take (begin, complete, interrupt, retake), register a manually verified Facebook publication URL (publish), or append manually observed metrics. Requires current revision; completing requires takeId, file and measured duration. Retake preserves completed files and only resets unpublished parts. Publish records a receipt only: it does not upload or verify Facebook. Never mark a part published before observing the actual public reel. Publish in series order.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input campaign.ActionInput) (*mcp.CallToolResult, campaign.Campaign, error) {
		value, err := s.campaigns.Update(ctx, input)

		return nil, value, err
	})
}
