package config

import (
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
	log "github.com/mehditeymorian/etefagh/internal/logger"
	"github.com/mehditeymorian/etefagh/internal/redis"
	"github.com/mehditeymorian/etefagh/internal/stan"
	"github.com/mehditeymorian/etefagh/internal/telemetry"
)

func Default() Config {
	return Config{
		Api: handler.Config{
			Port: "3000",
		},
		DB: db.Config{
			Uri:  "mongodb://db:27017",
			Name: "events",
		},
		Logger: log.Config{
			Level: "debug",
		},
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: true,
				Url:     "http://jaeger:14268/api/traces",
			},
		},
		Nats: stan.Config{
			Url:         "http://nats:4222",
			ClusterName: "test-cluster",
			ClientId:    "client1",
		},

		Redis: redis.Config{
			Address:  "redis:6379",
			Password: "",
			DB:       0,
		},
	}
}
