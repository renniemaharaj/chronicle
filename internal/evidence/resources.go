package evidence

import (
	"net/url"
	"regexp"
	"strings"

	"chronicle/internal/domain"
	"golang.org/x/net/html"
)

var youtubeID = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

var markdownResource = regexp.MustCompile(`(!?)\[([^\]\n]{0,300})\]\((https?://[^\s)]+)\)`)

// Extraction preserves candidates without downloading media or following embeds.
func ExtractResources(document, baseAddress string) []domain.Resource {
	base, err := url.Parse(baseAddress)
	if err != nil {
		return nil
	}
	result := []domain.Resource{}
	add := func(kind, address, label string) {
		relative, err := url.Parse(strings.TrimSpace(address))
		if err != nil || address == "" {
			return
		}
		resolved := base.ResolveReference(relative)
		resolved.Fragment = ""
		if !domain.SafeWebURL(resolved.String()) {
			return
		}
		if kind == "link" && videoAddress(resolved.String()) {
			kind = "video"
		}
		result = keepResource(result, domain.Resource{Kind: kind, URL: resolved.String(), Label: NormalizePassage(label)})
	}
	tokens := html.NewTokenizer(strings.NewReader(document))
	var anchor, label string
	for {
		tokenType := tokens.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := tokens.Token()
		attributes := map[string]string{}
		for _, attribute := range token.Attr {
			attributes[attribute.Key] = attribute.Val
		}
		if tokenType == html.TextToken && anchor != "" && len(label) < 300 {
			label += " " + token.Data
		}
		if tokenType == html.EndTagToken && token.Data == "a" {
			add("link", anchor, label)
			anchor, label = "", ""
		}
		if tokenType != html.StartTagToken && tokenType != html.SelfClosingTagToken {
			continue
		}
		// Some documented players expose their YouTube ID instead of an iframe URL.
		id := attributes["data-youtube-id"]
		if id == "" && strings.Contains(document, "youtube.com") {
			id = attributes["data-video"]
		}
		if youtubeID.MatchString(id) {
			add("video", "https://www.youtube.com/watch?v="+id, attributes["data-title"])
		}
		switch token.Data {
		case "a":
			anchor, label = attributes["href"], attributes["aria-label"]
		case "img":
			if attributes["width"] == "1" || attributes["height"] == "1" {
				continue
			}
			address := attributes["src"]
			if address == "" {
				address = attributes["data-src"]
			}
			add("image", address, attributes["alt"])
		case "meta":
			property := attributes["property"] + attributes["name"]
			if property == "og:image" || property == "twitter:image" {
				add("image", attributes["content"], "Article lead image")
			}
			if property == "og:video" || property == "og:video:url" {
				add("video", attributes["content"], "Article video")
			}
		case "video":
			add("video", attributes["src"], attributes["title"])
		case "iframe":
			parsed, _ := url.Parse(attributes["src"])
			if parsed != nil && videoHost(parsed.Hostname()) {
				add("video", attributes["src"], attributes["title"])
			}
		}
	}
	for _, match := range markdownResource.FindAllStringSubmatch(document, 180) {
		kind := "link"
		if match[1] == "!" {
			kind = "image"
		}
		add(kind, match[3], match[2])
	}
	return result
}

func videoHost(host string) bool {
	for _, allowed := range []string{"youtube.com", "youtu.be", "youtube-nocookie.com", "vimeo.com"} {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}
