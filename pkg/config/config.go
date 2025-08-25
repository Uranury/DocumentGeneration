package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr        string `envconfig:"LISTEN_ADDR" default:":8080"`
	PandocPath        string `envconfig:"PANDOC_PATH" required:"true"`
	PDFConverterURL   string `envconfig:"PDF_CONVERTER_URL" default:"http://localhost:3000"`
	StaticToken       string `envconfig:"STATIC_TOKEN" default:"default_token"`
	ServiceContextURL string `envconfig:"SERVICE_CONTEXT_URL" default:"/document-generator"`
	TemplateDir       string `envconfig:"TEMPLATE_DIR" default:"./templates"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
