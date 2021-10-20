package config

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
	"github.com/mehditeymorian/etefagh/internal/logger"
	"github.com/mehditeymorian/etefagh/internal/redis"
	"github.com/mehditeymorian/etefagh/internal/stan"
	"github.com/mehditeymorian/etefagh/internal/telemetry"
	"os"
)

// Config Struct of all the configuration in the app
type Config struct {
	Api       handler.Config   `yaml:"api"`
	DB        db.Config        `yaml:"db"`
	Logger    logger.Config    `yaml:"logger"`
	Telemetry telemetry.Config `yaml:"telemetry"`
	Nats      stan.Config      `yaml:"nats"`
	Redis     redis.Config     `yaml:"redis"`
}

func ReadConfig() (*Config, error) {
	var cfg Config

	path, _ := os.Getwd()
	bytes, err := os.ReadFile(path + "/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
