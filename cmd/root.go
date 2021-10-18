package main

import (
	"github.com/mehditeymorian/etefagh/internal/cmd/server"
	"github.com/mehditeymorian/etefagh/internal/config"
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
	log "github.com/mehditeymorian/etefagh/internal/logger"
	"github.com/mehditeymorian/etefagh/internal/redis"
	"github.com/mehditeymorian/etefagh/internal/stan"
	"github.com/mehditeymorian/etefagh/internal/telemetry"
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
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: true,
				Url:     "http://localhost:14268/api/traces",
			},
		},
		Nats: stan.Config{
			Url:         "http://localhost:4222",
			ClusterName: "test-cluster",
			ClientId:    "client1",
		},

		Redis: redis.Config{
			Address:  "localhost:6379",
			Password: "",
			DB:       0,
		},
	}

	logger := log.New(cfg.Logger)

	tracer := telemetry.New(cfg.Telemetry)
	server.Main(cfg, logger, tracer)
}
