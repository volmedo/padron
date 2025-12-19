package config

import (
	"fmt"

	"github.com/volmedo/padron/pkg/config/app"
)

type Config struct {
	Identity IdentityConfig `mapstructure:"identity" toml:"identity"`
	Server   ServerConfig   `mapstructure:"server" toml:"server"`
	Stores   StoreConfig    `mapstructure:"stores" toml:"stores"`
}

func (f Config) Validate() error {
	return validateConfig(f)
}

// Normalize applies compatibility fixes before validation.
func (f *Config) Normalize() {}

func (f Config) ToAppConfig() (app.AppConfig, error) {
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

	out.Stores, err = f.Stores.ToAppConfig()
	if err != nil {
		return app.AppConfig{}, fmt.Errorf("converting stores config to app config: %s", err)
	}

	return out, nil
}
