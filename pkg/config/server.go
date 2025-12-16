package config

import (
	"github.com/volmedo/padron/pkg/config/app"
)

type ServerConfig struct {
	Port uint   `mapstructure:"port" validate:"required,min=1,max=65535" flag:"port" toml:"port"`
	Host string `mapstructure:"host" validate:"required" flag:"host" toml:"host"`
}

func (s ServerConfig) Validate() error {
	return validateConfig(s)
}

func (s ServerConfig) ToAppConfig() (app.ServerConfig, error) {
	return app.ServerConfig{
		Host: s.Host,
		Port: s.Port,
	}, nil
}
