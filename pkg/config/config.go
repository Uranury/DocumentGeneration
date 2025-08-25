package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr        string `env:"LISTEN_ADDR" default:":8080"`
	PandocPath        string `env:"PANDOC_PATH" required:"true"`
	PDFConverterURL   string `env:"PDF_CONVERTER_URL" default:"http://localhost:3100"`
	StaticToken       string `env:"STATIC_TOKEN" default:"default_token"`
	ServiceContextURL string `env:"SERVICE_CONTEXT_URL" default:"/document-generator"`
	TemplateDir       string `env:"TEMPLATE_DIR" default:"./templates"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
