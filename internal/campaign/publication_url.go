package campaign

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var reelPath = regexp.MustCompile(`^/(?:reel/[0-9]+/?|(?:[^/]+/)?videos/[0-9]+/?)$`)
var youtubeID = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

func publicationURL(address, platform string) (string, string, error) {
	parsed, err := url.Parse(address)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Port() != "" {
		return "", "", fmt.Errorf("provide a direct HTTPS publication URL")
	}

	host := strings.ToLower(parsed.Hostname())
	path := strings.TrimRight(parsed.Path, "/")
	id := path[strings.LastIndex(path, "/")+1:]
	if platform == "youtube" {
		valid := false
		if host == "youtu.be" {
			valid = strings.Count(path, "/") == 1
		} else if host == "youtube.com" || host == "www.youtube.com" || host == "m.youtube.com" {
			valid = strings.HasPrefix(path, "/shorts/") && strings.Count(path, "/") == 2
			if path == "/watch" {
				id, valid = parsed.Query().Get("v"), true
			}
		}

		if !valid || !youtubeID.MatchString(id) {
			return "", "", fmt.Errorf("provide the direct YouTube Short or watch URL")
		}

		return "https://www.youtube.com/shorts/" + id, id, nil
	}

	if platform != "facebook" || (host != "facebook.com" && host != "www.facebook.com" && host != "m.facebook.com") || !reelPath.MatchString(path) {
		return "", "", fmt.Errorf("provide the direct numeric Facebook reel or video URL")
	}

	return "https://www.facebook.com" + path, id, nil
}
