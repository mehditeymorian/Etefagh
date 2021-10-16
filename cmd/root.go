package main

import (
	"github.com/mehditeymorian/etefagh/internal/cmd/server"
	"github.com/mehditeymorian/etefagh/internal/config"
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
	log "github.com/mehditeymorian/etefagh/internal/logger"
)

func main() {
	var cfg config.Config = config.Config{
		Api: handler.Config{
			Port: "3000",
		},
		DB: db.Config{
			Uri:  "mongodb://localhost:27017",
			Name: "events",
		},
		Logger: log.Config{
			Level: "debug",
		},
	}

	logger := log.New(cfg.Logger)
	server.Main(cfg, logger)
}
