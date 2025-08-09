package config

import (
	"os"
)

type Config struct {
	HTTPListen string
	LogLevel   string
	// UDS path for MediaCore gRPC when enabled
	MediaCoreUDS string
}

func Load() Config {
	cfg := Config{
		HTTPListen:  ":8080",
		LogLevel:    "info",
		MediaCoreUDS: "/var/run/mediacore.sock",
	}
	if v := os.Getenv("HTTP_LISTEN"); v != "" {
		cfg.HTTPListen = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("MEDIA_CORE_UDS"); v != "" {
		cfg.MediaCoreUDS = v
	}
	return cfg
}