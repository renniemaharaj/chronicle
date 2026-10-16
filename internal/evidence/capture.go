package evidence

import (
	"context"
	"crypto/sha256"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"chronicle/internal/domain"
)

type Capture struct {
	Source domain.Source `json:"source"`
	Text   string        `json:"text"`
}

var hiddenHTML = regexp.MustCompile(`(?is)<(?:script|style|noscript)\b[^>]*>.*?</(?:script|style|noscript)>|<!--.*?-->`)
var htmlTags = regexp.MustCompile(`<[^>]+>`)

func NormalizePassage(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// CaptureDocument records the retrieved bytes' hash and a conservative text extraction.
// It never executes scripts or treats document content as tool instructions.
func CaptureDocument(ctx context.Context, client *http.Client, source domain.Source) (Capture, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)
	if err != nil {
		return Capture{}, err
	}
	if err := checkRobots(ctx, client, request.URL); err != nil {
		return Capture{}, err
	}
	request.Header.Set("User-Agent", "Chronicle/1.0 (local academic research)")
	retrieval := *client
	retrieval.CheckRedirect = func(next *http.Request, previous []*http.Request) error {
		if len(previous) >= 5 {
			return fmt.Errorf("source exceeded five redirects")
		}
		return checkRobots(ctx, client, next.URL)
	}
	response, err := retrieval.Do(request)
	if err != nil {
		return Capture{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Capture{}, fmt.Errorf("source returned HTTP %d; no evidence captured", response.StatusCode)
	}
	contentType := strings.ToLower(response.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "text/html") && !strings.Contains(contentType, "text/plain") && !strings.Contains(contentType, "application/xhtml+xml") {
		return Capture{}, fmt.Errorf("source is not HTML or plain text; use a configured document retrieval tool")
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, (2<<20)+1))
	if err != nil {
		return Capture{}, err
	}
	if len(body) > 2<<20 {
		return Capture{}, fmt.Errorf("source exceeds the 2 MiB capture limit")
	}
	source.Resources = ExtractResources(string(body), response.Request.URL.String())
	hash := fmt.Sprintf("%x", sha256.Sum256(body))
	text := string(body)
	if !strings.Contains(contentType, "text/plain") {
		text = html.UnescapeString(htmlTags.ReplaceAllString(hiddenHTML.ReplaceAllString(text, " "), " "))
	}
	text = NormalizePassage(text)
	if text == "" {
		return Capture{}, fmt.Errorf("no readable source passage was retrieved")
	}
	source.ID = fmt.Sprintf("source-%x", sha256.Sum256([]byte(source.URL+hash)))[:31]
	source.ContentHash = hash
	source.RetrievedAt = time.Now().UTC().Format(time.RFC3339)
	source.Verification = "retrieved_text"
	source.IndependenceReviewed = false
	source.Limitations = []string{"Basic text extraction does not preserve page layout. Bibliographic metadata requires researcher review. Retrieval does not establish historical truth."}
	return Capture{Source: source, Text: text}, nil
}
