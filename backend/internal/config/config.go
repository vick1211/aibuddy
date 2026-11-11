package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	LogLevel    slog.Level
}

func Load(getenv func(string) string) (Config, error) {
	c := Config{
		HTTPAddr: orDefault(getenv("HTTP_ADDR"), ":8080"),
		DatabaseURL: strings.TrimSpace(getenv("DATABASE_URL")),
	}

	var errs []error
	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if lvl := strings.TrimSpace(getenv("LOG_LEVEL")); lvl != "" {
		if err := c.LogLevel.UnmarshalText([]byte(lvl)); err != nil {
			errs = append(errs,	fmt.Errorf("LOG_LEVEL: %w", err))
		}
	}
}

func orDefault(v, def string) string {
	if v = strings.TrimSpace(v); v != "" {
		return v
	}
	return def
}