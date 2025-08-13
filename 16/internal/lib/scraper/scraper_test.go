package scraper

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"testing"
	"wget/internal/config"
	"wget/internal/models"
)

type mockDownloader struct{}

func (m mockDownloader) Download(ctx context.Context, u *url.URL) (io.ReadCloser, string, error) {
	if strings.Contains(u.Path, "error") {
		return nil, "", errors.New("fail")
	}
	return io.NopCloser(strings.NewReader("<html></html>")), "text/html", nil
}

type mockParser struct{}

func (m mockParser) HTML(r io.Reader, u *url.URL) ([]*url.URL, []models.Resource, error) {
	next, _ := url.Parse("https://example.com/next")
	return []*url.URL{next}, nil, nil
}
func (m mockParser) MakeLocal(r io.ReadCloser, base *url.URL, ct, out string) io.ReadCloser {
	return r
}

func TestShouldFollowLink(t *testing.T) {
	cfg := &config.Config{StartURL: "https://example.com/"}
	s := New(cfg)
	u, _ := url.Parse("https://example.com/path")
	if !s.shouldFollowLink(u) {
		t.Errorf("should follow same domain")
	}
}

func TestShouldFollowLinkExternal(t *testing.T) {
	cfg := &config.Config{StartURL: "https://example.com"}
	s := New(cfg)
	u, _ := url.Parse("https://external.com/path")
	if s.shouldFollowLink(u) {
		t.Errorf("should not follow external domain")
	}
}

func TestShouldDownloadResource(t *testing.T) {
	cfg := &config.Config{StartURL: "https://example.com"}
	s := New(cfg)
	u1, _ := url.Parse("https://example.com/resource.png")
	u2, _ := url.Parse("https://other.com/resource.png")
	if !s.shouldDownloadResource(u1) {
		t.Errorf("expected same domain resource to be downloaded")
	}
	if s.shouldDownloadResource(u2) {
		t.Errorf("expected different domain resource not to be downloaded")
	}
}

func TestHasAnySuffix(t *testing.T) {
	if !hasAnySuffix("file.CSS", ".css") {
		t.Errorf("suffix match should be case-insensitive")
	}
	if hasAnySuffix("file.txt", ".css", ".js") {
		t.Errorf("should not match non-listed suffix")
	}
}

func TestCanonRemovesFragment(t *testing.T) {
	cfg := &config.Config{StartURL: "https://example.com"}
	s := New(cfg)
	u, _ := url.Parse("https://example.com/page#section")
	got := s.canon(u)
	if strings.Contains(got, "#") {
		t.Errorf("expected fragment removed, got %q", got)
	}
}

func TestAddTaskAndAvoidDuplicate(t *testing.T) {
	cfg := &config.Config{StartURL: "https://example.com"}
	s := New(cfg)
	u, _ := url.Parse("https://example.com/page")
	task := &models.Task{URL: u, Depth: 0, Type: models.HTML}
	s.addTask(task)
	if len(s.workQueue) != 1 {
		t.Errorf("expected 1 task in queue")
	}
	// try to add duplicate
	s.addTask(task)
	if len(s.workQueue) != 1 {
		t.Errorf("duplicate task should not be enqueued")
	}
}
