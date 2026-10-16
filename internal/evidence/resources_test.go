package evidence

import (
	"fmt"
	"strings"
	"testing"
)

func TestResourceExtractionPreservesProvenanceWithoutFetchingMedia(t *testing.T) {
	document := `<meta property="og:image" content="/lead.jpg"><img src="/portrait.jpg" alt="Portrait of a documented subject"><img src="javascript:alert(1)">
 <a href="https://artist.example/">Official website</a><a href="https://www.youtube.com/channel/abc">YouTube</a>
 <button data-video="abcdefghijk" data-title="A documented performance"></button>
 <iframe src="https://www.youtube.com/embed/abcdefghijk"></iframe><img src="/tracker" width="1">`
	resources := ExtractResources(document, "https://archive.example/article")
	found := map[string]string{}
	for _, resource := range resources {
		found[resource.URL] = resource.Kind
	}
	if found["https://archive.example/portrait.jpg"] != "image" || found["https://artist.example/"] != "link" {
		t.Fatal("lost image or website reference")
	}
	if found["https://www.youtube.com/channel/abc"] != "link" {
		t.Fatal("a channel is not an individual video")
	}
	if found["https://www.youtube.com/watch?v=abcdefghijk"] != "video" {
		t.Fatal("player identifier not preserved")
	}
	if _, exists := found["javascript:alert(1)"]; exists {
		t.Fatal("unsafe resource accepted")
	}
	if _, exists := found["https://archive.example/tracker"]; exists {
		t.Fatal("tracking pixel accepted")
	}
}

func TestLongArticleRetainsOfficialReferenceAtFooter(t *testing.T) {
	var document strings.Builder
	for index := range 500 {
		fmt.Fprintf(&document, `<a href="/navigation/%d">Navigation</a>`, index)
	}
	document.WriteString(`<a href="https://official.example/">Official website</a>`)
	resources := ExtractResources(document.String(), "https://archive.example/")
	if len(resources) > 180 {
		t.Fatal("resource budget exceeded")
	}
	for _, resource := range resources {
		if resource.URL == "https://official.example/" {
			return
		}
	}
	t.Fatal("long article hid its official reference")
}
