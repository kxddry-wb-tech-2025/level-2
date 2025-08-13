package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Downloader implements the scraper.Downloader interface.
type Downloader struct {
	client    *http.Client
	userAgent string
}

// New initializes and creates a Downloader
func New(timeout time.Duration, userAgent string) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: userAgent,
	}
}

// Download downloads a file from a link
func (d *Downloader) Download(ctx context.Context, url *url.URL) (io.ReadCloser, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", d.userAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error downloading %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, "", fmt.Errorf("HTTP %s: status %d", url, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		if idx := strings.Index(contentType, ";"); idx != -1 {
			contentType = strings.TrimSpace(contentType[:idx])
		}
	}

	return resp.Body, contentType, nil
}
