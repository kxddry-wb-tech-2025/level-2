package config

import "time"

// Options denote the options parsed by Cobra
type Options struct {
	Timeout time.Duration
	Address string
}
