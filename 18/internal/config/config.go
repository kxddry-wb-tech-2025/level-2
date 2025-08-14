package config

import "time"

// Config is a structure parsed in the config.yaml file
type Config struct {
	Env    string `yaml:"env" env-default:"local"` // local, dev, prod
	Server server `yaml:"server"`
}

type server struct {
	Host    string        `yaml:"host" env-default:"0.0.0.0"`
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}
