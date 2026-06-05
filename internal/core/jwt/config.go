package corejwt

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AccessSecret   string        `envconfig:"ACCESS_SECRET" required:"true"`
	AccessTTL      time.Duration `envconfig:"ACCESS_TTL" default:"30m"`
	RefreshByteLen int           `envconfig:"REFRESH_BYTE_LEN" default:"32"`
	RefreshTTL     time.Duration `envconfig:"REFRESH_TTL" default:"720h"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get JWT config: %w", err)
		panic(err)
	}

	return config
}
