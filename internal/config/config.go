package config

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	SourceType string
	FileName   string
	CodePage   string
	ApiUrl     string
	ApiKey     string
}

func ParseCommandLineArgs() (*Config, error) {
	cfg := &Config{CodePage: "CP1251"}

	flag.StringVar(&cfg.SourceType, "type", "file", "Data source type: 'file' or 'api'")
	flag.StringVar(&cfg.FileName, "filename", "", "Path to data file")
	flag.StringVar(&cfg.CodePage, "codepage", "CP1251", "File encoding (CP1251 | UTF-8)")
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

	cfg.CodePage = strings.ToUpper(cfg.CodePage)
	if cfg.CodePage != "CP1251" && cfg.CodePage != "UTF-8" {
		return cfg, fmt.Errorf("invalid code page: %s", cfg.CodePage)
	}

	return cfg, nil
}
