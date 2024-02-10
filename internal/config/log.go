package config

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

const (
	logLevelKey     = "log_level"
	logLevelDefault = "info"

	logDevelopmentKey = "log_development"
)

// LogLevel returns the log level.
// If the log level is not set, it returns the default log level.
// If the log level is invalid, it panics.
func (c *Config) LogLevel() zerolog.Level {
	c.v.SetDefault(logLevelKey, logLevelDefault)
	l := c.v.GetString(logLevelKey)

	level, err := zerolog.ParseLevel(l)
	if err != nil {
		panic(err)
	}

	return level
}

// SetLogLevel sets the log level.
func (c *Config) SetLogLevel(l zerolog.Level) {
	c.v.Set(logLevelKey, l.String())
}

// LoggerDevelopment returns true if the logger is in development mode.
func (c *Config) LoggerDevelopment() bool {
	return c.v.GetBool(logDevelopmentKey)
}

// SetLoggerDevelopment sets the logger development mode.
func (c *Config) SetLoggerDevelopment(v bool) {
	c.v.Set(logDevelopmentKey, v)
}

// Logger returns a new logger.
func (c *Config) Logger() *zerolog.Logger {
	if c.l != nil {
		return c.l
	}

	logger := zerolog.New(os.Stdout).Level(c.LogLevel())
	if c.LoggerDevelopment() {
		buildInfo, _ := debug.ReadBuildInfo()

		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Int("pid", os.Getpid()).
			Str("go_version", buildInfo.GoVersion).
			Logger()

	}

	c.l = &logger
	return c.l
}
