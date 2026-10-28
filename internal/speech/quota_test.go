package speech

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProviderLimitCodes(t *testing.T) {
	for _, test := range []struct {
		status int
		body   string
		code   string
	}{
		{401, `{"detail":{"status":"quota_exceeded"}}`, "insufficient_credits"},
		{401, `{"detail":{"code":"insufficient_credits","status":"obsolete"}}`, "insufficient_credits"},
		{402, `{"detail":"private data"}`, "insufficient_credits"},
		{429, `{"detail":{"status":"too_many_concurrent_requests"}}`, "rate_limited"},
		{429, `{"detail":{"code":"concurrent_limit_exceeded"}}`, "rate_limited"},
		{429, `invalid`, "rate_limited"},
	} {
		response := httptest.NewRecorder()
		response.WriteHeader(test.status)
		response.WriteString(test.body)
		var provider *ProviderError
		err := elevenLabsError(response.Result())
		if !errors.As(err, &provider) || provider.Code != test.code || strings.Contains(err.Error(), "private data") {
			t.Fatalf("limit classification failed: %v", err)
		}
	}

}
