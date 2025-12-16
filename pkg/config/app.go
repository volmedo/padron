package config

import (
	"fmt"

	"github.com/volmedo/padron/pkg/config/app"
)

type AppConfig struct {
	Identity IdentityConfig `mapstructure:"identity" toml:"identity"`
	Server   ServerConfig   `mapstructure:"server" toml:"server"`
}

func (f AppConfig) Validate() error {
	return validateConfig(f)
}

// Normalize applies compatibility fixes before validation.
func (f *AppConfig) Normalize() {}

func (f AppConfig) ToAppConfig() (app.AppConfig, error) {
	var (
		err error
		out app.AppConfig
	)

	out.Identity, err = f.Identity.ToAppConfig()
	if err != nil {
		return app.AppConfig{}, fmt.Errorf("converting identity to app config: %s", err)
	}

	out.Server, err = f.Server.ToAppConfig()
	if err != nil {
		return app.AppConfig{}, fmt.Errorf("converting server config to app config: %s", err)
	}

	return out, nil
}
