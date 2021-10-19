package main

import (
	"fmt"
	"github.com/mehditeymorian/etefagh/internal/cmd/server"
	"github.com/mehditeymorian/etefagh/internal/config"
	log "github.com/mehditeymorian/etefagh/internal/logger"
	"github.com/mehditeymorian/etefagh/internal/telemetry"
)

func main() {
	cfg := config.Default()

	fmt.Printf("%+v", cfg)

	logger := log.New(cfg.Logger)

	tracer := telemetry.New(cfg.Telemetry)
	server.Main(cfg, logger, tracer)
}
