package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func (c Config) Validate() error {
	_, port, err := net.SplitHostPort(c.Address)
	if err != nil {
		return fmt.Errorf("address must contain a host and port")
	}

	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return fmt.Errorf("address must use a port between 1 and 65535")
	}

	if c.Workspace == "" || c.StaticDir == "" {
		return fmt.Errorf("workspace and staticDir must not be empty")
	}

	endpoints := map[string]string{
		"geocoderUrl":          c.GeocoderURL,
		"mapTileUrl":           strings.NewReplacer("{z}", "0", "{x}", "0", "{y}", "0", "{s}", "a").Replace(c.MapTileURL),
		"speech.elevenLabsUrl": c.Speech.ElevenLabsURL,
		"speech.endpoint":      c.Speech.Endpoint,
		"research.baseUrl":     c.Research.BaseURL,
	}
	if c.Research.MCPURL != "" {
		endpoints["research.mcpUrl"] = c.Research.MCPURL
	}

	for name, address := range endpoints {
		if !httpURL(address) {
			return fmt.Errorf("%s must be an HTTP(S) URL without embedded credentials", name)
		}
	}

	for _, origin := range c.AllowedOrigins {
		parsed, err := url.Parse(origin)
		if !httpURL(origin) || err != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("allowedOrigins must contain exact HTTP(S) origins")
		}
	}

	if c.Research.Provider != "ollama" && c.Research.Provider != "openai" {
		return fmt.Errorf("research.provider must be ollama or openai")
	}

	if c.Research.MaxRequests < 4 || c.Research.MaxRequests > 100 || c.Research.Concurrency < 1 || c.Research.Concurrency > 4 {
		return fmt.Errorf("research budget must be 4-100 and concurrency 1-4")
	}

	return nil
}

func httpURL(value string) bool {
	parsed, err := url.Parse(value)

	return err == nil && parsed.Host != "" && parsed.User == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
}
