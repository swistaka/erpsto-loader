package config

import (
	"flag"
	"fmt"
)

type Config struct {
	SourceType string
	FileName   string
	ApiUrl     string
	ApiKey     string
}

func ParseCommandLineArgs() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.SourceType, "type", "file", "Data source type: 'file' or 'api'")
	flag.StringVar(&cfg.FileName, "filename", "", "Path to data file")
	flag.StringVar(&cfg.ApiUrl, "url", "", "API endpoint URL")
	flag.StringVar(&cfg.ApiKey, "token", "", "API authentication key")
	flag.Parse()

	if cfg.SourceType != "file" && cfg.SourceType != "api" {
		return cfg, fmt.Errorf("invalid source type: %s", cfg.SourceType)
	}

	if cfg.SourceType == "file" && cfg.FileName == "" {
		return cfg, fmt.Errorf("file name is required for file mode")
	}

	if cfg.SourceType == "api" && (cfg.ApiUrl == "" || cfg.ApiKey == "") {
		return cfg, fmt.Errorf("API mode requires URL and token")
	}

	return cfg, nil
}
