package config

import (
	"github.com/mehditeymorian/etefagh/internal/db"
	"github.com/mehditeymorian/etefagh/internal/handler"
	"github.com/mehditeymorian/etefagh/internal/logger"
	"github.com/mehditeymorian/etefagh/internal/redis"
	"github.com/mehditeymorian/etefagh/internal/stan"
	"github.com/mehditeymorian/etefagh/internal/telemetry"
)

// Config Struct of all the configuration in the app
type Config struct {
	Api       handler.Config
	DB        db.Config
	Logger    logger.Config
	Telemetry telemetry.Config
	Nats      stan.Config
	Redis     redis.Config
}
