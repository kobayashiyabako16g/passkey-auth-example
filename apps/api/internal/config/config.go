package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/kvstore"
)

type Config struct {
	Port                 string `env:"PORT" envDefault:"8080"`
	AllowDomain          string `env:"ALLOW_DOMAIN" envDefault:"localhost"`
	AllowOrigin          string `env:"ALLOW_ORIGIN" envDefault:"http://localhost:5173"`
	kvstore.ValKeyConfig `envPrefix:"KV_"`
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
