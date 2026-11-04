package campaign

import (
	"fmt"
	"net/url"
	"strings"
)

type Destination struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	RecordedAt string `json:"recordedAt,omitempty"`
}

func saveDestination(value *Campaign, input ActionInput, now string) error {
	platform := Platform(input.Platform)
	if input.Destination == nil {
		return fmt.Errorf("a destination account is required")
	}

	account := *input.Destination
	parsed, err := url.Parse(account.URL)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || account.ID == "" || strings.TrimSpace(account.Name) == "" {
		return fmt.Errorf("destination needs its observed account identifier, name and HTTPS URL")
	}

	selected := false

	for _, target := range Targets(*value) {
		selected = selected || target == platform
	}

	facebook := platform == "facebook" && parsed.Host == "www.facebook.com" && parsed.Path == "/profile.php" && parsed.Query().Get("id") == account.ID
	youtube := platform == "youtube" && parsed.Host == "www.youtube.com" && parsed.Path == "/channel/"+account.ID
	if !selected || (!facebook && !youtube) {
		return fmt.Errorf("destination does not match the selected platform account")
	}

	if existing, present := value.Accounts[platform]; present && existing.ID != account.ID {
		for _, part := range value.Parts {
			if ReceiptFor(part, platform) != nil {
				return fmt.Errorf("a destination with posted episodes cannot be replaced")
			}
		}
	}

	if value.Accounts == nil {
		value.Accounts = map[string]Destination{}
	}

	account.RecordedAt = now
	value.Accounts[platform] = account

	return nil
}
