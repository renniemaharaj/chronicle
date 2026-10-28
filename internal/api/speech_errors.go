package api

import (
	"errors"
	"net/http"

	"chronicle/internal/speech"
)

func failSpeech(w http.ResponseWriter, err error) {
	var provider *speech.ProviderError
	if errors.As(err, &provider) {
		status := http.StatusTooManyRequests
		if provider.Code == "insufficient_credits" {
			status = http.StatusPaymentRequired
		}

		write(w, status, map[string]string{
			"error": provider.Message,
			"code":  provider.Code,
		})

		return
	}

	fail(w, err, http.StatusBadGateway)
}
