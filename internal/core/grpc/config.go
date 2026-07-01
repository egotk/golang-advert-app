package coregrpc

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr                string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout     time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	ShouldUseReflection bool          `envconfig:"REFLECTION" default:"false"`
}

func NewConfigMust() Config {
	var config Config

	if err := envconfig.Process("GRPC", &config); err != nil {
		err := fmt.Errorf("get GRPC server config: %w", err)
		panic(err)
	}

	return config
}
