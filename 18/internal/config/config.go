package config

import "time"

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Host        string        `yaml:"host" env-default:"0.0.0.0"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
}
