package config

import (
	"fmt"
	"net/url"

	"github.com/volmedo/padron/pkg/config/app"
)

type ServerConfig struct {
	Port      uint   `mapstructure:"port" validate:"required,min=1,max=65535" flag:"port" toml:"port"`
	Host      string `mapstructure:"host" validate:"required" flag:"host" toml:"host"`
	PublicURL string `mapstructure:"public_url" flag:"public-url" toml:"public_url"`
}

func (s ServerConfig) Validate() error {
	return validateConfig(s)
}

func (s ServerConfig) ToAppConfig() (app.ServerConfig, error) {
	publicURL, err := url.Parse(s.PublicURL)
	if err != nil {
		return app.ServerConfig{}, fmt.Errorf("invalid public URL: %w", err)
	}

	return app.ServerConfig{
		Host:      s.Host,
		Port:      s.Port,
		PublicURL: publicURL,
	}, nil
}
