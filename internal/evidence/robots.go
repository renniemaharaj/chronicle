package evidence

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
)

func checkRobots(ctx context.Context, client *http.Client, target *url.URL) error {
	robotsURL := target.Scheme + "://" + target.Host + "/robots.txt"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Chronicle/1.0 (local academic research)")
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("cannot check source access rules: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 404 || response.StatusCode == 410 {
		return nil
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("robots.txt returned HTTP %d; source access cannot be established", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, (512<<10)+1))
	if err != nil || len(body) > 512<<10 {
		return fmt.Errorf("robots.txt could not be read within the access-rule budget")
	}
	rules, err := robotstxt.FromBytes(body)
	if err != nil {
		return fmt.Errorf("source access rules could not be parsed: %w", err)
	}
	group := rules.FindGroup("Chronicle")
	if !group.Test(target.RequestURI()) {
		return fmt.Errorf("source access is disallowed for Chronicle by robots.txt")
	}
	if group.CrawlDelay > 20*time.Second {
		return fmt.Errorf("source crawl delay exceeds capture budget; use an appropriate archive retrieval adapter")
	}
	if group.CrawlDelay > 0 {
		timer := time.NewTimer(group.CrawlDelay)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	return nil
}
