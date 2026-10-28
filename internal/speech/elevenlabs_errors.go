package speech

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ProviderError struct {
	Code    string
	Message string
}

func (e *ProviderError) Error() string {
	return e.Message
}

func elevenLabsError(response *http.Response) error {
	var output struct {
		Detail struct {
			Status string `json:"status"`
			Code   string `json:"code"`
		} `json:"detail"`
	}
	_ = json.NewDecoder(io.LimitReader(response.Body, 16<<10)).Decode(&output)

	// Translate known codes only; provider messages may include credentials or private data.
	code := output.Detail.Code
	if code == "" {
		code = output.Detail.Status
	}

	switch code {
	case "api_key_id_used_as_api_key":
		return fmt.Errorf("ElevenLabs needs the secret API key beginning with sk_; the configured value is an API key ID")
	case "invalid_api_key":
		return fmt.Errorf("ElevenLabs rejected the API key; replace it with a valid secret key")
	case "missing_permissions", "insufficient_permissions":
		return fmt.Errorf("the ElevenLabs key lacks permission for this request; check its API key permissions")
	case "quota_exceeded", "insufficient_credits", "credit_balance_too_low", "character_limit_exceeded":
		return &ProviderError{
			Code:    "insufficient_credits",
			Message: "ElevenLabs credits or the API key spending limit cannot cover this passage. Add credits, raise the key limit, or wait for renewal, then resume missing audio.",
		}
	case "rate_limit_exceeded", "concurrent_limit_exceeded", "too_many_concurrent_requests", "system_busy":
		return &ProviderError{
			Code:    "rate_limited",
			Message: "ElevenLabs is temporarily rate limited. Wait before resuming missing audio.",
		}
	}

	if response.StatusCode == http.StatusPaymentRequired {
		return &ProviderError{
			Code:    "insufficient_credits",
			Message: "ElevenLabs requires credits to continue. Add credits or wait for renewal, then resume missing audio.",
		}
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return &ProviderError{
			Code:    "rate_limited",
			Message: "ElevenLabs is temporarily rate limited. Wait before resuming missing audio.",
		}
	}

	return fmt.Errorf("ElevenLabs returned HTTP %d; check key permissions, credits, model and voice", response.StatusCode)
}
