package campaign

import "fmt"

func Platform(name string) string {
	if name == "" {
		return "facebook"
	}

	return name
}

func Targets(value Campaign) []string {
	if len(value.Platforms) == 0 {
		return []string{"facebook"}
	}

	return value.Platforms
}

func ValidatePlatforms(platforms []string) error {
	seen := map[string]bool{}

	for _, platform := range platforms {
		if (platform != "facebook" && platform != "youtube") || seen[platform] {
			return fmt.Errorf("choose Facebook Reels and/or YouTube Shorts without duplicate targets")
		}

		seen[platform] = true
	}

	return nil
}

func ReceiptFor(part Part, platform string) *Receipt {
	platform = Platform(platform)
	if receipt, exists := part.Publications[platform]; exists {
		return &receipt
	}

	if platform == "facebook" {
		return part.Publication
	}

	return nil
}

func anyPublication(part Part) bool {
	return part.Publication != nil || len(part.Publications) > 0
}
