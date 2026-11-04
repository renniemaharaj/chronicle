package api

import (
	"context"

	"chronicle/internal/campaign"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) addAudioCampaignTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "chronicle_create_audio_campaign",
		Description: "Create episodes only from the consecutive saved chronological narration available for a selected voice, without generating speech. Stops at the first missing audio file and measures existing files with ffprobe. Platforms may contain facebook and youtube. FFmpeg export uses scripts/render-reels.mjs to record the real story map with Follow location enabled and reuse the saved narration.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input campaign.AudioInput) (*mcp.CallToolResult, campaign.Campaign, error) {
		value, err := s.createAudioCampaign(ctx, input)

		return nil, value, err
	})
}
