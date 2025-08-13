package config

import "time"

type Options struct {
	Timeout time.Duration
	Address string
}
