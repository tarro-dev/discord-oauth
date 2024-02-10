package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/tarro-dev/discord-oauth/pkg/option"
)

// Config is the configuration.
type Config struct {
	v *viper.Viper
	l *zerolog.Logger
}

// NewConfig returns a new config.
func NewConfig(file option.Option[string]) *Config {
	v := viper.New()
	v.AutomaticEnv()

	if f, ok := file.Unwrap(); ok {
		v.SetConfigFile(f)

		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	return &Config{v: v}
}

const (
	addrKey     = "addr"
	addrDefault = ":8080"
)

// Addr returns the address to listen on.
func (c *Config) Addr() string {
	c.v.SetDefault(addrKey, addrDefault)
	return c.v.GetString(addrKey)
}
