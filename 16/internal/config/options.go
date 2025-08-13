package config

import "time"

// Config is used for a scraper, initialized on flags
type Config struct {
	StartURL     string
	MaxDepth     int
	OutputDir    string
	Workers      int
	Timeout      time.Duration
	UserAgent    string
	IgnoreRobots bool
}
