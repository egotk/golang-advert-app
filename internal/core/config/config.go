package coreconfig

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
	Root     string
}

func New() (*Config, error) {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("load time zone: %s: %w", tz, err)
	}

	root := os.Getenv("PROJECT_ROOT")

	return &Config{
		TimeZone: zone,
		Root:     root,
	}, nil
}

func NewMust() *Config {
	config, err := New()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
