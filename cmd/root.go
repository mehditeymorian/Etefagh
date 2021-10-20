package main

import (
	"fmt"
	"github.com/mehditeymorian/etefagh/internal/cmd/server"
	"github.com/mehditeymorian/etefagh/internal/config"
	log "github.com/mehditeymorian/etefagh/internal/logger"
	"github.com/mehditeymorian/etefagh/internal/telemetry"
	log2 "log"
)

func main() {
	readConfig, err := config.ReadConfig()
	if err != nil {
		log2.Fatalf("failed to read config: %w", err)
	}

	cfg := *readConfig

	fmt.Printf("%+v", cfg)

	logger := log.New(cfg.Logger)

	tracer := telemetry.New(cfg.Telemetry)
	server.Main(cfg, logger, tracer)
}
