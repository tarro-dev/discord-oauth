package config

import (
	"context"
)

type configKeyType string

const configKey configKeyType = "config"

// Ctx returns the config from the context.
func Ctx(ctx context.Context) *Config {
	val := ctx.Value(configKey)
	if val == nil {
		panic("config not found in context") // TODO: do this better
	}
	cfg, ok := ctx.Value(configKey).(*Config)
	if !ok {
		panic("config not found in context") // TODO: do this better
	}

	return cfg
}

// WithContext returns a new context with the config.
func (c *Config) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, configKey, c)
}
