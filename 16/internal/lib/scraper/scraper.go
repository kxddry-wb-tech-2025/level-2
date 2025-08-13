package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"wget/internal/config"
	"wget/internal/lib/downloader"
	"wget/internal/lib/parser"
	"wget/internal/lib/set"
	"wget/internal/models"

	"github.com/temoto/robotstxt"
)

// Parser is an interface that dictates what functions must a parser have
type Parser interface {
	HTML(html io.Reader, URL *url.URL) ([]*url.URL, []models.Resource, error)
	MakeLocal(html io.ReadCloser, base *url.URL, contentType, outputDir string) (newHTML io.ReadCloser)
}

// Downloader is an interface that dictates what functions must a downloader have
type Downloader interface {
	Download(ctx context.Context, url *url.URL) (content io.ReadCloser, contentType string, err error)
}

// Scraper is the main entity of the application
type Scraper struct {
	config      *config.Config
	downloader  Downloader
	parser      Parser
	visited     *set.Set
	enqueued    *set.Set
	baseURL     *url.URL
	baseDomain  string
	scopePrefix string
	workQueue   chan *models.Task
	wg          *sync.WaitGroup
	robots      *robotstxt.RobotsData
}

// New creates and initializes a Scraper
func New(config *config.Config) *Scraper {

	if len(config.StartURL) < 7 || config.StartURL[:4] != "http" {
		config.StartURL = "https://" + config.StartURL
	}

	baseURL, err := url.Parse(config.StartURL)
	if err != nil {
		println("invalid start url", config.StartURL)
		os.Exit(1)
		return nil
	}

	p := baseURL.Path
	if p == "" {
		p = "/"
	}
	currDir, _ := os.Getwd()
	if config.OutputDir == "" {
		config.OutputDir = baseURL.Host
	}

	if !filepath.IsAbs(config.OutputDir) {
		config.OutputDir = filepath.Join(currDir, config.OutputDir)
	}

	return &Scraper{
		config:      config,
		downloader:  downloader.New(config.Timeout, config.UserAgent),
		parser:      parser.New(),
		visited:     set.New(),
		enqueued:    set.New(),
		baseURL:     baseURL,
		baseDomain:  baseURL.Hostname(),
		scopePrefix: p,
		workQueue:   make(chan *models.Task, 1000),
		wg:          new(sync.WaitGroup),
	}
}

// Start starts a scraper
func (s *Scraper) Start(ctx context.Context) error {
	if err := os.MkdirAll(s.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if !s.config.IgnoreRobots {
		u := *s.baseURL
		u.Path = "/robots.txt"

		bots, contentType, err := s.downloader.Download(ctx, &u)
		if err == nil && contentType != "" {
			defer bots.Close()
			byteBots, err := io.ReadAll(bots)
			if err == nil {
				data, err := robotstxt.FromBytes(byteBots)
				if err == nil {
					s.robots = data
				}
			}
		}
	}

	for i := 0; i < s.config.Workers; i++ {
		go s.worker(ctx)
	}

	parsed, err := url.Parse(s.config.StartURL)
	if err != nil {
		return fmt.Errorf("failed to parse start url: %w", err)
	}

	s.addTask(&models.Task{
		URL:   parsed,
		Depth: 0,
		Type:  models.HTML,
	})

	s.wg.Wait()
	close(s.workQueue)
	return nil
}

func (s *Scraper) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-s.workQueue:
			if !ok {
				return
			}

			func(t *models.Task) {
				defer s.wg.Done()

				if s.robots != nil {
					g := s.robots.FindGroup(s.config.UserAgent)
					if !g.Test(t.URL.Path) {
						fmt.Println("blocked by robots.txt:", task.URL)
						return
					}
				}

				if task.Depth > s.config.MaxDepth {
					return
				}

				key := s.canon(task.URL)
				if s.visited.Contains(key) {
					return
				}
				s.visited.Add(key)
				s.enqueued.Remove(key)

				if err := s.processTask(ctx, task); err != nil {
					println("error processing task", task.URL, err.Error())
					return
				}
			}(task)
		}
	}
}

func (s *Scraper) processTask(ctx context.Context, task *models.Task) error {
	fmt.Printf("downloading %s (depth %d)\n", task.URL, task.Depth)

	content, contentType, err := s.downloader.Download(ctx, task.URL)
	if err != nil {
		return err
	}

	if task.Type == models.HTML && contentType == "text/html" {
		data, err := io.ReadAll(content)
		if err != nil {
			return err
		}
		links, resources, err := s.parser.HTML(bytes.NewReader(data), task.URL)
		_ = content.Close()
		if err != nil {
			return err
		}

		for _, link := range links {
			if s.shouldFollowLink(link) {
				s.addTask(&models.Task{
					URL:   link,
					Depth: task.Depth + 1,
					Type:  models.HTML,
				})
			}
		}

		for _, resource := range resources {
			if s.shouldDownloadResource(resource.URL) {
				s.addTask(&models.Task{
					URL:   resource.URL,
					Depth: task.Depth,
					Type:  resource.Type,
				})
			}
		}

		upd := s.parser.MakeLocal(io.NopCloser(bytes.NewReader(data)), task.URL, contentType, s.config.OutputDir)
		if _, err := saveFile(task.URL, upd, contentType, s.config.OutputDir); err != nil {
			return err
		}
	} else {
		// Stream non-HTML content directly to disk; saveFile will close the body.
		if _, err := saveFile(task.URL, content, contentType, s.config.OutputDir); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scraper) shouldFollowLink(parsedURL *url.URL) bool {
	// 1. Skip non-HTTP(S) schemes
	if parsedURL.Scheme != "" && parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// 2. Normalize host for comparison
	host := strings.ToLower(parsedURL.Host)
	baseHost := strings.ToLower(s.baseDomain)

	// 3. Internal links OR resource files
	sameDomain := host == "" || host == baseHost
	isStaticAsset := hasAnySuffix(parsedURL.Path, ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp")

	// 4. wget -m generally grabs assets, but only follows HTML links in-domain
	return sameDomain || isStaticAsset
}

func hasAnySuffix(s string, suffixes ...string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(strings.ToLower(s), suf) {
			return true
		}
	}
	return false
}

func (s *Scraper) shouldDownloadResource(parsedURL *url.URL) bool {
	return parsedURL.Host == s.baseDomain || parsedURL.Host == ""
}

func (s *Scraper) addTask(task *models.Task) {
	key := s.canon(task.URL)

	// Avoid queue blow-ups: don't enqueue if we've seen or already enqueued it.
	if s.visited.Contains(key) || s.enqueued.Contains(key) {
		return
	}

	s.enqueued.Add(key)
	s.wg.Add(1)
	select {
	case s.workQueue <- task:
		// enqueued
	default:
		// queue full, drop and undo bookkeeping to avoid leaks
		s.enqueued.Remove(key)
		s.wg.Done()
	}
}

func (s *Scraper) hasActiveWorkers() bool {
	return len(s.workQueue) > 0
}

// canon represents the canonical string version of a link
func (s *Scraper) canon(u *url.URL) string {
	v := *u
	v.Fragment = ""
	return v.String()
}
